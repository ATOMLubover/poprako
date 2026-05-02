package app_impl

import (
	"context"

	app_iface "poprako-s/internal/app"
	app_res "poprako-s/internal/app/res"
	app_util "poprako-s/internal/app/util"
	"poprako-s/internal/app/val"

	"go.uber.org/zap"
)

// `memberInvLogAppImpl` is logging decorator for `MemberInvApp`.
type memberInvLogAppImpl struct {
	inner app_iface.MemberInvApp
}

// `NewMemberInvLogApp` creates logging decorator for `MemberInvApp`.
func NewMemberInvLogApp(inner app_iface.MemberInvApp) app_iface.MemberInvApp {
	if inner == nil {
		zap.L().Panic("[NewMemberInvLogApp] nil dependency", zap.Bool("inner", inner == nil))
	}

	return &memberInvLogAppImpl{inner: inner}
}

// `List` enriches logger context and forwards call.
func (a *memberInvLogAppImpl) List(cx context.Context, currUid string, args *val.ListMemberInvArgs) app_res.AppRes[[]val.MemberInvVal] {
	if cx == nil {
		cx = context.Background()
	}

	if args == nil {
		return app_res.Reject[[]val.MemberInvVal](app_res.BadRequest, "分页参数不能为空")
	}

	lgr := app_util.TakeLgr(cx).With(zap.String("curr_uid", currUid), zap.String("team_id", args.TeamId), zap.Int("offset", args.Offset), zap.Int("limit", args.Limit))

	cx = app_util.SaveLgr(cx, lgr)

	return a.inner.List(cx, currUid, args)
}

// `Create` enriches logger context and forwards call.
func (a *memberInvLogAppImpl) Create(cx context.Context, currUid string, args *val.CreateMemberInvArgs) app_res.AppRes[val.CreateMemberInvRes] {
	if cx == nil {
		cx = context.Background()
	}

	if args == nil {
		return app_res.Reject[val.CreateMemberInvRes](app_res.BadRequest, "创建参数不能为空")
	}

	lgr := app_util.TakeLgr(cx).With(zap.String("curr_uid", currUid), zap.String("team_id", args.TeamId), zap.String("invitee_qid", args.InviteeQid), zap.Uint32("role_mask", uint32(args.RoleMask)))

	cx = app_util.SaveLgr(cx, lgr)

	return a.inner.Create(cx, currUid, args)
}

// `Update` enriches logger context and forwards call.
func (a *memberInvLogAppImpl) Update(cx context.Context, currUid string, args *val.MemberInvUpdArgs) app_res.AppRes[app_res.None] {
	if cx == nil {
		cx = context.Background()
	}

	if args == nil {
		return app_res.Reject[app_res.None](app_res.BadRequest, "更新参数不能为空")
	}

	lgr := app_util.TakeLgr(cx).With(zap.String("curr_uid", currUid), zap.String("invitation_id", args.Id), zap.String("team_id", args.TeamId), zap.Uint32("role_mask", uint32(args.RoleMask)))

	cx = app_util.SaveLgr(cx, lgr)

	return a.inner.Update(cx, currUid, args)
}

// `Delete` enriches logger context and forwards call.
func (a *memberInvLogAppImpl) Delete(cx context.Context, currUid string, invId string) app_res.AppRes[app_res.None] {
	if cx == nil {
		cx = context.Background()
	}

	lgr := app_util.TakeLgr(cx).With(zap.String("curr_uid", currUid), zap.String("invitation_id", invId))

	cx = app_util.SaveLgr(cx, lgr)

	return a.inner.Delete(cx, currUid, invId)
}
