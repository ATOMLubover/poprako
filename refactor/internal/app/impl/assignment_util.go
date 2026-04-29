package app_impl

import (
	app_res "poprako-s/internal/app/res"
	app_util "poprako-s/internal/app/util"
	"poprako-s/internal/app/val"
	"poprako-s/internal/domain/model/aggr"
	"poprako-s/internal/domain/model/query"
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

func vfyAssignmentListArgs(offset int, limit *int) app_res.AppRes[app_res.None] {
	return app_util.ClampOffsetLimit(offset, limit)
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
