package repo

import (
	"context"

	"poprako-s/internal/domain/model"
)

// PageRepo 是页面仓库的接口
type PageRepo interface {
	// GetByID 根据页面 ID 获取页面信息；若不存在返回 error
	GetByID(id string) (*model.PageInfo, error)
	// List 根据筛选条件返回页面信息列表
	List(opt model.PageQueryOpt) ([]model.PageInfo, error)
	// LockByID 在当前事务中锁定页面行，用于串行化同一页面的并发修改
	LockByID(id string) error

	// CreateBatch 批量创建页面
	CreateBatch(pages []*model.PageCreation) error
	// Update 更新页面信息
	Update(u *model.PageUpdate) error
	// UpdateStats 使用原子增量更新页面的 unit 统计字段
	UpdateStats(id string, totalDelta, translatedDelta, proofreadDelta int) error
	// Delete 删除页面（硬删除）
	Delete(id string) error
	// DeleteBatch 批量删除页面
	DeleteBatch(ids []string) error

	// FromTxnCx 从上下文中获取事务，并分离出一个带事务的 PageRepo 实例
	FromTxnCx(cx context.Context) (PageRepo, error)
}
