package model

import "time"

type TeamInfo struct {
	ID string

	Name             string
	Description      string
	AvatarOSSKey     string
	IsAvatarUploaded bool

	CreatedAt time.Time
	UpdatedAt time.Time
}