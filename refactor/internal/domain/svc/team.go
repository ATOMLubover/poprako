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

// `TeamSvc` provides stateless helpers for team aggregate.
type TeamSvc struct{}

// `NewTeamSvc` returns one ready-to-use `TeamSvc`.
func NewTeamSvc() TeamSvc {
	return TeamSvc{}
}

// `CanListTeam` validates whether current user can list all teams.
func (TeamSvc) CanListTeam(currUid string, userRepo repo_iface.UserRepo, clsf repo_iface.ErrClsf) svc_res.SvcRes {
	currUser, err := userRepo.GetById(currUid)
	if err != nil {
		zap.L().Error(
			"[TeamSvc.CanListTeam] failed to get current user",
			zap.String("currUid", currUid),
			zap.Error(err),
		)

		return classifyRepoErr(err, clsf, svc_res.Forbidden, "仅超级管理员可查看所有团队", "权限校验失败", "权限校验服务暂不可用", "权限校验失败")
	}

	if currUser == nil || !currUser.IsSuperAdmin {
		return svc_res.Reject(svc_res.Forbidden, "仅超级管理员可查看所有团队")
	}

	return svc_res.Accept()
}

// `CanAdminTeam` validates whether current user can admin one team.
func (TeamSvc) CanAdminTeam(currUid string, teamId string, memberRepo repo_iface.MemberRepo, clsf repo_iface.ErrClsf) svc_res.SvcRes {
	member, err := memberRepo.GetByUserTeamId(currUid, teamId)
	if err != nil {
		zap.L().Error(
			"[TeamSvc.CanAdminTeam] failed to verify member",
			zap.String("currUid", currUid),
			zap.String("teamId", teamId),
			zap.Error(err),
		)

		return classifyRepoErr(err, clsf, svc_res.Forbidden, "仅团队管理员可执行该操作", "权限校验失败", "权限校验服务暂不可用", "权限校验失败")
	}

	if member == nil || !member.HasAnyRole(enum.RoleAdmin) {
		return svc_res.Reject(svc_res.Forbidden, "仅团队管理员可执行该操作")
	}

	return svc_res.Accept()
}

// `NewTeamCre` builds one team creation payload.
func (TeamSvc) NewTeamCre(name string, desc string) (*aggr.TeamCre, error) {
	if name == "" {
		return nil, errors.New("团队名称不能为空")
	}

	return &aggr.TeamCre{
		Id:   util.GenId("team"),
		Name: name,
		Desc: desc,
	}, nil
}

// `GenAvatarKey` generates one team avatar object key by team id and extension.
func (TeamSvc) GenAvatarKey(teamId string, fileExt string) string {
	return "team_avatar/" + teamId + "." + fileExt
}
