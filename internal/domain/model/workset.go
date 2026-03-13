package model

import "time"

type WorksetInfo struct {
	ID     string
	TeamID string
	Team   *TeamInfo

	Index int

	Name        string
	Description string
	ComicCount  int

	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewWorksetInfo(
	id string,
	teamID string,
	team *TeamInfo,
	index int,
	name string,
	description string,
	comicCount int,
	createdAt time.Time,
	updatedAt time.Time,
) WorksetInfo {
	return WorksetInfo{
		ID:          id,
		TeamID:      teamID,
		Team:        team,
		Index:       index,
		Name:        name,
		Description: description,
		ComicCount:  comicCount,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}
}

type WorksetCreation struct {
	TeamID      string
	Index       int
	Name        string
	Description string
}

func NewWorksetCreation(
	teamID string,
	index int,
	name string,
	description string,
) *WorksetCreation {
	return &WorksetCreation{
		TeamID:      teamID,
		Index:       index,
		Name:        name,
		Description: description,
	}
}

type WorksetUpdate struct {
	ID          string
	Name        string
	Description *string
}

func NewWorksetUpdate(
	id string,
	name string,
	description *string,
) WorksetUpdate {
	return WorksetUpdate{
		ID:          id,
		Name:        name,
		Description: description,
	}
}
