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

// `MemberSvc` provides domain services for member creation
type MemberSvc struct{}

// `NewMemberSvc` creates a new `MemberSvc`
func NewMemberSvc() MemberSvc {
	return MemberSvc{}
}

// `NewMemberCre` builds a `MemberCre` aggregate with a generated id and the given role mask
func (MemberSvc) NewMemberCre(userId string, teamId string, roleMask aggr.RoleMask) *aggr.MemberCre {
	id := util.GenId("member")

	return &aggr.MemberCre{
		Id:       id,
		UserId:   userId,
		TeamId:   teamId,
		RoleMask: roleMask,
	}
}

// `CanListMember` validates whether current user can list members in a team.
func (MemberSvc) CanListMember(currUid string, teamId string, memberRepo repo_iface.MemberRepo, clsf repo_iface.ErrClsf) svc_res.SvcRes {
	ok, err := memberRepo.ExistByUserTeamId(currUid, teamId)
	if err != nil {
		zap.L().Error(
			"[MemberSvc.CanListMember] failed to verify member",
			zap.String("currUid", currUid),
			zap.String("teamId", teamId),
			zap.Error(err),
		)

		return classifyRepoErr(err, clsf, svc_res.Forbidden, "仅团队成员可查看成员列表", "权限校验失败", "权限校验服务暂不可用", "权限校验失败")
	}

	if !ok {
		return svc_res.Reject(svc_res.Forbidden, "仅团队成员可查看成员列表")
	}

	return svc_res.Accept()
}

// `CanAdminMember` validates whether current user can admin member resources.
func (MemberSvc) CanAdminMember(currUid string, teamId string, memberRepo repo_iface.MemberRepo, clsf repo_iface.ErrClsf) svc_res.SvcRes {
	currMember, err := memberRepo.GetByUserTeamId(currUid, teamId)
	if err != nil {
		zap.L().Error(
			"[MemberSvc.CanAdminMember] failed to verify admin member",
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

// `CanCreateMember` validates whether current user can create members for one team.
func (s MemberSvc) CanCreateMember(currUid string, teamId string, memberRepo repo_iface.MemberRepo, clsf repo_iface.ErrClsf) svc_res.SvcRes {
	return s.CanAdminMember(currUid, teamId, memberRepo, clsf)
}

// `CanUpdateMember` validates whether current user can update one target member roles.
func (s MemberSvc) CanUpdateMember(currUid string, target *aggr.Member, memberRepo repo_iface.MemberRepo, clsf repo_iface.ErrClsf) svc_res.SvcRes {
	if target == nil {
		return svc_res.Reject(svc_res.BadRequest, "成员不存在")
	}

	return s.CanAdminMember(currUid, target.TeamId, memberRepo, clsf)
}

// `CanDeleteMember` validates whether current user can delete one target member.
func (s MemberSvc) CanDeleteMember(currUid string, target *aggr.Member, memberRepo repo_iface.MemberRepo, clsf repo_iface.ErrClsf) svc_res.SvcRes {
	if target == nil {
		return svc_res.Reject(svc_res.BadRequest, "成员不存在")
	}

	return s.CanAdminMember(currUid, target.TeamId, memberRepo, clsf)
}

// `CanJoinTeam` validates join operation by invitation and target membership state.
func (s MemberSvc) CanJoinTeam(currUid string, currQid string, inv *aggr.MemberInv, memberRepo repo_iface.MemberRepo, clsf repo_iface.ErrClsf) svc_res.SvcRes {
	if inv == nil || !inv.Pending {
		return svc_res.Reject(svc_res.BadRequest, "邀请码无效或已被使用")
	}

	if inv.InviteeQid != currQid {
		return svc_res.Reject(svc_res.BadRequest, "邀请码无效或已被使用")
	}

	ok, err := memberRepo.ExistByUserTeamId(currUid, inv.TeamId)
	if err != nil {
		zap.L().Error(
			"[MemberSvc.CanJoinTeam] failed to verify existing member",
			zap.String("currUid", currUid),
			zap.String("teamId", inv.TeamId),
			zap.Error(err),
		)

		return classifyRepoErr(err, clsf, svc_res.ServerError, "加入团队失败", "加入团队失败", "加入团队失败", "加入团队失败")
	}

	if ok {
		return svc_res.Reject(svc_res.BadRequest, "您已经是该团队成员")
	}

	return svc_res.Accept()
}

// `NewMemberRoleUpd` builds role update payload with put semantics.
func (MemberSvc) NewMemberRoleUpd(memberId string, roleMask aggr.RoleMask) (*aggr.MemberRoleUpd, error) {
	if memberId == "" {
		return nil, errors.New("成员 ID 不能为空")
	}

	if roleMask == 0 {
		return nil, errors.New("至少指定一个角色")
	}

	return &aggr.MemberRoleUpd{
		Id:       memberId,
		RoleMask: roleMask,
	}, nil
}
