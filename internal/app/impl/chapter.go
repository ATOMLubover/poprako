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
	assignmentRepo repo_iface.AssignmentRepo

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
	assignmentRepo repo_iface.AssignmentRepo,
	chapterSvc svc.ChapterSvc,
	assignmentSvc svc.AssignmentSvc,
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
		memberRepo := prov.MemberRepo()
		worksetRepo := prov.WorksetRepo()
		comicRepo := prov.ComicRepo()
		chapterRepo := prov.ChapterRepo()
		assignmentRepo := prov.AssignmentRepo()

		chapter, err := chapterRepo.GetById(args.Id)
		if err != nil {
			return app_res.Reject[app_res.None](app_res.Forbidden, "仅汉化组管理员可更新章节"), app_res.DefErr()
		}

		comic, err := comicRepo.GetById(chapter.ComicId)
		if err != nil {
			return app_res.Reject[app_res.None](app_res.Forbidden, "仅汉化组管理员可更新章节"), app_res.DefErr()
		}

		workset, err := worksetRepo.GetById(comic.WorksetId)
		if err != nil {
			return app_res.Reject[app_res.None](app_res.Forbidden, "仅汉化组管理员可更新章节"), app_res.DefErr()
		}

		if args.WorkflowTransition != nil {
			// Workflow transition requires chapter-level role permission.
			if re := a.chapterSvc.CanTransiteWorkflow(currUid, args.Id, *args.WorkflowTransition, assignmentRepo, a.errClsf); re.IsReject() {
				return app_res.Reject[app_res.None](app_res.ErrCode(re.Code()), re.Msg()), app_res.DefErr()
			}
		} else {
			// Metadata-only update requires team admin permission.
			if re := a.chapterSvc.CanAdminChapter(currUid, workset.TeamId, memberRepo, a.errClsf); re.IsReject() {
				return app_res.Reject[app_res.None](app_res.ErrCode(re.Code()), re.Msg()), app_res.DefErr()
			}
		}

		if args.WorkflowTransition != nil {
			wasPublished := chapter.PublishedAt != nil

			if err := chapter.TransiteWorkflow(*args.WorkflowTransition); err != nil {
				return app_res.Reject[app_res.None](app_res.BadRequest, "无效的工作流状态转换"), app_res.DefErr()
			}

			if !wasPublished && chapter.PublishedAt != nil {
				assignments, listErr := listAllAssignments(assignmentRepo, chapter.Id)
				if listErr != nil {
					return app_res.Reject[app_res.None](app_res.ServerError, "更新章节失败"), listErr
				}

				assignedUserIds := collectAssignedUserIds(assignments)

				ev = append(ev, event.NewChapterPublishedEv(chapter.Id, assignedUserIds))
			}
		}

		upd := mkChapterUpd(args, chapter)
		if err := chapterRepo.Update(upd); err != nil {
			return app_res.Reject[app_res.None](app_res.ServerError, "更新章节失败"), err
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
