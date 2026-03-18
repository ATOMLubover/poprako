package model

import "time"

type PageInfo struct {
	ID string

	ChapterID  string
	Index      int
	OSSKey     string
	IsUploaded bool
	CreatorID  string
	// Creator 仅在 includes 指定时填充。
	Creator *UserInfo

	TotalUnitCount      int
	TranslatedUnitCount int
	ProofreadUnitCount  int

	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewPageInfo(
	id string,
	chapterID string,
	index int,
	ossKey string,
	isUploaded bool,
	creatorID string,
	creator *UserInfo,
	totalUnitCount int,
	translatedUnitCount int,
	proofreadUnitCount int,
	createdAt time.Time,
	updatedAt time.Time,
) PageInfo {
	return PageInfo{
		ID:                  id,
		ChapterID:           chapterID,
		Index:               index,
		OSSKey:              ossKey,
		IsUploaded:          isUploaded,
		CreatorID:           creatorID,
		Creator:             creator,
		TotalUnitCount:      totalUnitCount,
		TranslatedUnitCount: translatedUnitCount,
		ProofreadUnitCount:  proofreadUnitCount,
		CreatedAt:           createdAt,
		UpdatedAt:           updatedAt,
	}
}

type PageStats struct {
	PageID string

	TotalUnitCount      int
	TranslatedUnitCount int
	ProofreadUnitCount  int
}

func NewPageStats(
	pageID string,
	totalUnitCount int,
	translatedUnitCount int,
	proofreadUnitCount int,
) PageStats {
	return PageStats{
		PageID:              pageID,
		TotalUnitCount:      totalUnitCount,
		TranslatedUnitCount: translatedUnitCount,
		ProofreadUnitCount:  proofreadUnitCount,
	}
}

type PageCreation struct {
	ID string

	ChapterID string
	Index     int // 0-based
	OSSKey    string
	CreatorID string
}

func NewPageCreation(
	id string,
	chapterID string,
	index int,
	ossKey string,
	creatorID string,
) PageCreation {
	return PageCreation{
		ID:        id,
		ChapterID: chapterID,
		Index:     index,
		OSSKey:    ossKey,
		CreatorID: creatorID,
	}
}

type PageUpdate struct {
	ID string

	Index  int
	OSSKey string

	IsUploaded bool

	// 仅用于 unit save 时更新页面的统计数据
	TotalUnitCount      int
	TranslatedUnitCount int
	ProofreadUnitCount  int
}

func NewPageUpdate(
	id string,
	index int,
	ossKey string,
	isUploaded bool,
	totalUnitCount int,
	translatedUnitCount int,
	proofreadUnitCount int,
) PageUpdate {
	return PageUpdate{
		ID:                  id,
		Index:               index,
		OSSKey:              ossKey,
		IsUploaded:          isUploaded,
		TotalUnitCount:      totalUnitCount,
		TranslatedUnitCount: translatedUnitCount,
		ProofreadUnitCount:  proofreadUnitCount,
	}
}
