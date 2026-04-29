package app_impl

import (
	"context"

	app_iface "poprako-s/internal/app"
	"poprako-s/internal/app/res"
	app_util "poprako-s/internal/app/util"
	"poprako-s/internal/app/val"

	"go.uber.org/zap"
)

// `assignmentInvLogAppImpl` is logging decorator for `AssignmentInvApp`.
type assignmentInvLogAppImpl struct {
	inner app_iface.AssignmentInvApp
}

// `NewAssignmentInvLogApp` creates logging decorator for `AssignmentInvApp`.
func NewAssignmentInvLogApp(inner app_iface.AssignmentInvApp) app_iface.AssignmentInvApp {
	if inner == nil {
		zap.L().Panic("[NewAssignmentInvLogApp] nil dependency", zap.Bool("inner", inner == nil))
	}

	return &assignmentInvLogAppImpl{inner: inner}
}

// `ListByChapter` enriches logger context and forwards call.
func (a *assignmentInvLogAppImpl) ListByChapter(cx context.Context, currUid string, args *val.ListAssignmentInvArgs) app_res.AppRes[[]val.AssignmentInvVal] {
	if cx == nil {
		cx = context.Background()
	}

	if args == nil {
		return app_res.Reject[[]val.AssignmentInvVal](app_res.BadRequest, "分页参数不能为空")
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

// `Create` enriches logger context and forwards call.
func (a *assignmentInvLogAppImpl) Create(cx context.Context, currUid string, args *val.CreateAssignmentInvArgs) app_res.AppRes[val.CreateAssignmentInvRes] {
	if cx == nil {
		cx = context.Background()
	}

	if args == nil {
		return app_res.Reject[val.CreateAssignmentInvRes](app_res.BadRequest, "创建参数不能为空")
	}

	lgr := app_util.TakeLgr(cx).With(
		zap.String("curr_uid", currUid),
		zap.String("chapter_id", args.ChapterId),
		zap.String("invitee_qid", args.InviteeQid),
		zap.Uint32("role_mask", uint32(args.RoleMask)),
	)

	cx = app_util.SaveLgr(cx, lgr)

	return a.inner.Create(cx, currUid, args)
}

// `Remove` enriches logger context and forwards call.
func (a *assignmentInvLogAppImpl) Remove(cx context.Context, currUid string, invId string) app_res.AppRes[app_res.None] {
	if cx == nil {
		cx = context.Background()
	}

	lgr := app_util.TakeLgr(cx).With(
		zap.String("curr_uid", currUid),
		zap.String("invitation_id", invId),
	)

	cx = app_util.SaveLgr(cx, lgr)

	return a.inner.Remove(cx, currUid, invId)
}

// `JoinByInvCode` enriches logger context and forwards call.
func (a *assignmentInvLogAppImpl) JoinByInvCode(cx context.Context, currUid string, args *val.JoinAssignmentInvArgs) app_res.AppRes[app_res.None] {
	if cx == nil {
		cx = context.Background()
	}

	if args == nil {
		return app_res.Reject[app_res.None](app_res.BadRequest, "请求参数不能为空")
	}

	lgr := app_util.TakeLgr(cx).With(
		zap.String("curr_uid", currUid),
		zap.String("invitation_code", args.InvCode),
	)

	cx = app_util.SaveLgr(cx, lgr)

	return a.inner.JoinByInvCode(cx, currUid, args)
}
