package app_impl

import (
	"context"

	app_iface "poprako-s/internal/app"
	app_res "poprako-s/internal/app/res"
	app_util "poprako-s/internal/app/util"
	"poprako-s/internal/app/val"

	"go.uber.org/zap"
)

// `memberLogAppImpl` is logging decorator for `MemberApp`.
type memberLogAppImpl struct {
	inner app_iface.MemberApp
}

// `NewMemberLogApp` creates logging decorator for `MemberApp`.
func NewMemberLogApp(inner app_iface.MemberApp) app_iface.MemberApp {
	if inner == nil {
		zap.L().Panic("[NewMemberLogApp] nil dependency", zap.Bool("inner", inner == nil))
	}

	return &memberLogAppImpl{inner: inner}
}

// `Create` enriches logger context and forwards call.
func (a *memberLogAppImpl) Create(cx context.Context, currUid string, args *val.CreateMemberArgs) app_res.AppRes[val.CreateMemberRes] {
	if cx == nil {
		cx = context.Background()
	}

	if args == nil {
		return app_res.Reject[val.CreateMemberRes](app_res.BadRequest, "创建参数不能为空")
	}

	lgr := app_util.TakeLgr(cx).With(zap.String("curr_uid", currUid), zap.String("user_id", args.UserId), zap.String("team_id", args.TeamId))

	cx = app_util.SaveLgr(cx, lgr)

	return a.inner.Create(cx, currUid, args)
}

// `ListByTeam` enriches logger context and forwards call.
func (a *memberLogAppImpl) ListByTeam(cx context.Context, currUid string, args *val.ListMemberByTeamArgs) app_res.AppRes[[]val.MemberVal] {
	if cx == nil {
		cx = context.Background()
	}

	if args == nil {
		return app_res.Reject[[]val.MemberVal](app_res.BadRequest, "分页参数不能为空")
	}

	lgr := app_util.TakeLgr(cx).With(zap.String("curr_uid", currUid), zap.String("team_id", args.TeamId), zap.Int("offset", args.Offset), zap.Int("limit", args.Limit))

	cx = app_util.SaveLgr(cx, lgr)

	return a.inner.ListByTeam(cx, currUid, args)
}

// `ListMine` enriches logger context and forwards call.
func (a *memberLogAppImpl) ListMine(cx context.Context, currUid string, args *val.ListMyMemberArgs) app_res.AppRes[[]val.MemberVal] {
	if cx == nil {
		cx = context.Background()
	}

	if args == nil {
		return app_res.Reject[[]val.MemberVal](app_res.BadRequest, "分页参数不能为空")
	}

	lgr := app_util.TakeLgr(cx).With(zap.String("curr_uid", currUid), zap.Int("offset", args.Offset), zap.Int("limit", args.Limit))

	cx = app_util.SaveLgr(cx, lgr)

	return a.inner.ListMine(cx, currUid, args)
}

// `UpdateRole` enriches logger context and forwards call.
func (a *memberLogAppImpl) UpdateRole(cx context.Context, currUid string, args *val.MemberRoleUpdArgs) app_res.AppRes[app_res.None] {
	if cx == nil {
		cx = context.Background()
	}

	if args == nil {
		return app_res.Reject[app_res.None](app_res.BadRequest, "更新参数不能为空")
	}

	lgr := app_util.TakeLgr(cx).With(zap.String("curr_uid", currUid), zap.String("member_id", args.Id), zap.Uint32("role_mask", uint32(args.RoleMask)))

	cx = app_util.SaveLgr(cx, lgr)

	return a.inner.UpdateRole(cx, currUid, args)
}

// `Delete` enriches logger context and forwards call.
func (a *memberLogAppImpl) Delete(cx context.Context, currUid string, memberId string) app_res.AppRes[app_res.None] {
	if cx == nil {
		cx = context.Background()
	}

	lgr := app_util.TakeLgr(cx).With(zap.String("curr_uid", currUid), zap.String("member_id", memberId))

	cx = app_util.SaveLgr(cx, lgr)

	return a.inner.Delete(cx, currUid, memberId)
}

// `JoinTeam` enriches logger context and forwards call.
func (a *memberLogAppImpl) JoinTeam(cx context.Context, currUid string, args *val.JoinTeamArgs) app_res.AppRes[app_res.None] {
	if cx == nil {
		cx = context.Background()
	}

	if args == nil {
		return app_res.Reject[app_res.None](app_res.BadRequest, "加入参数不能为空")
	}

	lgr := app_util.TakeLgr(cx).With(zap.String("curr_uid", currUid))

	cx = app_util.SaveLgr(cx, lgr)

	return a.inner.JoinTeam(cx, currUid, args)
}
