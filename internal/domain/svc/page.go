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

// `NewPageCre` creates one page create payload.
func (PageSvc) NewPageCre(chapterId string, index int, imageKey *string) *aggr.PageCre {
	return &aggr.PageCre{
		Id:        util.GenId("page"),
		ChapterId: chapterId,
		Index:     index,
		ImageKey:  imageKey,
	}
}

// `GenImageKey` generates one page image OSS key.
func (PageSvc) GenImageKey(chapterId string, pageId string, fileExt string) string {
	return "chapter_" + chapterId + "/page_" + pageId + "." + fileExt
}
