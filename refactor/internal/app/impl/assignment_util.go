package app_impl

import (
	"poprako-s/internal/app/res"
	app_util "poprako-s/internal/app/util"
	"poprako-s/internal/app/val"
	"poprako-s/internal/domain/model/aggr"
	"poprako-s/internal/domain/model/enum"
	"poprako-s/internal/domain/model/query"
	repo_iface "poprako-s/internal/domain/repo"
)




func asmAssignmentVal(a *aggr.Assignment) val.AssignmentVal {
	return val.AssignmentVal{
		Id:                    a.Id,
		ChapterId:             a.ChapterId,
		UserId:                a.UserId,
		RoleMask:              a.ToRoleMask(),
		AssignedRawProviderAt: toUnixMilliPtr(a.AssignedRawProviderAt),
		AssignedTranslatorAt:  toUnixMilliPtr(a.AssignedTranslatorAt),
		AssignedProofreaderAt: toUnixMilliPtr(a.AssignedProofreaderAt),
		AssignedTypesetterAt:  toUnixMilliPtr(a.AssignedTypesetterAt),
		AssignedRedrawerAt:    toUnixMilliPtr(a.AssignedRedrawerAt),
		AssignedReviewerAt:    toUnixMilliPtr(a.AssignedReviewerAt),
		AssignedPublisherAt:   toUnixMilliPtr(a.AssignedPublisherAt),
		CreatedAt:             a.CreatedAt.UnixMilli(),
		UpdatedAt:             a.UpdatedAt.UnixMilli(),
	}
}

func asmAssignmentInvVal(a *aggr.AssignmentInv) val.AssignmentInvVal {
	return val.AssignmentInvVal{
		Id:         a.Id,
		ChapterId:  a.ChapterId,
		InviterId:  a.InviterId,
		InviteeQid: a.InviteeQid,
		InvCode:    a.InvCode,
		Pending:    a.Pending,
		RoleMask:   a.RoleMask,
		CreatedAt:  a.CreatedAt.UnixMilli(),
		UpdatedAt:  a.UpdatedAt.UnixMilli(),
	}
}


func vfyAssignmentListArgs(offset int, limit *int) (res.ErrCode, string, bool) {
	return app_util.ClampOffsetLimit(offset, limit)
}

func ensureReviewerPermission(assignmentRepo repo_iface.AssignmentRepo, chapterId string, currUid string) (res.ErrCode, string, bool) {
	currAssignment, err := assignmentRepo.GetByChapterUserId(chapterId, currUid)
	if err != nil {
		return app_util.ClassifyRepoErr(err, res.Forbidden, "仅章节监修可执行该操作", "权限校验超时", "权限校验服务暂不可用", "权限校验失败")
	}

	if currAssignment == nil || !currAssignment.HasAnyRole(enum.RoleReviewer) {
		return res.Forbidden, "仅章节监修可执行该操作", true
	}

	return 0, "", false
}

func ensureUserCanTakeRoles(memberRepo repo_iface.MemberRepo, chapterRepo repo_iface.ChapterRepo, comicRepo repo_iface.ComicRepo, worksetRepo repo_iface.WorksetRepo, chapterId string, userId string, roleMask aggr.RoleMask) (res.ErrCode, string, bool) {
	if roleMask == 0 {
		return 0, "", false
	}

	if roleMask.HasAnyRole(enum.RoleAdmin) {
		return res.BadRequest, "章节分配不支持管理员角色", true
	}

	ch, err := chapterRepo.GetById(chapterId)
	if err != nil {
		return app_util.ClassifyRepoErr(err, res.BadRequest, "章节不存在", "章节信息查询超时", "章节信息服务暂不可用", "章节信息查询失败")
	}

	cm, err := comicRepo.GetById(ch.ComicId)
	if err != nil {
		return app_util.ClassifyRepoErr(err, res.BadRequest, "章节不存在", "章节信息查询超时", "章节信息服务暂不可用", "章节信息查询失败")
	}

	ws, err := worksetRepo.GetById(cm.WorksetId)
	if err != nil {
		return app_util.ClassifyRepoErr(err, res.BadRequest, "章节不存在", "章节信息查询超时", "章节信息服务暂不可用", "章节信息查询失败")
	}

	member, err := memberRepo.GetByUserTeamId(userId, ws.TeamId)
	if err != nil {
		return app_util.ClassifyRepoErr(err, res.BadRequest, "目标成员不在当前汉化组中", "成员信息查询超时", "成员信息服务暂不可用", "成员信息查询失败")
	}

	if member == nil {
		return res.BadRequest, "目标成员不在当前汉化组中", true
	}

	for _, role := range roleMask.ToRoleArr() {
		if role == enum.RoleAdmin {
			return res.BadRequest, "章节分配不支持管理员角色", true
		}

		if !member.HasAnyRole(role) {
			return res.BadRequest, "目标成员无法承担指定章节角色", true
		}
	}

	return 0, "", false
}

func mkListAssignmentOptByChapter(chapterId string, offset int, limit int) *query.ListAssignmentOpt {
	return &query.ListAssignmentOpt{
		ChapterId: &chapterId,
		Pagi: query.PagiOpt{
			Offset: offset,
			Limit:  limit,
		},
	}
}

func mkListAssignmentOptByUser(userId string, offset int, limit int) *query.ListAssignmentOpt {
	return &query.ListAssignmentOpt{
		UserId: &userId,
		Pagi: query.PagiOpt{
			Offset: offset,
			Limit:  limit,
		},
	}
}
