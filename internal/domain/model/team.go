package model

import "time"

type TeamInfo struct {
	ID string

	Name        string
	Description string
	AvatarOSSKey string
	IsAvatarUploaded bool

	CreatedAt time.Time
	UpdatedAt time.Time
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
