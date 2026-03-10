package model

import "time"

type ChapterInfo struct {
	ID string

	ComicID   string
	Index     int
	ChapterNo string

	PageCount           int
	TotalUnitCount      int
	TranslatedUnitCount int
	ProofreadUnitCount  int

	CoverURL string

	UploadedAt     *time.Time
	TransalatingAt *time.Time
	TranslatedAt   *time.Time
	ProofreadingAt *time.Time
	ProofreadAt    *time.Time
	TypesettingAt  *time.Time
	TypesetAt      *time.Time
	ReviewedAt     *time.Time
	PublishedAt    *time.Time

	CreatorID string

	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewChapterInfo(
	id string,
	comicID string,
	index int,
	chapterNo string,
	pageCount int,
	totalUnitCount int,
	translatedUnitCount int,
	proofreadUnitCount int,
	coverURL string,
	uploadedAt *time.Time,
	transalatingAt *time.Time,
	translatedAt *time.Time,
	proofreadingAt *time.Time,
	proofreadAt *time.Time,
	typesettingAt *time.Time,
	typesetAt *time.Time,
	reviewedAt *time.Time,
	publishedAt *time.Time,
	creatorID string,
	createdAt time.Time,
	updatedAt time.Time,
) ChapterInfo {
	return ChapterInfo{
		ID:                  id,
		ComicID:             comicID,
		Index:               index,
		ChapterNo:           chapterNo,
		PageCount:           pageCount,
		TotalUnitCount:      totalUnitCount,
		TranslatedUnitCount: translatedUnitCount,
		ProofreadUnitCount:  proofreadUnitCount,
		CoverURL:            coverURL,
		UploadedAt:          uploadedAt,
		TransalatingAt:      transalatingAt,
		TranslatedAt:        translatedAt,
		ProofreadingAt:      proofreadAt,
		ProofreadAt:         proofreadAt,
		TypesettingAt:       typesettingAt,
		TypesetAt:           typesetAt,
		ReviewedAt:          reviewedAt,
		PublishedAt:         publishedAt,
		CreatorID:           creatorID,
		CreatedAt:           createdAt,
		UpdatedAt:           updatedAt,
	}
}

type ChapterCreation struct {
	ComicID   string
	Index     int
	ChapterNo string
	CreatorID string
}

func NewChapterCreation(
	comicID string,
	index int,
	chapterNo string,
	creatorID string,
) ChapterCreation {
	return ChapterCreation{
		ComicID:   comicID,
		Index:     index,
		ChapterNo: chapterNo,
		CreatorID: creatorID,
	}
}
