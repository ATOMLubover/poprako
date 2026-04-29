package app_impl

import (
	"context"

	app_iface "poprako-s/internal/app"
	app_res "poprako-s/internal/app/res"
	app_util "poprako-s/internal/app/util"
	"poprako-s/internal/app/val"
	"poprako-s/internal/domain/model/aggr"
	"poprako-s/internal/domain/model/enum"
	"poprako-s/internal/domain/model/event"
	"poprako-s/internal/domain/model/query"
	repo_iface "poprako-s/internal/domain/repo"
	"poprako-s/internal/domain/svc"
	event_iface "poprako-s/internal/event"

	"go.uber.org/zap"
)

// `chapterAppImpl` is the default implementation of `ChapterApp`.
type chapterAppImpl struct {
	txnCtrl repo_iface.TxnCtrl

	chapterSvc    svc.ChapterSvc
	assignmentSvc svc.AssignmentSvc

	memberRepo     repo_iface.MemberRepo
	worksetRepo    repo_iface.WorksetRepo
	comicRepo      repo_iface.ComicRepo
	chapterRepo    repo_iface.ChapterRepo
	assignmentRepo repo_iface.AssignmentRepo

	evBus event_iface.EvBus

	errClsf repo_iface.ErrClsf
}

// `NewChapterApp` creates one `ChapterApp` implementation.
func NewChapterApp(
	txnCtrl repo_iface.TxnCtrl,
	chapterSvc svc.ChapterSvc,
	assignmentSvc svc.AssignmentSvc,
	memberRepo repo_iface.MemberRepo,
	worksetRepo repo_iface.WorksetRepo,
	comicRepo repo_iface.ComicRepo,
	chapterRepo repo_iface.ChapterRepo,
	assignmentRepo repo_iface.AssignmentRepo,
	evBus event_iface.EvBus,
	errClsf repo_iface.ErrClsf,
) app_iface.ChapterApp {
	if txnCtrl == nil ||
		memberRepo == nil ||
		worksetRepo == nil ||
		comicRepo == nil ||
		chapterRepo == nil ||
		assignmentRepo == nil ||
		evBus == nil ||
		errClsf == nil {
		zap.L().Panic(
			"[NewChapterApp] nil dependency",
			zap.Bool("txnCtrl", txnCtrl == nil),
			zap.Bool("memberRepo", memberRepo == nil),
			zap.Bool("worksetRepo", worksetRepo == nil),
			zap.Bool("comicRepo", comicRepo == nil),
			zap.Bool("chapterRepo", chapterRepo == nil),
			zap.Bool("assignmentRepo", assignmentRepo == nil),
			zap.Bool("evBus", evBus == nil),
			zap.Bool("errClsf", errClsf == nil),
		)
	}

	return &chapterAppImpl{
		txnCtrl:        txnCtrl,
		chapterSvc:     chapterSvc,
		assignmentSvc:  assignmentSvc,
		memberRepo:     memberRepo,
		worksetRepo:    worksetRepo,
		comicRepo:      comicRepo,
		chapterRepo:    chapterRepo,
		assignmentRepo: assignmentRepo,
		evBus:          evBus,
		errClsf:         errClsf,
	}
}

// `List` returns chapter list under target comic.
func (a *chapterAppImpl) List(cx context.Context, currUid string, args *val.ListChapterArgs) app_res.AppRes[[]val.ChapterVal] {
	lgr := app_util.TakeLgr(cx)

	if re := vfyListChapterArgs(args); re.IsReject() {
		return app_res.Reject[[]val.ChapterVal](re.Code(), re.Msg())
	}

	cm, err := a.comicRepo.GetById(args.ComicId, enum.ComicInclWorkset)
	if err != nil {
		return app_res.Reject[[]val.ChapterVal](app_res.BadRequest, "漫画不存在")
	}

	if re := a.chapterSvc.CanListChapter(currUid, cm.Workset.TeamId, a.memberRepo, a.errClsf); re.IsReject() {
		return app_res.Reject[[]val.ChapterVal](app_res.ErrCode(re.Code()), re.Msg())
	}

	chapters, err := a.chapterRepo.List(&query.ListChapterOpt{
		ComicId: &args.ComicId,
		Pagi: query.PagiOpt{
			Offset: args.Offset,
			Limit:  args.Limit,
		},
	})
	if err != nil {
		lgr.Error("[chapterAppImpl.List] failed to list chapters", zap.Error(err))
		return app_res.Reject[[]val.ChapterVal](app_res.ServerError, "获取章节列表失败")
	}

	vals := make([]val.ChapterVal, len(chapters))
	for i, ch := range chapters {
		vals[i] = asmChapterVal(ch)
	}

	return app_res.Accept(&vals)
}

// `GetPinned` returns pinned chapter under target comic.
func (a *chapterAppImpl) GetPinned(cx context.Context, currUid string, comicId string) app_res.AppRes[val.ChapterVal] {
	lgr := app_util.TakeLgr(cx)

	if comicId == "" {
		return app_res.Reject[val.ChapterVal](app_res.BadRequest, "comic_id 不能为空")
	}

	cm, err := a.comicRepo.GetById(comicId, enum.ComicInclWorkset)
	if err != nil {
		return app_res.Reject[val.ChapterVal](app_res.BadRequest, "漫画不存在")
	}

	if re := a.chapterSvc.CanListChapter(currUid, cm.Workset.TeamId, a.memberRepo, a.errClsf); re.IsReject() {
		return app_res.Reject[val.ChapterVal](app_res.ErrCode(re.Code()), re.Msg())
	}

	ch, err := a.chapterRepo.FindPinnedByComicId(comicId)
	if err != nil {
		lgr.Error("[chapterAppImpl.GetPinned] failed to get pinned chapter", zap.Error(err))
		return app_res.Reject[val.ChapterVal](app_res.ServerError, "获取置顶章节失败")
	}

	if ch == nil {
		return app_res.Accept[val.ChapterVal](nil)
	}

	v := asmChapterVal(ch)

	return app_res.Accept(&v)
}

// `Create` creates one chapter under target comic.
func (a *chapterAppImpl) Create(cx context.Context, currUid string, args *val.CreateChapterArgs) app_res.AppRes[val.ChapterCreatedRes] {
	lgr := app_util.TakeLgr(cx)

	if re := vfyCreateChapterArgs(args); re.IsReject() {
		return app_res.Reject[val.ChapterCreatedRes](re.Code(), re.Msg())
	}

	ev := make([]event_iface.Event, 0)

	re, err := repo_iface.RunWithTxn[app_res.AppRes[val.ChapterCreatedRes]](a.txnCtrl, func(prov repo_iface.Prov) (app_res.AppRes[val.ChapterCreatedRes], error) {
		memberRepo := prov.MemberRepo()
		worksetRepo := prov.WorksetRepo()
		comicRepo := prov.ComicRepo()
		chapterRepo := prov.ChapterRepo()
		assignmentRepo := prov.AssignmentRepo()

		cm, err := comicRepo.GetById(args.ComicId)
		if err != nil {
			return app_res.Reject[val.ChapterCreatedRes](app_res.Forbidden, "仅汉化组管理员可创建章节"), app_res.DefErr()
		}

		ws, err := worksetRepo.GetById(cm.WorksetId)
		if err != nil {
			return app_res.Reject[val.ChapterCreatedRes](app_res.Forbidden, "仅汉化组管理员可创建章节"), app_res.DefErr()
		}

		if rj := a.chapterSvc.CanAdminChapter(currUid, ws.TeamId, memberRepo, a.errClsf); rj.IsReject() {
			return app_res.Reject[val.ChapterCreatedRes](app_res.ErrCode(rj.Code()), rj.Msg()), app_res.DefErr()
		}

		count, err := chapterRepo.Count(&query.ListChapterOpt{ComicId: &args.ComicId})
		if err != nil {
			return app_res.Reject[val.ChapterCreatedRes](app_res.ServerError, "创建章节失败"), err
		}

		cre := a.chapterSvc.NewChapterCre(args.ComicId, int(count), args.Subtitle, currUid)
		ch, err := chapterRepo.Create(cre)
		if err != nil {
			return app_res.Reject[val.ChapterCreatedRes](app_res.ServerError, "创建章节失败"), err
		}

		if err := comicRepo.UpdateChapterCount(args.ComicId, 1); err != nil {
			return app_res.Reject[val.ChapterCreatedRes](app_res.ServerError, "创建章节失败"), err
		}

		if err := comicRepo.TouchLastActive(args.ComicId); err != nil {
			return app_res.Reject[val.ChapterCreatedRes](app_res.ServerError, "创建章节失败"), err
		}

		reviewerCre := a.assignmentSvc.NewAssignmentCre(ch.Id, currUid, aggr.RoleMask(enum.RoleReviewer))
		if _, err := assignmentRepo.Create(reviewerCre); err != nil {
			return app_res.Reject[val.ChapterCreatedRes](app_res.ServerError, "创建章节失败"), err
		}

		ev = append(ev, event.NewAssignmentCreatedEv(currUid, ch.Id))

		return app_res.Accept(&val.ChapterCreatedRes{Id: ch.Id}), nil
	})
	if err != nil {
		lgr.Error("[chapterAppImpl.Create] failed to run create chapter transaction", zap.Error(err))

		return re
	}

	a.evBus.Pub(context.Background(), ev)

	return re
}

// `Update` updates one chapter.
func (a *chapterAppImpl) Update(cx context.Context, currUid string, args *val.ChapterUpdArgs) app_res.AppRes[app_res.None] {
	lgr := app_util.TakeLgr(cx)

	if re := vfyUpdateChapterArgs(args); re.IsReject() {
		return app_res.Reject[app_res.None](re.Code(), re.Msg())
	}

	ev := make([]event_iface.Event, 0)

	re, err := repo_iface.RunWithTxn[app_res.AppRes[app_res.None]](a.txnCtrl, func(prov repo_iface.Prov) (app_res.AppRes[app_res.None], error) {
		memberRepo := prov.MemberRepo()
		worksetRepo := prov.WorksetRepo()
		comicRepo := prov.ComicRepo()
		chapterRepo := prov.ChapterRepo()
		assignmentRepo := prov.AssignmentRepo()

		ch, err := chapterRepo.GetById(args.Id)
		if err != nil {
			return app_res.Reject[app_res.None](app_res.Forbidden, "仅汉化组管理员可更新章节"), app_res.DefErr()
		}

		cm, err := comicRepo.GetById(ch.ComicId)
		if err != nil {
			return app_res.Reject[app_res.None](app_res.Forbidden, "仅汉化组管理员可更新章节"), app_res.DefErr()
		}

		ws, err := worksetRepo.GetById(cm.WorksetId)
		if err != nil {
			return app_res.Reject[app_res.None](app_res.Forbidden, "仅汉化组管理员可更新章节"), app_res.DefErr()
		}

		if rj := a.chapterSvc.CanAdminChapter(currUid, ws.TeamId, memberRepo, a.errClsf); rj.IsReject() {
			return app_res.Reject[app_res.None](app_res.ErrCode(rj.Code()), rj.Msg()), app_res.DefErr()
		}

		if args.WorkflowTransition != nil {
			wasPublished := ch.PublishedAt != nil

			if err := ch.TransiteWorkflow(*args.WorkflowTransition); err != nil {
				return app_res.Reject[app_res.None](app_res.BadRequest, "无效的工作流状态转换"), app_res.DefErr()
			}

			if !wasPublished && ch.PublishedAt != nil {
				assignments, lerr := assignmentRepo.List(&query.ListAssignmentOpt{ChapterId: &ch.Id, Pagi: query.PagiOpt{Limit: 500}})
				if lerr != nil {
					return app_res.Reject[app_res.None](app_res.ServerError, "更新章节失败"), lerr
				}

				assignedUserIds := make([]string, 0, len(assignments))
				for i := range assignments {
					assignedUserIds = append(assignedUserIds, assignments[i].UserId)
				}

				ev = append(ev, event.NewChapterPublishedEv(ch.Id, assignedUserIds))
			}
		}

		upd := mkChapterUpd(args, ch)
		if err := chapterRepo.Update(upd); err != nil {
			return app_res.Reject[app_res.None](app_res.ServerError, "更新章节失败"), err
		}

		if err := comicRepo.TouchLastActive(ch.ComicId); err != nil {
			return app_res.Reject[app_res.None](app_res.ServerError, "更新章节失败"), err
		}

		return app_res.Accept(&app_res.None{}), nil
	})
	if err != nil {
		lgr.Error("[chapterAppImpl.Update] failed to run update chapter transaction", zap.Error(err))

		return re
	}

	a.evBus.Pub(context.Background(), ev)

	return re
}

// `Remove` soft-deletes one chapter.
func (a *chapterAppImpl) Remove(cx context.Context, currUid string, chapterId string) app_res.AppRes[app_res.None] {
	lgr := app_util.TakeLgr(cx)

	if re := vfyRemoveChapterId(chapterId); re.IsReject() {
		return app_res.Reject[app_res.None](re.Code(), re.Msg())
	}

	ev := make([]event_iface.Event, 0)

	re, err := repo_iface.RunWithTxn[app_res.AppRes[app_res.None]](a.txnCtrl, func(prov repo_iface.Prov) (app_res.AppRes[app_res.None], error) {
		memberRepo := prov.MemberRepo()
		worksetRepo := prov.WorksetRepo()
		comicRepo := prov.ComicRepo()
		chapterRepo := prov.ChapterRepo()
		assignmentRepo := prov.AssignmentRepo()

		ch, err := chapterRepo.GetById(chapterId)
		if err != nil {
			return app_res.Reject[app_res.None](app_res.Forbidden, "仅汉化组管理员可删除章节"), app_res.DefErr()
		}

		cm, err := comicRepo.GetById(ch.ComicId)
		if err != nil {
			return app_res.Reject[app_res.None](app_res.Forbidden, "仅汉化组管理员可删除章节"), app_res.DefErr()
		}

		ws, err := worksetRepo.GetById(cm.WorksetId)
		if err != nil {
			return app_res.Reject[app_res.None](app_res.Forbidden, "仅汉化组管理员可删除章节"), app_res.DefErr()
		}

		if rj := a.chapterSvc.CanAdminChapter(currUid, ws.TeamId, memberRepo, a.errClsf); rj.IsReject() {
			return app_res.Reject[app_res.None](app_res.ErrCode(rj.Code()), rj.Msg()), app_res.DefErr()
		}

		assignments, err := assignmentRepo.List(&query.ListAssignmentOpt{ChapterId: &ch.Id, Pagi: query.PagiOpt{Limit: 500}})
		if err != nil {
			return app_res.Reject[app_res.None](app_res.ServerError, "删除章节失败"), err
		}

		assignedUserIds := make([]string, 0, len(assignments))
		for i := range assignments {
			assignedUserIds = append(assignedUserIds, assignments[i].UserId)
		}

		if err := assignmentRepo.DeleteByChapterId(ch.Id); err != nil {
			return app_res.Reject[app_res.None](app_res.ServerError, "删除章节失败"), err
		}

		if err := chapterRepo.Remove(chapterId); err != nil {
			return app_res.Reject[app_res.None](app_res.ServerError, "删除章节失败"), err
		}

		if err := comicRepo.UpdateChapterCount(ch.ComicId, -1); err != nil {
			return app_res.Reject[app_res.None](app_res.ServerError, "删除章节失败"), err
		}

		if err := comicRepo.TouchLastActive(ch.ComicId); err != nil {
			return app_res.Reject[app_res.None](app_res.ServerError, "删除章节失败"), err
		}

		ev = append(ev, event.NewChapterRemovedEv(ch.Id, ch.PublishedAt != nil, assignedUserIds))

		return app_res.Accept(&app_res.None{}), nil
	})
	if err != nil {
		lgr.Error("[chapterAppImpl.Remove] failed to run remove chapter transaction", zap.Error(err))

		return re
	}

	a.evBus.Pub(context.Background(), ev)

	return re
}
