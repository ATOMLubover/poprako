package app_impl

import (
	"context"

	app_iface "poprako-s/internal/app"
	app_res "poprako-s/internal/app/res"
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

// `Create` enriches logger context then forwards the call.
func (a *teamLogAppImpl) Create(cx context.Context, currUid string, args *val.TeamCreArgs) app_res.AppRes[val.TeamCreRes] {
	if cx == nil {
		cx = context.Background()
	}

	if args == nil {
		return app_res.Reject[val.TeamCreRes](app_res.BadRequest, "创建参数不能为空")
	}

	lgr := app_util.TakeLgr(cx)
	if lgr == nil {
		lgr = zap.L()
	}

	lgr = lgr.With(zap.String("curr_uid", currUid), zap.String("name", args.Name))

	cx = app_util.SaveLgr(cx, lgr)

	return a.inner.Create(cx, currUid, args)
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

// `List` enriches logger context then forwards the call.
func (a *teamLogAppImpl) List(cx context.Context, currUid string, args *val.ListTeamArgs) app_res.AppRes[[]val.TeamVal] {
	if cx == nil {
		cx = context.Background()
	}

	if args == nil {
		return app_res.Reject[[]val.TeamVal](app_res.BadRequest, "分页参数不能为空")
	}

	lgr := app_util.TakeLgr(cx)
	if lgr == nil {
		lgr = zap.L()
	}

	lgr = lgr.With(zap.String("curr_uid", currUid), zap.Int("offset", args.Offset), zap.Int("limit", args.Limit))

	cx = app_util.SaveLgr(cx, lgr)

	return a.inner.List(cx, currUid, args)
}

// `ListByUser` enriches logger context then forwards the call.
func (a *teamLogAppImpl) ListByUser(cx context.Context, userId string, args *val.ListTeamArgs) app_res.AppRes[[]val.TeamVal] {
	if cx == nil {
		cx = context.Background()
	}

	if args == nil {
		return app_res.Reject[[]val.TeamVal](app_res.BadRequest, "分页参数不能为空")
	}

	lgr := app_util.TakeLgr(cx)
	if lgr == nil {
		lgr = zap.L()
	}

	lgr = lgr.With(zap.String("user_id", userId), zap.Int("offset", args.Offset), zap.Int("limit", args.Limit))

	cx = app_util.SaveLgr(cx, lgr)

	return a.inner.ListByUser(cx, userId, args)
}

// `Update` enriches logger context then forwards the call.
func (a *teamLogAppImpl) Update(cx context.Context, currUid string, args *val.TeamUpdArgs) app_res.AppRes[app_res.None] {
	if cx == nil {
		cx = context.Background()
	}

	if args == nil {
		return app_res.Reject[app_res.None](app_res.BadRequest, "更新参数不能为空")
	}

	lgr := app_util.TakeLgr(cx)
	if lgr == nil {
		lgr = zap.L()
	}

	lgr = lgr.With(zap.String("curr_uid", currUid), zap.String("team_id", args.Id))

	cx = app_util.SaveLgr(cx, lgr)

	return a.inner.Update(cx, currUid, args)
}

// `ResvAvatar` enriches logger context then forwards the call.
func (a *teamLogAppImpl) ResvAvatar(cx context.Context, currUid string, args *val.ResvTeamAvatarArgs) app_res.AppRes[val.ResvTeamAvatarRes] {
	if cx == nil {
		cx = context.Background()
	}

	if args == nil {
		return app_res.Reject[val.ResvTeamAvatarRes](app_res.BadRequest, "预留参数不能为空")
	}

	lgr := app_util.TakeLgr(cx)
	if lgr == nil {
		lgr = zap.L()
	}

	lgr = lgr.With(zap.String("curr_uid", currUid), zap.String("team_id", args.TeamId), zap.String("file_ext", args.FileExt))

	cx = app_util.SaveLgr(cx, lgr)

	return a.inner.ResvAvatar(cx, currUid, args)
}

// `MarkAvatarUploaded` enriches logger context then forwards the call.
func (a *teamLogAppImpl) MarkAvatarUploaded(cx context.Context, currUid string, teamId string) app_res.AppRes[app_res.None] {
	if cx == nil {
		cx = context.Background()
	}

	lgr := app_util.TakeLgr(cx)
	if lgr == nil {
		lgr = zap.L()
	}

	lgr = lgr.With(zap.String("curr_uid", currUid), zap.String("team_id", teamId))

	cx = app_util.SaveLgr(cx, lgr)

	return a.inner.MarkAvatarUploaded(cx, currUid, teamId)
}
