package entity

import (
	"time"

	"poprako-s/internal/domain/model/aggr"
)

// `PAGE_TABLE` is the table name for page entity.
const PAGE_TABLE = "t_page"

// `PageRow` maps one page row for read query.
type PageRow struct {
	// `Id` is the primary key.
	Id string `gorm:"column:id;primaryKey"`

	// `ChapterId` is the owner chapter identifier.
	ChapterId string `gorm:"column:chapter_id"`

	// `Index` is the chapter-scoped page order.
	Index int `gorm:"column:index"`

	// `ImageKey` is the reserved page image key.
	ImageKey *string `gorm:"column:image_key"`

	// `ImageUploaded` marks whether the page image upload is completed.
	ImageUploaded bool `gorm:"column:image_uploaded"`

	// `TotalUnitCount` is the total unit count on the page.
	TotalUnitCount int `gorm:"column:total_unit_count"`

	// `TranslatedUnitCount` is the translated unit count on the page.
	TranslatedUnitCount int `gorm:"column:translated_unit_count"`

	// `ProofreadUnitCount` is the proofread unit count on the page.
	ProofreadUnitCount int `gorm:"column:proofread_unit_count"`

	// `CreatedAt` is the creation timestamp.
	CreatedAt time.Time `gorm:"column:created_at"`

	// `UpdatedAt` is the update timestamp.
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

// `TableName` returns the table name for `PageRow`.
func (*PageRow) TableName() string {
	return PAGE_TABLE
}

// `ToPageAggr` converts one `PageRow` to `Page` aggregate.
func (r *PageRow) ToPageAggr() *aggr.Page {
	if r == nil {
		return nil
	}

	return &aggr.Page{
		Id:                  r.Id,
		ChapterId:           r.ChapterId,
		Index:               r.Index,
		ImageKey:            r.ImageKey,
		ImageUploaded:       r.ImageUploaded,
		TotalUnitCount:      r.TotalUnitCount,
		TranslatedUnitCount: r.TranslatedUnitCount,
		ProofreadUnitCount:  r.ProofreadUnitCount,
		CreatedAt:           r.CreatedAt,
		UpdatedAt:           r.UpdatedAt,
	}
}

// `PageCreRow` is the write model for page batch create.
type PageCreRow struct {
	// `Id` is the page identifier.
	Id string `gorm:"column:id;primaryKey"`

	// `ChapterId` is the owner chapter identifier.
	ChapterId string `gorm:"column:chapter_id"`

	// `Index` is the chapter-scoped page order.
	Index int `gorm:"column:index"`

	// `ImageKey` is the reserved page image key.
	ImageKey *string `gorm:"column:image_key"`

	// `ImageUploaded` marks whether the image upload is completed.
	ImageUploaded bool `gorm:"column:image_uploaded"`

	// `TotalUnitCount` is the total unit count on the page.
	TotalUnitCount int `gorm:"column:total_unit_count"`

	// `TranslatedUnitCount` is the translated unit count on the page.
	TranslatedUnitCount int `gorm:"column:translated_unit_count"`

	// `ProofreadUnitCount` is the proofread unit count on the page.
	ProofreadUnitCount int `gorm:"column:proofread_unit_count"`
}

// `TableName` returns the table name for `PageCreRow`.
func (*PageCreRow) TableName() string {
	return PAGE_TABLE
}

// `NewPageCreRowFromAggr` builds `PageCreRow` from `PageCre`.
func NewPageCreRowFromAggr(cre *aggr.PageCre) *PageCreRow {
	return &PageCreRow{
		Id:                  cre.Id,
		ChapterId:           cre.ChapterId,
		Index:               cre.Index,
		ImageKey:            cre.ImageKey,
		ImageUploaded:       false,
		TotalUnitCount:      0,
		TranslatedUnitCount: 0,
		ProofreadUnitCount:  0,
	}
}

// `PageImageUploadedUpdRow` is the write model for upload confirmation.
type PageImageUploadedUpdRow struct {
	// `ImageUploaded` marks that the image upload is completed.
	ImageUploaded bool `gorm:"column:image_uploaded"`

	// `UpdatedAt` is the update timestamp.
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

// `TableName` returns the table name for `PageImageUploadedUpdRow`.
func (*PageImageUploadedUpdRow) TableName() string {
	return PAGE_TABLE
}
