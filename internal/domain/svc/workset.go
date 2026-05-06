package svc

import (
	"poprako-s/internal/domain/model/aggr"
	"poprako-s/internal/domain/model/enum"
	repo_iface "poprako-s/internal/domain/repo"
	svc_res "poprako-s/internal/domain/svc/res"
	"poprako-s/pkg/util"

	"go.uber.org/zap"
)

// `WorksetSvc` provides stateless domain services for the `Workset` aggregate.
type WorksetSvc struct{}

// `NewWorksetSvc` returns a ready-to-use `WorksetSvc`.
func NewWorksetSvc() WorksetSvc {
	return WorksetSvc{}
}

// `NewWorksetCre` builds a `WorksetCre` aggregate with a generated id.
// The `index` must be pre-computed by the caller from a transactional count
// of active worksets for the given team.
func (WorksetSvc) NewWorksetCre(teamId string, index int, name string, desc *string) *aggr.WorksetCre {
	return &aggr.WorksetCre{
		Id:     util.GenId("workset"),
		TeamId: teamId,
		Index:  index,
		Name:   name,
		Desc:   desc,
	}
}

func (WorksetSvc) CanAdminWorkset(currUid string, teamId string, memberRepo repo_iface.MemberRepo, clsf repo_iface.ErrClsf) svc_res.SvcRes {
	member, err := memberRepo.GetByUserTeamId(currUid, teamId)
	if err != nil {
		zap.L().Error(
			"[WorksetSvc.CanAdminWorkset] failed to verify member",
			zap.String("teamId", teamId),
			zap.String("currUid", currUid),
			zap.Error(err),
		)

		return classifyRepoErr(err, clsf, svc_res.Forbidden, "仅汉化组管理员可执行该操作", "获取成员信息失败", "成员信息服务暂不可用", "权限校验失败")
	}
	if member == nil || !member.HasAnyRole(enum.RoleAdmin) {
		return svc_res.Reject(svc_res.Forbidden, "仅汉化组管理员可执行该操作")
	}

	return svc_res.Accept()
}

func (WorksetSvc) CanListWorkset(currUid string, teamId string, memberRepo repo_iface.MemberRepo, clsf repo_iface.ErrClsf) svc_res.SvcRes {
	ok, err := memberRepo.ExistByUserTeamId(currUid, teamId)
	if err != nil {
		zap.L().Error(
			"[WorksetSvc.CanListWorkset] failed to verify member",
			zap.String("teamId", teamId),
			zap.String("currUid", currUid),
			zap.Error(err),
		)

		return classifyRepoErr(err, clsf, svc_res.Forbidden, "仅汉化组成员可查看作品集", "获取成员信息失败", "成员信息服务暂不可用", "获取成员信息失败")
	}
	if !ok {
		return svc_res.Reject(svc_res.Forbidden, "仅汉化组成员可查看作品集")
	}

	return svc_res.Accept()
}
