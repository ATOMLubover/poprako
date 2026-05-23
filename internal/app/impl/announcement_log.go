package app_impl

import (
	"context"

	app_iface "poprako-s/internal/app"
	app_res "poprako-s/internal/app/res"
	app_util "poprako-s/internal/app/util"
	"poprako-s/internal/app/val"

	"go.uber.org/zap"
)

// `announcementLogAppImpl` is a logging decorator for `AnnouncementApp`.
type announcementLogAppImpl struct {
	inner app_iface.AnnouncementApp
}

// `NewAnnouncementLogApp` creates logging decorator for `AnnouncementApp`.
func NewAnnouncementLogApp(inner app_iface.AnnouncementApp) app_iface.AnnouncementApp {
	if inner == nil {
		zap.L().Panic("[NewAnnouncementLogApp] nil dependency", zap.Bool("inner", inner == nil))
	}

	return &announcementLogAppImpl{inner: inner}
}

// `List` enriches logger context and forwards call.
func (a *announcementLogAppImpl) List(cx context.Context, currUid string, args *val.ListAnnouncementArgs) app_res.AppRes[[]val.AnnouncementVal] {
	if cx == nil {
		cx = context.Background()
	}

	if args == nil {
		return app_res.Reject[[]val.AnnouncementVal](app_res.BadRequest, "分页参数不能为空")
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

// `Create` enriches logger context and forwards call.
func (a *announcementLogAppImpl) Create(cx context.Context, currUid string, args *val.CreateAnnouncementArgs) app_res.AppRes[val.AnnouncementCreatedRes] {
	if cx == nil {
		cx = context.Background()
	}

	if args == nil {
		return app_res.Reject[val.AnnouncementCreatedRes](app_res.BadRequest, "创建参数不能为空")
	}

	lgr := app_util.TakeLgr(cx)
	lgr = lgr.With(
		zap.String("curr_uid", currUid),
		zap.String("team_id", args.TeamId),
	)

	cx = app_util.SaveLgr(cx, lgr)

	return a.inner.Create(cx, currUid, args)
}
