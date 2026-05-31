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

// `chapterAppImpl` is the default implementation of `ChapterApp`.
type chapterAppImpl struct {
	txnCtrl repo_iface.TxnCtrl

	chapterSvc    svc.ChapterSvc
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

// `NewChapterApp` creates one `ChapterApp` implementation.
func NewChapterApp(
	txnCtrl repo_iface.TxnCtrl,
	memberRepo repo_iface.MemberRepo,
	worksetRepo repo_iface.WorksetRepo,
	comicRepo repo_iface.ComicRepo,
	chapterRepo repo_iface.ChapterRepo,
	pageRepo repo_iface.PageRepo,
	assignmentRepo repo_iface.AssignmentRepo,
	chapterSvc svc.ChapterSvc,
	assignmentSvc svc.AssignmentSvc,
	ossSigner oss_iface.Signer,
	evBus event_iface.EvBus,
	errClsf repo_iface.ErrClsf,
) app_iface.ChapterApp {
	if txnCtrl == nil ||
		memberRepo == nil ||
		worksetRepo == nil ||
		comicRepo == nil ||
		chapterRepo == nil ||
		pageRepo == nil ||
		assignmentRepo == nil ||
		ossSigner == nil ||
		evBus == nil ||
		errClsf == nil {
		zap.L().Panic(
			"[NewChapterApp] nil dependency",
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

	return &chapterAppImpl{
		txnCtrl:        txnCtrl,
		chapterSvc:     chapterSvc,
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

// `mkChapterRepoIncl` composes repo include options from app args.
func mkChapterRepoIncl(includes []enum.ChapterIncl) []enum.ChapterIncl {
	if len(includes) == 0 {
		return nil
	}

	repoIncls := make([]enum.ChapterIncl, 0, len(includes))

	for i := range includes {
		repoIncls = append(repoIncls, includes[i])

		switch includes[i] {
		case enum.ChapterInclComicWorkset:
			repoIncls = append(repoIncls, enum.ChapterInclComic)

		case enum.ChapterInclComicWorksetTeam:
			repoIncls = append(repoIncls, enum.ChapterInclComic, enum.ChapterInclComicWorkset)

		case enum.ChapterInclComicCreator:
			repoIncls = append(repoIncls, enum.ChapterInclComic)
		}
	}

	return repoIncls
}

// `List` returns chapter list under target comic.
func (a *chapterAppImpl) List(cx context.Context, currUid string, args *val.ListChapterArgs) app_res.AppRes[[]val.ChapterVal] {
	lgr := app_util.TakeLgr(cx)

	if re := vfyListChapterArgs(args); re.IsReject() {
		return app_res.Reject[[]val.ChapterVal](re.Code(), re.Msg())
	}

	comic, err := a.comicRepo.GetById(args.ComicId, enum.ComicInclWorkset)
	if err != nil {
		return app_res.Reject[[]val.ChapterVal](app_res.BadRequest, "漫画不存在")
	}

	if re := a.chapterSvc.CanListChapter(currUid, comic.Workset.TeamId, a.memberRepo, a.errClsf); re.IsReject() {
		return app_res.Reject[[]val.ChapterVal](app_res.ErrCode(re.Code()), re.Msg())
	}

	chapters, err := a.chapterRepo.List(&query.ListChapterOpt{
		ComicId: &args.ComicId,
		Pagi: query.PagiOpt{
			Offset: args.Offset,
			Limit:  args.Limit,
		},
	}, mkChapterRepoIncl(args.Includes)...)
	if err != nil {
		lgr.Error("[chapterAppImpl.List] failed to list chapters", zap.Error(err))
		return app_res.Reject[[]val.ChapterVal](app_res.ServerError, "获取章节列表失败")
	}

	chapterVals := make([]val.ChapterVal, len(chapters))
	for i, chapter := range chapters {
		chapterVals[i] = asmChapterVal(chapter)
	}

	if err := tryFillCoverForChapters(
		chapterVals,
		a.ossSigner,
		memoizeComicCoverKeyGetter(mkPinnedFirstPageImageKeyGetter(a.chapterRepo, a.pageRepo, lgr)),
		lgr,
	); err != nil {
		errCode := mapComicFallbackErrCode(err, a.errClsf)

		lgr.Error(
			"[chapterAppImpl.List] failed to fill chapter comic covers",
			zap.Error(err),
			zap.Int("errCode", int(errCode)),
		)

		return app_res.Reject[[]val.ChapterVal](errCode, "获取章节列表失败")
	}

	return app_res.Accept(&chapterVals)
}

// `GetById` returns one chapter by id.
func (a *chapterAppImpl) GetById(cx context.Context, currUid string, args *val.GetChapterByIdArgs) app_res.AppRes[val.ChapterVal] {
	lgr := app_util.TakeLgr(cx)

	if args == nil {
		return app_res.Reject[val.ChapterVal](app_res.BadRequest, "查询参数不能为空")
	}

	if re := vfyChapterId(args.ChapterId); re.IsReject() {
		return app_res.Reject[val.ChapterVal](re.Code(), re.Msg())
	}

	chapter, err := a.chapterRepo.GetById(args.ChapterId, mkChapterRepoIncl(args.Includes)...)
	if err != nil {
		if repo_infra.IsNotFound(err) {
			return app_res.Reject[val.ChapterVal](app_res.NotFound, "章节不存在")
		}

		lgr.Error("[chapterAppImpl.GetById] failed to get chapter", zap.Error(err))

		return app_res.Reject[val.ChapterVal](app_res.ServerError, "获取章节失败")
	}

	comic, err := a.comicRepo.GetById(chapter.ComicId, enum.ComicInclWorkset)
	if err != nil {
		lgr.Error("[chapterAppImpl.GetById] failed to get comic", zap.Error(err))

		return app_res.Reject[val.ChapterVal](app_res.ServerError, "获取章节失败")
	}

	if re := a.chapterSvc.CanListChapter(currUid, comic.Workset.TeamId, a.memberRepo, a.errClsf); re.IsReject() {
		return app_res.Reject[val.ChapterVal](app_res.ErrCode(re.Code()), re.Msg())
	}

	chapterVal := asmChapterVal(chapter)
	if err := tryFillCoverOnChapter(
		&chapterVal,
		a.ossSigner,
		mkPinnedFirstPageImageKeyGetter(a.chapterRepo, a.pageRepo, lgr),
		lgr,
	); err != nil {
		errCode := mapComicFallbackErrCode(err, a.errClsf)

		lgr.Error(
			"[chapterAppImpl.GetById] failed to fill chapter comic cover",
			zap.Error(err),
			zap.Int("errCode", int(errCode)),
		)

		return app_res.Reject[val.ChapterVal](errCode, "获取章节失败")
	}

	return app_res.Accept(&chapterVal)
}

// `GetPinned` returns pinned chapter under target comic.
func (a *chapterAppImpl) GetPinned(cx context.Context, currUid string, comicId string) app_res.AppRes[val.ChapterVal] {
	lgr := app_util.TakeLgr(cx)

	if comicId == "" {
		return app_res.Reject[val.ChapterVal](app_res.BadRequest, "comic_id 不能为空")
	}

	comic, err := a.comicRepo.GetById(comicId, enum.ComicInclWorkset)
	if err != nil {
		return app_res.Reject[val.ChapterVal](app_res.BadRequest, "漫画不存在")
	}

	if re := a.chapterSvc.CanListChapter(currUid, comic.Workset.TeamId, a.memberRepo, a.errClsf); re.IsReject() {
		return app_res.Reject[val.ChapterVal](app_res.ErrCode(re.Code()), re.Msg())
	}

	chapter, err := a.chapterRepo.FindPinnedByComicId(comicId)
	if err != nil {
		lgr.Error("[chapterAppImpl.GetPinned] failed to get pinned chapter", zap.Error(err))
		return app_res.Reject[val.ChapterVal](app_res.ServerError, "获取置顶章节失败")
	}

	if chapter == nil {
		return app_res.Accept[val.ChapterVal](nil)
	}

	chapterVal := asmChapterVal(chapter)
	if err := tryFillCoverOnChapter(
		&chapterVal,
		a.ossSigner,
		mkPinnedFirstPageImageKeyGetter(a.chapterRepo, a.pageRepo, lgr),
		lgr,
	); err != nil {
		errCode := mapComicFallbackErrCode(err, a.errClsf)

		lgr.Error(
			"[chapterAppImpl.GetPinned] failed to fill chapter comic cover",
			zap.Error(err),
			zap.Int("errCode", int(errCode)),
		)

		return app_res.Reject[val.ChapterVal](errCode, "获取置顶章节失败")
	}

	return app_res.Accept(&chapterVal)
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

		comic, err := comicRepo.GetById(args.ComicId)
		if err != nil {
			return app_res.Reject[val.ChapterCreatedRes](app_res.Forbidden, "仅汉化组管理员可创建章节"), app_res.DefErr()
		}

		workset, err := worksetRepo.GetById(comic.WorksetId)
		if err != nil {
			return app_res.Reject[val.ChapterCreatedRes](app_res.Forbidden, "仅汉化组管理员可创建章节"), app_res.DefErr()
		}

		if re := a.chapterSvc.CanAdminChapter(currUid, workset.TeamId, memberRepo, a.errClsf); re.IsReject() {
			return app_res.Reject[val.ChapterCreatedRes](app_res.ErrCode(re.Code()), re.Msg()), app_res.DefErr()
		}

		index, err := comicRepo.IncrementChapterNextIndex(args.ComicId)
		if err != nil {
			return app_res.Reject[val.ChapterCreatedRes](app_res.ServerError, "创建章节失败"), err
		}

		cre := a.chapterSvc.NewChapterCre(args.ComicId, index, args.Subtitle, currUid)
		chapter, err := chapterRepo.Create(cre)
		if err != nil {
			return app_res.Reject[val.ChapterCreatedRes](app_res.ServerError, "创建章节失败"), err
		}

		if err := comicRepo.UpdateChapterCount(args.ComicId, 1); err != nil {
			return app_res.Reject[val.ChapterCreatedRes](app_res.ServerError, "创建章节失败"), err
		}

		if err := comicRepo.TouchLastActive(args.ComicId); err != nil {
			return app_res.Reject[val.ChapterCreatedRes](app_res.ServerError, "创建章节失败"), err
		}

		reviewerCre := a.assignmentSvc.NewAssignmentCre(chapter.Id, currUid, aggr.RoleMask(enum.RoleReviewer))
		if _, err := assignmentRepo.Create(reviewerCre); err != nil {
			return app_res.Reject[val.ChapterCreatedRes](app_res.ServerError, "创建章节失败"), err
		}

		ev = append(ev, event.NewAssignmentCreatedEv(currUid, chapter.Id))

		return app_res.Accept(&val.ChapterCreatedRes{Id: chapter.Id}), nil
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
		comicRepo := prov.ComicRepo()
		chapterRepo := prov.ChapterRepo()
		pageRepo := prov.PageRepo()
		assignmentRepo := prov.AssignmentRepo()
		ossMsgRepo := prov.OssMsgRepo()

		chapter, err := chapterRepo.GetById(args.Id)
		if err != nil {
			return app_res.Reject[app_res.None](app_res.Forbidden, "仅汉化组管理员可更新章节"), app_res.DefErr()
		}

		if args.WorkflowTransition != nil {
			// Workflow transition requires chapter-level role permission.
			if re := a.chapterSvc.CanTransiteWorkflow(currUid, args.Id, *args.WorkflowTransition, assignmentRepo, a.errClsf); re.IsReject() {
				return app_res.Reject[app_res.None](app_res.ErrCode(re.Code()), re.Msg()), app_res.DefErr()
			}
		} else if args.RevertTransition != nil {
			// Revert transition requires chapter-level role permission scoped to revert rules
			if re := a.chapterSvc.CanRevertWorkflow(currUid, args.Id, *args.RevertTransition, assignmentRepo, a.errClsf); re.IsReject() {
				return app_res.Reject[app_res.None](app_res.ErrCode(re.Code()), re.Msg()), app_res.DefErr()
			}
		} else {
			// Metadata-only update requires chapter-level `RoleReviewer`.
			if re := a.chapterSvc.CanUpdateChapter(currUid, args.Id, assignmentRepo, a.errClsf); re.IsReject() {
				return app_res.Reject[app_res.None](app_res.ErrCode(re.Code()), re.Msg()), app_res.DefErr()
			}
		}

		clearPublishedImages := false

		if args.WorkflowTransition != nil {
			wasPublished := chapter.PublishedAt != nil

			if err := chapter.TransiteWorkflow(*args.WorkflowTransition); err != nil {
				return app_res.Reject[app_res.None](app_res.BadRequest, "无效的工作流状态转换"), app_res.DefErr()
			}

			ev = append(ev, chapter.PullEv()...)

			if !wasPublished && chapter.PublishedAt != nil {
				clearPublishedImages = true
			}
		}

		if args.RevertTransition != nil {
			if err := chapter.RevertWorkflow(*args.RevertTransition); err != nil {
				return app_res.Reject[app_res.None](app_res.BadRequest, "无效的工作流回退转换"), app_res.DefErr()
			}

			ev = append(ev, chapter.PullEv()...)
		}

		upd := mkChapterUpd(args, chapter)
		if err := chapterRepo.Update(upd); err != nil {
			if repo_infra.IsConditionalUpdateFailed(err) {
				return app_res.Reject[app_res.None](app_res.BadRequest, "workflow 状态已变更，请刷新后重试"), app_res.DefErr()
			}

			return app_res.Reject[app_res.None](app_res.ServerError, "更新章节失败"), err
		}

		if clearPublishedImages {
			limit := chapter.PageCount
			if limit <= 0 {
				limit = 10000
			}

			pages, listErr := pageRepo.List(&query.ListPageOpt{
				ChapterId: &chapter.Id,
				Pagi:      query.PagiOpt{Limit: limit},
			})
			if listErr != nil {
				return app_res.Reject[app_res.None](app_res.ServerError, "更新章节失败"), listErr
			}

			ossMsgSvc := svc.NewOssMsgSvc()

			for i := range pages {
				if pages[i].ImageKey == nil || *pages[i].ImageKey == "" {
					continue
				}

				if err := ossMsgSvc.SavePendingDel(ossMsgRepo, enum.OssResPageImage, pages[i].Id, []string{*pages[i].ImageKey}); err != nil {
					return app_res.Reject[app_res.None](app_res.ServerError, "更新章节失败"), err
				}
			}

			if err := pageRepo.ClearImagesByChapterId(chapter.Id); err != nil {
				return app_res.Reject[app_res.None](app_res.ServerError, "更新章节失败"), err
			}
		}

		if err := comicRepo.TouchLastActive(chapter.ComicId); err != nil {
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

// `Join` adds current user to chapter assignment by role-mask union.
func (a *chapterAppImpl) Join(cx context.Context, currUid string, args val.JoinChapterArgs) app_res.AppRes[val.AssignmentVal] {
	lgr := app_util.TakeLgr(cx)

	if re := vfyJoinChapterArgs(args); re.IsReject() {
		return app_res.Reject[val.AssignmentVal](re.Code(), re.Msg())
	}

	ev := make([]event_iface.Event, 0)

	re, err := repo_iface.RunWithTxn[app_res.AppRes[val.AssignmentVal]](a.txnCtrl, func(prov repo_iface.Prov) (app_res.AppRes[val.AssignmentVal], error) {
		memberRepo := prov.MemberRepo()
		worksetRepo := prov.WorksetRepo()
		comicRepo := prov.ComicRepo()
		chapterRepo := prov.ChapterRepo()
		assignmentRepo := prov.AssignmentRepo()

		if re := a.assignmentSvc.CanTakeAssignmentRoles(currUid, args.ChapterId, args.RoleMask, memberRepo, chapterRepo, comicRepo, worksetRepo, a.errClsf); re.IsReject() {
			return app_res.Reject[val.AssignmentVal](app_res.ErrCode(re.Code()), re.Msg()), app_res.DefErr()
		}

		currAssignment, err := assignmentRepo.GetByChapterUserId(args.ChapterId, currUid)
		if err != nil {
			if !repo_infra.IsNotFound(err) {
				return app_res.Reject[val.AssignmentVal](app_res.ServerError, "加入章节失败"), err
			}

			cre := a.assignmentSvc.NewAssignmentCre(args.ChapterId, currUid, args.RoleMask)
			created, err := assignmentRepo.Create(cre)
			if err != nil {
				return app_res.Reject[val.AssignmentVal](app_res.ServerError, "加入章节失败"), err
			}

			ev = append(ev, event.NewAssignmentCreatedEv(currUid, args.ChapterId))

			assignmentVal := asmAssignmentVal(created)
			if err := tryFillCoverOnAssignment(
				&assignmentVal,
				a.ossSigner,
				mkPinnedFirstPageImageKeyGetter(a.chapterRepo, a.pageRepo, lgr),
				lgr,
			); err != nil {
				return app_res.Reject[val.AssignmentVal](mapComicFallbackErrCode(err, a.errClsf), "加入章节失败"), err
			}

			return app_res.Accept(&assignmentVal), nil
		}

		mergedMask := currAssignment.ToRoleMask() | args.RoleMask
		put := a.assignmentSvc.NewAssignmentPut(currAssignment, mergedMask)
		if err := assignmentRepo.Put(put); err != nil {
			return app_res.Reject[val.AssignmentVal](app_res.ServerError, "加入章节失败"), err
		}

		updated, err := assignmentRepo.GetById(currAssignment.Id)
		if err != nil {
			return app_res.Reject[val.AssignmentVal](app_res.ServerError, "加入章节失败"), err
		}

		assignmentVal := asmAssignmentVal(updated)
		if err := tryFillCoverOnAssignment(
			&assignmentVal,
			a.ossSigner,
			mkPinnedFirstPageImageKeyGetter(a.chapterRepo, a.pageRepo, lgr),
			lgr,
		); err != nil {
			return app_res.Reject[val.AssignmentVal](mapComicFallbackErrCode(err, a.errClsf), "加入章节失败"), err
		}

		return app_res.Accept(&assignmentVal), nil
	})
	if err != nil {
		lgr.Error("[chapterAppImpl.Join] failed to run join chapter transaction", zap.Error(err))

		return re
	}

	a.evBus.Pub(context.Background(), ev)

	return re
}

// `Delete` hard-deletes one chapter.
func (a *chapterAppImpl) Delete(cx context.Context, currUid string, chapterId string) app_res.AppRes[app_res.None] {
	lgr := app_util.TakeLgr(cx)

	if re := vfyDeleteChapterId(chapterId); re.IsReject() {
		return app_res.Reject[app_res.None](re.Code(), re.Msg())
	}

	ev := make([]event_iface.Event, 0)

	re, err := repo_iface.RunWithTxn[app_res.AppRes[app_res.None]](a.txnCtrl, func(prov repo_iface.Prov) (app_res.AppRes[app_res.None], error) {
		memberRepo := prov.MemberRepo()
		worksetRepo := prov.WorksetRepo()
		comicRepo := prov.ComicRepo()
		chapterRepo := prov.ChapterRepo()
		pageRepo := prov.PageRepo()
		ossMsgRepo := prov.OssMsgRepo()
		assignmentRepo := prov.AssignmentRepo()

		chapter, err := chapterRepo.GetById(chapterId)
		if err != nil {
			return app_res.Reject[app_res.None](app_res.Forbidden, "仅汉化组管理员可删除章节"), app_res.DefErr()
		}

		comic, err := comicRepo.GetById(chapter.ComicId)
		if err != nil {
			return app_res.Reject[app_res.None](app_res.Forbidden, "仅汉化组管理员可删除章节"), app_res.DefErr()
		}

		workset, err := worksetRepo.GetById(comic.WorksetId)
		if err != nil {
			return app_res.Reject[app_res.None](app_res.Forbidden, "仅汉化组管理员可删除章节"), app_res.DefErr()
		}

		if re := a.chapterSvc.CanAdminChapter(currUid, workset.TeamId, memberRepo, a.errClsf); re.IsReject() {
			return app_res.Reject[app_res.None](app_res.ErrCode(re.Code()), re.Msg()), app_res.DefErr()
		}

		// Load assignments before deleting rows so chapter-removed payload stays complete.
		assignments, err := listAllAssignments(assignmentRepo, chapter.Id)
		if err != nil {
			return app_res.Reject[app_res.None](app_res.ServerError, "删除章节失败"), err
		}

		assignedUserIds := collectAssignedUserIds(assignments)

		// Load pages before deleting rows so OSS cleanup stays complete.
		pages, err := listAllPages(pageRepo, chapter.Id)
		if err != nil {
			return app_res.Reject[app_res.None](app_res.ServerError, "删除章节失败"), err
		}

		// Enqueue page-image cleanup before deleting local rows.
		ossMsgSvc := svc.NewOssMsgSvc()
		if err := savePageImageDeleteMsgs(ossMsgSvc, ossMsgRepo, pages); err != nil {
			return app_res.Reject[app_res.None](app_res.ServerError, "删除章节失败"), err
		}

		// Delete descendant rows before deleting the chapter row itself.
		if err := pageRepo.DeleteByChapterId(chapter.Id); err != nil {
			return app_res.Reject[app_res.None](app_res.ServerError, "删除章节失败"), err
		}

		if err := chapterRepo.SetPageCount(chapter.Id, 0); err != nil {
			return app_res.Reject[app_res.None](app_res.ServerError, "删除章节失败"), err
		}

		if err := assignmentRepo.DeleteByChapterId(chapter.Id); err != nil {
			return app_res.Reject[app_res.None](app_res.ServerError, "删除章节失败"), err
		}

		if err := chapterRepo.Delete(chapterId); err != nil {
			return app_res.Reject[app_res.None](app_res.ServerError, "删除章节失败"), err
		}

		// Promote one remaining chapter to pinned when the deleted chapter was pinned.
		if chapter.IsPinned {
			remainingChapters, err := chapterRepo.List(&query.ListChapterOpt{
				ComicId: &chapter.ComicId,
				Pagi: query.PagiOpt{
					Limit: 1,
				},
			})
			if err != nil {
				return app_res.Reject[app_res.None](app_res.ServerError, "删除章节失败"), err
			}

			if len(remainingChapters) > 0 {
				isPinned := true

				if err := chapterRepo.Update(&aggr.ChapterUpd{
					Id:       remainingChapters[0].Id,
					IsPinned: &isPinned,
				}); err != nil {
					return app_res.Reject[app_res.None](app_res.ServerError, "删除章节失败"), err
				}
			}
		}

		if err := comicRepo.UpdateChapterCount(chapter.ComicId, -1); err != nil {
			return app_res.Reject[app_res.None](app_res.ServerError, "删除章节失败"), err
		}

		if err := comicRepo.TouchLastActive(chapter.ComicId); err != nil {
			return app_res.Reject[app_res.None](app_res.ServerError, "删除章节失败"), err
		}

		ev = append(ev, event.NewChapterRemovedEv(chapter.Id, chapter.PublishedAt != nil, assignedUserIds))

		return app_res.Accept(&app_res.None{}), nil
	})
	if err != nil {
		lgr.Error("[chapterAppImpl.Delete] failed to run delete chapter transaction", zap.Error(err))

		return re
	}

	a.evBus.Pub(context.Background(), ev)

	return re
}
