package app_impl

import (
	"context"

	app_iface "poprako-s/internal/app"
	app_res "poprako-s/internal/app/res"
	app_util "poprako-s/internal/app/util"
	"poprako-s/internal/app/val"
	oss_iface "poprako-s/internal/domain/ext/oss"
	"poprako-s/internal/domain/model/event"
	repo_iface "poprako-s/internal/domain/repo"
	"poprako-s/internal/domain/svc"
	event_iface "poprako-s/internal/event"
	repo_infra "poprako-s/internal/infra/repo"

	"go.uber.org/zap"
)

// `assignmentAppImpl` is the default implementation of `AssignmentApp`.
type assignmentAppImpl struct {
	txnCtrl repo_iface.TxnCtrl

	assignmentSvc svc.AssignmentSvc

	memberRepo     repo_iface.MemberRepo
	worksetRepo    repo_iface.WorksetRepo
	comicRepo      repo_iface.ComicRepo
	chapterRepo    repo_iface.ChapterRepo
	pageRepo       repo_iface.PageRepo
	assignmentRepo repo_iface.AssignmentRepo
	ossSigner      oss_iface.Signer

	evBus event_iface.EvBus

	errClsf repo_iface.ErrClsf
}

// `NewAssignmentApp` creates one `AssignmentApp` implementation.
func NewAssignmentApp(
	txnCtrl repo_iface.TxnCtrl,
	memberRepo repo_iface.MemberRepo,
	worksetRepo repo_iface.WorksetRepo,
	comicRepo repo_iface.ComicRepo,
	chapterRepo repo_iface.ChapterRepo,
	pageRepo repo_iface.PageRepo,
	assignmentRepo repo_iface.AssignmentRepo,
	assignmentSvc svc.AssignmentSvc,
	ossSigner oss_iface.Signer,
	evBus event_iface.EvBus,
	errClsf repo_iface.ErrClsf,
) app_iface.AssignmentApp {
	if txnCtrl == nil || memberRepo == nil || worksetRepo == nil || comicRepo == nil || chapterRepo == nil || pageRepo == nil || assignmentRepo == nil || ossSigner == nil || evBus == nil || errClsf == nil {
		zap.L().Panic(
			"[NewAssignmentApp] nil dependency",
			zap.Bool("txnCtrl", txnCtrl == nil),
			zap.Bool("memberRepo", memberRepo == nil),
			zap.Bool("worksetRepo", worksetRepo == nil),
			zap.Bool("comicRepo", comicRepo == nil),
			zap.Bool("chapterRepo", chapterRepo == nil),
			zap.Bool("pageRepo", pageRepo == nil),
			zap.Bool("assignmentRepo", assignmentRepo == nil),
			zap.Bool("ossSigner", ossSigner == nil),
			zap.Bool("evBus", evBus == nil),
			zap.Bool("errClsf", errClsf == nil),
		)
	}

	return &assignmentAppImpl{
		txnCtrl:        txnCtrl,
		assignmentSvc:  assignmentSvc,
		memberRepo:     memberRepo,
		worksetRepo:    worksetRepo,
		comicRepo:      comicRepo,
		chapterRepo:    chapterRepo,
		pageRepo:       pageRepo,
		assignmentRepo: assignmentRepo,
		ossSigner:      ossSigner,
		evBus:          evBus,
		errClsf:        errClsf,
	}
}

// `ListByChapter` lists assignments under one chapter.
func (a *assignmentAppImpl) ListByChapter(cx context.Context, currUid string, args *val.ListAssignmentByChapterArgs) app_res.AppRes[[]val.AssignmentVal] {
	lgr := app_util.TakeLgr(cx)

	if args == nil || args.ChapterId == "" {
		return app_res.Reject[[]val.AssignmentVal](app_res.BadRequest, "chapter_id 不能为空")
	}

	if re := vfyAssignmentListArgs(args.Offset, &args.Limit); re.IsReject() {
		return app_res.Reject[[]val.AssignmentVal](re.Code(), re.Msg())
	}

	if re := a.assignmentSvc.CanListByChapter(
		currUid,
		args.ChapterId,
		a.memberRepo,
		a.worksetRepo,
		a.comicRepo,
		a.chapterRepo,
		a.assignmentRepo,
		a.errClsf,
	); re.IsReject() {
		return app_res.Reject[[]val.AssignmentVal](app_res.ErrCode(re.Code()), re.Msg())
	}

	items, err := a.assignmentRepo.List(
		mkListAssignmentOptByChapter(args.ChapterId, args.Includes, args.Offset, args.Limit),
		mkAssignmentRepoIncl(args.Includes)...,
	)
	if err != nil {
		lgr.Error("[assignmentAppImpl.ListByChapter] failed to list assignments", zap.Error(err))
		return app_res.Reject[[]val.AssignmentVal](app_res.ServerError, "获取分配列表失败")
	}

	assignmentVals := make([]val.AssignmentVal, len(items))
	for i, item := range items {
		assignmentVals[i] = asmAssignmentVal(item)
	}

	if err := tryFillCoverForAssignments(
		assignmentVals,
		a.ossSigner,
		memoizeComicCoverKeyGetter(mkPinnedFirstPageImageKeyGetter(a.chapterRepo, a.pageRepo, lgr)),
		lgr,
	); err != nil {
		errCode := mapComicFallbackErrCode(err, a.errClsf)

		lgr.Error(
			"[assignmentAppImpl.ListByChapter] failed to fill assignment comic covers",
			zap.Error(err),
			zap.Int("errCode", int(errCode)),
		)

		return app_res.Reject[[]val.AssignmentVal](errCode, "获取分配列表失败")
	}

	return app_res.Accept(&assignmentVals)
}

func (a *assignmentAppImpl) ListByUser(cx context.Context, currUid string, args *val.ListAssignmentByUserArgs) app_res.AppRes[[]val.AssignmentVal] {
	lgr := app_util.TakeLgr(cx)

	if args == nil {
		return app_res.Reject[[]val.AssignmentVal](app_res.BadRequest, "分页参数不能为空")
	}

	if re := vfyAssignmentListArgs(args.Offset, &args.Limit); re.IsReject() {
		return app_res.Reject[[]val.AssignmentVal](re.Code(), re.Msg())
	}

	items, err := a.assignmentRepo.List(
		mkListAssignmentOptByUser(currUid, args.Includes, args.Offset, args.Limit),
		mkAssignmentRepoIncl(args.Includes)...,
	)
	if err != nil {
		lgr.Error("[assignmentAppImpl.ListMy] failed to list assignments", zap.Error(err))
		return app_res.Reject[[]val.AssignmentVal](app_res.ServerError, "获取分配列表失败")
	}

	assignmentVals := make([]val.AssignmentVal, len(items))
	for i, item := range items {
		assignmentVals[i] = asmAssignmentVal(item)
	}

	if err := tryFillCoverForAssignments(
		assignmentVals,
		a.ossSigner,
		memoizeComicCoverKeyGetter(mkPinnedFirstPageImageKeyGetter(a.chapterRepo, a.pageRepo, lgr)),
		lgr,
	); err != nil {
		errCode := mapComicFallbackErrCode(err, a.errClsf)

		lgr.Error(
			"[assignmentAppImpl.ListByUser] failed to fill assignment comic covers",
			zap.Error(err),
			zap.Int("errCode", int(errCode)),
		)

		return app_res.Reject[[]val.AssignmentVal](errCode, "获取分配列表失败")
	}

	return app_res.Accept(&assignmentVals)
}

// `Upsert` executes put-semantics upsert for assignment roles.
func (a *assignmentAppImpl) Upsert(cx context.Context, currUid string, args *val.UpsertAssignmentArgs) app_res.AppRes[app_res.None] {
	lgr := app_util.TakeLgr(cx)

	if args == nil || args.ChapterId == "" || args.UserId == "" {
		return app_res.Reject[app_res.None](app_res.BadRequest, "chapter_id 和 user_id 不能为空")
	}

	ev := make([]event_iface.Event, 0)

	re, err := repo_iface.RunWithTxn[app_res.AppRes[app_res.None]](a.txnCtrl, func(prov repo_iface.Prov) (app_res.AppRes[app_res.None], error) {
		memberRepo := prov.MemberRepo()
		worksetRepo := prov.WorksetRepo()
		comicRepo := prov.ComicRepo()
		chapterRepo := prov.ChapterRepo()
		assignmentRepo := prov.AssignmentRepo()

		if re := a.assignmentSvc.CanReviewAssignment(currUid, args.ChapterId, assignmentRepo, a.errClsf); re.IsReject() {
			return app_res.Reject[app_res.None](app_res.ErrCode(re.Code()), re.Msg()), app_res.DefErr()
		}

		if re := a.assignmentSvc.CanTakeAssignmentRoles(args.UserId, args.ChapterId, args.RoleMask, memberRepo, chapterRepo, comicRepo, worksetRepo, a.errClsf); re.IsReject() {
			return app_res.Reject[app_res.None](app_res.ErrCode(re.Code()), re.Msg()), app_res.DefErr()
		}

		if args.RoleMask == 0 {
			_, err := assignmentRepo.GetByChapterUserId(args.ChapterId, args.UserId)
			if err != nil {
				if repo_infra.IsNotFound(err) {
					return app_res.Accept(&app_res.None{}), nil
				}
				return app_res.Reject[app_res.None](app_res.ServerError, "保存分配失败"), err
			}

			chapter, err := chapterRepo.GetById(args.ChapterId)
			if err != nil {
				return app_res.Reject[app_res.None](app_res.ServerError, "删除分配失败"), err
			}

			if err := assignmentRepo.DeleteByChapterUserId(args.ChapterId, args.UserId); err != nil {
				return app_res.Reject[app_res.None](app_res.ServerError, "删除分配失败"), err
			}

			ev = append(ev, event.NewAssignmentRemovedEv(args.UserId, args.ChapterId, chapter.PublishedAt != nil))

			return app_res.Accept(&app_res.None{}), nil
		}

		curr, err := assignmentRepo.GetByChapterUserId(args.ChapterId, args.UserId)
		if err != nil {
			if !repo_infra.IsNotFound(err) {
				return app_res.Reject[app_res.None](app_res.ServerError, "保存分配失败"), err
			}

			cre := a.assignmentSvc.NewAssignmentCre(args.ChapterId, args.UserId, args.RoleMask)
			if _, err := assignmentRepo.Create(cre); err != nil {
				return app_res.Reject[app_res.None](app_res.ServerError, "保存分配失败"), err
			}

			ev = append(ev, event.NewAssignmentCreatedEv(args.UserId, args.ChapterId))

			return app_res.Accept(&app_res.None{}), nil
		}

		put := a.assignmentSvc.NewAssignmentPut(curr, args.RoleMask)
		if err := assignmentRepo.Put(put); err != nil {
			return app_res.Reject[app_res.None](app_res.ServerError, "保存分配失败"), err
		}

		return app_res.Accept(&app_res.None{}), nil
	})
	if err != nil {
		lgr.Error("[assignmentAppImpl.Upsert] failed to run transaction", zap.Error(err))
		return re
	}

	a.evBus.Pub(context.Background(), ev)

	return re
}

// `Delete` deletes one assignment by id.
func (a *assignmentAppImpl) Delete(cx context.Context, currUid string, assignmentId string) app_res.AppRes[app_res.None] {
	lgr := app_util.TakeLgr(cx)

	if assignmentId == "" {
		return app_res.Reject[app_res.None](app_res.BadRequest, "assignment_id 不能为空")
	}

	ev := make([]event_iface.Event, 0)

	re, err := repo_iface.RunWithTxn[app_res.AppRes[app_res.None]](a.txnCtrl, func(prov repo_iface.Prov) (app_res.AppRes[app_res.None], error) {
		assignmentRepo := prov.AssignmentRepo()
		chapterRepo := prov.ChapterRepo()

		target, err := assignmentRepo.GetById(assignmentId)
		if err != nil {
			if repo_infra.IsNotFound(err) {
				return app_res.Reject[app_res.None](app_res.NotFound, "分配不存在"), app_res.DefErr()
			}

			return app_res.Reject[app_res.None](app_res.ServerError, "删除分配失败"), err
		}

		if re := a.assignmentSvc.CanReviewAssignment(currUid, target.ChapterId, assignmentRepo, a.errClsf); re.IsReject() {
			return app_res.Reject[app_res.None](app_res.ErrCode(re.Code()), re.Msg()), app_res.DefErr()
		}

		chapter, err := chapterRepo.GetById(target.ChapterId)
		if err != nil {
			return app_res.Reject[app_res.None](app_res.ServerError, "删除分配失败"), err
		}

		if err := assignmentRepo.Delete(assignmentId); err != nil {
			return app_res.Reject[app_res.None](app_res.ServerError, "删除分配失败"), err
		}

		ev = append(ev, event.NewAssignmentRemovedEv(target.UserId, target.ChapterId, chapter.PublishedAt != nil))

		return app_res.Accept(&app_res.None{}), nil
	})
	if err != nil {
		lgr.Error("[assignmentAppImpl.Delete] failed to run transaction", zap.Error(err))
		return re
	}

	a.evBus.Pub(context.Background(), ev)

	return re
}
