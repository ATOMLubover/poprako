package model

import "time"

type ComicInfo struct {
	ID     string
	TeamID string

	Index       int
	Title       string
	Author      string
	Description string

	CoverURL string

	ChapterCount int

	CreatorID string

	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewComicInfo(
	id string,
	teamID string,
	index int,
	title string,
	author string,
	description string,
	coverURL string,
	chapterCount int,
	creatorID string,
	createdAt time.Time,
	updatedAt time.Time,
) *ComicInfo {
	return &ComicInfo{
		ID:           id,
		TeamID:       teamID,
		Index:        index,
		Title:        title,
		Author:       author,
		Description:  description,
		CoverURL:     coverURL,
		ChapterCount: chapterCount,
		CreatorID:    creatorID,
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
	}
}

type ComicCreation struct {
	TeamID      string
	Title       string
	Author      string
	Description string
	CreatorID   string
}

func NewComicCreation(
	teamID string,
	title string,
	author string,
	description string,
	creatorID string,
) *ComicCreation {
	return &ComicCreation{
		TeamID:      teamID,
		Title:       title,
		Author:      author,
		Description: description,
		CreatorID:   creatorID,
	}
}

type ComicUpdate struct {
	ID          string
	Title       *string
	Author      *string
	Description *string
}

func NewComicUpdate(
	id string,
	title *string,
	author *string,
	description *string,
) *ComicUpdate {
	return &ComicUpdate{
		ID:          id,
		Title:       title,
		Author:      author,
		Description: description,
	}
}
