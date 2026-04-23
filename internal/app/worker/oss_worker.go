package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"poprako-s/internal/domain/ext/oss"
	"poprako-s/internal/domain/model"
	"poprako-s/internal/domain/repo"

	"go.uber.org/zap"
)

const (
	// workerScanInterval worker 扫描消息表的间隔
	workerScanInterval = 5 * time.Minute
	// workerStuckThreshold processing 超过 3 个扫描间隔视为卡死
	workerStuckThreshold = 3 * workerScanInterval
	// workerCompletedPurgeInterval completed 消息清理间隔
	workerCompletedPurgeInterval = 24 * time.Hour
	// workerRetryBackoff 失败后下次重试的等待时间
	workerRetryBackoff = 5 * time.Minute
)

// OSSWorker 后台消费 oss_message_table，负责执行远程 OSS 删除与超时清理
type OSSWorker struct {
	msgRepo repo.OSSMessageRepo
	deleter oss.Deleter
}

// NewOSSWorker 返回一个 OSSWorker 实例
func NewOSSWorker(
	msgRepo repo.OSSMessageRepo,
	deleter oss.Deleter,
) *OSSWorker {
	if msgRepo == nil || deleter == nil {
		zap.L().Panic(
			"NewOSSWorker: 依赖项不能为空",
			zap.Bool("msgRepo_nil", msgRepo == nil),
			zap.Bool("deleter_nil", deleter == nil),
		)
	}

	return &OSSWorker{
		msgRepo: msgRepo,
		deleter: deleter,
	}
}

// Start 启动后台 worker，阻塞直到 ctx 取消
func (w *OSSWorker) Start(ctx context.Context) {
	zap.L().Info("[OSSWorker] 后台 worker 启动")

	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		defer wg.Done()
		w.consumeLoop(ctx)
	}()

	go func() {
		defer wg.Done()
		w.reclaimLoop(ctx)
	}()

	go func() {
		defer wg.Done()
		w.purgeLoop(ctx)
	}()

	<-ctx.Done()
	wg.Wait()
	zap.L().Info("[OSSWorker] 后台 worker 停止")
}

func (w *OSSWorker) consumeLoop(ctx context.Context) {
	ticker := time.NewTicker(workerScanInterval)
	defer ticker.Stop()

	w.consumeAll(ctx, model.OSSOperationDeletePending)
	w.consumeAll(ctx, model.OSSOperationCreatePending)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			w.consumeAll(ctx, model.OSSOperationDeletePending)
			w.consumeAll(ctx, model.OSSOperationCreatePending)
		}
	}
}

func (w *OSSWorker) reclaimLoop(ctx context.Context) {
	ticker := time.NewTicker(workerScanInterval)
	defer ticker.Stop()

	w.reclaimStuckProcessing()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			w.reclaimStuckProcessing()
		}
	}
}

func (w *OSSWorker) purgeLoop(ctx context.Context) {
	ticker := time.NewTicker(workerCompletedPurgeInterval)
	defer ticker.Stop()

	w.purgeCompletedMessages()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			w.purgeCompletedMessages()
		}
	}
}

// consumeAll 循环抢占并消费指定类型的 pending 消息，直到没有可消费消息为止
func (w *OSSWorker) consumeAll(ctx context.Context, op model.OSSOperation) {
	for {
		if ctx.Err() != nil {
			return
		}

		msg, err := w.msgRepo.ClaimPending(op)
		if err != nil {
			zap.L().Error(
				"[OSSWorker] 抢占消息失败",
				zap.String("operation", string(op)),
				zap.Error(err),
			)

			return
		}

		if msg == nil {
			return
		}

		w.processMessage(msg)
	}
}

// processMessage 处理单条消息
func (w *OSSWorker) processMessage(msg *model.OSSMessage) {
	lgr := zap.L().With(
		zap.String("msg_id", msg.ID),
		zap.String("resource_type", string(msg.ResourceType)),
		zap.String("resource_id", msg.ResourceID),
		zap.String("operation", string(msg.Operation)),
		zap.Int("attempt_count", msg.AttemptCount),
	)

	var processErr error

	switch msg.Operation {
	case model.OSSOperationDeletePending:
		processErr = w.processDelete(msg)
	case model.OSSOperationCreatePending:
		processErr = w.processCreateTimeout(msg)
	default:
		lgr.Warn("[OSSWorker] 未知操作类型，跳过", zap.String("operation", string(msg.Operation)))
		processErr = fmt.Errorf("未知操作类型: %s", msg.Operation)
	}

	if processErr != nil {
		lgr.Warn(
			"[OSSWorker] 消息处理失败，回退为 pending",
			zap.Error(processErr),
		)

		nextVisible := time.Now().Add(workerRetryBackoff)
		if resetErr := w.msgRepo.ResetToRetry(msg.ID, processErr.Error(), nextVisible); resetErr != nil {
			lgr.Error("[OSSWorker] 回退消息状态失败", zap.Error(resetErr))
		}

		return
	}

	if err := w.msgRepo.MarkCompleted(msg.ID); err != nil {
		lgr.Error("[OSSWorker] 标记消息完成失败", zap.Error(err))
	} else {
		lgr.Info("[OSSWorker] 消息处理完成")
	}
}

func (w *OSSWorker) reclaimStuckProcessing() {
	cutoff := time.Now().Add(-workerStuckThreshold)
	if err := w.msgRepo.ResetStuckProcessing(cutoff); err != nil {
		zap.L().Error("[OSSWorker] 回收卡死消息失败", zap.Error(err))
	}
}

func (w *OSSWorker) purgeCompletedMessages() {
	cutoff := time.Now().Add(-workerCompletedPurgeInterval)
	if err := w.msgRepo.DeleteCompletedBefore(cutoff); err != nil {
		zap.L().Error("[OSSWorker] 清理已完成消息失败", zap.Error(err))
	}
}

// processDelete 执行 delete_pending 消息：单删或批量删除远程 OSS 对象
func (w *OSSWorker) processDelete(msg *model.OSSMessage) error {
	// 优先使用 payload_json 中的批量 keys
	if msg.PayloadJSON != "" {
		var payload struct {
			ObjectKeys []string `json:"object_keys"`
		}

		if err := json.Unmarshal([]byte(msg.PayloadJSON), &payload); err != nil {
			return fmt.Errorf("解析 payload_json 失败: %w", err)
		}

		if len(payload.ObjectKeys) > 0 {
			return w.deleter.DeleteBatch(payload.ObjectKeys)
		}
	}

	// 回退到单 key 删除
	if msg.ObjectKey != "" {
		return w.deleter.Delete(msg.ObjectKey)
	}

	// 没有任何 key，视为成功
	return nil
}

// processCreateTimeout 执行 create_pending 超时清理：删除未被 confirm 的上传对象
func (w *OSSWorker) processCreateTimeout(msg *model.OSSMessage) error {
	// 仅当 expire_at 已过期才执行清理
	if msg.ExpireAt == nil || time.Now().Before(*msg.ExpireAt) {
		// 未过期，不处理
		return nil
	}

	if msg.ObjectKey == "" {
		return nil
	}

	return w.deleter.Delete(msg.ObjectKey)
}
