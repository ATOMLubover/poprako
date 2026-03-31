package repo

import "poprako-s/internal/domain/model"

// ChapterRepo 定义章节数据持久化操作的抽象接口
type ChapterRepo interface {
	// Create 创建一个新的章节，返回创建后的章节信息
	Create(c *model.ChapterCreation) (*model.ChapterInfo, error)
}
