package repo

import (
	"context"

	"poprako-s/internal/domain/model"
)

// ComicRepo 是漫画仓库的接口
type ComicRepo interface {
	// GetByID 根据漫画 ID 获取漫画信息；若不存在返回 error
	GetByID(id string) (*model.ComicInfo, error)
	// List 根据筛选条件返回漫画信息列表
	List(opt model.ComicQueryOpt) ([]model.ComicInfo, error)
	// Count 根据筛选条件返回漫画数量
	Count(opt model.ComicQueryOpt) (int64, error)

	// Create 持久化一个新的漫画
	Create(c *model.ComicCreation) (*model.ComicInfo, error)
	// Update 更新漫画信息（PUT 语义）
	Update(u *model.ComicUpdate) error
	// UpdateChapterCount 按 delta 更新漫画下的章节数量
	UpdateChapterCount(id string, delta int) error
	// PreFillCoverOSSKey 预填充封面对象 Key，并将上传状态重置为未上传
	PreFillCoverOSSKey(id string, coverOSSKey string) error
	// ConfirmCoverUploaded 将漫画封面状态标记为已上传
	ConfirmCoverUploaded(id string) error
	// Delete 删除漫画（硬删除）
	Delete(id string) error

	// FromTxnCx 从上下文中获取事务，并分离出一个带事务的 ComicRepo 实例
	FromTxnCx(cx context.Context) (ComicRepo, error)
}
