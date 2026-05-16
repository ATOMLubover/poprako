package svc

import (
	"time"

	"poprako-s/internal/domain/model/aggr"
	"poprako-s/internal/domain/model/enum"
	repo_iface "poprako-s/internal/domain/repo"
	svc_res "poprako-s/internal/domain/svc/res"
	"poprako-s/pkg/util"
)

// `AssignmentSvc` provides stateless domain services for assignment.
type AssignmentSvc struct{}

// `NewAssignmentSvc` returns a ready-to-use `AssignmentSvc`.
func NewAssignmentSvc() AssignmentSvc {
	return AssignmentSvc{}
}

func (AssignmentSvc) CanReviewAssignment(currUid string, chapterId string, assignmentRepo repo_iface.AssignmentRepo, clsf repo_iface.ErrClsf) svc_res.SvcRes {
	currAssignment, err := assignmentRepo.GetByChapterUserId(chapterId, currUid)
	if err != nil {
		return classifyRepoErr(err, clsf, svc_res.Forbidden, "仅章节监修可执行该操作", "权限校验超时", "权限校验服务暂不可用", "权限校验失败")
	}
	if currAssignment == nil || !currAssignment.HasAnyRole(enum.RoleReviewer) {
		return svc_res.Reject(svc_res.Forbidden, "仅章节监修可执行该操作")
	}

	return svc_res.Accept()
}

// `CanSelfReduceAssignment` allows a user to reduce or remove their own assignment roles.
// Self-assignment (adding new roles to oneself) is rejected — that requires a reviewer.
func (AssignmentSvc) CanSelfReduceAssignment(currUid string, targetUserId string, chapterId string, newMask aggr.RoleMask, assignmentRepo repo_iface.AssignmentRepo, clsf repo_iface.ErrClsf) svc_res.SvcRes {
	if currUid != targetUserId {
		return svc_res.Reject(svc_res.Forbidden, "仅可对自己的分配执行此操作")
	}

	curr, err := assignmentRepo.GetByChapterUserId(chapterId, targetUserId)
	if err != nil {
		return classifyRepoErr(err, clsf, svc_res.Forbidden, "仅可缩减或退出自己的分配", "权限校验超时", "权限校验服务暂不可用", "权限校验失败")
	}

	if curr == nil {
		if newMask == 0 {
			return svc_res.Accept()
		}
		return svc_res.Reject(svc_res.Forbidden, "仅章节监修可分配新角色")
	}

	currMask := curr.ToRoleMask()
	if newMask&^currMask != 0 {
		return svc_res.Reject(svc_res.Forbidden, "仅可缩减或退出自己的分配，不可新增角色")
	}

	return svc_res.Accept()
}

// `CanListByChapter` validates whether caller can list assignments under one chapter.
// Legacy-compatible rule: allow either team member access or chapter-assignment fallback access.
func (s AssignmentSvc) CanListByChapter(
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

		return classifyRepoErr(err, clsf, svc_res.Forbidden, "无权查看该章节的分配列表", "权限校验超时", "权限校验服务暂不可用", "权限校验失败")
	}

	if ok {
		return svc_res.Accept()
	}

	if re := s.canListByChapterFallback(currUid, chapterId, assignmentRepo, clsf); !re.IsReject() {
		return re
	}

	return svc_res.Reject(svc_res.Forbidden, "无权查看该章节的分配列表")
}

// `canListByChapterFallback` validates legacy fallback access by existing chapter assignment.
func (AssignmentSvc) canListByChapterFallback(currUid string, chapterId string, assignmentRepo repo_iface.AssignmentRepo, clsf repo_iface.ErrClsf) svc_res.SvcRes {
	assignment, err := assignmentRepo.GetByChapterUserId(chapterId, currUid)
	if err != nil {
		return classifyRepoErr(err, clsf, svc_res.Forbidden, "无权查看该章节的分配列表", "权限校验超时", "权限校验服务暂不可用", "权限校验失败")
	}

	if assignment == nil {
		return svc_res.Reject(svc_res.Forbidden, "无权查看该章节的分配列表")
	}

	return svc_res.Accept()
}

func (AssignmentSvc) CanTakeAssignmentRoles(userId string, chapterId string, roleMask aggr.RoleMask, memberRepo repo_iface.MemberRepo, chapterRepo repo_iface.ChapterRepo, comicRepo repo_iface.ComicRepo, worksetRepo repo_iface.WorksetRepo, clsf repo_iface.ErrClsf) svc_res.SvcRes {
	if roleMask == 0 {
		return svc_res.Accept()
	}

	if roleMask.HasAnyRole(enum.RoleAdmin) {
		return svc_res.Reject(svc_res.BadRequest, "章节分配不支持管理员角色")
	}

	ch, err := chapterRepo.GetById(chapterId)
	if err != nil {
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

	member, err := memberRepo.GetByUserTeamId(userId, ws.TeamId)
	if err != nil {
		return classifyRepoErr(err, clsf, svc_res.BadRequest, "目标成员不在当前汉化组中", "成员信息查询超时", "成员信息服务暂不可用", "成员信息查询失败")
	}

	if member == nil {
		return svc_res.Reject(svc_res.BadRequest, "目标成员不在当前汉化组中")
	}

	for _, role := range roleMask.ToRoleArr() {
		if !member.HasAnyRole(role) {
			return svc_res.Reject(svc_res.BadRequest, "目标成员无法承担指定章节角色")
		}
	}

	return svc_res.Accept()
}

// `NewAssignmentCre` builds one assignment create payload from role mask.
func (AssignmentSvc) NewAssignmentCre(chapterId string, userId string, roleMask aggr.RoleMask) *aggr.AssignmentCre {
	now := time.Now()
	roles := mkTimedRolesFromMask(roleMask, now)

	return &aggr.AssignmentCre{
		Id:         util.GenId("assignment"),
		ChapterId:  chapterId,
		UserId:     userId,
		TimedRoles: roles,
	}
}

// `NewAssignmentPut` builds put payload that preserves existing role timestamps.
func (AssignmentSvc) NewAssignmentPut(curr *aggr.Assignment, roleMask aggr.RoleMask) *aggr.AssignmentPut {
	now := time.Now()
	roles := mkTimedRolesFromCurr(curr, roleMask, now)

	return &aggr.AssignmentPut{
		Id:         curr.Id,
		TimedRoles: roles,
	}
}

func mkTimedRolesFromMask(mask aggr.RoleMask, now time.Time) aggr.TimedRoles {
	roles := aggr.TimedRoles{}

	if mask.HasAnyRole(enum.RoleRawProvider) {
		t := now
		roles.AssignedRawProviderAt = &t
	}
	if mask.HasAnyRole(enum.RoleTranslator) {
		t := now
		roles.AssignedTranslatorAt = &t
	}
	if mask.HasAnyRole(enum.RoleProofreader) {
		t := now
		roles.AssignedProofreaderAt = &t
	}
	if mask.HasAnyRole(enum.RoleTypesetter) {
		t := now
		roles.AssignedTypesetterAt = &t
	}
	if mask.HasAnyRole(enum.RoleRedrawer) {
		t := now
		roles.AssignedRedrawerAt = &t
	}
	if mask.HasAnyRole(enum.RoleReviewer) {
		t := now
		roles.AssignedReviewerAt = &t
	}
	if mask.HasAnyRole(enum.RolePublisher) {
		t := now
		roles.AssignedPublisherAt = &t
	}

	return roles
}

func mkTimedRolesFromCurr(curr *aggr.Assignment, mask aggr.RoleMask, now time.Time) aggr.TimedRoles {
	keepOrNew := func(has bool, prev *time.Time) *time.Time {
		if !has {
			return nil
		}

		if prev != nil {
			t := *prev
			return &t
		}

		t := now
		return &t
	}

	return aggr.TimedRoles{
		AssignedRawProviderAt: keepOrNew(mask.HasAnyRole(enum.RoleRawProvider), curr.AssignedRawProviderAt),
		AssignedTranslatorAt:  keepOrNew(mask.HasAnyRole(enum.RoleTranslator), curr.AssignedTranslatorAt),
		AssignedProofreaderAt: keepOrNew(mask.HasAnyRole(enum.RoleProofreader), curr.AssignedProofreaderAt),
		AssignedTypesetterAt:  keepOrNew(mask.HasAnyRole(enum.RoleTypesetter), curr.AssignedTypesetterAt),
		AssignedRedrawerAt:    keepOrNew(mask.HasAnyRole(enum.RoleRedrawer), curr.AssignedRedrawerAt),
		AssignedReviewerAt:    keepOrNew(mask.HasAnyRole(enum.RoleReviewer), curr.AssignedReviewerAt),
		AssignedPublisherAt:   keepOrNew(mask.HasAnyRole(enum.RolePublisher), curr.AssignedPublisherAt),
	}
}
