package entity

import (
	"time"

	"labelplus-next-web-be/internal/domain/model"
)

type TeamRow struct {
	ID          string     `gorm:"column:id"`
	Name        string     `gorm:"column:name"`
	Description string     `gorm:"column:description"`
	CreatedAt   time.Time  `gorm:"column:created_at"`
	UpdatedAt   time.Time  `gorm:"column:updated_at"`
	DeletedAt   *time.Time `gorm:"column:deleted_at"`
}

func (TeamRow) TableName() string { return "team_table" }

func ToTeamInfo(row TeamRow) *model.TeamInfo {
	return &model.TeamInfo{
		ID:          row.ID,
		Name:        row.Name,
		Description: row.Description,
		CreatedAt:   row.CreatedAt,
		UpdatedAt:   row.UpdatedAt,
	}
}
