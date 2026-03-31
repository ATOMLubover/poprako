package model

import "time"

type WorksetInfo struct {
	ID string

	TeamID string
	// Team 仅在 includes 指定时填充
	Team *TeamInfo

	// 在汉化组内部的序号
	Index int

	Name        string
	Description string
	ComicCount  int

	CreatedAt time.Time
	UpdatedAt time.Time
}
