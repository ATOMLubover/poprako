package app_impl

import (
	"context"

	app_iface "poprako-s/internal/app"
	app_res "poprako-s/internal/app/res"
	app_util "poprako-s/internal/app/util"
	"poprako-s/internal/app/val"

	"go.uber.org/zap"
)

// `assignmentLogAppImpl` is logging decorator for `AssignmentApp`.
type assignmentLogAppImpl struct {
	inner app_iface.AssignmentApp
}

// `NewAssignmentLogApp` creates logging decorator for `AssignmentApp`.
func NewAssignmentLogApp(inner app_iface.AssignmentApp) app_iface.AssignmentApp {
	if inner == nil {
		zap.L().Panic("[NewAssignmentLogApp] nil dependency", zap.Bool("inner", inner == nil))
	}

	return &assignmentLogAppImpl{inner: inner}
}

// `ListByChapter` enriches logger context and forwards call.
func (a *assignmentLogAppImpl) ListByChapter(cx context.Context, currUid string, args *val.ListAssignmentByChapterArgs) app_res.AppRes[[]val.AssignmentVal] {
	if cx == nil {
		cx = context.Background()
	}

	if args == nil {
		return app_res.Reject[[]val.AssignmentVal](app_res.BadRequest, "分页参数不能为空")
	}

	lgr := app_util.TakeLgr(cx).With(
		zap.String("curr_uid", currUid),
		zap.String("chapter_id", args.ChapterId),
		zap.Int("offset", args.Offset),
		zap.Int("limit", args.Limit),
	)

	cx = app_util.SaveLgr(cx, lgr)

	return a.inner.ListByChapter(cx, currUid, args)
}

// `ListByUser` enriches logger context and forwards call.
func (a *assignmentLogAppImpl) ListByUser(cx context.Context, currUid string, args *val.ListAssignmentByUserArgs) app_res.AppRes[[]val.AssignmentVal] {
	if cx == nil {
		cx = context.Background()
	}

	if args == nil {
		return app_res.Reject[[]val.AssignmentVal](app_res.BadRequest, "分页参数不能为空")
	}

	lgr := app_util.TakeLgr(cx).With(
		zap.String("curr_uid", currUid),
		zap.Int("offset", args.Offset),
		zap.Int("limit", args.Limit),
	)

	cx = app_util.SaveLgr(cx, lgr)

	return a.inner.ListByUser(cx, currUid, args)
}

// `Upsert` enriches logger context and forwards call.
func (a *assignmentLogAppImpl) Upsert(cx context.Context, currUid string, args *val.UpsertAssignmentArgs) app_res.AppRes[app_res.None] {
	if cx == nil {
		cx = context.Background()
	}

	if args == nil {
		return app_res.Reject[app_res.None](app_res.BadRequest, "更新参数不能为空")
	}

	lgr := app_util.TakeLgr(cx).With(
		zap.String("curr_uid", currUid),
		zap.String("chapter_id", args.ChapterId),
		zap.String("user_id", args.UserId),
		zap.Uint32("role_mask", uint32(args.RoleMask)),
	)

	cx = app_util.SaveLgr(cx, lgr)

	return a.inner.Upsert(cx, currUid, args)
}

// `Delete` enriches logger context and forwards call.
func (a *assignmentLogAppImpl) Delete(cx context.Context, currUid string, assignmentId string) app_res.AppRes[app_res.None] {
	if cx == nil {
		cx = context.Background()
	}

	lgr := app_util.TakeLgr(cx).With(
		zap.String("curr_uid", currUid),
		zap.String("assignment_id", assignmentId),
	)

	cx = app_util.SaveLgr(cx, lgr)

	return a.inner.Delete(cx, currUid, assignmentId)
}
