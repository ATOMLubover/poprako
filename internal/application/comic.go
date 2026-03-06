package application

import (
	"errors"

	"labelplus-next-web-be/internal/domain/model"
	"labelplus-next-web-be/internal/domain/repository"
	"labelplus-next-web-be/internal/domain/service"
	"labelplus-next-web-be/internal/repository/query_option"
	"labelplus-next-web-be/internal/util"
	"labelplus-next-web-be/internal/value"

	"go.uber.org/zap"
)

type ComicApplication interface {
	ListTeamComics(
		scope util.TraceScope,
		currentUserID string,
		args *value.ListTeamComicArgs,
	) ([]*value.ComicInfo, error)
	CreateComic(
		scope util.TraceScope,
		currentUserID string,
		args *value.CreateComicArgs,
	) (*value.CreateComicResult, error)
	UpdateComic(
		scope util.TraceScope,
		currentUserID string,
		args *value.UpdateComicArgs,
	) error
	DeleteComic(
		scope util.TraceScope,
		currentUserID string,
		comicID string,
	) error
}

type comicApplication struct {
	userRepository   repository.UserRepository
	memberRepository repository.MemberRepository
	comicRepository  repository.ComicRepository
}

func NewComicApplication(
	userRepository repository.UserRepository,
	memberRepository repository.MemberRepository,
	comicRepository repository.ComicRepository,
) ComicApplication {
	if userRepository == nil ||
		memberRepository == nil ||
		comicRepository == nil {
		zap.L().Panic(
			"NewComicApplication: 依赖项不能为空",
			zap.Bool("userRepository_nil", userRepository == nil),
			zap.Bool("memberRepository_nil", memberRepository == nil),
			zap.Bool("comicRepository_nil", comicRepository == nil),
		)
	}

	return &comicApplication{
		userRepository:   userRepository,
		memberRepository: memberRepository,
		comicRepository:  comicRepository,
	}
}

func (ca *comicApplication) ListTeamComics(
	scope util.TraceScope,
	currentUserID string,
	args *value.ListTeamComicArgs,
) ([]*value.ComicInfo, error) {
	const fn = "ComicApplication.ListTeamComics"

	if currentUserID == "" {
		scope.Logger().Warn(fn + ": currentUserID 为空")
		return nil, errors.New(ErrInternalError)
	}

	if args == nil {
		scope.Logger().Warn(fn + ": args 为空")
		return nil, errors.New(ErrInternalError)
	}

	if err := args.Validate(); err != nil {
		scope.Logger().Warn(fn+": 参数验证失败", zap.Error(err))
		return nil, errors.New("参数错误: " + err.Error())
	}

	scope.WithFields(
		zap.String("current_user_id", currentUserID),
		zap.Any("args", args),
	)

	scope.Logger().Debug(fn + ": 被调用")

	// 鉴权：获取当前用户在各汉化组的成员信息
	currentUserMemberships, err := ca.memberRepository.List(
		nil,
		query_option.MemberQuery().FilterByUserID(currentUserID),
	)
	if err != nil {
		scope.Logger().Error(fn+": 获取当前用户的成员信息失败", zap.Error(err))
		return nil, errors.New("无法获取成员信息")
	}

	// 鉴权：检查当前用户在指定汉化组是否有权限查看漫画列表
	if !service.CheckComicPermission(args.TeamID, currentUserMemberships, model.PermissionComicList) {
		return nil, errors.New("没有权限查看漫画列表")
	}

	// 获取漫画列表
	comicList, err := ca.comicRepository.List(
		nil,
		query_option.ComicQuery().FilterByTeamID(args.TeamID),
		query_option.CreatedAtDesc(),
		query_option.Paginate(args.Offset, args.Limit),
	)
	if err != nil {
		scope.Logger().Error(fn+": 获取漫画列表失败", zap.Error(err))
		return nil, errors.New("无法获取漫画列表")
	}

	result := make([]*value.ComicInfo, len(comicList))
	for i, comic := range comicList {
		result[i] = value.NewComicInfoFromModel(comic)
	}

	return result, nil
}

func (ca *comicApplication) CreateComic(
	scope util.TraceScope,
	currentUserID string,
	args *value.CreateComicArgs,
) (*value.CreateComicResult, error) {
	const fn = "ComicApplication.CreateComic"

	if currentUserID == "" {
		scope.Logger().Warn(fn + ": currentUserID 为空")
		return nil, errors.New(ErrInternalError)
	}

	if args == nil {
		scope.Logger().Warn(fn + ": args 为空")
		return nil, errors.New(ErrInternalError)
	}

	if err := args.Validate(); err != nil {
		scope.Logger().Warn(fn+": 参数验证失败", zap.Error(err))
		return nil, errors.New("参数错误: " + err.Error())
	}

	scope.WithFields(
		zap.String("current_user_id", currentUserID),
		zap.Any("args", args),
	)

	scope.Logger().Debug(fn + ": 被调用")

	// 鉴权：获取当前用户在各汉化组的成员信息
	currentUserMemberships, err := ca.memberRepository.List(
		nil,
		query_option.MemberQuery().FilterByUserID(currentUserID),
	)
	if err != nil {
		scope.Logger().Error(fn+": 获取当前用户的成员信息失败", zap.Error(err))
		return nil, errors.New("无法获取成员信息")
	}

	// 鉴权：仅对应汉化组管理员有权限创建漫画
	if !service.CheckComicPermission(args.TeamID, currentUserMemberships, model.PermissionComicCreate) {
		return nil, errors.New("没有权限创建漫画")
	}

	transactionExecutor := ca.comicRepository.BeginTransaction()
	if transactionExecutor.Error != nil {
		scope.Logger().Error(fn+": 开启事务失败", zap.Error(transactionExecutor.Error))
		return nil, errors.New("创建漫画失败")
	}

	var transactionErr error

	defer func() {
		if transactionErr != nil {
			if rollbackErr := transactionExecutor.Rollback().Error; rollbackErr != nil {
				scope.Logger().Error(fn+": 回滚事务失败", zap.Error(rollbackErr))
			}
		}
	}()

	transactionErr = ca.comicRepository.LockByTeamID(transactionExecutor, args.TeamID)
	if transactionErr != nil {
		scope.Logger().Error(fn+": 锁定漫画记录失败", zap.Error(transactionErr))
		return nil, errors.New("创建漫画失败")
	}

	comicCount, transactionErr := ca.comicRepository.CountByTeamID(transactionExecutor, args.TeamID)
	if transactionErr != nil {
		scope.Logger().Error(fn+": 统计漫画数量失败", zap.Error(transactionErr))
		return nil, errors.New("创建漫画失败")
	}

	// 创建漫画
	comicCreation := model.NewComicCreation(
		args.TeamID,
		int(comicCount)+1,
		args.Title,
		args.Author,
		args.Description,
		currentUserID,
	)

	comicID, transactionErr := ca.comicRepository.Create(transactionExecutor, comicCreation)
	if transactionErr != nil {
		scope.Logger().Error(fn+": 创建漫画失败", zap.Error(transactionErr))
		return nil, errors.New("创建漫画失败")
	}

	if commitErr := transactionExecutor.Commit().Error; commitErr != nil {
		scope.Logger().Error(fn+": 提交事务失败", zap.Error(commitErr))
		return nil, errors.New("创建漫画失败")
	}
	transactionErr = nil

	return &value.CreateComicResult{ID: comicID}, nil
}

func (ca *comicApplication) UpdateComic(
	scope util.TraceScope,
	currentUserID string,
	args *value.UpdateComicArgs,
) error {
	const fn = "ComicApplication.UpdateComic"

	if currentUserID == "" {
		scope.Logger().Warn(fn + ": currentUserID 为空")
		return errors.New(ErrInternalError)
	}

	if args == nil {
		scope.Logger().Warn(fn + ": args 为空")
		return errors.New(ErrInternalError)
	}

	if err := args.Validate(); err != nil {
		scope.Logger().Warn(fn+": 参数验证失败", zap.Error(err))
		return errors.New("参数错误: " + err.Error())
	}

	scope.WithFields(
		zap.String("current_user_id", currentUserID),
		zap.Any("args", args),
	)

	scope.Logger().Debug(fn + ": 被调用")

	targetComic, err := ca.comicRepository.GetByID(nil, args.ID)
	if err != nil {
		scope.Logger().Error(fn+": 获取目标漫画信息失败", zap.Error(err))
		return errors.New("无法获取漫画信息")
	}

	// 鉴权：获取当前用户在各汉化组的成员信息
	currentUserMemberships, err := ca.memberRepository.List(
		nil,
		query_option.MemberQuery().FilterByUserID(currentUserID),
	)
	if err != nil {
		scope.Logger().Error(fn+": 获取当前用户的成员信息失败", zap.Error(err))
		return errors.New("无法获取成员信息")
	}

	// 鉴权：仅对应汉化组管理员有权限更新漫画
	if !service.CheckComicPermission(targetComic.TeamID, currentUserMemberships, model.PermissionComicUpdate) {
		return errors.New("没有权限更新漫画")
	}

	// 构建更新对象
	comicUpdate := model.NewComicUpdate(
		args.ID,
		&args.Title,
		&args.Author,
		&args.Description,
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

	if currentUserID == "" {
		scope.Logger().Warn(fn + ": currentUserID 为空")
		return errors.New(ErrInternalError)
	}

	if comicID == "" {
		scope.Logger().Warn(fn + ": comicID 为空")
		return errors.New(ErrInternalError)
	}

	scope.WithFields(
		zap.String("current_user_id", currentUserID),
		zap.String("comic_id", comicID),
	)

	scope.Logger().Debug(fn + ": 被调用")

	targetComic, err := ca.comicRepository.GetByID(nil, comicID)
	if err != nil {
		scope.Logger().Error(fn+": 获取目标漫画信息失败", zap.Error(err))
		return errors.New("无法获取漫画信息")
	}

	// 鉴权：获取当前用户在各汉化组的成员信息
	currentUserMemberships, err := ca.memberRepository.List(
		nil,
		query_option.MemberQuery().FilterByUserID(currentUserID),
	)
	if err != nil {
		scope.Logger().Error(fn+": 获取当前用户的成员信息失败", zap.Error(err))
		return errors.New("无法获取成员信息")
	}

	// 鉴权：仅对应汉化组管理员有权限删除漫画
	if !service.CheckComicPermission(targetComic.TeamID, currentUserMemberships, model.PermissionComicDelete) {
		return errors.New("没有权限删除漫画")
	}

	// 删除漫画
	if err := ca.comicRepository.Delete(nil, comicID); err != nil {
		scope.Logger().Error(fn+": 删除漫画失败", zap.Error(err))
		return errors.New("删除漫画失败")
	}

	return nil
}
