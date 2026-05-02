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

// `worksetAppImpl` is the default implementation of `WorksetApp`.
type worksetAppImpl struct {
	txnCtrl     repo_iface.TxnCtrl
	worksetSvc  svc.WorksetSvc
	memberRepo  repo_iface.MemberRepo
	worksetRepo repo_iface.WorksetRepo
	evBus       event_iface.EvBus
	errClsf     repo_iface.ErrClsf
}

// `NewWorksetApp` creates a `WorksetApp` implementation.
func NewWorksetApp(
	txnCtrl repo_iface.TxnCtrl,
	memberRepo repo_iface.MemberRepo,
	worksetRepo repo_iface.WorksetRepo,
	worksetSvc svc.WorksetSvc,
	evBus event_iface.EvBus,
	errClsf repo_iface.ErrClsf,
) app_iface.WorksetApp {
	if txnCtrl == nil ||
		memberRepo == nil ||
		worksetRepo == nil ||
		evBus == nil ||
		errClsf == nil {
		zap.L().Panic(
			"[NewWorksetApp] nil dependency",
			zap.Bool("txnCtrl", txnCtrl == nil),
			zap.Bool("memberRepo", memberRepo == nil),
			zap.Bool("worksetRepo", worksetRepo == nil),
			zap.Bool("evBus", evBus == nil),
			zap.Bool("errClsf", errClsf == nil),
		)
	}

	return &worksetAppImpl{
		txnCtrl:     txnCtrl,
		worksetSvc:  worksetSvc,
		memberRepo:  memberRepo,
		worksetRepo: worksetRepo,
		evBus:       evBus,
		errClsf:     errClsf,
	}
}

// `List` returns all worksets for the given team.
func (a *worksetAppImpl) List(cx context.Context, currUid string, args *val.ListWorksetArgs) app_res.AppRes[[]val.WorksetVal] {
	lgr := app_util.TakeLgr(cx)

	if re := vfyListWorksetArgs(args); re.IsReject() {
		return app_res.Reject[[]val.WorksetVal](re.Code(), re.Msg())
	}

	// Permission check: only team members can list the team's worksets.
	if re := a.worksetSvc.CanListWorkset(currUid, args.TeamId, a.memberRepo, a.errClsf); re.IsReject() {
		return app_res.Reject[[]val.WorksetVal](app_res.ErrCode(re.Code()), re.Msg())
	}

	// Retrieve all worksets for the team.
	worksets, err := a.worksetRepo.List(
		&query.ListWorksetOpt{
			TeamId: &args.TeamId,
			Pagi: query.PagiOpt{
				Offset: args.Offset,
				Limit:  args.Limit,
			},
		},
		enum.WorksetInclTeam,
	)
	if err != nil {
		lgr.Error(
			"[worksetAppImpl.List] failed to list worksets",
			zap.Error(err),
		)

		return app_res.Reject[[]val.WorksetVal](app_res.ServerError, "获取作品集列表失败")
	}

	// Assemble the value-object slice.
	worksetVals := make([]val.WorksetVal, len(worksets))
	for i, workset := range worksets {
		worksetVals[i] = asmWorksetVal(workset)
	}

	return app_res.Accept(&worksetVals)
}

// `Create` creates a new workset inside a team.
func (a *worksetAppImpl) Create(cx context.Context, currUid string, args *val.CreateWorksetArgs) app_res.AppRes[val.WorksetCreatedRes] {
	lgr := app_util.TakeLgr(cx)

	if args == nil || args.TeamId == "" || args.Name == "" {
		return app_res.Reject[val.WorksetCreatedRes](app_res.BadRequest, "team_id name 不能为空")
	}

	// Open a transaction to atomically verify permission and create the new workset.
	re, err := repo_iface.RunWithTxn(a.txnCtrl, func(prov repo_iface.Prov) (app_res.AppRes[val.WorksetCreatedRes], error) {
		memberRepo := prov.MemberRepo()
		teamRepo := prov.TeamRepo()
		worksetRepo := prov.WorksetRepo()

		// Verify admin role via domain service.
		if re := a.worksetSvc.CanAdminWorkset(currUid, args.TeamId, memberRepo, a.errClsf); re.IsReject() {
			return app_res.Reject[val.WorksetCreatedRes](app_res.ErrCode(re.Code()), re.Msg()), app_res.DefErr()
		}

		// Allocate the next index atomically from team-level sequence.
		index, err := teamRepo.IncrementWorksetNextIndex(args.TeamId)
		if err != nil {
			return app_res.Reject[val.WorksetCreatedRes](app_res.ServerError, "创建作品集失败"), err
		}

		// Build the creation input via the domain service.
		cre := a.worksetSvc.NewWorksetCre(args.TeamId, index, args.Name, args.Desc)

		// Persist the new workset.
		workset, err := worksetRepo.Create(cre)
		if err != nil {
			return app_res.Reject[val.WorksetCreatedRes](app_res.ServerError, "创建作品集失败"), err
		}

		return app_res.Accept(&val.WorksetCreatedRes{Id: workset.Id}), nil
	})
	if err != nil {
		lgr.Error(
			"[worksetAppImpl.Create] failed to run create workset transaction",
			zap.Error(err),
		)

		return re
	}

	return re
}

// `Update` updates the name and/or description of an existing workset.
func (a *worksetAppImpl) Update(cx context.Context, currUid string, args *val.WorksetUpdArgs) app_res.AppRes[app_res.None] {
	lgr := app_util.TakeLgr(cx)

	if args == nil || args.Id == "" || args.Name == "" {
		return app_res.Reject[app_res.None](app_res.BadRequest, "id name 不能为空")
	}

	// Load the target workset first.
	workset, err := a.worksetRepo.GetById(args.Id)
	if err != nil {
		lgr.Error("[worksetAppImpl.Update] failed to get workset", zap.Error(err))

		return app_res.Reject[app_res.None](app_res.BadRequest, "作品集不存在")
	}

	// Verify admin role via domain service.
	if re := a.worksetSvc.CanAdminWorkset(currUid, workset.TeamId, a.memberRepo, a.errClsf); re.IsReject() {
		return app_res.Reject[app_res.None](app_res.ErrCode(re.Code()), re.Msg())
	}

	// Apply the `PUT` update.
	if err := a.worksetRepo.Update(&aggr.WorksetUpd{Id: args.Id, Name: args.Name, Desc: args.Desc}); err != nil {
		lgr.Error("[worksetAppImpl.Update] failed to update workset", zap.Error(err))

		return app_res.Reject[app_res.None](app_res.ServerError, "更新作品集失败")
	}

	return app_res.Accept(&app_res.None{})
}

// `Delete` hard-deletes a workset by id.
func (a *worksetAppImpl) Delete(cx context.Context, currUid string, worksetId string) app_res.AppRes[app_res.None] {
	lgr := app_util.TakeLgr(cx)

	if worksetId == "" {
		return app_res.Reject[app_res.None](app_res.BadRequest, "workset_id 不能为空")
	}

	ev := make([]event_iface.Event, 0)

	re, err := repo_iface.RunWithTxn[app_res.AppRes[app_res.None]](a.txnCtrl, func(prov repo_iface.Prov) (app_res.AppRes[app_res.None], error) {
		// Load all repo dependencies needed by the cascade delete flow.
		memberRepo := prov.MemberRepo()
		worksetRepo := prov.WorksetRepo()
		comicRepo := prov.ComicRepo()
		chapterRepo := prov.ChapterRepo()
		pageRepo := prov.PageRepo()
		assignmentRepo := prov.AssignmentRepo()
		ossMsgRepo := prov.OssMsgRepo()
		ossMsgSvc := svc.NewOssMsgSvc()

		// Load the target workset and verify admin permission first.
		workset, err := worksetRepo.GetById(worksetId)
		if err != nil {
			return app_res.Reject[app_res.None](app_res.BadRequest, "作品集不存在"), app_res.DefErr()
		}

		if re := a.worksetSvc.CanAdminWorkset(currUid, workset.TeamId, memberRepo, a.errClsf); re.IsReject() {
			return app_res.Reject[app_res.None](app_res.ErrCode(re.Code()), re.Msg()), app_res.DefErr()
		}

		// Load all comics explicitly so the cascade flow is complete for any valid cardinality.
		comics, err := listAllComics(comicRepo, workset.Id)
		if err != nil {
			return app_res.Reject[app_res.None](app_res.ServerError, "删除作品集失败"), err
		}

		for i := range comics {
			// Load all descendant chapters before deleting the comic row.
			chapters, err := listAllChapters(chapterRepo, comics[i].Id)
			if err != nil {
				return app_res.Reject[app_res.None](app_res.ServerError, "删除作品集失败"), err
			}

			for j := range chapters {
				// Delete one chapter subtree and preserve chapter-removed event payload.
				assignedUserIds, err := deleteChapterCascade(
					chapters[j],
					pageRepo,
					assignmentRepo,
					chapterRepo,
					ossMsgRepo,
					ossMsgSvc,
				)
				if err != nil {
					return app_res.Reject[app_res.None](app_res.ServerError, "删除作品集失败"), err
				}

				ev = append(ev, event.NewChapterRemovedEv(chapters[j].Id, chapters[j].PublishedAt != nil, assignedUserIds))
			}

			// Enqueue comic-cover cleanup before deleting the comic row.
			if comics[i].CoverKey != nil && *comics[i].CoverKey != "" {
				if err := ossMsgSvc.SavePendingDel(ossMsgRepo, enum.OssResComicCover, comics[i].Id, []string{*comics[i].CoverKey}); err != nil {
					return app_res.Reject[app_res.None](app_res.ServerError, "删除作品集失败"), err
				}
			}

			// Delete the comic row after all descendants have been removed.
			if err := comicRepo.Delete(comics[i].Id); err != nil {
				return app_res.Reject[app_res.None](app_res.ServerError, "删除作品集失败"), err
			}
		}

		// Delete the workset row last after the whole cascade is complete.
		if err := worksetRepo.Delete(worksetId); err != nil {
			return app_res.Reject[app_res.None](app_res.ServerError, "删除作品集失败"), err
		}

		return app_res.Accept(&app_res.None{}), nil
	})
	if err != nil {
		lgr.Error("[worksetAppImpl.Delete] failed to delete workset", zap.Error(err))

		return re
	}

	a.evBus.Pub(context.Background(), ev)

	return re
}
