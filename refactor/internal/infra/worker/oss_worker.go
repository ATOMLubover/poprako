package worker_infra

import (
	"context"
	"time"

	oss_iface "poprako-s/internal/domain/ext/oss"
	"poprako-s/internal/domain/model/aggr"
	"poprako-s/internal/domain/model/enum"
	repo_iface "poprako-s/internal/domain/repo"

	"go.uber.org/zap"
)

// Worker timing constants used by `OssWorker` loops.
const (
	// `ossConsumeInterval` is the interval of pending-message consumption.
	ossConsumeInterval = 5 * time.Minute

	// `ossReclaimInterval` is the interval of stuck-message reclaiming.
	ossReclaimInterval = 5 * time.Minute

	// `ossPurgeInterval` is the interval of completed-message cleanup.
	ossPurgeInterval = 24 * time.Hour

	// `ossProcTimeout` is the max processing duration before a message is reclaimed.
	ossProcTimeout = 15 * time.Minute

	// `ossRetryDelay` is the retry delay after one processing failure.
	ossRetryDelay = 5 * time.Minute
)

// `OssWorker` defines the runtime contract of OSS background worker.
type OssWorker interface {
	// `Run` starts worker loops and blocks until `cx` is cancelled.
	Run(cx context.Context)
}

// `ossWorkerImpl` is the default implementation of `OssWorker`.
type ossWorkerImpl struct {
	msgRepo repo_iface.OssMsgRepo
	cleaner oss_iface.Cleaner
}

// `NewOssWorker` creates one ready-to-run `OssWorker` implementation.
func NewOssWorker(msgRepo repo_iface.OssMsgRepo, cleaner oss_iface.Cleaner) OssWorker {
	if msgRepo == nil || cleaner == nil {
		zap.L().Panic(
			"[NewOssWorker] nil dependency",
			zap.Bool("msgRepo", msgRepo == nil),
			zap.Bool("cleaner", cleaner == nil),
		)
	}

	return &ossWorkerImpl{
		msgRepo: msgRepo,
		cleaner: cleaner,
	}
}

// `Run` starts all loops and waits for cancellation.
func (w *ossWorkerImpl) Run(cx context.Context) {
	// Start periodic consume loop.
	go w.consumeLoop(cx)

	// Start periodic stuck-message reclaim loop.
	go w.reclaimLoop(cx)

	// Start periodic completed-message purge loop.
	go w.purgeLoop(cx)

	// Block until caller requests shutdown.
	<-cx.Done()
}

// `consumeLoop` periodically consumes pending delete and create messages.
func (w *ossWorkerImpl) consumeLoop(cx context.Context) {
	ticker := time.NewTicker(ossConsumeInterval)

	// Run one immediate pass at startup.
	w.consumeOnce()

	for {
		select {
		case <-cx.Done():
			ticker.Stop()

			return

		case <-ticker.C:
			w.consumeOnce()
		}
	}
}

// `consumeOnce` drains both delete and create pending queues.
func (w *ossWorkerImpl) consumeOnce() {
	// Consume delete queue first.
	w.consumeByOp(enum.OssOpDel)

	// Consume create queue second.
	w.consumeByOp(enum.OssOpCre)
}

// `consumeByOp` claims and handles one operation queue until empty.
func (w *ossWorkerImpl) consumeByOp(op enum.OssOp) {
	for {
		msg, err := w.msgRepo.ClaimPending(op)
		if err != nil {
			zap.L().Error(
				"[ossWorkerImpl.consumeByOp] failed to claim pending message",
				zap.String("op", string(op)),
				zap.Error(err),
			)

			return
		}

		if msg == nil {
			return
		}

		w.handleMsg(msg)
	}
}

// `handleMsg` dispatches one claimed message by operation.
func (w *ossWorkerImpl) handleMsg(msg *aggr.OssMsg) {
	if msg.Op == enum.OssOpDel {
		w.handleDelMsg(msg)

		return
	}

	if msg.Op == enum.OssOpCre {
		w.handleCreMsg(msg)

		return
	}

	w.markPending(msg.Id, "unknown operation", time.Now().Add(ossRetryDelay))
}

// `handleDelMsg` handles one pending delete message.
func (w *ossWorkerImpl) handleDelMsg(msg *aggr.OssMsg) {
	// Empty key set already satisfies deletion expectation.
	if len(msg.ObjKeys) == 0 {
		w.markCompleted(msg.Id)

		return
	}

	// Execute remote batch deletion.
	if err := w.cleaner.DelBatch(msg.ObjKeys); err != nil {
		w.markPending(msg.Id, err.Error(), time.Now().Add(ossRetryDelay))

		return
	}

	// Mark done after successful deletion.
	w.markCompleted(msg.Id)
}

// `handleCreMsg` handles one pending create message.
func (w *ossWorkerImpl) handleCreMsg(msg *aggr.OssMsg) {
	now := time.Now()

	// Requeue not-yet-expired create message and avoid premature deletion.
	if now.Before(msg.ExpireAt) {
		w.markPending(msg.Id, "not yet expired", msg.ExpireAt)

		return
	}

	// Empty key set already satisfies deletion expectation.
	if len(msg.ObjKeys) == 0 {
		w.markCompleted(msg.Id)

		return
	}

	// Delete expired unconfirmed objects.
	if err := w.cleaner.DelBatch(msg.ObjKeys); err != nil {
		w.markPending(msg.Id, err.Error(), time.Now().Add(ossRetryDelay))

		return
	}

	// Mark done after successful deletion.
	w.markCompleted(msg.Id)
}

// `markCompleted` marks one queue message as completed.
func (w *ossWorkerImpl) markCompleted(msgId string) {
	if err := w.msgRepo.MarkCompleted(msgId); err != nil {
		zap.L().Error(
			"[ossWorkerImpl.markCompleted] failed to mark message completed",
			zap.String("msg_id", msgId),
			zap.Error(err),
		)
	}
}

// `markPending` resets one queue message to pending with retry metadata.
func (w *ossWorkerImpl) markPending(msgId string, errMsg string, nextVisibleAt time.Time) {
	if err := w.msgRepo.MarkPending(msgId, errMsg, nextVisibleAt); err != nil {
		zap.L().Error(
			"[ossWorkerImpl.markPending] failed to mark message pending",
			zap.String("msg_id", msgId),
			zap.String("err_msg", errMsg),
			zap.Time("next_visible_at", nextVisibleAt),
			zap.Error(err),
		)
	}
}

// `reclaimLoop` periodically resets stuck processing messages.
func (w *ossWorkerImpl) reclaimLoop(cx context.Context) {
	ticker := time.NewTicker(ossReclaimInterval)

	// Run one immediate pass at startup.
	w.reclaimOnce()

	for {
		select {
		case <-cx.Done():
			ticker.Stop()

			return

		case <-ticker.C:
			w.reclaimOnce()
		}
	}
}

// `reclaimOnce` resets processing messages stuck beyond timeout.
func (w *ossWorkerImpl) reclaimOnce() {
	cutoff := time.Now().Add(-ossProcTimeout)

	if err := w.msgRepo.ResetStuck(cutoff); err != nil {
		zap.L().Error(
			"[ossWorkerImpl.reclaimOnce] failed to reset stuck messages",
			zap.Error(err),
		)
	}
}

// `purgeLoop` periodically purges old completed messages.
func (w *ossWorkerImpl) purgeLoop(cx context.Context) {
	ticker := time.NewTicker(ossPurgeInterval)

	// Run one immediate pass at startup.
	w.purgeOnce()

	for {
		select {
		case <-cx.Done():
			ticker.Stop()

			return

		case <-ticker.C:
			w.purgeOnce()
		}
	}
}

// `purgeOnce` deletes completed messages older than retention window.
func (w *ossWorkerImpl) purgeOnce() {
	cutoff := time.Now().Add(-ossPurgeInterval)

	if err := w.msgRepo.CleanCompleted(cutoff); err != nil {
		zap.L().Error(
			"[ossWorkerImpl.purgeOnce] failed to purge completed messages",
			zap.Error(err),
		)
	}
}
