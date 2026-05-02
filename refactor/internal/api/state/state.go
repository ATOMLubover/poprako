package state

import (
	app_iface "poprako-s/internal/app"
	"poprako-s/internal/cfg"
)

// `AppState` holds all application-wide dependencies shared across HTTP handlers.
type AppState struct {
	// `Cfg` is the application configuration.
	Cfg *cfg.AppCfg

	// `UserApp` handles user-related use-cases.
	UserApp app_iface.UserApp

	// `TeamApp` handles team-related use-cases.
	TeamApp app_iface.TeamApp

	// `MemberApp` handles member-related use-cases.
	MemberApp app_iface.MemberApp

	// `MemberInvApp` handles member invitation use-cases.
	MemberInvApp app_iface.MemberInvApp

	// `WorksetApp` handles workset-related use-cases.
	WorksetApp app_iface.WorksetApp

	// `ComicApp` handles comic-related use-cases.
	ComicApp app_iface.ComicApp

	// `ChapterApp` handles chapter-related use-cases.
	ChapterApp app_iface.ChapterApp

	// `ChapterPortApp` handles chapter import and export use-cases.
	ChapterPortApp app_iface.ChapterPortApp

	// `PageApp` handles page-related use-cases.
	PageApp app_iface.PageApp

	// `UnitApp` handles unit-related use-cases.
	UnitApp app_iface.UnitApp

	// `SysMailApp` handles system mail related use-cases.
	SysMailApp app_iface.SysMailApp

	// `AssignmentInvApp` handles assignment invitation use-cases.
	AssignmentInvApp app_iface.AssignmentInvApp

	// `AssignmentApp` handles assignment use-cases.
	AssignmentApp app_iface.AssignmentApp
}

// `NewAppState` constructs an `AppState` from its dependencies.
func NewAppState(
	appCfg *cfg.AppCfg,
	userApp app_iface.UserApp,
	teamApp app_iface.TeamApp,
	memberApp app_iface.MemberApp,
	memberInvApp app_iface.MemberInvApp,
	worksetApp app_iface.WorksetApp,
	comicApp app_iface.ComicApp,
	chapterApp app_iface.ChapterApp,
	chapterPortApp app_iface.ChapterPortApp,
	pageApp app_iface.PageApp,
	unitApp app_iface.UnitApp,
	sysMailApp app_iface.SysMailApp,
	assignmentInvApp app_iface.AssignmentInvApp,
	assignmentApp app_iface.AssignmentApp,
) *AppState {
	return &AppState{
		Cfg:              appCfg,
		UserApp:          userApp,
		TeamApp:          teamApp,
		MemberApp:        memberApp,
		MemberInvApp:     memberInvApp,
		WorksetApp:       worksetApp,
		ComicApp:         comicApp,
		ChapterApp:       chapterApp,
		ChapterPortApp:   chapterPortApp,
		PageApp:          pageApp,
		UnitApp:          unitApp,
		SysMailApp:       sysMailApp,
		AssignmentInvApp: assignmentInvApp,
		AssignmentApp:    assignmentApp,
	}
}
