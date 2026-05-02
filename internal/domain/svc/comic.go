package svc

import (
	"poprako-s/internal/domain/model/aggr"
	"poprako-s/internal/domain/model/enum"
	repo_iface "poprako-s/internal/domain/repo"
	svc_res "poprako-s/internal/domain/svc/res"
	"poprako-s/pkg/util"

	"go.uber.org/zap"
)

// `ComicSvc` provides stateless helpers for comic aggregate workflows
// It only builds domain payloads and does not access repositories directly
// This keeps business construction logic centralized and reusable
type ComicSvc struct{}

// `NewComicSvc` returns a ready-to-use `ComicSvc`
func NewComicSvc() ComicSvc {
	return ComicSvc{}
}

func (ComicSvc) CanListComic(currUid string, teamId string, memberRepo repo_iface.MemberRepo, clsf repo_iface.ErrClsf) svc_res.SvcRes {
	ok, err := memberRepo.ExistByUserTeamId(currUid, teamId)
	if err != nil {
		zap.L().Error(
			"[ComicSvc.CanListComic] failed to verify member",
			zap.String("teamId", teamId),
			zap.String("currUid", currUid),
			zap.Error(err),
		)

		return classifyRepoErr(err, clsf, svc_res.Forbidden, "仅汉化组成员可查看该作品集的漫画", "获取成员信息失败", "成员信息服务暂不可用", "获取成员信息失败")
	}
	if !ok {
		return svc_res.Reject(svc_res.Forbidden, "无权访问该作品集的漫画")
	}

	return svc_res.Accept()
}

func (ComicSvc) CanAdminComic(currUid string, teamId string, memberRepo repo_iface.MemberRepo, clsf repo_iface.ErrClsf) svc_res.SvcRes {
	member, err := memberRepo.GetByUserTeamId(currUid, teamId)
	if err != nil {
		zap.L().Error(
			"[ComicSvc.CanAdminComic] failed to verify member",
			zap.String("teamId", teamId),
			zap.String("currUid", currUid),
			zap.Error(err),
		)

		return classifyRepoErr(err, clsf, svc_res.Forbidden, "仅汉化组管理员可操作该作品集的漫画", "获取成员信息失败", "成员信息服务暂不可用", "权限校验失败")
	}
	if member == nil || !member.HasAnyRole(enum.RoleAdmin) {
		return svc_res.Reject(svc_res.Forbidden, "仅汉化组管理员可操作该作品集的漫画")
	}

	return svc_res.Accept()
}

// `NewComicCre` builds a `ComicCre` aggregate with generated id
// Caller should pre-compute `index` in transaction for consistency
// `desc` keeps nullable semantics to support put-style updates later
func (ComicSvc) NewComicCre(
	worksetId string,
	index int,
	title string,
	author string,
	desc *string,
	creatorId string,
) *aggr.ComicCre {
	return &aggr.ComicCre{
		Id:        util.GenId("comic"),
		WorksetId: worksetId,
		Index:     index,
		Title:     title,
		Author:    author,
		Desc:      desc,
		CreatorId: creatorId,
	}
}
