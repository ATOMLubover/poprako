package app_impl

import (
	"strings"

	app_res "poprako-s/internal/app/res"
	app_util "poprako-s/internal/app/util"
	"poprako-s/internal/app/val"
	"poprako-s/internal/domain/model/aggr"
	"poprako-s/internal/domain/model/enum"
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
func vfyListComicArgs(args *val.ListComicArgs) app_res.AppRes[app_res.None] {
	if args == nil {
		return app_res.Reject[app_res.None](app_res.BadRequest, "分页参数不能为空")
	}

	if args.WorksetId == "" {
		return app_res.Reject[app_res.None](app_res.BadRequest, "workset_id 不能为空")
	}

	if re := app_util.ClampOffsetLimit(args.Offset, &args.Limit); re.IsReject() {
		return re
	}

	args.FuzzyTitle = strings.TrimSpace(args.FuzzyTitle)

	if re := vfyWorkflowPhase(args.UploadPhase, "upload_phase"); re.IsReject() {
		return re
	}

	if re := vfyWorkflowPhase(args.TranslatePhase, "translate_phase"); re.IsReject() {
		return re
	}

	if re := vfyWorkflowPhase(args.ProofreadPhase, "proofread_phase"); re.IsReject() {
		return re
	}

	if re := vfyWorkflowPhase(args.TypesetPhase, "typeset_phase"); re.IsReject() {
		return re
	}

	if re := vfyWorkflowPhase(args.ReviewPhase, "review_phase"); re.IsReject() {
		return re
	}

	if re := vfyWorkflowPhase(args.PublishPhase, "publish_phase"); re.IsReject() {
		return re
	}

	return app_res.Accept(&app_res.None{})
}

// `vfyWorkflowPhase` validates one workflow phase value.
func vfyWorkflowPhase(phase *enum.WorkflowPhase, field string) app_res.AppRes[app_res.None] {
	if phase == nil {
		return app_res.Accept(&app_res.None{})
	}

	if *phase < enum.WorkflowPending || *phase > enum.WorkflowCompleted {
		return app_res.Reject[app_res.None](app_res.BadRequest, field+" 参数不合法")
	}

	return app_res.Accept(&app_res.None{})
}

// `vfyCreateComicArgs` validates create arguments.
func vfyCreateComicArgs(args *val.CreateComicArgs) app_res.AppRes[app_res.None] {
	if args == nil {
		return app_res.Reject[app_res.None](app_res.BadRequest, "创建参数不能为空")
	}

	if args.WorksetId == "" || args.Title == "" || args.Author == "" {
		return app_res.Reject[app_res.None](app_res.BadRequest, "workset_id title author 不能为空")
	}

	return app_res.Accept(&app_res.None{})
}

// `vfyUpdateComicArgs` validates update arguments.
func vfyUpdateComicArgs(args *val.ComicUpdArgs) app_res.AppRes[app_res.None] {
	if args == nil {
		return app_res.Reject[app_res.None](app_res.BadRequest, "更新参数不能为空")
	}

	if args.Id == "" || args.Title == "" || args.Author == "" {
		return app_res.Reject[app_res.None](app_res.BadRequest, "id title author 不能为空")
	}

	return app_res.Accept(&app_res.None{})
}

// `vfyComicId` validates comic id.
func vfyComicId(comicId string) app_res.AppRes[app_res.None] {
	if comicId == "" {
		return app_res.Reject[app_res.None](app_res.BadRequest, "comic_id 不能为空")
	}

	return app_res.Accept(&app_res.None{})
}

// `vfyRemoveComicId` validates remove arguments.
func vfyRemoveComicId(comicId string) app_res.AppRes[app_res.None] {
	if comicId == "" {
		return app_res.Reject[app_res.None](app_res.BadRequest, "comic_id 不能为空")
	}

	return app_res.Accept(&app_res.None{})
}
