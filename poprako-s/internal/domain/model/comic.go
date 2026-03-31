package model

import "time"

// ComicInfo 包含作品集/漫画的元信息，用于聚合其章节和统计数据
type ComicInfo struct {
	ID string

	WorksetID string
	// Workset 仅在 includes 指定时填充
	Workset *WorksetInfo

	// 在作品集内部的序号
	Index        int
	Title        string
	Author       string
	Description  string
	ChapterCount int

	CreatorID string
	// Creator 仅在 includes 指定时填充
	Creator *UserInfo

	LastActiveAt time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
}
