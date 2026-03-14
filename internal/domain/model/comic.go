package model

import "time"

type ComicInfo struct {
	ID        string
	WorksetID string
	// TeamID 通过 JOIN workset_table 填充，不是 comic_table 的直接列。
	TeamID string

	// Workset 仅在 includes 指定时填充。
	Workset *WorksetInfo

	Index       int
	Title       string
	Author      string
	Description string

	CoverURL string

	ChapterCount int
	CreatorID    string
	// Creator 仅在 includes 指定时填充。
	Creator *UserInfo

	LastActiveAt time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewComicInfo(
	id string,
	worksetID string,
	teamID string,
	workset *WorksetInfo,
	index int,
	title string,
	author string,
	description string,
	coverURL string,
	chapterCount int,
	creatorID string,
	creator *UserInfo,
	lastActiveAt time.Time,
	createdAt time.Time,
	updatedAt time.Time,
) ComicInfo {
	return ComicInfo{
		ID:           id,
		WorksetID:    worksetID,
		TeamID:       teamID,
		Workset:      workset,
		Index:        index,
		Title:        title,
		Author:       author,
		Description:  description,
		CoverURL:     coverURL,
		ChapterCount: chapterCount,
		CreatorID:    creatorID,
		Creator:      creator,
		LastActiveAt: lastActiveAt,
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
	}
}

type ComicCreation struct {
	WorksetID   string
	Index       int
	Title       string
	Author      string
	Description string
	CreatorID   string
}

func NewComicCreation(
	worksetID string,
	index int,
	title string,
	author string,
	description string,
	creatorID string,
) *ComicCreation {
	return &ComicCreation{
		WorksetID:   worksetID,
		Index:       index,
		Title:       title,
		Author:      author,
		Description: description,
		CreatorID:   creatorID,
	}
}

type ComicUpdate struct {
	ID          string
	Title       string
	Author      string
	Description string
}

func NewComicUpdate(
	id string,
	title string,
	author string,
	description string,
) ComicUpdate {
	return ComicUpdate{
		ID:          id,
		Title:       title,
		Author:      author,
		Description: description,
	}
}
