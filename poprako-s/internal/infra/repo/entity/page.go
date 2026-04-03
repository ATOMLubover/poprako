package repo_entity

import (
	"time"

	"poprako-s/internal/domain/model"
)

const PageTable = "page_table"

type PageInfoRow struct {
	ID string `gorm:"column:id"`

	ChapterID string `gorm:"column:chapter_id"`

	Index     int     `gorm:"column:index"`
	OSSKey    *string `gorm:"column:oss_key"`
	Uploaded  bool    `gorm:"column:uploaded"`
	CreatorID string  `gorm:"column:creator_id"`

	TotalUnitCount      int `gorm:"column:total_unit_count"`
	TranslatedUnitCount int `gorm:"column:translated_unit_count"`
	ProofreadUnitCount  int `gorm:"column:proofread_unit_count"`

	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func ToPageInfo(row PageInfoRow) model.PageInfo {
	info := model.PageInfo{
		ID:                  row.ID,
		ChapterID:           row.ChapterID,
		Index:               row.Index,
		IsUploaded:          row.Uploaded,
		CreatorID:           row.CreatorID,
		TotalUnitCount:      row.TotalUnitCount,
		TranslatedUnitCount: row.TranslatedUnitCount,
		ProofreadUnitCount:  row.ProofreadUnitCount,
		CreatedAt:           row.CreatedAt,
		UpdatedAt:           row.UpdatedAt,
	}

	if row.OSSKey != nil {
		info.OSSKey = *row.OSSKey
	}

	return info
}
