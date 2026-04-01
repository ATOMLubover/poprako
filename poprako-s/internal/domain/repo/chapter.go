package repo

import "poprako-s/internal/domain/model"

// ChapterRepo 定义章节数据持久化操作的抽象接口
type ChapterRepo interface {
	// GetByID 根据章节 ID 获取章节信息；若不存在返回 error
	GetByID(id string) (*model.ChapterInfo, error)
	// FindPinnedByComicID 获取漫画下被置顶的章节信息
	FindPinnedByComicID(comicID string) (*model.ChapterInfo, error)
	// List 根据筛选条件返回章节信息列表
	List(opt model.ChapterQueryOpt) ([]model.ChapterInfo, error)
	// Count 根据筛选条件返回章节数量
	Count(opt model.ChapterQueryOpt) (int64, error)

	// Create 持久化一个新的章节
	Create(c *model.ChapterCreation) (*model.ChapterInfo, error)
	// UpdateStats 仅更新章节的统计字段（unit 数量等）
	UpdateStats(stats *model.ChapterStats) error
	// Delete 删除章节（硬删除）
	Delete(id string) error
}
