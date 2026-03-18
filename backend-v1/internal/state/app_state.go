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
	WorksetApplication    application.WorksetApplication
	ComicApplication      application.ComicApplication
	ChapterApplication    application.ChapterApplication
	PageApplication       application.PageApplication
	AssignmentApplication application.AssignmentApplication
	UnitApplication       application.UnitApplication
}

func NewAppState(
	appConfig *config.AppConfig,
	userApplication application.UserApplication,
	invitationApplication application.InvitationApplication,
	teamApplication application.TeamApplication,
	memberApplication application.MemberApplication,
	worksetApplication application.WorksetApplication,
	comicApplication application.ComicApplication,
	chapterApplication application.ChapterApplication,
	pageApplication application.PageApplication,
	assignmentApplication application.AssignmentApplication,
	unitApplication application.UnitApplication,
) *AppState {
	return &AppState{
		AppConfig:             appConfig,
		UserApplication:       userApplication,
		InvitationApplication: invitationApplication,
		TeamApplication:       teamApplication,
		MemberApplication:     memberApplication,
		WorksetApplication:    worksetApplication,
		ComicApplication:      comicApplication,
		ChapterApplication:    chapterApplication,
		PageApplication:       pageApplication,
		AssignmentApplication: assignmentApplication,
		UnitApplication:       unitApplication,
	}
}
