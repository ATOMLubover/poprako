package app_impl

import (
	"context"

	app_iface "poprako-s/internal/app"
	"poprako-s/internal/app/res"
	app_util "poprako-s/internal/app/util"
	"poprako-s/internal/app/val"

	"go.uber.org/zap"
)

// `worksetLogAppImpl` is a logging decorator for `WorksetApp`.
type worksetLogAppImpl struct {
	inner app_iface.WorksetApp
}

// `NewWorksetLogApp` creates the logging decorator for `WorksetApp`.
func NewWorksetLogApp(inner app_iface.WorksetApp) app_iface.WorksetApp {
	if inner == nil {
		zap.L().Panic("[NewWorksetLogApp] nil dependency", zap.Bool("inner", inner == nil))
	}

	return &worksetLogAppImpl{inner: inner}
}

// `List` enriches the logger context then forwards the call.
func (a *worksetLogAppImpl) List(cx context.Context, currUid string, args *val.ListWorksetArgs) res.AppRes[[]val.WorksetVal] {
	if cx == nil {
		cx = context.Background()
	}

	if args == nil {
		return res.Reject[[]val.WorksetVal](res.BadRequest, "分页参数不能为空")
	}

	lgr := app_util.TakeLgr(cx)

	lgr = lgr.With(
		zap.String("curr_uid", currUid),
		zap.String("team_id", args.TeamId),
		zap.Int("offset", args.Offset),
		zap.Int("limit", args.Limit),
	)

	cx = app_util.SaveLgr(cx, lgr)

	return a.inner.List(cx, currUid, args)
}

// `Create` enriches the logger context then forwards the call.
func (a *worksetLogAppImpl) Create(cx context.Context, currUid string, args *val.CreateWorksetArgs) res.AppRes[val.WorksetCreatedRes] {
	if cx == nil {
		cx = context.Background()
	}

	if args == nil {
		return res.Reject[val.WorksetCreatedRes](res.BadRequest, "创建参数不能为空")
	}

	lgr := app_util.TakeLgr(cx)

	lgr = lgr.With(
		zap.String("curr_uid", currUid),
		zap.String("team_id", args.TeamId),
	)

	cx = app_util.SaveLgr(cx, lgr)

	return a.inner.Create(cx, currUid, args)
}

// `Update` enriches the logger context then forwards the call.
func (a *worksetLogAppImpl) Update(cx context.Context, currUid string, args *val.WorksetUpdArgs) res.AppRes[res.None] {
	if cx == nil {
		cx = context.Background()
	}

	if args == nil {
		return res.Reject[res.None](res.BadRequest, "更新参数不能为空")
	}

	lgr := app_util.TakeLgr(cx)

	lgr = lgr.With(
		zap.String("curr_uid", currUid),
		zap.String("workset_id", args.Id),
	)

	cx = app_util.SaveLgr(cx, lgr)

	return a.inner.Update(cx, currUid, args)
}

// `Remove` enriches the logger context then forwards the call.
func (a *worksetLogAppImpl) Remove(cx context.Context, currUid string, worksetId string) res.AppRes[res.None] {
	if cx == nil {
		cx = context.Background()
	}

	lgr := app_util.TakeLgr(cx)

	lgr = lgr.With(
		zap.String("curr_uid", currUid),
		zap.String("workset_id", worksetId),
	)

	cx = app_util.SaveLgr(cx, lgr)

	return a.inner.Remove(cx, currUid, worksetId)
}
