package repo

import "poprako-s/internal/domain/model"

// PageRepo 是页面仓库的接口
type PageRepo interface {
	// GetByID 根据页面 ID 获取页面信息；若不存在返回 error
	GetByID(id string) (*model.PageInfo, error)
	// List 根据筛选条件返回页面信息列表
	List(opt model.PageQueryOpt) ([]model.PageInfo, error)
	// GetStatsByID 根据页面 ID 获取页面统计数据
	GetStatsByID(pageID string) (*model.PageStats, error)

	// CreateBatch 批量创建页面
	CreateBatch(pages []*model.PageCreation) error
	// Update 更新页面信息
	Update(u *model.PageUpdate) error
	// UpdateStats 仅更新页面的统计字段
	UpdateStats(stats *model.PageStats) error
	// Delete 删除页面（硬删除）
	Delete(id string) error
	// DeleteBatch 批量删除页面
	DeleteBatch(ids []string) error
}
