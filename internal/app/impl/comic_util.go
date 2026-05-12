package app_impl

import (
	"strings"

	app_res "poprako-s/internal/app/res"
	app_util "poprako-s/internal/app/util"
	"poprako-s/internal/app/val"
	oss_iface "poprako-s/internal/domain/ext/oss"
	"poprako-s/internal/domain/model/aggr"
	"poprako-s/internal/domain/model/enum"
	repo_iface "poprako-s/internal/domain/repo"

	"go.uber.org/zap"
)

// `asmComicVal` converts a `Comic` aggregate to app-facing `ComicVal`
func asmComicVal(comic *aggr.Comic) val.ComicVal {
	var worksetVal *val.WorksetVal
	var creatorVal *val.UserVal

	if comic.Workset != nil {
		v := asmWorksetVal(comic.Workset)

		worksetVal = &v
	}

	if comic.Creator != nil {
		creatorVal = &val.UserVal{
			Id:             comic.Creator.Id,
			Qid:            comic.Creator.Qid,
			Nickname:       comic.Creator.Nickname,
			AvatarUploaded: comic.Creator.AvatarUploaded,
			IsSuperAdmin:   comic.Creator.IsSuperAdmin,
			LastActiveAt:   comic.Creator.LastActiveAt.UnixMilli(),
			CreatedAt:      comic.Creator.CreatedAt.UnixMilli(),
			UpdatedAt:      comic.Creator.UpdatedAt.UnixMilli(),
		}
	}

	return val.ComicVal{
		Id:            comic.Id,
		WorksetId:     comic.WorksetId,
		Workset:       worksetVal,
		Index:         comic.Index,
		Title:         comic.Title,
		Author:        comic.Author,
		Desc:          comic.Desc,
		IsCompleted:   comic.IsCompleted,
		CoverUploaded: comic.CoverUploaded,
		ChapterCount:  comic.ChapterCount,
		CreatorId:     comic.CreatorId,
		Creator:       creatorVal,
		LastActiveAt:  comic.LastActiveAt.UnixMilli(),
		CreatedAt:     comic.CreatedAt.UnixMilli(),
		UpdatedAt:     comic.UpdatedAt.UnixMilli(),
	}
}

// `resolveComicCoverUrl` resolves one cover url with comic-cover and pinned-page fallback.
func resolveComicCoverUrl(
	comic *aggr.Comic,
	signer oss_iface.Signer,
	chapterRepo repo_iface.ChapterRepo,
	pageRepo repo_iface.PageRepo,
	lgr *zap.Logger,
) (string, repo_iface.RepoErr) {
	// Prefer the uploaded comic cover key when available.
	if comic.CoverUploaded && comic.CoverKey != nil && *comic.CoverKey != "" {
		coverUrl, err := signer.GenGetUrl(*comic.CoverKey)
		if err != nil {
			lgr.Error("[resolveComicCoverUrl] failed to generate comic cover url", zap.String("comicId", comic.Id), zap.Error(err))

			return "", nil
		}

		return coverUrl, nil
	}

	// Find pinned chapter for fallback source.
	pinnedChapter, err := chapterRepo.FindPinnedByComicId(comic.Id)
	if err != nil {
		return "", err
	}

	if pinnedChapter == nil {
		lgr.Warn("[resolveComicCoverUrl] pinned chapter not found", zap.String("comicId", comic.Id))

		return "", nil
	}

	// Read the first page under pinned chapter as fallback candidate.
	firstPages, err := pageRepo.FindFirstPageByChapters([]string{pinnedChapter.Id})
	if err != nil {
		return "", err
	}

	if len(firstPages) == 0 || firstPages[0] == nil {
		lgr.Warn("[resolveComicCoverUrl] pinned chapter has no pages", zap.String("chapterId", pinnedChapter.Id))

		return "", nil
	}

	firstPage := firstPages[0]
	if !firstPage.ImageUploaded || firstPage.ImageKey == nil || *firstPage.ImageKey == "" {
		lgr.Warn("[resolveComicCoverUrl] pinned first page image not uploaded", zap.String("pageId", firstPage.Id))

		return "", nil
	}

	// Generate get url for pinned first page image.
	coverUrl, err := signer.GenGetUrl(*firstPage.ImageKey)
	if err != nil {
		lgr.Error("[resolveComicCoverUrl] failed to generate fallback page image url", zap.String("pageId", firstPage.Id), zap.Error(err))

		return "", nil
	}

	return coverUrl, nil
}

// `resolveComicFallbackCoverUrls` resolves fallback urls for many comics by batch queries.
func resolveComicFallbackCoverUrls(
	comicIds []string,
	chapterRepo repo_iface.ChapterRepo,
	pageRepo repo_iface.PageRepo,
	signer oss_iface.Signer,
	lgr *zap.Logger,
) (map[string]string, repo_iface.RepoErr) {
	if len(comicIds) == 0 {
		return map[string]string{}, nil
	}

	// Batch load pinned chapters aligned to `comicIds`.
	pinnedChapters, err := chapterRepo.FindPinnedByComics(comicIds)
	if err != nil {
		return nil, err
	}

	chapterIds := make([]string, 0, len(pinnedChapters))
	chapterByComicId := make(map[string]*aggr.Chapter, len(pinnedChapters))
	for i := range pinnedChapters {
		if pinnedChapters[i] == nil {
			lgr.Warn("[resolveComicFallbackCoverUrls] pinned chapter not found", zap.String("comicId", comicIds[i]))

			continue
		}

		chapterByComicId[comicIds[i]] = pinnedChapters[i]
		chapterIds = append(chapterIds, pinnedChapters[i].Id)
	}

	if len(chapterIds) == 0 {
		return map[string]string{}, nil
	}

	// Batch load first pages aligned to `chapterIds`.
	firstPages, err := pageRepo.FindFirstPageByChapters(chapterIds)
	if err != nil {
		return nil, err
	}

	pageByChapterId := make(map[string]*aggr.Page, len(firstPages))
	for i := range firstPages {
		pageByChapterId[chapterIds[i]] = firstPages[i]
	}

	coverUrls := make(map[string]string, len(comicIds))
	for i := range comicIds {
		pinnedChapter := chapterByComicId[comicIds[i]]
		if pinnedChapter == nil {
			continue
		}

		firstPage := pageByChapterId[pinnedChapter.Id]
		if firstPage == nil {
			lgr.Warn("[resolveComicFallbackCoverUrls] pinned chapter has no pages", zap.String("chapterId", pinnedChapter.Id))

			continue
		}

		if !firstPage.ImageUploaded || firstPage.ImageKey == nil || *firstPage.ImageKey == "" {
			lgr.Warn("[resolveComicFallbackCoverUrls] pinned first page image not uploaded", zap.String("pageId", firstPage.Id))

			continue
		}

		coverUrl, err := signer.GenGetUrl(*firstPage.ImageKey)
		if err != nil {
			lgr.Error("[resolveComicFallbackCoverUrls] failed to generate fallback page image url", zap.String("pageId", firstPage.Id), zap.Error(err))

			continue
		}

		coverUrls[comicIds[i]] = coverUrl
	}

	return coverUrls, nil
}

// `mapComicFallbackErrCode` maps fallback-query repo errors to app-level error code.
func mapComicFallbackErrCode(err repo_iface.RepoErr, errClsf repo_iface.ErrClsf) app_res.ErrCode {
	if errClsf.IsTimeout(err) {
		return app_res.Timeout
	}

	if errClsf.IsUnavailable(err) {
		return app_res.Unavailable
	}

	return app_res.ServerError
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

// `vfyDeleteComicId` validates delete arguments.
func vfyDeleteComicId(comicId string) app_res.AppRes[app_res.None] {
	if comicId == "" {
		return app_res.Reject[app_res.None](app_res.BadRequest, "comic_id 不能为空")
	}

	return app_res.Accept(&app_res.None{})
}

// `vfyResvComicCoverArgs` validates reserve-cover arguments.
func vfyResvComicCoverArgs(args *val.ResvComicCoverArgs) app_res.AppRes[app_res.None] {
	if args == nil {
		return app_res.Reject[app_res.None](app_res.BadRequest, "预留参数不能为空")
	}

	if args.ComicId == "" || args.FileExt == "" {
		return app_res.Reject[app_res.None](app_res.BadRequest, "comic_id file_extension 不能为空")
	}

	return app_res.Accept(&app_res.None{})
}
