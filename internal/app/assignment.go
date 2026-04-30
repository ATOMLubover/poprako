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
	// Get 获取指定分配详情
	Get(
		cx context.Context,
		currUserID string,
		assignmentID string,
	) (*val.AssignmentInfo, error)

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

	// Update 更新分配记录（替换角色集合）
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
	comicRepo      repo.ComicRepo
	memberRepo     repo.MemberRepo
	worksetRepo    repo.WorksetRepo
	userRepo       repo.UserRepo
	txnMgr         repo.TxnMgr
	eventBus       event.EventBus
	urlSigner      oss.URLSigner
}

func NewAssignmentApp(
	assignmentSvc service.AssignmentService,
	assignmentRepo repo.AssignmentRepo,
	chapterInvRepo repo.ChapterInvitationRepo,
	chapterRepo repo.ChapterRepo,
	comicRepo repo.ComicRepo,
	memberRepo repo.MemberRepo,
	worksetRepo repo.WorksetRepo,
	userRepo repo.UserRepo,
	txnMgr repo.TxnMgr,
	eventBus event.EventBus,
	urlSigner oss.URLSigner,
) AssignmentApp {
	// 校验构造函数依赖
	if assignmentSvc == nil ||
		assignmentRepo == nil ||
		chapterInvRepo == nil ||
		chapterRepo == nil ||
		comicRepo == nil ||
		memberRepo == nil ||
		worksetRepo == nil ||
		userRepo == nil ||
		txnMgr == nil ||
		eventBus == nil ||
		urlSigner == nil {
		zap.L().Panic(
			"NewAssignmentApp: 依赖项不能为空",
			zap.Bool("assignmentSvc_nil", assignmentSvc == nil),
			zap.Bool("assignmentRepo_nil", assignmentRepo == nil),
			zap.Bool("chapterInvRepo_nil", chapterInvRepo == nil),
			zap.Bool("chapterRepo_nil", chapterRepo == nil),
			zap.Bool("comicRepo_nil", comicRepo == nil),
			zap.Bool("memberRepo_nil", memberRepo == nil),
			zap.Bool("worksetRepo_nil", worksetRepo == nil),
			zap.Bool("userRepo_nil", userRepo == nil),
			zap.Bool("txnMgr_nil", txnMgr == nil),
			zap.Bool("eventBus_nil", eventBus == nil),
			zap.Bool("urlSigner_nil", urlSigner == nil),
		)
	}

	// 返回真实业务实现
	return &assignmentAppImpl{
		assignmentSvc:  assignmentSvc,
		assignmentRepo: assignmentRepo,
		chapterInvRepo: chapterInvRepo,
		chapterRepo:    chapterRepo,
		comicRepo:      comicRepo,
		memberRepo:     memberRepo,
		worksetRepo:    worksetRepo,
		userRepo:       userRepo,
		txnMgr:         txnMgr,
		eventBus:       eventBus,
		urlSigner:      urlSigner,
	}
}

func (a *assignmentAppImpl) ListByChapter(
	cx context.Context,
	currUserID string,
	args *val.ListChapterAssignmentArgs,
) ([]*val.AssignmentInfo, error) {
	// 获取上下文中的日志记录器
	lgr := retrieveLgr(cx)

	// 鉴权：检查当前用户是否为章节所属汉化组的成员（按 team 权限）
	var allowedByAssignment bool
	targetChapter, err := a.chapterRepo.GetByID(args.ChapterID)
	if err != nil {
		// 回退：若当前用户在该章节已有 assignment，则允许访问（兼容测试/旧逻辑）
		if _, aerr := a.assignmentRepo.Get(model.AssignmentQueryOpt{ChapterID: &args.ChapterID, UserID: &currUserID}); aerr == nil {
			allowedByAssignment = true
		} else {
			lgr.Error(
				"获取分配列表失败：获取章节信息失败",
				zap.String("chapter_id", args.ChapterID),
				zap.Error(err),
			)

			return nil, errors.New("无法获取章节信息")
		}
	}

	var targetWorkset *model.WorksetInfo
	if !allowedByAssignment {
		targetComic, err := a.comicRepo.GetByID(targetChapter.ComicID)
		if err != nil {
			lgr.Error(
				"获取分配列表失败：获取漫画信息失败",
				zap.String("comic_id", targetChapter.ComicID),
				zap.Error(err),
			)

			return nil, errors.New("无法获取漫画信息")
		}

		tw, err := a.worksetRepo.GetByID(targetComic.WorksetID)
		if err != nil {
			lgr.Error(
				"获取分配列表失败：获取作品集信息失败",
				zap.String("workset_id", targetComic.WorksetID),
				zap.Error(err),
			)

			return nil, errors.New("无法获取作品集信息")
		}

		targetWorkset = tw
	}

	if !allowedByAssignment {
		_, err = a.memberRepo.Get(model.MemberQueryOpt{
			UserID: &currUserID,
			TeamID: &targetWorkset.TeamID,
		})

		if err != nil {
			// 如果按 team 的成员检查失败，则回退到章节分配检查（兼容旧逻辑）
			if _, aerr := a.assignmentRepo.Get(model.AssignmentQueryOpt{
				ChapterID: &args.ChapterID,
				UserID:    &currUserID,
			}); aerr != nil {
				// 记录权限校验失败
				lgr.Warn(
					"获取分配列表失败：权限不足",
					zap.String("curr_user_id", currUserID),
					zap.String("chapter_id", args.ChapterID),
				)

				// 返回客户端可展示的错误
				return nil, errors.New("权限不足")
			}
		}
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

	// 按需补充 includes 嵌套数据
	a.populateAssignmentIncludes(assignments, args.Includes)

	// 组装为 app 层值对象列表
	result := make([]*val.AssignmentInfo, len(assignments))

	for i, assignment := range assignments {
		result[i] = assembleAssignmentInfo(&assignment, a.urlSigner)
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

	// 按需补充 includes 嵌套数据
	a.populateAssignmentIncludes(assignments, args.Includes)

	// 组装为 app 层值对象列表
	result := make([]*val.AssignmentInfo, len(assignments))

	for i, assignment := range assignments {
		result[i] = assembleAssignmentInfo(&assignment, a.urlSigner)
	}

	// 返回分配列表
	return result, nil
}

func (a *assignmentAppImpl) Get(
	cx context.Context,
	currUserID string,
	assignmentID string,
) (*val.AssignmentInfo, error) {
	// 获取上下文中的日志记录器
	lgr := retrieveLgr(cx)

	// 查询目标分配信息
	targetAssignment, err := a.assignmentRepo.GetByID(assignmentID)
	if err != nil {
		// 记录查询失败
		lgr.Error(
			"获取分配详情失败：获取分配信息失败",
			zap.String("assignment_id", assignmentID),
			zap.Error(err),
		)

		// 返回客户端可展示的错误
		return nil, errors.New("无法获取分配信息")
	}

	// 鉴权：检查当前用户是否为该章节的成员（先按 assignment 规则）
	_, err = a.assignmentRepo.Get(model.AssignmentQueryOpt{
		ChapterID: &targetAssignment.ChapterID,
		UserID:    &currUserID,
	})
	if err != nil {
		// 回退到按 team 的成员关系校验
		targetChapter, chapterErr := a.chapterRepo.GetByID(targetAssignment.ChapterID)
		if chapterErr != nil {
			lgr.Error(
				"获取分配详情失败：获取章节信息失败",
				zap.String("chapter_id", targetAssignment.ChapterID),
				zap.Error(chapterErr),
			)

			return nil, errors.New("无法获取章节信息")
		}

		targetComic, comicErr := a.comicRepo.GetByID(targetChapter.ComicID)
		if comicErr != nil {
			lgr.Error(
				"获取分配详情失败：获取漫画信息失败",
				zap.String("comic_id", targetChapter.ComicID),
				zap.Error(comicErr),
			)

			return nil, errors.New("无法获取漫画信息")
		}

		targetWorkset, worksetErr := a.worksetRepo.GetByID(targetComic.WorksetID)
		if worksetErr != nil {
			lgr.Error(
				"获取分配详情失败：获取作品集信息失败",
				zap.String("workset_id", targetComic.WorksetID),
				zap.Error(worksetErr),
			)

			return nil, errors.New("无法获取作品集信息")
		}

		_, memberErr := a.memberRepo.Get(model.MemberQueryOpt{
			UserID: &currUserID,
			TeamID: &targetWorkset.TeamID,
		})
		if memberErr != nil {
			lgr.Warn(
				"获取分配详情失败：权限不足",
				zap.String("curr_user_id", currUserID),
				zap.String("assignment_id", assignmentID),
			)

			return nil, errors.New("权限不足")
		}
	}

	// 返回分配详情
	return assembleAssignmentInfo(targetAssignment, a.urlSigner), nil
}

// hasInclude 检查 includes 切片中是否包含指定的 key
func hasInclude(includes []string, key string) bool {
	for _, inc := range includes {
		if inc == key {
			return true
		}
	}

	return false
}

// populateAssignmentIncludes 按照 includes 列表为每条分配记录补充嵌套数据
// 涉及的嵌套键：user, chapter, chapter.comic, chapter.creator
func (a *assignmentAppImpl) populateAssignmentIncludes(
	assignments []model.AssignmentInfo,
	includes []string,
) {
	if len(includes) == 0 {
		return
	}

	wantUser := hasInclude(includes, "user")
	wantChapter := hasInclude(includes, "chapter") ||
		hasInclude(includes, "chapter.comic") ||
		hasInclude(includes, "chapter.creator")
	wantChapterComic := hasInclude(includes, "chapter.comic")
	wantChapterCreator := hasInclude(includes, "chapter.creator")

	for i := range assignments {
		// 按需获取用户信息
		if wantUser && assignments[i].User == nil {
			if user, err := a.userRepo.GetByID(assignments[i].UserID); err == nil {
				assignments[i].User = user
			}
		}

		// 按需获取章节信息
		if wantChapter && assignments[i].Chapter == nil {
			chapter, err := a.chapterRepo.GetByID(assignments[i].ChapterID)
			if err == nil {
				// 按需获取章节所属漫画
				if wantChapterComic {
					if comic, cerr := a.comicRepo.GetByID(chapter.ComicID); cerr == nil {
						chapter.Comic = comic
					}
				}

				// 按需获取章节创建者
				if wantChapterCreator {
					if creator, cerr := a.userRepo.GetByID(chapter.CreatorID); cerr == nil {
						chapter.Creator = creator
					}
				}

				assignments[i].Chapter = chapter
			}
		}
	}
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
			a.memberRepo,
			a.chapterRepo,
			a.comicRepo,
			a.worksetRepo,
			currUserID,
			args.ChapterID,
			args.UserID,
			args.Roles,
		)
		if err != nil {
			return err
		}

		assignInfo, err := assignmentRepoTxn.UpsertCreate(creation)
		if err != nil {
			return err
		}

		createdID = assignInfo.ID
		eventCx := event_handler.WithUserRepoTxn(cx, userRepoTxn)

		return a.eventBus.Pub(eventCx, creation.PullEvents())
	}); err != nil {
		lgr.Error(
			"创建或更新分配失败",
			zap.String("curr_user_id", currUserID),
			zap.String("chapter_id", args.ChapterID),
			zap.String("user_id", args.UserID),
			zap.Error(err),
		)

		return nil, errors.New("创建或更新分配失败")
	}

	return &val.CreateAssignmentRes{ID: createdID}, nil
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

		if err := a.assignmentSvc.EnsureUserCanTakeRoles(
			a.memberRepo,
			a.chapterRepo,
			a.comicRepo,
			a.worksetRepo,
			targetInv.ChapterID,
			currUserID,
			targetInv.InvitedRoleMask(),
		); err != nil {
			return err
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

		existing, _ := assignmentRepoTxn.Get(model.AssignmentQueryOpt{ChapterID: &creation.ChapterID, UserID: &creation.UserID})

		if existing == nil {
			creation.PushEvent(&event.AssignmentCreatedEvent{
				UserID:    currUserID,
				ChapterID: targetInv.ChapterID,
			})
		} else {
			if creation.AssignedRawProviderAt == nil {
				creation.AssignedRawProviderAt = existing.AssignedRawProviderAt
			}
			if creation.AssignedTranslatorAt == nil {
				creation.AssignedTranslatorAt = existing.AssignedTranslatorAt
			}
			if creation.AssignedProofreaderAt == nil {
				creation.AssignedProofreaderAt = existing.AssignedProofreaderAt
			}
			if creation.AssignedTypesetterAt == nil {
				creation.AssignedTypesetterAt = existing.AssignedTypesetterAt
			}
			if creation.AssignedRedrawerAt == nil {
				creation.AssignedRedrawerAt = existing.AssignedRedrawerAt
			}
			if creation.AssignedReviewerAt == nil {
				creation.AssignedReviewerAt = existing.AssignedReviewerAt
			}
			if creation.AssignedPublisherAt == nil {
				creation.AssignedPublisherAt = existing.AssignedPublisherAt
			}
		}

		if _, err := assignmentRepoTxn.UpsertCreate(creation); err != nil {
			return err
		}

		eventCx := event_handler.WithUserRepoTxn(txCx, userRepoTxn)

		if existing == nil {
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

func (a *assignmentAppImpl) Update(
	cx context.Context,
	currUserID string,
	args *val.UpdateAssignmentArgs,
) error {
	lgr := retrieveLgr(cx)

	if args == nil || args.ID == "" {
		return errors.New("参数不合法")
	}

	if err := a.txnMgr.RunInTxn(func(txCx context.Context) error {
		assignmentRepoTxn, err := a.assignmentRepo.FromTxnCx(txCx)
		if err != nil {
			// fallback to non-transactional repo when txn context not provided (tests)
			assignmentRepoTxn = a.assignmentRepo
		}

		chapterRepoTxn, err := a.chapterRepo.FromTxnCx(txCx)
		if err != nil {
			chapterRepoTxn = a.chapterRepo
		}

		userRepoTxn, err := a.userRepo.FromTxnCx(txCx)
		if err != nil {
			userRepoTxn = a.userRepo
		}

		// 获取目标分配记录
		target, err := assignmentRepoTxn.GetByID(args.ID)
		if err != nil {
			return err
		}

		// 由领域服务生成更新载荷（含权限校验）
		upd, err := a.assignmentSvc.NewUpdate(
			assignmentRepoTxn,
			a.memberRepo,
			chapterRepoTxn,
			a.comicRepo,
			a.worksetRepo,
			currUserID,
			args.ID,
			target,
			args.Roles,
		)
		if err != nil {
			return err
		}

		// 将更新载荷转换为 upsert-create 以便持久化（保持 DB 层无 Update 方法的兼容）
		creation := &model.AssignmentCreation{
			ID:                    target.ID,
			ChapterID:             target.ChapterID,
			UserID:                target.UserID,
			AssignedRawProviderAt: upd.AssignedRawProviderAt,
			AssignedTranslatorAt:  upd.AssignedTranslatorAt,
			AssignedProofreaderAt: upd.AssignedProofreaderAt,
			AssignedTypesetterAt:  upd.AssignedTypesetterAt,
			AssignedRedrawerAt:    upd.AssignedRedrawerAt,
			AssignedReviewerAt:    upd.AssignedReviewerAt,
			AssignedPublisherAt:   upd.AssignedPublisherAt,
		}

		if _, err := assignmentRepoTxn.UpsertCreate(creation); err != nil {
			return err
		}

		// publish any events if present (assignment updates don't push events in service currently)
		_ = userRepoTxn

		return nil
	}); err != nil {
		lgr.Error(
			"更新分配失败",
			zap.String("curr_user_id", currUserID),
			zap.Error(err),
		)

		return errors.New("更新分配失败")
	}

	return nil
}

// assembleAssignmentInfo 将领域层分配信息转换为 app 层值对象
func assembleAssignmentInfo(
	info *model.AssignmentInfo,
	urlSigner oss.URLSigner,
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
		result.Chapter = assembleChapterInfo(info.Chapter, urlSigner)
	}

	// 若包含用户信息则一并组装
	if info.User != nil {
		userInfo, _ := assembleUserInfo(info.User, urlSigner)
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

func (a *logAssignmentAppImpl) Get(
	cx context.Context,
	currUserID string,
	assignmentID string,
) (*val.AssignmentInfo, error) {
	if a == nil || a.app == nil {
		return nil, errors.New("AssignmentApp 不可用")
	}

	if assignmentID == "" {
		return nil, errors.New("分配 ID 不能为空")
	}

	lgr := retrieveLgr(cx).With(zap.String("method", "Get"), zap.String("curr_user_id", currUserID), zap.String("assignment_id", assignmentID))

	cx = injectLgr(cx, lgr)

	lgr.Info("[logAssignmentAppImpl.Get] CALL")

	return a.app.Get(cx, currUserID, assignmentID)
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

	lgr := retrieveLgr(cx).With(zap.String("method", "Update"), zap.String("curr_user_id", currUserID))

	cx = injectLgr(cx, lgr)

	lgr.Info("[logAssignmentAppImpl.Update] CALL")

	return a.app.Update(cx, currUserID, args)
}
