package repo

import "poprako-s/internal/domain/model"

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
	// Delete 删除漫画（硬删除）
	Delete(id string) error
}
