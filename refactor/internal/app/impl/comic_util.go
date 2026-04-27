package app_impl

import (
	"strings"

	"poprako-s/internal/app/res"
	"poprako-s/internal/app/val"
	"poprako-s/internal/domain/model/aggr"
	"poprako-s/internal/domain/model/enum"
)

const (
	comicListDefLimit = 20
	comicListMaxLimit = 200
)

// `asmComicVal` converts a `Comic` aggregate to app-facing `ComicVal`
func asmComicVal(cm *aggr.Comic) val.ComicVal {
	return val.ComicVal{
		Id:           cm.Id,
		WorksetId:    cm.WorksetId,
		Index:        cm.Index,
		Title:        cm.Title,
		Author:       cm.Author,
		Desc:         cm.Desc,
		IsCompleted:  cm.IsCompleted,
		ChapterCount: cm.ChapterCount,
		CreatorId:    cm.CreatorId,
		LastActiveAt: cm.LastActiveAt.UnixMilli(),
		CreatedAt:    cm.CreatedAt.UnixMilli(),
		UpdatedAt:    cm.UpdatedAt.UnixMilli(),
	}
}

// `vfyListComicArgs` validates and normalizes list arguments.
func vfyListComicArgs(args *val.ListComicArgs) (res.ErrCode, string, bool) {
	if args == nil {
		return res.BadRequest, "分页参数不能为空", true
	}

	if args.WorksetId == "" {
		return res.BadRequest, "workset_id 不能为空", true
	}

	if args.Offset < 0 {
		return res.BadRequest, "offset 不能小于 0", true
	}

	if args.Limit <= 0 {
		args.Limit = comicListDefLimit
	}

	if args.Limit > comicListMaxLimit {
		args.Limit = comicListMaxLimit
	}

	args.FuzzyTitle = strings.TrimSpace(args.FuzzyTitle)

	if code, msg, reject := vfyWorkflowPhase(args.UploadPhase, "upload_phase"); reject {
		return code, msg, true
	}

	if code, msg, reject := vfyWorkflowPhase(args.TranslatePhase, "translate_phase"); reject {
		return code, msg, true
	}

	if code, msg, reject := vfyWorkflowPhase(args.ProofreadPhase, "proofread_phase"); reject {
		return code, msg, true
	}

	if code, msg, reject := vfyWorkflowPhase(args.TypesetPhase, "typeset_phase"); reject {
		return code, msg, true
	}

	if code, msg, reject := vfyWorkflowPhase(args.ReviewPhase, "review_phase"); reject {
		return code, msg, true
	}

	if code, msg, reject := vfyWorkflowPhase(args.PublishPhase, "publish_phase"); reject {
		return code, msg, true
	}

	return 0, "", false
}

// `vfyWorkflowPhase` validates one workflow phase value.
func vfyWorkflowPhase(phase *enum.WorkflowPhase, field string) (res.ErrCode, string, bool) {
	if phase == nil {
		return 0, "", false
	}

	if *phase < enum.WorkflowPending || *phase > enum.WorkflowCompleted {
		return res.BadRequest, field + " 参数不合法", true
	}

	return 0, "", false
}

// `vfyCreateComicArgs` validates create arguments.
func vfyCreateComicArgs(args *val.CreateComicArgs) (res.ErrCode, string, bool) {
	if args == nil {
		return res.BadRequest, "创建参数不能为空", true
	}

	if args.WorksetId == "" || args.Title == "" || args.Author == "" {
		return res.BadRequest, "workset_id title author 不能为空", true
	}

	return 0, "", false
}

// `vfyUpdateComicArgs` validates update arguments.
func vfyUpdateComicArgs(args *val.ComicUpdArgs) (res.ErrCode, string, bool) {
	if args == nil {
		return res.BadRequest, "更新参数不能为空", true
	}

	if args.Id == "" || args.Title == "" || args.Author == "" {
		return res.BadRequest, "id title author 不能为空", true
	}

	return 0, "", false
}

// `vfyRemoveComicId` validates remove arguments.
func vfyRemoveComicId(comicId string) (res.ErrCode, string, bool) {
	if comicId == "" {
		return res.BadRequest, "comic_id 不能为空", true
	}

	return 0, "", false
}
