package app_impl

import (
	"context"

	app_iface "poprako-s/internal/app"
	"poprako-s/internal/app/res"
	app_util "poprako-s/internal/app/util"
	"poprako-s/internal/app/val"

	"go.uber.org/zap"
)

// `sysMailLogAppImpl` is a logging decorator for `SysMailApp`.
type sysMailLogAppImpl struct {
	inner app_iface.SysMailApp
}

// `NewSysMailLogApp` creates logging decorator for `SysMailApp`.
func NewSysMailLogApp(inner app_iface.SysMailApp) app_iface.SysMailApp {
	if inner == nil {
		zap.L().Panic("[NewSysMailLogApp] nil dependency", zap.Bool("inner", inner == nil))
	}

	return &sysMailLogAppImpl{inner: inner}
}

// `List` enriches logger context and forwards call.
func (a *sysMailLogAppImpl) List(cx context.Context, currUid string, args *val.ListSysMailArgs) app_res.AppRes[[]val.SysMailVal] {
	if cx == nil {
		cx = context.Background()
	}

	if args == nil {
		return app_res.Reject[[]val.SysMailVal](app_res.BadRequest, "分页参数不能为空")
	}

	lgr := app_util.TakeLgr(cx)
	lgr = lgr.With(
		zap.String("curr_uid", currUid),
		zap.Int("offset", args.Offset),
		zap.Int("limit", args.Limit),
	)

	cx = app_util.SaveLgr(cx, lgr)

	return a.inner.List(cx, currUid, args)
}

// `MarkRead` enriches logger context and forwards call.
func (a *sysMailLogAppImpl) MarkRead(cx context.Context, currUid string, id string) app_res.AppRes[app_res.None] {
	if cx == nil {
		cx = context.Background()
	}

	lgr := app_util.TakeLgr(cx)
	lgr = lgr.With(
		zap.String("curr_uid", currUid),
		zap.String("sys_mail_id", id),
	)

	cx = app_util.SaveLgr(cx, lgr)

	return a.inner.MarkRead(cx, currUid, id)
}
