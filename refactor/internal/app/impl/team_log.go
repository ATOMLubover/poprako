package app_impl

import (
	"context"

	app_iface "poprako-s/internal/app"
	"poprako-s/internal/app/res"
	app_util "poprako-s/internal/app/util"
	"poprako-s/internal/app/val"

	"go.uber.org/zap"
)

// `teamLogAppImpl` is a logging decorator for `TeamApp`.
type teamLogAppImpl struct {
	inner app_iface.TeamApp
}

// `NewTeamLogApp` creates the logging decorator for `TeamApp`.
func NewTeamLogApp(inner app_iface.TeamApp) app_iface.TeamApp {
	if inner == nil {
		zap.L().Panic("[NewTeamLogApp] nil dependency", zap.Bool("inner", inner == nil))
	}

	return &teamLogAppImpl{inner: inner}
}

// `GetInfo` enriches logger context then forwards the call.
func (a *teamLogAppImpl) GetInfo(cx context.Context, id string) app_res.AppRes[val.TeamVal] {
	if cx == nil {
		cx = context.Background()
	}

	lgr := app_util.TakeLgr(cx)
	if lgr == nil {
		lgr = zap.L()
	}

	lgr = lgr.With(zap.String("team_id", id))

	cx = app_util.SaveLgr(cx, lgr)

	return a.inner.GetInfo(cx, id)
}

// `ListByUser` enriches logger context then forwards the call.
func (a *teamLogAppImpl) ListByUser(cx context.Context, userId string) app_res.AppRes[[]val.TeamVal] {
	if cx == nil {
		cx = context.Background()
	}

	lgr := app_util.TakeLgr(cx)
	if lgr == nil {
		lgr = zap.L()
	}

	lgr = lgr.With(zap.String("user_id", userId))

	cx = app_util.SaveLgr(cx, lgr)

	return a.inner.ListByUser(cx, userId)
}

// `Update` enriches logger context then forwards the call.
func (a *teamLogAppImpl) Update(cx context.Context, args *val.TeamUpdArgs) app_res.AppRes[app_res.None] {
	if cx == nil {
		cx = context.Background()
	}

	if args == nil {
		return app_res.Reject[app_res.None](app_res.ServerError, "更新团队功能暂未实现")
	}

	lgr := app_util.TakeLgr(cx)
	if lgr == nil {
		lgr = zap.L()
	}

	lgr = lgr.With(zap.String("team_id", args.Id))

	cx = app_util.SaveLgr(cx, lgr)

	return a.inner.Update(cx, args)
}

// `ResvAvatar` enriches logger context then forwards the call.
func (a *teamLogAppImpl) ResvAvatar(cx context.Context, args *val.ResvTeamAvatarArgs) app_res.AppRes[val.ResvTeamAvatarRes] {
	if cx == nil {
		cx = context.Background()
	}

	if args == nil {
		return app_res.Reject[val.ResvTeamAvatarRes](app_res.ServerError, "团队头像预留功能暂未实现")
	}

	lgr := app_util.TakeLgr(cx)
	if lgr == nil {
		lgr = zap.L()
	}

	lgr = lgr.With(zap.String("team_id", args.TeamId), zap.String("file_ext", args.FileExt))

	cx = app_util.SaveLgr(cx, lgr)

	return a.inner.ResvAvatar(cx, args)
}

// `MarkAvatarUploaded` enriches logger context then forwards the call.
func (a *teamLogAppImpl) MarkAvatarUploaded(cx context.Context, teamId string) app_res.AppRes[app_res.None] {
	if cx == nil {
		cx = context.Background()
	}

	lgr := app_util.TakeLgr(cx)
	if lgr == nil {
		lgr = zap.L()
	}

	lgr = lgr.With(zap.String("team_id", teamId))

	cx = app_util.SaveLgr(cx, lgr)

	return a.inner.MarkAvatarUploaded(cx, teamId)
}
