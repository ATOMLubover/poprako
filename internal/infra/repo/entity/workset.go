package repo_entity

import (
	"time"

	"poprako-s/internal/domain/model"
)

const WorksetTable = "workset_table"

type WorksetInfoRow struct {
	ID string `gorm:"column:id"`

	TeamID string `gorm:"column:team_id"`
	Index  int    `gorm:"column:index"`

	Name       string  `gorm:"column:name"`
	Desc       *string `gorm:"column:description"`
	ComicCount int     `gorm:"column:comic_count"`

	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func ToWorksetInfo(row WorksetInfoRow) model.WorksetInfo {
	info := model.WorksetInfo{
		ID:         row.ID,
		TeamID:     row.TeamID,
		Index:      row.Index,
		Name:       row.Name,
		ComicCount: row.ComicCount,
		CreatedAt:  row.CreatedAt,
		UpdatedAt:  row.UpdatedAt,
	}

	if row.Desc != nil {
		info.Desc = *row.Desc
	}

	return info
}
