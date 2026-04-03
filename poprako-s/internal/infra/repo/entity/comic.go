package repo_entity

import (
	"time"

	"poprako-s/internal/domain/model"
)

const ComicTable = "comic_table"

type ComicInfoRow struct {
	ID string `gorm:"column:id"`

	WorksetID string `gorm:"column:workset_id"`

	Index         int     `gorm:"column:index"`
	Title         string  `gorm:"column:title"`
	Author        string  `gorm:"column:author"`
	ComposedTitle string  `gorm:"column:composed_title"`
	Description   *string `gorm:"column:description"`

	ChapterCount int `gorm:"column:chapter_count"`

	HasPinnedChapter     bool       `gorm:"column:has_pinned_chapter"`
	PinnedUploadedAt     *time.Time `gorm:"column:pinned_uploaded_at"`
	PinnedTransalatingAt *time.Time `gorm:"column:pinned_transalating_at"`
	PinnedTranslatedAt   *time.Time `gorm:"column:pinned_translated_at"`
	PinnedProofreadingAt *time.Time `gorm:"column:pinned_proofreading_at"`
	PinnedProofreadAt    *time.Time `gorm:"column:pinned_proofread_at"`
	PinnedTypesettingAt  *time.Time `gorm:"column:pinned_typesetting_at"`
	PinnedTypesetAt      *time.Time `gorm:"column:pinned_typeset_at"`
	PinnedReviewedAt     *time.Time `gorm:"column:pinned_reviewed_at"`
	PinnedPublishedAt    *time.Time `gorm:"column:pinned_published_at"`

	CreatorID    string    `gorm:"column:creator_id"`
	LastActiveAt time.Time `gorm:"column:last_active_at"`

	CreatedAt time.Time  `gorm:"column:created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at"`
}

func ToComicInfo(row ComicInfoRow) model.ComicInfo {
	info := model.ComicInfo{
		ID:           row.ID,
		WorksetID:    row.WorksetID,
		Index:        row.Index,
		Title:        row.Title,
		Author:       row.Author,
		ChapterCount: row.ChapterCount,
		CreatorID:    row.CreatorID,
		LastActiveAt: row.LastActiveAt,
		CreatedAt:    row.CreatedAt,
		UpdatedAt:    row.UpdatedAt,
	}

	if row.Description != nil {
		info.Description = *row.Description
	}

	return info
}
