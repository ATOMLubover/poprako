package app

import (
	"context"
	"errors"
	"time"

	event_handler "poprako-s/internal/app/event_handler"
	"poprako-s/internal/app/val"
	"poprako-s/internal/domain/event"
	"poprako-s/internal/domain/ext/oss"
	"poprako-s/internal/domain/model"
	"poprako-s/internal/domain/repo"
	"poprako-s/internal/domain/service"

	"go.uber.org/zap"
)

type AssignmentApp interface {
	// ListByChapter 获取指定章节的分配列表
	ListByChapter(
		cx context.Context,
		currUserID string,
		args *val.ListChapterAssignmentArgs,
	) ([]*val.AssignmentInfo, error)

	// ListMy 获取当前用户的所有分配列表
	ListMy(
		cx context.Context,
		currUserID string,
		args *val.ListMyAssignmentArgs,
	) ([]*val.AssignmentInfo, error)

	// Create 创建分配记录
	Create(
		cx context.Context,
		currUserID string,
		args *val.CreateAssignmentArgs,
	) (*val.CreateAssignmentRes, error)

	// Update 更新分配记录的角色
	Update(
		cx context.Context,
		currUserID string,
		args *val.UpdateAssignmentArgs,
	) error

	// Remove 删除分配记录
	Remove(
		cx context.Context,
		currUserID string,
		assignmentID string,
	) error

	// JoinInvitorChapter 通过章节邀请加入协作
	JoinInvitorChapter(
		cx context.Context,
		currUserID string,
		args *val.JoinInvitorChapterArgs,
	) error
}

type assignmentAppImpl struct {
	assignmentSvc service.AssignmentService

	assignmentRepo repo.AssignmentRepo
	chapterInvRepo repo.ChapterInvitationRepo
	chapterRepo    repo.ChapterRepo
	userRepo       repo.UserRepo
	txnMgr         repo.TxnMgr
	eventBus       event.EventBus
	ossClient      oss.Client
}

func NewAssignmentApp(
	assignmentSvc service.AssignmentService,
	assignmentRepo repo.AssignmentRepo,
	chapterInvRepo repo.ChapterInvitationRepo,
	chapterRepo repo.ChapterRepo,
	userRepo repo.UserRepo,
	txnMgr repo.TxnMgr,
	eventBus event.EventBus,
	ossClient oss.Client,
) AssignmentApp {
	// 校验构造函数依赖
	if assignmentSvc == nil ||
		assignmentRepo == nil ||
		chapterInvRepo == nil ||
		chapterRepo == nil ||
		userRepo == nil ||
		txnMgr == nil ||
		eventBus == nil ||
		ossClient == nil {
		zap.L().Panic(
			"NewAssignmentApp: 依赖项不能为空",
			zap.Bool("assignmentSvc_nil", assignmentSvc == nil),
			zap.Bool("assignmentRepo_nil", assignmentRepo == nil),
			zap.Bool("chapterInvRepo_nil", chapterInvRepo == nil),
			zap.Bool("chapterRepo_nil", chapterRepo == nil),
			zap.Bool("userRepo_nil", userRepo == nil),
			zap.Bool("txnMgr_nil", txnMgr == nil),
			zap.Bool("eventBus_nil", eventBus == nil),
			zap.Bool("ossClient_nil", ossClient == nil),
		)
	}

	// 返回真实业务实现
	return &assignmentAppImpl{
		assignmentSvc:  assignmentSvc,
		assignmentRepo: assignmentRepo,
		chapterInvRepo: chapterInvRepo,
		chapterRepo:    chapterRepo,
		userRepo:       userRepo,
		txnMgr:         txnMgr,
		eventBus:       eventBus,
		ossClient:      ossClient,
	}
}

func (a *assignmentAppImpl) ListByChapter(
	cx context.Context,
	currUserID string,
	args *val.ListChapterAssignmentArgs,
) ([]*val.AssignmentInfo, error) {
	// 获取上下文中的日志记录器
	lgr := retrieveLgr(cx)

	// 鉴权：检查当前用户是否为该章节的分配人员
	_, err := a.assignmentRepo.Get(model.AssignmentQueryOpt{
		ChapterID: &args.ChapterID,
		UserID:    &currUserID,
	})
	if err != nil {
		// 记录权限校验失败
		lgr.Warn(
			"获取分配列表失败：权限不足",
			zap.String("curr_user_id", currUserID),
			zap.String("chapter_id", args.ChapterID),
		)

		// 返回客户端可展示的错误
		return nil, errors.New("权限不足")
	}

	// 查询分配列表
	assignments, err := a.assignmentRepo.List(model.AssignmentQueryOpt{
		ChapterID: &args.ChapterID,
	})
	if err != nil {
		// 记录查询失败
		lgr.Error(
			"获取分配列表失败",
			zap.String("chapter_id", args.ChapterID),
			zap.Error(err),
		)

		// 返回客户端可展示的错误
		return nil, errors.New("获取分配列表失败")
	}

	// 组装为 app 层值对象列表
	result := make([]*val.AssignmentInfo, len(assignments))

	for i, assignment := range assignments {
		result[i] = assembleAssignmentInfo(&assignment, a.ossClient)
	}

	// 返回分配列表
	return result, nil
}

func (a *assignmentAppImpl) ListMy(
	cx context.Context,
	currUserID string,
	args *val.ListMyAssignmentArgs,
) ([]*val.AssignmentInfo, error) {
	// 获取上下文中的日志记录器
	lgr := retrieveLgr(cx)

	// 查询当前用户的所有分配列表
	assignments, err := a.assignmentRepo.List(model.AssignmentQueryOpt{
		UserID: &currUserID,
	})
	if err != nil {
		// 记录查询失败
		lgr.Error(
			"获取我的分配列表失败",
			zap.String("curr_user_id", currUserID),
			zap.Error(err),
		)

		// 返回客户端可展示的错误
		return nil, errors.New("获取分配列表失败")
	}

	// 组装为 app 层值对象列表
	result := make([]*val.AssignmentInfo, len(assignments))

	for i, assignment := range assignments {
		result[i] = assembleAssignmentInfo(&assignment, a.ossClient)
	}

	// 返回分配列表
	return result, nil
}

func (a *assignmentAppImpl) Create(
	cx context.Context,
	currUserID string,
	args *val.CreateAssignmentArgs,
) (*val.CreateAssignmentRes, error) {
	// 获取上下文中的日志记录器
	lgr := retrieveLgr(cx)

	var createdID string

	if err := a.txnMgr.RunInTxn(func(cx context.Context) error {
		assignmentRepoTxn, err := a.assignmentRepo.FromTxnCx(cx)
		if err != nil {
			return err
		}

		userRepoTxn, err := a.userRepo.FromTxnCx(cx)
		if err != nil {
			return err
		}

		creation, err := a.assignmentSvc.NewCreation(
			assignmentRepoTxn,
			currUserID,
			args.ChapterID,
			args.UserID,
			args.Roles,
		)
		if err != nil {
			return err
		}

		assignInfo, err := assignmentRepoTxn.Create(creation)
		if err != nil {
			return err
		}

		createdID = assignInfo.ID
		eventCx := event_handler.WithUserRepoTxn(cx, userRepoTxn)

		return a.eventBus.Pub(eventCx, creation.PullEvents())
	}); err != nil {
		lgr.Error(
			"创建分配失败",
			zap.String("curr_user_id", currUserID),
			zap.String("chapter_id", args.ChapterID),
			zap.String("user_id", args.UserID),
			zap.Error(err),
		)

		return nil, errors.New("创建分配失败")
	}

	return &val.CreateAssignmentRes{ID: createdID}, nil
}

func (a *assignmentAppImpl) Update(
	cx context.Context,
	currUserID string,
	args *val.UpdateAssignmentArgs,
) error {
	// 获取上下文中的日志记录器
	lgr := retrieveLgr(cx)

	// 查询当前分配信息
	currAssignInfo, err := a.assignmentRepo.GetByID(args.ID)
	if err != nil {
		// 记录查询失败
		lgr.Error(
			"更新分配失败：获取当前分配信息失败",
			zap.String("assignment_id", args.ID),
			zap.Error(err),
		)

		// 返回客户端可展示的错误
		return errors.New("无法获取分配信息")
	}

	// 通过领域服务构造更新载荷（含权限校验）
	update, err := a.assignmentSvc.NewUpdate(
		a.assignmentRepo,
		currUserID,
		args.ID,
		currAssignInfo,
		args.Roles,
	)
	if err != nil {
		// 记录权限校验失败
		lgr.Warn(
			"更新分配失败：权限不足",
			zap.String("curr_user_id", currUserID),
			zap.Error(err),
		)

		// 返回领域服务返回的错误
		return err
	}

	// 持久化更新
	if err := a.assignmentRepo.Update(update); err != nil {
		// 记录更新失败
		lgr.Error(
			"更新分配失败",
			zap.String("assignment_id", args.ID),
			zap.Error(err),
		)

		// 返回客户端可展示的错误
		return errors.New("更新分配失败")
	}

	// 返回更新成功
	return nil
}

func (a *assignmentAppImpl) Remove(
	cx context.Context,
	currUserID string,
	assignmentID string,
) error {
	// 获取上下文中的日志记录器
	lgr := retrieveLgr(cx)

	if err := a.txnMgr.RunInTxn(func(cx context.Context) error {
		assignmentRepoTxn, err := a.assignmentRepo.FromTxnCx(cx)
		if err != nil {
			return err
		}

		chapterRepoTxn, err := a.chapterRepo.FromTxnCx(cx)
		if err != nil {
			return err
		}

		userRepoTxn, err := a.userRepo.FromTxnCx(cx)
		if err != nil {
			return err
		}

		targetAssignment, err := assignmentRepoTxn.GetByID(assignmentID)
		if err != nil {
			return err
		}

		currAssignment, err := assignmentRepoTxn.Get(model.AssignmentQueryOpt{
			ChapterID: &targetAssignment.ChapterID,
			UserID:    &currUserID,
		})
		if err != nil || !currAssignment.HasAnyRole(model.RoleReviewer) {
			return errors.New("权限不足")
		}

		targetChapter, err := chapterRepoTxn.GetByID(targetAssignment.ChapterID)
		if err != nil {
			return err
		}

		if err := assignmentRepoTxn.Delete(assignmentID); err != nil {
			return err
		}

		eventCx := event_handler.WithUserRepoTxn(cx, userRepoTxn)

		removalEvent := a.assignmentSvc.NewRemovalEvent(targetAssignment, targetChapter.PublishedAt != nil)

		return a.eventBus.Pub(eventCx, []event.Event{removalEvent})
	}); err != nil {
		lgr.Error(
			"删除分配失败",
			zap.String("assignment_id", assignmentID),
			zap.Error(err),
		)

		if err.Error() == "权限不足" {
			return err
		}

		return errors.New("删除分配失败")
	}

	return nil
}

func (a *assignmentAppImpl) JoinInvitorChapter(
	cx context.Context,
	currUserID string,
	args *val.JoinInvitorChapterArgs,
) error {
	if args == nil || args.InvitationCode == "" {
		return errors.New("参数不合法")
	}

	// 获取上下文中的日志记录器
	lgr := retrieveLgr(cx)

	// 查询当前用户信息
	currUser, err := a.userRepo.GetByID(currUserID)
	if err != nil {
		lgr.Error(
			"加入章节协作失败：无法获取用户信息",
			zap.String("curr_user_id", currUserID),
			zap.Error(err),
		)

		return errors.New("加入章节协作失败：无法获取用户信息")
	}

	if err := a.txnMgr.RunInTxn(func(txCx context.Context) error {
		assignmentRepoTxn, err := a.assignmentRepo.FromTxnCx(txCx)
		if err != nil {
			return err
		}

		chapterInvRepoTxn, err := a.chapterInvRepo.FromTxnCx(txCx)
		if err != nil {
			return err
		}

		userRepoTxn, err := a.userRepo.FromTxnCx(txCx)
		if err != nil {
			return err
		}

		invs, err := chapterInvRepoTxn.List(model.ChapterInvitationQueryOpt{
			InviteeQQ:       &currUser.QQ,
			OnlyPendingTrue: true,
		})
		if err != nil {
			return err
		}

		var targetInv *model.ChapterInvitationInfo

		for i := range invs {
			if invs[i].InvitationCode == args.InvitationCode {
				targetInv = &invs[i]

				break
			}
		}

		if targetInv == nil || !targetInv.Pending {
			return errors.New("邀请码无效或已被使用")
		}

		now := time.Now()

		toAssign := func(flag bool, currAt *time.Time) *time.Time {
			if !flag {
				return currAt
			}

			if currAt != nil {
				t := *currAt
				return &t
			}

			t := now
			return &t
		}

		exists, err := assignmentRepoTxn.Exist(model.AssignmentQueryOpt{
			ChapterID: &targetInv.ChapterID,
			UserID:    &currUserID,
		})
		if err != nil {
			return err
		}

		if exists {
			existing, err := assignmentRepoTxn.Get(model.AssignmentQueryOpt{
				ChapterID: &targetInv.ChapterID,
				UserID:    &currUserID,
			})
			if err != nil {
				return err
			}

			update := &model.AssignmentUpdate{
				ID:                    existing.ID,
				AssignedRawProviderAt: toAssign(targetInv.ToBeRawProvider, existing.AssignedRawProviderAt),
				AssignedTranslatorAt:  toAssign(targetInv.ToBeTranslator, existing.AssignedTranslatorAt),
				AssignedProofreaderAt: toAssign(targetInv.ToBeProofreader, existing.AssignedProofreaderAt),
				AssignedTypesetterAt:  toAssign(targetInv.ToBeTypesetter, existing.AssignedTypesetterAt),
				AssignedRedrawerAt:    toAssign(targetInv.ToBeRedrawer, existing.AssignedRedrawerAt),
				AssignedReviewerAt:    toAssign(targetInv.ToBeReviewer, existing.AssignedReviewerAt),
				AssignedPublisherAt:   toAssign(targetInv.ToBePublisher, existing.AssignedPublisherAt),
			}

			if err := assignmentRepoTxn.Update(update); err != nil {
				return err
			}
		} else {
			creation := &model.AssignmentCreation{
				ID:                    service.GenID("assignment"),
				ChapterID:             targetInv.ChapterID,
				UserID:                currUserID,
				AssignedRawProviderAt: toAssign(targetInv.ToBeRawProvider, nil),
				AssignedTranslatorAt:  toAssign(targetInv.ToBeTranslator, nil),
				AssignedProofreaderAt: toAssign(targetInv.ToBeProofreader, nil),
				AssignedTypesetterAt:  toAssign(targetInv.ToBeTypesetter, nil),
				AssignedRedrawerAt:    toAssign(targetInv.ToBeRedrawer, nil),
				AssignedReviewerAt:    toAssign(targetInv.ToBeReviewer, nil),
				AssignedPublisherAt:   toAssign(targetInv.ToBePublisher, nil),
			}

			creation.PushEvent(&event.AssignmentCreatedEvent{
				UserID:    currUserID,
				ChapterID: targetInv.ChapterID,
			})

			if _, err := assignmentRepoTxn.Create(creation); err != nil {
				return err
			}

			eventCx := event_handler.WithUserRepoTxn(txCx, userRepoTxn)

			if err := a.eventBus.Pub(eventCx, creation.PullEvents()); err != nil {
				return err
			}
		}

		if err := chapterInvRepoTxn.Invalidate(targetInv.ID); err != nil {
			return err
		}

		return nil
	}); err != nil {
		if err.Error() == "邀请码无效或已被使用" {
			return err
		}

		lgr.Error(
			"加入章节协作失败",
			zap.String("curr_user_id", currUserID),
			zap.Error(err),
		)

		return errors.New("加入章节协作失败")
	}

	return nil
}

// assembleAssignmentInfo 将领域层分配信息转换为 app 层值对象
func assembleAssignmentInfo(
	info *model.AssignmentInfo,
	ossClient oss.Client,
) *val.AssignmentInfo {
	result := &val.AssignmentInfo{
		ID:        info.ID,
		ChapterID: info.ChapterID,
		UserID:    info.UserID,
		Roles:     info.AssignedRoleMask(),
		CreatedAt: info.CreatedAt.UnixMilli(),
		UpdatedAt: info.UpdatedAt.UnixMilli(),
	}

	// 若包含章节信息则一并组装
	if info.Chapter != nil {
		result.Chapter = assembleChapterInfo(info.Chapter, ossClient)
	}

	// 若包含用户信息则一并组装
	if info.User != nil {
		userInfo, _ := assembleUserInfo(info.User, ossClient)
		result.User = userInfo
	}

	return result
}

// logAssignmentAppImpl 是 AssignmentApp 的日志包装实现
type logAssignmentAppImpl struct {
	app AssignmentApp
}

func NewLogAssignmentApp(
	app AssignmentApp,
) AssignmentApp {
	if app == nil {
		zap.L().Panic(
			"NewLogAssignmentApp: 依赖项不能为空",
			zap.Bool("app_nil", app == nil),
		)
	}

	return &logAssignmentAppImpl{app: app}
}

func (a *logAssignmentAppImpl) ListByChapter(
	cx context.Context,
	currUserID string,
	args *val.ListChapterAssignmentArgs,
) ([]*val.AssignmentInfo, error) {
	if a == nil || a.app == nil {
		return nil, errors.New("AssignmentApp 不可用")
	}

	if args == nil || args.ChapterID == "" {
		return nil, errors.New("参数不合法")
	}

	lgr := retrieveLgr(cx).With(zap.String("method", "ListByChapter"), zap.String("curr_user_id", currUserID), zap.String("chapter_id", args.ChapterID))

	cx = injectLgr(cx, lgr)

	lgr.Info("[logAssignmentAppImpl.ListByChapter] CALL")

	return a.app.ListByChapter(cx, currUserID, args)
}

func (a *logAssignmentAppImpl) ListMy(
	cx context.Context,
	currUserID string,
	args *val.ListMyAssignmentArgs,
) ([]*val.AssignmentInfo, error) {
	if a == nil || a.app == nil {
		return nil, errors.New("AssignmentApp 不可用")
	}

	if args == nil {
		return nil, errors.New("参数不合法")
	}

	lgr := retrieveLgr(cx).With(zap.String("method", "ListMy"), zap.String("curr_user_id", currUserID))

	cx = injectLgr(cx, lgr)

	lgr.Info("[logAssignmentAppImpl.ListMy] CALL")

	return a.app.ListMy(cx, currUserID, args)
}

func (a *logAssignmentAppImpl) Create(
	cx context.Context,
	currUserID string,
	args *val.CreateAssignmentArgs,
) (*val.CreateAssignmentRes, error) {
	if a == nil || a.app == nil {
		return nil, errors.New("AssignmentApp 不可用")
	}

	if args == nil || args.ChapterID == "" || args.UserID == "" {
		return nil, errors.New("参数不合法")
	}

	lgr := retrieveLgr(cx).With(zap.String("method", "Create"), zap.String("curr_user_id", currUserID))

	cx = injectLgr(cx, lgr)

	lgr.Info("[logAssignmentAppImpl.Create] CALL")

	return a.app.Create(cx, currUserID, args)
}

func (a *logAssignmentAppImpl) Update(
	cx context.Context,
	currUserID string,
	args *val.UpdateAssignmentArgs,
) error {
	if a == nil || a.app == nil {
		return errors.New("AssignmentApp 不可用")
	}

	if args == nil || args.ID == "" {
		return errors.New("参数不合法")
	}

	lgr := retrieveLgr(cx).With(zap.String("method", "Update"), zap.String("curr_user_id", currUserID), zap.String("assignment_id", args.ID))

	cx = injectLgr(cx, lgr)

	lgr.Info("[logAssignmentAppImpl.Update] CALL")

	return a.app.Update(cx, currUserID, args)
}

func (a *logAssignmentAppImpl) Remove(
	cx context.Context,
	currUserID string,
	assignmentID string,
) error {
	if a == nil || a.app == nil {
		return errors.New("AssignmentApp 不可用")
	}

	if assignmentID == "" {
		return errors.New("分配 ID 不能为空")
	}

	lgr := retrieveLgr(cx).With(zap.String("method", "Remove"), zap.String("curr_user_id", currUserID), zap.String("assignment_id", assignmentID))

	cx = injectLgr(cx, lgr)

	lgr.Info("[logAssignmentAppImpl.Remove] CALL")

	return a.app.Remove(cx, currUserID, assignmentID)
}

func (a *logAssignmentAppImpl) JoinInvitorChapter(
	cx context.Context,
	currUserID string,
	args *val.JoinInvitorChapterArgs,
) error {
	if a == nil || a.app == nil {
		return errors.New("AssignmentApp 不可用")
	}

	if args == nil || args.InvitationCode == "" {
		return errors.New("参数不合法")
	}

	lgr := retrieveLgr(cx).With(zap.String("method", "JoinInvitorChapter"), zap.String("curr_user_id", currUserID))

	cx = injectLgr(cx, lgr)

	lgr.Info("[logAssignmentAppImpl.JoinInvitorChapter] CALL")

	return a.app.JoinInvitorChapter(cx, currUserID, args)
}
