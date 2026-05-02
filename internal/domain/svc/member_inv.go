package svc

import (
	"errors"

	"poprako-s/internal/domain/model/aggr"
	"poprako-s/internal/domain/model/enum"
	repo_iface "poprako-s/internal/domain/repo"
	svc_res "poprako-s/internal/domain/svc/res"
	"poprako-s/pkg/util"

	"go.uber.org/zap"
)

// `MemberInvSvc` provides stateless helpers for member invitation aggregate.
type MemberInvSvc struct{}

// `NewMemberInvSvc` returns one ready-to-use `MemberInvSvc`.
func NewMemberInvSvc() MemberInvSvc {
	return MemberInvSvc{}
}

// `CanListMemberInv` validates whether current user can list team invitations.
func (MemberInvSvc) CanListMemberInv(currUid string, teamId string, memberRepo repo_iface.MemberRepo, clsf repo_iface.ErrClsf) svc_res.SvcRes {
	ok, err := memberRepo.ExistByUserTeamId(currUid, teamId)
	if err != nil {
		zap.L().Error(
			"[MemberInvSvc.CanListMemberInv] failed to verify member",
			zap.String("currUid", currUid),
			zap.String("teamId", teamId),
			zap.Error(err),
		)

		return classifyRepoErr(err, clsf, svc_res.Forbidden, "仅团队成员可查看邀请列表", "权限校验失败", "权限校验服务暂不可用", "权限校验失败")
	}

	if !ok {
		return svc_res.Reject(svc_res.Forbidden, "仅团队成员可查看邀请列表")
	}

	return svc_res.Accept()
}

// `CanAdminMemberInv` validates whether current user can create update delete team invitations.
func (MemberInvSvc) CanAdminMemberInv(currUid string, teamId string, memberRepo repo_iface.MemberRepo, clsf repo_iface.ErrClsf) svc_res.SvcRes {
	currMember, err := memberRepo.GetByUserTeamId(currUid, teamId)
	if err != nil {
		zap.L().Error(
			"[MemberInvSvc.CanAdminMemberInv] failed to verify admin member",
			zap.String("currUid", currUid),
			zap.String("teamId", teamId),
			zap.Error(err),
		)

		return classifyRepoErr(err, clsf, svc_res.Forbidden, "仅团队管理员可执行该操作", "权限校验失败", "权限校验服务暂不可用", "权限校验失败")
	}

	if currMember == nil || !currMember.HasAnyRole(enum.RoleAdmin) {
		return svc_res.Reject(svc_res.Forbidden, "仅团队管理员可执行该操作")
	}

	return svc_res.Accept()
}

// `NewMemberInvCre` builds one member invitation create payload.
func (MemberInvSvc) NewMemberInvCre(invitorId string, teamId string, inviteeQid string, roleMask aggr.RoleMask) (*aggr.MemberInvCre, error) {
	if inviteeQid == "" {
		return nil, errors.New("invitee_qid 不能为空")
	}

	if roleMask == 0 {
		return nil, errors.New("至少指定一个角色")
	}

	id := util.GenId("member_invitation")

	return &aggr.MemberInvCre{
		Id:         id,
		InvitorId:  invitorId,
		TeamId:     teamId,
		InviteeQid: inviteeQid,
		InvCode:    genMemberInvCode(),
		RoleMask:   roleMask,
	}, nil
}

// `VfyUpdateRoleMask` validates invitation role mask for put update.
func (MemberInvSvc) VfyUpdateRoleMask(roleMask aggr.RoleMask) error {
	if roleMask == 0 {
		return errors.New("至少指定一个角色")
	}

	return nil
}

func genMemberInvCode() string {
	id := util.GenId("member_inv_code")

	if len(id) >= 6 {
		return id[len(id)-6:]
	}

	return id
}
