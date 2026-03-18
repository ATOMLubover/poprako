package model

import "time"

type TeamInfo struct {
	ID string

	Name             string
	Description      string
	AvatarOSSKey     string
	IsAvatarUploaded bool

	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewTeamInfo(
	id string,
	name string,
	description string,
	avatarOSSKey string,
	isAvatarUploaded bool,
	createdAt time.Time,
	updatedAt time.Time,
) TeamInfo {
	return TeamInfo{
		ID:               id,
		Name:             name,
		Description:      description,
		AvatarOSSKey:     avatarOSSKey,
		IsAvatarUploaded: isAvatarUploaded,
		CreatedAt:        createdAt,
		UpdatedAt:        updatedAt,
	}
}

type TeamCreation struct {
	Name        string
	Description string
}

func NewTeamCreation(name, description string) *TeamCreation {
	return &TeamCreation{
		Name:        name,
		Description: description,
	}
}

type TeamUpdate struct {
	ID string

	Name        string
	Description string
}

func NewTeamUpdate(id string, name string, description string) TeamUpdate {
	return TeamUpdate{
		ID:          id,
		Name:        name,
		Description: description,
	}
}
