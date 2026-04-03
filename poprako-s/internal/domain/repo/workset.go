package repo

import (
	"context"

	"poprako-s/internal/domain/model"
)

// WorksetRepo 是作品集仓库的接口
type WorksetRepo interface {
	// GetByID 根据作品集 ID 获取作品集信息；若不存在返回 error
	GetByID(id string) (*model.WorksetInfo, error)
	// List 根据筛选条件返回作品集信息列表
	List(opt model.WorksetQueryOpt) ([]model.WorksetInfo, error)
	// Count 根据筛选条件返回作品集数量
	Count(opt model.WorksetQueryOpt) (int64, error)

	// Create 持久化一个新的作品集
	Create(c *model.WorksetCreation) (*model.WorksetInfo, error)
	// Update 更新作品集信息
	Update(u *model.WorksetUpdate) error
	// UpdateComicCount 按 delta 更新作品集下的漫画数量
	UpdateComicCount(id string, delta int) error
	// Delete 删除作品集（硬删除）
	Delete(id string) error

	// FromTxnCx 从上下文中获取事务，并分离出一个带事务的 WorksetRepo 实例
	FromTxnCx(cx context.Context) (WorksetRepo, error)
}
