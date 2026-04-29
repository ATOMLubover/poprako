package svc

import (
	"fmt"

	"poprako-s/internal/domain/model/aggr"
	"poprako-s/internal/domain/model/enum"
	repo_iface "poprako-s/internal/domain/repo"
	svc_res "poprako-s/internal/domain/svc/res"
	"poprako-s/pkg/util"

	"go.uber.org/zap"
)

// `ChapterSvc` provides stateless helpers for chapter aggregate.
type ChapterSvc struct{}

// `NewChapterSvc` returns a ready-to-use `ChapterSvc`.
func NewChapterSvc() ChapterSvc {
	return ChapterSvc{}
}

func (ChapterSvc) CanListChapter(currUid string, teamId string, memberRepo repo_iface.MemberRepo, clsf repo_iface.ErrClsf) svc_res.SvcRes {
	ok, err := memberRepo.ExistByUserTeamId(currUid, teamId)
	if err != nil {
		zap.L().Error(
			"[ChapterSvc.CanListChapter] failed to verify member",
			zap.String("teamId", teamId),
			zap.String("currUid", currUid),
			zap.Error(err),
		)

		return classifyRepoErr(err, clsf, svc_res.Forbidden, "仅汉化组成员可查看该漫画的章节", "获取成员信息失败", "成员信息服务暂不可用", "获取成员信息失败")
	}
	if !ok {
		return svc_res.Reject(svc_res.Forbidden, "无权访问该漫画的章节")
	}

	return svc_res.Accept()
}

func (ChapterSvc) CanAdminChapter(currUid string, teamId string, memberRepo repo_iface.MemberRepo, clsf repo_iface.ErrClsf) svc_res.SvcRes {
	member, err := memberRepo.GetByUserTeamId(currUid, teamId)
	if err != nil {
		zap.L().Error(
			"[ChapterSvc.CanAdminChapter] failed to verify member",
			zap.String("teamId", teamId),
			zap.String("currUid", currUid),
			zap.Error(err),
		)

		return classifyRepoErr(err, clsf, svc_res.Forbidden, "仅汉化组管理员可操作该漫画的章节", "获取成员信息失败", "成员信息服务暂不可用", "权限校验失败")
	}
	if member == nil || !member.HasAnyRole(enum.RoleAdmin) {
		return svc_res.Reject(svc_res.Forbidden, "仅汉化组管理员可操作该漫画的章节")
	}

	return svc_res.Accept()
}

// `NewChapterCre` builds a chapter creation payload with generated id.
func (ChapterSvc) NewChapterCre(
	comicId string,
	index int,
	subtitle *string,
	creatorId string,
) *aggr.ChapterCre {
	return &aggr.ChapterCre{
		Id:        util.GenId("chapter"),
		ComicId:   comicId,
		Index:     index,
		Subtitle:  subtitle,
		CreatorId: creatorId,
	}
}

// `DefSubtitle` returns default subtitle for chapter index.
func (ChapterSvc) DefSubtitle(index int) string {
	return fmt.Sprintf("Ch.%d", index)
}
