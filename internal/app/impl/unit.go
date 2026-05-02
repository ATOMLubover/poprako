package app_impl

import (
	"context"

	app_iface "poprako-s/internal/app"
	app_res "poprako-s/internal/app/res"
	app_util "poprako-s/internal/app/util"
	"poprako-s/internal/app/val"
	repo_iface "poprako-s/internal/domain/repo"
	"poprako-s/internal/domain/svc"
	repo_infra "poprako-s/internal/infra/repo"

	"go.uber.org/zap"
)

// `unitAppImpl` is the default implementation of `UnitApp`.
type unitAppImpl struct {
	// `txnCtrl` controls application transactions.
	txnCtrl repo_iface.TxnCtrl

	// `unitSvc` provides unit-specific business helpers.
	unitSvc svc.UnitSvc
	// `comicRepo` reads comic information for ownership traversal.
	comicRepo repo_iface.ComicRepo
	// `chapterRepo` reads and updates chapter rows.
	chapterRepo repo_iface.ChapterRepo
	// `pageRepo` reads and updates page rows.
	pageRepo repo_iface.PageRepo
	// `unitRepo` reads and mutates unit rows.
	unitRepo repo_iface.UnitRepo
	// `assignmentRepo` reads assignment information for permission checks.
	assignmentRepo repo_iface.AssignmentRepo

	// `errClsf` classifies repository errors for permission services.
	errClsf repo_iface.ErrClsf
}

// `NewUnitApp` creates one `UnitApp` implementation.
func NewUnitApp(
	txnCtrl repo_iface.TxnCtrl,
	comicRepo repo_iface.ComicRepo,
	chapterRepo repo_iface.ChapterRepo,
	pageRepo repo_iface.PageRepo,
	unitRepo repo_iface.UnitRepo,
	assignmentRepo repo_iface.AssignmentRepo,
	unitSvc svc.UnitSvc,
	errClsf repo_iface.ErrClsf,
) app_iface.UnitApp {
	if txnCtrl == nil || comicRepo == nil || chapterRepo == nil || pageRepo == nil || unitRepo == nil || assignmentRepo == nil || errClsf == nil {
		zap.L().Panic(
			"[NewUnitApp] nil dependency",
			zap.Bool("txnCtrl", txnCtrl == nil),
			zap.Bool("comicRepo", comicRepo == nil),
			zap.Bool("chapterRepo", chapterRepo == nil),
			zap.Bool("pageRepo", pageRepo == nil),
			zap.Bool("unitRepo", unitRepo == nil),
			zap.Bool("assignmentRepo", assignmentRepo == nil),
			zap.Bool("errClsf", errClsf == nil),
		)
	}

	return &unitAppImpl{
		txnCtrl:        txnCtrl,
		unitSvc:        unitSvc,
		comicRepo:      comicRepo,
		chapterRepo:    chapterRepo,
		pageRepo:       pageRepo,
		unitRepo:       unitRepo,
		assignmentRepo: assignmentRepo,
		errClsf:        errClsf,
	}
}

// `ListByPage` returns units under one page.
func (a *unitAppImpl) ListByPage(cx context.Context, currUid string, args *val.ListPageUnitsArgs) app_res.AppRes[val.ListPageUnitsRes] {
	lgr := app_util.TakeLgr(cx)

	// Validate the target page id before querying repositories.
	if re := vfyListPageUnitsArgs(args); re.IsReject() {
		return app_res.Reject[val.ListPageUnitsRes](re.Code(), re.Msg())
	}

	// Load page first so chapter-scoped permission can be checked.
	page, err := a.pageRepo.GetById(args.PageId)
	if err != nil {
		if repo_infra.IsNotFound(err) {
			return app_res.Reject[val.ListPageUnitsRes](app_res.BadRequest, "页面不存在")
		}

		lgr.Error("[unitAppImpl.ListByPage] failed to get page", zap.Error(err))

		return app_res.Reject[val.ListPageUnitsRes](app_res.ServerError, "获取 unit 列表失败")
	}

	// Require chapter assignment to list page units.
	if re := a.unitSvc.CanListPageUnits(currUid, page.ChapterId, a.assignmentRepo, a.errClsf); re.IsReject() {
		return app_res.Reject[val.ListPageUnitsRes](app_res.ErrCode(re.Code()), re.Msg())
	}

	// Load units ordered by page index and assemble app-facing values.
	units, err := a.unitRepo.ListByPage(args.PageId)
	if err != nil {
		lgr.Error("[unitAppImpl.ListByPage] failed to list units", zap.Error(err))

		return app_res.Reject[val.ListPageUnitsRes](app_res.ServerError, "获取 unit 列表失败")
	}

	unitVals := make([]val.UnitVal, 0, len(units))
	for i := range units {
		unitVals = append(unitVals, asmUnitVal(units[i]))
	}

	return app_res.Accept(&val.ListPageUnitsRes{
		Units:               unitVals,
		TotalUnitCount:      page.TotalUnitCount,
		TranslatedUnitCount: page.TranslatedUnitCount,
		ProofreadUnitCount:  page.ProofreadUnitCount,
	})
}

// `SaveByPage` applies one page unit diff and synchronizes page and chapter counts.
func (a *unitAppImpl) SaveByPage(cx context.Context, currUid string, args *val.SavePageUnitsArgs) app_res.AppRes[val.SavePageUnitsRes] {
	lgr := app_util.TakeLgr(cx)

	// Validate the incoming diff and convert it to a domain-safe payload.
	diff, vfyRe := vfySavePageUnitsArgs(args)
	if vfyRe.IsReject() {
		return app_res.Reject[val.SavePageUnitsRes](vfyRe.Code(), vfyRe.Msg())
	}

	// Apply unit mutations and counter synchronization in one transaction.
	re, err := repo_iface.RunWithTxn[app_res.AppRes[val.SavePageUnitsRes]](a.txnCtrl, func(prov repo_iface.Prov) (app_res.AppRes[val.SavePageUnitsRes], error) {
		comicRepo := prov.ComicRepo()
		chapterRepo := prov.ChapterRepo()
		pageRepo := prov.PageRepo()
		unitRepo := prov.UnitRepo()
		assignmentRepo := prov.AssignmentRepo()

		page, err := pageRepo.GetById(args.PageId)
		if err != nil {
			if repo_infra.IsNotFound(err) {
				return app_res.Reject[val.SavePageUnitsRes](app_res.BadRequest, "页面不存在"), app_res.DefErr()
			}

			return app_res.Reject[val.SavePageUnitsRes](app_res.ServerError, "保存 unit 失败"), err
		}

		oldTotal := page.TotalUnitCount
		oldTranslated := page.TranslatedUnitCount
		oldProofread := page.ProofreadUnitCount

		// Require translator or proofreader role before mutating units.
		if re := a.unitSvc.CanEditPageUnits(currUid, page.ChapterId, assignmentRepo, a.errClsf); re.IsReject() {
			return app_res.Reject[val.SavePageUnitsRes](app_res.ErrCode(re.Code()), re.Msg()), app_res.DefErr()
		}

		// Load chapter after permission check so comic last-active can be touched.
		chapter, err := chapterRepo.GetById(page.ChapterId)
		if err != nil {
			return app_res.Reject[val.SavePageUnitsRes](app_res.ServerError, "保存 unit 失败"), err
		}

		if re := a.unitSvc.ApplyOps(diff, unitRepo); re.IsReject() {
			return app_res.Reject[val.SavePageUnitsRes](app_res.ErrCode(re.Code()), re.Msg()), app_res.DefErr()
		}

		newTotal, newTranslated, newProofread, err := unitRepo.CountByPage(args.PageId)
		if err != nil {
			return app_res.Reject[val.SavePageUnitsRes](app_res.ServerError, "保存 unit 失败"), err
		}

		if err := pageRepo.SetUnitCounts(args.PageId, newTotal, newTranslated, newProofread); err != nil {
			return app_res.Reject[val.SavePageUnitsRes](app_res.ServerError, "保存 unit 失败"), err
		}

		deltaTotal := newTotal - oldTotal
		deltaTranslated := newTranslated - oldTranslated
		deltaProofread := newProofread - oldProofread

		if err := chapterRepo.AdjustUnitCounts(page.ChapterId, deltaTotal, deltaTranslated, deltaProofread); err != nil {
			return app_res.Reject[val.SavePageUnitsRes](app_res.ServerError, "保存 unit 失败"), err
		}

		if err := comicRepo.TouchLastActive(chapter.ComicId); err != nil {
			return app_res.Reject[val.SavePageUnitsRes](app_res.ServerError, "保存 unit 失败"), err
		}

		return app_res.Accept(&val.SavePageUnitsRes{
			TotalUnitCount:      newTotal,
			TranslatedUnitCount: newTranslated,
			ProofreadUnitCount:  newProofread,
		}), nil
	})
	if err != nil {
		lgr.Error("[unitAppImpl.SaveByPage] failed to save page units", zap.Error(err))

		return re
	}

	return re
}
