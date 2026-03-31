package repo

import "poprako-s/internal/domain/model"

type ChapterRepo interface {
	// Create 创建一个新的章节，返回创建后的章节信息
	Create(c *model.ChapterCreation) (*model.ChapterInfo, error)
}
