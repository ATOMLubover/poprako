package entity

import (
	"time"

	"poprako-s/internal/domain/model/aggr"
)

// `CHAPTER_TABLE` is table name for chapter entity.
const CHAPTER_TABLE = "t_chapter"

// `ChapterRow` maps one chapter row for read queries.
type ChapterRow struct {
	// `Id` is chapter id.
	Id string `gorm:"column:id;primaryKey"`

	// `ComicId` is owner comic id.
	ComicId string `gorm:"column:comic_id"`
	// `Comic` is included when includes contain `comic`.
	Comic *ComicRow `gorm:"foreignKey:ComicId"`

	// `IsPinned` indicates pinned status.
	IsPinned bool `gorm:"column:pinned"`

	// `Index` is comic-scoped chapter index.
	Index int `gorm:"column:index"`
	// `Subtitle` is chapter subtitle.
	Subtitle string `gorm:"column:subtitle"`

	// Counters.
	PageCount           int `gorm:"column:page_count"`
	TotalUnitCount      int `gorm:"column:total_unit_count"`
	TranslatedUnitCount int `gorm:"column:translated_unit_count"`
	ProofreadUnitCount  int `gorm:"column:proofread_unit_count"`

	// Workflow timestamps.
	UploadedAt *time.Time `gorm:"column:uploaded_at"`

	TransalatingAt *time.Time `gorm:"column:transalating_at"`
	TranslatedAt   *time.Time `gorm:"column:translated_at"`

	ProofreadingAt *time.Time `gorm:"column:proofreading_at"`
	ProofreadAt    *time.Time `gorm:"column:proofread_at"`

	TypesettingAt *time.Time `gorm:"column:typesetting_at"`
	TypesetAt     *time.Time `gorm:"column:typeset_at"`

	ReviewedAt  *time.Time `gorm:"column:reviewed_at"`
	PublishedAt *time.Time `gorm:"column:published_at"`

	// `CreatorId` is creator user id.
	CreatorId string `gorm:"column:creator_id"`

	// Timestamps.
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

// `TableName` returns table name for `ChapterRow`.
func (*ChapterRow) TableName() string {
	return CHAPTER_TABLE
}

// `ToChapterAggr` converts one row to chapter aggregate.
func (r *ChapterRow) ToChapterAggr() *aggr.Chapter {
	if r == nil {
		return nil
	}

	var comic *aggr.Comic
	if r.Comic != nil {
		comic = r.Comic.ToComicAggr()
	}

	return &aggr.Chapter{
		Id:                  r.Id,
		ComicId:             r.ComicId,
		Comic:               comic,
		IsPinned:            r.IsPinned,
		Index:               r.Index,
		Subtitle:            r.Subtitle,
		PageCount:           r.PageCount,
		TotalUnitCount:      r.TotalUnitCount,
		TranslatedUnitCount: r.TranslatedUnitCount,
		ProofreadUnitCount:  r.ProofreadUnitCount,
		UploadedAt:          r.UploadedAt,
		TransalatingAt:      r.TransalatingAt,
		TranslatedAt:        r.TranslatedAt,
		ProofreadingAt:      r.ProofreadingAt,
		ProofreadAt:         r.ProofreadAt,
		TypesettingAt:       r.TypesettingAt,
		TypesetAt:           r.TypesetAt,
		ReviewedAt:          r.ReviewedAt,
		PublishedAt:         r.PublishedAt,
		CreatorId:           r.CreatorId,
		CreatedAt:           r.CreatedAt,
		UpdatedAt:           r.UpdatedAt,
	}
}

// `ChapterCreRow` is write model for chapter create.
type ChapterCreRow struct {
	// `Id` is generated chapter identifier.
	Id string `gorm:"column:id;primaryKey"`

	// `ComicId` is owner comic identifier.
	ComicId string `gorm:"column:comic_id"`
	// `IsPinned` is always true for a newly created chapter.
	IsPinned bool `gorm:"column:pinned"`

	// `Index` is comic-scoped chapter sequence index.
	Index int `gorm:"column:index"`
	// `Subtitle` is resolved chapter subtitle.
	Subtitle string `gorm:"column:subtitle"`

	// `CreatorId` is creator user identifier.
	CreatorId string `gorm:"column:creator_id"`
}

// `TableName` returns table name for `ChapterCreRow`.
func (*ChapterCreRow) TableName() string {
	return CHAPTER_TABLE
}

// `NewChapterCreRowFromAggr` builds chapter create row from aggregate.
func NewChapterCreRowFromAggr(cre *aggr.ChapterCre, subtitle string) *ChapterCreRow {
	return &ChapterCreRow{
		Id:        cre.Id,
		ComicId:   cre.ComicId,
		IsPinned:  true,
		Index:     cre.Index,
		Subtitle:  subtitle,
		CreatorId: cre.CreatorId,
	}
}

// `ChapterUpdRow` is write model for chapter update.
// Only columns listed in `Select` are written; pointer fields carry nil-as-NULL semantics.
type ChapterUpdRow struct {
	// `Subtitle` is the new subtitle value.
	Subtitle *string `gorm:"column:subtitle"`
	// `IsPinned` is the new pinned flag.
	IsPinned *bool `gorm:"column:pinned"`

	// `UploadedAt` is the upload completion timestamp.
	UploadedAt *time.Time `gorm:"column:uploaded_at"`

	// `TransalatingAt` is the translate-start timestamp.
	TransalatingAt *time.Time `gorm:"column:transalating_at"`
	// `TranslatedAt` is the translate-complete timestamp.
	TranslatedAt *time.Time `gorm:"column:translated_at"`

	// `ProofreadingAt` is the proofread-start timestamp.
	ProofreadingAt *time.Time `gorm:"column:proofreading_at"`
	// `ProofreadAt` is the proofread-complete timestamp.
	ProofreadAt *time.Time `gorm:"column:proofread_at"`

	// `TypesettingAt` is the typeset-start timestamp.
	TypesettingAt *time.Time `gorm:"column:typesetting_at"`
	// `TypesetAt` is the typeset-complete timestamp.
	TypesetAt *time.Time `gorm:"column:typeset_at"`

	// `ReviewedAt` is the review-complete timestamp.
	ReviewedAt *time.Time `gorm:"column:reviewed_at"`
	// `PublishedAt` is the publish-complete timestamp.
	PublishedAt *time.Time `gorm:"column:published_at"`

	// `UpdatedAt` is the row update timestamp.
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

// `TableName` returns table name for `ChapterUpdRow`.
func (*ChapterUpdRow) TableName() string {
	return CHAPTER_TABLE
}

// `ChapterPinUpdRow` is write model for bulk pin reset.
type ChapterPinUpdRow struct {
	// `IsPinned` is the new pinned value written to all matched rows.
	IsPinned bool `gorm:"column:pinned"`
	// `UpdatedAt` is the row update timestamp.
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

// `TableName` returns table name for `ChapterPinUpdRow`.
func (*ChapterPinUpdRow) TableName() string {
	return CHAPTER_TABLE
}

// `ChapterRemoveUpdRow` is write model for chapter soft delete.
type ChapterRemoveUpdRow struct {
	// `DeletedAt` is the soft-delete timestamp written on remove.
	DeletedAt time.Time `gorm:"column:deleted_at"`
	// `UpdatedAt` is the row update timestamp.
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

// `TableName` returns table name for `ChapterRemoveUpdRow`.
func (*ChapterRemoveUpdRow) TableName() string {
	return CHAPTER_TABLE
}

// `ChapterPageCountUpdRow` is write model for page count overwrite.
type ChapterPageCountUpdRow struct {
	// `PageCount` is the new page count value.
	PageCount int `gorm:"column:page_count"`

	// `UpdatedAt` is the row update timestamp.
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

// `TableName` returns table name for `ChapterPageCountUpdRow`.
func (*ChapterPageCountUpdRow) TableName() string {
	return CHAPTER_TABLE
}
