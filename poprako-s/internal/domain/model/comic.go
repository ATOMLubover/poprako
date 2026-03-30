package model

import "time"

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
