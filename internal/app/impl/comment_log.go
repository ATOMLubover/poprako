package app_impl

import (
	"context"

	app_iface "poprako-s/internal/app"
	app_res "poprako-s/internal/app/res"
	app_util "poprako-s/internal/app/util"
	"poprako-s/internal/app/val"

	"go.uber.org/zap"
)

// `commentLogAppImpl` is a logging decorator for `CommentApp`.
type commentLogAppImpl struct {
	inner app_iface.CommentApp
}

// `NewCommentLogApp` creates logging decorator for `CommentApp`.
func NewCommentLogApp(inner app_iface.CommentApp) app_iface.CommentApp {
	if inner == nil {
		zap.L().Panic("[NewCommentLogApp] nil dependency", zap.Bool("inner", inner == nil))
	}

	return &commentLogAppImpl{inner: inner}
}

// `List` enriches logger context and forwards call.
func (a *commentLogAppImpl) List(cx context.Context, currUid string, args *val.ListCommentArgs) app_res.AppRes[[]val.CommentVal] {
	if cx == nil {
		cx = context.Background()
	}

	if args == nil {
		return app_res.Reject[[]val.CommentVal](app_res.BadRequest, "分页参数不能为空")
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
func (a *commentLogAppImpl) Create(cx context.Context, currUid string, args *val.CreateCommentArgs) app_res.AppRes[val.CommentCreatedRes] {
	if cx == nil {
		cx = context.Background()
	}

	if args == nil {
		return app_res.Reject[val.CommentCreatedRes](app_res.BadRequest, "创建参数不能为空")
	}

	lgr := app_util.TakeLgr(cx)
	lgr = lgr.With(
		zap.String("curr_uid", currUid),
		zap.String("team_id", args.TeamId),
	)

	cx = app_util.SaveLgr(cx, lgr)

	return a.inner.Create(cx, currUid, args)
}
