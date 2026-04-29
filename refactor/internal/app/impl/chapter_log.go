package app_impl

import (
	"context"

	app_iface "poprako-s/internal/app"
	"poprako-s/internal/app/res"
	app_util "poprako-s/internal/app/util"
	"poprako-s/internal/app/val"

	"go.uber.org/zap"
)

// `chapterLogAppImpl` is a logging decorator for `ChapterApp`.
type chapterLogAppImpl struct {
	inner app_iface.ChapterApp
}

// `NewChapterLogApp` creates logging decorator for `ChapterApp`.
func NewChapterLogApp(inner app_iface.ChapterApp) app_iface.ChapterApp {
	if inner == nil {
		zap.L().Panic("[NewChapterLogApp] nil dependency", zap.Bool("inner", inner == nil))
	}

	return &chapterLogAppImpl{inner: inner}
}

// `List` enriches logger context and forwards call.
func (a *chapterLogAppImpl) List(cx context.Context, currUid string, args *val.ListChapterArgs) app_res.AppRes[[]val.ChapterVal] {
	if cx == nil {
		cx = context.Background()
	}

	if args == nil {
		return app_res.Reject[[]val.ChapterVal](app_res.BadRequest, "分页参数不能为空")
	}

	lgr := app_util.TakeLgr(cx)
	lgr = lgr.With(
		zap.String("curr_uid", currUid),
		zap.String("comic_id", args.ComicId),
		zap.Int("offset", args.Offset),
		zap.Int("limit", args.Limit),
	)

	cx = app_util.SaveLgr(cx, lgr)

	return a.inner.List(cx, currUid, args)
}

// `GetById` enriches logger context and forwards call.
func (a *chapterLogAppImpl) GetById(cx context.Context, currUid string, chapterId string) app_res.AppRes[val.ChapterVal] {
	if cx == nil {
		cx = context.Background()
	}

	lgr := app_util.TakeLgr(cx)
	lgr = lgr.With(
		zap.String("curr_uid", currUid),
		zap.String("chapter_id", chapterId),
	)

	cx = app_util.SaveLgr(cx, lgr)

	return a.inner.GetById(cx, currUid, chapterId)
}

// `GetPinned` enriches logger context and forwards call.
func (a *chapterLogAppImpl) GetPinned(cx context.Context, currUid string, comicId string) app_res.AppRes[val.ChapterVal] {
	if cx == nil {
		cx = context.Background()
	}

	lgr := app_util.TakeLgr(cx)
	lgr = lgr.With(
		zap.String("curr_uid", currUid),
		zap.String("comic_id", comicId),
	)

	cx = app_util.SaveLgr(cx, lgr)

	return a.inner.GetPinned(cx, currUid, comicId)
}

// `Create` enriches logger context and forwards call.
func (a *chapterLogAppImpl) Create(cx context.Context, currUid string, args *val.CreateChapterArgs) app_res.AppRes[val.ChapterCreatedRes] {
	if cx == nil {
		cx = context.Background()
	}

	if args == nil {
		return app_res.Reject[val.ChapterCreatedRes](app_res.BadRequest, "创建参数不能为空")
	}

	lgr := app_util.TakeLgr(cx)
	lgr = lgr.With(
		zap.String("curr_uid", currUid),
		zap.String("comic_id", args.ComicId),
	)

	cx = app_util.SaveLgr(cx, lgr)

	return a.inner.Create(cx, currUid, args)
}

// `Update` enriches logger context and forwards call.
func (a *chapterLogAppImpl) Update(cx context.Context, currUid string, args *val.ChapterUpdArgs) app_res.AppRes[app_res.None] {
	if cx == nil {
		cx = context.Background()
	}

	if args == nil {
		return app_res.Reject[app_res.None](app_res.BadRequest, "更新参数不能为空")
	}

	lgr := app_util.TakeLgr(cx)
	lgr = lgr.With(
		zap.String("curr_uid", currUid),
		zap.String("chapter_id", args.Id),
	)

	cx = app_util.SaveLgr(cx, lgr)

	return a.inner.Update(cx, currUid, args)
}

// `Remove` enriches logger context and forwards call.
func (a *chapterLogAppImpl) Remove(cx context.Context, currUid string, chapterId string) app_res.AppRes[app_res.None] {
	if cx == nil {
		cx = context.Background()
	}

	lgr := app_util.TakeLgr(cx)
	lgr = lgr.With(
		zap.String("curr_uid", currUid),
		zap.String("chapter_id", chapterId),
	)

	cx = app_util.SaveLgr(cx, lgr)

	return a.inner.Remove(cx, currUid, chapterId)
}
