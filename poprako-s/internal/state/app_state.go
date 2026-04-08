package state

import (
	"poprako-s/internal/app"
	"poprako-s/internal/cfg"
)

// AppState 聚合了应用层的全部依赖，用于向 HTTP 处理器注入服务
type AppState struct {
	Cfg *cfg.AppCfg

	UserApp       app.UserApp
	TeamApp       app.TeamApp
	MemberApp     app.MemberApp
	InvitationApp app.MemberInvitationApp
	WorksetApp    app.WorksetApp
	ComicApp      app.ComicApp
	ChapterApp    app.ChapterApp
	PageApp       app.PageApp
	AssignmentApp app.AssignmentApp
	UnitApp       app.UnitApp
}

// NewAppState 创建一个新的 AppState 实例
func NewAppState(
	cfg *cfg.AppCfg,
	userApp app.UserApp,
	teamApp app.TeamApp,
	memberApp app.MemberApp,
	invitationApp app.MemberInvitationApp,
	worksetApp app.WorksetApp,
	comicApp app.ComicApp,
	chapterApp app.ChapterApp,
	pageApp app.PageApp,
	assignmentApp app.AssignmentApp,
	unitApp app.UnitApp,
) *AppState {
	return &AppState{
		Cfg:           cfg,
		UserApp:       userApp,
		TeamApp:       teamApp,
		MemberApp:     memberApp,
		InvitationApp: invitationApp,
		WorksetApp:    worksetApp,
		ComicApp:      comicApp,
		ChapterApp:    chapterApp,
		PageApp:       pageApp,
		AssignmentApp: assignmentApp,
		UnitApp:       unitApp,
	}
}
