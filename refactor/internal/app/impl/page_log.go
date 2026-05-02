package app_impl

import (
	"context"

	app_iface "poprako-s/internal/app"
	app_res "poprako-s/internal/app/res"
	app_util "poprako-s/internal/app/util"
	"poprako-s/internal/app/val"

	"go.uber.org/zap"
)

// `pageLogAppImpl` is a logging decorator for `PageApp`.
type pageLogAppImpl struct {
	// `inner` is the wrapped page application.
	inner app_iface.PageApp
}

// `NewPageLogApp` creates the logging decorator for `PageApp`.
func NewPageLogApp(inner app_iface.PageApp) app_iface.PageApp {
	if inner == nil {
		zap.L().Panic("[NewPageLogApp] nil dependency", zap.Bool("inner", inner == nil))
	}

	return &pageLogAppImpl{inner: inner}
}

// `ResvChapterPages` enriches logger context then forwards the call.
func (a *pageLogAppImpl) ResvChapterPages(cx context.Context, currUid string, args *val.ResvChapterPagesArgs) app_res.AppRes[val.ResvChapterPagesRes] {
	if cx == nil {
		cx = context.Background()
	}

	if args == nil {
		return app_res.Reject[val.ResvChapterPagesRes](app_res.BadRequest, "预留参数不能为空")
	}

	lgr := app_util.TakeLgr(cx)
	lgr = lgr.With(zap.String("curr_uid", currUid), zap.String("chapter_id", args.ChapterId), zap.Int("page_count", args.PageCount))
	cx = app_util.SaveLgr(cx, lgr)

	return a.inner.ResvChapterPages(cx, currUid, args)
}

// `List` enriches logger context then forwards the call.
func (a *pageLogAppImpl) List(cx context.Context, currUid string, args *val.ListChapterPageArgs) app_res.AppRes[[]val.PageVal] {
	if cx == nil {
		cx = context.Background()
	}

	if args == nil {
		return app_res.Reject[[]val.PageVal](app_res.BadRequest, "分页参数不能为空")
	}

	lgr := app_util.TakeLgr(cx)
	lgr = lgr.With(zap.String("curr_uid", currUid), zap.String("chapter_id", args.ChapterId), zap.Int("offset", args.Offset), zap.Int("limit", args.Limit))
	cx = app_util.SaveLgr(cx, lgr)

	return a.inner.List(cx, currUid, args)
}

// `MarkImageUploaded` enriches logger context then forwards the call.
func (a *pageLogAppImpl) MarkImageUploaded(cx context.Context, currUid string, args *val.MarkPageImageUploadedArgs) app_res.AppRes[app_res.None] {
	if cx == nil {
		cx = context.Background()
	}

	if args == nil {
		return app_res.Reject[app_res.None](app_res.BadRequest, "上传确认参数不能为空")
	}

	lgr := app_util.TakeLgr(cx)
	lgr = lgr.With(zap.String("curr_uid", currUid), zap.String("page_id", args.PageId))
	cx = app_util.SaveLgr(cx, lgr)

	return a.inner.MarkImageUploaded(cx, currUid, args)
}

// `DeleteByChapterId` enriches logger context then forwards the call.
func (a *pageLogAppImpl) DeleteByChapterId(cx context.Context, currUid string, chapterId string) app_res.AppRes[app_res.None] {
	if cx == nil {
		cx = context.Background()
	}

	lgr := app_util.TakeLgr(cx)
	lgr = lgr.With(zap.String("curr_uid", currUid), zap.String("chapter_id", chapterId))
	cx = app_util.SaveLgr(cx, lgr)

	return a.inner.DeleteByChapterId(cx, currUid, chapterId)
}
