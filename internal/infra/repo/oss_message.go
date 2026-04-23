package repo_infra

import (
	"context"
	"errors"
	"time"

	"poprako-s/internal/domain/model"
	iface "poprako-s/internal/domain/repo"
	entity "poprako-s/internal/infra/repo/entity"

	"gorm.io/gorm"
)

type ossMessageRepoImpl struct {
	gdb *gorm.DB
}

// NewOSSMessageRepo 返回 OSSMessageRepo 的 GORM 实现
func NewOSSMessageRepo(gdb *gorm.DB) iface.OSSMessageRepo {
	return &ossMessageRepoImpl{gdb: gdb}
}

func newOSSMessageRepoFromCx(cx context.Context) (iface.OSSMessageRepo, error) {
	gdb, ok := cx.Value(txnKey).(*gorm.DB)
	if !ok {
		return nil, errors.New("[newOSSMessageRepoFromCx]: 无法从上下文中获取事务数据库连接")
	}

	return &ossMessageRepoImpl{gdb: gdb}, nil
}

func (r *ossMessageRepoImpl) FromTxnCx(cx context.Context) (iface.OSSMessageRepo, error) {
	return newOSSMessageRepoFromCx(cx)
}

// UpsertCreatePending 对同一逻辑资源的 create_pending 消息 upsert
func (r *ossMessageRepoImpl) UpsertCreatePending(msg *model.OSSMessage) error {
	now := time.Now()

	vals := map[string]any{
		"id":            msg.ID,
		"resource_type": string(msg.ResourceType),
		"resource_id":   msg.ResourceID,
		"operation":     string(model.OSSOperationCreatePending),
		"status":        string(model.OSSMessageStatusPending),
		"object_key":    msg.ObjectKey,
		"payload_json":  msg.PayloadJSON,
		"visible_at":    msg.VisibleAt,
		"expire_at":     msg.ExpireAt,
		"attempt_count": 0,
		"last_error":    "",
		"created_at":    now,
		"updated_at":    now,
	}

	// 先尝试按 (resource_type, resource_id, operation) 更新已有的 pending 消息
	result := r.gdb.Table(entity.OSSMessageTable).
		Where(
			"resource_type = ? AND resource_id = ? AND operation = ? AND status != ?",
			string(msg.ResourceType),
			msg.ResourceID,
			string(model.OSSOperationCreatePending),
			string(model.OSSMessageStatusCompleted),
		).
		Updates(map[string]any{
			"id":            msg.ID,
			"object_key":    msg.ObjectKey,
			"status":        string(model.OSSMessageStatusPending),
			"attempt_count": 0,
			"last_error":    "",
			"expire_at":     msg.ExpireAt,
			"visible_at":    msg.VisibleAt,
			"updated_at":    now,
		})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected > 0 {
		return nil
	}

	// 否则插入新记录
	return r.gdb.Table(entity.OSSMessageTable).Create(vals).Error
}

// InsertDeletePending 插入一条 delete_pending 消息
func (r *ossMessageRepoImpl) InsertDeletePending(msg *model.OSSMessage) error {
	now := time.Now()

	vals := map[string]any{
		"id":            msg.ID,
		"resource_type": string(msg.ResourceType),
		"resource_id":   msg.ResourceID,
		"operation":     string(model.OSSOperationDeletePending),
		"status":        string(model.OSSMessageStatusPending),
		"object_key":    msg.ObjectKey,
		"payload_json":  msg.PayloadJSON,
		"visible_at":    msg.VisibleAt,
		"expire_at":     nil,
		"attempt_count": 0,
		"last_error":    "",
		"created_at":    now,
		"updated_at":    now,
	}

	return r.gdb.Table(entity.OSSMessageTable).Create(vals).Error
}

// MarkCompleted 将指定消息标记为 completed
func (r *ossMessageRepoImpl) MarkCompleted(id string) error {
	return r.gdb.Table(entity.OSSMessageTable).
		Where("id = ?", id).
		Updates(map[string]any{
			"status":     string(model.OSSMessageStatusCompleted),
			"updated_at": time.Now(),
		}).Error
}

// CompleteCreatePendingByResource 将指定资源的 create_pending 消息标记为 completed
func (r *ossMessageRepoImpl) CompleteCreatePendingByResource(
	resourceType model.OSSResourceType,
	resourceID string,
) error {
	return r.gdb.Table(entity.OSSMessageTable).
		Where(
			"resource_type = ? AND resource_id = ? AND operation = ? AND status != ?",
			string(resourceType),
			resourceID,
			string(model.OSSOperationCreatePending),
			string(model.OSSMessageStatusCompleted),
		).
		Updates(map[string]any{
			"status":     string(model.OSSMessageStatusCompleted),
			"updated_at": time.Now(),
		}).Error
}

// ClaimPending 用 CAS 抢占一条 pending 消息并返回
func (r *ossMessageRepoImpl) ClaimPending(operation model.OSSOperation) (*model.OSSMessage, error) {
	now := time.Now()

	var row entity.OSSMessageRow

	// 空队列是正常情况，这里避免用 First 返回 ErrRecordNotFound 触发错误日志。
	result := r.gdb.Table(entity.OSSMessageTable).
		Where(
			"operation = ? AND status = ? AND visible_at <= ?",
			string(operation),
			string(model.OSSMessageStatusPending),
			now,
		).
		Order("visible_at ASC").
		Limit(1).
		Find(&row)
	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, nil
	}

	// CAS：将 pending -> processing
	result = r.gdb.Table(entity.OSSMessageTable).
		Where("id = ? AND status = ?", row.ID, string(model.OSSMessageStatusPending)).
		Updates(map[string]any{
			"status":        string(model.OSSMessageStatusProcessing),
			"processing_at": now,
			"updated_at":    now,
		})
	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		// 其他 worker 抢先，视为未抢到
		return nil, nil
	}

	row.Status = string(model.OSSMessageStatusProcessing)
	row.ProcessingAt = &now

	msg := entity.ToOSSMessage(row)

	return &msg, nil
}

// ResetToRetry 将消息回退为 pending 并记录错误
func (r *ossMessageRepoImpl) ResetToRetry(
	id string,
	lastError string,
	nextVisibleAt interface{},
) error {
	return r.gdb.Table(entity.OSSMessageTable).
		Where("id = ?", id).
		Updates(map[string]any{
			"status":        string(model.OSSMessageStatusPending),
			"last_error":    lastError,
			"visible_at":    nextVisibleAt,
			"attempt_count": gorm.Expr("attempt_count + 1"),
			"updated_at":    time.Now(),
		}).Error
}

// ResetStuckProcessing 将超时卡住的 processing 消息回退为 pending
func (r *ossMessageRepoImpl) ResetStuckProcessing(before time.Time) error {
	return r.gdb.Table(entity.OSSMessageTable).
		Where(
			"status = ? AND processing_at IS NOT NULL AND processing_at <= ?",
			string(model.OSSMessageStatusProcessing),
			before,
		).
		Updates(map[string]any{
			"status":        string(model.OSSMessageStatusPending),
			"processing_at": nil,
			"visible_at":    time.Now(),
			"updated_at":    time.Now(),
		}).Error
}

// DeleteCompletedBefore 删除指定时间之前已经 completed 的消息
func (r *ossMessageRepoImpl) DeleteCompletedBefore(before time.Time) error {
	return r.gdb.Table(entity.OSSMessageTable).
		Where("status = ? AND updated_at <= ?", string(model.OSSMessageStatusCompleted), before).
		Delete(&entity.OSSMessageRow{}).Error
}
