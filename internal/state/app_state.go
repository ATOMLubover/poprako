package state

import (
	"labelplus-next-web-be/internal/application"
	"labelplus-next-web-be/internal/config"
)

type AppState struct {
	AppConfig *config.AppConfig

	UserApplication       application.UserApplication
	InvitationApplication application.InvitationApplication
	TeamApplication       application.TeamApplication
	MemberApplication     application.MemberApplication
	ComicApplication      application.ComicApplication
}

func NewAppState(
	appConfig *config.AppConfig,
	userApplication application.UserApplication,
	invitationApplication application.InvitationApplication,
	teamApplication application.TeamApplication,
	memberApplication application.MemberApplication,
	comicApplication application.ComicApplication,
) *AppState {
	return &AppState{
		AppConfig:             appConfig,
		UserApplication:       userApplication,
		InvitationApplication: invitationApplication,
		TeamApplication:       teamApplication,
		MemberApplication:     memberApplication,
		ComicApplication:      comicApplication,
	}
}
