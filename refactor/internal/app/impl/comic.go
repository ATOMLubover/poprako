package app_impl

import (
	"context"

	app_iface "poprako-s/internal/app"
	app_res "poprako-s/internal/app/res"
	app_util "poprako-s/internal/app/util"
	"poprako-s/internal/app/val"
	oss_iface "poprako-s/internal/domain/ext/oss"
	"poprako-s/internal/domain/model/aggr"
	"poprako-s/internal/domain/model/enum"
	"poprako-s/internal/domain/model/event"
	"poprako-s/internal/domain/model/query"
	repo_iface "poprako-s/internal/domain/repo"
	"poprako-s/internal/domain/svc"
	event_iface "poprako-s/internal/event"
	repo_infra "poprako-s/internal/infra/repo"

	"go.uber.org/zap"
)

// `comicAppImpl` is the default implementation of `ComicApp`
type comicAppImpl struct {
	txnCtrl repo_iface.TxnCtrl

	comicSvc      svc.ComicSvc
	chapterSvc    svc.ChapterSvc
	assignmentSvc svc.AssignmentSvc

	memberRepo  repo_iface.MemberRepo
	worksetRepo repo_iface.WorksetRepo
	comicRepo   repo_iface.ComicRepo
	ossSigner   oss_iface.Signer

	evBus   event_iface.EvBus
	errClsf repo_iface.ErrClsf
}

// `NewComicApp` creates a `ComicApp` implementation
func NewComicApp(
	txnCtrl repo_iface.TxnCtrl,
	memberRepo repo_iface.MemberRepo,
	worksetRepo repo_iface.WorksetRepo,
	comicRepo repo_iface.ComicRepo,
	ossSigner oss_iface.Signer,
	comicSvc svc.ComicSvc,
	chapterSvc svc.ChapterSvc,
	assignmentSvc svc.AssignmentSvc,
	evBus event_iface.EvBus,
	errClsf repo_iface.ErrClsf,
) app_iface.ComicApp {
	if txnCtrl == nil ||
		memberRepo == nil ||
		worksetRepo == nil ||
		comicRepo == nil ||
		ossSigner == nil ||
		evBus == nil ||
		errClsf == nil {
		zap.L().Panic(
			"[NewComicApp] nil dependency",
			zap.Bool("txnCtrl", txnCtrl == nil),
			zap.Bool("memberRepo", memberRepo == nil),
			zap.Bool("worksetRepo", worksetRepo == nil),
			zap.Bool("comicRepo", comicRepo == nil),
			zap.Bool("ossSigner", ossSigner == nil),
			zap.Bool("evBus", evBus == nil),
			zap.Bool("errClsf", errClsf == nil),
		)
	}

	return &comicAppImpl{
		txnCtrl:       txnCtrl,
		comicSvc:      comicSvc,
		chapterSvc:    chapterSvc,
		assignmentSvc: assignmentSvc,
		memberRepo:    memberRepo,
		worksetRepo:   worksetRepo,
		comicRepo:     comicRepo,
		ossSigner:     ossSigner,
		evBus:         evBus,
		errClsf:       errClsf,
	}
}

// `List` returns all comics for one workset
func (a *comicAppImpl) List(cx context.Context, currUid string, args *val.ListComicArgs) app_res.AppRes[[]val.ComicVal] {
	lgr := app_util.TakeLgr(cx)

	if re := vfyListComicArgs(args); re.IsReject() {
		return app_res.Reject[[]val.ComicVal](re.Code(), re.Msg())
	}

	workset, err := a.worksetRepo.GetById(args.WorksetId)
	if err != nil {
		return app_res.Reject[[]val.ComicVal](app_res.BadRequest, "作品集不存在")
	}

	if re := a.comicSvc.CanListComic(currUid, workset.TeamId, a.memberRepo, a.errClsf); re.IsReject() {
		return app_res.Reject[[]val.ComicVal](app_res.ErrCode(re.Code()), re.Msg())
	}

	listOpt := &query.ListComicOpt{
		WorksetId:      &args.WorksetId,
		UploadPhase:    args.UploadPhase,
		TranslatePhase: args.TranslatePhase,
		ProofreadPhase: args.ProofreadPhase,
		TypesetPhase:   args.TypesetPhase,
		ReviewPhase:    args.ReviewPhase,
		PublishPhase:   args.PublishPhase,
		Pagi: query.PagiOpt{
			Offset: args.Offset,
			Limit:  args.Limit,
		},
	}

	if args.FuzzyTitle != "" {
		listOpt.FuzzyTitle = &args.FuzzyTitle
	}

	comics, err := a.comicRepo.List(listOpt)
	if err != nil {
		lgr.Error(
			"[comicAppImpl.List] failed to list comics",
			zap.Error(err),
		)

		return app_res.Reject[[]val.ComicVal](app_res.ServerError, "获取漫画列表失败")
	}

	comicVals := make([]val.ComicVal, len(comics))
	for i, comic := range comics {
		comicVals[i] = asmComicVal(comic)
	}

	return app_res.Accept(&comicVals)
}

// `Create` creates a comic in target workset and auto-creates its first chapter.
func (a *comicAppImpl) Create(cx context.Context, currUid string, args *val.CreateComicArgs) app_res.AppRes[val.ComicCreatedRes] {
	lgr := app_util.TakeLgr(cx)

	if re := vfyCreateComicArgs(args); re.IsReject() {
		return app_res.Reject[val.ComicCreatedRes](re.Code(), re.Msg())
	}

	ev := make([]event_iface.Event, 0)

	re, err := repo_iface.RunWithTxn[app_res.AppRes[val.ComicCreatedRes]](a.txnCtrl, func(prov repo_iface.Prov) (app_res.AppRes[val.ComicCreatedRes], error) {
		memberRepo := prov.MemberRepo()
		worksetRepo := prov.WorksetRepo()
		comicRepo := prov.ComicRepo()

		workset, err := worksetRepo.GetById(args.WorksetId)
		if err != nil {
			return app_res.Reject[val.ComicCreatedRes](app_res.Forbidden, "仅汉化组管理员可创建漫画"), app_res.DefErr()
		}

		if re := a.comicSvc.CanAdminComic(currUid, workset.TeamId, memberRepo, a.errClsf); re.IsReject() {
			return app_res.Reject[val.ComicCreatedRes](app_res.ErrCode(re.Code()), re.Msg()), app_res.DefErr()
		}

		index, err := worksetRepo.IncrementComicNextIndex(args.WorksetId)
		if err != nil {
			return app_res.Reject[val.ComicCreatedRes](app_res.ServerError, "创建漫画失败"), err
		}

		cre := a.comicSvc.NewComicCre(
			args.WorksetId,
			index,
			args.Title,
			args.Author,
			args.Desc,
			currUid,
		)

		comic, err := comicRepo.Create(cre)
		if err != nil {
			return app_res.Reject[val.ComicCreatedRes](app_res.ServerError, "创建漫画失败"), err
		}

		if err := worksetRepo.UpdateComicCount(args.WorksetId, 1); err != nil {
			return app_res.Reject[val.ComicCreatedRes](app_res.ServerError, "创建漫画失败"), err
		}

		// Create first chapter with allocated chapter index.
		chapterIndex, err := comicRepo.IncrementChapterNextIndex(comic.Id)
		if err != nil {
			return app_res.Reject[val.ComicCreatedRes](app_res.ServerError, "创建漫画失败"), err
		}

		chapterCre := a.chapterSvc.NewChapterCre(comic.Id, chapterIndex, nil, currUid)
		chapter, err := prov.ChapterRepo().Create(chapterCre)
		if err != nil {
			return app_res.Reject[val.ComicCreatedRes](app_res.ServerError, "创建漫画失败"), err
		}

		if err := comicRepo.UpdateChapterCount(comic.Id, 1); err != nil {
			return app_res.Reject[val.ComicCreatedRes](app_res.ServerError, "创建漫画失败"), err
		}

		if err := comicRepo.TouchLastActive(comic.Id); err != nil {
			return app_res.Reject[val.ComicCreatedRes](app_res.ServerError, "创建漫画失败"), err
		}

		reviewerCre := a.assignmentSvc.NewAssignmentCre(chapter.Id, currUid, aggr.RoleMask(enum.RoleReviewer))
		if _, err := prov.AssignmentRepo().Create(reviewerCre); err != nil {
			return app_res.Reject[val.ComicCreatedRes](app_res.ServerError, "创建漫画失败"), err
		}

		ev = append(ev, event.NewAssignmentCreatedEv(currUid, chapter.Id))

		return app_res.Accept(&val.ComicCreatedRes{Id: comic.Id}), nil
	})
	if err != nil {
		lgr.Error(
			"[comicAppImpl.Create] failed to run create comic transaction",
			zap.Error(err),
		)

		return re
	}

	a.evBus.Pub(context.Background(), ev)

	return re
}

// `Update` updates title author and description of a comic
func (a *comicAppImpl) Update(cx context.Context, currUid string, args *val.ComicUpdArgs) app_res.AppRes[app_res.None] {
	lgr := app_util.TakeLgr(cx)

	if re := vfyUpdateComicArgs(args); re.IsReject() {
		return app_res.Reject[app_res.None](re.Code(), re.Msg())
	}

	comic, err := a.comicRepo.GetById(args.Id)
	if err != nil {
		lgr.Error("[comicAppImpl.Update] failed to get comic", zap.Error(err))

		return app_res.Reject[app_res.None](app_res.BadRequest, "漫画不存在")
	}

	workset, err := a.worksetRepo.GetById(comic.WorksetId)
	if err != nil {
		lgr.Error("[comicAppImpl.Update] failed to get workset", zap.Error(err))

		return app_res.Reject[app_res.None](app_res.ServerError, "更新漫画失败")
	}

	if re := a.comicSvc.CanAdminComic(currUid, workset.TeamId, a.memberRepo, a.errClsf); re.IsReject() {
		return app_res.Reject[app_res.None](app_res.ErrCode(re.Code()), re.Msg())
	}

	upd := &aggr.ComicUpd{
		Id:     args.Id,
		Title:  args.Title,
		Author: args.Author,
		Desc:   args.Desc,
	}

	if err := a.comicRepo.Update(upd); err != nil {
		lgr.Error("[comicAppImpl.Update] failed to update comic", zap.Error(err))

		return app_res.Reject[app_res.None](app_res.ServerError, "更新漫画失败")
	}

	return app_res.Accept(&app_res.None{})
}

// `GetById` returns one comic by id
func (a *comicAppImpl) GetById(cx context.Context, currUid string, comicId string) app_res.AppRes[val.ComicVal] {
	lgr := app_util.TakeLgr(cx)

	if re := vfyComicId(comicId); re.IsReject() {
		return app_res.Reject[val.ComicVal](re.Code(), re.Msg())
	}

	comic, err := a.comicRepo.GetById(comicId)
	if err != nil {
		if repo_infra.IsNotFound(err) {
			return app_res.Reject[val.ComicVal](app_res.NotFound, "漫画不存在")
		}

		lgr.Error("[comicAppImpl.GetById] failed to get comic", zap.Error(err))

		return app_res.Reject[val.ComicVal](app_res.ServerError, "获取漫画失败")
	}

	workset, err := a.worksetRepo.GetById(comic.WorksetId)
	if err != nil {
		lgr.Error("[comicAppImpl.GetById] failed to get workset", zap.Error(err))

		return app_res.Reject[val.ComicVal](app_res.ServerError, "获取漫画失败")
	}

	if re := a.comicSvc.CanListComic(currUid, workset.TeamId, a.memberRepo, a.errClsf); re.IsReject() {
		return app_res.Reject[val.ComicVal](app_res.ErrCode(re.Code()), re.Msg())
	}

	comicVal := asmComicVal(comic)

	if comic.CoverUploaded && comic.CoverKey != nil && *comic.CoverKey != "" {
		coverUrl, err := a.ossSigner.GenGetUrl(*comic.CoverKey)
		if err != nil {
			lgr.Error("[comicAppImpl.GetById] failed to generate comic cover url", zap.Error(err))
		} else {
			comicVal.CoverUrl = coverUrl
		}
	}

	return app_res.Accept(&comicVal)
}

// `ResvCover` reserves signed upload url for one comic cover.
func (a *comicAppImpl) ResvCover(cx context.Context, currUid string, args *val.ResvComicCoverArgs) app_res.AppRes[val.ResvComicCoverRes] {
	lgr := app_util.TakeLgr(cx)

	if re := vfyResvComicCoverArgs(args); re.IsReject() {
		return app_res.Reject[val.ResvComicCoverRes](re.Code(), re.Msg())
	}

	newKey := "comic_cover/" + args.ComicId + "." + args.FileExt

	if re, err := repo_iface.RunWithTxn[app_res.AppRes[val.ResvComicCoverRes]](a.txnCtrl, func(prov repo_iface.Prov) (app_res.AppRes[val.ResvComicCoverRes], error) {
		memberRepo := prov.MemberRepo()
		worksetRepo := prov.WorksetRepo()
		comicRepo := prov.ComicRepo()
		ossMsgRepo := prov.OssMsgRepo()
		ossMsgSvc := svc.NewOssMsgSvc()

		comic, err := comicRepo.GetById(args.ComicId)
		if err != nil {
			return app_res.Reject[val.ResvComicCoverRes](app_res.BadRequest, "漫画不存在"), app_res.DefErr()
		}

		workset, err := worksetRepo.GetById(comic.WorksetId)
		if err != nil {
			return app_res.Reject[val.ResvComicCoverRes](app_res.Forbidden, "仅汉化组管理员可预留漫画封面"), app_res.DefErr()
		}

		if re := a.comicSvc.CanAdminComic(currUid, workset.TeamId, memberRepo, a.errClsf); re.IsReject() {
			return app_res.Reject[val.ResvComicCoverRes](app_res.ErrCode(re.Code()), re.Msg()), app_res.DefErr()
		}

		oldKey := ""
		if comic.CoverKey != nil {
			oldKey = *comic.CoverKey
		}

		if err := comicRepo.PrefillCoverKey(args.ComicId, newKey); err != nil {
			return app_res.Reject[val.ResvComicCoverRes](app_res.ServerError, "预留漫画封面失败"), err
		}

		if oldKey != "" && oldKey != newKey {
			if err := ossMsgSvc.SavePendingDel(ossMsgRepo, enum.OssResComicCover, args.ComicId, []string{oldKey}); err != nil {
				return app_res.Reject[val.ResvComicCoverRes](app_res.ServerError, "预留漫画封面失败"), err
			}
		}

		if err := ossMsgSvc.SavePendingCre(ossMsgRepo, enum.OssResComicCover, args.ComicId, []string{newKey}); err != nil {
			return app_res.Reject[val.ResvComicCoverRes](app_res.ServerError, "预留漫画封面失败"), err
		}

		return app_res.Accept(&val.ResvComicCoverRes{}), nil
	}); err != nil {
		lgr.Error("[comicAppImpl.ResvCover] failed to reserve comic cover", zap.Error(err))

		return re
	}

	putUrl, err := a.ossSigner.GenPutUrl(newKey)
	if err != nil {
		lgr.Error("[comicAppImpl.ResvCover] failed to generate cover upload url", zap.Error(err))

		return app_res.Reject[val.ResvComicCoverRes](app_res.ServerError, "预留漫画封面失败")
	}

	return app_res.Accept(&val.ResvComicCoverRes{PutUrl: putUrl})
}

// `MarkCoverUploaded` confirms one comic cover upload.
func (a *comicAppImpl) MarkCoverUploaded(cx context.Context, currUid string, comicId string) app_res.AppRes[app_res.None] {
	lgr := app_util.TakeLgr(cx)

	if comicId == "" {
		return app_res.Reject[app_res.None](app_res.BadRequest, "comic_id 不能为空")
	}

	if re, err := repo_iface.RunWithTxn[app_res.AppRes[app_res.None]](a.txnCtrl, func(prov repo_iface.Prov) (app_res.AppRes[app_res.None], error) {
		memberRepo := prov.MemberRepo()
		worksetRepo := prov.WorksetRepo()
		comicRepo := prov.ComicRepo()
		ossMsgRepo := prov.OssMsgRepo()

		comic, err := comicRepo.GetById(comicId)
		if err != nil {
			if repo_infra.IsNotFound(err) {
				return app_res.Reject[app_res.None](app_res.NotFound, "漫画不存在"), app_res.DefErr()
			}

			return app_res.Reject[app_res.None](app_res.ServerError, "确认漫画封面上传失败"), err
		}

		workset, err := worksetRepo.GetById(comic.WorksetId)
		if err != nil {
			return app_res.Reject[app_res.None](app_res.Forbidden, "仅汉化组管理员可确认漫画封面"), app_res.DefErr()
		}

		if re := a.comicSvc.CanAdminComic(currUid, workset.TeamId, memberRepo, a.errClsf); re.IsReject() {
			return app_res.Reject[app_res.None](app_res.ErrCode(re.Code()), re.Msg()), app_res.DefErr()
		}

		if err := comicRepo.MarkCoverUploaded(comicId); err != nil {
			return app_res.Reject[app_res.None](app_res.ServerError, "确认漫画封面上传失败"), err
		}

		if err := ossMsgRepo.MarkCompletedByRes(enum.OssResComicCover, comicId); err != nil {
			if repo_infra.IsNotFound(err) {
				return app_res.Reject[app_res.None](app_res.BadRequest, "无效的漫画封面上传状态"), app_res.DefErr()
			}

			return app_res.Reject[app_res.None](app_res.ServerError, "确认漫画封面上传失败"), err
		}

		return app_res.Accept(&app_res.None{}), nil
	}); err != nil {
		lgr.Error("[comicAppImpl.MarkCoverUploaded] failed to confirm comic cover upload", zap.Error(err))

		return re
	}

	return app_res.Accept(&app_res.None{})
}

// `Delete` hard-deletes a comic and its descendants.
func (a *comicAppImpl) Delete(cx context.Context, currUid string, comicId string) app_res.AppRes[app_res.None] {
	lgr := app_util.TakeLgr(cx)

	if re := vfyDeleteComicId(comicId); re.IsReject() {
		return app_res.Reject[app_res.None](re.Code(), re.Msg())
	}

	ev := make([]event_iface.Event, 0)

	re, err := repo_iface.RunWithTxn[app_res.AppRes[app_res.None]](a.txnCtrl, func(prov repo_iface.Prov) (app_res.AppRes[app_res.None], error) {
		// Load all repo dependencies required by the cascade delete flow.
		memberRepo := prov.MemberRepo()
		worksetRepo := prov.WorksetRepo()
		comicRepo := prov.ComicRepo()
		chapterRepo := prov.ChapterRepo()
		pageRepo := prov.PageRepo()
		assignmentRepo := prov.AssignmentRepo()
		ossMsgRepo := prov.OssMsgRepo()
		ossMsgSvc := svc.NewOssMsgSvc()

		// Load the target comic and verify admin permission first.
		comic, err := comicRepo.GetById(comicId)
		if err != nil {
			return app_res.Reject[app_res.None](app_res.Forbidden, "仅汉化组管理员可删除漫画"), app_res.DefErr()
		}

		workset, err := worksetRepo.GetById(comic.WorksetId)
		if err != nil {
			return app_res.Reject[app_res.None](app_res.Forbidden, "仅汉化组管理员可删除漫画"), app_res.DefErr()
		}

		if re := a.comicSvc.CanAdminComic(currUid, workset.TeamId, memberRepo, a.errClsf); re.IsReject() {
			return app_res.Reject[app_res.None](app_res.ErrCode(re.Code()), re.Msg()), app_res.DefErr()
		}

		// Load all chapters explicitly so the cascade flow is complete for any valid cardinality.
		chapters, err := listAllChapters(chapterRepo, comic.Id)
		if err != nil {
			return app_res.Reject[app_res.None](app_res.ServerError, "删除漫画失败"), err
		}

		for i := range chapters {
			// Delete one chapter subtree and preserve chapter-removed event payload.
			assignedUserIds, err := deleteChapterCascade(
				chapters[i],
				pageRepo,
				assignmentRepo,
				chapterRepo,
				ossMsgRepo,
				ossMsgSvc,
			)
			if err != nil {
				return app_res.Reject[app_res.None](app_res.ServerError, "删除漫画失败"), err
			}

			ev = append(ev, event.NewChapterRemovedEv(chapters[i].Id, chapters[i].PublishedAt != nil, assignedUserIds))
		}

		// Enqueue comic-cover cleanup before deleting the comic row.
		if comic.CoverKey != nil && *comic.CoverKey != "" {
			if err := ossMsgSvc.SavePendingDel(ossMsgRepo, enum.OssResComicCover, comic.Id, []string{*comic.CoverKey}); err != nil {
				return app_res.Reject[app_res.None](app_res.ServerError, "删除漫画失败"), err
			}
		}

		// Delete the comic row after all descendants have been removed.
		if err := comicRepo.Delete(comicId); err != nil {
			return app_res.Reject[app_res.None](app_res.ServerError, "删除漫画失败"), err
		}

		// Decrement the workset comic counter after the comic row is removed.
		if err := worksetRepo.UpdateComicCount(workset.Id, -1); err != nil {
			return app_res.Reject[app_res.None](app_res.ServerError, "删除漫画失败"), err
		}

		return app_res.Accept(&app_res.None{}), nil
	})
	if err != nil {
		lgr.Error(
			"[comicAppImpl.Delete] failed to run delete comic transaction",
			zap.Error(err),
		)

		return re
	}

	a.evBus.Pub(context.Background(), ev)

	return re
}
