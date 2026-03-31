package model

import "time"

// TeamInfo 表示团队（如汉化组）的元信息
type TeamInfo struct {
	ID string

	Name             string
	Description      string
	AvatarOSSKey     string
	IsAvatarUploaded bool

	CreatedAt time.Time
	UpdatedAt time.Time
}
