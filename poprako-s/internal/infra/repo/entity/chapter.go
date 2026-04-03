package repo_entity

import (
	"time"

	"poprako-s/internal/domain/model"
)

const ChapterTable = "chapter_table"

type ChapterInfoRow struct {
	ID string `gorm:"column:id"`

	ComicID string `gorm:"column:comic_id"`
	Pinned  bool   `gorm:"column:pinned"`

	Index    int    `gorm:"column:index"`
	Subtitle string `gorm:"column:subtitle"`

	PageCount           int `gorm:"column:page_count"`
	TotalUnitCount      int `gorm:"column:total_unit_count"`
	TranslatedUnitCount int `gorm:"column:translated_unit_count"`
	ProofreadUnitCount  int `gorm:"column:proofread_unit_count"`

	UploadedAt     *time.Time `gorm:"column:uploaded_at"`
	TransalatingAt *time.Time `gorm:"column:transalating_at"`
	TranslatedAt   *time.Time `gorm:"column:translated_at"`
	ProofreadingAt *time.Time `gorm:"column:proofreading_at"`
	ProofreadAt    *time.Time `gorm:"column:proofread_at"`
	TypesettingAt  *time.Time `gorm:"column:typesetting_at"`
	TypesetAt      *time.Time `gorm:"column:typeset_at"`
	ReviewedAt     *time.Time `gorm:"column:reviewed_at"`
	PublishedAt    *time.Time `gorm:"column:published_at"`

	CreatorID string `gorm:"column:creator_id"`

	CreatedAt time.Time  `gorm:"column:created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at"`
}

func ToChapterInfo(row ChapterInfoRow) model.ChapterInfo {
	return model.ChapterInfo{
		ID:                  row.ID,
		ComicID:             row.ComicID,
		IsPinned:            row.Pinned,
		Index:               row.Index,
		Subtitle:            row.Subtitle,
		PageCount:           row.PageCount,
		TotalUnitCount:      row.TotalUnitCount,
		TranslatedUnitCount: row.TranslatedUnitCount,
		ProofreadUnitCount:  row.ProofreadUnitCount,
		CreatorID:           row.CreatorID,
		CreatedAt:           row.CreatedAt,
		UpdatedAt:           row.UpdatedAt,
		UploadedAt:          row.UploadedAt,
		TransalatingAt:      row.TransalatingAt,
		TranslatedAt:        row.TranslatedAt,
		ProofreadingAt:      row.ProofreadingAt,
		ProofreadAt:         row.ProofreadAt,
		TypesettingAt:       row.TypesettingAt,
		TypesetAt:           row.TypesetAt,
		ReviewedAt:          row.ReviewedAt,
		PublishedAt:         row.PublishedAt,
	}
}
