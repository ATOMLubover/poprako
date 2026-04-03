package repo_entity

import (
	"time"

	"poprako-s/internal/domain/model"
)

const UnitTable = "unit_table"

type UnitInfoRow struct {
	ID string `gorm:"column:id"`

	PageID string `gorm:"column:page_id"`

	XCoord float64 `gorm:"column:x_coord"`
	YCoord float64 `gorm:"column:y_coord"`

	Index       int  `gorm:"column:index"`
	InBubble    bool `gorm:"column:in_bubble"`
	IsProofread bool `gorm:"column:is_proofread"`

	TranslatedText    *string `gorm:"column:translated_text"`
	TranslatorID      *string `gorm:"column:translator_id"`
	TranslatorComment *string `gorm:"column:translator_comment"`

	ProofreaderText    *string `gorm:"column:proofreader_text"`
	ProofreaderID      *string `gorm:"column:proofreader_id"`
	ProofreaderComment *string `gorm:"column:proofreader_comment"`

	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func ToUnitInfo(row UnitInfoRow) model.UnitInfo {
	return model.UnitInfo{
		ID:                 row.ID,
		PageID:             row.PageID,
		Index:              row.Index,
		XCoord:             int(row.XCoord),
		YCoord:             int(row.YCoord),
		IsBubble:           row.InBubble,
		TranslatedText:     row.TranslatedText,
		TranslatorID:       row.TranslatorID,
		TranslatorComment:  row.TranslatorComment,
		IsProofread:        row.IsProofread,
		ProofreadText:      row.ProofreaderText,
		ProofreaderID:      row.ProofreaderID,
		ProofreaderComment: row.ProofreaderComment,
	}
}
