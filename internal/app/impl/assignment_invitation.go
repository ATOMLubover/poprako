package app_impl

import (
	"context"

	app_iface "poprako-s/internal/app"
	app_res "poprako-s/internal/app/res"
	app_util "poprako-s/internal/app/util"
	"poprako-s/internal/app/val"
	"poprako-s/internal/domain/model/aggr"
	"poprako-s/internal/domain/model/event"
	"poprako-s/internal/domain/model/query"
	repo_iface "poprako-s/internal/domain/repo"
	"poprako-s/internal/domain/svc"
	event_iface "poprako-s/internal/event"
	repo_infra "poprako-s/internal/infra/repo"

	"go.uber.org/zap"
)

// `assignmentInvAppImpl` is the default implementation of `AssignmentInvApp`.
type assignmentInvAppImpl struct {
	txnCtrl repo_iface.TxnCtrl

	assignmentInvSvc svc.AssignmentInvSvc
	assignmentSvc    svc.AssignmentSvc

	userRepo          repo_iface.UserRepo
	memberRepo        repo_iface.MemberRepo
	worksetRepo       repo_iface.WorksetRepo
	comicRepo         repo_iface.ComicRepo
	chapterRepo       repo_iface.ChapterRepo
	assignmentInvRepo repo_iface.AssignmentInvRepo
	assignmentRepo    repo_iface.AssignmentRepo

	evBus event_iface.EvBus

	errClsf repo_iface.ErrClsf
}

// `NewAssignmentInvApp` creates one `AssignmentInvApp` implementation.
func NewAssignmentInvApp(
	txnCtrl repo_iface.TxnCtrl,
	userRepo repo_iface.UserRepo,
	memberRepo repo_iface.MemberRepo,
	worksetRepo repo_iface.WorksetRepo,
	comicRepo repo_iface.ComicRepo,
	chapterRepo repo_iface.ChapterRepo,
	assignmentInvRepo repo_iface.AssignmentInvRepo,
	assignmentRepo repo_iface.AssignmentRepo,
	assignmentInvSvc svc.AssignmentInvSvc,
	assignmentSvc svc.AssignmentSvc,
	evBus event_iface.EvBus,
	errClsf repo_iface.ErrClsf,
) app_iface.AssignmentInvApp {
	if txnCtrl == nil || userRepo == nil || memberRepo == nil || worksetRepo == nil || comicRepo == nil || chapterRepo == nil || assignmentInvRepo == nil || assignmentRepo == nil || evBus == nil || errClsf == nil {
		zap.L().Panic(
			"[NewAssignmentInvApp] nil dependency",
			zap.Bool("txnCtrl", txnCtrl == nil),
			zap.Bool("userRepo", userRepo == nil),
			zap.Bool("memberRepo", memberRepo == nil),
			zap.Bool("worksetRepo", worksetRepo == nil),
			zap.Bool("comicRepo", comicRepo == nil),
			zap.Bool("chapterRepo", chapterRepo == nil),
			zap.Bool("assignmentInvRepo", assignmentInvRepo == nil),
			zap.Bool("assignmentRepo", assignmentRepo == nil),
			zap.Bool("evBus", evBus == nil),
			zap.Bool("errClsf", errClsf == nil),
		)
	}

	return &assignmentInvAppImpl{
		txnCtrl:           txnCtrl,
		assignmentInvSvc:  assignmentInvSvc,
		assignmentSvc:     assignmentSvc,
		userRepo:          userRepo,
		memberRepo:        memberRepo,
		worksetRepo:       worksetRepo,
		comicRepo:         comicRepo,
		chapterRepo:       chapterRepo,
		assignmentInvRepo: assignmentInvRepo,
		assignmentRepo:    assignmentRepo,
		evBus:             evBus,
		errClsf:           errClsf,
	}
}

// `ListByChapter` lists invitations under one chapter.
func (a *assignmentInvAppImpl) ListByChapter(cx context.Context, currUid string, args *val.ListAssignmentInvArgs) app_res.AppRes[[]val.AssignmentInvVal] {
	lgr := app_util.TakeLgr(cx)

	if args == nil || args.ChapterId == "" {
		return app_res.Reject[[]val.AssignmentInvVal](app_res.BadRequest, "chapter_id 不能为空")
	}

	if re := vfyAssignmentListArgs(args.Offset, &args.Limit); re.IsReject() {
		return app_res.Reject[[]val.AssignmentInvVal](re.Code(), re.Msg())
	}

	if re := a.assignmentSvc.CanReviewAssignment(currUid, args.ChapterId, a.assignmentRepo, a.errClsf); re.IsReject() {
		return app_res.Reject[[]val.AssignmentInvVal](app_res.ErrCode(re.Code()), re.Msg())
	}

	items, err := a.assignmentInvRepo.List(query.ListAssignmentInvOpt{
		ChapterId: args.ChapterId,
		Pending:   args.Pending,
		Pagi: query.PagiOpt{
			Offset: args.Offset,
			Limit:  args.Limit,
		},
	})
	if err != nil {
		lgr.Error("[assignmentInvAppImpl.ListByChapter] failed to list invitations", zap.Error(err))
		return app_res.Reject[[]val.AssignmentInvVal](app_res.ServerError, "获取邀请列表失败")
	}

	invitationVals := make([]val.AssignmentInvVal, len(items))
	for i := range items {
		invitationVal := asmAssignmentInvVal(&items[i])
		invitationVals[i] = invitationVal
	}

	return app_res.Accept(&invitationVals)
}

// `Create` creates one invitation for assignment.
func (a *assignmentInvAppImpl) Create(cx context.Context, currUid string, args *val.CreateAssignmentInvArgs) app_res.AppRes[val.CreateAssignmentInvRes] {
	lgr := app_util.TakeLgr(cx)

	if args == nil || args.ChapterId == "" || args.InviteeQid == "" {
		return app_res.Reject[val.CreateAssignmentInvRes](app_res.BadRequest, "chapter_id 和 invitee_qid 不能为空")
	}

	re, err := repo_iface.RunWithTxn[app_res.AppRes[val.CreateAssignmentInvRes]](a.txnCtrl, func(prov repo_iface.Prov) (app_res.AppRes[val.CreateAssignmentInvRes], error) {
		userRepo := prov.UserRepo()
		assignmentRepo := prov.AssignmentRepo()
		assignmentInvRepo := prov.AssignmentInvRepo()

		if re := a.assignmentSvc.CanReviewAssignment(currUid, args.ChapterId, assignmentRepo, a.errClsf); re.IsReject() {
			return app_res.Reject[val.CreateAssignmentInvRes](app_res.ErrCode(re.Code()), re.Msg()), app_res.DefErr()
		}

		if re := a.assignmentInvSvc.VfyInviteeNotAssigned(args.InviteeQid, args.ChapterId, userRepo, assignmentRepo, a.errClsf); re.IsReject() {
			return app_res.Reject[val.CreateAssignmentInvRes](app_res.ErrCode(re.Code()), re.Msg()), app_res.DefErr()
		}

		cre, err := a.assignmentInvSvc.NewAssignmentInvCre(currUid, args.ChapterId, args.InviteeQid, args.RoleMask)
		if err != nil {
			return app_res.Reject[val.CreateAssignmentInvRes](app_res.BadRequest, err.Error()), app_res.DefErr()
		}

		inv, err := assignmentInvRepo.Create(cre)
		if err != nil {
			return app_res.Reject[val.CreateAssignmentInvRes](app_res.ServerError, "创建邀请失败"), err
		}

		return app_res.Accept(&val.CreateAssignmentInvRes{Id: inv.Id, InvCode: inv.InvCode}), nil
	})
	if err != nil {
		lgr.Error("[assignmentInvAppImpl.Create] failed to run transaction", zap.Error(err))
		return re
	}

	return re
}

// `Delete` deletes one invitation by id.
func (a *assignmentInvAppImpl) Delete(cx context.Context, currUid string, invId string) app_res.AppRes[app_res.None] {
	lgr := app_util.TakeLgr(cx)

	if invId == "" {
		return app_res.Reject[app_res.None](app_res.BadRequest, "invitation_id 不能为空")
	}

	re, err := repo_iface.RunWithTxn[app_res.AppRes[app_res.None]](a.txnCtrl, func(prov repo_iface.Prov) (app_res.AppRes[app_res.None], error) {
		assignmentRepo := prov.AssignmentRepo()
		assignmentInvRepo := prov.AssignmentInvRepo()

		inv, err := assignmentInvRepo.GetById(invId)
		if err != nil {
			if repo_infra.IsNotFound(err) {
				return app_res.Reject[app_res.None](app_res.NotFound, "邀请不存在"), app_res.DefErr()
			}

			return app_res.Reject[app_res.None](app_res.ServerError, "删除邀请失败"), err
		}

		if re := a.assignmentSvc.CanReviewAssignment(currUid, inv.ChapterId, assignmentRepo, a.errClsf); re.IsReject() {
			return app_res.Reject[app_res.None](app_res.ErrCode(re.Code()), re.Msg()), app_res.DefErr()
		}

		if err := assignmentInvRepo.Delete(invId); err != nil {
			return app_res.Reject[app_res.None](app_res.ServerError, "删除邀请失败"), err
		}

		return app_res.Accept(&app_res.None{}), nil
	})
	if err != nil {
		lgr.Error("[assignmentInvAppImpl.Delete] failed to run transaction", zap.Error(err))
		return re
	}

	return re
}

// `JoinByInvCode` joins one chapter by invitation code.
func (a *assignmentInvAppImpl) JoinByInvCode(cx context.Context, currUid string, args *val.JoinAssignmentInvArgs) app_res.AppRes[app_res.None] {
	lgr := app_util.TakeLgr(cx)

	if args == nil || args.InvCode == "" {
		return app_res.Reject[app_res.None](app_res.BadRequest, "invitation_code 不能为空")
	}

	ev := make([]event_iface.Event, 0)

	re, err := repo_iface.RunWithTxn[app_res.AppRes[app_res.None]](a.txnCtrl, func(prov repo_iface.Prov) (app_res.AppRes[app_res.None], error) {
		userRepo := prov.UserRepo()
		memberRepo := prov.MemberRepo()
		worksetRepo := prov.WorksetRepo()
		comicRepo := prov.ComicRepo()
		chapterRepo := prov.ChapterRepo()
		assignmentInvRepo := prov.AssignmentInvRepo()
		assignmentRepo := prov.AssignmentRepo()

		currUser, err := userRepo.GetById(currUid)
		if err != nil {
			if repo_infra.IsNotFound(err) {
				return app_res.Reject[app_res.None](app_res.NotFound, "用户不存在"), app_res.DefErr()
			}

			return app_res.Reject[app_res.None](app_res.ServerError, "加入章节协作失败"), err
		}

		invs, err := assignmentInvRepo.ListPendingByInviteeQid(currUser.Qid)
		if err != nil {
			return app_res.Reject[app_res.None](app_res.ServerError, "加入章节协作失败"), err
		}

		var target *aggr.AssignmentInv
		for i := range invs {
			if invs[i].InvCode == args.InvCode {
				target = &invs[i]
				break
			}
		}

		if target == nil || !target.Pending {
			return app_res.Reject[app_res.None](app_res.BadRequest, "邀请码无效或已被使用"), app_res.DefErr()
		}

		if re := a.assignmentSvc.CanTakeAssignmentRoles(currUid, target.ChapterId, target.RoleMask, memberRepo, chapterRepo, comicRepo, worksetRepo, a.errClsf); re.IsReject() {
			return app_res.Reject[app_res.None](app_res.ErrCode(re.Code()), re.Msg()), app_res.DefErr()
		}

		currAssignment, err := assignmentRepo.GetByChapterUserId(target.ChapterId, currUid)
		if err != nil {
			if !repo_infra.IsNotFound(err) {
				return app_res.Reject[app_res.None](app_res.ServerError, "加入章节协作失败"), err
			}

			cre := a.assignmentSvc.NewAssignmentCre(target.ChapterId, currUid, target.RoleMask)
			if _, err := assignmentRepo.Create(cre); err != nil {
				return app_res.Reject[app_res.None](app_res.ServerError, "加入章节协作失败"), err
			}

			ev = append(ev, event.NewAssignmentCreatedEv(currUid, target.ChapterId))
		} else {
			mergedMask := currAssignment.ToRoleMask() | target.RoleMask
			put := a.assignmentSvc.NewAssignmentPut(currAssignment, mergedMask)
			if err := assignmentRepo.Put(put); err != nil {
				return app_res.Reject[app_res.None](app_res.ServerError, "加入章节协作失败"), err
			}
		}

		if err := assignmentInvRepo.MarkCompleted(target.Id); err != nil {
			return app_res.Reject[app_res.None](app_res.ServerError, "加入章节协作失败"), err
		}

		return app_res.Accept(&app_res.None{}), nil
	})
	if err != nil {
		lgr.Error("[assignmentInvAppImpl.JoinByInvCode] failed to run transaction", zap.Error(err))
		return re
	}

	a.evBus.Pub(context.Background(), ev)

	return re
}
