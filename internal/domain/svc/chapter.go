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

// `CanTransiteWorkflow` validates that the current user's chapter-level assignment
// authorises the requested forward workflow transition
// `RoleReviewer` can execute any transition
// Other roles can only execute their corresponding transitions (see `CanTransiteWorkflow` doc in `docs/perm-list.md`)
// For revert transitions, use `CanRevertWorkflow` instead
func (ChapterSvc) CanTransiteWorkflow(currUid string, chapterId string, t enum.WorkflowTransition, assignmentRepo repo_iface.AssignmentRepo, clsf repo_iface.ErrClsf) svc_res.SvcRes {
	assignment, err := assignmentRepo.GetByChapterUserId(chapterId, currUid)
	if err != nil {
		return classifyRepoErr(err, clsf, svc_res.Forbidden, "仅章节成员可推进工作流", "权限校验超时", "权限校验服务暂不可用", "权限校验失败")
	}

	if assignment == nil {
		return svc_res.Reject(svc_res.Forbidden, "仅章节成员可推进工作流")
	}

	// `RoleReviewer` can execute any workflow transition.
	if assignment.HasAnyRole(enum.RoleReviewer) {
		return svc_res.Accept()
	}

	// Map transition to required chapter-level role.
	switch t {
	case enum.WorkflowUploadComplete:
		if !assignment.HasAnyRole(enum.RoleRawProvider) {
			return svc_res.Reject(svc_res.Forbidden, "只有图源可以标记上传完成")
		}

	case enum.WorkflowTranslateStart, enum.WorkflowTranslateComplete:
		if !assignment.HasAnyRole(enum.RoleTranslator) {
			return svc_res.Reject(svc_res.Forbidden, "只有翻译可以标记翻译开始或完成")
		}

	case enum.WorkflowProofreadStart, enum.WorkflowProofreadComplete:
		if !assignment.HasAnyRole(enum.RoleProofreader) {
			return svc_res.Reject(svc_res.Forbidden, "只有校对可以标记校对开始或完成")
		}

	case enum.WorkflowTypesetStart, enum.WorkflowTypesetComplete:
		if !assignment.HasAnyRole(enum.RoleTypesetter) {
			return svc_res.Reject(svc_res.Forbidden, "只有嵌字可以标记嵌字开始或完成")
		}

	case enum.WorkflowReviewComplete:
		if !assignment.HasAnyRole(enum.RoleReviewer) {
			return svc_res.Reject(svc_res.Forbidden, "只有监修可以标记监修完成")
		}

	case enum.WorkflowPublishComplete:
		if !assignment.HasAnyRole(enum.RolePublisher) {
			return svc_res.Reject(svc_res.Forbidden, "只有发布可以标记发布完成")
		}
	}

	return svc_res.Accept()
}

// `CanRevertWorkflow` validates that the current user's chapter-level assignment
// authorises the requested revert transition
// Publish-complete cannot be reverted by anyone
// `RoleReviewer` can revert any revertible transition
// `RoleProofreader` can additionally revert translate-start and translate-complete
// Other roles can only revert their own phase transitions
// Permission matrix:
//   - `upload_revert`:          `RoleRawProvider`, `RoleReviewer`
//   - `translate_start_revert`: `RoleTranslator`, `RoleProofreader`, `RoleReviewer`
//   - `translate_revert`:       `RoleTranslator`, `RoleProofreader`, `RoleReviewer`
//   - `proofread_start_revert`: `RoleProofreader`, `RoleReviewer`
//   - `proofread_revert`:       `RoleProofreader`, `RoleReviewer`
//   - `typeset_start_revert`:   `RoleTypesetter`, `RoleReviewer`
//   - `typeset_revert`:         `RoleTypesetter`, `RoleReviewer`
//   - `review_revert`:          `RoleReviewer`
func (ChapterSvc) CanRevertWorkflow(currUid string, chapterId string, t enum.WorkflowTransition, assignmentRepo repo_iface.AssignmentRepo, clsf repo_iface.ErrClsf) svc_res.SvcRes {
	assignment, err := assignmentRepo.GetByChapterUserId(chapterId, currUid)
	if err != nil {
		return classifyRepoErr(err, clsf, svc_res.Forbidden, "仅章节成员可回退工作流", "权限校验超时", "权限校验服务暂不可用", "权限校验失败")
	}

	if assignment == nil {
		return svc_res.Reject(svc_res.Forbidden, "仅章节成员可回退工作流")
	}

	// `RoleReviewer` can revert any transition
	if assignment.HasAnyRole(enum.RoleReviewer) {
		return svc_res.Accept()
	}

	// Map revert transition to required chapter-level role
	switch t {
	case enum.WorkflowUploadRevert:
		if !assignment.HasAnyRole(enum.RoleRawProvider) {
			return svc_res.Reject(svc_res.Forbidden, "只有图源可以回退上传完成")
		}

	case enum.WorkflowTranslateStartRevert, enum.WorkflowTranslateRevert:
		if !assignment.HasAnyRole(enum.RoleTranslator, enum.RoleProofreader) {
			return svc_res.Reject(svc_res.Forbidden, "只有翻译或校对可以回退翻译进度")
		}

	case enum.WorkflowProofreadStartRevert, enum.WorkflowProofreadRevert:
		if !assignment.HasAnyRole(enum.RoleProofreader) {
			return svc_res.Reject(svc_res.Forbidden, "只有校对可以回退校对进度")
		}

	case enum.WorkflowTypesetStartRevert, enum.WorkflowTypesetRevert:
		if !assignment.HasAnyRole(enum.RoleTypesetter) {
			return svc_res.Reject(svc_res.Forbidden, "只有嵌字可以回退嵌字进度")
		}

	case enum.WorkflowReviewRevert:
		// Already handled above — only `RoleReviewer` reaches here and was already accepted
		return svc_res.Reject(svc_res.Forbidden, "只有监修可以回退监修进度")
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
	return fmt.Sprintf("第%d话", index+1)
}
