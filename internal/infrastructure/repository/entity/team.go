package entity

import (
	"time"

	"labelplus-next-web-be/internal/domain/model"
)

const TeamTable = "team_table"

// TeamInfoRow 用于 List 查询汉化组信息。
type TeamInfoRow struct {
	ID          string     `gorm:"column:id"`
	Name        string     `gorm:"column:name"`
	Description string     `gorm:"column:description"`
	CreatedAt   time.Time  `gorm:"column:created_at"`
	UpdatedAt   time.Time  `gorm:"column:updated_at"`
	DeletedAt   *time.Time `gorm:"column:deleted_at"`
}

func (TeamInfoRow) TableName() string { return TeamTable }

// TeamInsertRow 用于 Create，仅包含写入所需字段。
type TeamInsertRow struct {
	ID          string `gorm:"column:id"`
	Name        string `gorm:"column:name"`
	Description string `gorm:"column:description"`
}

func (TeamInsertRow) TableName() string { return TeamTable }

func ToTeamInfo(row TeamInfoRow) *model.TeamInfo {
	return &model.TeamInfo{
		ID:          row.ID,
		Name:        row.Name,
		Description: row.Description,
		CreatedAt:   row.CreatedAt,
		UpdatedAt:   row.UpdatedAt,
	}
}
