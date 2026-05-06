package app_impl

import (
	app_res "poprako-s/internal/app/res"
	app_util "poprako-s/internal/app/util"
	"poprako-s/internal/app/val"
	"poprako-s/internal/domain/model/aggr"
	"poprako-s/internal/domain/model/enum"
	"poprako-s/internal/domain/model/query"
)

func asmAssignmentVal(a *aggr.Assignment) val.AssignmentVal {
	var userVal *val.UserVal
	if a.User != nil {
		userVal = &val.UserVal{
			Id:             a.User.Id,
			Qid:            a.User.Qid,
			Nickname:       a.User.Nickname,
			AvatarUploaded: a.User.AvatarUploaded,
			IsSuperAdmin:   a.User.IsSuperAdmin,
			LastActiveAt:   a.User.LastActiveAt.UnixMilli(),
			CreatedAt:      a.User.CreatedAt.UnixMilli(),
			UpdatedAt:      a.User.UpdatedAt.UnixMilli(),
		}
	}

	var chapterVal *val.ChapterVal
	if a.Chapter != nil {
		v := asmChapterVal(a.Chapter)

		chapterVal = &v
	}

	return val.AssignmentVal{
		Id:                    a.Id,
		ChapterId:             a.ChapterId,
		UserId:                a.UserId,
		Chapter:               chapterVal,
		User:                  userVal,
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

func mkListAssignmentOptByChapter(chapterId string, includes []enum.AssignmentIncl, offset int, limit int) *query.ListAssignmentOpt {
	return &query.ListAssignmentOpt{
		ChapterId: &chapterId,
		Includes:  includes,
		Pagi: query.PagiOpt{
			Offset: offset,
			Limit:  limit,
		},
	}
}

func mkListAssignmentOptByUser(userId string, includes []enum.AssignmentIncl, offset int, limit int) *query.ListAssignmentOpt {
	return &query.ListAssignmentOpt{
		UserId:   &userId,
		Includes: includes,
		Pagi: query.PagiOpt{
			Offset: offset,
			Limit:  limit,
		},
	}
}

func mkAssignmentRepoIncl(includes []enum.AssignmentIncl) []enum.AssignmentIncl {
	if len(includes) == 0 {
		return nil
	}

	repoIncls := make([]enum.AssignmentIncl, 0, len(includes))

	for i := range includes {
		repoIncls = append(repoIncls, includes[i])

		switch includes[i] {
		case enum.AssignmentInclChapterComic:
			repoIncls = append(repoIncls, enum.AssignmentInclChapter)

		case enum.AssignmentInclChapterComicWorkset:
			repoIncls = append(repoIncls, enum.AssignmentInclChapter, enum.AssignmentInclChapterComic)

		case enum.AssignmentInclChapterComicWorksetTeam:
			repoIncls = append(repoIncls, enum.AssignmentInclChapter, enum.AssignmentInclChapterComic, enum.AssignmentInclChapterComicWorkset)

		case enum.AssignmentInclChapterCreator:
			repoIncls = append(repoIncls, enum.AssignmentInclChapter)

		case enum.AssignmentInclChapterComicCreator:
			repoIncls = append(repoIncls, enum.AssignmentInclChapter, enum.AssignmentInclChapterComic)
		}
	}

	return repoIncls
}
