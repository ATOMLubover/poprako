package app

import (
	"context"
	"errors"

	"poprako-s/internal/app/val"
	"poprako-s/internal/domain/ext/oss"
	"poprako-s/internal/domain/model"
	"poprako-s/internal/domain/repo"
	"poprako-s/internal/domain/service"

	"go.uber.org/zap"
)

type PageApp interface {
	// Reserve 预留指定数量的页面并返回上传预签名 URL
	Reserve(
		cx context.Context,
		currUserID string,
		args *val.ReserveChapterPagesArgs,
	) (*val.ReserveChapterPagesRes, error)

	// List 获取指定章节的页面列表
	List(
		cx context.Context,
		currUserID string,
		args *val.ListChapterPageArgs,
	) ([]*val.PageInfo, error)

	// Update 更新页面信息
	Update(
		cx context.Context,
		currUserID string,
		args *val.UpdatePageArgs,
	) error

	// Remove 删除页面
	Remove(
		cx context.Context,
		currUserID string,
		pageID string,
	) error
}

type pageAppImpl struct {
	pageSvc service.PageService

	assignmentRepo repo.AssignmentRepo
	chapterRepo    repo.ChapterRepo
	comicRepo      repo.ComicRepo
	worksetRepo    repo.WorksetRepo
	memberRepo     repo.MemberRepo
	pageRepo       repo.PageRepo
	txnMgr         repo.TxnMgr
	msgRepo        repo.OSSMessageRepo
	urlSigner      oss.URLSigner
}

func NewPageApp(
	pageSvc service.PageService,
	assignmentRepo repo.AssignmentRepo,
	chapterRepo repo.ChapterRepo,
	comicRepo repo.ComicRepo,
	worksetRepo repo.WorksetRepo,
	memberRepo repo.MemberRepo,
	pageRepo repo.PageRepo,
	txnMgr repo.TxnMgr,
	msgRepo repo.OSSMessageRepo,
	urlSigner oss.URLSigner,
) PageApp {
	// 校验构造函数依赖
	if pageSvc == nil ||
		assignmentRepo == nil ||
		chapterRepo == nil ||
		comicRepo == nil ||
		worksetRepo == nil ||
		memberRepo == nil ||
		pageRepo == nil ||
		txnMgr == nil ||
		msgRepo == nil ||
		urlSigner == nil {
		zap.L().Panic(
			"NewPageApp: 依赖项不能为空",
			zap.Bool("pageSvc_nil", pageSvc == nil),
			zap.Bool("assignmentRepo_nil", assignmentRepo == nil),
			zap.Bool("chapterRepo_nil", chapterRepo == nil),
			zap.Bool("comicRepo_nil", comicRepo == nil),
			zap.Bool("worksetRepo_nil", worksetRepo == nil),
			zap.Bool("memberRepo_nil", memberRepo == nil),
			zap.Bool("pageRepo_nil", pageRepo == nil),
			zap.Bool("txnMgr_nil", txnMgr == nil),
			zap.Bool("msgRepo_nil", msgRepo == nil),
			zap.Bool("urlSigner_nil", urlSigner == nil),
		)
	}

	// 返回真实业务实现
	return &pageAppImpl{
		pageSvc:        pageSvc,
		assignmentRepo: assignmentRepo,
		chapterRepo:    chapterRepo,
		comicRepo:      comicRepo,
		worksetRepo:    worksetRepo,
		memberRepo:     memberRepo,
		pageRepo:       pageRepo,
		txnMgr:         txnMgr,
		msgRepo:        msgRepo,
		urlSigner:      urlSigner,
	}
}

func (a *pageAppImpl) Reserve(
	cx context.Context,
	currUserID string,
	args *val.ReserveChapterPagesArgs,
) (*val.ReserveChapterPagesRes, error) {
	// 获取上下文中的日志记录器
	lgr := retrieveLgr(cx)

	// 获取章节信息（含 chapterID，用于生成 OSS Key）
	targetChapter, err := a.chapterRepo.GetByID(args.ChapterID)
	if err != nil {
		lgr.Error(
			"预留页面失败：获取章节信息失败",
			zap.String("chapter_id", args.ChapterID),
			zap.Error(err),
		)

		return nil, errors.New("无法获取章节信息")
	}

	// 批量预留页面：每个页面独立事务（page 行写入 + create_pending 消息写入原子完成）
	creationResults := make([]val.PageCreationResult, 0, args.PageCount)

	for i := 0; i < args.PageCount; i++ {
		// 通过领域服务构造页面创建载荷（含权限校验）
		creation, err := a.pageSvc.NewCreation(
			a.assignmentRepo,
			currUserID,
			args.ChapterID,
			i,
			"", // ossKey 将由 ReservePageImage 填入
			currUserID,
		)
		if err != nil {
			lgr.Warn(
				"预留页面失败：权限不足或参数不合法",
				zap.String("curr_user_id", currUserID),
				zap.String("chapter_id", args.ChapterID),
				zap.Error(err),
			)

			return nil, err
		}

		ossKey := a.pageSvc.GenOSSKey(targetChapter.ID, creation.ID)

		putURL, err := a.urlSigner.GeneratePutPresignedURL(ossKey)
		if err != nil {
			lgr.Error(
				"预留页面失败：生成预签名 URL 失败",
				zap.String("chapter_id", args.ChapterID),
				zap.Int("index", i),
				zap.Error(err),
			)

			return nil, errors.New("生成上传地址失败")
		}

		if err := a.txnMgr.RunInTxn(func(txCx context.Context) error {
			pageRepoTxn, txErr := a.pageRepo.FromTxnCx(txCx)
			if txErr != nil {
				return txErr
			}

			creation.OSSKey = ossKey

			if txErr := pageRepoTxn.CreateBatch([]*model.PageCreation{creation}); txErr != nil {
				return txErr
			}

			return upsertCreatePendingMessage(a.msgRepo, txCx, model.OSSResourcePageImage, creation.ID, ossKey)
		}); err != nil {
			lgr.Error(
				"预留页面失败：写入页面或消息失败",
				zap.String("chapter_id", args.ChapterID),
				zap.Int("index", i),
				zap.Error(err),
			)

			return nil, errors.New("预留页面失败")
		}

		creationResults = append(creationResults, val.PageCreationResult{
			PageID: creation.ID,
			PutURL: putURL,
		})
	}

	// 回写章节页面数
	if err := a.txnMgr.RunInTxn(func(cx context.Context) error {
		chapterRepoTxn, err := a.chapterRepo.FromTxnCx(cx)
		if err != nil {
			return err
		}

		return chapterRepoTxn.UpdatePageCount(args.ChapterID, args.PageCount)
	}); err != nil {
		lgr.Error(
			"预留页面失败：回写章节页面数失败",
			zap.String("chapter_id", args.ChapterID),
			zap.Error(err),
		)

		return nil, errors.New("预留页面失败")
	}

	// 返回创建结果
	return &val.ReserveChapterPagesRes{
		Creations: creationResults,
	}, nil
}

func (a *pageAppImpl) List(
	cx context.Context,
	currUserID string,
	args *val.ListChapterPageArgs,
) ([]*val.PageInfo, error) {
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
				"获取页面列表失败：获取章节信息失败",
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
				"获取页面列表失败：获取漫画信息失败",
				zap.String("comic_id", targetChapter.ComicID),
				zap.Error(err),
			)

			return nil, errors.New("无法获取漫画信息")
		}

		tw, err := a.worksetRepo.GetByID(targetComic.WorksetID)
		if err != nil {
			lgr.Error(
				"获取页面列表失败：获取作品集信息失败",
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
					"获取页面列表失败：权限不足",
					zap.String("curr_user_id", currUserID),
					zap.String("chapter_id", args.ChapterID),
				)

				// 返回客户端可展示的错误
				return nil, errors.New("权限不足")
			}
		}
	}

	// 查询页面列表
	pages, err := a.pageRepo.List(model.PageQueryOpt{
		ChapterID: &args.ChapterID,
	})
	if err != nil {
		// 记录查询失败
		lgr.Error(
			"获取页面列表失败",
			zap.String("chapter_id", args.ChapterID),
			zap.Error(err),
		)

		// 返回客户端可展示的错误
		return nil, errors.New("获取页面列表失败")
	}

	// 组装为 app 层值对象列表
	result := make([]*val.PageInfo, len(pages))

	for i, page := range pages {
		result[i] = assemblePageInfo(&page, a.urlSigner)
	}

	// 返回页面列表
	return result, nil
}

func (a *pageAppImpl) Update(
	cx context.Context,
	currUserID string,
	args *val.UpdatePageArgs,
) error {
	// 获取上下文中的日志记录器
	lgr := retrieveLgr(cx)

	// 查询目标页面信息以获取所属章节
	targetPage, err := a.pageRepo.GetByID(args.ID)
	if err != nil {
		// 记录查询失败
		lgr.Error(
			"更新页面失败：获取目标页面信息失败",
			zap.String("page_id", args.ID),
			zap.Error(err),
		)

		// 返回客户端可展示的错误
		return errors.New("无法获取页面信息")
	}

	// 鉴权：检查当前用户是否为该章节的监修或图源
	currAssignment, err := a.assignmentRepo.Get(model.AssignmentQueryOpt{
		ChapterID: &targetPage.ChapterID,
		UserID:    &currUserID,
	})
	if err != nil || !currAssignment.HasAnyRole(model.RoleReviewer, model.RoleRawProvider) {
		// 记录权限校验失败
		lgr.Warn(
			"更新页面失败：权限不足",
			zap.String("curr_user_id", currUserID),
		)

		// 返回客户端可展示的错误
		return errors.New("权限不足")
	}

	// 若此次更新将 is_uploaded 标记为已上传，则同时完成 create_pending 消息
	if args.IsUploaded && !targetPage.IsUploaded {
		if err := a.txnMgr.RunInTxn(func(txCx context.Context) error {
			pageRepoTxn, txErr := a.pageRepo.FromTxnCx(txCx)
			if txErr != nil {
				return txErr
			}

			update := &model.PageUpdate{
				ID:                  args.ID,
				Index:               targetPage.Index,
				OSSKey:              targetPage.OSSKey,
				IsUploaded:          args.IsUploaded,
				TotalUnitCount:      targetPage.TotalUnitCount,
				TranslatedUnitCount: targetPage.TranslatedUnitCount,
				ProofreadUnitCount:  targetPage.ProofreadUnitCount,
			}

			if txErr := pageRepoTxn.Update(update); txErr != nil {
				return txErr
			}

			return completeCreatePendingMessage(a.msgRepo, txCx, model.OSSResourcePageImage, targetPage.ID)
		}); err != nil {
			lgr.Error(
				"更新页面失败",
				zap.String("page_id", args.ID),
				zap.Error(err),
			)

			return errors.New("更新页面失败")
		}

		return nil
	}

	// 构造更新载荷
	update := &model.PageUpdate{
		ID:                  args.ID,
		Index:               targetPage.Index,
		OSSKey:              targetPage.OSSKey,
		IsUploaded:          args.IsUploaded,
		TotalUnitCount:      targetPage.TotalUnitCount,
		TranslatedUnitCount: targetPage.TranslatedUnitCount,
		ProofreadUnitCount:  targetPage.ProofreadUnitCount,
	}

	// 持久化更新
	if err := a.pageRepo.Update(update); err != nil {
		// 记录更新失败
		lgr.Error(
			"更新页面失败",
			zap.String("page_id", args.ID),
			zap.Error(err),
		)

		// 返回客户端可展示的错误
		return errors.New("更新页面失败")
	}

	// 返回更新成功
	return nil
}

func (a *pageAppImpl) Remove(
	cx context.Context,
	currUserID string,
	pageID string,
) error {
	// 获取上下文中的日志记录器
	lgr := retrieveLgr(cx)

	// 查询目标页面信息以获取所属章节
	targetPage, err := a.pageRepo.GetByID(pageID)
	if err != nil {
		// 记录查询失败
		lgr.Error(
			"删除页面失败：获取目标页面信息失败",
			zap.String("page_id", pageID),
			zap.Error(err),
		)

		// 返回客户端可展示的错误
		return errors.New("无法获取页面信息")
	}

	// 鉴权：检查当前用户是否为该章节的监修或图源
	currAssignment, err := a.assignmentRepo.Get(model.AssignmentQueryOpt{
		ChapterID: &targetPage.ChapterID,
		UserID:    &currUserID,
	})
	if err != nil || !currAssignment.HasAnyRole(model.RoleReviewer, model.RoleRawProvider) {
		// 记录权限校验失败
		lgr.Warn(
			"删除页面失败：权限不足",
			zap.String("curr_user_id", currUserID),
		)

		// 返回客户端可展示的错误
		return errors.New("权限不足")
	}

	// 在事务中同时删除页面、入队 OSS 删除并回写章节页面数
	if err := a.txnMgr.RunInTxn(func(cx context.Context) error {
		chapterRepoTxn, err := a.chapterRepo.FromTxnCx(cx)
		if err != nil {
			return err
		}

		pageRepoTxn, err := a.pageRepo.FromTxnCx(cx)
		if err != nil {
			return err
		}

		if err := pageRepoTxn.Delete(pageID); err != nil {
			return err
		}

		if err := enqueueDeleteSingleMessage(a.msgRepo, cx, model.OSSResourcePageImage, pageID, targetPage.OSSKey); err != nil {
			return err
		}

		return chapterRepoTxn.UpdatePageCount(targetPage.ChapterID, -1)
	}); err != nil {
		// 记录删除失败
		lgr.Error(
			"删除页面失败",
			zap.String("page_id", pageID),
			zap.Error(err),
		)

		// 返回客户端可展示的错误
		return errors.New("删除页面失败")
	}

	// 返回删除成功
	return nil
}

// assemblePageInfo 将领域层页面信息转换为 app 层值对象
func assemblePageInfo(
	info *model.PageInfo,
	urlSigner oss.URLSigner,
) *val.PageInfo {
	// 默认图片地址为空字符串
	imageURL := ""

	// 仅在页面已上传且存在 OSS Key 时生成访问地址
	if info.IsUploaded && info.OSSKey != "" {
		url, err := urlSigner.GenerateGetPresignedURL(info.OSSKey)
		if err == nil {
			imageURL = url
		}
	}

	result := &val.PageInfo{
		ID:                  info.ID,
		ChapterID:           info.ChapterID,
		Index:               info.Index,
		ImageURL:            imageURL,
		IsUploaded:          info.IsUploaded,
		CreatorID:           info.CreatorID,
		TotalUnitCount:      info.TotalUnitCount,
		TranslatedUnitCount: info.TranslatedUnitCount,
		ProofreadUnitCount:  info.ProofreadUnitCount,
		CreatedAt:           info.CreatedAt.UnixMilli(),
		UpdatedAt:           info.UpdatedAt.UnixMilli(),
	}

	// 若包含创建者信息则一并组装
	if info.Creator != nil {
		userInfo, _ := assembleUserInfo(info.Creator, urlSigner)
		result.Creator = userInfo
	}

	return result
}

// logPageAppImpl 是 PageApp 的日志包装实现
type logPageAppImpl struct {
	app PageApp
}

func NewLogPageApp(
	app PageApp,
) PageApp {
	if app == nil {
		zap.L().Panic(
			"NewLogPageApp: 依赖项不能为空",
			zap.Bool("app_nil", app == nil),
		)
	}

	return &logPageAppImpl{app: app}
}

func (a *logPageAppImpl) Reserve(
	cx context.Context,
	currUserID string,
	args *val.ReserveChapterPagesArgs,
) (*val.ReserveChapterPagesRes, error) {
	if a == nil || a.app == nil {
		return nil, errors.New("PageApp 不可用")
	}

	if args == nil || args.ChapterID == "" || args.PageCount <= 0 {
		return nil, errors.New("参数不合法")
	}

	lgr := retrieveLgr(cx).With(zap.String("method", "Reserve"), zap.String("curr_user_id", currUserID), zap.String("chapter_id", args.ChapterID))

	cx = injectLgr(cx, lgr)

	lgr.Info("[logPageAppImpl.Reserve] CALL")

	return a.app.Reserve(cx, currUserID, args)
}

func (a *logPageAppImpl) List(
	cx context.Context,
	currUserID string,
	args *val.ListChapterPageArgs,
) ([]*val.PageInfo, error) {
	if a == nil || a.app == nil {
		return nil, errors.New("PageApp 不可用")
	}

	if args == nil || args.ChapterID == "" {
		return nil, errors.New("参数不合法")
	}

	lgr := retrieveLgr(cx).With(zap.String("method", "List"), zap.String("curr_user_id", currUserID), zap.String("chapter_id", args.ChapterID))

	cx = injectLgr(cx, lgr)

	lgr.Info("[logPageAppImpl.List] CALL")

	return a.app.List(cx, currUserID, args)
}

func (a *logPageAppImpl) Update(
	cx context.Context,
	currUserID string,
	args *val.UpdatePageArgs,
) error {
	if a == nil || a.app == nil {
		return errors.New("PageApp 不可用")
	}

	if args == nil || args.ID == "" {
		return errors.New("参数不合法")
	}

	lgr := retrieveLgr(cx).With(zap.String("method", "Update"), zap.String("curr_user_id", currUserID), zap.String("page_id", args.ID))

	cx = injectLgr(cx, lgr)

	lgr.Info("[logPageAppImpl.Update] CALL")

	return a.app.Update(cx, currUserID, args)
}

func (a *logPageAppImpl) Remove(
	cx context.Context,
	currUserID string,
	pageID string,
) error {
	if a == nil || a.app == nil {
		return errors.New("PageApp 不可用")
	}

	if pageID == "" {
		return errors.New("页面 ID 不能为空")
	}

	lgr := retrieveLgr(cx).With(zap.String("method", "Remove"), zap.String("curr_user_id", currUserID), zap.String("page_id", pageID))

	cx = injectLgr(cx, lgr)

	lgr.Info("[logPageAppImpl.Remove] CALL")

	return a.app.Remove(cx, currUserID, pageID)
}
