package application

import (
	"errors"

	"labelplus-next-web-be/internal/application/adapter"
	"labelplus-next-web-be/internal/application/assembler"
	"labelplus-next-web-be/internal/domain/external"
	"labelplus-next-web-be/internal/domain/model"
	"labelplus-next-web-be/internal/domain/repository"
	"labelplus-next-web-be/internal/domain/service"
	repository_infra "labelplus-next-web-be/internal/infrastructure/repository"
	"labelplus-next-web-be/internal/infrastructure/repository/query_option"
	"labelplus-next-web-be/internal/util"
	"labelplus-next-web-be/internal/value"

	"go.uber.org/zap"
)

type ComicApplication interface {
	ListComics(
		scope util.TraceScope,
		currentUserID string,
		args value.ListComicArgs,
	) ([]value.ComicInfo, error)
	CreateComic(
		scope util.TraceScope,
		currentUserID string,
		args value.CreateComicArgs,
	) (value.CreateComicResult, error)
	UpdateComic(
		scope util.TraceScope,
		currentUserID string,
		args value.UpdateComicArgs,
	) error
	DeleteComic(
		scope util.TraceScope,
		currentUserID string,
		comicID string,
	) error
	GetComicCover(
		scope util.TraceScope,
		currentUserID string,
		comicID string,
	) (string, error)
}

type comicApplication struct {
	ossClient         external.OSSClient
	userRepository    repository.UserRepository
	memberRepository  repository.MemberRepository
	worksetRepository repository.WorksetRepository
	comicRepository   repository.ComicRepository
}

func NewComicApplication(
	ossClient external.OSSClient,
	userRepository repository.UserRepository,
	memberRepository repository.MemberRepository,
	worksetRepository repository.WorksetRepository,
	comicRepository repository.ComicRepository,
) ComicApplication {
	if ossClient == nil ||
		userRepository == nil ||
		memberRepository == nil ||
		worksetRepository == nil ||
		comicRepository == nil {
		zap.L().Panic(
			"NewComicApplication: 依赖项不能为空",
			zap.Bool("ossClient_nil", ossClient == nil),
			zap.Bool("userRepository_nil", userRepository == nil),
			zap.Bool("memberRepository_nil", memberRepository == nil),
			zap.Bool("worksetRepository_nil", worksetRepository == nil),
			zap.Bool("comicRepository_nil", comicRepository == nil),
		)
	}

	return &comicApplication{
		ossClient:         ossClient,
		userRepository:    userRepository,
		memberRepository:  memberRepository,
		worksetRepository: worksetRepository,
		comicRepository:   comicRepository,
	}
}

func (ca *comicApplication) ListComics(
	scope util.TraceScope,
	currentUserID string,
	args value.ListComicArgs,
) ([]value.ComicInfo, error) {
	const fn = "ComicApplication.ListComics"

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

	// 通过工作集获取汉化组 ID，用于鉴权
	targetWorkset, err := ca.worksetRepository.Get(
		nil,
		query_option.FilterByID(repository_infra.WorksetTable, args.WorksetID),
	)
	if err != nil {
		scope.Logger().Error(fn+": 获取工作集信息失败", zap.Error(err))
		return nil, errors.New("无法获取工作集信息")
	}

	// 鉴权：检查当前用户在工作集所属汉化组是否有权限查看漫画列表
	if !model.PermComicList().Check(
		currentUserID,
		targetWorkset.TeamID,
		adapter.HandleLoadMemberInfo(ca.memberRepository),
	) {
		scope.Logger().Warn(fn + ": 权限检查失败")
		return nil, errors.New("权限不足")
	}

	includeSpec := service.ResolveComicListIncludeSpec(args.Includes)

	queryOptions := []repository.QueryOption{
		query_option.ComicQuery().FilterByWorksetID(args.WorksetID),
		query_option.ComicQuery().OrderByLastActiveAtDesc(),
		query_option.Paginate(args.Offset, args.Limit),
	}
	switch {
	case includeSpec.NeedWorkset && includeSpec.NeedCreator:
		queryOptions = append(queryOptions, query_option.ComicQuery().IncludeWorksetAndCreatorInfo())
	case includeSpec.NeedWorkset:
		queryOptions = append(queryOptions, query_option.ComicQuery().IncludeWorksetInfo())
	case includeSpec.NeedCreator:
		queryOptions = append(queryOptions, query_option.ComicQuery().IncludeCreatorInfo())
	}

	// 获取漫画列表
	comicList, err := ca.comicRepository.List(nil, queryOptions...)
	if err != nil {
		scope.Logger().Error(fn+": 获取漫画列表失败", zap.Error(err))
		return nil, errors.New("无法获取漫画列表")
	}

	result := make([]value.ComicInfo, len(comicList))
	for i, comic := range comicList {
		result[i] = assembler.AssembleComicInfo(comic, func(string) (string, error) {
			return "", nil
		})
	}

	return result, nil
}

func (ca *comicApplication) CreateComic(
	scope util.TraceScope,
	currentUserID string,
	args value.CreateComicArgs,
) (value.CreateComicResult, error) {
	const fn = "ComicApplication.CreateComic"

	if err := args.Validate(); err != nil {
		scope.Logger().Warn(fn+": 参数验证失败", zap.Error(err))
		return value.CreateComicResult{}, errors.New("参数错误: " + err.Error())
	}

	scope.
		WithFields(
			zap.String("current_user_id", currentUserID),
			zap.Any("args", args),
		).
		Logger().
		Debug(fn + ": 被调用")

	// 通过工作集获取汉化组 ID，用于鉴权
	targetWorkset, err := ca.worksetRepository.Get(
		nil,
		query_option.FilterByID(repository_infra.WorksetTable, args.WorksetID),
	)
	if err != nil {
		scope.Logger().Error(fn+": 获取工作集信息失败", zap.Error(err))
		return value.CreateComicResult{}, errors.New("无法获取工作集信息")
	}

	// 鉴权：检查当前用户在工作集所属汉化组是否有创建漫画权限
	if !model.PermComicCreate().Check(
		currentUserID,
		targetWorkset.TeamID,
		adapter.HandleLoadMemberInfo(ca.memberRepository),
	) {
		scope.Logger().Warn(fn + ": 权限检查失败")
		return value.CreateComicResult{}, errors.New("权限不足")
	}

	transactionExecutor := ca.comicRepository.BeginTransaction()
	if transactionExecutor.Error != nil {
		scope.Logger().Error(fn+": 开启事务失败", zap.Error(transactionExecutor.Error))
		return value.CreateComicResult{}, errors.New("创建漫画失败")
	}

	var transactionErr error

	defer func() {
		if transactionErr != nil {
			if rollbackErr := transactionExecutor.Rollback().Error; rollbackErr != nil {
				scope.Logger().Error(fn+": 回滚事务失败", zap.Error(rollbackErr))
			}
		}
	}()

	// FIXME：其实都没必要 lock，因为数据库有 UNIQUE (workset_id, index) 约束了，并发创建漫画时必然有一个会失败
	transactionErr = ca.comicRepository.LockByWorksetID(transactionExecutor, args.WorksetID)
	if transactionErr != nil {
		scope.Logger().Error(fn+": 锁定漫画记录失败", zap.Error(transactionErr))
		return value.CreateComicResult{}, errors.New("创建漫画失败")
	}

	comicCount, transactionErr := ca.comicRepository.Count(
		transactionExecutor,
		query_option.ComicQuery().FilterByWorksetID(args.WorksetID),
	)
	if transactionErr != nil {
		scope.Logger().Error(fn+": 统计漫画数量失败", zap.Error(transactionErr))
		return value.CreateComicResult{}, errors.New("创建漫画失败")
	}

	// 创建漫画
	comicCreation := model.NewComicCreation(
		args.WorksetID,
		// 因为是 0-based index，所以新漫画的 index 就是当前漫画数量
		int(comicCount),
		args.Title,
		args.Author,
		args.Description,
		currentUserID,
	)

	comicID, transactionErr := ca.comicRepository.Create(transactionExecutor, *comicCreation)
	if transactionErr != nil {
		scope.Logger().Error(fn+": 创建漫画失败", zap.Error(transactionErr))
		return value.CreateComicResult{}, errors.New("创建漫画失败")
	}

	if commitErr := transactionExecutor.Commit().Error; commitErr != nil {
		scope.Logger().Error(fn+": 提交事务失败", zap.Error(commitErr))
		return value.CreateComicResult{}, errors.New("创建漫画失败")
	}

	result := value.CreateComicResult{ID: comicID}

	return result, nil
}

func (ca *comicApplication) UpdateComic(
	scope util.TraceScope,
	currentUserID string,
	args value.UpdateComicArgs,
) error {
	const fn = "ComicApplication.UpdateComic"

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

	_, err := ca.comicRepository.Get(
		nil,
		query_option.FilterByID(repository_infra.ComicTable, args.ID),
	)
	if err != nil {
		scope.Logger().Error(fn+": 获取目标漫画信息失败", zap.Error(err))
		return errors.New("无法获取漫画信息")
	}

	// 鉴权：检查当前用户在目标漫画所属汉化组是否有更新权限
	if !model.PermComicUpdate().Check(
		currentUserID,
		args.ID,
		adapter.HandleLoadComicInfo(ca.comicRepository),
		adapter.HandleLoadWorksetInfo(ca.worksetRepository),
		adapter.HandleLoadMemberInfo(ca.memberRepository),
	) {
		scope.Logger().Warn(fn + ": 权限检查失败")
		return errors.New("权限不足")
	}

	// 构建更新对象
	comicUpdate := model.NewComicUpdate(
		args.ID,
		args.Title,
		args.Author,
		args.Description,
	)

	if err := ca.comicRepository.Update(nil, comicUpdate); err != nil {
		scope.Logger().Error(fn+": 更新漫画失败", zap.Error(err))
		return errors.New("更新漫画失败")
	}

	return nil
}

func (ca *comicApplication) DeleteComic(
	scope util.TraceScope,
	currentUserID string,
	comicID string,
) error {
	const fn = "ComicApplication.DeleteComic"

	if comicID == "" {
		scope.Logger().Warn(fn + ": comicID 为空")
		return errors.New(ErrInternalError)
	}

	scope.
		WithFields(
			zap.String("current_user_id", currentUserID),
			zap.String("comic_id", comicID),
		).
		Logger().
		Debug(fn + ": 被调用")

	_, err := ca.comicRepository.Get(
		nil,
		query_option.FilterByID(repository_infra.ComicTable, comicID),
	)
	if err != nil {
		scope.Logger().Error(fn+": 获取目标漫画信息失败", zap.Error(err))
		return errors.New("无法获取漫画信息")
	}

	// 鉴权：检查当前用户在目标漫画所属汉化组是否有删除权限
	if !model.PermComicDelete().Check(
		currentUserID,
		comicID,
		adapter.HandleLoadComicInfo(ca.comicRepository),
		adapter.HandleLoadWorksetInfo(ca.worksetRepository),
		adapter.HandleLoadMemberInfo(ca.memberRepository),
	) {
		scope.Logger().Warn(fn + ": 权限检查失败")
		return errors.New("权限不足")
	}

	// 删除漫画
	if err := ca.comicRepository.Delete(nil, comicID); err != nil {
		scope.Logger().Error(fn+": 删除漫画失败", zap.Error(err))
		return errors.New("删除漫画失败")
	}

	return nil
}

func (ca *comicApplication) GetComicCover(
	scope util.TraceScope,
	currentUserID string,
	comicID string,
) (string, error) {
	const fn = "ComicApplication.GetComicCover"

	if comicID == "" {
		scope.Logger().Warn(fn + ": comicID 为空")
		return "", errors.New("漫画 ID 不能为空")
	}

	scope.
		WithFields(
			zap.String("current_user_id", currentUserID),
			zap.String("comic_id", comicID),
		).
		Logger().
		Debug(fn + ": 被调用")

	comicInfo, err := ca.comicRepository.Get(
		nil,
		query_option.FilterByID(repository_infra.ComicTable, comicID),
	)
	if err != nil {
		scope.Logger().Error(fn+": 获取目标漫画信息失败", zap.Error(err))
		return "", errors.New("无法获取漫画信息")
	}

	worksetInfo, err := ca.worksetRepository.Get(
		nil,
		query_option.FilterByID(repository_infra.WorksetTable, comicInfo.WorksetID),
	)
	if err != nil {
		scope.Logger().Error(fn+": 获取目标工作集信息失败", zap.Error(err))
		return "", errors.New("无法获取工作集信息")
	}

	if !model.PermComicList().Check(
		currentUserID,
		worksetInfo.TeamID,
		adapter.HandleLoadMemberInfo(ca.memberRepository),
	) {
		scope.Logger().Warn(fn + ": 权限检查失败")
		return "", errors.New("权限不足")
	}

	coverOSSKey, err := ca.comicRepository.GetLatestChapterFirstPageOSSKey(nil, comicID)
	if err != nil {
		scope.Logger().Error(fn+": 查询封面 OSS Key 失败", zap.Error(err))
		return "", errors.New("无法获取漫画封面")
	}

	if coverOSSKey == nil {
		return "", nil
	}

	coverURL, err := ca.ossClient.GenerateGetPresignedURL(*coverOSSKey)
	if err != nil {
		scope.Logger().Error(fn+": 生成封面访问链接失败", zap.Error(err))
		return "", errors.New("无法获取漫画封面")
	}

	return coverURL, nil
}
