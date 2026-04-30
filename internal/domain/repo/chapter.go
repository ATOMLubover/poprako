package repo

import (
	"context"

	"poprako-s/internal/domain/model"
)

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
	// Update 更新章节信息（PUT 语义）
	Update(u *model.ChapterUpdate) error
	// Delete 删除章节（硬删除）
	Remove(id string) error

	// LockByID 在当前事务中锁定章节行，用于串行化同一章节的并发修改
	LockByID(id string) error
	// UpdateStats 使用原子增量更新章节的 unit 统计字段
	UpdateStats(id string, totalDelta, translatedDelta, proofreadDelta int) error
	// UpdatePageCount 更新章节的页面数量统计
	UpdatePageCount(id string, delta int) error

	// FromTxnCx 从上下文中获取事务，并分离出一个带事务的 ChapterRepo 实例
	FromTxnCx(cx context.Context) (ChapterRepo, error)
}
