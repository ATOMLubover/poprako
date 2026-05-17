package svc

import (
	"poprako-s/internal/domain/model/aggr"
	"poprako-s/internal/domain/model/enum"
	repo_iface "poprako-s/internal/domain/repo"
	svc_res "poprako-s/internal/domain/svc/res"
	"poprako-s/pkg/util"
)

// `PageSvc` provides stateless helpers for page aggregate.
type PageSvc struct{}

// `NewPageSvc` creates one `PageSvc`.
func NewPageSvc() PageSvc {
	return PageSvc{}
}

// `CanResvPages` validates chapter-level permission for page reservation.
func (PageSvc) CanResvPages(currUid string, chapterId string, assignmentRepo repo_iface.AssignmentRepo, clsf repo_iface.ErrClsf) svc_res.SvcRes {
	assignment, err := assignmentRepo.GetByChapterUserId(chapterId, currUid)
	if err != nil {
		return classifyRepoErr(err, clsf, svc_res.Forbidden, "仅图源或监修可预留页面", "权限校验失败", "权限校验服务暂不可用", "权限校验失败")
	}
	if assignment == nil || !assignment.HasAnyRole(enum.RoleReviewer, enum.RoleRawProvider) {
		return svc_res.Reject(svc_res.Forbidden, "仅图源或监修可预留页面")
	}

	return svc_res.Accept()
}

// `CanListByChapter` validates whether caller can list pages under one chapter.
// Legacy-compatible rule: allow either team member access or chapter-assignment fallback access.
func (s PageSvc) CanListByChapter(
	currUid string,
	chapterId string,
	memberRepo repo_iface.MemberRepo,
	worksetRepo repo_iface.WorksetRepo,
	comicRepo repo_iface.ComicRepo,
	chapterRepo repo_iface.ChapterRepo,
	assignmentRepo repo_iface.AssignmentRepo,
	clsf repo_iface.ErrClsf,
) svc_res.SvcRes {
	ch, err := chapterRepo.GetById(chapterId)
	if err != nil {
		if re := s.canListByChapterFallback(currUid, chapterId, assignmentRepo, clsf); !re.IsReject() {
			return re
		}

		return classifyRepoErr(err, clsf, svc_res.BadRequest, "章节不存在", "章节信息查询超时", "章节信息服务暂不可用", "章节信息查询失败")
	}

	cm, err := comicRepo.GetById(ch.ComicId)
	if err != nil {
		return classifyRepoErr(err, clsf, svc_res.BadRequest, "章节不存在", "章节信息查询超时", "章节信息服务暂不可用", "章节信息查询失败")
	}

	ws, err := worksetRepo.GetById(cm.WorksetId)
	if err != nil {
		return classifyRepoErr(err, clsf, svc_res.BadRequest, "章节不存在", "章节信息查询超时", "章节信息服务暂不可用", "章节信息查询失败")
	}

	ok, err := memberRepo.ExistByUserTeamId(currUid, ws.TeamId)
	if err != nil {
		if re := s.canListByChapterFallback(currUid, chapterId, assignmentRepo, clsf); !re.IsReject() {
			return re
		}

		return classifyRepoErr(err, clsf, svc_res.Forbidden, "无权查看该章节的页面", "权限校验超时", "权限校验服务暂不可用", "权限校验失败")
	}

	if ok {
		return svc_res.Accept()
	}

	if re := s.canListByChapterFallback(currUid, chapterId, assignmentRepo, clsf); !re.IsReject() {
		return re
	}

	return svc_res.Reject(svc_res.Forbidden, "无权查看该章节的页面")
}

// `CanMarkImageUploaded` validates chapter-level permission for upload confirmation.
func (PageSvc) CanMarkImageUploaded(currUid string, chapterId string, assignmentRepo repo_iface.AssignmentRepo, clsf repo_iface.ErrClsf) svc_res.SvcRes {
	assignment, err := assignmentRepo.GetByChapterUserId(chapterId, currUid)
	if err != nil {
		return classifyRepoErr(err, clsf, svc_res.Forbidden, "仅图源可确认页面上传", "权限校验失败", "权限校验服务暂不可用", "权限校验失败")
	}
	if assignment == nil || !assignment.HasAnyRole(enum.RoleRawProvider) {
		return svc_res.Reject(svc_res.Forbidden, "仅图源可确认页面上传")
	}

	return svc_res.Accept()
}

// `canListByChapterFallback` validates legacy fallback access by existing chapter assignment.
func (PageSvc) canListByChapterFallback(currUid string, chapterId string, assignmentRepo repo_iface.AssignmentRepo, clsf repo_iface.ErrClsf) svc_res.SvcRes {
	assignment, err := assignmentRepo.GetByChapterUserId(chapterId, currUid)
	if err != nil {
		return classifyRepoErr(err, clsf, svc_res.Forbidden, "无权查看该章节的页面", "权限校验超时", "权限校验服务暂不可用", "权限校验失败")
	}

	if assignment == nil {
		return svc_res.Reject(svc_res.Forbidden, "无权查看该章节的页面")
	}

	return svc_res.Accept()
}

// `NewPageCre` creates one page create payload.
func (PageSvc) NewPageCre(chapterId string, index int, imageKey *string) *aggr.PageCre {
	return &aggr.PageCre{
		Id:        util.GenId("page"),
		ChapterId: chapterId,
		Index:     index,
		ImageKey:  imageKey,
	}
}
