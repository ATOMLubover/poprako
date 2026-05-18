package app_impl

import (
	"context"

	app_iface "poprako-s/internal/app"
	app_res "poprako-s/internal/app/res"
	app_util "poprako-s/internal/app/util"
	"poprako-s/internal/app/val"
	oss_iface "poprako-s/internal/domain/ext/oss"
	repo_iface "poprako-s/internal/domain/repo"
	"poprako-s/internal/domain/svc"
	repo_infra "poprako-s/internal/infra/repo"

	"go.uber.org/zap"
)

// `memberAppImpl` is the default implementation of `MemberApp`.
type memberAppImpl struct {
	txnCtrl repo_iface.TxnCtrl

	memberSvc svc.MemberSvc

	userRepo      repo_iface.UserRepo
	memberRepo    repo_iface.MemberRepo
	memberInvRepo repo_iface.MemberInvRepo

	ossSigner oss_iface.Signer
	errClsf   repo_iface.ErrClsf
}

// `NewMemberApp` creates one `MemberApp` implementation.
func NewMemberApp(
	txnCtrl repo_iface.TxnCtrl,
	userRepo repo_iface.UserRepo,
	memberRepo repo_iface.MemberRepo,
	memberInvRepo repo_iface.MemberInvRepo,
	memberSvc svc.MemberSvc,
	ossSigner oss_iface.Signer,
	errClsf repo_iface.ErrClsf,
) app_iface.MemberApp {
	if txnCtrl == nil || userRepo == nil || memberRepo == nil || memberInvRepo == nil || ossSigner == nil || errClsf == nil {
		zap.L().Panic(
			"[NewMemberApp] nil dependency",
			zap.Bool("txnCtrl", txnCtrl == nil),
			zap.Bool("userRepo", userRepo == nil),
			zap.Bool("memberRepo", memberRepo == nil),
			zap.Bool("memberInvRepo", memberInvRepo == nil),
			zap.Bool("ossSigner", ossSigner == nil),
			zap.Bool("errClsf", errClsf == nil),
		)
	}

	return &memberAppImpl{
		txnCtrl:       txnCtrl,
		memberSvc:     memberSvc,
		userRepo:      userRepo,
		memberRepo:    memberRepo,
		memberInvRepo: memberInvRepo,
		ossSigner:     ossSigner,
		errClsf:       errClsf,
	}
}

// `Create` creates one member under one team.
func (a *memberAppImpl) Create(cx context.Context, currUid string, args *val.CreateMemberArgs) app_res.AppRes[val.CreateMemberRes] {
	lgr := app_util.TakeLgr(cx)

	if args == nil || args.UserId == "" || args.TeamId == "" {
		return app_res.Reject[val.CreateMemberRes](app_res.BadRequest, "user_id 和 team_id 不能为空")
	}

	re, err := repo_iface.RunWithTxn[app_res.AppRes[val.CreateMemberRes]](a.txnCtrl, func(prov repo_iface.Prov) (app_res.AppRes[val.CreateMemberRes], error) {
		userRepo := prov.UserRepo()
		memberRepo := prov.MemberRepo()

		if re := a.memberSvc.CanCreateMember(currUid, prov.UserRepo(), a.errClsf); re.IsReject() {
			return app_res.Reject[val.CreateMemberRes](app_res.ErrCode(re.Code()), re.Msg()), app_res.DefErr()
		}

		exists, err := memberRepo.ExistByUserTeamId(args.UserId, args.TeamId)
		if err != nil {
			return app_res.Reject[val.CreateMemberRes](app_res.ServerError, "创建成员失败"), err
		}

		if exists {
			return app_res.Reject[val.CreateMemberRes](app_res.Conflict, "该用户已经加入该团队"), app_res.DefErr()
		}

		user, err := userRepo.GetById(args.UserId)
		if err != nil {
			if repo_infra.IsNotFound(err) {
				return app_res.Reject[val.CreateMemberRes](app_res.BadRequest, "用户不存在"), app_res.DefErr()
			}

			return app_res.Reject[val.CreateMemberRes](app_res.ServerError, "创建成员失败"), err
		}

		memberCre := a.memberSvc.NewMemberCre(args.UserId, user.Nickname, args.TeamId, args.RoleMask)
		member, err := memberRepo.Create(memberCre)
		if err != nil {
			return app_res.Reject[val.CreateMemberRes](app_res.ServerError, "创建成员失败"), err
		}

		return app_res.Accept(&val.CreateMemberRes{Id: member.Id}), nil
	})
	if err != nil {
		lgr.Error("[memberAppImpl.Create] failed to run transaction", zap.Error(err))

		return re
	}

	return re
}

// `ListByTeam` lists members under one team.
func (a *memberAppImpl) ListByTeam(cx context.Context, currUid string, args *val.ListMemberByTeamArgs) app_res.AppRes[[]val.MemberVal] {
	lgr := app_util.TakeLgr(cx)

	if args == nil || args.TeamId == "" {
		return app_res.Reject[[]val.MemberVal](app_res.BadRequest, "team_id 不能为空")
	}

	if re := vfyListMemberArgs(args.Offset, &args.Limit); re.IsReject() {
		return app_res.Reject[[]val.MemberVal](re.Code(), re.Msg())
	}

	if re := a.memberSvc.CanListMember(currUid, args.TeamId, a.memberRepo, a.errClsf); re.IsReject() {
		return app_res.Reject[[]val.MemberVal](app_res.ErrCode(re.Code()), re.Msg())
	}

	members, err := a.memberRepo.List(mkListMemberOptByTeam(args.TeamId, args.UserNicknameKeyword, args.Includes, args.Offset, args.Limit), args.Includes...)
	if err != nil {
		lgr.Error("[memberAppImpl.ListByTeam] failed to list members", zap.Error(err))

		return app_res.Reject[[]val.MemberVal](app_res.ServerError, "获取成员列表失败")
	}

	memberVals := make([]val.MemberVal, len(members))
	for i := range members {
		memberVal, err := asmMemberVal(members[i], a.ossSigner)
		if err != nil {
			lgr.Error("[memberAppImpl.ListByTeam] failed to assemble member value", zap.Error(err))

			return app_res.Reject[[]val.MemberVal](app_res.ServerError, "获取成员列表失败")
		}

		memberVals[i] = *memberVal
	}

	return app_res.Accept(&memberVals)
}

// `ListMine` lists memberships of current user.
func (a *memberAppImpl) ListMine(cx context.Context, currUid string, args *val.ListMyMemberArgs) app_res.AppRes[[]val.MemberVal] {
	lgr := app_util.TakeLgr(cx)

	if args == nil {
		return app_res.Reject[[]val.MemberVal](app_res.BadRequest, "分页参数不能为空")
	}

	if re := vfyListMemberArgs(args.Offset, &args.Limit); re.IsReject() {
		return app_res.Reject[[]val.MemberVal](re.Code(), re.Msg())
	}

	members, err := a.memberRepo.List(mkListMemberOptByUser(currUid, args.Includes, args.Offset, args.Limit), args.Includes...)
	if err != nil {
		lgr.Error("[memberAppImpl.ListMine] failed to list members", zap.Error(err))

		return app_res.Reject[[]val.MemberVal](app_res.ServerError, "获取成员列表失败")
	}

	memberVals := make([]val.MemberVal, len(members))
	for i := range members {
		memberVal, err := asmMemberVal(members[i], a.ossSigner)
		if err != nil {
			lgr.Error("[memberAppImpl.ListMine] failed to assemble member value", zap.Error(err))

			return app_res.Reject[[]val.MemberVal](app_res.ServerError, "获取成员列表失败")
		}

		memberVals[i] = *memberVal
	}

	return app_res.Accept(&memberVals)
}

// `UpdateRole` updates one member role mask by put semantics.
func (a *memberAppImpl) UpdateRole(cx context.Context, currUid string, args *val.MemberRoleUpdArgs) app_res.AppRes[app_res.None] {
	lgr := app_util.TakeLgr(cx)

	if args == nil || args.Id == "" {
		return app_res.Reject[app_res.None](app_res.BadRequest, "id 不能为空")
	}

	re, err := repo_iface.RunWithTxn[app_res.AppRes[app_res.None]](a.txnCtrl, func(prov repo_iface.Prov) (app_res.AppRes[app_res.None], error) {
		memberRepo := prov.MemberRepo()

		target, err := memberRepo.GetById(args.Id)
		if err != nil {
			if repo_infra.IsNotFound(err) {
				return app_res.Reject[app_res.None](app_res.BadRequest, "成员不存在"), app_res.DefErr()
			}

			return app_res.Reject[app_res.None](app_res.ServerError, "更新成员角色失败"), err
		}

		if re := a.memberSvc.CanUpdateMember(currUid, target, memberRepo, a.errClsf); re.IsReject() {
			return app_res.Reject[app_res.None](app_res.ErrCode(re.Code()), re.Msg()), app_res.DefErr()
		}

		memberUpd, err := a.memberSvc.NewMemberRoleUpd(args.Id, args.RoleMask)
		if err != nil {
			return app_res.Reject[app_res.None](app_res.BadRequest, err.Error()), app_res.DefErr()
		}

		if err := memberRepo.UpdateRoles(memberUpd); err != nil {
			return app_res.Reject[app_res.None](app_res.ServerError, "更新成员角色失败"), err
		}

		return app_res.Accept(&app_res.None{}), nil
	})
	if err != nil {
		lgr.Error("[memberAppImpl.UpdateRole] failed to run transaction", zap.Error(err))

		return re
	}

	return re
}

// `Delete` deletes one member by id.
func (a *memberAppImpl) Delete(cx context.Context, currUid string, memberId string) app_res.AppRes[app_res.None] {
	lgr := app_util.TakeLgr(cx)

	if memberId == "" {
		return app_res.Reject[app_res.None](app_res.BadRequest, "member_id 不能为空")
	}

	re, err := repo_iface.RunWithTxn[app_res.AppRes[app_res.None]](a.txnCtrl, func(prov repo_iface.Prov) (app_res.AppRes[app_res.None], error) {
		memberRepo := prov.MemberRepo()

		target, err := memberRepo.GetById(memberId)
		if err != nil {
			if repo_infra.IsNotFound(err) {
				return app_res.Reject[app_res.None](app_res.BadRequest, "成员不存在"), app_res.DefErr()
			}

			return app_res.Reject[app_res.None](app_res.ServerError, "删除成员失败"), err
		}

		if re := a.memberSvc.CanDeleteMember(currUid, target, memberRepo, a.errClsf); re.IsReject() {
			return app_res.Reject[app_res.None](app_res.ErrCode(re.Code()), re.Msg()), app_res.DefErr()
		}

		if err := memberRepo.Delete(memberId); err != nil {
			return app_res.Reject[app_res.None](app_res.ServerError, "删除成员失败"), err
		}

		return app_res.Accept(&app_res.None{}), nil
	})
	if err != nil {
		lgr.Error("[memberAppImpl.Delete] failed to run transaction", zap.Error(err))

		return re
	}

	return re
}

// `JoinTeam` creates one membership by invitation code.
func (a *memberAppImpl) JoinTeam(cx context.Context, currUid string, args *val.JoinTeamArgs) app_res.AppRes[app_res.None] {
	lgr := app_util.TakeLgr(cx)

	if args == nil || args.InvCode == "" {
		return app_res.Reject[app_res.None](app_res.BadRequest, "invitation_code 不能为空")
	}

	re, err := repo_iface.RunWithTxn[app_res.AppRes[app_res.None]](a.txnCtrl, func(prov repo_iface.Prov) (app_res.AppRes[app_res.None], error) {
		userRepo := prov.UserRepo()
		memberRepo := prov.MemberRepo()
		memberInvRepo := prov.MemberInvRepo()

		currUser, err := userRepo.GetById(currUid)
		if err != nil {
			if repo_infra.IsNotFound(err) {
				return app_res.Reject[app_res.None](app_res.NotFound, "用户不存在"), app_res.DefErr()
			}

			return app_res.Reject[app_res.None](app_res.ServerError, "加入团队失败"), err
		}

		targetInv, err := memberInvRepo.GetPendingByInviteeQid(currUser.Qid)
		if err != nil {
			if repo_infra.IsNotFound(err) {
				return app_res.Reject[app_res.None](app_res.BadRequest, "邀请码无效或已被使用"), app_res.DefErr()
			}

			return app_res.Reject[app_res.None](app_res.ServerError, "加入团队失败"), err
		}

		if targetInv.InvCode != args.InvCode {
			return app_res.Reject[app_res.None](app_res.BadRequest, "邀请码无效或已被使用"), app_res.DefErr()
		}

		if re := a.memberSvc.CanJoinTeam(currUid, currUser.Qid, targetInv, memberRepo, a.errClsf); re.IsReject() {
			return app_res.Reject[app_res.None](app_res.ErrCode(re.Code()), re.Msg()), app_res.DefErr()
		}

		memberCre := a.memberSvc.NewMemberCre(currUid, currUser.Nickname, targetInv.TeamId, targetInv.RoleMask)
		if _, err := memberRepo.Create(memberCre); err != nil {
			return app_res.Reject[app_res.None](app_res.ServerError, "加入团队失败"), err
		}

		if err := memberInvRepo.MarkCompleted(targetInv.Id); err != nil {
			return app_res.Reject[app_res.None](app_res.ServerError, "加入团队失败"), err
		}

		return app_res.Accept(&app_res.None{}), nil
	})
	if err != nil {
		lgr.Error("[memberAppImpl.JoinTeam] failed to run transaction", zap.Error(err))

		return re
	}

	return re
}
