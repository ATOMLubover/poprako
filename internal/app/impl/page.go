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
	"poprako-s/internal/domain/model/query"
	repo_iface "poprako-s/internal/domain/repo"
	"poprako-s/internal/domain/svc"
	repo_infra "poprako-s/internal/infra/repo"

	"go.uber.org/zap"
)

// `pageAppImpl` is the default implementation of `PageApp`.
type pageAppImpl struct {
	// `txnCtrl` controls application transactions.
	txnCtrl repo_iface.TxnCtrl

	// `pageSvc` provides page-specific business helpers.
	pageSvc svc.PageSvc
	// `chapterSvc` provides chapter-level permission checks.
	chapterSvc svc.ChapterSvc
	// `ossMsgSvc` builds pending OSS create messages.
	ossMsgSvc svc.OssMsgSvc

	// `memberRepo` reads member information for team permission checks.
	memberRepo repo_iface.MemberRepo
	// `worksetRepo` reads workset information for chapter ownership traversal.
	worksetRepo repo_iface.WorksetRepo
	// `comicRepo` reads comic information for chapter ownership traversal.
	comicRepo repo_iface.ComicRepo
	// `chapterRepo` reads and updates chapter counters.
	chapterRepo repo_iface.ChapterRepo
	// `pageRepo` reads and writes page rows.
	pageRepo repo_iface.PageRepo
	// `assignmentRepo` reads chapter assignment information.
	assignmentRepo repo_iface.AssignmentRepo

	// `ossSigner` generates signed OSS urls.
	ossSigner oss_iface.Signer

	// `errClsf` classifies repository errors for permission services.
	errClsf repo_iface.ErrClsf
}

// `pageResvHolder` carries one reserved page id and image key across commit boundary.
type pageResvHolder struct {
	// `PageId` is the reserved page identifier.
	PageId string

	// `ImageKey` is the reserved page image key.
	ImageKey string
}

// `NewPageApp` creates one `PageApp` implementation.
func NewPageApp(
	txnCtrl repo_iface.TxnCtrl,
	memberRepo repo_iface.MemberRepo,
	worksetRepo repo_iface.WorksetRepo,
	comicRepo repo_iface.ComicRepo,
	chapterRepo repo_iface.ChapterRepo,
	pageRepo repo_iface.PageRepo,
	assignmentRepo repo_iface.AssignmentRepo,
	pageSvc svc.PageSvc,
	chapterSvc svc.ChapterSvc,
	ossMsgSvc svc.OssMsgSvc,
	ossSigner oss_iface.Signer,
	errClsf repo_iface.ErrClsf,
) app_iface.PageApp {
	if txnCtrl == nil || memberRepo == nil || worksetRepo == nil || comicRepo == nil || chapterRepo == nil || pageRepo == nil || assignmentRepo == nil || ossSigner == nil || errClsf == nil {
		zap.L().Panic(
			"[NewPageApp] nil dependency",
			zap.Bool("txnCtrl", txnCtrl == nil),
			zap.Bool("memberRepo", memberRepo == nil),
			zap.Bool("worksetRepo", worksetRepo == nil),
			zap.Bool("comicRepo", comicRepo == nil),
			zap.Bool("chapterRepo", chapterRepo == nil),
			zap.Bool("pageRepo", pageRepo == nil),
			zap.Bool("assignmentRepo", assignmentRepo == nil),
			zap.Bool("ossSigner", ossSigner == nil),
			zap.Bool("errClsf", errClsf == nil),
		)
	}

	return &pageAppImpl{
		txnCtrl:        txnCtrl,
		pageSvc:        pageSvc,
		chapterSvc:     chapterSvc,
		ossMsgSvc:      ossMsgSvc,
		memberRepo:     memberRepo,
		worksetRepo:    worksetRepo,
		comicRepo:      comicRepo,
		chapterRepo:    chapterRepo,
		pageRepo:       pageRepo,
		assignmentRepo: assignmentRepo,
		ossSigner:      ossSigner,
		errClsf:        errClsf,
	}
}

// `ResvChapterPages` reserves upload slots for chapter pages.
func (a *pageAppImpl) ResvChapterPages(cx context.Context, currUid string, args *val.ResvChapterPagesArgs) app_res.AppRes[val.ResvChapterPagesRes] {
	lgr := app_util.TakeLgr(cx)

	// Validate the incoming reservation request before touching storage.
	if re := vfyResvChapterPagesArgs(args); re.IsReject() {
		return app_res.Reject[val.ResvChapterPagesRes](re.Code(), re.Msg())
	}

	holders := make([]pageResvHolder, 0, args.PageCount)

	// Reserve pages, enqueue OSS create messages, and update chapter page count atomically.
	re, err := repo_iface.RunWithTxn[app_res.AppRes[val.ResvChapterPagesRes]](a.txnCtrl, func(prov repo_iface.Prov) (app_res.AppRes[val.ResvChapterPagesRes], error) {
		chapterRepo := prov.ChapterRepo()
		comicRepo := prov.ComicRepo()
		pageRepo := prov.PageRepo()
		assignmentRepo := prov.AssignmentRepo()
		ossMsgRepo := prov.OssMsgRepo()

		chapter, err := chapterRepo.GetById(args.ChapterId)
		if err != nil {
			if repo_infra.IsNotFound(err) {
				return app_res.Reject[val.ResvChapterPagesRes](app_res.BadRequest, "章节不存在"), app_res.DefErr()
			}

			return app_res.Reject[val.ResvChapterPagesRes](app_res.ServerError, "预留页面失败"), err
		}

		if re := a.pageSvc.CanResvPages(currUid, chapter.Id, assignmentRepo, a.errClsf); re.IsReject() {
			return app_res.Reject[val.ResvChapterPagesRes](app_res.ErrCode(re.Code()), re.Msg()), app_res.DefErr()
		}

		if chapter.PageCount != 0 {
			return app_res.Reject[val.ResvChapterPagesRes](app_res.Conflict, "当前章节已存在页面，不允许再次上传"), app_res.DefErr()
		}

		txnHolders := make([]pageResvHolder, 0, args.PageCount)
		batch := make([]*aggr.PageCre, 0, args.PageCount)
		for i := 0; i < args.PageCount; i++ {
			cre := a.pageSvc.NewPageCre(chapter.Id, i, nil)
			imageKey := a.pageSvc.GenImageKey(chapter.Id, cre.Id, args.FileExt)
			cre.ImageKey = &imageKey

			batch = append(batch, cre)
			txnHolders = append(txnHolders, pageResvHolder{PageId: cre.Id, ImageKey: imageKey})
		}

		if err := pageRepo.CreateBatch(batch); err != nil {
			return app_res.Reject[val.ResvChapterPagesRes](app_res.ServerError, "预留页面失败"), err
		}

		for i := range txnHolders {
			if err := a.ossMsgSvc.SavePendingCre(ossMsgRepo, enum.OssResPageImage, txnHolders[i].PageId, []string{txnHolders[i].ImageKey}); err != nil {
				return app_res.Reject[val.ResvChapterPagesRes](app_res.ServerError, "预留页面失败"), err
			}
		}

		if err := chapterRepo.SetPageCount(chapter.Id, len(txnHolders)); err != nil {
			return app_res.Reject[val.ResvChapterPagesRes](app_res.ServerError, "预留页面失败"), err
		}

		if err := comicRepo.TouchLastActive(chapter.ComicId); err != nil {
			return app_res.Reject[val.ResvChapterPagesRes](app_res.ServerError, "预留页面失败"), err
		}

		holders = txnHolders

		return app_res.Accept(&val.ResvChapterPagesRes{}), nil
	})
	if err != nil {
		lgr.Error("[pageAppImpl.ResvChapterPages] failed to reserve chapter pages", zap.Error(err))

		return re
	}

	// Generate signed upload urls only after the database transaction commits successfully.
	creations := make([]val.PageCreationRes, 0, len(holders))
	for i := range holders {
		putUrl, err := a.ossSigner.GenPutUrl(holders[i].ImageKey)
		if err != nil {
			lgr.Error("[pageAppImpl.ResvChapterPages] failed to generate page upload url", zap.String("page_id", holders[i].PageId), zap.Error(err))

			return app_res.Reject[val.ResvChapterPagesRes](app_res.ServerError, "生成页面上传地址失败")
		}

		creations = append(creations, val.PageCreationRes{PageId: holders[i].PageId, PutUrl: putUrl})
	}

	return app_res.Accept(&val.ResvChapterPagesRes{Creations: creations})
}

// `List` returns pages under one chapter.
func (a *pageAppImpl) List(cx context.Context, currUid string, args *val.ListChapterPageArgs) app_res.AppRes[[]val.PageVal] {
	lgr := app_util.TakeLgr(cx)

	// Validate and normalize list arguments before querying repositories.
	if re := vfyListChapterPageArgs(args); re.IsReject() {
		return app_res.Reject[[]val.PageVal](re.Code(), re.Msg())
	}

	// Verify caller access via legacy-compatible permission path in domain service.

	if re := a.pageSvc.CanListByChapter(
		currUid,
		args.ChapterId,
		a.memberRepo,
		a.worksetRepo,
		a.comicRepo,
		a.chapterRepo,
		a.assignmentRepo,
		a.errClsf,
	); re.IsReject() {
		return app_res.Reject[[]val.PageVal](app_res.ErrCode(re.Code()), re.Msg())
	}

	// Load pages and assemble app-facing values with signed image urls.
	pages, err := a.pageRepo.List(&query.ListPageOpt{ChapterId: &args.ChapterId, Pagi: query.PagiOpt{Offset: args.Offset, Limit: args.Limit}})
	if err != nil {
		lgr.Error("[pageAppImpl.List] failed to list pages", zap.Error(err))

		return app_res.Reject[[]val.PageVal](app_res.ServerError, "获取页面列表失败")
	}

	pageVals := make([]val.PageVal, 0, len(pages))
	for i := range pages {
		pageVal, err := asmPageVal(pages[i], a.ossSigner)
		if err != nil {
			lgr.Error("[pageAppImpl.List] failed to assemble page value", zap.String("page_id", pages[i].Id), zap.Error(err))

			return app_res.Reject[[]val.PageVal](app_res.ServerError, "获取页面列表失败")
		}

		pageVals = append(pageVals, *pageVal)
	}

	return app_res.Accept(&pageVals)
}

// `MarkImageUploaded` confirms one page image upload.
func (a *pageAppImpl) MarkImageUploaded(cx context.Context, currUid string, args *val.MarkPageImageUploadedArgs) app_res.AppRes[app_res.None] {
	lgr := app_util.TakeLgr(cx)

	// Validate the target page id before starting the transaction.
	if re := vfyMarkPageImageUploadedArgs(args); re.IsReject() {
		return app_res.Reject[app_res.None](re.Code(), re.Msg())
	}

	// Mark the page as uploaded and complete the corresponding OSS create message atomically.
	re, err := repo_iface.RunWithTxn[app_res.AppRes[app_res.None]](a.txnCtrl, func(prov repo_iface.Prov) (app_res.AppRes[app_res.None], error) {
		pageRepo := prov.PageRepo()
		assignmentRepo := prov.AssignmentRepo()
		ossMsgRepo := prov.OssMsgRepo()

		page, err := pageRepo.GetById(args.PageId)
		if err != nil {
			if repo_infra.IsNotFound(err) {
				return app_res.Reject[app_res.None](app_res.NotFound, "页面不存在"), app_res.DefErr()
			}

			return app_res.Reject[app_res.None](app_res.ServerError, "确认页面上传失败"), err
		}

		if re := a.pageSvc.CanMarkImageUploaded(currUid, page.ChapterId, assignmentRepo, a.errClsf); re.IsReject() {
			return app_res.Reject[app_res.None](app_res.ErrCode(re.Code()), re.Msg()), app_res.DefErr()
		}

		if page.ImageUploaded {
			return app_res.Accept(&app_res.None{}), nil
		}

		if err := pageRepo.MarkImageUploaded(page.Id); err != nil {
			return app_res.Reject[app_res.None](app_res.ServerError, "确认页面上传失败"), err
		}

		if err := ossMsgRepo.MarkCompletedByRes(enum.OssResPageImage, page.Id); err != nil {
			if repo_infra.IsNotFound(err) {
				return app_res.Reject[app_res.None](app_res.BadRequest, "无效的页面上传状态"), app_res.DefErr()
			}

			return app_res.Reject[app_res.None](app_res.ServerError, "确认页面上传失败"), err
		}

		return app_res.Accept(&app_res.None{}), nil
	})
	if err != nil {
		lgr.Error("[pageAppImpl.MarkImageUploaded] failed to confirm page upload", zap.Error(err))

		return re
	}

	return re
}

// `DeleteByChapterId` deletes all pages under one chapter.
func (a *pageAppImpl) DeleteByChapterId(cx context.Context, currUid string, chapterId string) app_res.AppRes[app_res.None] {
	lgr := app_util.TakeLgr(cx)

	// Validate the target chapter id before starting destructive operations.
	if re := vfyDeleteByChapterId(chapterId); re.IsReject() {
		return app_res.Reject[app_res.None](re.Code(), re.Msg())
	}

	// Enqueue OSS delete messages, remove pages, and clear chapter page count atomically.
	re, err := repo_iface.RunWithTxn[app_res.AppRes[app_res.None]](a.txnCtrl, func(prov repo_iface.Prov) (app_res.AppRes[app_res.None], error) {
		memberRepo := prov.MemberRepo()
		worksetRepo := prov.WorksetRepo()
		comicRepo := prov.ComicRepo()
		chapterRepo := prov.ChapterRepo()
		pageRepo := prov.PageRepo()
		ossMsgRepo := prov.OssMsgRepo()

		chapter, err := chapterRepo.GetById(chapterId)
		if err != nil {
			if repo_infra.IsNotFound(err) {
				return app_res.Reject[app_res.None](app_res.BadRequest, "章节不存在"), app_res.DefErr()
			}

			return app_res.Reject[app_res.None](app_res.ServerError, "删除页面失败"), err
		}

		comic, err := comicRepo.GetById(chapter.ComicId)
		if err != nil {
			return app_res.Reject[app_res.None](app_res.Forbidden, "仅汉化组管理员可删除页面"), app_res.DefErr()
		}

		workset, err := worksetRepo.GetById(comic.WorksetId)
		if err != nil {
			return app_res.Reject[app_res.None](app_res.Forbidden, "仅汉化组管理员可删除页面"), app_res.DefErr()
		}

		if re := a.chapterSvc.CanAdminChapter(currUid, workset.TeamId, memberRepo, a.errClsf); re.IsReject() {
			return app_res.Reject[app_res.None](app_res.ErrCode(re.Code()), re.Msg()), app_res.DefErr()
		}

		limit := chapter.PageCount
		if limit <= 0 {
			limit = 10000
		}

		pages, err := pageRepo.List(&query.ListPageOpt{ChapterId: &chapterId, Pagi: query.PagiOpt{Limit: limit}})
		if err != nil {
			return app_res.Reject[app_res.None](app_res.ServerError, "删除页面失败"), err
		}

		ossMsgSvc := svc.NewOssMsgSvc()
		for i := range pages {
			if pages[i].ImageKey == nil || *pages[i].ImageKey == "" {
				continue
			}

			if err := ossMsgSvc.SavePendingDel(ossMsgRepo, enum.OssResPageImage, pages[i].Id, []string{*pages[i].ImageKey}); err != nil {
				return app_res.Reject[app_res.None](app_res.ServerError, "删除页面失败"), err
			}
		}

		if err := pageRepo.DeleteByChapterId(chapterId); err != nil {
			return app_res.Reject[app_res.None](app_res.ServerError, "删除页面失败"), err
		}

		if err := chapterRepo.SetPageCount(chapterId, 0); err != nil {
			return app_res.Reject[app_res.None](app_res.ServerError, "删除页面失败"), err
		}

		if err := comicRepo.TouchLastActive(chapter.ComicId); err != nil {
			return app_res.Reject[app_res.None](app_res.ServerError, "删除页面失败"), err
		}

		return app_res.Accept(&app_res.None{}), nil
	})
	if err != nil {
		lgr.Error("[pageAppImpl.DeleteByChapterId] failed to delete chapter pages", zap.Error(err))

		return re
	}

	return re
}
