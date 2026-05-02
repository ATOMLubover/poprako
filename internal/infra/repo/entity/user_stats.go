package entity

import (
	"time"

	"poprako-s/internal/domain/model/aggr"
)

// `USER_STATS_TABLE` is table name for user stats.
const USER_STATS_TABLE = "t_user_stats"

// `UserStatsRow` maps one user stats record for read queries.
type UserStatsRow struct {
	Id string `gorm:"column:id;primaryKey"`

	UserId string `gorm:"column:user_id"`

	TotalAssignmentCnt    int `gorm:"column:total_assignment_count"`
	ActiveAssignmentCnt   int `gorm:"column:active_assignment_count"`
	FinishedAssignmentCnt int `gorm:"column:finished_assignment_count"`

	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

// `TableName` returns table name for `UserStatsRow`.
func (*UserStatsRow) TableName() string {
	return USER_STATS_TABLE
}

// `ToUserStatsAggr` converts row to user stats aggregate.
func (r *UserStatsRow) ToUserStatsAggr() *aggr.UserStats {
	if r == nil {
		return nil
	}

	return &aggr.UserStats{
		UserId:                r.UserId,
		TotalAssignmentCnt:    r.TotalAssignmentCnt,
		ActiveAssignmentCnt:   r.ActiveAssignmentCnt,
		FinishedAssignmentCnt: r.FinishedAssignmentCnt,
	}
}

// `UserStatsCreRow` maps user stats columns for create queries.
type UserStatsCreRow struct {
	Id string `gorm:"column:id;primaryKey"`

	UserId string `gorm:"column:user_id"`

	TotalAssignmentCnt    int `gorm:"column:total_assignment_count"`
	ActiveAssignmentCnt   int `gorm:"column:active_assignment_count"`
	FinishedAssignmentCnt int `gorm:"column:finished_assignment_count"`

	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

// `TableName` returns table name for `UserStatsCreRow`.
func (*UserStatsCreRow) TableName() string {
	return USER_STATS_TABLE
}

// `NewUserStatsCreRow` builds create row for one user stats record.
func NewUserStatsCreRow(userId string) *UserStatsCreRow {
	now := time.Now()

	return &UserStatsCreRow{
		Id:                    userId,
		UserId:                userId,
		TotalAssignmentCnt:    0,
		ActiveAssignmentCnt:   0,
		FinishedAssignmentCnt: 0,
		CreatedAt:             now,
		UpdatedAt:             now,
	}
}
