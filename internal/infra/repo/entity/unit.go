package entity

import (
	"time"

	"poprako-s/internal/domain/model/aggr"
)

// `UNIT_TABLE` is the table name for unit entity.
const UNIT_TABLE = "t_unit"

// `UnitRow` maps one unit row for read query.
type UnitRow struct {
	// `Id` is the unit identifier.
	Id string `gorm:"column:id;primaryKey"`

	// `PageId` is the owner page identifier.
	PageId string `gorm:"column:page_id"`

	// `Index` is the page-scoped unit order.
	Index int `gorm:"column:index"`

	// `IsBubble` marks whether the unit belongs to a bubble.
	IsBubble bool `gorm:"column:is_bubble"`

	// `IsProofread` marks whether the unit has been proofread.
	IsProofread bool `gorm:"column:is_proofread"`

	// `XCoord` is the unit x coordinate.
	XCoord float64 `gorm:"column:x_coord"`

	// `YCoord` is the unit y coordinate.
	YCoord float64 `gorm:"column:y_coord"`

	// `TranslatedText` is translated text content.
	TranslatedText *string `gorm:"column:translated_text"`

	// `TranslatorComment` is translator note.
	TranslatorComment *string `gorm:"column:translator_comment"`

	// `LastTranslatorId` is the last translator identifier.
	LastTranslatorId *string `gorm:"column:last_translator_id"`

	// `ProofreadText` is proofread text content.
	ProofreadText *string `gorm:"column:proofread_text"`

	// `ProofreaderComment` is proofreader note.
	ProofreaderComment *string `gorm:"column:proofreader_comment"`

	// `LastProofreaderId` is the last proofreader identifier.
	LastProofreaderId *string `gorm:"column:last_proofreader_id"`

	// `CreatedAt` is the creation timestamp.
	CreatedAt time.Time `gorm:"column:created_at"`

	// `UpdatedAt` is the update timestamp.
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

// `TableName` returns the table name for `UnitRow`.
func (*UnitRow) TableName() string {
	return UNIT_TABLE
}

// `ToUnitAggr` converts one `UnitRow` to `Unit` aggregate.
func (r *UnitRow) ToUnitAggr() *aggr.Unit {
	if r == nil {
		return nil
	}

	return &aggr.Unit{
		Id:                 r.Id,
		PageId:             r.PageId,
		Index:              r.Index,
		IsBubble:           r.IsBubble,
		IsProofread:        r.IsProofread,
		XCoord:             r.XCoord,
		YCoord:             r.YCoord,
		TranslatedText:     r.TranslatedText,
		TranslatorComment:  r.TranslatorComment,
		LastTranslatorId:   r.LastTranslatorId,
		ProofreadText:      r.ProofreadText,
		ProofreaderComment: r.ProofreaderComment,
		LastProofreaderId:  r.LastProofreaderId,
		CreatedAt:          r.CreatedAt,
		UpdatedAt:          r.UpdatedAt,
	}
}

// `UnitCreRow` is the write model for unit insert.
type UnitCreRow struct {
	// `Id` is the unit identifier.
	Id string `gorm:"column:id;primaryKey"`

	// `PageId` is the owner page identifier.
	PageId string `gorm:"column:page_id"`

	// `Index` is the page-scoped unit order.
	Index int `gorm:"column:index"`

	// `IsBubble` marks whether the unit belongs to a bubble.
	IsBubble bool `gorm:"column:is_bubble"`

	// `IsProofread` marks whether the unit has been proofread.
	IsProofread bool `gorm:"column:is_proofread"`

	// `XCoord` is the unit x coordinate.
	XCoord float64 `gorm:"column:x_coord"`

	// `YCoord` is the unit y coordinate.
	YCoord float64 `gorm:"column:y_coord"`

	// `TranslatedText` is translated text content.
	TranslatedText *string `gorm:"column:translated_text"`

	// `TranslatorComment` is translator note.
	TranslatorComment *string `gorm:"column:translator_comment"`

	// `LastTranslatorId` is the last translator identifier.
	LastTranslatorId *string `gorm:"column:last_translator_id"`

	// `ProofreadText` is proofread text content.
	ProofreadText *string `gorm:"column:proofread_text"`

	// `ProofreaderComment` is proofreader note.
	ProofreaderComment *string `gorm:"column:proofreader_comment"`

	// `LastProofreaderId` is the last proofreader identifier.
	LastProofreaderId *string `gorm:"column:last_proofreader_id"`

	// `CreatedAt` is the creation timestamp.
	CreatedAt time.Time `gorm:"column:created_at"`

	// `UpdatedAt` is the update timestamp.
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

// `TableName` returns the table name for `UnitCreRow`.
func (*UnitCreRow) TableName() string {
	return UNIT_TABLE
}

// `NewUnitCreRowFromAggr` builds `UnitCreRow` from `UnitCre`.
func NewUnitCreRowFromAggr(cre *aggr.UnitCre) *UnitCreRow {
	now := time.Now()

	return &UnitCreRow{
		Id:                 cre.LocalId,
		PageId:             cre.PageId,
		Index:              0,
		IsBubble:           cre.IsBubble,
		IsProofread:        cre.IsProofread,
		XCoord:             cre.XCoord,
		YCoord:             cre.YCoord,
		TranslatedText:     cre.TranslatedText,
		TranslatorComment:  cre.TranslatorComment,
		LastTranslatorId:   cre.LastTranslatorId,
		ProofreadText:      cre.ProofreadText,
		ProofreaderComment: cre.ProofreaderComment,
		LastProofreaderId:  cre.LastProofreaderId,
		CreatedAt:          now,
		UpdatedAt:          now,
	}
}

// `UnitUpdRow` is the write model for unit upsert.
type UnitUpdRow struct {
	// `Id` is the unit identifier.
	Id string `gorm:"column:id;primaryKey"`

	// `PageId` is the owner page identifier.
	PageId string `gorm:"column:page_id"`

	// `Index` is the page-scoped unit order.
	Index int `gorm:"column:index"`

	// `IsBubble` marks whether the unit belongs to a bubble.
	IsBubble bool `gorm:"column:is_bubble"`

	// `IsProofread` marks whether the unit has been proofread.
	IsProofread bool `gorm:"column:is_proofread"`

	// `XCoord` is the unit x coordinate.
	XCoord float64 `gorm:"column:x_coord"`

	// `YCoord` is the unit y coordinate.
	YCoord float64 `gorm:"column:y_coord"`

	// `TranslatedText` is translated text content.
	TranslatedText *string `gorm:"column:translated_text"`

	// `TranslatorComment` is translator note.
	TranslatorComment *string `gorm:"column:translator_comment"`

	// `LastTranslatorId` is the last translator identifier.
	LastTranslatorId *string `gorm:"column:last_translator_id"`

	// `ProofreadText` is proofread text content.
	ProofreadText *string `gorm:"column:proofread_text"`

	// `ProofreaderComment` is proofreader note.
	ProofreaderComment *string `gorm:"column:proofreader_comment"`

	// `LastProofreaderId` is the last proofreader identifier.
	LastProofreaderId *string `gorm:"column:last_proofreader_id"`

	// `CreatedAt` is the creation timestamp.
	CreatedAt time.Time `gorm:"column:created_at"`

	// `UpdatedAt` is the update timestamp.
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

// `TableName` returns the table name for `UnitUpdRow`.
func (*UnitUpdRow) TableName() string {
	return UNIT_TABLE
}

// `NewUnitUpdRowFromAggr` builds `UnitUpdRow` from `UnitSave`.
func NewUnitUpdRowFromAggr(sv *aggr.UnitSave) *UnitUpdRow {
	now := time.Now()

	return &UnitUpdRow{
		Id:                 sv.Id,
		PageId:             sv.PageId,
		Index:              0,
		IsBubble:           sv.IsBubble,
		IsProofread:        sv.IsProofread,
		XCoord:             sv.XCoord,
		YCoord:             sv.YCoord,
		TranslatedText:     sv.TranslatedText,
		TranslatorComment:  sv.TranslatorComment,
		LastTranslatorId:   sv.LastTranslatorId,
		ProofreadText:      sv.ProofreadText,
		ProofreaderComment: sv.ProofreaderComment,
		LastProofreaderId:  sv.LastProofreaderId,
		CreatedAt:          now,
		UpdatedAt:          now,
	}
}

// `UnitReindexUpdRow` is the write model for single unit reindex.
type UnitReindexUpdRow struct {
	// `Index` is the new page-scoped order.
	Index int `gorm:"column:index"`

	// `UpdatedAt` is the update timestamp.
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

// `TableName` returns the table name for `UnitReindexUpdRow`.
func (*UnitReindexUpdRow) TableName() string {
	return UNIT_TABLE
}
