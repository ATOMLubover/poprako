package app_impl

import (
	"strings"
	"time"

	app_res "poprako-s/internal/app/res"
	app_util "poprako-s/internal/app/util"
	"poprako-s/internal/app/val"
	"poprako-s/internal/domain/model/aggr"
	"poprako-s/internal/domain/model/enum"
)

// `asmChapterVal` converts chapter aggregate to app value object.
func asmChapterVal(chapter *aggr.Chapter) val.ChapterVal {
	var comicVal *val.ComicVal
	var creatorVal *val.UserVal

	if chapter.Comic != nil {
		v := asmComicVal(chapter.Comic)

		comicVal = &v
	}

	if chapter.Creator != nil {
		creatorVal = &val.UserVal{
			Id:             chapter.Creator.Id,
			Qid:            chapter.Creator.Qid,
			Nickname:       chapter.Creator.Nickname,
			AvatarUploaded: chapter.Creator.AvatarUploaded,
			IsSuperAdmin:   chapter.Creator.IsSuperAdmin,
			LastActiveAt:   chapter.Creator.LastActiveAt.UnixMilli(),
			CreatedAt:      chapter.Creator.CreatedAt.UnixMilli(),
			UpdatedAt:      chapter.Creator.UpdatedAt.UnixMilli(),
		}
	}

	return val.ChapterVal{
		Id:                  chapter.Id,
		ComicId:             chapter.ComicId,
		Comic:               comicVal,
		IsPinned:            chapter.IsPinned,
		Index:               chapter.Index,
		Subtitle:            chapter.Subtitle,
		PageCount:           chapter.PageCount,
		TotalUnitCount:      chapter.TotalUnitCount,
		TranslatedUnitCount: chapter.TranslatedUnitCount,
		ProofreadUnitCount:  chapter.ProofreadUnitCount,
		CreatorId:           chapter.CreatorId,
		Creator:             creatorVal,
		UploadedAt:          toUnixMilliPtr(chapter.UploadedAt),
		TransalatingAt:      toUnixMilliPtr(chapter.TransalatingAt),
		TranslatedAt:        toUnixMilliPtr(chapter.TranslatedAt),
		ProofreadingAt:      toUnixMilliPtr(chapter.ProofreadingAt),
		ProofreadAt:         toUnixMilliPtr(chapter.ProofreadAt),
		TypesettingAt:       toUnixMilliPtr(chapter.TypesettingAt),
		TypesetAt:           toUnixMilliPtr(chapter.TypesetAt),
		ReviewedAt:          toUnixMilliPtr(chapter.ReviewedAt),
		PublishedAt:         toUnixMilliPtr(chapter.PublishedAt),
		CreatedAt:           chapter.CreatedAt.UnixMilli(),
		UpdatedAt:           chapter.UpdatedAt.UnixMilli(),
	}
}

// `toUnixMilliPtr` converts time pointer to unix milli pointer.
func toUnixMilliPtr(t *time.Time) *int64 {
	if t == nil {
		return nil
	}

	unixMilli := t.UnixMilli()

	return &unixMilli
}

// `vfyListChapterArgs` validates and normalizes chapter list arguments.
func vfyListChapterArgs(args *val.ListChapterArgs) app_res.AppRes[app_res.None] {
	if args == nil {
		return app_res.Reject[app_res.None](app_res.BadRequest, "分页参数不能为空")
	}

	if args.ComicId == "" {
		return app_res.Reject[app_res.None](app_res.BadRequest, "comic_id 不能为空")
	}

	if re := app_util.ClampOffsetLimit(args.Offset, &args.Limit); re.IsReject() {
		return re
	}

	return app_res.Accept(&app_res.None{})
}

// `vfyCreateChapterArgs` validates chapter create arguments.
func vfyCreateChapterArgs(args *val.CreateChapterArgs) app_res.AppRes[app_res.None] {
	if args == nil {
		return app_res.Reject[app_res.None](app_res.BadRequest, "创建参数不能为空")
	}

	if args.ComicId == "" {
		return app_res.Reject[app_res.None](app_res.BadRequest, "comic_id 不能为空")
	}

	if args.Subtitle != nil {
		subtitle := strings.TrimSpace(*args.Subtitle)
		args.Subtitle = &subtitle
	}

	return app_res.Accept(&app_res.None{})
}

// `vfyUpdateChapterArgs` validates chapter update arguments.
func vfyUpdateChapterArgs(args *val.ChapterUpdArgs) app_res.AppRes[app_res.None] {
	if args == nil {
		return app_res.Reject[app_res.None](app_res.BadRequest, "更新参数不能为空")
	}

	if args.Id == "" {
		return app_res.Reject[app_res.None](app_res.BadRequest, "chapter_id 不能为空")
	}

	if args.Subtitle != nil {
		subtitle := strings.TrimSpace(*args.Subtitle)
		args.Subtitle = &subtitle
	}

	if args.WorkflowTransition != nil {
		if !isWorkflowTransitionValid(*args.WorkflowTransition) {
			return app_res.Reject[app_res.None](app_res.BadRequest, "workflow_transition 参数不合法")
		}
	}

	return app_res.Accept(&app_res.None{})
}

// `vfyChapterId` validates chapter id.
func vfyChapterId(chapterId string) app_res.AppRes[app_res.None] {
	if chapterId == "" {
		return app_res.Reject[app_res.None](app_res.BadRequest, "chapter_id 不能为空")
	}

	return app_res.Accept(&app_res.None{})
}

// `vfyDeleteChapterId` validates chapter delete id.
func vfyDeleteChapterId(chapterId string) app_res.AppRes[app_res.None] {
	if chapterId == "" {
		return app_res.Reject[app_res.None](app_res.BadRequest, "chapter_id 不能为空")
	}

	return app_res.Accept(&app_res.None{})
}

// `vfyJoinChapterArgs` validates chapter join arguments.
func vfyJoinChapterArgs(args val.JoinChapterArgs) app_res.AppRes[app_res.None] {
	if args.ChapterId == "" {
		return app_res.Reject[app_res.None](app_res.BadRequest, "chapter_id 不能为空")
	}

	if args.RoleMask == 0 {
		return app_res.Reject[app_res.None](app_res.BadRequest, "role_mask 不能为空")
	}

	return app_res.Accept(&app_res.None{})
}

// `isWorkflowTransitionValid` checks transition value in closed set.
func isWorkflowTransitionValid(t enum.WorkflowTransition) bool {
	switch t {
	case enum.WorkflowUploadComplete,
		enum.WorkflowTranslateStart,
		enum.WorkflowTranslateComplete,
		enum.WorkflowProofreadStart,
		enum.WorkflowProofreadComplete,
		enum.WorkflowTypesetStart,
		enum.WorkflowTypesetComplete,
		enum.WorkflowReviewComplete,
		enum.WorkflowPublishComplete:
		return true
	default:
		return false
	}
}

// `mkChapterUpd` builds chapter update payload from args and chapter state.
func mkChapterUpd(args *val.ChapterUpdArgs, ch *aggr.Chapter) *aggr.ChapterUpd {
	upd := &aggr.ChapterUpd{
		Id:                 args.Id,
		Subtitle:           args.Subtitle,
		IsPinned:           args.IsPinned,
		WorkflowTransition: args.WorkflowTransition,
	}

	if args.WorkflowTransition == nil {
		return upd
	}

	switch *args.WorkflowTransition {
	case enum.WorkflowUploadComplete:
		upd.UploadedAt = toTimePtrPtr(ch.UploadedAt)
	case enum.WorkflowTranslateStart:
		upd.TransalatingAt = toTimePtrPtr(ch.TransalatingAt)
	case enum.WorkflowTranslateComplete:
		upd.TranslatedAt = toTimePtrPtr(ch.TranslatedAt)
	case enum.WorkflowProofreadStart:
		upd.ProofreadingAt = toTimePtrPtr(ch.ProofreadingAt)
	case enum.WorkflowProofreadComplete:
		upd.ProofreadAt = toTimePtrPtr(ch.ProofreadAt)
	case enum.WorkflowTypesetStart:
		upd.TypesettingAt = toTimePtrPtr(ch.TypesettingAt)
	case enum.WorkflowTypesetComplete:
		upd.TypesetAt = toTimePtrPtr(ch.TypesetAt)
	case enum.WorkflowReviewComplete:
		upd.ReviewedAt = toTimePtrPtr(ch.ReviewedAt)
	case enum.WorkflowPublishComplete:
		upd.PublishedAt = toTimePtrPtr(ch.PublishedAt)
	}

	return upd
}

// `toTimePtrPtr` wraps one time pointer into pointer-to-pointer.
func toTimePtrPtr(t *time.Time) **time.Time {
	return &t
}
