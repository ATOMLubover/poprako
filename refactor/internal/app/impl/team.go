package app_impl

import (
	"context"

	app_iface "poprako-s/internal/app"
	"poprako-s/internal/app/res"
	app_util "poprako-s/internal/app/util"
	"poprako-s/internal/app/val"
	oss_iface "poprako-s/internal/domain/ext/oss"
	repo_iface "poprako-s/internal/domain/repo"

	"go.uber.org/zap"
)

// `teamAppImpl` is the default implementation of `TeamApp`.
type teamAppImpl struct {
	teamRepo  repo_iface.TeamRepo
	ossSigner oss_iface.Signer
}

// `NewTeamApp` creates an app implementation for team use-cases.
func NewTeamApp(teamRepo repo_iface.TeamRepo, ossSigner oss_iface.Signer) app_iface.TeamApp {
	if teamRepo == nil || ossSigner == nil {
		zap.L().Panic(
			"[NewTeamApp] nil dependency",
			zap.Bool("teamRepo", teamRepo == nil),
			zap.Bool("ossSigner", ossSigner == nil),
		)
	}

	return &teamAppImpl{
		teamRepo:  teamRepo,
		ossSigner: ossSigner,
	}
}

// `GetInfo` gets team information by id.
func (a *teamAppImpl) GetInfo(cx context.Context, id string) app_res.AppRes[val.TeamVal] {
	lgr := app_util.TakeLgr(cx)

	team, err := a.teamRepo.GetById(id)
	if err != nil {
		lgr.Error(
			"[teamAppImpl.GetInfo] failed to get team by id",
			zap.String("team_id", id),
			zap.Error(err),
		)

		return app_res.Reject[val.TeamVal](app_res.BadRequest, "团队不存在")
	}

	teamVal, err := asmTeamVal(team, a.ossSigner)
	if err != nil {
		lgr.Error(
			"[teamAppImpl.GetInfo] failed to assemble team info",
			zap.String("team_id", id),
			zap.Error(err),
		)

		return app_res.Reject[val.TeamVal](app_res.ServerError, "获取团队信息失败")
	}

	return app_res.Accept(teamVal)
}

// `ListByUser` lists teams that a user belongs to.
func (a *teamAppImpl) ListByUser(cx context.Context, userId string) app_res.AppRes[[]val.TeamVal] {
	lgr := app_util.TakeLgr(cx)

	// Keep explicit not-implemented feedback until member/team list query is introduced.
	lgr.Warn("[teamAppImpl.ListUserTeams] not implemented", zap.String("user_id", userId))

	return app_res.Reject[[]val.TeamVal](app_res.ServerError, "团队列表功能暂未实现")
}

// `Update` updates one team profile.
func (a *teamAppImpl) Update(cx context.Context, args *val.TeamUpdArgs) app_res.AppRes[app_res.None] {
	lgr := app_util.TakeLgr(cx)

	// Keep explicit not-implemented feedback until team update repo capability is introduced.
	lgr.Warn("[teamAppImpl.Update] not implemented")

	return app_res.Reject[app_res.None](app_res.ServerError, "更新团队功能暂未实现")
}

// `ResvAvatar` reserves avatar upload for team.
func (a *teamAppImpl) ResvAvatar(cx context.Context, args *val.ResvTeamAvatarArgs) app_res.AppRes[val.ResvTeamAvatarRes] {
	lgr := app_util.TakeLgr(cx)

	// Keep explicit not-implemented feedback until team avatar flow is introduced.
	lgr.Warn("[teamAppImpl.ResvAvatar] not implemented")

	return app_res.Reject[val.ResvTeamAvatarRes](app_res.ServerError, "团队头像预留功能暂未实现")
}

// `MarkAvatarUploaded` marks team avatar upload as completed.
func (a *teamAppImpl) MarkAvatarUploaded(cx context.Context, teamId string) app_res.AppRes[app_res.None] {
	lgr := app_util.TakeLgr(cx)

	// Keep explicit not-implemented feedback until team avatar flow is introduced.
	lgr.Warn("[teamAppImpl.MarkAvatarUploaded] not implemented", zap.String("team_id", teamId))

	return app_res.Reject[app_res.None](app_res.ServerError, "团队头像确认功能暂未实现")
}
