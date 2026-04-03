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
	pageRepo       repo.PageRepo
	ossClient      oss.Client
}

func NewPageApp(
	pageSvc service.PageService,
	assignmentRepo repo.AssignmentRepo,
	chapterRepo repo.ChapterRepo,
	pageRepo repo.PageRepo,
	ossClient oss.Client,
) PageApp {
	// 校验构造函数依赖
	if pageSvc == nil ||
		assignmentRepo == nil ||
		chapterRepo == nil ||
		pageRepo == nil ||
		ossClient == nil {
		zap.L().Panic(
			"NewPageApp: 依赖项不能为空",
			zap.Bool("pageSvc_nil", pageSvc == nil),
			zap.Bool("assignmentRepo_nil", assignmentRepo == nil),
			zap.Bool("chapterRepo_nil", chapterRepo == nil),
			zap.Bool("pageRepo_nil", pageRepo == nil),
			zap.Bool("ossClient_nil", ossClient == nil),
		)
	}

	// 返回真实业务实现
	return &pageAppImpl{
		pageSvc:        pageSvc,
		assignmentRepo: assignmentRepo,
		chapterRepo:    chapterRepo,
		pageRepo:       pageRepo,
		ossClient:      ossClient,
	}
}

func (a *pageAppImpl) Reserve(
	cx context.Context,
	currUserID string,
	args *val.ReserveChapterPagesArgs,
) (*val.ReserveChapterPagesRes, error) {
	// 获取上下文中的日志记录器
	lgr := retrieveLgr(cx)

	// 批量构建页面创建载荷
	creations := make([]*model.PageCreation, args.PageCount)

	creationResults := make([]val.PageCreationResult, args.PageCount)

	for i := 0; i < args.PageCount; i++ {
		// 通过领域服务生成 OSS Key
		ossKey := a.pageSvc.GenOSSKey(i)

		// 通过领域服务构造页面创建载荷（含权限校验）
		creation, err := a.pageSvc.NewCreation(
			a.assignmentRepo,
			currUserID,
			args.ChapterID,
			i,
			ossKey,
			currUserID,
		)
		if err != nil {
			// 记录权限校验失败
			lgr.Warn(
				"预留页面失败：权限不足或参数不合法",
				zap.String("curr_user_id", currUserID),
				zap.String("chapter_id", args.ChapterID),
				zap.Error(err),
			)

			// 返回领域服务返回的错误
			return nil, err
		}

		// 生成上传预签名 URL
		putURL, err := a.ossClient.GeneratePutPresignedURL(ossKey)
		if err != nil {
			// 记录生成预签名 URL 失败
			lgr.Error(
				"预留页面失败：生成预签名 URL 失败",
				zap.String("oss_key", ossKey),
				zap.Error(err),
			)

			// 返回客户端可展示的错误
			return nil, errors.New("生成上传地址失败")
		}

		creations[i] = creation

		creationResults[i] = val.PageCreationResult{
			PageID: creation.ID,
			PutURL: putURL,
		}
	}

	// 批量持久化页面
	if err := a.pageRepo.CreateBatch(creations); err != nil {
		// 记录创建失败
		lgr.Error(
			"预留页面失败：批量创建失败",
			zap.String("chapter_id", args.ChapterID),
			zap.Error(err),
		)

		// 返回客户端可展示的错误
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

	// 鉴权：检查当前用户是否为该章节的分配人员
	_, err := a.assignmentRepo.Get(model.AssignmentQueryOpt{
		ChapterID: &args.ChapterID,
		UserID:    &currUserID,
	})
	if err != nil {
		// 记录权限校验失败
		lgr.Warn(
			"获取页面列表失败：权限不足",
			zap.String("curr_user_id", currUserID),
			zap.String("chapter_id", args.ChapterID),
		)

		// 返回客户端可展示的错误
		return nil, errors.New("权限不足")
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
		result[i] = assemblePageInfo(&page, a.ossClient)
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

	// 删除 OSS 资源
	if targetPage.OSSKey != "" {
		if err := a.ossClient.Delete(targetPage.OSSKey); err != nil {
			// 记录删除 OSS 资源失败（继续删除数据库记录）
			lgr.Warn(
				"删除页面 OSS 资源失败",
				zap.String("page_id", pageID),
				zap.String("oss_key", targetPage.OSSKey),
				zap.Error(err),
			)
		}
	}

	// 执行删除
	if err := a.pageRepo.Delete(pageID); err != nil {
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
	ossClient oss.Client,
) *val.PageInfo {
	// 默认图片地址为空字符串
	imageURL := ""

	// 仅在存在 OSS Key 时生成访问地址
	if info.OSSKey != "" {
		url, err := ossClient.GenerateGetPresignedURL(info.OSSKey)
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
		userInfo, _ := assembleUserInfo(info.Creator, ossClient)
		result.Creator = userInfo
	}

	return result
}

// logPageAppImpl 是 PageApp 的日志包装实现
type logPageAppImpl struct {
	lgr *zap.Logger
	app PageApp
}

func NewLogPageApp(
	lgr *zap.Logger,
	app PageApp,
) PageApp {
	if lgr == nil || app == nil {
		zap.L().Panic(
			"NewLogPageApp: 依赖项不能为空",
			zap.Bool("lgr_nil", lgr == nil),
			zap.Bool("app_nil", app == nil),
		)
	}

	return &logPageAppImpl{lgr: lgr, app: app}
}

func (a *logPageAppImpl) Reserve(
	cx context.Context,
	currUserID string,
	args *val.ReserveChapterPagesArgs,
) (*val.ReserveChapterPagesRes, error) {
	if a == nil || a.app == nil || a.lgr == nil {
		return nil, errors.New("PageApp 不可用")
	}

	if args == nil || args.ChapterID == "" || args.PageCount <= 0 {
		return nil, errors.New("参数不合法")
	}

	lgr := a.lgr.With(zap.String("method", "Reserve"), zap.String("curr_user_id", currUserID), zap.String("chapter_id", args.ChapterID))

	cx = injectLgr(cx, lgr)

	lgr.Info("[logPageAppImpl.Reserve] CALL")

	return a.app.Reserve(cx, currUserID, args)
}

func (a *logPageAppImpl) List(
	cx context.Context,
	currUserID string,
	args *val.ListChapterPageArgs,
) ([]*val.PageInfo, error) {
	if a == nil || a.app == nil || a.lgr == nil {
		return nil, errors.New("PageApp 不可用")
	}

	if args == nil || args.ChapterID == "" {
		return nil, errors.New("参数不合法")
	}

	lgr := a.lgr.With(zap.String("method", "List"), zap.String("curr_user_id", currUserID), zap.String("chapter_id", args.ChapterID))

	cx = injectLgr(cx, lgr)

	lgr.Info("[logPageAppImpl.List] CALL")

	return a.app.List(cx, currUserID, args)
}

func (a *logPageAppImpl) Update(
	cx context.Context,
	currUserID string,
	args *val.UpdatePageArgs,
) error {
	if a == nil || a.app == nil || a.lgr == nil {
		return errors.New("PageApp 不可用")
	}

	if args == nil || args.ID == "" {
		return errors.New("参数不合法")
	}

	lgr := a.lgr.With(zap.String("method", "Update"), zap.String("curr_user_id", currUserID), zap.String("page_id", args.ID))

	cx = injectLgr(cx, lgr)

	lgr.Info("[logPageAppImpl.Update] CALL")

	return a.app.Update(cx, currUserID, args)
}

func (a *logPageAppImpl) Remove(
	cx context.Context,
	currUserID string,
	pageID string,
) error {
	if a == nil || a.app == nil || a.lgr == nil {
		return errors.New("PageApp 不可用")
	}

	if pageID == "" {
		return errors.New("页面 ID 不能为空")
	}

	lgr := a.lgr.With(zap.String("method", "Remove"), zap.String("curr_user_id", currUserID), zap.String("page_id", pageID))

	cx = injectLgr(cx, lgr)

	lgr.Info("[logPageAppImpl.Remove] CALL")

	return a.app.Remove(cx, currUserID, pageID)
}
