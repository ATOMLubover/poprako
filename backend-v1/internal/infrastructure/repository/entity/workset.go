package entity

import (
	"time"

	"labelplus-next-web-be/internal/domain/model"
)

const WorksetTable = "workset_table"

// WorksetInfoRow 用于 List、Get，映射完整工作集信息。
type WorksetInfoRow struct {
	ID     string `gorm:"column:id"`
	TeamID string `gorm:"column:team_id"`

	Index int `gorm:"column:index"`

	Name        string `gorm:"column:name"`
	Description string `gorm:"column:description"`
	ComicCount  int    `gorm:"column:comic_count"`

	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`

	TeamName             string    `gorm:"column:team_name"`
	TeamDescription      string    `gorm:"column:team_description"`
	TeamAvatarOSSKey     string    `gorm:"column:team_avatar_oss_key"`
	TeamIsAvatarUploaded bool      `gorm:"column:team_is_avatar_uploaded"`
	TeamCreatedAt        time.Time `gorm:"column:team_created_at"`
	TeamUpdatedAt        time.Time `gorm:"column:team_updated_at"`
}

func (WorksetInfoRow) TableName() string { return WorksetTable }

// WorksetInsertRow 用于 Create，仅包含写入所需字段。
type WorksetInsertRow struct {
	ID          string `gorm:"column:id"`
	TeamID      string `gorm:"column:team_id"`
	Index       int    `gorm:"column:index"`
	Name        string `gorm:"column:name"`
	Description string `gorm:"column:description"`
}

func (WorksetInsertRow) TableName() string { return WorksetTable }

func ToWorksetInfo(row WorksetInfoRow) model.WorksetInfo {
	var teamInfo *model.TeamInfo
	if !row.TeamCreatedAt.IsZero() {
		teamData := model.NewTeamInfo(
			row.TeamID,
			row.TeamName,
			row.TeamDescription,
			row.TeamAvatarOSSKey,
			row.TeamIsAvatarUploaded,
			row.TeamCreatedAt,
			row.TeamUpdatedAt,
		)
		teamInfo = &teamData
	}

	return model.NewWorksetInfo(
		row.ID,
		row.TeamID,
		teamInfo,
		row.Index,
		row.Name,
		row.Description,
		row.ComicCount,
		row.CreatedAt,
		row.UpdatedAt,
	)
}
