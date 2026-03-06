package model

import (
	"time"

	"labelplus-next-web-be/internal/util"
)

type TeamInfo struct {
	ID string

	Name        string
	Description string

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

	Name        util.Option[string]
	Description util.Option[string]
}

func NewTeamUpdate(id string) *TeamUpdate {
	return &TeamUpdate{
		ID: id,
	}
}
