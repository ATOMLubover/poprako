package repo_entity

import (
	"time"

	"poprako-s/internal/domain/model"
)

const TeamTable = "team_table"

type TeamInfoRow struct {
	ID string `gorm:"column:id"`

	Name             string  `gorm:"column:name"`
	Desc             *string `gorm:"column:description"`
	AvatarOSSKey     *string `gorm:"column:avatar_oss_key"`
	IsAvatarUploaded bool    `gorm:"column:is_avatar_uploaded"`

	CreatedAt time.Time  `gorm:"column:created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at"`
}

func ToTeamInfo(row TeamInfoRow) model.TeamInfo {
	info := model.TeamInfo{
		ID:               row.ID,
		Name:             row.Name,
		IsAvatarUploaded: row.IsAvatarUploaded,
		CreatedAt:        row.CreatedAt,
		UpdatedAt:        row.UpdatedAt,
	}

	if row.Desc != nil {
		info.Desc = *row.Desc
	}
	if row.AvatarOSSKey != nil {
		info.AvatarOSSKey = *row.AvatarOSSKey
	}

	return info
}
