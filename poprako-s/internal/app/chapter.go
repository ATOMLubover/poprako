package app

import (
	"context"
	"errors"

	"poprako-s/internal/app/val"
	"poprako-s/internal/domain/event"
	"poprako-s/internal/domain/ext/oss"
	"poprako-s/internal/domain/model"
	"poprako-s/internal/domain/repo"
	"poprako-s/internal/domain/service"

	"go.uber.org/zap"
)

type ChapterApp interface {
	// List 获取指定漫画的章节列表
	List(
		cx context.Context,
		currUserID string,
		args *val.ListChapterArgs,
	) ([]*val.ChapterInfo, error)

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
}

type chapterAppImpl struct {
	chapterSvc service.ChapterService

	memberRepo     repo.MemberRepo
	worksetRepo    repo.WorksetRepo
	comicRepo      repo.ComicRepo
	chapterRepo    repo.ChapterRepo
	assignmentRepo repo.AssignmentRepo
	pageRepo       repo.PageRepo
	txnMgr         repo.TxnMgr
	eventBus       event.EventBus
	ossClient      oss.Client
}

func NewChapterApp(
	chapterSvc service.ChapterService,
	memberRepo repo.MemberRepo,
	worksetRepo repo.WorksetRepo,
	comicRepo repo.ComicRepo,
	chapterRepo repo.ChapterRepo,
	assignmentRepo repo.AssignmentRepo,
	pageRepo repo.PageRepo,
	txnMgr repo.TxnMgr,
	eventBus event.EventBus,
	ossClient oss.Client,
) ChapterApp {
	// 校验构造函数依赖
	if chapterSvc == nil ||
		memberRepo == nil ||
		worksetRepo == nil ||
		comicRepo == nil ||
		chapterRepo == nil ||
		assignmentRepo == nil ||
		pageRepo == nil ||
		txnMgr == nil ||
		eventBus == nil ||
		ossClient == nil {
		zap.L().Panic(
			"NewChapterApp: 依赖项不能为空",
			zap.Bool("chapterSvc_nil", chapterSvc == nil),
			zap.Bool("memberRepo_nil", memberRepo == nil),
			zap.Bool("worksetRepo_nil", worksetRepo == nil),
			zap.Bool("comicRepo_nil", comicRepo == nil),
			zap.Bool("chapterRepo_nil", chapterRepo == nil),
			zap.Bool("assignmentRepo_nil", assignmentRepo == nil),
			zap.Bool("pageRepo_nil", pageRepo == nil),
			zap.Bool("txnMgr_nil", txnMgr == nil),
			zap.Bool("eventBus_nil", eventBus == nil),
			zap.Bool("ossClient_nil", ossClient == nil),
		)
	}

	// 返回真实业务实现
	return &chapterAppImpl{
		chapterSvc:     chapterSvc,
		memberRepo:     memberRepo,
		worksetRepo:    worksetRepo,
		comicRepo:      comicRepo,
		chapterRepo:    chapterRepo,
		assignmentRepo: assignmentRepo,
		pageRepo:       pageRepo,
		txnMgr:         txnMgr,
		eventBus:       eventBus,
		ossClient:      ossClient,
	}
}

func (a *chapterAppImpl) List(
	cx context.Context,
	currUserID string,
	args *val.ListChapterArgs,
) ([]*val.ChapterInfo, error) {
	// 获取上下文中的日志记录器
	lgr := retrieveLgr(cx)

	// 通过漫画获取所属作品集，再获取所属团队 ID 用于鉴权
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

	// 通过作品集获取所属团队 ID
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

	// 鉴权：检查当前用户是否为该团队成员
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
		result[i] = assembleChapterInfo(&ch, a.ossClient)
	}

	// 返回章节列表
	return result, nil
}

func (a *chapterAppImpl) Create(
	cx context.Context,
	currUserID string,
	args *val.CreateChapterArgs,
) (*val.CreateChapterRes, error) {
	// 获取上下文中的日志记录器
	lgr := retrieveLgr(cx)

	// 通过漫画获取所属作品集，再获取所属团队 ID 用于鉴权
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

	// 通过作品集获取所属团队 ID
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

	// 鉴权：检查当前用户在漫画所属团队中是否为管理员
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
		// 统计当前漫画下的章节数量以确定 index
		count, err := a.chapterRepo.Count(model.ChapterQueryOpt{
			ComicID: &args.ComicID,
		})
		if err != nil {
			return err
		}

		// 构造章节创建载荷
		creation := &model.ChapterCreation{
			ID:        service.GenID("chapter"),
			ComicID:   args.ComicID,
			Index:     int(count),
			Subtitle:  args.Subtitle,
			CreatorID: currUserID,
		}

		// 持久化章节
		chInfo, err := a.chapterRepo.Create(creation)
		if err != nil {
			return err
		}

		createdID = chInfo.ID

		return nil
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

	// 若需要执行工作流转换
	if args.WorkflowTransition != nil {
		// 鉴权：检查当前用户在章节中是否有对应角色权限
		currAssignment, err := a.assignmentRepo.Get(model.AssignmentQueryOpt{
			ChapterID: &args.ChapterID,
			UserID:    &currUserID,
		})
		if err != nil {
			// 记录权限校验失败
			lgr.Warn(
				"更新章节失败：当前用户无章节分配",
				zap.String("curr_user_id", currUserID),
				zap.String("chapter_id", args.ChapterID),
			)

			// 返回客户端可展示的错误
			return errors.New("权限不足")
		}

		// 通过领域服务执行工作流转换（含权限校验和状态变更）
		if err := a.chapterSvc.TransiteWorkflow(
			*args.WorkflowTransition,
			targetChapter,
			currAssignment,
		); err != nil {
			// 记录工作流转换失败
			lgr.Warn(
				"更新章节失败：工作流转换失败",
				zap.String("chapter_id", args.ChapterID),
				zap.Error(err),
			)

			// 返回领域服务返回的错误
			return err
		}

		// 发布领域事件
		if events := targetChapter.Events(); len(events) > 0 {
			if err := a.eventBus.PubAsync(events); err != nil {
				// 记录事件发布失败（非阻塞）
				lgr.Error(
					"章节工作流事件发布失败",
					zap.String("chapter_id", args.ChapterID),
					zap.Error(err),
				)
			}
		}
	}

	// 构造章节更新载荷
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

	// 持久化更新
	if err := a.chapterRepo.UpdateStats(&model.ChapterStats{
		ChapterID:           args.ChapterID,
		TotalUnitCount:      targetChapter.TotalUnitCount,
		TranslatedUnitCount: targetChapter.TranslatedUnitCount,
		ProofreadUnitCount:  targetChapter.ProofreadUnitCount,
	}); err != nil {
		// 记录更新统计失败（非致命）
		lgr.Warn(
			"更新章节统计失败",
			zap.String("chapter_id", args.ChapterID),
			zap.Error(err),
		)
	}

	_ = update

	// 返回更新成功
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

	// 通过漫画获取所属作品集，再获取所属团队 ID 用于鉴权
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

	// 通过作品集获取所属团队 ID
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

	// 鉴权：检查当前用户在章节所属团队中是否为管理员
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

	// 执行删除
	if err := a.chapterRepo.Remove(chapterID); err != nil {
		// 记录删除失败
		lgr.Error(
			"删除章节失败",
			zap.String("chapter_id", chapterID),
			zap.Error(err),
		)

		// 返回客户端可展示的错误
		return errors.New("删除章节失败")
	}

	// 异步清理章节页面的 OSS 资源
	go a.cleanupChapterPages(chapterID)

	// 返回删除成功
	return nil
}

// cleanupChapterPages 异步清理章节下所有页面的 OSS 资源
func (a *chapterAppImpl) cleanupChapterPages(chapterID string) {
	// 查询章节下所有页面
	pages, err := a.pageRepo.List(model.PageQueryOpt{
		ChapterID: &chapterID,
	})
	if err != nil {
		// 记录查询失败
		zap.L().Error(
			"清理章节页面失败：查询页面失败",
			zap.String("chapter_id", chapterID),
			zap.Error(err),
		)

		return
	}

	// 逐个删除页面 OSS 资源
	for _, page := range pages {
		if page.OSSKey != "" {
			if err := a.ossClient.Delete(page.OSSKey); err != nil {
				// 记录删除 OSS 资源失败（继续处理其他页面）
				zap.L().Error(
					"清理章节页面失败：删除 OSS 资源失败",
					zap.String("chapter_id", chapterID),
					zap.String("page_id", page.ID),
					zap.String("oss_key", page.OSSKey),
					zap.Error(err),
				)
			}
		}
	}
}

// assembleChapterInfo 将领域层章节信息转换为 app 层值对象
func assembleChapterInfo(
	info *model.ChapterInfo,
	ossClient oss.Client,
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
		userInfo, _ := assembleUserInfo(info.Creator, ossClient)
		result.Creator = userInfo
	}

	return result
}

// logChapterAppImpl 是 ChapterApp 的日志包装实现
type logChapterAppImpl struct {
	lgr *zap.Logger
	app ChapterApp
}

func NewLogChapterApp(
	lgr *zap.Logger,
	app ChapterApp,
) ChapterApp {
	if lgr == nil || app == nil {
		zap.L().Panic(
			"NewLogChapterApp: 依赖项不能为空",
			zap.Bool("lgr_nil", lgr == nil),
			zap.Bool("app_nil", app == nil),
		)
	}

	return &logChapterAppImpl{lgr: lgr, app: app}
}

func (a *logChapterAppImpl) List(
	cx context.Context,
	currUserID string,
	args *val.ListChapterArgs,
) ([]*val.ChapterInfo, error) {
	if a == nil || a.app == nil || a.lgr == nil {
		return nil, errors.New("ChapterApp 不可用")
	}

	if args == nil || args.ComicID == "" {
		return nil, errors.New("参数不合法")
	}

	lgr := a.lgr.With(zap.String("method", "List"), zap.String("curr_user_id", currUserID), zap.String("comic_id", args.ComicID))

	cx = injectLgr(cx, lgr)

	lgr.Info("[logChapterAppImpl.List] CALL")

	return a.app.List(cx, currUserID, args)
}

func (a *logChapterAppImpl) Create(
	cx context.Context,
	currUserID string,
	args *val.CreateChapterArgs,
) (*val.CreateChapterRes, error) {
	if a == nil || a.app == nil || a.lgr == nil {
		return nil, errors.New("ChapterApp 不可用")
	}

	if args == nil || args.ComicID == "" {
		return nil, errors.New("参数不合法")
	}

	lgr := a.lgr.With(zap.String("method", "Create"), zap.String("curr_user_id", currUserID))

	cx = injectLgr(cx, lgr)

	lgr.Info("[logChapterAppImpl.Create] CALL")

	return a.app.Create(cx, currUserID, args)
}

func (a *logChapterAppImpl) Update(
	cx context.Context,
	currUserID string,
	args *val.UpdateChapterArgs,
) error {
	if a == nil || a.app == nil || a.lgr == nil {
		return errors.New("ChapterApp 不可用")
	}

	if args == nil || args.ChapterID == "" {
		return errors.New("参数不合法")
	}

	lgr := a.lgr.With(zap.String("method", "Update"), zap.String("curr_user_id", currUserID), zap.String("chapter_id", args.ChapterID))

	cx = injectLgr(cx, lgr)

	lgr.Info("[logChapterAppImpl.Update] CALL")

	return a.app.Update(cx, currUserID, args)
}

func (a *logChapterAppImpl) Remove(
	cx context.Context,
	currUserID string,
	chapterID string,
) error {
	if a == nil || a.app == nil || a.lgr == nil {
		return errors.New("ChapterApp 不可用")
	}

	if chapterID == "" {
		return errors.New("章节 ID 不能为空")
	}

	lgr := a.lgr.With(zap.String("method", "Remove"), zap.String("curr_user_id", currUserID), zap.String("chapter_id", chapterID))

	cx = injectLgr(cx, lgr)

	lgr.Info("[logChapterAppImpl.Remove] CALL")

	return a.app.Remove(cx, currUserID, chapterID)
}
