package app

import (
	"context"
	"errors"
	"path/filepath"

	"poprako-s/internal/app/val"
	"poprako-s/internal/cfg"
	"poprako-s/internal/domain/event"
	"poprako-s/internal/domain/ext/oss"
	"poprako-s/internal/domain/model"
	"poprako-s/internal/domain/repo"
	"poprako-s/internal/domain/service"

	"go.uber.org/zap"
)

type UserApp interface {
	// ParseToken 解析访问令牌并返回用户 ID
	ParseToken(
		cx context.Context,
		tokenStr string,
	) (string, error)

	// Login 处理用户登录逻辑，验证用户凭证并返回登录结果
	Login(
		cx context.Context,
		args *val.LoginUserArgs,
	) (*val.LoginUserRes, error)

	// Reg 处理用户注册逻辑，创建新用户并返回注册结果
	// 注意，该处的注册仅包含“用户从未注册过”的情况
	// 如果需要已有用户加入新的汉化组，必须使用 TeamApp.Join 函数
	Reg(
		cx context.Context,
		args *val.RegUserArgs,
	) (*val.RegUserRes, error)

	// GetInfo 获取指定用户的信息，返回用户的基本资料
	GetInfo(
		cx context.Context,
		userID string,
	) (*val.UserInfo, error)

	// GetMyInfo 获取当前登录用户的信息，返回用户的基本资料
	GetMyInfo(
		cx context.Context,
		currUserID string,
	) (*val.UserInfo, error)

	// UpdateMyInfo 更新当前登录用户的信息，允许修改用户的基本资料
	UpdateMyInfo(
		cx context.Context,
		currUserID string,
		args *val.UpdateUserArgs,
	) error

	// GetMyStats 获取当前登录用户的统计信息
	GetMyStats(
		cx context.Context,
		currUserID string,
	) (*val.UserStatsInfo, error)

	// ReserveMyAvatar 预留头像上传所需的预签名 URL
	ReserveMyAvatar(
		cx context.Context,
		currUserID string,
		args *val.ReserveUserAvatarArgs,
	) (*val.ReserveUserAvatarRes, error)

	// ConfirmMyAvatarUploaded 确认当前用户头像已经完成上传
	ConfirmMyAvatarUploaded(
		cx context.Context,
		currUserID string,
	) error

	// Remove 删除目标用户账号
	Remove(
		cx context.Context,
		currUserID string,
		targetUserID string,
	) error
}

type userAppImpl struct {
	userSvc   service.UserService
	memberSvc service.MemberService

	userRepo   repo.UserRepo
	invRepo    repo.MemberInvitationRepo
	memberRepo repo.MemberRepo
	txnMgr     repo.TxnMgr

	eventBus  event.EventBus
	ossClient oss.Client

	authCfg *cfg.AuthCfg
}

func NewUserApp(
	userSvc service.UserService,
	memberSvc service.MemberService,
	userRepo repo.UserRepo,
	invRepo repo.MemberInvitationRepo,
	memberRepo repo.MemberRepo,
	txnMgr repo.TxnMgr,
	eventBus event.EventBus,
	ossClient oss.Client,
	authCfg *cfg.AuthCfg,
) UserApp {
	// 校验构造函数依赖，避免在运行期出现空指针问题
	if userSvc == nil ||
		memberSvc == nil ||
		userRepo == nil ||
		invRepo == nil ||
		memberRepo == nil ||
		txnMgr == nil ||
		eventBus == nil ||
		ossClient == nil ||
		authCfg == nil {
		zap.L().Panic(
			"NewUserApp: 依赖项不能为空",
			zap.Bool("userSvc_nil", userSvc == nil),
			zap.Bool("memberSvc_nil", memberSvc == nil),
			zap.Bool("userRepo_nil", userRepo == nil),
			zap.Bool("invRepo_nil", invRepo == nil),
			zap.Bool("memberRepo_nil", memberRepo == nil),
			zap.Bool("txnMgr_nil", txnMgr == nil),
			zap.Bool("eventBus_nil", eventBus == nil),
			zap.Bool("ossClient_nil", ossClient == nil),
			zap.Bool("authCfg_nil", authCfg == nil),
		)
	}

	// 返回真实业务实现
	return &userAppImpl{
		userSvc:    userSvc,
		memberSvc:  memberSvc,
		userRepo:   userRepo,
		invRepo:    invRepo,
		memberRepo: memberRepo,
		txnMgr:     txnMgr,
		eventBus:   eventBus,
		ossClient:  ossClient,
		authCfg:    authCfg,
	}
}

func (a *userAppImpl) ParseToken(
	cx context.Context,
	tokenStr string,
) (string, error) {
	lgr := retrieveLgr(cx)

	if tokenStr == "" {
		return "", errors.New("访问令牌不能为空")
	}

	claims, err := a.userSvc.ParseToken(tokenStr, []byte(a.authCfg.SecretKey))
	if err != nil {
		lgr.Warn("解析访问令牌失败", zap.Error(err))
		return "", errors.New("无效的访问令牌")
	}

	if claims == nil || claims.UserID == "" {
		return "", errors.New("访问令牌不包含用户信息")
	}

	return claims.UserID, nil
}

func (a *userAppImpl) Login(
	cx context.Context,
	args *val.LoginUserArgs,
) (*val.LoginUserRes, error) {
	// 获取上下文中的日志记录器
	lgr := retrieveLgr(cx)

	// 查询用户登录凭证
	creds, err := a.userRepo.GetCredsByQQ(args.QQ)
	if err != nil {
		// 记录失败原因，但对外保持统一错误语义
		lgr.Warn(
			"登录失败：无法获取用户凭证",
			zap.String("qq", args.QQ),
			zap.Error(err),
		)

		// 返回统一错误，避免暴露用户存在性信息
		return nil, errors.New("登录失败：用户不存在或密码错误")
	}

	// 校验用户密码
	if err := creds.Authenticate(args.Pwd); err != nil {
		// 记录认证失败信息
		lgr.Warn(
			"登录失败：用户凭证验证失败",
			zap.String("qq", args.QQ),
			zap.Error(err),
		)

		// 返回统一错误，避免暴露用户存在性信息
		return nil, errors.New("登录失败：用户不存在或密码错误")
	}

	// 发布登录成功后产生的领域事件
	if err := a.eventBus.Pub(cx, creds.PullEvents()); err != nil {
		// 记录领域事件处理失败
		lgr.Error(
			"登录失败：处理领域事件失败",
			zap.String("qq", args.QQ),
			zap.Error(err),
		)

		// 返回内部错误提示
		return nil, errors.New("登录失败：内部错误")
	}

	// 记录登录成功日志
	lgr.Info(
		"用户登录成功",
		zap.String("qq", args.QQ),
	)

	// 查询登录用户的完整信息
	info, err := a.userRepo.GetByQQ(args.QQ)
	if err != nil {
		// 记录获取用户信息失败
		lgr.Warn(
			"登录失败：无法获取用户信息",
			zap.String("qq", args.QQ),
			zap.Error(err),
		)

		// 返回通用错误提示
		return nil, errors.New("登录失败：无法获取用户信息")
	}

	// 基于用户信息生成访问令牌
	token, err := info.GenToken(a.authCfg.SecretKey, a.authCfg.ExpHrs)
	if err != nil {
		// 记录生成令牌失败
		lgr.Error(
			"登录失败：生成访问令牌失败",
			zap.String("qq", args.QQ),
			zap.Error(err),
		)

		// 返回通用错误提示
		return nil, errors.New("登录失败：生成访问令牌失败")
	}

	// 返回登录结果
	return &val.LoginUserRes{
		UserID:      info.ID,
		AccessToken: token,
	}, nil
}

func (a *userAppImpl) Reg(
	cx context.Context,
	args *val.RegUserArgs,
) (*val.RegUserRes, error) {
	// 获取上下文中的日志记录器
	lgr := retrieveLgr(cx)

	// 在事务外声明创建结果，便于事务提交后继续生成 token

	var createdUser *model.UserInfo

	// 在事务中完成邀请码校验、用户创建和成员创建
	if err := a.txnMgr.RunInTxn(func(cx context.Context) error {
		// 从事务上下文中构造事务版用户仓库
		userRepoTxn, err := a.userRepo.FromTxnCx(cx)
		if err != nil {
			// 直接返回底层错误，由外层统一终止流程
			return err
		}

		// 从事务上下文中构造事务版邀请仓库
		invRepoTxn, err := a.invRepo.FromTxnCx(cx)
		if err != nil {
			// 直接返回底层错误，由外层统一终止流程
			return err
		}

		memberRepoTxn, err := a.memberRepo.FromTxnCx(cx)
		if err != nil {
			// 直接返回底层错误，由外层统一终止流程
			return err
		}

		// 根据 QQ 查找对应的邀请码信息
		inv, err := invRepoTxn.GetByInviteeQQ(args.QQ)
		if err != nil {
			// 记录邀请码无效
			lgr.Warn(
				"注册失败：无效的邀请码",
				zap.String("qq", args.QQ),
				zap.Error(err),
			)

			// 返回面向客户端的错误
			return errors.New("注册失败：无效的邀请码")
		}

		// 构造用户创建领域模型
		userCreation, err := a.userSvc.NewCreation(args.Name, args.QQ, inv)
		if err != nil {
			// 记录领域模型构造失败
			lgr.Error(
				"注册失败：创建用户领域模型失败",
				zap.String("qq", args.QQ),
				zap.Error(err),
			)

			// 返回通用错误提示
			return errors.New("注册失败：用户创建异常")
		}

		// 发布用户创建时产生的领域事件
		if err := a.eventBus.Pub(cx, userCreation.PullEvents()); err != nil {
			// 记录领域事件处理失败
			lgr.Error(
				"注册失败：处理领域事件失败",
				zap.String("qq", args.QQ),
				zap.Error(err),
			)

			// 返回通用错误提示
			return errors.New("注册失败：内部错误")
		}

		// 持久化用户信息
		userInfo, err := userRepoTxn.Create(userCreation)
		if err != nil {
			// 记录创建用户失败
			lgr.Error(
				"注册失败：创建用户失败",
				zap.String("qq", args.QQ),
				zap.Error(err),
			)

			// 返回通用错误提示
			return errors.New("注册失败：创建用户失败")
		}

		// 根据邀请信息构造成员创建领域模型
		memberCreation, err := a.memberSvc.NewCreationFromInvitation(userInfo, inv)
		if err != nil {
			// 记录创建成员领域模型失败
			lgr.Error(
				"注册失败：将用户加入汉化组失败",
				zap.String("qq", args.QQ),
				zap.Error(err),
			)

			// 返回通用错误提示
			return errors.New("注册失败：加入汉化组失败")
		}

		// 持久化成员关系
		if _, err := memberRepoTxn.Create(memberCreation); err != nil {
			// 记录创建成员记录失败
			lgr.Error(
				"注册失败：创建成员记录失败",
				zap.String("qq", args.QQ),
				zap.Error(err),
			)

			// 返回通用错误提示
			return errors.New("注册失败：加入汉化组失败")
		}

		// 将对应的 invitation 标记为已使用
		if err := invRepoTxn.Invalidate(inv.ID); err != nil {
			// 记录标记邀请码失败
			lgr.Error(
				"注册失败：标记邀请码失败",
				zap.String("qq", args.QQ),
				zap.Error(err),
			)

			// 返回通用错误提示
			return errors.New("注册失败：内部错误")
		}

		// 保存已创建的用户信息供事务外继续使用
		createdUser = userInfo

		// 返回 nil 表示事务逻辑执行成功
		return nil
	}); err != nil {
		// 直接返回事务中已经包装好的错误
		return nil, err
	}

	// 基于创建后的用户信息生成访问令牌
	token, err := createdUser.GenToken(a.authCfg.SecretKey, a.authCfg.ExpHrs)
	if err != nil {
		// 记录生成令牌失败
		lgr.Error(
			"注册失败：生成访问令牌失败",
			zap.String("qq", args.QQ),
			zap.Error(err),
		)

		// 返回通用错误提示
		return nil, errors.New("注册失败：生成访问令牌失败")
	}

	// 返回注册结果
	return &val.RegUserRes{
		UserID:      createdUser.ID,
		AccessToken: token,
	}, nil
}

func (a *userAppImpl) GetInfo(
	cx context.Context,
	userID string,
) (*val.UserInfo, error) {
	// 获取上下文中的日志记录器
	lgr := retrieveLgr(cx)

	// 根据用户 ID 查询用户信息
	info, err := a.userRepo.GetByID(userID)
	if err != nil {
		// 记录查询失败
		lgr.Warn(
			"获取用户信息失败",
			zap.String("user_id", userID),
			zap.Error(err),
		)

		// 返回客户端可展示的错误
		return nil, errors.New("获取用户信息失败")
	}

	// 将领域模型组装为 app 层值对象
	infoVal, err := assembleUserInfo(info, a.ossClient)
	if err != nil {
		// 记录头像访问地址生成失败
		lgr.Error(
			"获取用户信息失败：生成头像访问链接失败",
			zap.String("user_id", userID),
			zap.Error(err),
		)

		// 返回客户端可展示的错误
		return nil, errors.New("获取用户信息失败")
	}

	// 返回组装后的用户信息
	return infoVal, nil
}

func (a *userAppImpl) GetMyInfo(
	cx context.Context,
	currUserID string,
) (*val.UserInfo, error) {
	// 直接复用查询用户信息的主逻辑
	return a.GetInfo(cx, currUserID)
}

func (a *userAppImpl) UpdateMyInfo(
	cx context.Context,
	currUserID string,
	args *val.UpdateUserArgs,
) error {
	// 获取上下文中的日志记录器
	lgr := retrieveLgr(cx)

	// 组装用户更新载荷
	update := &model.UserUpdate{
		ID:   currUserID,
		Name: args.Name,
		QQ:   args.QQ,
	}

	// 持久化用户信息更新
	if err := a.userRepo.Update(update); err != nil {
		// 记录更新失败
		lgr.Error(
			"更新用户信息失败",
			zap.String("user_id", currUserID),
			zap.Error(err),
		)

		// 返回客户端可展示的错误
		return errors.New("更新用户信息失败")
	}

	// 返回更新成功
	return nil
}

func (a *userAppImpl) GetMyStats(
	cx context.Context,
	currUserID string,
) (*val.UserStatsInfo, error) {
	// 获取上下文中的日志记录器
	lgr := retrieveLgr(cx)

	// 读取或懒创建当前用户统计信息
	stats, err := a.userRepo.GetOrCreateStats(currUserID)
	if err != nil {
		// 记录统计信息获取失败
		lgr.Error(
			"获取用户统计信息失败",
			zap.String("user_id", currUserID),
			zap.Error(err),
		)

		// 返回客户端可展示的错误
		return nil, errors.New("获取用户统计信息失败")
	}

	// 组装并返回统计信息
	return assembleUserStatsInfo(stats), nil
}

func (a *userAppImpl) ReserveMyAvatar(
	cx context.Context,
	currUserID string,
	args *val.ReserveUserAvatarArgs,
) (*val.ReserveUserAvatarRes, error) {
	// 获取上下文中的日志记录器
	lgr := retrieveLgr(cx)

	// 提取文件扩展名，与 OSS Key 拼接以便 OSS 正确识别 Content-Type
	ext := filepath.Ext(args.FileName)

	// 基于用户 ID 生成头像对象 Key，并附加扩展名
	avatarOSSKey := a.userSvc.GenAvatarOSSKey(currUserID) + ext

	// 为客户端生成预签名上传链接
	putURL, err := a.ossClient.GeneratePutPresignedURL(avatarOSSKey)
	if err != nil {
		// 记录上传链接生成失败
		lgr.Error(
			"预留头像失败：生成上传链接失败",
			zap.String("user_id", currUserID),
			zap.Error(err),
		)

		// 返回客户端可展示的错误
		return nil, errors.New("预留头像失败")
	}

	// 在用户记录中预填充头像对象 Key
	if err := a.userRepo.PreFillAvatarOSSKey(currUserID, avatarOSSKey); err != nil {
		// 记录预写头像 Key 失败
		lgr.Error(
			"预留头像失败：写入头像 OSS Key 失败",
			zap.String("user_id", currUserID),
			zap.Error(err),
		)

		// 返回客户端可展示的错误
		return nil, errors.New("预留头像失败")
	}

	// 返回预留结果
	return &val.ReserveUserAvatarRes{PutURL: putURL}, nil
}

func (a *userAppImpl) ConfirmMyAvatarUploaded(
	cx context.Context,
	currUserID string,
) error {
	// 获取上下文中的日志记录器
	lgr := retrieveLgr(cx)

	// 将头像状态标记为已上传
	if err := a.userRepo.ConfirmAvatarUploaded(currUserID); err != nil {
		// 记录确认失败
		lgr.Error(
			"确认头像上传失败",
			zap.String("user_id", currUserID),
			zap.Error(err),
		)

		// 返回客户端可展示的错误
		return errors.New("确认头像上传失败")
	}

	// 返回确认成功
	return nil
}

func (a *userAppImpl) Remove(
	cx context.Context,
	currUserID string,
	targetUserID string,
) error {
	// 获取上下文中的日志记录器
	lgr := retrieveLgr(cx)

	// 拦截用户自删场景
	if currUserID == targetUserID {
		// 返回客户端可展示的错误
		return errors.New("无法删除自己")
	}

	targetUser, err := a.userRepo.GetByID(targetUserID)
	if err != nil {
		lgr.Error(
			"删除用户失败：获取目标用户信息失败",
			zap.String("curr_user_id", currUserID),
			zap.String("target_user_id", targetUserID),
			zap.Error(err),
		)

		return errors.New("删除用户失败")
	}

	if err := newOSSDeleteExecutor(a.ossClient).deleteOne(targetUser.AvatarKey); err != nil {
		lgr.Error(
			"删除用户失败：删除头像 OSS 资源失败",
			zap.String("curr_user_id", currUserID),
			zap.String("target_user_id", targetUserID),
			zap.String("avatar_oss_key", targetUser.AvatarKey),
			zap.Error(err),
		)

		return errors.New("删除用户失败")
	}

	// 执行目标用户删除
	if err := a.userRepo.Remove(targetUserID); err != nil {
		// 记录删除失败
		lgr.Error(
			"删除用户失败",
			zap.String("curr_user_id", currUserID),
			zap.String("target_user_id", targetUserID),
			zap.Error(err),
		)

		// 返回客户端可展示的错误
		return errors.New("删除用户失败")
	}

	// 返回删除成功
	return nil
}

// assembleUserInfo 将领域层用户信息转换为 app 层值对象
func assembleUserInfo(
	info *model.UserInfo,
	ossClient oss.Client,
) (*val.UserInfo, error) {
	// 默认头像地址为空字符串
	avatarURL := ""

	// 仅在头像已上传且存在对象 Key 时生成访问地址
	if info.IsAvatarUploaded && info.AvatarKey != "" {
		// 生成头像下载链接
		url, err := ossClient.GenerateGetPresignedURL(info.AvatarKey)
		if err != nil {
			// 将错误返回给调用方处理
			return nil, err
		}

		// 保存生成后的头像地址
		avatarURL = url
	}

	// 返回组装后的 app 层值对象
	return &val.UserInfo{
		ID:               info.ID,
		Name:             info.Name,
		QQ:               info.QQ,
		AvatarURL:        avatarURL,
		IsAvatarUploaded: info.IsAvatarUploaded,
		IsSuperAdmin:     info.IsSuperAdmin,
		LastLoginAt:      info.LastLoginAt.UnixMilli(),
		CreatedAt:        info.CreatedAt.UnixMilli(),
		UpdatedAt:        info.UpdatedAt.UnixMilli(),
	}, nil
}

// assembleUserStatsInfo 将领域层统计信息转换为 app 层值对象
func assembleUserStatsInfo(stats *model.UserStats) *val.UserStatsInfo {
	// 返回组装后的统计信息
	return &val.UserStatsInfo{
		UserID:                  stats.UserID,
		TotalAssignmentCount:    stats.TotalAssignmentCount,
		ActiveAssignmentCount:   stats.ActiveAssignmentCount,
		FinishedAssignmentCount: stats.FinishedAssignmentCount,
	}
}

// logUserAppImpl 是 UserApp 的日志包装实现
type logUserAppImpl struct {
	// app 是被包装的真实 UserApp 实现
	app UserApp
}

func NewLogUserApp(
	app UserApp,
) UserApp {
	// 校验构造函数依赖，避免运行期空指针
	if app == nil {
		zap.L().Panic(
			"NewLogUserApp: 依赖项不能为空",
			zap.Bool("app_nil", app == nil),
		)
	}

	// 返回日志包装实现
	return &logUserAppImpl{
		app: app,
	}
}

func (a *logUserAppImpl) Login(
	cx context.Context,
	args *val.LoginUserArgs,
) (*val.LoginUserRes, error) {
	// 校验包装器实例本身是否合法
	if a == nil || a.app == nil {
		// 返回客户端可展示的错误
		return nil, errors.New("UserApp 不可用")
	}

	// 校验登录参数
	if err := a.validateLoginArgs(args); err != nil {
		// 直接返回参数错误
		return nil, err
	}

	// 为当前调用构造带方法名的日志记录器
	lgr := retrieveLgr(cx).With(
		zap.String("method", "Login"),
	)

	// 将日志记录器注入上下文
	cx = injectLgr(cx, lgr)

	// 记录方法调用日志
	lgr.Info("[logUserAppImpl.Login] CALL")

	// 转发调用到真实实现
	return a.app.Login(cx, args)
}

func (a *logUserAppImpl) ParseToken(
	cx context.Context,
	tokenStr string,
) (string, error) {
	if a == nil || a.app == nil {
		return "", errors.New("UserApp 不可用")
	}

	if tokenStr == "" {
		return "", errors.New("访问令牌不能为空")
	}

	lgr := retrieveLgr(cx).With(zap.String("method", "ParseToken"))
	cx = injectLgr(cx, lgr)
	lgr.Info("[logUserAppImpl.ParseToken] CALL")

	return a.app.ParseToken(cx, tokenStr)
}

func (a *logUserAppImpl) Reg(
	cx context.Context,
	args *val.RegUserArgs,
) (*val.RegUserRes, error) {
	// 校验包装器实例本身是否合法
	if a == nil || a.app == nil {
		// 返回客户端可展示的错误
		return nil, errors.New("UserApp 不可用")
	}

	// 校验注册参数
	if err := a.validateRegArgs(args); err != nil {
		// 直接返回参数错误
		return nil, err
	}

	// 为当前调用构造带方法名的日志记录器
	lgr := retrieveLgr(cx).With(
		zap.String("method", "Reg"),
	)

	// 将日志记录器注入上下文
	cx = injectLgr(cx, lgr)

	// 记录方法调用日志
	lgr.Info("[logUserAppImpl.Reg] CALL")

	// 转发调用到真实实现
	return a.app.Reg(cx, args)
}

func (a *logUserAppImpl) GetInfo(
	cx context.Context,
	userID string,
) (*val.UserInfo, error) {
	// 校验包装器实例本身是否合法
	if a == nil || a.app == nil {
		// 返回客户端可展示的错误
		return nil, errors.New("UserApp 不可用")
	}

	// 校验目标用户 ID
	if userID == "" {
		// 返回客户端可展示的错误
		return nil, errors.New("用户 ID 不能为空")
	}

	// 为当前调用构造带上下文的日志记录器
	lgr := retrieveLgr(cx).With(
		zap.String("method", "GetInfo"),
		zap.String("user_id", userID),
	)

	// 将日志记录器注入上下文
	cx = injectLgr(cx, lgr)

	// 记录方法调用日志
	lgr.Info("[logUserAppImpl.GetInfo] CALL")

	// 转发调用到真实实现
	return a.app.GetInfo(cx, userID)
}

func (a *logUserAppImpl) GetMyInfo(
	cx context.Context,
	currUserID string,
) (*val.UserInfo, error) {
	// 校验包装器实例本身是否合法
	if a == nil || a.app == nil {
		// 返回客户端可展示的错误
		return nil, errors.New("UserApp 不可用")
	}

	// 校验当前用户 ID
	if currUserID == "" {
		// 返回客户端可展示的错误
		return nil, errors.New("用户 ID 不能为空")
	}

	// 为当前调用构造带上下文的日志记录器
	lgr := retrieveLgr(cx).With(
		zap.String("method", "GetMyInfo"),
		zap.String("curr_user_id", currUserID),
	)

	// 将日志记录器注入上下文
	cx = injectLgr(cx, lgr)

	// 记录方法调用日志
	lgr.Info("[logUserAppImpl.GetMyInfo] CALL")

	// 转发调用到真实实现
	return a.app.GetMyInfo(cx, currUserID)
}

func (a *logUserAppImpl) UpdateMyInfo(
	cx context.Context,
	currUserID string,
	args *val.UpdateUserArgs,
) error {
	// 校验包装器实例本身是否合法
	if a == nil || a.app == nil {
		// 返回客户端可展示的错误
		return errors.New("UserApp 不可用")
	}

	// 校验当前用户 ID
	if currUserID == "" {
		// 返回客户端可展示的错误
		return errors.New("用户 ID 不能为空")
	}

	// 校验更新参数
	if err := a.validateUpdateUserArgs(args); err != nil {
		// 直接返回参数错误
		return err
	}

	// 为当前调用构造带上下文的日志记录器
	lgr := retrieveLgr(cx).With(
		zap.String("method", "UpdateMyInfo"),
		zap.String("curr_user_id", currUserID),
	)

	// 将日志记录器注入上下文
	cx = injectLgr(cx, lgr)

	// 记录方法调用日志
	lgr.Info("[logUserAppImpl.UpdateMyInfo] CALL")

	// 转发调用到真实实现
	return a.app.UpdateMyInfo(cx, currUserID, args)
}

func (a *logUserAppImpl) GetMyStats(
	cx context.Context,
	currUserID string,
) (*val.UserStatsInfo, error) {
	// 校验包装器实例本身是否合法
	if a == nil || a.app == nil {
		// 返回客户端可展示的错误
		return nil, errors.New("UserApp 不可用")
	}

	// 校验当前用户 ID
	if currUserID == "" {
		// 返回客户端可展示的错误
		return nil, errors.New("用户 ID 不能为空")
	}

	// 为当前调用构造带上下文的日志记录器
	lgr := retrieveLgr(cx).With(
		zap.String("method", "GetMyStats"),
		zap.String("curr_user_id", currUserID),
	)

	// 将日志记录器注入上下文
	cx = injectLgr(cx, lgr)

	// 记录方法调用日志
	lgr.Info("[logUserAppImpl.GetMyStats] CALL")

	// 转发调用到真实实现
	return a.app.GetMyStats(cx, currUserID)
}

func (a *logUserAppImpl) ReserveMyAvatar(
	cx context.Context,
	currUserID string,
	args *val.ReserveUserAvatarArgs,
) (*val.ReserveUserAvatarRes, error) {
	// 校验包装器实例本身是否合法
	if a == nil || a.app == nil {
		// 返回客户端可展示的错误
		return nil, errors.New("UserApp 不可用")
	}

	// 校验当前用户 ID
	if currUserID == "" {
		// 返回客户端可展示的错误
		return nil, errors.New("用户 ID 不能为空")
	}

	// 校验上传参数
	if args == nil || args.FileName == "" {
		// 返回客户端可展示的错误
		return nil, errors.New("文件名不能为空")
	}

	// 为当前调用构造带上下文的日志记录器
	lgr := retrieveLgr(cx).With(
		zap.String("method", "ReserveMyAvatar"),
		zap.String("curr_user_id", currUserID),
	)

	// 将日志记录器注入上下文
	cx = injectLgr(cx, lgr)

	// 记录方法调用日志
	lgr.Info("[logUserAppImpl.ReserveMyAvatar] CALL")

	// 转发调用到真实实现
	return a.app.ReserveMyAvatar(cx, currUserID, args)
}

func (a *logUserAppImpl) ConfirmMyAvatarUploaded(
	cx context.Context,
	currUserID string,
) error {
	// 校验包装器实例本身是否合法
	if a == nil || a.app == nil {
		// 返回客户端可展示的错误
		return errors.New("UserApp 不可用")
	}

	// 校验当前用户 ID
	if currUserID == "" {
		// 返回客户端可展示的错误
		return errors.New("用户 ID 不能为空")
	}

	// 为当前调用构造带上下文的日志记录器
	lgr := retrieveLgr(cx).With(
		zap.String("method", "ConfirmMyAvatarUploaded"),
		zap.String("curr_user_id", currUserID),
	)

	// 将日志记录器注入上下文
	cx = injectLgr(cx, lgr)

	// 记录方法调用日志
	lgr.Info("[logUserAppImpl.ConfirmMyAvatarUploaded] CALL")

	// 转发调用到真实实现
	return a.app.ConfirmMyAvatarUploaded(cx, currUserID)
}

func (a *logUserAppImpl) Remove(
	cx context.Context,
	currUserID string,
	targetUserID string,
) error {
	// 校验包装器实例本身是否合法
	if a == nil || a.app == nil {
		// 返回客户端可展示的错误
		return errors.New("UserApp 不可用")
	}

	// 校验当前用户 ID
	if currUserID == "" {
		// 返回客户端可展示的错误
		return errors.New("用户 ID 不能为空")
	}

	// 校验目标用户 ID
	if targetUserID == "" {
		// 返回客户端可展示的错误
		return errors.New("目标用户 ID 不能为空")
	}

	// 为当前调用构造带上下文的日志记录器
	lgr := retrieveLgr(cx).With(
		zap.String("method", "Remove"),
		zap.String("curr_user_id", currUserID),
		zap.String("target_user_id", targetUserID),
	)

	// 将日志记录器注入上下文
	cx = injectLgr(cx, lgr)

	// 记录方法调用日志
	lgr.Info("[logUserAppImpl.Remove] CALL")

	// 转发调用到真实实现
	return a.app.Remove(cx, currUserID, targetUserID)
}

func (a *logUserAppImpl) validateLoginArgs(args *val.LoginUserArgs) error {
	// 校验参数对象本身不能为空
	if args == nil {
		// 返回客户端可展示的错误
		return errors.New("LoginUserArgs 不能为空")
	}

	// 校验 QQ 长度范围
	if args.QQ == "" || len(args.QQ) > 14 || len(args.QQ) < 6 {
		// 返回客户端可展示的错误
		return errors.New("QQ 号长度不合法")
	}

	// 校验密码长度范围
	if args.Pwd == "" || len(args.Pwd) < 6 {
		// 返回客户端可展示的错误
		return errors.New("密码长度不能少于 6 位")
	}

	// 返回校验通过
	return nil
}

func (a *logUserAppImpl) validateRegArgs(args *val.RegUserArgs) error {
	// 校验参数对象本身不能为空
	if args == nil {
		// 返回客户端可展示的错误
		return errors.New("RegUserArgs 不能为空")
	}

	// 校验 QQ 长度范围
	if args.QQ == "" || len(args.QQ) > 14 || len(args.QQ) < 6 {
		// 返回客户端可展示的错误
		return errors.New("QQ 号长度不合法")
	}

	// 校验密码长度范围
	if args.Pwd == "" || len(args.Pwd) < 6 {
		// 返回客户端可展示的错误
		return errors.New("密码长度不能少于 6 位")
	}

	// 校验用户名不能为空
	if args.Name == "" {
		// 返回客户端可展示的错误
		return errors.New("用户名不能为空")
	}

	// 校验邀请码不能为空
	if args.InvCode == "" {
		// 返回客户端可展示的错误
		return errors.New("邀请码不能为空")
	}

	// 返回校验通过
	return nil
}

func (a *logUserAppImpl) validateUpdateUserArgs(args *val.UpdateUserArgs) error {
	// 校验参数对象本身不能为空
	if args == nil {
		// 返回客户端可展示的错误
		return errors.New("UpdateUserArgs 不能为空")
	}

	// 校验用户名不能为空
	if args.Name == "" {
		// 返回客户端可展示的错误
		return errors.New("用户名不能为空")
	}

	// 校验 QQ 长度范围
	if args.QQ == "" || len(args.QQ) > 14 || len(args.QQ) < 6 {
		// 返回客户端可展示的错误
		return errors.New("QQ 号长度不合法")
	}

	// 返回校验通过
	return nil
}
