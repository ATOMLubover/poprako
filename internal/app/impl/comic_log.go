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
func (a *comicLogAppImpl) List(cx context.Context, currUid string, args *val.ListComicArgs) app_res.AppRes[[]val.ComicVal] {
	if cx == nil {
		cx = context.Background()
	}

	if args == nil {
		return app_res.Reject[[]val.ComicVal](app_res.BadRequest, "分页参数不能为空")
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
func (a *comicLogAppImpl) Create(cx context.Context, currUid string, args *val.CreateComicArgs) app_res.AppRes[val.ComicCreatedRes] {
	if cx == nil {
		cx = context.Background()
	}

	if args == nil {
		return app_res.Reject[val.ComicCreatedRes](app_res.BadRequest, "创建参数不能为空")
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
func (a *comicLogAppImpl) Update(cx context.Context, currUid string, args *val.ComicUpdArgs) app_res.AppRes[app_res.None] {
	if cx == nil {
		cx = context.Background()
	}

	if args == nil {
		return app_res.Reject[app_res.None](app_res.BadRequest, "更新参数不能为空")
	}

	lgr := app_util.TakeLgr(cx)

	lgr = lgr.With(
		zap.String("curr_uid", currUid),
		zap.String("comic_id", args.Id),
	)

	cx = app_util.SaveLgr(cx, lgr)

	return a.inner.Update(cx, currUid, args)
}

// `GetById` enriches logger context and forwards call.
func (a *comicLogAppImpl) GetById(cx context.Context, currUid string, args *val.GetComicByIdArgs) app_res.AppRes[val.ComicVal] {
	if cx == nil {
		cx = context.Background()
	}

	if args == nil {
		return app_res.Reject[val.ComicVal](app_res.BadRequest, "查询参数不能为空")
	}

	lgr := app_util.TakeLgr(cx)

	lgr = lgr.With(
		zap.String("curr_uid", currUid),
		zap.String("comic_id", args.ComicId),
	)

	cx = app_util.SaveLgr(cx, lgr)

	return a.inner.GetById(cx, currUid, args)
}

// `ResvCover` enriches logger context and forwards call
func (a *comicLogAppImpl) ResvCover(cx context.Context, currUid string, args *val.ResvComicCoverArgs) app_res.AppRes[val.ResvComicCoverRes] {
	if cx == nil {
		cx = context.Background()
	}

	if args == nil {
		return app_res.Reject[val.ResvComicCoverRes](app_res.BadRequest, "预留参数不能为空")
	}

	lgr := app_util.TakeLgr(cx)

	lgr = lgr.With(
		zap.String("curr_uid", currUid),
		zap.String("comic_id", args.ComicId),
		zap.String("file_extension", args.FileExt),
	)

	cx = app_util.SaveLgr(cx, lgr)

	return a.inner.ResvCover(cx, currUid, args)
}

// `MarkCoverUploaded` enriches logger context and forwards call
func (a *comicLogAppImpl) MarkCoverUploaded(cx context.Context, currUid string, comicId string) app_res.AppRes[app_res.None] {
	if cx == nil {
		cx = context.Background()
	}

	lgr := app_util.TakeLgr(cx)

	lgr = lgr.With(
		zap.String("curr_uid", currUid),
		zap.String("comic_id", comicId),
	)

	cx = app_util.SaveLgr(cx, lgr)

	return a.inner.MarkCoverUploaded(cx, currUid, comicId)
}

// `Delete` enriches logger context and forwards call.
func (a *comicLogAppImpl) Delete(cx context.Context, currUid string, comicId string) app_res.AppRes[app_res.None] {
	if cx == nil {
		cx = context.Background()
	}

	lgr := app_util.TakeLgr(cx)

	lgr = lgr.With(
		zap.String("curr_uid", currUid),
		zap.String("comic_id", comicId),
	)

	cx = app_util.SaveLgr(cx, lgr)

	return a.inner.Delete(cx, currUid, comicId)
}
