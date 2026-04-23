package repo

import (
	"context"
	"time"

	"poprako-s/internal/domain/model"
)

// OSSMessageRepo 定义本地 OSS 消息表的仓储接口
type OSSMessageRepo interface {
	// UpsertCreatePending 对同一逻辑资源的 create_pending 消息进行 upsert：
	// 若已存在则重置状态与字段，若不存在则插入。
	UpsertCreatePending(msg *model.OSSMessage) error

	// InsertDeletePending 插入一条 delete_pending 消息（快照语义，不需要 upsert）
	InsertDeletePending(msg *model.OSSMessage) error

	// MarkCompleted 将消息状态由任意状态改为 completed
	MarkCompleted(id string) error

	// CompleteCreatePendingByResource 将指定资源的 create_pending 消息标记为 completed
	// 用于 confirm 场景
	CompleteCreatePendingByResource(resourceType model.OSSResourceType, resourceID string) error

	// ClaimPending 用 CAS 将一条满足条件的 pending 消息改为 processing 并返回
	// 若没有可消费的消息则返回 nil, nil
	ClaimPending(operation model.OSSOperation) (*model.OSSMessage, error)

	// ResetToRetry 将消息从 processing 回退为 pending 并记录错误信息
	ResetToRetry(id string, lastError string, nextVisibleAt interface{}) error

	// ResetStuckProcessing 将超时卡住的 processing 消息回退为 pending
	ResetStuckProcessing(before time.Time) error

	// DeleteCompletedBefore 删除指定时间之前已经 completed 的消息
	DeleteCompletedBefore(before time.Time) error

	// FromTxnCx 从事务上下文中提取事务 DB 并返回绑定事务的仓储实例
	FromTxnCx(cx context.Context) (OSSMessageRepo, error)
}
