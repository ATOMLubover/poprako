package application

import (
	"errors"

	"labelplus-next-web-be/internal/application/adapter"
	"labelplus-next-web-be/internal/config"
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

type UserApplication interface {
	LoginUser(
		scope util.TraceScope,
		args value.LoginUserArgs,
	) (value.LoginUserResult, error)
	RegisterUser(
		scope util.TraceScope,
		args value.RegisterUserArgs,
	) (value.RegisterUserResult, error)
	GetUser(
		scope util.TraceScope,
		userID string,
	) (value.UserInfo, error)
	ListUsers(
		scope util.TraceScope,
		currentUserID string,
		args value.ListUserArgs,
	) ([]value.UserInfo, error)
	ReserveUserAvatar(
		scope util.TraceScope,
		currentUserID string,
		targetUserID string,
	) (value.ReserveUserAvatarResult, error)
	UpdateUser(
		scope util.TraceScope,
		currentUserID string,
		args value.UpdateUserArgs,
	) error
	RemoveUser(
		scope util.TraceScope,
		currentUserID string,
		targetUserID string,
	) error
}

type userApplication struct {
	authConfig *config.AuthConfig
	ossClient   external.OSSClient

	userRepository       repository.UserRepository
	memberRepository     repository.MemberRepository
	invitationRepository repository.InvitationRepository
}

func NewUserApplication(
	authConfig *config.AuthConfig,
	ossClient external.OSSClient,
	userRepository repository.UserRepository,
	memberRepository repository.MemberRepository,
	invitationRepository repository.InvitationRepository,
) UserApplication {
	if authConfig == nil ||
		ossClient == nil ||
		userRepository == nil ||
		memberRepository == nil ||
		invitationRepository == nil {
		zap.L().Panic(
			"NewUserApplication: 依赖项不能为空",
			zap.Bool("authConfig_nil", authConfig == nil),
			zap.Bool("ossClient_nil", ossClient == nil),
			zap.Bool("userRepository_nil", userRepository == nil),
			zap.Bool("memberRepository_nil", memberRepository == nil),
			zap.Bool("invitationRepository_nil", invitationRepository == nil),
		)
	}

	return &userApplication{
		authConfig:           authConfig,
		ossClient:            ossClient,
		userRepository:       userRepository,
		memberRepository:     memberRepository,
		invitationRepository: invitationRepository,
	}
}

func (ua *userApplication) LoginUser(
	scope util.TraceScope,
	args value.LoginUserArgs,
) (value.LoginUserResult, error) {
	const fn = "UserApplication.LoginUser"

	if err := args.Validate(); err != nil {
		scope.Logger().Warn(fn+": 参数验证失败", zap.Error(err))
		return value.LoginUserResult{}, err
	}

	scope.
		WithFields(
			zap.String("qq", args.QQ),
		).
		Logger().
		Debug(fn + ": 被调用")

	// 获取用户登录凭根（含密码哈希）
	credentials, err := ua.userRepository.GetCredentials(
		nil,
		query_option.UserQuery().FilterByQQ(args.QQ),
	)
	if err != nil && !errors.Is(err, repository_infra.ErrRecordNotFound) {
		scope.Logger().Error(fn+": 获取用户凭证失败", zap.Error(err))
		return value.LoginUserResult{}, errors.New("用户不存在或密码错误")
	}

	// 用户不存在或密码错误，故意不区分以防止用户枚举
	if !service.VerifyPassword(args.Password, credentials.PasswordHash) {
		return value.LoginUserResult{}, errors.New("用户不存在或密码错误")
	}

	// 生成 access token
	accessToken, err := service.GenerateAccessToken(
		credentials.UserID,
		[]byte(ua.authConfig.JWTSecretKey),
		ua.authConfig.ExpirationHours,
	)
	if err != nil {
		scope.Logger().Error(fn+": 生成访问令牌失败", zap.Error(err))
		return value.LoginUserResult{}, errors.New("生成访问令牌失败")
	}

	return value.NewLoginUserResult(credentials.UserID, accessToken), nil
}

func (ua *userApplication) RegisterUser(
	scope util.TraceScope,
	args value.RegisterUserArgs,
) (value.RegisterUserResult, error) {
	const fn = "UserApplication.RegisterUser"

	if err := args.Validate(); err != nil {
		scope.Logger().Error(fn+": 参数验证失败", zap.Error(err))
		return value.RegisterUserResult{}, err
	}

	scope.
		WithFields(
			zap.String("qq", args.QQ),
			zap.String("invitation_code", args.InvitationCode),
		).
		Logger().
		Debug(fn + ": 被调用")

	// 根据 QQ 和邀请码精确查找邀请，同时支持多邀请场景
	invitationInfo, err := ua.invitationRepository.Get(
		nil,
		query_option.InvitationQuery().FilterByInviteeQQ(args.QQ),
		query_option.InvitationQuery().FilterByCode(args.InvitationCode),
		query_option.InvitationQuery().FilterPending(true),
	)
	if err != nil && !errors.Is(err, repository_infra.ErrRecordNotFound) {
		scope.Logger().Error(fn+": 查询邀请信息失败", zap.Error(err))
		return value.RegisterUserResult{}, errors.New("获取邀请信息失败")
	}

	if err != nil {
		scope.Logger().Error(fn + ": 没有找到对应的邀请信息")
		return value.RegisterUserResult{}, errors.New("没有找到对应的邀请信息，请确保邀请码正确")
	}

	// 通过检查后，对密码进行哈希处理
	passwordHash, err := service.HashPassword(args.Password)
	if err != nil {
		scope.Logger().Error(fn+": 密码哈希失败", zap.Error(err))
		return value.RegisterUserResult{}, errors.New("创建用户失败")
	}

	// 通过检查后，在一个事务中，先创建用户信息，再创建成员信息，最后标记邀请信息为已使用
	transactionExecutor := ua.userRepository.BeginTransaction()

	var transactionErr error

	defer func() {
		if transactionErr != nil {
			if rollbackErr := transactionExecutor.Rollback().Error; rollbackErr != nil {
				scope.Logger().Error(
					fn+": 事务回滚失败",
					zap.Error(transactionErr),
					zap.Error(rollbackErr),
				)
			}
		}
	}()

	userRegistration := model.NewUserRegistration(
		args.Name,
		args.QQ,
		passwordHash,
		model.UnmaskRoles(invitationInfo.RoleMask())...,
	)

	userID, transactionErr := ua.userRepository.Create(
		transactionExecutor,
		userRegistration,
	)
	if transactionErr != nil {
		scope.Logger().Error(fn+": 创建用户信息失败", zap.Error(err))
		return value.RegisterUserResult{}, errors.New("创建用户失败")
	}

	// 创建成员信息
	memberCreation := model.NewMemberCreation(
		userID,
		invitationInfo.TeamID,
		model.UnmaskRoles(invitationInfo.RoleMask())...,
	)

	_, transactionErr = ua.memberRepository.Create(
		transactionExecutor,
		memberCreation,
	)
	if transactionErr != nil {
		zap.L().Error(fn+": 创建成员信息失败", zap.Error(transactionErr))
		return value.RegisterUserResult{}, errors.New("创建用户失败")
	}

	// 标记邀请信息为已使用
	transactionErr = ua.invitationRepository.Invalidate(
		transactionExecutor,
		invitationInfo.ID,
	)
	if transactionErr != nil {
		zap.L().Error(fn+": 标记邀请信息已使用失败", zap.Error(transactionErr))
		return value.RegisterUserResult{}, errors.New("创建用户失败")
	}

	if commitErr := transactionExecutor.Commit().Error; commitErr != nil {
		scope.Logger().Error(fn+": 事务提交失败", zap.Error(commitErr))
		return value.RegisterUserResult{}, errors.New("创建用户失败")
	}

	// 生成 access token
	accessToken, err := service.GenerateAccessToken(
		userID,
		[]byte(ua.authConfig.JWTSecretKey),
		ua.authConfig.ExpirationHours,
	)
	if err != nil {
		scope.Logger().Error(fn+": 生成访问令牌失败", zap.Error(err))
		return value.RegisterUserResult{}, errors.New("生成访问令牌失败")
	}

	return value.NewRegisterUserResult(userID, accessToken), nil
}

func (ua *userApplication) GetUser(
	scope util.TraceScope,
	userID string,
) (value.UserInfo, error) {
	const fn = "UserApplication.GetUser"

	scope.
		WithFields(
			zap.String("user_id", userID),
		).
		Logger().
		Debug(fn + ": 被调用")

	// 鉴权：用户只能查看自己的信息，超级管理员可以查看所有用户的信息
	if !model.PermUserView().Check(
		userID,
		userID,
	) {
		scope.Logger().Warn(fn + ":权限检查不通过")
		return value.UserInfo{}, errors.New("没有权限查看用户信息")
	}

	userInfo, err := ua.userRepository.Get(
		nil,
		query_option.FilterByID(repository_infra.UserTable, userID),
	)
	if err != nil {
		scope.Logger().Error(fn+": 获取用户信息失败", zap.Error(err))
		return value.UserInfo{}, errors.New("无法获取用户信息")
	}

	if userInfo.ID == "" {
		return value.UserInfo{}, errors.New("用户不存在")
	}

	return value.NewUserInfoFromModel(userInfo), nil
}

func (ua *userApplication) ListUsers(
	scope util.TraceScope,
	currentUserID string,
	args value.ListUserArgs,
) ([]value.UserInfo, error) {
	const fn = "UserApplication.ListUsers"

	if err := args.Validate(); err != nil {
		scope.Logger().Warn(fn+": 参数验证失败", zap.Error(err))
		return nil, err
	}

	scope.
		WithFields(
			zap.Any("args", args),
		).
		Logger().
		Debug(fn + ": 被调用")

	// 鉴权：仅超级管理员有权限查看用户列表
	if !model.PermUserList().Check(
		currentUserID,
		adapter.HandleLoadUserInfo(ua.userRepository),
	) {
		scope.Logger().Warn(fn + ": 权限检查不通过")
		return nil, errors.New("没有权限查看用户列表")
	}

	// 构建动态查询选项
	opts := make([]repository.QueryOption, 0)

	if args.QQ != "" {
		opts = append(opts, query_option.UserQuery().FilterByQQ(args.QQ))
	}
	if args.FuzzyName != "" {
		opts = append(opts, query_option.UserQuery().FilterByFuzzyName(args.FuzzyName))
	}

	opts = append(opts, query_option.CreatedAtDesc(repository_infra.UserTable))
	opts = append(opts, query_option.Paginate(args.Offset, args.Limit))

	// 查询用户列表
	userList, err := ua.userRepository.List(nil, opts...)
	if err != nil {
		scope.Logger().Error(fn+": 获取用户列表失败", zap.Error(err))
		return nil, errors.New("无法获取用户列表")
	}

	// 转换为应用层的值对象
	result := make([]value.UserInfo, len(userList))
	for i := range userList {
		result[i] = value.NewUserInfoFromModel(userList[i])
	}

	return result, nil
}

func (ua *userApplication) ReserveUserAvatar(
	scope util.TraceScope,
	currentUserID string,
	targetUserID string,
) (value.ReserveUserAvatarResult, error) {
	const fn = "UserApplication.ReserveUserAvatar"

	scope.
		WithFields(
			zap.String("current_user_id", currentUserID),
			zap.String("target_user_id", targetUserID),
		).
		Logger().
		Debug(fn + ": 被调用")

	if targetUserID == "" {
		return value.ReserveUserAvatarResult{}, errors.New("用户 ID 不能为空")
	}

	if !model.PermUserUpdate().Check(
		currentUserID,
		targetUserID,
		adapter.HandleLoadUserInfo(ua.userRepository),
	) {
		scope.Logger().Warn(fn + ": 权限检查不通过")
		return value.ReserveUserAvatarResult{}, errors.New("没有权限预留用户头像")
	}

	avatarOSSKey := service.GenerateUserAvatarOSSKey(targetUserID)

	putURL, err := ua.ossClient.GeneratePutPresignedURL(avatarOSSKey)
	if err != nil {
		scope.Logger().Error(fn+": 生成预签名 URL 失败", zap.Error(err))
		return value.ReserveUserAvatarResult{}, errors.New("预留用户头像失败")
	}

	if err := ua.userRepository.ReserveAvatar(nil, targetUserID, avatarOSSKey); err != nil {
		scope.Logger().Error(fn+": 写入头像 OSS Key 失败", zap.Error(err))
		return value.ReserveUserAvatarResult{}, errors.New("预留用户头像失败")
	}

	return value.NewReserveUserAvatarResult(avatarOSSKey, putURL), nil
}

func (ua *userApplication) UpdateUser(
	scope util.TraceScope,
	currentUserID string,
	args value.UpdateUserArgs,
) error {
	const fn = "UserApplication.UpdateUser"

	if err := args.Validate(); err != nil {
		scope.Logger().Warn(fn+": 参数验证失败", zap.Error(err))
		return err
	}

	scope.
		WithFields(
			zap.String("current_user_id", currentUserID),
			zap.String("target_user_id", args.UserID),
			zap.Any("args", args),
		).
		Logger().
		Debug(fn + ": 被调用")

	if !model.PermUserUpdate().Check(
		currentUserID,
		args.UserID,
		adapter.HandleLoadUserInfo(ua.userRepository),
	) {
		scope.Logger().Warn(fn + ": 权限检查不通过")
		return errors.New("没有权限更新用户")
	}

	hashedPassword, err := service.HashPassword(args.Password)
	if err != nil {
		scope.Logger().Error(fn+": 密码哈希失败", zap.Error(err))
		return errors.New("更新用户失败")
	}

	userUpdate := model.NewUserUpdate(
		args.UserID,
		args.Name,
		args.QQ,
		hashedPassword,
		args.IsAvatarUploaded,
	)

	if err := ua.userRepository.Update(nil, userUpdate); err != nil {
		scope.Logger().Error(fn+": 更新用户失败", zap.Error(err))
		return errors.New("更新用户失败")
	}

	return nil
}

func (ua *userApplication) RemoveUser(
	scope util.TraceScope,
	currentUserID string,
	targetUserID string,
) error {
	const fn = "UserApplication.RemoveUserByID"

	scope.
		WithFields(
			zap.String("current_user_id", currentUserID),
			zap.String("target_user_id", targetUserID),
		).
		Logger().
		Debug(fn + ": 被调用")

	// 鉴权：仅超级管理员可删除用户
	if !model.PermUserRemove().Check(
		currentUserID,
		targetUserID,
		adapter.HandleLoadUserInfo(ua.userRepository),
	) {
		scope.Logger().Warn(fn + ": 权限检查不通过")
		return errors.New("没有权限删除用户")
	}

	// 防止删除自身
	if currentUserID == targetUserID {
		return errors.New("无法删除自己")
	}

	if err := ua.userRepository.Delete(nil, targetUserID); err != nil {
		scope.Logger().Error(fn+": 删除用户失败", zap.Error(err))
		return errors.New("删除用户失败")
	}

	return nil
}
