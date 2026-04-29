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

// `comicAppImpl` is the default implementation of `ComicApp`
type comicAppImpl struct {
	txnCtrl repo_iface.TxnCtrl

	comicSvc      svc.ComicSvc
	chapterSvc    svc.ChapterSvc
	assignmentSvc svc.AssignmentSvc

	memberRepo  repo_iface.MemberRepo
	worksetRepo repo_iface.WorksetRepo
	comicRepo   repo_iface.ComicRepo

	evBus  event_iface.EvBus
	errClsf repo_iface.ErrClsf
}

// `NewComicApp` creates a `ComicApp` implementation
func NewComicApp(
	txnCtrl repo_iface.TxnCtrl,
	comicSvc svc.ComicSvc,
	chapterSvc svc.ChapterSvc,
	assignmentSvc svc.AssignmentSvc,
	memberRepo repo_iface.MemberRepo,
	worksetRepo repo_iface.WorksetRepo,
	comicRepo repo_iface.ComicRepo,
	evBus event_iface.EvBus,
	errClsf repo_iface.ErrClsf,
) app_iface.ComicApp {
	if txnCtrl == nil ||
		memberRepo == nil ||
		worksetRepo == nil ||
		comicRepo == nil ||
		evBus == nil ||
		errClsf == nil {
		zap.L().Panic(
			"[NewComicApp] nil dependency",
			zap.Bool("txnCtrl", txnCtrl == nil),
			zap.Bool("memberRepo", memberRepo == nil),
			zap.Bool("worksetRepo", worksetRepo == nil),
			zap.Bool("comicRepo", comicRepo == nil),
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
		evBus:         evBus,
		errClsf:       errClsf,
	}
}

// `List` returns all active comics for one workset
func (a *comicAppImpl) List(cx context.Context, currUid string, args *val.ListComicArgs) app_res.AppRes[[]val.ComicVal] {
	lgr := app_util.TakeLgr(cx)

	if re := vfyListComicArgs(args); re.IsReject() {
		return app_res.Reject[[]val.ComicVal](re.Code(), re.Msg())
	}

	ws, err := a.worksetRepo.GetById(args.WorksetId)
	if err != nil {
		return app_res.Reject[[]val.ComicVal](app_res.BadRequest, "作品集不存在")
	}

	if re := a.comicSvc.CanListComic(currUid, ws.TeamId, a.memberRepo, a.errClsf); re.IsReject() {
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

	vals := make([]val.ComicVal, len(comics))
	for i, cm := range comics {
		vals[i] = asmComicVal(cm)
	}

	return app_res.Accept(&vals)
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

		ws, err := worksetRepo.GetById(args.WorksetId)
		if err != nil {
			return app_res.Reject[val.ComicCreatedRes](app_res.Forbidden, "仅汉化组管理员可创建漫画"), app_res.DefErr()
		}

		if re := a.comicSvc.CanAdminComic(currUid, ws.TeamId, memberRepo, a.errClsf); re.IsReject() {
			return app_res.Reject[val.ComicCreatedRes](app_res.ErrCode(re.Code()), re.Msg()), app_res.DefErr()
		}

		count, err := comicRepo.Count(&query.ListComicOpt{WorksetId: &args.WorksetId})
		if err != nil {
			return app_res.Reject[val.ComicCreatedRes](app_res.ServerError, "创建漫画失败"), err
		}

		cre := a.comicSvc.NewComicCre(
			args.WorksetId,
			int(count),
			args.Title,
			args.Author,
			args.Desc,
			currUid,
		)

		cm, err := comicRepo.Create(cre)
		if err != nil {
			return app_res.Reject[val.ComicCreatedRes](app_res.ServerError, "创建漫画失败"), err
		}

		if err := worksetRepo.UpdateComicCount(args.WorksetId, 1); err != nil {
			return app_res.Reject[val.ComicCreatedRes](app_res.ServerError, "创建漫画失败"), err
		}

		// Create first chapter with index 0 as default pinned chapter.
		chCre := a.chapterSvc.NewChapterCre(cm.Id, 0, nil, currUid)
		ch, err := prov.ChapterRepo().Create(chCre)
		if err != nil {
			return app_res.Reject[val.ComicCreatedRes](app_res.ServerError, "创建漫画失败"), err
		}

		if err := comicRepo.UpdateChapterCount(cm.Id, 1); err != nil {
			return app_res.Reject[val.ComicCreatedRes](app_res.ServerError, "创建漫画失败"), err
		}

		if err := comicRepo.TouchLastActive(cm.Id); err != nil {
			return app_res.Reject[val.ComicCreatedRes](app_res.ServerError, "创建漫画失败"), err
		}

		reviewerCre := a.assignmentSvc.NewAssignmentCre(ch.Id, currUid, aggr.RoleMask(enum.RoleReviewer))
		if _, err := prov.AssignmentRepo().Create(reviewerCre); err != nil {
			return app_res.Reject[val.ComicCreatedRes](app_res.ServerError, "创建漫画失败"), err
		}

		ev = append(ev, event.NewAssignmentCreatedEv(currUid, ch.Id))

		return app_res.Accept(&val.ComicCreatedRes{Id: cm.Id}), nil
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

	cm, err := a.comicRepo.GetById(args.Id)
	if err != nil {
		lgr.Error("[comicAppImpl.Update] failed to get comic", zap.Error(err))

		return app_res.Reject[app_res.None](app_res.BadRequest, "漫画不存在")
	}

	ws, err := a.worksetRepo.GetById(cm.WorksetId)
	if err != nil {
		lgr.Error("[comicAppImpl.Update] failed to get workset", zap.Error(err))

		return app_res.Reject[app_res.None](app_res.ServerError, "更新漫画失败")
	}

	if re := a.comicSvc.CanAdminComic(currUid, ws.TeamId, a.memberRepo, a.errClsf); re.IsReject() {
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

// `Remove` soft-deletes a comic and updates workset comic counter
func (a *comicAppImpl) Remove(cx context.Context, currUid string, comicId string) app_res.AppRes[app_res.None] {
	lgr := app_util.TakeLgr(cx)

	if re := vfyRemoveComicId(comicId); re.IsReject() {
		return app_res.Reject[app_res.None](re.Code(), re.Msg())
	}

	re, err := repo_iface.RunWithTxn[app_res.AppRes[app_res.None]](a.txnCtrl, func(prov repo_iface.Prov) (app_res.AppRes[app_res.None], error) {
		memberRepo := prov.MemberRepo()
		worksetRepo := prov.WorksetRepo()
		comicRepo := prov.ComicRepo()

		cm, err := comicRepo.GetById(comicId)
		if err != nil {
			return app_res.Reject[app_res.None](app_res.Forbidden, "仅汉化组管理员可删除漫画"), app_res.DefErr()
		}

		ws, err := worksetRepo.GetById(cm.WorksetId)
		if err != nil {
			return app_res.Reject[app_res.None](app_res.Forbidden, "仅汉化组管理员可删除漫画"), app_res.DefErr()
		}

		if re := a.comicSvc.CanAdminComic(currUid, ws.TeamId, memberRepo, a.errClsf); re.IsReject() {
			return app_res.Reject[app_res.None](app_res.ErrCode(re.Code()), re.Msg()), app_res.DefErr()
		}

		if err := comicRepo.Remove(comicId); err != nil {
			return app_res.Reject[app_res.None](app_res.ServerError, "删除漫画失败"), err
		}

		if err := worksetRepo.UpdateComicCount(ws.Id, -1); err != nil {
			return app_res.Reject[app_res.None](app_res.ServerError, "删除漫画失败"), err
		}

		return app_res.Accept(&app_res.None{}), nil
	})
	if err != nil {
		lgr.Error(
			"[comicAppImpl.Remove] failed to run remove comic transaction",
			zap.Error(err),
		)

		return re
	}

	return re
}
