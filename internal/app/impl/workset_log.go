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
func (a *worksetLogAppImpl) List(cx context.Context, currUid string, args *val.ListWorksetArgs) app_res.AppRes[[]val.WorksetVal] {
	if cx == nil {
		cx = context.Background()
	}

	if args == nil {
		return app_res.Reject[[]val.WorksetVal](app_res.BadRequest, "分页参数不能为空")
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
func (a *worksetLogAppImpl) Create(cx context.Context, currUid string, args *val.CreateWorksetArgs) app_res.AppRes[val.WorksetCreatedRes] {
	if cx == nil {
		cx = context.Background()
	}

	if args == nil {
		return app_res.Reject[val.WorksetCreatedRes](app_res.BadRequest, "创建参数不能为空")
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
func (a *worksetLogAppImpl) Update(cx context.Context, currUid string, args *val.WorksetUpdArgs) app_res.AppRes[app_res.None] {
	if cx == nil {
		cx = context.Background()
	}

	if args == nil {
		return app_res.Reject[app_res.None](app_res.BadRequest, "更新参数不能为空")
	}

	lgr := app_util.TakeLgr(cx)

	lgr = lgr.With(
		zap.String("curr_uid", currUid),
		zap.String("workset_id", args.Id),
	)

	cx = app_util.SaveLgr(cx, lgr)

	return a.inner.Update(cx, currUid, args)
}

// `Delete` enriches the logger context then forwards the call.
func (a *worksetLogAppImpl) Delete(cx context.Context, currUid string, worksetId string) app_res.AppRes[app_res.None] {
	if cx == nil {
		cx = context.Background()
	}

	lgr := app_util.TakeLgr(cx)

	lgr = lgr.With(
		zap.String("curr_uid", currUid),
		zap.String("workset_id", worksetId),
	)

	cx = app_util.SaveLgr(cx, lgr)

	return a.inner.Delete(cx, currUid, worksetId)
}
