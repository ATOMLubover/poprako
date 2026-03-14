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
	// Creator 仅在 includes 指定时填充。
	Creator *UserInfo

	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewChapterDetail(
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

type ChapterWithComicInfo struct {
	ID string

	Comic     ComicInfo
	Index     int
	ChapterNo string

	CoverURL string

	PageCount           int
	TotalUnitCount      int
	TranslatedUnitCount int
	ProofreadUnitCount  int

	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewChapterWithComicInfo(
	id string,
	comicInfo ComicInfo,
	index int,
	chapterNo string,
	coverURL string,
	pageCount int,
	totalUnitCount int,
	translatedUnitCount int,
	proofreadUnitCount int,
	createdAt time.Time,
	updatedAt time.Time,
) ChapterWithComicInfo {
	return ChapterWithComicInfo{
		ID:                  id,
		Comic:               comicInfo,
		Index:               index,
		ChapterNo:           chapterNo,
		CoverURL:            coverURL,
		PageCount:           pageCount,
		TotalUnitCount:      totalUnitCount,
		TranslatedUnitCount: translatedUnitCount,
		ProofreadUnitCount:  proofreadUnitCount,
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

type ChapterUpdate struct {
	ID        string
	ChapterNo string

	UploadedAt     *time.Time
	TransalatingAt *time.Time
	TranslatedAt   *time.Time
	ProofreadingAt *time.Time
	ProofreadAt    *time.Time
	TypesettingAt  *time.Time
	TypesetAt      *time.Time
	ReviewedAt     *time.Time
	PublishedAt    *time.Time
}

func NewChapterUpdate(
	id string,
	chapterNo *string,
	current ChapterInfo,
	uploadStatus *WorkflowStatus,
	translateStatus *WorkflowStatus,
	proofreadStatus *WorkflowStatus,
	typesetStatus *WorkflowStatus,
	reviewStatus *WorkflowStatus,
	publishStatus *WorkflowStatus,
) ChapterUpdate {
	now := time.Now()

	deref := func(s *WorkflowStatus) WorkflowStatus {
		if s == nil {
			return WorkflowUnset
		}
		return *s
	}

	chapterNoVal := current.ChapterNo
	if chapterNo != nil {
		chapterNoVal = *chapterNo
	}

	update := ChapterUpdate{
		ID:        id,
		ChapterNo: chapterNoVal,
	}

	// Upload 仅支持 pending / completed
	switch deref(uploadStatus) {
	case WorkflowPending:
		update.UploadedAt = nil
	case WorkflowCompleted:
		t := now
		update.UploadedAt = &t
	default: // WorkflowUnset：保留当前值
		update.UploadedAt = current.UploadedAt
	}

	// Translate 支持 pending / in_progress / completed
	switch deref(translateStatus) {
	case WorkflowPending:
		update.TransalatingAt = nil
		update.TranslatedAt = nil
	case WorkflowInProgress:
		t := now
		update.TransalatingAt = &t
		update.TranslatedAt = nil
	case WorkflowCompleted:
		update.TransalatingAt = current.TransalatingAt // 保留开始时间
		t := now
		update.TranslatedAt = &t
	default: // WorkflowUnset：保留当前值
		update.TransalatingAt = current.TransalatingAt
		update.TranslatedAt = current.TranslatedAt
	}

	// Proofread 支持 pending / in_progress / completed
	switch deref(proofreadStatus) {
	case WorkflowPending:
		update.ProofreadingAt = nil
		update.ProofreadAt = nil
	case WorkflowInProgress:
		t := now
		update.ProofreadingAt = &t
		update.ProofreadAt = nil
	case WorkflowCompleted:
		update.ProofreadingAt = current.ProofreadingAt // 保留开始时间
		t := now
		update.ProofreadAt = &t
	default: // WorkflowUnset：保留当前值
		update.ProofreadingAt = current.ProofreadingAt
		update.ProofreadAt = current.ProofreadAt
	}

	// Typeset 支持 pending / in_progress / completed
	switch deref(typesetStatus) {
	case WorkflowPending:
		update.TypesettingAt = nil
		update.TypesetAt = nil
	case WorkflowInProgress:
		t := now
		update.TypesettingAt = &t
		update.TypesetAt = nil
	case WorkflowCompleted:
		update.TypesettingAt = current.TypesettingAt // 保留开始时间
		t := now
		update.TypesetAt = &t
	default: // WorkflowUnset：保留当前值
		update.TypesettingAt = current.TypesettingAt
		update.TypesetAt = current.TypesetAt
	}

	// Review 仅支持 pending / completed
	switch deref(reviewStatus) {
	case WorkflowPending:
		update.ReviewedAt = nil
	case WorkflowCompleted:
		t := now
		update.ReviewedAt = &t
	default: // WorkflowUnset：保留当前值
		update.ReviewedAt = current.ReviewedAt
	}

	// Publish 仅支持 pending / completed
	switch deref(publishStatus) {
	case WorkflowPending:
		update.PublishedAt = nil
	case WorkflowCompleted:
		t := now
		update.PublishedAt = &t
	default: // WorkflowUnset：保留当前值
		update.PublishedAt = current.PublishedAt
	}

	return update
}
