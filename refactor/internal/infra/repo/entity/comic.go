package entity

import (
	"time"

	"poprako-s/internal/domain/model/aggr"
)

// `COMIC_TABLE` is the table name for comic entity
const COMIC_TABLE = "t_comic"

// `ComicRow` maps a full comic row for read query
// It is used for active rows with deleted filter in repo layer
// Relation fields are optional and include-driven
// This struct reflects migration columns exactly
// Keep nullable semantics for `Desc`
type ComicRow struct {
	// `Id` is the primary key
	Id string `gorm:"column:id;primaryKey"`

	// `WorksetId` is the owner workset identifier
	WorksetId string `gorm:"column:workset_id"`
	// `Workset` is included when `includes` contains `workset`
	Workset *WorksetRow `gorm:"foreignKey:WorksetId"`

	// `Index` is workset-scoped comic sequence
	Index int `gorm:"column:index"`

	// `Title` is comic title
	Title string `gorm:"column:title"`
	// `Author` is comic author
	Author string `gorm:"column:author"`
	// `FuzzyTitle` is dedicated fuzzy-match text for comic search.
	FuzzyTitle string `gorm:"column:fuzzy_title"`
	// `Desc` is optional comic description
	Desc *string `gorm:"column:description"`
	// `IsCompleted` marks whether the comic is completed.
	IsCompleted bool `gorm:"column:is_completed"`

	// `ChapterCount` is denormalized chapter count
	ChapterCount int `gorm:"column:chapter_count"`

	// `CreatorId` is creator user id
	CreatorId string `gorm:"column:creator_id"`

	// `LastActiveAt` is last activity timestamp
	LastActiveAt time.Time `gorm:"column:last_active_at"`
	// `CreatedAt` is create time
	CreatedAt time.Time `gorm:"column:created_at"`
	// `UpdatedAt` is update time
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

// `TableName` returns table name for `ComicRow`
func (*ComicRow) TableName() string {
	return COMIC_TABLE
}

// `ToComicAggr` converts one `ComicRow` to `Comic` aggregate
func (r *ComicRow) ToComicAggr() *aggr.Comic {
	if r == nil {
		return nil
	}

	var workset *aggr.Workset

	if r.Workset != nil {
		workset = r.Workset.ToWorksetAggr()
	}

	return &aggr.Comic{
		Id:           r.Id,
		WorksetId:    r.WorksetId,
		Workset:      workset,
		Index:        r.Index,
		Title:        r.Title,
		Author:       r.Author,
		Desc:         r.Desc,
		IsCompleted:  r.IsCompleted,
		ChapterCount: r.ChapterCount,
		CreatorId:    r.CreatorId,
		LastActiveAt: r.LastActiveAt,
		CreatedAt:    r.CreatedAt,
		UpdatedAt:    r.UpdatedAt,
	}
}

// `ComicCreRow` is write model for comic create
// Fields are minimal and align to migration
// Server-side default columns are not duplicated here
type ComicCreRow struct {
	// `Id` is generated comic id
	Id string `gorm:"column:id;primaryKey"`

	// `WorksetId` is owner workset id
	WorksetId string `gorm:"column:workset_id"`
	// `Index` is workset-scoped sequence index
	Index int `gorm:"column:index"`

	// `Title` is comic title
	Title string `gorm:"column:title"`
	// `Author` is comic author
	Author string `gorm:"column:author"`
	// `FuzzyTitle` is dedicated fuzzy-match text for comic search.
	FuzzyTitle string `gorm:"column:fuzzy_title"`
	// `Desc` is optional comic description
	Desc *string `gorm:"column:description"`

	// `CreatorId` is creator user id
	CreatorId string `gorm:"column:creator_id"`
}

// `TableName` returns table name for `ComicCreRow`
func (*ComicCreRow) TableName() string {
	return COMIC_TABLE
}

// `NewComicCreRowFromAggr` builds `ComicCreRow` from `ComicCre`
func NewComicCreRowFromAggr(cre *aggr.ComicCre, fuzzyTitle string) *ComicCreRow {
	return &ComicCreRow{
		Id:         cre.Id,
		WorksetId:  cre.WorksetId,
		Index:      cre.Index,
		Title:      cre.Title,
		Author:     cre.Author,
		FuzzyTitle: fuzzyTitle,
		Desc:       cre.Desc,
		CreatorId:  cre.CreatorId,
	}
}

// `ComicUpdRow` is write model for comic put update
// Nil `Desc` should be written as SQL NULL through explicit select
// `UpdatedAt` is set by repo when update executes
type ComicUpdRow struct {
	// `Title` is new comic title
	Title string `gorm:"column:title"`
	// `Author` is new comic author
	Author string `gorm:"column:author"`
	// `FuzzyTitle` is refreshed fuzzy-match text.
	FuzzyTitle string `gorm:"column:fuzzy_title"`
	// `Desc` is new optional description
	Desc *string `gorm:"column:description"`

	// `UpdatedAt` is update timestamp
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

// `TableName` returns table name for `ComicUpdRow`
func (*ComicUpdRow) TableName() string {
	return COMIC_TABLE
}

// `ComicChapterCountUpdRow` is write model for chapter counter update.
type ComicChapterCountUpdRow struct {
	// `ChapterCount` is updated chapter counter.
	ChapterCount int `gorm:"column:chapter_count"`

	// `UpdatedAt` is update timestamp.
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

// `TableName` returns table name for `ComicChapterCountUpdRow`.
func (*ComicChapterCountUpdRow) TableName() string {
	return COMIC_TABLE
}

// `ComicLastActiveUpdRow` is write model for last active update.
type ComicLastActiveUpdRow struct {
	// `LastActiveAt` is refreshed activity timestamp.
	LastActiveAt time.Time `gorm:"column:last_active_at"`

	// `UpdatedAt` is update timestamp.
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

// `TableName` returns table name for `ComicLastActiveUpdRow`.
func (*ComicLastActiveUpdRow) TableName() string {
	return COMIC_TABLE
}
