package application

import (
	"errors"

	"labelplus-next-web-be/internal/application/adapter"
	"labelplus-next-web-be/internal/domain/external"
	"labelplus-next-web-be/internal/domain/model"
	"labelplus-next-web-be/internal/domain/repository"
	"labelplus-next-web-be/internal/domain/service"
	repository_infra "labelplus-next-web-be/internal/infrastructure/repository"
	"labelplus-next-web-be/internal/infrastructure/repository/query_option"
	"labelplus-next-web-be/internal/util"
	"labelplus-next-web-be/internal/value"

	"labelplus-next-web-be/internal/application/assembler"

	"go.uber.org/zap"
)

type PageApplication interface {
	ReserveChapterPages(
		scope util.TraceScope,
		currentUserID string,
		args value.ReserveChapterPagesArgs,
	) (value.ReserveChapterPagesResult, error)
	ListChapterPages(
		scope util.TraceScope,
		currentUserID string,
		args value.ListChapterPageArgs,
	) ([]value.PageInfo, error)
	UpdatePage(
		scope util.TraceScope,
		currentUserID string,
		args value.UpdatePageArgs,
	) error
	DeletePages(
		scope util.TraceScope,
		currentUserID string,
		chapterID string,
	) error
}

type pageApplication struct {
	ossClient            external.OSSClient
	comicRepository      repository.ComicRepository
	memberRepository     repository.MemberRepository
	pageRepository       repository.PageRepository
	assignmentRepository repository.AssignmentRepository
}

func NewPageApplication(
	ossClient external.OSSClient,
	comicRepository repository.ComicRepository,
	memberRepository repository.MemberRepository,
	pageRepository repository.PageRepository,
	assignmentRepository repository.AssignmentRepository,
) PageApplication {
	if ossClient == nil ||
		comicRepository == nil ||
		memberRepository == nil ||
		assignmentRepository == nil ||
		pageRepository == nil {
		zap.L().Panic(
			"NewPageApplication: 依赖项不能为空",
			zap.Bool("ossClient_nil", ossClient == nil),
			zap.Bool("comicRepository_nil", comicRepository == nil),
			zap.Bool("memberRepository_nil", memberRepository == nil),
			zap.Bool("pageRepository_nil", pageRepository == nil),
			zap.Bool("assignmentRepository_nil", assignmentRepository == nil),
		)
	}

	return &pageApplication{
		ossClient:            ossClient,
		comicRepository:      comicRepository,
		memberRepository:     memberRepository,
		pageRepository:       pageRepository,
		assignmentRepository: assignmentRepository,
	}
}

func (pa *pageApplication) ReserveChapterPages(
	scope util.TraceScope,
	currentUserID string,
	args value.ReserveChapterPagesArgs,
) (value.ReserveChapterPagesResult, error) {
	const fn = "PageApplication.CreateChapterPages"

	if err := args.Validate(); err != nil {
		scope.Logger().Error(fn+": 参数验证失败", zap.Error(err))
		return value.ReserveChapterPagesResult{}, err
	}

	scope.
		WithFields(
			zap.Any("args", args),
		).
		Logger().
		Debug(fn + ": 被调用")

	// 鉴权：判定当前用户是否有权限在当前漫画下创建章节
	if !model.PermPageCreate().Check(
		currentUserID,
		args.ChapterID,
		adapter.HandleLoadAssignmentInfo(pa.assignmentRepository),
	) {
		scope.Logger().Warn(fn + ": 权限检查失败")
		return value.ReserveChapterPagesResult{}, errors.New("权限不足")
	}

	// 先在数据库中创建页面记录，随后再根据 ID 创建预签名 PUT URL，最后返回页面 ID 和预签名 URL 列表
	transactionExecutor := pa.pageRepository.BeginTransaction()
	if transactionExecutor.Error != nil {
		scope.Logger().Error(fn+": 开始事务失败", zap.Error(transactionExecutor.Error))
		return value.ReserveChapterPagesResult{}, errors.New("创建漫画页失败")
	}

	var transactionErr error

	defer func() {
		if transactionErr != nil {
			if rollbackErr := transactionExecutor.Rollback().Error; rollbackErr != nil {
				scope.Logger().Error(fn+": 回滚事务失败", zap.Error(rollbackErr))
			}
		}
	}()

	transactionErr = pa.pageRepository.LockByChapterID(transactionExecutor, args.ChapterID)
	if transactionErr != nil {
		scope.Logger().Error(fn+": 锁定章节记录失败", zap.Error(transactionErr))
		return value.ReserveChapterPagesResult{}, errors.New("创建漫画页失败")
	}

	// 生成所有的主键和 OSS Key
	pageCreations := make([]model.PageCreation, args.PageCount)

	for i := 0; i < args.PageCount; i++ {
		pageCreations[i] = model.NewPageCreation(
			util.GenerateUUID(),
			args.ChapterID,
			i,
			service.GeneratePageOSSKey(i),
			currentUserID,
		)
	}

	transactionErr = pa.pageRepository.CreateBatch(transactionExecutor, pageCreations)
	if transactionErr != nil {
		scope.Logger().Error(fn+": 创建页面记录失败", zap.Error(transactionErr))
		return value.ReserveChapterPagesResult{}, errors.New("创建漫画页失败")
	}

	// 生成预签名 URL
	creationResults := make([]value.PageCreationResult, args.PageCount)
	for i, pageCreation := range pageCreations {
		presignedURL, err := pa.ossClient.GeneratePutPresignedURL(pageCreation.OSSKey)
		if err != nil {
			// 防止变量掩蔽导致 defer 中无法正确回滚事务
			transactionErr = err
			scope.Logger().Error(fn+": 生成预签名 URL 失败", zap.Error(transactionErr))

			return value.ReserveChapterPagesResult{}, errors.New("创建漫画页失败")
		}

			creationResults[i] = value.PageCreationResult{PageID: pageCreations[i].ID, PutURL: presignedURL}
	}

	if commitErr := transactionExecutor.Commit().Error; commitErr != nil {
		scope.Logger().Error(fn+": 提交事务失败", zap.Error(commitErr))
		return value.ReserveChapterPagesResult{}, errors.New("创建漫画页失败")
	}

	result := value.ReserveChapterPagesResult{Creations: creationResults}

	return result, nil
}

func (pa *pageApplication) ListChapterPages(
	scope util.TraceScope,
	currentUserID string,
	args value.ListChapterPageArgs,
) ([]value.PageInfo, error) {
	const fn = "PageApplication.ListChapterPages"

	if err := args.Validate(); err != nil {
		scope.Logger().Warn(fn+": 参数验证失败", zap.Error(err))
		return nil, errors.New("参数错误: " + err.Error())
	}

	scope.
		WithFields(
			zap.String("current_user_id", currentUserID),
			zap.Any("args", args),
		).
		Logger().
		Debug(fn + ": 被调用")

	if !model.PermPageList().Check(
		currentUserID,
		args.ChapterID,
		func(chapterID string) (model.ChapterInfo, error) {
			return pa.pageRepository.GetChapterByID(nil, chapterID)
		},
		adapter.HandleLoadComicInfo(pa.comicRepository),
		adapter.HandleLoadMemberInfo(pa.memberRepository),
	) {
		scope.Logger().Warn(fn + ": 权限检查失败")
		return nil, errors.New("权限不足")
	}

	includeSpec := service.ResolvePageListIncludeSpec(args.Includes)

	queryOptions := []repository.QueryOption{
		query_option.PageQuery().FilterByChapterID(args.ChapterID),
		query_option.PageQuery().OrderByIndexAsc(),
		query_option.Paginate(args.Offset, args.Limit),
	}

	if includeSpec.NeedCreator {
		queryOptions = append(queryOptions, query_option.PageQuery().IncludeCreatorInfo())
	}

	pageInfos, err := pa.pageRepository.List(
		nil,
		queryOptions...,
	)
	if err != nil {
		scope.Logger().Error(fn+": 获取页面列表失败", zap.Error(err))
		return nil, errors.New("获取页面列表失败")
	}

	result := make([]value.PageInfo, len(pageInfos))
	for i, pageInfo := range pageInfos {
		result[i] = assembler.AssemblePageInfo(pageInfo, pa.ossClient.GenerateGetPresignedURL)
	}

	return result, nil
}

func (pa *pageApplication) UpdatePage(
	scope util.TraceScope,
	currentUserID string,
	args value.UpdatePageArgs,
) error {
	const fn = "PageApplication.UpdatePage"

	if err := args.Validate(); err != nil {
		scope.Logger().Warn(fn+": 参数验证失败", zap.Error(err))
		return errors.New("参数错误: " + err.Error())
	}

	scope.
		WithFields(
			zap.String("current_user_id", currentUserID),
			zap.Any("args", args),
		).
		Logger().
		Debug(fn + ": 被调用")

	if !model.PermPageUpdate().Check(
		currentUserID,
		args.ID,
		func(pageID string) (model.PageInfo, error) {
			return pa.pageRepository.Get(
				nil,
				query_option.FilterByID(repository_infra.PageTable, pageID),
			)
		},
		func(chapterID string) (model.ChapterInfo, error) {
			return pa.pageRepository.GetChapterByID(nil, chapterID)
		},
		adapter.HandleLoadAssignmentInfo(pa.assignmentRepository),
	) {
		scope.Logger().Warn(fn + ": 权限检查失败")
		return errors.New("权限不足")
	}

	pageInfo, err := pa.pageRepository.Get(
		nil,
		query_option.FilterByID(repository_infra.PageTable, args.ID),
	)
	if err != nil {
		scope.Logger().Error(fn+": 获取目标页面信息失败", zap.Error(err))
		return errors.New("无法获取页面信息")
	}

	pageUpdate := model.NewPageUpdate(
		args.ID,
		pageInfo.Index,
		pageInfo.OSSKey,
		args.IsUploaded,
		pageInfo.TotalUnitCount,
		pageInfo.TranslatedUnitCount,
		pageInfo.ProofreadUnitCount,
	)

	if err := pa.pageRepository.Update(nil, pageUpdate); err != nil {
		scope.Logger().Error(fn+": 更新页面失败", zap.Error(err))
		return errors.New("更新页面失败")
	}

	return nil
}

func (pa *pageApplication) DeletePages(
	scope util.TraceScope,
	currentUserID string,
	chapterID string,
) error {
	const fn = "PageApplication.DeletePages"

	if chapterID == "" {
		scope.Logger().Warn(fn + ": chapterID 不能为空")
		return errors.New("参数错误: chapterID 不能为空")
	}

	scope.
		WithFields(
			zap.String("current_user_id", currentUserID),
			zap.String("chapter_id", chapterID),
		).
		Logger().
		Debug(fn + ": 被调用")

	// 删除是 chapter 维度的批量操作，这里沿用 chapter 维度的权限校验。
	if !model.PermPageCreate().Check(
		currentUserID,
		chapterID,
		adapter.HandleLoadAssignmentInfo(pa.assignmentRepository),
	) {
		scope.Logger().Warn(fn + ": 权限检查失败")
		return errors.New("权限不足")
	}

	transactionExecutor := pa.pageRepository.BeginTransaction()
	if transactionExecutor.Error != nil {
		scope.Logger().Error(fn+": 开始事务失败", zap.Error(transactionExecutor.Error))
		return errors.New("删除漫画页失败")
	}

	var transactionErr error

	defer func() {
		if transactionErr != nil {
			if rollbackErr := transactionExecutor.Rollback().Error; rollbackErr != nil {
				scope.Logger().Error(fn+": 回滚事务失败", zap.Error(rollbackErr))
			}
		}
	}()

	transactionErr = pa.pageRepository.LockByChapterID(transactionExecutor, chapterID)
	if transactionErr != nil {
		scope.Logger().Error(fn+": 锁定章节记录失败", zap.Error(transactionErr))
		return errors.New("删除漫画页失败")
	}

	pages, transactionErr := pa.pageRepository.List(
		transactionExecutor,
		query_option.PageQuery().FilterByChapterID(chapterID),
	)
	if transactionErr != nil {
		scope.Logger().Error(fn+": 查询页面信息失败", zap.Error(transactionErr))
		return errors.New("删除漫画页失败")
	}

	pageIDs := make([]string, 0, len(pages))
	ossKeys := make([]string, 0, len(pages))
	for _, page := range pages {
		pageIDs = append(pageIDs, page.ID)
		if page.OSSKey != "" {
			ossKeys = append(ossKeys, page.OSSKey)
		}
	}

	transactionErr = pa.pageRepository.DeleteBatch(transactionExecutor, pageIDs)
	if transactionErr != nil {
		scope.Logger().Error(fn+": 删除页面记录失败", zap.Error(transactionErr))
		return errors.New("删除漫画页失败")
	}

	if err := pa.ossClient.DeleteBatch(ossKeys); err != nil {
		transactionErr = err
		scope.Logger().Error(fn+": 删除 OSS 页面资源失败", zap.Error(transactionErr))
		return errors.New("删除漫画页失败")
	}

	if commitErr := transactionExecutor.Commit().Error; commitErr != nil {
		scope.Logger().Error(fn+": 提交事务失败", zap.Error(commitErr))
		return errors.New("删除漫画页失败")
	}

	return nil
}
