package app

import (
	"context"
	"errors"

	event_handler "poprako-s/internal/app/event_handler"
	"poprako-s/internal/app/val"
	"poprako-s/internal/domain/event"
	"poprako-s/internal/domain/ext/oss"
	"poprako-s/internal/domain/model"
	"poprako-s/internal/domain/repo"
	"poprako-s/internal/domain/service"

	"go.uber.org/zap"
)

type ChapterApp interface {
	// Get 获取指定章节详情
	Get(
		cx context.Context,
		currUserID string,
		chapterID string,
	) (*val.ChapterInfo, error)

	// List 获取指定漫画的章节列表
	List(
		cx context.Context,
		currUserID string,
		args *val.ListChapterArgs,
	) ([]*val.ChapterInfo, error)

	// GetComicPinned 获取指定漫画的置顶章节信息
	// 如果没有置顶章节则返回 nil 而非错误
	GetComicPinned(
		cx context.Context,
		currUserID string,
		comicID string,
	) (*val.ChapterInfo, error)

	// Create 创建一个新的章节
	Create(
		cx context.Context,
		currUserID string,
		args *val.CreateChapterArgs,
	) (*val.CreateChapterRes, error)

	// Update 更新章节信息（含工作流转换）
	Update(
		cx context.Context,
		currUserID string,
		args *val.UpdateChapterArgs,
	) error

	// Remove 删除章节
	Remove(
		cx context.Context,
		currUserID string,
		chapterID string,
	) error

	// InviteAssignee 创建章节协作邀请
	InviteAssignee(
		cx context.Context,
		currUserID string,
		args *val.InviteChapterAssigneeArgs,
	) (*val.InviteChapterAssigneeRes, error)
}

type chapterAppImpl struct {
	chapterSvc    service.ChapterService
	chapterInvSvc service.ChapterInvitationService

	memberRepo     repo.MemberRepo
	worksetRepo    repo.WorksetRepo
	comicRepo      repo.ComicRepo
	chapterRepo    repo.ChapterRepo
	assignmentRepo repo.AssignmentRepo
	userRepo       repo.UserRepo
	pageRepo       repo.PageRepo
	chapterInvRepo repo.ChapterInvitationRepo
	txnMgr         repo.TxnMgr
	msgRepo        repo.OSSMessageRepo
	eventBus       event.EventBus
	urlSigner      oss.URLSigner
}

func NewChapterApp(
	chapterSvc service.ChapterService,
	chapterInvSvc service.ChapterInvitationService,
	memberRepo repo.MemberRepo,
	worksetRepo repo.WorksetRepo,
	comicRepo repo.ComicRepo,
	chapterRepo repo.ChapterRepo,
	assignmentRepo repo.AssignmentRepo,
	userRepo repo.UserRepo,
	pageRepo repo.PageRepo,
	chapterInvRepo repo.ChapterInvitationRepo,
	txnMgr repo.TxnMgr,
	msgRepo repo.OSSMessageRepo,
	eventBus event.EventBus,
	urlSigner oss.URLSigner,
) ChapterApp {
	// 校验构造函数依赖
	if chapterSvc == nil ||
		chapterInvSvc == nil ||
		memberRepo == nil ||
		worksetRepo == nil ||
		comicRepo == nil ||
		chapterRepo == nil ||
		assignmentRepo == nil ||
		userRepo == nil ||
		pageRepo == nil ||
		chapterInvRepo == nil ||
		txnMgr == nil ||
		msgRepo == nil ||
		eventBus == nil ||
		urlSigner == nil {
		zap.L().Panic(
			"NewChapterApp: 依赖项不能为空",
			zap.Bool("chapterSvc_nil", chapterSvc == nil),
			zap.Bool("chapterInvSvc_nil", chapterInvSvc == nil),
			zap.Bool("memberRepo_nil", memberRepo == nil),
			zap.Bool("worksetRepo_nil", worksetRepo == nil),
			zap.Bool("comicRepo_nil", comicRepo == nil),
			zap.Bool("chapterRepo_nil", chapterRepo == nil),
			zap.Bool("assignmentRepo_nil", assignmentRepo == nil),
			zap.Bool("userRepo_nil", userRepo == nil),
			zap.Bool("pageRepo_nil", pageRepo == nil),
			zap.Bool("chapterInvRepo_nil", chapterInvRepo == nil),
			zap.Bool("txnMgr_nil", txnMgr == nil),
			zap.Bool("msgRepo_nil", msgRepo == nil),
			zap.Bool("eventBus_nil", eventBus == nil),
			zap.Bool("urlSigner_nil", urlSigner == nil),
		)
	}

	// 返回真实业务实现
	return &chapterAppImpl{
		chapterSvc:     chapterSvc,
		chapterInvSvc:  chapterInvSvc,
		memberRepo:     memberRepo,
		worksetRepo:    worksetRepo,
		comicRepo:      comicRepo,
		chapterRepo:    chapterRepo,
		assignmentRepo: assignmentRepo,
		userRepo:       userRepo,
		pageRepo:       pageRepo,
		chapterInvRepo: chapterInvRepo,
		txnMgr:         txnMgr,
		msgRepo:        msgRepo,
		eventBus:       eventBus,
		urlSigner:      urlSigner,
	}
}

func (a *chapterAppImpl) List(
	cx context.Context,
	currUserID string,
	args *val.ListChapterArgs,
) ([]*val.ChapterInfo, error) {
	// 获取上下文中的日志记录器
	lgr := retrieveLgr(cx)

	// 通过漫画获取所属作品集，再获取所属汉化组 ID 用于鉴权
	targetComic, err := a.comicRepo.GetByID(args.ComicID)
	if err != nil {
		// 记录查询失败
		lgr.Error(
			"获取章节列表失败：获取漫画信息失败",
			zap.String("comic_id", args.ComicID),
			zap.Error(err),
		)

		// 返回客户端可展示的错误
		return nil, errors.New("无法获取漫画信息")
	}

	// 通过作品集获取所属汉化组 ID
	targetWorkset, err := a.worksetRepo.GetByID(targetComic.WorksetID)
	if err != nil {
		// 记录查询失败
		lgr.Error(
			"获取章节列表失败：获取作品集信息失败",
			zap.String("workset_id", targetComic.WorksetID),
			zap.Error(err),
		)

		// 返回客户端可展示的错误
		return nil, errors.New("无法获取作品集信息")
	}

	// 鉴权：检查当前用户是否为该汉化组成员
	_, err = a.memberRepo.Get(model.MemberQueryOpt{
		UserID: &currUserID,
		TeamID: &targetWorkset.TeamID,
	})
	if err != nil {
		// 记录权限校验失败
		lgr.Warn(
			"获取章节列表失败：权限不足",
			zap.String("curr_user_id", currUserID),
		)

		// 返回客户端可展示的错误
		return nil, errors.New("权限不足")
	}

	// 查询章节列表
	chapters, err := a.chapterRepo.List(model.ChapterQueryOpt{
		ComicID: &args.ComicID,
	})
	if err != nil {
		// 记录查询失败
		lgr.Error(
			"获取章节列表失败",
			zap.String("comic_id", args.ComicID),
			zap.Error(err),
		)

		// 返回客户端可展示的错误
		return nil, errors.New("获取章节列表失败")
	}

	// 组装为 app 层值对象列表
	result := make([]*val.ChapterInfo, len(chapters))

	for i, ch := range chapters {
		result[i] = assembleChapterInfo(&ch, a.urlSigner)
	}

	// 返回章节列表
	return result, nil
}

func (a *chapterAppImpl) Get(
	cx context.Context,
	currUserID string,
	chapterID string,
) (*val.ChapterInfo, error) {
	// 获取上下文中的日志记录器
	lgr := retrieveLgr(cx)

	// 查询目标章节信息
	targetChapter, err := a.chapterRepo.GetByID(chapterID)
	if err != nil {
		// 记录查询失败
		lgr.Error(
			"获取章节详情失败：获取章节信息失败",
			zap.String("chapter_id", chapterID),
			zap.Error(err),
		)

		// 返回客户端可展示的错误
		return nil, errors.New("无法获取章节信息")
	}

	// 通过漫画获取所属作品集，再获取所属汉化组 ID 用于鉴权
	targetComic, err := a.comicRepo.GetByID(targetChapter.ComicID)
	if err != nil {
		// 记录查询失败
		lgr.Error(
			"获取章节详情失败：获取漫画信息失败",
			zap.String("comic_id", targetChapter.ComicID),
			zap.Error(err),
		)

		// 返回客户端可展示的错误
		return nil, errors.New("无法获取漫画信息")
	}

	// 通过作品集获取所属汉化组 ID
	targetWorkset, err := a.worksetRepo.GetByID(targetComic.WorksetID)
	if err != nil {
		// 记录查询失败
		lgr.Error(
			"获取章节详情失败：获取作品集信息失败",
			zap.String("workset_id", targetComic.WorksetID),
			zap.Error(err),
		)

		// 返回客户端可展示的错误
		return nil, errors.New("无法获取作品集信息")
	}

	// 鉴权：检查当前用户是否为该汉化组成员
	_, err = a.memberRepo.Get(model.MemberQueryOpt{
		UserID: &currUserID,
		TeamID: &targetWorkset.TeamID,
	})
	if err != nil {
		// 记录权限校验失败
		lgr.Warn(
			"获取章节详情失败：权限不足",
			zap.String("curr_user_id", currUserID),
			zap.String("chapter_id", chapterID),
		)

		// 返回客户端可展示的错误
		return nil, errors.New("权限不足")
	}

	// 返回章节详情
	return assembleChapterInfo(targetChapter, a.urlSigner), nil
}

func (a *chapterAppImpl) GetComicPinned(
	cx context.Context,
	currUserID string,
	comicID string,
) (*val.ChapterInfo, error) {
	// 获取上下文中的日志记录器
	lgr := retrieveLgr(cx)

	// 通过漫画获取所属作品集，再获取所属汉化组 ID 用于鉴权
	targetComic, err := a.comicRepo.GetByID(comicID)
	if err != nil {
		// 记录查询失败
		lgr.Error(
			"获取置顶章节失败：获取漫画信息失败",
			zap.String("comic_id", comicID),
			zap.Error(err),
		)

		// 返回客户端可展示的错误
		return nil, errors.New("无法获取漫画信息")
	}

	// 通过作品集获取所属汉化组 ID
	targetWorkset, err := a.worksetRepo.GetByID(targetComic.WorksetID)
	if err != nil {
		// 记录查询失败
		lgr.Error(
			"获取置顶章节失败：获取作品集信息失败",
			zap.String("workset_id", targetComic.WorksetID),
			zap.Error(err),
		)

		// 返回客户端可展示的错误
		return nil, errors.New("无法获取作品集信息")
	}

	// 鉴权：检查当前用户是否为该汉化组成员
	_, err = a.memberRepo.Get(model.MemberQueryOpt{
		UserID: &currUserID,
		TeamID: &targetWorkset.TeamID,
	})
	if err != nil {
		// 记录权限校验失败
		lgr.Warn(
			"获取置顶章节失败：权限不足",
			zap.String("curr_user_id", currUserID),
		)

		// 返回客户端可展示的错误
		return nil, errors.New("权限不足")
	}

	// 查询置顶章节；若漫画尚无置顶章节则返回 nil 而非错误
	chapter, err := a.chapterRepo.FindPinnedByComicID(comicID)
	if err != nil {
		// 记录查询失败
		lgr.Error(
			"获取置顶章节失败",
			zap.String("comic_id", comicID),
			zap.Error(err),
		)

		// 返回客户端可展示的错误
		return nil, errors.New("获取置顶章节失败")
	}

	// 置顶章节不存在时返回 nil
	if chapter == nil {
		return nil, nil
	}

	// 组装为 app 层值对象并返回
	return assembleChapterInfo(chapter, a.urlSigner), nil
}

func (a *chapterAppImpl) Create(
	cx context.Context,
	currUserID string,
	args *val.CreateChapterArgs,
) (*val.CreateChapterRes, error) {
	// 获取上下文中的日志记录器
	lgr := retrieveLgr(cx)

	// 通过漫画获取所属作品集，再获取所属汉化组 ID 用于鉴权
	targetComic, err := a.comicRepo.GetByID(args.ComicID)
	if err != nil {
		// 记录查询失败
		lgr.Error(
			"创建章节失败：获取漫画信息失败",
			zap.String("comic_id", args.ComicID),
			zap.Error(err),
		)

		// 返回客户端可展示的错误
		return nil, errors.New("无法获取漫画信息")
	}

	// 通过作品集获取所属汉化组 ID
	targetWorkset, err := a.worksetRepo.GetByID(targetComic.WorksetID)
	if err != nil {
		// 记录查询失败
		lgr.Error(
			"创建章节失败：获取作品集信息失败",
			zap.String("workset_id", targetComic.WorksetID),
			zap.Error(err),
		)

		// 返回客户端可展示的错误
		return nil, errors.New("无法获取作品集信息")
	}

	// 鉴权：检查当前用户在漫画所属汉化组中是否为管理员
	currMember, err := a.memberRepo.Get(model.MemberQueryOpt{
		UserID: &currUserID,
		TeamID: &targetWorkset.TeamID,
	})
	if err != nil || !currMember.HasAnyRole(model.RoleAdmin) {
		// 记录权限校验失败
		lgr.Warn(
			"创建章节失败：权限不足",
			zap.String("curr_user_id", currUserID),
		)

		// 返回客户端可展示的错误
		return nil, errors.New("权限不足")
	}

	// 在事务中创建章节（需要 count 获取 index）
	var createdID string

	if err := a.txnMgr.RunInTxn(func(cx context.Context) error {
		chapterRepoTxn, err := a.chapterRepo.FromTxnCx(cx)
		if err != nil {
			return err
		}

		comicRepoTxn, err := a.comicRepo.FromTxnCx(cx)
		if err != nil {
			return err
		}

		assignmentRepoTxn, err := a.assignmentRepo.FromTxnCx(cx)
		if err != nil {
			return err
		}

		// 统计当前漫画下的章节数量以确定 index，并由领域服务构造创建载荷
		creation, err := a.chapterSvc.NewCreation(chapterRepoTxn, args.ComicID, args.Subtitle, currUserID)
		if err != nil {
			return err
		}

		// 持久化章节
		chInfo, err := chapterRepoTxn.Create(creation)
		if err != nil {
			return err
		}

		// 章节创建者必须立即成为该章节的监修，不能依赖事件侧路补齐。
		reviewerCreation := service.NewAssignmentService().NewInitReviewerCreation(chInfo.ID, currUserID)
		if _, err := assignmentRepoTxn.UpsertCreate(reviewerCreation); err != nil {
			return err
		}

		createdID = chInfo.ID

		eventCx := event_handler.WithComicRepoTxn(cx, comicRepoTxn)

		return a.eventBus.Pub(eventCx, creation.PullEvents())
	}); err != nil {
		// 记录创建失败
		lgr.Error(
			"创建章节失败",
			zap.String("curr_user_id", currUserID),
			zap.String("comic_id", args.ComicID),
			zap.Error(err),
		)

		// 返回客户端可展示的错误
		return nil, errors.New("创建章节失败")
	}

	// 返回创建结果
	return &val.CreateChapterRes{ID: createdID}, nil
}

func (a *chapterAppImpl) Update(
	cx context.Context,
	currUserID string,
	args *val.UpdateChapterArgs,
) error {
	// 获取上下文中的日志记录器
	lgr := retrieveLgr(cx)

	// 查询目标章节信息
	targetChapter, err := a.chapterRepo.GetByID(args.ChapterID)
	if err != nil {
		// 记录查询失败
		lgr.Error(
			"更新章节失败：获取目标章节信息失败",
			zap.String("chapter_id", args.ChapterID),
			zap.Error(err),
		)

		// 返回客户端可展示的错误
		return errors.New("无法获取章节信息")
	}

	var currAssignment *model.AssignmentInfo

	if args.WorkflowTransition != nil {
		currAssignment, err = a.assignmentRepo.Get(model.AssignmentQueryOpt{
			ChapterID: &args.ChapterID,
			UserID:    &currUserID,
		})
		if err != nil {
			lgr.Warn(
				"更新章节失败：当前用户无章节分配",
				zap.String("curr_user_id", currUserID),
				zap.String("chapter_id", args.ChapterID),
			)

			return errors.New("权限不足")
		}

		if err := a.chapterSvc.TransiteWorkflow(
			*args.WorkflowTransition,
			targetChapter,
			currAssignment,
		); err != nil {
			lgr.Warn(
				"更新章节失败：工作流转换失败",
				zap.String("chapter_id", args.ChapterID),
				zap.Error(err),
			)

			return err
		}
	}

	update := &model.ChapterUpdate{
		ID:             args.ChapterID,
		Subtitle:       args.Subtitle,
		IsPinned:       args.IsPinned,
		UploadedAt:     targetChapter.UploadedAt,
		TransalatingAt: targetChapter.TransalatingAt,
		TranslatedAt:   targetChapter.TranslatedAt,
		ProofreadingAt: targetChapter.ProofreadingAt,
		ProofreadAt:    targetChapter.ProofreadAt,
		TypesettingAt:  targetChapter.TypesettingAt,
		TypesetAt:      targetChapter.TypesetAt,
		ReviewedAt:     targetChapter.ReviewedAt,
		PublishedAt:    targetChapter.PublishedAt,
	}

	if args.WorkflowTransition != nil && *args.WorkflowTransition == model.WorkflowPublishComplete {
		if err := a.txnMgr.RunInTxn(func(cx context.Context) error {
			chapterRepoTxn, err := a.chapterRepo.FromTxnCx(cx)
			if err != nil {
				return err
			}

			assignmentRepoTxn, err := a.assignmentRepo.FromTxnCx(cx)
			if err != nil {
				return err
			}

			userRepoTxn, err := a.userRepo.FromTxnCx(cx)
			if err != nil {
				return err
			}

			if err := chapterRepoTxn.Update(update); err != nil {
				return err
			}

			eventCx := event_handler.WithUserRepoTxn(
				event_handler.WithAssignmentRepoTxn(cx, assignmentRepoTxn),
				userRepoTxn,
			)

			return a.eventBus.Pub(eventCx, targetChapter.PullEvents())
		}); err != nil {
			lgr.Error(
				"更新章节失败",
				zap.String("chapter_id", args.ChapterID),
				zap.Error(err),
			)

			return errors.New("更新章节失败")
		}
	} else {
		if err := a.chapterRepo.Update(update); err != nil {
			lgr.Error(
				"更新章节失败",
				zap.String("chapter_id", args.ChapterID),
				zap.Error(err),
			)

			return errors.New("更新章节失败")
		}

		if events := targetChapter.PullEvents(); len(events) > 0 {
			if err := a.eventBus.Pub(cx, events); err != nil {
				lgr.Error(
					"章节工作流事件发布失败",
					zap.String("chapter_id", args.ChapterID),
					zap.Error(err),
				)
			}
		}
	}

	return nil
}

func (a *chapterAppImpl) Remove(
	cx context.Context,
	currUserID string,
	chapterID string,
) error {
	// 获取上下文中的日志记录器
	lgr := retrieveLgr(cx)

	// 查询目标章节信息以获取所属漫画
	targetChapter, err := a.chapterRepo.GetByID(chapterID)
	if err != nil {
		// 记录查询失败
		lgr.Error(
			"删除章节失败：获取目标章节信息失败",
			zap.String("chapter_id", chapterID),
			zap.Error(err),
		)

		// 返回客户端可展示的错误
		return errors.New("无法获取章节信息")
	}

	// 通过漫画获取所属作品集，再获取所属汉化组 ID 用于鉴权
	targetComic, err := a.comicRepo.GetByID(targetChapter.ComicID)
	if err != nil {
		// 记录查询失败
		lgr.Error(
			"删除章节失败：获取漫画信息失败",
			zap.String("comic_id", targetChapter.ComicID),
			zap.Error(err),
		)

		// 返回客户端可展示的错误
		return errors.New("无法获取漫画信息")
	}

	// 通过作品集获取所属汉化组 ID
	targetWorkset, err := a.worksetRepo.GetByID(targetComic.WorksetID)
	if err != nil {
		// 记录查询失败
		lgr.Error(
			"删除章节失败：获取作品集信息失败",
			zap.String("workset_id", targetComic.WorksetID),
			zap.Error(err),
		)

		// 返回客户端可展示的错误
		return errors.New("无法获取作品集信息")
	}

	// 鉴权：检查当前用户在章节所属汉化组中是否为管理员
	currMember, err := a.memberRepo.Get(model.MemberQueryOpt{
		UserID: &currUserID,
		TeamID: &targetWorkset.TeamID,
	})
	if err != nil || !currMember.HasAnyRole(model.RoleAdmin) {
		// 记录权限校验失败
		lgr.Warn(
			"删除章节失败：权限不足",
			zap.String("curr_user_id", currUserID),
		)

		// 返回客户端可展示的错误
		return errors.New("权限不足")
	}

	pages, err := a.pageRepo.List(model.PageQueryOpt{
		ChapterID: &chapterID,
	})
	if err != nil {
		lgr.Error(
			"删除章节失败：获取页面信息失败",
			zap.String("chapter_id", chapterID),
			zap.Error(err),
		)

		return errors.New("删除章节失败")
	}

	pageOSSKeys := make([]string, 0, len(pages))
	for _, page := range pages {
		pageOSSKeys = append(pageOSSKeys, page.OSSKey)
	}

	if err := a.txnMgr.RunInTxn(func(cx context.Context) error {
		chapterRepoTxn, err := a.chapterRepo.FromTxnCx(cx)
		if err != nil {
			return err
		}

		assignmentRepoTxn, err := a.assignmentRepo.FromTxnCx(cx)
		if err != nil {
			return err
		}

		comicRepoTxn, err := a.comicRepo.FromTxnCx(cx)
		if err != nil {
			return err
		}

		userRepoTxn, err := a.userRepo.FromTxnCx(cx)
		if err != nil {
			return err
		}

		assignments, err := assignmentRepoTxn.List(model.AssignmentQueryOpt{
			ChapterID: &chapterID,
		})
		if err != nil {
			return err
		}

		assignedUserIDs := make([]string, 0, len(assignments))
		for _, assignment := range assignments {
			assignedUserIDs = append(assignedUserIDs, assignment.UserID)
		}

		if err := chapterRepoTxn.Remove(chapterID); err != nil {
			return err
		}

		if err := enqueueDeleteBatchMessage(a.msgRepo, cx, model.OSSResourcePageImage, chapterID, pageOSSKeys); err != nil {
			return err
		}

		eventCx := event_handler.WithUserRepoTxn(
			event_handler.WithComicRepoTxn(cx, comicRepoTxn),
			userRepoTxn,
		)

		removalEvent := a.chapterSvc.NewRemovalEvent(targetChapter, assignedUserIDs)

		return a.eventBus.Pub(eventCx, []event.Event{removalEvent})
	}); err != nil {
		lgr.Error(
			"删除章节失败",
			zap.String("chapter_id", chapterID),
			zap.Error(err),
		)

		return errors.New("删除章节失败")
	}

	// 返回删除成功
	return nil
}

func (a *chapterAppImpl) InviteAssignee(
	cx context.Context,
	currUserID string,
	args *val.InviteChapterAssigneeArgs,
) (*val.InviteChapterAssigneeRes, error) {
	if args == nil || args.ChapterID == "" || args.InviteeQQ == "" {
		return nil, errors.New("参数不合法")
	}

	// 获取上下文中的日志记录器
	lgr := retrieveLgr(cx)

	errNoChapter := errors.New("无法获取章节信息")
	errForbidden := errors.New("权限不足")

	invCode := ""

	if err := a.txnMgr.RunInTxn(func(txCx context.Context) error {
		chapterRepoTxn, err := a.chapterRepo.FromTxnCx(txCx)
		if err != nil {
			return err
		}

		assignmentRepoTxn, err := a.assignmentRepo.FromTxnCx(txCx)
		if err != nil {
			return err
		}

		chapterInvRepoTxn, err := a.chapterInvRepo.FromTxnCx(txCx)
		if err != nil {
			return err
		}

		targetChapter, err := chapterRepoTxn.GetByID(args.ChapterID)
		if err != nil {
			return errNoChapter
		}

		currAssignment, err := assignmentRepoTxn.Get(model.AssignmentQueryOpt{
			ChapterID: &targetChapter.ID,
			UserID:    &currUserID,
		})
		if err != nil || !currAssignment.HasAnyRole(model.RoleReviewer) {
			return errForbidden
		}

		creation, err := a.chapterInvSvc.NewCreation(
			currUserID,
			args.ChapterID,
			args.InviteeQQ,
			model.UnmaskRoles(args.Roles)...,
		)
		if err != nil {
			return err
		}

		created, err := chapterInvRepoTxn.Create(creation)
		if err != nil {
			return err
		}

		invCode = created.InvitationCode

		return nil
	}); err != nil {
		switch {
		case errors.Is(err, errNoChapter), errors.Is(err, errForbidden):
			if errors.Is(err, errForbidden) {
				lgr.Warn(
					"创建章节邀请失败：权限不足",
					zap.String("curr_user_id", currUserID),
					zap.String("chapter_id", args.ChapterID),
				)
			}

			return nil, err
		default:
			lgr.Error(
				"创建章节邀请失败",
				zap.String("curr_user_id", currUserID),
				zap.String("chapter_id", args.ChapterID),
				zap.String("invitee_qq", args.InviteeQQ),
				zap.Error(err),
			)

			return nil, errors.New("创建章节邀请失败")
		}
	}

	return &val.InviteChapterAssigneeRes{InvCode: invCode}, nil
}

// assembleChapterInfo 将领域层章节信息转换为 app 层值对象
func assembleChapterInfo(
	info *model.ChapterInfo,
	urlSigner oss.URLSigner,
) *val.ChapterInfo {
	result := &val.ChapterInfo{
		ID:                  info.ID,
		ComicID:             info.ComicID,
		IsPinned:            info.IsPinned,
		Index:               info.Index,
		Subtitle:            info.Subtitle,
		PageCount:           info.PageCount,
		TotalUnitCount:      info.TotalUnitCount,
		TranslatedUnitCount: info.TranslatedUnitCount,
		ProofreadUnitCount:  info.ProofreadUnitCount,
		CreatorID:           info.CreatorID,
		CreatedAt:           info.CreatedAt.UnixMilli(),
		UpdatedAt:           info.UpdatedAt.UnixMilli(),
	}

	// 组装工作流时间戳
	if info.UploadedAt != nil {
		ms := info.UploadedAt.UnixMilli()
		result.UploadedAt = &ms
	}

	if info.TransalatingAt != nil {
		ms := info.TransalatingAt.UnixMilli()
		result.TransalatingAt = &ms
	}

	if info.TranslatedAt != nil {
		ms := info.TranslatedAt.UnixMilli()
		result.TranslatedAt = &ms
	}

	if info.ProofreadingAt != nil {
		ms := info.ProofreadingAt.UnixMilli()
		result.ProofreadingAt = &ms
	}

	if info.ProofreadAt != nil {
		ms := info.ProofreadAt.UnixMilli()
		result.ProofreadAt = &ms
	}

	if info.TypesettingAt != nil {
		ms := info.TypesettingAt.UnixMilli()
		result.TypesettingAt = &ms
	}

	if info.TypesetAt != nil {
		ms := info.TypesetAt.UnixMilli()
		result.TypesetAt = &ms
	}

	if info.ReviewedAt != nil {
		ms := info.ReviewedAt.UnixMilli()
		result.ReviewedAt = &ms
	}

	if info.PublishedAt != nil {
		ms := info.PublishedAt.UnixMilli()
		result.PublishedAt = &ms
	}

	// 若包含创建者信息则一并组装
	if info.Creator != nil {
		userInfo, _ := assembleUserInfo(info.Creator, urlSigner)
		result.Creator = userInfo
	}

	// 若包含漫画信息则一并组装
	if info.Comic != nil {
		result.Comic = assembleComicInfo(info.Comic, urlSigner)
	}

	return result
}

// logChapterAppImpl 是 ChapterApp 的日志包装实现
type logChapterAppImpl struct {
	app ChapterApp
}

func NewLogChapterApp(
	app ChapterApp,
) ChapterApp {
	if app == nil {
		zap.L().Panic(
			"NewLogChapterApp: 依赖项不能为空",
			zap.Bool("app_nil", app == nil),
		)
	}

	return &logChapterAppImpl{app: app}
}

func (a *logChapterAppImpl) GetComicPinned(
	cx context.Context,
	currUserID string,
	comicID string,
) (*val.ChapterInfo, error) {
	if a == nil || a.app == nil {
		return nil, errors.New("ChapterApp 不可用")
	}

	if comicID == "" {
		return nil, errors.New("漫画 ID 不能为空")
	}

	lgr := retrieveLgr(cx).With(
		zap.String("method", "GetComicPinned"),
		zap.String("curr_user_id", currUserID),
		zap.String("comic_id", comicID),
	)

	cx = injectLgr(cx, lgr)

	lgr.Info("[logChapterAppImpl.GetComicPinned] CALL")

	return a.app.GetComicPinned(cx, currUserID, comicID)
}

func (a *logChapterAppImpl) Get(
	cx context.Context,
	currUserID string,
	chapterID string,
) (*val.ChapterInfo, error) {
	if a == nil || a.app == nil {
		return nil, errors.New("ChapterApp 不可用")
	}

	if chapterID == "" {
		return nil, errors.New("章节 ID 不能为空")
	}

	lgr := retrieveLgr(cx).With(
		zap.String("method", "Get"),
		zap.String("curr_user_id", currUserID),
		zap.String("chapter_id", chapterID),
	)

	cx = injectLgr(cx, lgr)

	lgr.Info("[logChapterAppImpl.Get] CALL")

	return a.app.Get(cx, currUserID, chapterID)
}

func (a *logChapterAppImpl) List(
	cx context.Context,
	currUserID string,
	args *val.ListChapterArgs,
) ([]*val.ChapterInfo, error) {
	if a == nil || a.app == nil {
		return nil, errors.New("ChapterApp 不可用")
	}

	if args == nil || args.ComicID == "" {
		return nil, errors.New("参数不合法")
	}

	lgr := retrieveLgr(cx).With(zap.String("method", "List"), zap.String("curr_user_id", currUserID), zap.String("comic_id", args.ComicID))

	cx = injectLgr(cx, lgr)

	lgr.Info("[logChapterAppImpl.List] CALL")

	return a.app.List(cx, currUserID, args)
}

func (a *logChapterAppImpl) Create(
	cx context.Context,
	currUserID string,
	args *val.CreateChapterArgs,
) (*val.CreateChapterRes, error) {
	if a == nil || a.app == nil {
		return nil, errors.New("ChapterApp 不可用")
	}

	if args == nil || args.ComicID == "" {
		return nil, errors.New("参数不合法")
	}

	lgr := retrieveLgr(cx).With(zap.String("method", "Create"), zap.String("curr_user_id", currUserID))

	cx = injectLgr(cx, lgr)

	lgr.Info("[logChapterAppImpl.Create] CALL")

	return a.app.Create(cx, currUserID, args)
}

func (a *logChapterAppImpl) Update(
	cx context.Context,
	currUserID string,
	args *val.UpdateChapterArgs,
) error {
	if a == nil || a.app == nil {
		return errors.New("ChapterApp 不可用")
	}

	if args == nil || args.ChapterID == "" {
		return errors.New("参数不合法")
	}

	lgr := retrieveLgr(cx).With(zap.String("method", "Update"), zap.String("curr_user_id", currUserID), zap.String("chapter_id", args.ChapterID))

	cx = injectLgr(cx, lgr)

	lgr.Info("[logChapterAppImpl.Update] CALL")

	return a.app.Update(cx, currUserID, args)
}

func (a *logChapterAppImpl) Remove(
	cx context.Context,
	currUserID string,
	chapterID string,
) error {
	if a == nil || a.app == nil {
		return errors.New("ChapterApp 不可用")
	}

	if chapterID == "" {
		return errors.New("章节 ID 不能为空")
	}

	lgr := retrieveLgr(cx).With(zap.String("method", "Remove"), zap.String("curr_user_id", currUserID), zap.String("chapter_id", chapterID))

	cx = injectLgr(cx, lgr)

	lgr.Info("[logChapterAppImpl.Remove] CALL")

	return a.app.Remove(cx, currUserID, chapterID)
}

func (a *logChapterAppImpl) InviteAssignee(
	cx context.Context,
	currUserID string,
	args *val.InviteChapterAssigneeArgs,
) (*val.InviteChapterAssigneeRes, error) {
	if a == nil || a.app == nil {
		return nil, errors.New("ChapterApp 不可用")
	}

	if args == nil || args.ChapterID == "" || args.InviteeQQ == "" {
		return nil, errors.New("参数不合法")
	}

	lgr := retrieveLgr(cx).With(
		zap.String("method", "InviteAssignee"),
		zap.String("curr_user_id", currUserID),
		zap.String("chapter_id", args.ChapterID),
		zap.String("invitee_qq", args.InviteeQQ),
	)

	cx = injectLgr(cx, lgr)

	lgr.Info("[logChapterAppImpl.InviteAssignee] CALL")

	return a.app.InviteAssignee(cx, currUserID, args)
}
