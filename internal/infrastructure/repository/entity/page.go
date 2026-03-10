package entity

import (
	"time"

	"labelplus-next-web-be/internal/domain/model"
)

const (
	ChapterTable = "chapter_table"
	PageTable    = "page_table"
)

type ChapterInfoRow struct {
	ID string `gorm:"column:id"`

	ComicID   string `gorm:"column:comic_id"`
	Index     int    `gorm:"column:index"`
	ChapterNo string `gorm:"column:subtitle"`

	PageCount int `gorm:"column:page_count"`

	CoverURL string `gorm:"column:cover_url"`

	CreatorID string `gorm:"column:creator_id"`

	CreatedAt time.Time  `gorm:"column:created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at"`
}

func (ChapterInfoRow) TableName() string { return ChapterTable }

type ChapterInsertRow struct {
	ID        string `gorm:"column:id"`
	ComicID   string `gorm:"column:comic_id"`
	Index     int    `gorm:"column:index"`
	ChapterNo string `gorm:"column:subtitle"`
	CreatorID string `gorm:"column:creator_id"`
}

func (ChapterInsertRow) TableName() string { return ChapterTable }

type PageInfoRow struct {
	ID string `gorm:"column:id"`

	ChapterID string `gorm:"column:chapter_id"`
	Index     int    `gorm:"column:index"`
	OSSKey    string `gorm:"column:oss_key"`

	IsUploaded bool `gorm:"column:uploaded"`

	TotalUnitCount      int `gorm:"column:total_unit_count"`
	TranslatedUnitCount int `gorm:"column:translated_unit_count"`
	ProofreadUnitCount  int `gorm:"column:proved_unit_count"`

	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (PageInfoRow) TableName() string { return PageTable }

type PageInsertRow struct {
	ID string `gorm:"column:id"`

	ChapterID string `gorm:"column:chapter_id"`
	Index     int    `gorm:"column:index"`
	OSSKey    string `gorm:"column:oss_key"`
	CreatorID string `gorm:"column:creator_id"`
}

func (PageInsertRow) TableName() string { return PageTable }

func ToPageInfo(row PageInfoRow) model.PageInfo {
	return model.NewPageInfo(
		row.ID,
		row.ChapterID,
		row.Index,
		row.OSSKey,
		row.TotalUnitCount,
		row.TranslatedUnitCount,
		row.ProofreadUnitCount,
		row.CreatedAt,
		row.UpdatedAt,
	)
}

func ToChapterInfo(row ChapterInfoRow) model.ChapterInfo {
	return model.NewChapterInfo(
		row.ID,
		row.ComicID,
		row.Index,
		row.ChapterNo,
		row.PageCount,
		0,
		0,
		0,
		row.CoverURL,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		row.CreatorID,
		row.CreatedAt,
		row.UpdatedAt,
	)
}
