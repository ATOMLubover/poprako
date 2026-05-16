package app_impl

import (
	"context"

	app_iface "poprako-s/internal/app"
	app_res "poprako-s/internal/app/res"
	app_util "poprako-s/internal/app/util"
	"poprako-s/internal/app/val"
	"poprako-s/internal/domain/model/aggr"
	"poprako-s/internal/domain/model/enum"
	"poprako-s/internal/domain/model/query"
	repo_iface "poprako-s/internal/domain/repo"
	"poprako-s/internal/domain/svc"
	repo_infra "poprako-s/internal/infra/repo"

	"go.uber.org/zap"
)

func mkMemberInvRepoIncl(includes []enum.MemberInvIncl) []enum.MemberInvIncl {
	if len(includes) == 0 {
		return nil
	}

	repoIncls := make([]enum.MemberInvIncl, 0, len(includes))

	for i := range includes {
		repoIncls = append(repoIncls, includes[i])
	}

	return repoIncls
}

// `memberInvAppImpl` is the default implementation of `MemberInvApp`.
type memberInvAppImpl struct {
	txnCtrl repo_iface.TxnCtrl

	memberInvSvc svc.MemberInvSvc

	memberRepo    repo_iface.MemberRepo
	memberInvRepo repo_iface.MemberInvRepo

	errClsf repo_iface.ErrClsf
}

// `NewMemberInvApp` creates one `MemberInvApp` implementation.
func NewMemberInvApp(
	txnCtrl repo_iface.TxnCtrl,
	memberRepo repo_iface.MemberRepo,
	memberInvRepo repo_iface.MemberInvRepo,
	memberInvSvc svc.MemberInvSvc,
	errClsf repo_iface.ErrClsf,
) app_iface.MemberInvApp {
	if txnCtrl == nil || memberRepo == nil || memberInvRepo == nil || errClsf == nil {
		zap.L().Panic(
			"[NewMemberInvApp] nil dependency",
			zap.Bool("txnCtrl", txnCtrl == nil),
			zap.Bool("memberRepo", memberRepo == nil),
			zap.Bool("memberInvRepo", memberInvRepo == nil),
			zap.Bool("errClsf", errClsf == nil),
		)
	}

	return &memberInvAppImpl{
		txnCtrl:       txnCtrl,
		memberInvSvc:  memberInvSvc,
		memberRepo:    memberRepo,
		memberInvRepo: memberInvRepo,
		errClsf:       errClsf,
	}
}

// `List` lists invitations under one team.
func (a *memberInvAppImpl) List(cx context.Context, currUid string, args *val.ListMemberInvArgs) app_res.AppRes[[]val.MemberInvVal] {
	lgr := app_util.TakeLgr(cx)

	if re := vfyListMemberInvArgs(args); re.IsReject() {
		return app_res.Reject[[]val.MemberInvVal](re.Code(), re.Msg())
	}

	if re := a.memberInvSvc.CanListMemberInv(currUid, args.TeamId, a.memberRepo, a.errClsf); re.IsReject() {
		return app_res.Reject[[]val.MemberInvVal](app_res.ErrCode(re.Code()), re.Msg())
	}

	invs, err := a.memberInvRepo.List(query.ListMemberInvOpt{
		TeamId:   args.TeamId,
		Includes: args.Includes,
		Pending:  args.Pending,
		Pagi: query.PagiOpt{
			Offset: args.Offset,
			Limit:  args.Limit,
		},
	}, mkMemberInvRepoIncl(args.Includes)...)
	if err != nil {
		lgr.Error("[memberInvAppImpl.List] failed to list invitations", zap.Error(err))

		return app_res.Reject[[]val.MemberInvVal](app_res.ServerError, "获取邀请列表失败")
	}

	invVals := make([]val.MemberInvVal, len(invs))
	for i := range invs {
		invVals[i] = asmMemberInvVal(&invs[i])
	}

	return app_res.Accept(&invVals)
}

// `Create` creates one invitation under one team.
func (a *memberInvAppImpl) Create(cx context.Context, currUid string, args *val.CreateMemberInvArgs) app_res.AppRes[val.CreateMemberInvRes] {
	lgr := app_util.TakeLgr(cx)

	if args == nil || args.TeamId == "" || args.InviteeQid == "" {
		return app_res.Reject[val.CreateMemberInvRes](app_res.BadRequest, "team_id 和 invitee_qid 不能为空")
	}

	re, err := repo_iface.RunWithTxn[app_res.AppRes[val.CreateMemberInvRes]](a.txnCtrl, func(prov repo_iface.Prov) (app_res.AppRes[val.CreateMemberInvRes], error) {
		userRepo := prov.UserRepo()
		memberRepo := prov.MemberRepo()
		memberInvRepo := prov.MemberInvRepo()

		if re := a.memberInvSvc.CanAdminMemberInv(currUid, args.TeamId, memberRepo, a.errClsf); re.IsReject() {
			return app_res.Reject[val.CreateMemberInvRes](app_res.ErrCode(re.Code()), re.Msg()), app_res.DefErr()
		}

		if re := a.memberInvSvc.VfyInviteeNotMember(args.InviteeQid, args.TeamId, userRepo, memberRepo, a.errClsf); re.IsReject() {
			return app_res.Reject[val.CreateMemberInvRes](app_res.ErrCode(re.Code()), re.Msg()), app_res.DefErr()
		}

		memberInvCre, err := a.memberInvSvc.NewMemberInvCre(currUid, args.TeamId, args.InviteeQid, args.RoleMask)
		if err != nil {
			return app_res.Reject[val.CreateMemberInvRes](app_res.BadRequest, err.Error()), app_res.DefErr()
		}

		inv, err := memberInvRepo.Create(memberInvCre)
		if err != nil {
			if repo_infra.IsDupKey(err) {
				return app_res.Reject[val.CreateMemberInvRes](app_res.Conflict, "该用户在当前团队已有待处理邀请"), app_res.DefErr()
			}

			return app_res.Reject[val.CreateMemberInvRes](app_res.ServerError, "创建邀请失败"), err
		}

		return app_res.Accept(&val.CreateMemberInvRes{Id: inv.Id, InvCode: inv.InvCode}), nil
	})
	if err != nil {
		lgr.Error("[memberInvAppImpl.Create] failed to run transaction", zap.Error(err))

		return re
	}

	return re
}

// `Update` updates invitation role mask by put semantics.
func (a *memberInvAppImpl) Update(cx context.Context, currUid string, args *val.MemberInvUpdArgs) app_res.AppRes[app_res.None] {
	lgr := app_util.TakeLgr(cx)

	if args == nil || args.Id == "" {
		return app_res.Reject[app_res.None](app_res.BadRequest, "id 不能为空")
	}

	re, err := repo_iface.RunWithTxn[app_res.AppRes[app_res.None]](a.txnCtrl, func(prov repo_iface.Prov) (app_res.AppRes[app_res.None], error) {
		memberRepo := prov.MemberRepo()
		memberInvRepo := prov.MemberInvRepo()

		invitation, err := memberInvRepo.GetById(args.Id)
		if err != nil {
			if repo_infra.IsNotFound(err) {
				return app_res.Reject[app_res.None](app_res.BadRequest, "邀请不存在"), app_res.DefErr()
			}

			return app_res.Reject[app_res.None](app_res.ServerError, "更新邀请失败"), err
		}

		if re := a.memberInvSvc.CanAdminMemberInv(currUid, invitation.TeamId, memberRepo, a.errClsf); re.IsReject() {
			return app_res.Reject[app_res.None](app_res.ErrCode(re.Code()), re.Msg()), app_res.DefErr()
		}

		if err := a.memberInvSvc.VfyUpdateRoleMask(args.RoleMask); err != nil {
			return app_res.Reject[app_res.None](app_res.BadRequest, err.Error()), app_res.DefErr()
		}

		if err := memberInvRepo.Update(&aggr.MemberInvUpd{Id: args.Id, RoleMask: args.RoleMask}); err != nil {
			if repo_infra.IsNotFound(err) {
				return app_res.Reject[app_res.None](app_res.BadRequest, "邀请不存在"), app_res.DefErr()
			}

			return app_res.Reject[app_res.None](app_res.ServerError, "更新邀请失败"), err
		}

		return app_res.Accept(&app_res.None{}), nil
	})
	if err != nil {
		lgr.Error("[memberInvAppImpl.Update] failed to run transaction", zap.Error(err))

		return re
	}

	return re
}

// `Delete` deletes one invitation by id.
func (a *memberInvAppImpl) Delete(cx context.Context, currUid string, invId string) app_res.AppRes[app_res.None] {
	lgr := app_util.TakeLgr(cx)

	if invId == "" {
		return app_res.Reject[app_res.None](app_res.BadRequest, "invitation_id 不能为空")
	}

	re, err := repo_iface.RunWithTxn[app_res.AppRes[app_res.None]](a.txnCtrl, func(prov repo_iface.Prov) (app_res.AppRes[app_res.None], error) {
		memberRepo := prov.MemberRepo()
		memberInvRepo := prov.MemberInvRepo()

		invitation, err := memberInvRepo.GetById(invId)
		if err != nil {
			if repo_infra.IsNotFound(err) {
				return app_res.Reject[app_res.None](app_res.BadRequest, "邀请不存在"), app_res.DefErr()
			}

			return app_res.Reject[app_res.None](app_res.ServerError, "删除邀请失败"), err
		}

		if re := a.memberInvSvc.CanAdminMemberInv(currUid, invitation.TeamId, memberRepo, a.errClsf); re.IsReject() {
			return app_res.Reject[app_res.None](app_res.ErrCode(re.Code()), re.Msg()), app_res.DefErr()
		}

		if err := memberInvRepo.Delete(invId); err != nil {
			return app_res.Reject[app_res.None](app_res.ServerError, "删除邀请失败"), err
		}

		return app_res.Accept(&app_res.None{}), nil
	})
	if err != nil {
		lgr.Error("[memberInvAppImpl.Delete] failed to run transaction", zap.Error(err))

		return re
	}

	return re
}
