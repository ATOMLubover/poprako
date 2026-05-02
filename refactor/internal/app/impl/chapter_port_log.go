package app_impl

import (
	"context"

	app_iface "poprako-s/internal/app"
	app_res "poprako-s/internal/app/res"
	app_util "poprako-s/internal/app/util"
	"poprako-s/internal/app/val"

	"go.uber.org/zap"
)

// `chapterPortLogAppImpl` is a logging decorator for `ChapterPortApp`.
type chapterPortLogAppImpl struct {
	inner app_iface.ChapterPortApp
}

// `NewChapterPortLogApp` creates logging decorator for `ChapterPortApp`.
func NewChapterPortLogApp(inner app_iface.ChapterPortApp) app_iface.ChapterPortApp {
	if inner == nil {
		zap.L().Panic("[NewChapterPortLogApp] nil dependency", zap.Bool("inner", inner == nil))
	}

	return &chapterPortLogAppImpl{inner: inner}
}

// `Export` enriches logger context and forwards call.
func (a *chapterPortLogAppImpl) Export(cx context.Context, currUid string, chapterId string) app_res.AppRes[val.ChapterExportVal] {
	if cx == nil {
		cx = context.Background()
	}

	lgr := app_util.TakeLgr(cx)
	lgr = lgr.With(zap.String("curr_uid", currUid), zap.String("chapter_id", chapterId))

	cx = app_util.SaveLgr(cx, lgr)

	return a.inner.Export(cx, currUid, chapterId)
}

// `ExportLp` enriches logger context and forwards call.
func (a *chapterPortLogAppImpl) ExportLp(cx context.Context, currUid string, chapterId string) app_res.AppRes[string] {
	if cx == nil {
		cx = context.Background()
	}

	lgr := app_util.TakeLgr(cx)
	lgr = lgr.With(zap.String("curr_uid", currUid), zap.String("chapter_id", chapterId))

	cx = app_util.SaveLgr(cx, lgr)

	return a.inner.ExportLp(cx, currUid, chapterId)
}

// `Import` enriches logger context and forwards call.
func (a *chapterPortLogAppImpl) Import(cx context.Context, currUid string, args *val.ImportChapterArgs) app_res.AppRes[val.ImportChapterRes] {
	if cx == nil {
		cx = context.Background()
	}

	if args == nil {
		return app_res.Reject[val.ImportChapterRes](app_res.BadRequest, "导入参数不能为空")
	}

	lgr := app_util.TakeLgr(cx)
	lgr = lgr.With(
		zap.String("curr_uid", currUid),
		zap.String("chapter_id", args.ChapterId),
		zap.String("format", args.Format),
	)

	cx = app_util.SaveLgr(cx, lgr)

	return a.inner.Import(cx, currUid, args)
}
