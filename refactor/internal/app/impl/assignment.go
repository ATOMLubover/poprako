package app_impl

import (
	"context"

	app_iface "poprako-s/internal/app"
	"poprako-s/internal/app/res"
	app_util "poprako-s/internal/app/util"
	"poprako-s/internal/app/val"
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
	assignmentRepo repo_iface.AssignmentRepo

	evBus event_iface.EvBus
}

// `NewAssignmentApp` creates one `AssignmentApp` implementation.
func NewAssignmentApp(txnCtrl repo_iface.TxnCtrl, assignmentSvc svc.AssignmentSvc, memberRepo repo_iface.MemberRepo, worksetRepo repo_iface.WorksetRepo, comicRepo repo_iface.ComicRepo, chapterRepo repo_iface.ChapterRepo, assignmentRepo repo_iface.AssignmentRepo, evBus event_iface.EvBus) app_iface.AssignmentApp {
	if txnCtrl == nil || memberRepo == nil || worksetRepo == nil || comicRepo == nil || chapterRepo == nil || assignmentRepo == nil || evBus == nil {
		zap.L().Panic(
			"[NewAssignmentApp] nil dependency",
			zap.Bool("txnCtrl", txnCtrl == nil),
			zap.Bool("memberRepo", memberRepo == nil),
			zap.Bool("worksetRepo", worksetRepo == nil),
			zap.Bool("comicRepo", comicRepo == nil),
			zap.Bool("chapterRepo", chapterRepo == nil),
			zap.Bool("assignmentRepo", assignmentRepo == nil),
			zap.Bool("evBus", evBus == nil),
		)
	}

	return &assignmentAppImpl{
		txnCtrl:        txnCtrl,
		assignmentSvc:  assignmentSvc,
		memberRepo:     memberRepo,
		worksetRepo:    worksetRepo,
		comicRepo:      comicRepo,
		chapterRepo:    chapterRepo,
		assignmentRepo: assignmentRepo,
		evBus:          evBus,
	}
}

// `ListByChapter` lists assignments under one chapter.
func (a *assignmentAppImpl) ListByChapter(cx context.Context, currUid string, args *val.ListAssignmentByChapterArgs) res.AppRes[[]val.AssignmentVal] {
	lgr := app_util.TakeLgr(cx)

	if args == nil || args.ChapterId == "" {
		return res.Reject[[]val.AssignmentVal](res.BadRequest, "chapter_id 不能为空")
	}

	if code, msg, reject := vfyAssignmentListArgs(args.Offset, &args.Limit); reject {
		return res.Reject[[]val.AssignmentVal](code, msg)
	}

	if code, msg, reject := ensureReviewerPermission(a.assignmentRepo, args.ChapterId, currUid); reject {
		return res.Reject[[]val.AssignmentVal](code, msg)
	}

	items, err := a.assignmentRepo.List(mkListAssignmentOptByChapter(args.ChapterId, args.Offset, args.Limit))
	if err != nil {
		lgr.Error("[assignmentAppImpl.ListByChapter] failed to list assignments", zap.Error(err))
		return res.Reject[[]val.AssignmentVal](res.ServerError, "获取分配列表失败")
	}

	vals := make([]val.AssignmentVal, len(items))
	for i, item := range items {
		vals[i] = asmAssignmentVal(item)
	}

	return res.Accept(&vals)
}

// `ListByUser` lists all assignments of current user.
func (a *assignmentAppImpl) ListByUser(cx context.Context, currUid string, args *val.ListMyAssignmentArgs) res.AppRes[[]val.AssignmentVal] {
	lgr := app_util.TakeLgr(cx)

	if args == nil {
		return res.Reject[[]val.AssignmentVal](res.BadRequest, "分页参数不能为空")
	}

	if code, msg, reject := vfyAssignmentListArgs(args.Offset, &args.Limit); reject {
		return res.Reject[[]val.AssignmentVal](code, msg)
	}

	items, err := a.assignmentRepo.List(mkListAssignmentOptByUser(currUid, args.Offset, args.Limit))
	if err != nil {
		lgr.Error("[assignmentAppImpl.ListMy] failed to list assignments", zap.Error(err))
		return res.Reject[[]val.AssignmentVal](res.ServerError, "获取分配列表失败")
	}

	vals := make([]val.AssignmentVal, len(items))
	for i, item := range items {
		vals[i] = asmAssignmentVal(item)
	}

	return res.Accept(&vals)
}

// `Upsert` executes put-semantics upsert for assignment roles.
func (a *assignmentAppImpl) Upsert(cx context.Context, currUid string, args *val.UpsertAssignmentArgs) res.AppRes[res.None] {
	lgr := app_util.TakeLgr(cx)

	if args == nil || args.ChapterId == "" || args.UserId == "" {
		return res.Reject[res.None](res.BadRequest, "chapter_id 和 user_id 不能为空")
	}

	ev := make([]event_iface.Event, 0)

	re, err := repo_iface.RunWithTxn[res.AppRes[res.None]](a.txnCtrl, func(prov repo_iface.Prov) (res.AppRes[res.None], error) {
		memberRepo := prov.MemberRepo()
		worksetRepo := prov.WorksetRepo()
		comicRepo := prov.ComicRepo()
		chapterRepo := prov.ChapterRepo()
		assignmentRepo := prov.AssignmentRepo()

		if code, msg, reject := ensureReviewerPermission(assignmentRepo, args.ChapterId, currUid); reject {
			return res.Reject[res.None](code, msg), res.DefErr()
		}

		if code, msg, reject := ensureUserCanTakeRoles(memberRepo, chapterRepo, comicRepo, worksetRepo, args.ChapterId, args.UserId, args.RoleMask); reject {
			return res.Reject[res.None](code, msg), res.DefErr()
		}

		if args.RoleMask == 0 {
			_, err := assignmentRepo.GetByChapterUserId(args.ChapterId, args.UserId)
			if err != nil {
				if repo_infra.IsNotFound(err) {
					return res.Accept(&res.None{}), nil
				}
				return res.Reject[res.None](res.ServerError, "保存分配失败"), err
			}

			ch, err := chapterRepo.GetById(args.ChapterId)
			if err != nil {
				return res.Reject[res.None](res.ServerError, "删除分配失败"), err
			}

			if err := assignmentRepo.DeleteByChapterUserId(args.ChapterId, args.UserId); err != nil {
				return res.Reject[res.None](res.ServerError, "删除分配失败"), err
			}

			ev = append(ev, event.NewAssignmentRemovedEv(args.UserId, args.ChapterId, ch.PublishedAt != nil))

			return res.Accept(&res.None{}), nil
		}

		curr, err := assignmentRepo.GetByChapterUserId(args.ChapterId, args.UserId)
		if err != nil {
			if !repo_infra.IsNotFound(err) {
				return res.Reject[res.None](res.ServerError, "保存分配失败"), err
			}

			cre := a.assignmentSvc.NewAssignmentCre(args.ChapterId, args.UserId, args.RoleMask)
			if _, err := assignmentRepo.Create(cre); err != nil {
				return res.Reject[res.None](res.ServerError, "保存分配失败"), err
			}

			ev = append(ev, event.NewAssignmentCreatedEv(args.UserId, args.ChapterId))

			return res.Accept(&res.None{}), nil
		}

		put := a.assignmentSvc.NewAssignmentPut(curr, args.RoleMask)
		if err := assignmentRepo.Put(put); err != nil {
			return res.Reject[res.None](res.ServerError, "保存分配失败"), err
		}

		return res.Accept(&res.None{}), nil
	})
	if err != nil {
		lgr.Error("[assignmentAppImpl.Upsert] failed to run transaction", zap.Error(err))
		return re
	}

	a.evBus.Pub(context.Background(), ev)

	return re
}

// `Delete` deletes one assignment by id.
func (a *assignmentAppImpl) Delete(cx context.Context, currUid string, assignmentId string) res.AppRes[res.None] {
	lgr := app_util.TakeLgr(cx)

	if assignmentId == "" {
		return res.Reject[res.None](res.BadRequest, "assignment_id 不能为空")
	}

	ev := make([]event_iface.Event, 0)

	re, err := repo_iface.RunWithTxn[res.AppRes[res.None]](a.txnCtrl, func(prov repo_iface.Prov) (res.AppRes[res.None], error) {
		assignmentRepo := prov.AssignmentRepo()
		chapterRepo := prov.ChapterRepo()

		target, err := assignmentRepo.GetById(assignmentId)
		if err != nil {
			if repo_infra.IsNotFound(err) {
				return res.Reject[res.None](res.NotFound, "分配不存在"), res.DefErr()
			}

			return res.Reject[res.None](res.ServerError, "删除分配失败"), err
		}

		if code, msg, reject := ensureReviewerPermission(assignmentRepo, target.ChapterId, currUid); reject {
			return res.Reject[res.None](code, msg), res.DefErr()
		}

		ch, err := chapterRepo.GetById(target.ChapterId)
		if err != nil {
			return res.Reject[res.None](res.ServerError, "删除分配失败"), err
		}

		if err := assignmentRepo.Delete(assignmentId); err != nil {
			return res.Reject[res.None](res.ServerError, "删除分配失败"), err
		}

		ev = append(ev, event.NewAssignmentRemovedEv(target.UserId, target.ChapterId, ch.PublishedAt != nil))

		return res.Accept(&res.None{}), nil
	})
	if err != nil {
		lgr.Error("[assignmentAppImpl.Delete] failed to run transaction", zap.Error(err))
		return re
	}

	a.evBus.Pub(context.Background(), ev)

	return re
}
