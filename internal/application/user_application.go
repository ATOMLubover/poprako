package application

import (
	"errors"

	"labelplus-next-web-be/internal/config"
	"labelplus-next-web-be/internal/domain/model"
	"labelplus-next-web-be/internal/domain/repository"
	"labelplus-next-web-be/internal/domain/service"
	"labelplus-next-web-be/internal/repository/query_option"
	"labelplus-next-web-be/internal/util"
	"labelplus-next-web-be/internal/value"

	"go.uber.org/zap"
)

type UserApplication interface {
	LoginUser(
		scope util.TraceScope,
		args *value.LoginUserArgs,
	) (*value.LoginUserResult, error)
	RegisterUser(
		scope util.TraceScope,
		args *value.RegisterUserArgs,
	) (*value.RegisterUserResult, error)
	GetUserInfoByID(
		scope util.TraceScope,
		userID string,
	) (*value.UserInfo, error)
	ListUsers(
		scope util.TraceScope,
		args *value.ListUsersArgs,
	) ([]*value.UserInfo, error)
	RemoveUser(
		scope util.TraceScope,
		currentUserID string,
		targetUserID string,
	) error
}

type userApplication struct {
	authConfig *config.AuthConfig

	userService       service.UserService
	permissionService service.PermissionService

	userRepository       repository.UserRepository
	invitationRepository repository.InvitationRepository
}

func NewUserApplication(
	authConfig *config.AuthConfig,
	userService service.UserService,
	permissionService service.PermissionService,
	userRepository repository.UserRepository,
	invitationRepository repository.InvitationRepository,
) UserApplication {
	if authConfig == nil ||
		userService == nil ||
		permissionService == nil ||
		userRepository == nil ||
		invitationRepository == nil {
		zap.L().Panic(
			"NewUserApplication: 依赖项不能为空",
			zap.Bool("authConfig_nil", authConfig == nil),
			zap.Bool("userService_nil", userService == nil),
			zap.Bool("permissionService_nil", permissionService == nil),
			zap.Bool("userRepository_nil", userRepository == nil),
			zap.Bool("invitationRepository_nil", invitationRepository == nil),
		)
	}

	return &userApplication{
		authConfig:           authConfig,
		userService:          userService,
		permissionService:    permissionService,
		userRepository:       userRepository,
		invitationRepository: invitationRepository,
	}
}

func (ua *userApplication) LoginUser(
	scope util.TraceScope,
	args *value.LoginUserArgs,
) (*value.LoginUserResult, error) {
	const fn = "UserApplication.LoginUser"

	if args == nil {
		scope.Logger().Warn(fn + ": args 为空")
		return nil, errors.New(ErrInternalError)
	}

	if err := args.Validate(); err != nil {
		scope.Logger().Warn(fn+": 参数验证失败", zap.Error(err))
		return nil, err
	}

	scope.WithField(
		zap.String("qq", args.QQ),
	)

	scope.Logger().Debug(fn + ": 被调用")

	// 获取用户登录凭根（含密码哈希）
	credentials, err := ua.userRepository.GetCredentialsByQQ(nil, args.QQ)
	if err != nil {
		scope.Logger().Error(fn+": 获取用户凭证失败", zap.Error(err))
		return nil, errors.New("获取用户信息失败")
	}

	// 用户不存在或密码错误，故意不区分以防止用户枚举
	if credentials == nil || !ua.userService.VerifyPassword(args.Password, credentials.PasswordHash) {
		return nil, errors.New("用户不存在或密码错误")
	}

	// 生成 access token
	accessToken, err := ua.userService.GenerateAccessToken(
		credentials.UserID,
		[]byte(ua.authConfig.JWTSecretKey),
		ua.authConfig.ExpirationHours,
	)
	if err != nil {
		scope.Logger().Error(fn+": 生成访问令牌失败", zap.Error(err))
		return nil, errors.New("生成访问令牌失败")
	}

	return &value.LoginUserResult{
		UserID:      credentials.UserID,
		AccessToken: accessToken,
	}, nil
}

func (ua *userApplication) RegisterUser(
	scope util.TraceScope,
	args *value.RegisterUserArgs,
) (*value.RegisterUserResult, error) {
	const fn = "UserApplication.RegisterUser"

	if args == nil {
		scope.Logger().Error(fn + ": args 为空")
		return nil, errors.New(ErrInternalError)
	}

	if err := args.Validate(); err != nil {
		scope.Logger().Error(fn+": 参数验证失败", zap.Error(err))
		return nil, err
	}

	// 检查是否有对应的邀请
	invitationInfo, err := ua.invitationRepository.GetInfoByInviteeQQ(nil, args.QQ)
	if err != nil {
		scope.Logger().Error(fn+": 查询邀请信息失败", zap.Error(err))
		return nil, errors.New("获取邀请信息失败")
	}

	if invitationInfo == nil {
		scope.Logger().Error(fn + ": 没有找到对应的邀请信息")
		return nil, errors.New("没有找到对应的邀请信息，请确保您已被邀请")
	}

	// 检查邀请码是否匹配
	if invitationInfo.InvitationCode != args.InvitationCode {
		scope.Logger().Error(fn + ": 邀请码不匹配")
		return nil, errors.New("邀请码不正确，请检查后重试")
	}

	// 通过检查后，对密码进行 哈希处理
	passwordHash, err := ua.userService.HashPassword(args.Password)
	if err != nil {
		scope.Logger().Error(fn+": 密码哈希失败", zap.Error(err))
		return nil, errors.New("创建用户失败")
	}

	// 通过检查后，创建对应的用户信息
	userRegistration := model.NewUserRegistration(
		args.Name,
		args.QQ,
		passwordHash,
		model.UnmaskRoles(invitationInfo.RoleMask())...,
	)

	userID, err := ua.userRepository.Create(nil, userRegistration)
	if err != nil {
		scope.Logger().Error(fn+": 创建用户信息失败", zap.Error(err))
		return nil, errors.New("创建用户失败")
	}

	// 生成 access token
	accessToken, err := ua.userService.GenerateAccessToken(
		userID,
		[]byte(ua.authConfig.JWTSecretKey),
		ua.authConfig.ExpirationHours,
	)
	if err != nil {
		scope.Logger().Error(fn+": 生成访问令牌失败", zap.Error(err))
		return nil, errors.New("生成访问令牌失败")
	}

	// 构造注册成功的结果值对象
	result := value.NewRegisterUserResult(userID, accessToken)

	return result, nil
}

func (ua *userApplication) GetUserInfoByID(
	scope util.TraceScope,
	userID string,
) (*value.UserInfo, error) {
	const fn = "UserApplication.GetUserInfoByID"

	scope.WithField(
		zap.String("user_id", userID),
	)

	scope.Logger().Debug(fn + ": 被调用")

	userInfo, err := ua.userRepository.GetInfoByID(nil, userID)
	if err != nil {
		scope.Logger().Error(fn+": 获取用户信息失败", zap.Error(err))
		return nil, errors.New("无法获取用户信息")
	}

	if userInfo == nil {
		return nil, errors.New("用户不存在")
	}

	return value.NewUserInfoFromModel(userInfo), nil
}

func (ua *userApplication) ListUsers(
	scope util.TraceScope,
	args *value.ListUsersArgs,
) ([]*value.UserInfo, error) {
	const fn = "UserApplication.ListUsers"

	if args == nil {
		scope.Logger().Warn(fn + ": args 为空")
		return nil, errors.New(ErrInternalError)
	}

	if err := args.Validate(); err != nil {
		scope.Logger().Warn(fn+": 参数验证失败", zap.Error(err))
		return nil, err
	}

	scope.WithField(
		zap.Any("args", args),
	)

	scope.Logger().Debug(fn + ": 被调用")

	// 构建动态查询选项
	opts := make([]repository.QueryOption, 0)

	if args.QQ != "" {
		opts = append(opts, query_option.FilterByQQ(args.QQ))
	}

	if args.FuzzyName != "" {
		opts = append(opts, query_option.FuzzyFilterByName(args.FuzzyName))
	}

	if args.Role != 0 {
		opts = append(opts, query_option.FilterByUserRole(args.Role))
	}

	opts = append(opts, query_option.CreatedAtDesc()...)
	opts = append(opts, query_option.Paginate(args.Offset, args.Limit)...)

	// 查询用户列表
	userList, err := ua.userRepository.List(nil, opts...)
	if err != nil {
		scope.Logger().Error(fn+": 获取用户列表失败", zap.Error(err))
		return nil, errors.New("无法获取用户列表")
	}

	// 转换为应用层的值对象
	result := make([]*value.UserInfo, len(userList))
	for i := range userList {
		result[i] = value.NewUserInfoFromModel(&userList[i])
	}

	return result, nil
}

func (ua *userApplication) RemoveUser(
	scope util.TraceScope,
	currentUserID string,
	targetUserID string,
) error {
	const fn = "UserApplication.RemoveUser"

	scope.WithField(
		zap.String("current_user_id", currentUserID),
		zap.String("target_user_id", targetUserID),
	)

	scope.Logger().Debug(fn + ": 被调用")

	// 鉴权：仅超级管理员可删除用户
	currentUser, err := ua.userRepository.GetInfoByID(nil, currentUserID)
	if err != nil {
		scope.Logger().Error(fn+": 获取当前用户信息失败", zap.Error(err))
		return errors.New("无法获取用户信息")
	}

	if !ua.permissionService.CheckPermission(currentUser, model.PermissionUsersRemove) {
		return errors.New("没有权限删除用户")
	}

	// 防止删除自身
	if currentUserID == targetUserID {
		return errors.New("无法删除自己")
	}

	if err := ua.userRepository.DeleteByID(nil, targetUserID); err != nil {
		scope.Logger().Error(fn+": 删除用户失败", zap.Error(err))
		return errors.New("删除用户失败")
	}

	return nil
}
