package app_impl

import (
	app_res "poprako-s/internal/app/res"
	app_util "poprako-s/internal/app/util"
	"poprako-s/internal/app/val"
	"poprako-s/internal/domain/model/aggr"
)

// `asmMemberInvVal` converts one member invitation aggregate into app value object.
func asmMemberInvVal(inv *aggr.MemberInv) val.MemberInvVal {
	return val.MemberInvVal{
		Id:         inv.Id,
		InvitorId:  inv.InvitorId,
		TeamId:     inv.TeamId,
		InviteeQid: inv.InviteeQid,
		InvCode:    inv.InvCode,
		Pending:    inv.Pending,
		RoleMask:   inv.RoleMask,
		CreatedAt:  inv.CreatedAt.UnixMilli(),
	}
}

// `vfyListMemberInvArgs` validates and normalizes member invitation list args.
func vfyListMemberInvArgs(args *val.ListMemberInvArgs) app_res.AppRes[app_res.None] {
	if args == nil || args.TeamId == "" {
		return app_res.Reject[app_res.None](app_res.BadRequest, "team_id 不能为空")
	}

	if re := app_util.ClampOffsetLimit(args.Offset, &args.Limit); re.IsReject() {
		return re
	}

	return app_res.Accept(&app_res.None{})
}
