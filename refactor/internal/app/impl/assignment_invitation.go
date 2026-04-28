package app_impl

import (
	"context"

	app_iface "poprako-s/internal/app"
	"poprako-s/internal/app/res"
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
}

// `NewAssignmentInvApp` creates one `AssignmentInvApp` implementation.
func NewAssignmentInvApp(txnCtrl repo_iface.TxnCtrl, assignmentInvSvc svc.AssignmentInvSvc, assignmentSvc svc.AssignmentSvc, userRepo repo_iface.UserRepo, memberRepo repo_iface.MemberRepo, worksetRepo repo_iface.WorksetRepo, comicRepo repo_iface.ComicRepo, chapterRepo repo_iface.ChapterRepo, assignmentInvRepo repo_iface.AssignmentInvRepo, assignmentRepo repo_iface.AssignmentRepo, evBus event_iface.EvBus) app_iface.AssignmentInvApp {
	if txnCtrl == nil || userRepo == nil || memberRepo == nil || worksetRepo == nil || comicRepo == nil || chapterRepo == nil || assignmentInvRepo == nil || assignmentRepo == nil || evBus == nil {
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
	}
}

// `ListByChapter` lists invitations under one chapter.
func (a *assignmentInvAppImpl) ListByChapter(cx context.Context, currUid string, args *val.ListAssignmentInvArgs) res.AppRes[[]val.AssignmentInvVal] {
	lgr := app_util.TakeLgr(cx)

	if args == nil || args.ChapterId == "" {
		return res.Reject[[]val.AssignmentInvVal](res.BadRequest, "chapter_id 不能为空")
	}

	if code, msg, reject := vfyAssignmentListArgs(args.Offset, &args.Limit); reject {
		return res.Reject[[]val.AssignmentInvVal](code, msg)
	}

	if code, msg, reject := ensureReviewerPermission(a.assignmentRepo, args.ChapterId, currUid); reject {
		return res.Reject[[]val.AssignmentInvVal](code, msg)
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
		return res.Reject[[]val.AssignmentInvVal](res.ServerError, "获取邀请列表失败")
	}

	vals := make([]val.AssignmentInvVal, len(items))
	for i := range items {
		v := asmAssignmentInvVal(&items[i])
		vals[i] = v
	}

	return res.Accept(&vals)
}

// `Create` creates one invitation for assignment.
func (a *assignmentInvAppImpl) Create(cx context.Context, currUid string, args *val.CreateAssignmentInvArgs) res.AppRes[val.CreateAssignmentInvRes] {
	lgr := app_util.TakeLgr(cx)

	if args == nil || args.ChapterId == "" || args.InviteeQid == "" {
		return res.Reject[val.CreateAssignmentInvRes](res.BadRequest, "chapter_id 和 invitee_qid 不能为空")
	}

	re, err := repo_iface.RunWithTxn[res.AppRes[val.CreateAssignmentInvRes]](a.txnCtrl, func(prov repo_iface.Prov) (res.AppRes[val.CreateAssignmentInvRes], error) {
		assignmentRepo := prov.AssignmentRepo()
		assignmentInvRepo := prov.AssignmentInvRepo()

		if code, msg, reject := ensureReviewerPermission(assignmentRepo, args.ChapterId, currUid); reject {
			return res.Reject[val.CreateAssignmentInvRes](code, msg), res.DefErr()
		}

		if args.RoleMask == 0 {
			return res.Reject[val.CreateAssignmentInvRes](res.BadRequest, "至少指定一个角色"), res.DefErr()
		}

		if args.RoleMask.HasAnyRole(enum.RoleAdmin) {
			return res.Reject[val.CreateAssignmentInvRes](res.BadRequest, "章节邀请不支持管理员角色"), res.DefErr()
		}

		cre, err := a.assignmentInvSvc.NewAssignmentInvCre(currUid, args.ChapterId, args.InviteeQid, args.RoleMask)
		if err != nil {
			return res.Reject[val.CreateAssignmentInvRes](res.BadRequest, err.Error()), res.DefErr()
		}

		inv, err := assignmentInvRepo.Create(cre)
		if err != nil {
			return res.Reject[val.CreateAssignmentInvRes](res.ServerError, "创建邀请失败"), err
		}

		return res.Accept(&val.CreateAssignmentInvRes{Id: inv.Id, InvCode: inv.InvCode}), nil
	})
	if err != nil {
		lgr.Error("[assignmentInvAppImpl.Create] failed to run transaction", zap.Error(err))
		return re
	}

	return re
}

// `Remove` removes one invitation by id.
func (a *assignmentInvAppImpl) Remove(cx context.Context, currUid string, invId string) res.AppRes[res.None] {
	lgr := app_util.TakeLgr(cx)

	if invId == "" {
		return res.Reject[res.None](res.BadRequest, "invitation_id 不能为空")
	}

	re, err := repo_iface.RunWithTxn[res.AppRes[res.None]](a.txnCtrl, func(prov repo_iface.Prov) (res.AppRes[res.None], error) {
		assignmentRepo := prov.AssignmentRepo()
		assignmentInvRepo := prov.AssignmentInvRepo()

		inv, err := assignmentInvRepo.GetById(invId)
		if err != nil {
			if repo_infra.IsNotFound(err) {
				return res.Reject[res.None](res.NotFound, "邀请不存在"), res.DefErr()
			}

			return res.Reject[res.None](res.ServerError, "删除邀请失败"), err
		}

		if code, msg, reject := ensureReviewerPermission(assignmentRepo, inv.ChapterId, currUid); reject {
			return res.Reject[res.None](code, msg), res.DefErr()
		}

		if err := assignmentInvRepo.Delete(invId); err != nil {
			return res.Reject[res.None](res.ServerError, "删除邀请失败"), err
		}

		return res.Accept(&res.None{}), nil
	})
	if err != nil {
		lgr.Error("[assignmentInvAppImpl.Remove] failed to run transaction", zap.Error(err))
		return re
	}

	return re
}

// `JoinByInvCode` joins one chapter by invitation code.
func (a *assignmentInvAppImpl) JoinByInvCode(cx context.Context, currUid string, args *val.JoinAssignmentInvArgs) res.AppRes[res.None] {
	lgr := app_util.TakeLgr(cx)

	if args == nil || args.InvCode == "" {
		return res.Reject[res.None](res.BadRequest, "invitation_code 不能为空")
	}

	ev := make([]event_iface.Event, 0)

	re, err := repo_iface.RunWithTxn[res.AppRes[res.None]](a.txnCtrl, func(prov repo_iface.Prov) (res.AppRes[res.None], error) {
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
				return res.Reject[res.None](res.NotFound, "用户不存在"), res.DefErr()
			}

			return res.Reject[res.None](res.ServerError, "加入章节协作失败"), err
		}

		invs, err := assignmentInvRepo.ListPendingByInviteeQid(currUser.Qid)
		if err != nil {
			return res.Reject[res.None](res.ServerError, "加入章节协作失败"), err
		}

		var target *aggr.AssignmentInv
		for i := range invs {
			if invs[i].InvCode == args.InvCode {
				target = &invs[i]
				break
			}
		}

		if target == nil || !target.Pending {
			return res.Reject[res.None](res.BadRequest, "邀请码无效或已被使用"), res.DefErr()
		}

		if code, msg, reject := ensureUserCanTakeRoles(memberRepo, chapterRepo, comicRepo, worksetRepo, target.ChapterId, currUid, target.RoleMask); reject {
			return res.Reject[res.None](code, msg), res.DefErr()
		}

		currAssignment, err := assignmentRepo.GetByChapterUserId(target.ChapterId, currUid)
		if err != nil {
			if !repo_infra.IsNotFound(err) {
				return res.Reject[res.None](res.ServerError, "加入章节协作失败"), err
			}

			cre := a.assignmentSvc.NewAssignmentCre(target.ChapterId, currUid, target.RoleMask)
			if _, err := assignmentRepo.Create(cre); err != nil {
				return res.Reject[res.None](res.ServerError, "加入章节协作失败"), err
			}

			ev = append(ev, event.NewAssignmentCreatedEv(currUid, target.ChapterId))
		} else {
			mergedMask := currAssignment.ToRoleMask() | target.RoleMask
			put := a.assignmentSvc.NewAssignmentPut(currAssignment, mergedMask)
			if err := assignmentRepo.Put(put); err != nil {
				return res.Reject[res.None](res.ServerError, "加入章节协作失败"), err
			}
		}

		if err := assignmentInvRepo.MarkCompleted(target.Id); err != nil {
			return res.Reject[res.None](res.ServerError, "加入章节协作失败"), err
		}

		return res.Accept(&res.None{}), nil
	})
	if err != nil {
		lgr.Error("[assignmentInvAppImpl.JoinByInvCode] failed to run transaction", zap.Error(err))
		return re
	}

	a.evBus.Pub(context.Background(), ev)

	return re
}
