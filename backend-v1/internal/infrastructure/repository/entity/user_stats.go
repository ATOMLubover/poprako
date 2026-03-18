package entity

import (
	"time"

	"labelplus-next-web-be/internal/domain/model"
)

const UserStatsTable = "user_stats_table"

type UserStatsRow struct {
	ID string `gorm:"column:id"`

	UserID string `gorm:"column:user_id"`

	TotalAssignmentCount    int `gorm:"column:total_assignment_count"`
	ActiveAssignmentCount   int `gorm:"column:active_assignment_count"`
	FinishedAssignmentCount int `gorm:"column:finished_assignment_count"`

	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (UserStatsRow) TableName() string { return UserStatsTable }

func ToUserStats(row UserStatsRow) model.UserStats {
	return model.NewUserStats(
		row.UserID,
		row.TotalAssignmentCount,
		row.ActiveAssignmentCount,
		row.FinishedAssignmentCount,
		row.CreatedAt,
		row.UpdatedAt,
	)
}
