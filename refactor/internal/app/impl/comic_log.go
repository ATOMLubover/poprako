package app_impl

import (
	"context"

	app_iface "poprako-s/internal/app"
	"poprako-s/internal/app/res"
	app_util "poprako-s/internal/app/util"
	"poprako-s/internal/app/val"

	"go.uber.org/zap"
)

// `comicLogAppImpl` is a logging decorator for `ComicApp`
type comicLogAppImpl struct {
	inner app_iface.ComicApp
}

// `NewComicLogApp` creates a logging decorator for `ComicApp`
func NewComicLogApp(inner app_iface.ComicApp) app_iface.ComicApp {
	if inner == nil {
		zap.L().Panic("[NewComicLogApp] nil dependency", zap.Bool("inner", inner == nil))
	}

	return &comicLogAppImpl{inner: inner}
}

// `List` enriches logger context and forwards call
func (a *comicLogAppImpl) List(cx context.Context, currUid string, args *val.ListComicArgs) res.AppRes[[]val.ComicVal] {
	if cx == nil {
		cx = context.Background()
	}

	if args == nil {
		return res.Reject[[]val.ComicVal](res.BadRequest, "分页参数不能为空")
	}

	lgr := app_util.TakeLgr(cx)

	lgr = lgr.With(
		zap.String("curr_uid", currUid),
		zap.String("workset_id", args.WorksetId),
		zap.String("fuzzy_title", args.FuzzyTitle),
		zap.Int("offset", args.Offset),
		zap.Int("limit", args.Limit),
	)

	cx = app_util.SaveLgr(cx, lgr)

	return a.inner.List(cx, currUid, args)
}

// `Create` enriches logger context and forwards call
func (a *comicLogAppImpl) Create(cx context.Context, currUid string, args *val.CreateComicArgs) res.AppRes[val.ComicCreatedRes] {
	if cx == nil {
		cx = context.Background()
	}

	if args == nil {
		return res.Reject[val.ComicCreatedRes](res.BadRequest, "创建参数不能为空")
	}

	lgr := app_util.TakeLgr(cx)

	lgr = lgr.With(
		zap.String("curr_uid", currUid),
		zap.String("workset_id", args.WorksetId),
		zap.String("title", args.Title),
		zap.String("author", args.Author),
	)

	cx = app_util.SaveLgr(cx, lgr)

	return a.inner.Create(cx, currUid, args)
}

// `Update` enriches logger context and forwards call
func (a *comicLogAppImpl) Update(cx context.Context, currUid string, args *val.ComicUpdArgs) res.AppRes[res.None] {
	if cx == nil {
		cx = context.Background()
	}

	if args == nil {
		return res.Reject[res.None](res.BadRequest, "更新参数不能为空")
	}

	lgr := app_util.TakeLgr(cx)

	lgr = lgr.With(
		zap.String("curr_uid", currUid),
		zap.String("comic_id", args.Id),
	)

	cx = app_util.SaveLgr(cx, lgr)

	return a.inner.Update(cx, currUid, args)
}

// `Remove` enriches logger context and forwards call
func (a *comicLogAppImpl) Remove(cx context.Context, currUid string, comicId string) res.AppRes[res.None] {
	if cx == nil {
		cx = context.Background()
	}

	lgr := app_util.TakeLgr(cx)

	lgr = lgr.With(
		zap.String("curr_uid", currUid),
		zap.String("comic_id", comicId),
	)

	cx = app_util.SaveLgr(cx, lgr)

	return a.inner.Remove(cx, currUid, comicId)
}
