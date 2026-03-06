package entity

import (
	"time"

	"labelplus-next-web-be/internal/domain/model"
)

const ComicTable = "comic_table"

// ComicInfoRow 用于 List、GetByID，映射完整漫画信息。
type ComicInfoRow struct {
	ID string `gorm:"column:id"`

	TeamID string `gorm:"column:team_id"`

	Index       int    `gorm:"column:index"`
	Title       string `gorm:"column:title"`
	Author      string `gorm:"column:author"`
	Description string `gorm:"column:description"`

	CoverURL string `gorm:"column:cover_url"`

	ChapterCount int `gorm:"column:chapter_count"`

	CreatorID string `gorm:"column:creator_id"`

	CreatedAt time.Time  `gorm:"column:created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at"`
}

func (ComicInfoRow) TableName() string { return ComicTable }

// ComicInsertRow 用于 Create，仅包含写入所需字段。
type ComicInsertRow struct {
	ID          string `gorm:"column:id"`
	TeamID      string `gorm:"column:team_id"`
	Index       int    `gorm:"column:index"`
	Title       string `gorm:"column:title"`
	Author      string `gorm:"column:author"`
	Description string `gorm:"column:description"`
	CoverURL    string `gorm:"column:cover_url"`
	CreatorID   string `gorm:"column:creator_id"`
}

func (ComicInsertRow) TableName() string { return ComicTable }

func ToComicInfo(row ComicInfoRow) *model.ComicInfo {
	return &model.ComicInfo{
		ID:           row.ID,
		TeamID:       row.TeamID,
		Index:        row.Index,
		Title:        row.Title,
		Author:       row.Author,
		Description:  row.Description,
		CoverURL:     row.CoverURL,
		ChapterCount: row.ChapterCount,
		CreatorID:    row.CreatorID,
		CreatedAt:    row.CreatedAt,
		UpdatedAt:    row.UpdatedAt,
	}
}
