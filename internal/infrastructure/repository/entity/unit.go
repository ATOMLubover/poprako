package entity

import (
	"time"

	"labelplus-next-web-be/internal/domain/model"
)

const UnitTable = "unit_table"

// UnitInfoRow 对应数据库 unit_table 行。
// 注意：x_coord / y_coord 在 SQL 中为 REAL，in_bubble 对应 model 的 IsBubble，
// proofreader_text 对应 model 的 ProofreadText。
type UnitInfoRow struct {
	ID string `gorm:"column:id"`

	PageID string `gorm:"column:page_id"`
	Index  int    `gorm:"column:index"`

	XCoord   float64 `gorm:"column:x_coord"`
	YCoord   float64 `gorm:"column:y_coord"`
	InBubble bool    `gorm:"column:in_bubble"`

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

func (UnitInfoRow) TableName() string { return UnitTable }

// UnitInsertRow 仅包含 CreateBatch 时需要写入的字段。
type UnitInsertRow struct {
	ID string `gorm:"column:id"`

	PageID string `gorm:"column:page_id"`
	Index  int    `gorm:"column:index"`

	XCoord   float64 `gorm:"column:x_coord"`
	YCoord   float64 `gorm:"column:y_coord"`
	InBubble bool    `gorm:"column:in_bubble"`

	IsProofread bool `gorm:"column:is_proofread"`

	TranslatedText    *string `gorm:"column:translated_text"`
	TranslatorID      *string `gorm:"column:translator_id"`
	TranslatorComment *string `gorm:"column:translator_comment"`

	ProofreaderText    *string `gorm:"column:proofreader_text"`
	ProofreaderID      *string `gorm:"column:proofreader_id"`
	ProofreaderComment *string `gorm:"column:proofreader_comment"`
}

func (UnitInsertRow) TableName() string { return UnitTable }

func ToUnitInfo(row UnitInfoRow) model.UnitInfo {
	unit := model.NewUnitInfo(
		row.ID,
		row.Index,
		int(row.XCoord),
		int(row.YCoord),
		row.InBubble,
		row.TranslatedText,
		row.TranslatorID,
		row.TranslatorComment,
		row.IsProofread,
		row.ProofreaderText,
		row.ProofreaderID,
		row.ProofreaderComment,
	)
	unit.PageID = row.PageID

	return unit
}
