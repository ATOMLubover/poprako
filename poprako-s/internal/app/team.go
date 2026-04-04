package app

import (
	"context"
	"errors"
	"path/filepath"

	"poprako-s/internal/app/val"
	"poprako-s/internal/domain/ext/oss"
	"poprako-s/internal/domain/model"
	"poprako-s/internal/domain/repo"
	"poprako-s/internal/domain/service"

	"go.uber.org/zap"
)

type TeamApp interface {
	// Create 创建一个新的汉化组，仅超级管理员可操作
	Create(
		cx context.Context,
		currUserID string,
		args *val.CreateTeamArgs,
	) (*val.CreateTeamRes, error)

	// List 获取所有汉化组列表（超级管理员专用）
	List(
		cx context.Context,
		currUserID string,
		args *val.ListTeamArgs,
	) ([]*val.TeamInfo, error)

	// ListMy 获取当前用户所属的汉化组列表
	ListMy(
		cx context.Context,
		currUserID string,
		args *val.ListMyTeamArgs,
	) ([]*val.TeamInfo, error)

	// Update 更新汉化组信息
	Update(
		cx context.Context,
		currUserID string,
		args *val.UpdateTeamArgs,
	) error

	// Remove 删除汉化组
	Remove(
		cx context.Context,
		currUserID string,
		teamID string,
	) error

	// ReserveAvatar 预留汉化组头像上传所需的预签名 URL
	ReserveAvatar(
		cx context.Context,
		currUserID string,
		args *val.ReserveTeamAvatarArgs,
	) (*val.ReserveTeamAvatarRes, error)

	// ConfirmAvatarUploaded 确认汉化组头像已经完成上传
	ConfirmAvatarUploaded(
		cx context.Context,
		currUserID string,
		teamID string,
	) error
}

type teamAppImpl struct {
	teamSvc   service.TeamService
	memberSvc service.MemberService

	userRepo   repo.UserRepo
	teamRepo   repo.TeamRepo
	memberRepo repo.MemberRepo

	ossClient oss.Client
}

func NewTeamApp(
	teamSvc service.TeamService,
	memberSvc service.MemberService,
	userRepo repo.UserRepo,
	teamRepo repo.TeamRepo,
	memberRepo repo.MemberRepo,
	ossClient oss.Client,
) TeamApp {
	// 校验构造函数依赖，避免在运行期出现空指针问题
	if teamSvc == nil ||
		memberSvc == nil ||
		userRepo == nil ||
		teamRepo == nil ||
		memberRepo == nil ||
		ossClient == nil {
		zap.L().Panic(
			"NewTeamApp: 依赖项不能为空",
			zap.Bool("teamSvc_nil", teamSvc == nil),
			zap.Bool("memberSvc_nil", memberSvc == nil),
			zap.Bool("userRepo_nil", userRepo == nil),
			zap.Bool("teamRepo_nil", teamRepo == nil),
			zap.Bool("memberRepo_nil", memberRepo == nil),
			zap.Bool("ossClient_nil", ossClient == nil),
		)
	}

	// 返回真实业务实现
	return &teamAppImpl{
		teamSvc:    teamSvc,
		memberSvc:  memberSvc,
		userRepo:   userRepo,
		teamRepo:   teamRepo,
		memberRepo: memberRepo,
		ossClient:  ossClient,
	}
}

func (a *teamAppImpl) Create(
	cx context.Context,
	currUserID string,
	args *val.CreateTeamArgs,
) (*val.CreateTeamRes, error) {
	// 获取上下文中的日志记录器
	lgr := retrieveLgr(cx)

	// 查询当前用户信息，用于权限校验
	currUser, err := a.userRepo.GetByID(currUserID)
	if err != nil {
		// 记录查询失败
		lgr.Warn(
			"创建汉化组失败：无法获取当前用户信息",
			zap.String("curr_user_id", currUserID),
			zap.Error(err),
		)

		// 返回客户端可展示的错误
		return nil, errors.New("创建汉化组失败：无法获取用户信息")
	}

	// 通过领域服务构造汉化组创建载荷（含权限校验）
	creation, err := a.teamSvc.NewCreation(currUser, args.Name, args.Description)
	if err != nil {
		// 记录权限校验失败
		lgr.Warn(
			"创建汉化组失败：权限不足",
			zap.String("curr_user_id", currUserID),
			zap.Error(err),
		)

		// 返回领域服务返回的错误
		return nil, err
	}

	// 持久化汉化组信息
	teamInfo, err := a.teamRepo.Create(creation)
	if err != nil {
		// 记录创建失败
		lgr.Error(
			"创建汉化组失败：持久化失败",
			zap.String("curr_user_id", currUserID),
			zap.Error(err),
		)

		// 返回客户端可展示的错误
		return nil, errors.New("创建汉化组失败")
	}

	// 返回创建结果
	return &val.CreateTeamRes{ID: teamInfo.ID}, nil
}

func (a *teamAppImpl) List(
	cx context.Context,
	currUserID string,
	args *val.ListTeamArgs,
) ([]*val.TeamInfo, error) {
	// 获取上下文中的日志记录器
	lgr := retrieveLgr(cx)

	// 查询当前用户信息，用于超级管理员权限校验
	currUser, err := a.userRepo.GetByID(currUserID)
	if err != nil {
		// 记录查询失败
		lgr.Warn(
			"获取汉化组列表失败：无法获取当前用户信息",
			zap.String("curr_user_id", currUserID),
			zap.Error(err),
		)

		// 返回客户端可展示的错误
		return nil, errors.New("获取汉化组列表失败")
	}

	// 校验超级管理员权限
	if !currUser.IsSuperAdmin {
		// 记录权限校验失败
		lgr.Warn(
			"获取汉化组列表失败：仅超级管理员可以查看所有汉化组",
			zap.String("curr_user_id", currUserID),
		)

		// 返回客户端可展示的错误
		return nil, errors.New("权限不足")
	}

	// 查询所有汉化组列表
	teams, err := a.teamRepo.List(model.TeamQueryOpt{})
	if err != nil {
		// 记录查询失败
		lgr.Error(
			"获取汉化组列表失败",
			zap.Error(err),
		)

		// 返回客户端可展示的错误
		return nil, errors.New("获取汉化组列表失败")
	}

	// 组装为 app 层值对象列表
	result := make([]*val.TeamInfo, len(teams))

	for i, team := range teams {
		// 组装单个汉化组信息
		info, err := assembleTeamInfo(&team, a.ossClient)
		if err != nil {
			// 记录组装失败
			lgr.Error(
				"获取汉化组列表失败：组装汉化组信息失败",
				zap.String("team_id", team.ID),
				zap.Error(err),
			)

			// 返回客户端可展示的错误
			return nil, errors.New("获取汉化组列表失败")
		}

		result[i] = info
	}

	// 返回汉化组列表
	return result, nil
}

func (a *teamAppImpl) ListMy(
	cx context.Context,
	currUserID string,
	args *val.ListMyTeamArgs,
) ([]*val.TeamInfo, error) {
	// 获取上下文中的日志记录器
	lgr := retrieveLgr(cx)

	// 查询当前用户的所有成员记录
	members, err := a.memberRepo.List(model.MemberQueryOpt{
		UserID: &currUserID,
	})
	if err != nil {
		// 记录查询失败
		lgr.Error(
			"获取我的汉化组列表失败：查询成员记录失败",
			zap.String("curr_user_id", currUserID),
			zap.Error(err),
		)

		// 返回客户端可展示的错误
		return nil, errors.New("获取我的汉化组列表失败")
	}

	// 若用户未加入任何汉化组，直接返回空列表
	if len(members) == 0 {
		return nil, nil
	}

	// 逐个查询成员所属的汉化组信息
	result := make([]*val.TeamInfo, 0, len(members))

	for _, member := range members {
		// 根据汉化组 ID 查询汉化组信息
		team, err := a.teamRepo.GetByID(member.TeamID)
		if err != nil {
			// 记录查询失败
			lgr.Error(
				"获取我的汉化组列表失败：查询汉化组信息失败",
				zap.String("team_id", member.TeamID),
				zap.Error(err),
			)

			// 返回客户端可展示的错误
			return nil, errors.New("获取我的汉化组列表失败")
		}

		// 组装单个汉化组信息
		info, err := assembleTeamInfo(team, a.ossClient)
		if err != nil {
			// 记录组装失败
			lgr.Error(
				"获取我的汉化组列表失败：组装汉化组信息失败",
				zap.String("team_id", team.ID),
				zap.Error(err),
			)

			// 返回客户端可展示的错误
			return nil, errors.New("获取我的汉化组列表失败")
		}

		result = append(result, info)
	}

	// 返回汉化组列表
	return result, nil
}

func (a *teamAppImpl) Update(
	cx context.Context,
	currUserID string,
	args *val.UpdateTeamArgs,
) error {
	// 获取上下文中的日志记录器
	lgr := retrieveLgr(cx)

	// 查询当前用户在目标汉化组中的成员信息，用于鉴权
	currMember, err := a.memberRepo.Get(model.MemberQueryOpt{
		UserID: &currUserID,
		TeamID: &args.ID,
	})
	if err != nil || !currMember.HasAnyRole(model.RoleAdmin) {
		// 记录权限校验失败
		lgr.Warn(
			"更新汉化组失败：权限不足",
			zap.String("curr_user_id", currUserID),
			zap.String("team_id", args.ID),
		)

		// 返回客户端可展示的错误
		return errors.New("权限不足")
	}

	// 组装汉化组更新载荷
	update := &model.TeamUpdate{
		ID:          args.ID,
		Name:        args.Name,
		Description: args.Description,
	}

	// 持久化更新
	if err := a.teamRepo.Update(update); err != nil {
		// 记录更新失败
		lgr.Error(
			"更新汉化组失败",
			zap.String("team_id", args.ID),
			zap.Error(err),
		)

		// 返回客户端可展示的错误
		return errors.New("更新汉化组失败")
	}

	// 返回更新成功
	return nil
}

func (a *teamAppImpl) Remove(
	cx context.Context,
	currUserID string,
	teamID string,
) error {
	// 获取上下文中的日志记录器
	lgr := retrieveLgr(cx)

	// 查询当前用户信息，校验超级管理员权限
	currUser, err := a.userRepo.GetByID(currUserID)
	if err != nil {
		// 记录查询失败
		lgr.Warn(
			"删除汉化组失败：无法获取当前用户信息",
			zap.String("curr_user_id", currUserID),
			zap.Error(err),
		)

		// 返回客户端可展示的错误
		return errors.New("删除汉化组失败")
	}

	// 校验超级管理员权限
	if !currUser.IsSuperAdmin {
		// 记录权限校验失败
		lgr.Warn(
			"删除汉化组失败：权限不足",
			zap.String("curr_user_id", currUserID),
		)

		// 返回客户端可展示的错误
		return errors.New("权限不足")
	}

	// 执行删除
	if err := a.teamRepo.Delete(teamID); err != nil {
		// 记录删除失败
		lgr.Error(
			"删除汉化组失败",
			zap.String("team_id", teamID),
			zap.Error(err),
		)

		// 返回客户端可展示的错误
		return errors.New("删除汉化组失败")
	}

	// 返回删除成功
	return nil
}

func (a *teamAppImpl) ReserveAvatar(
	cx context.Context,
	currUserID string,
	args *val.ReserveTeamAvatarArgs,
) (*val.ReserveTeamAvatarRes, error) {
	// 获取上下文中的日志记录器
	lgr := retrieveLgr(cx)

	// 查询当前用户在目标汉化组中的成员信息，用于鉴权
	currMember, err := a.memberRepo.Get(model.MemberQueryOpt{
		UserID: &currUserID,
		TeamID: &args.TeamID,
	})
	if err != nil || !currMember.HasAnyRole(model.RoleAdmin) {
		// 记录权限校验失败
		lgr.Warn(
			"预留汉化组头像失败：权限不足",
			zap.String("curr_user_id", currUserID),
			zap.String("team_id", args.TeamID),
		)

		// 返回客户端可展示的错误
		return nil, errors.New("权限不足")
	}

	// 提取文件扩展名，与 OSS Key 拼接以便 OSS 正确识别 Content-Type
	ext := filepath.Ext(args.FileName)

	// 基于汉化组 ID 生成头像对象 Key，并附加扩展名
	avatarOSSKey := a.teamSvc.GenAvatarOSSKey(args.TeamID) + ext

	// 为客户端生成预签名上传链接
	putURL, err := a.ossClient.GeneratePutPresignedURL(avatarOSSKey)
	if err != nil {
		// 记录上传链接生成失败
		lgr.Error(
			"预留汉化组头像失败：生成上传链接失败",
			zap.String("team_id", args.TeamID),
			zap.Error(err),
		)

		// 返回客户端可展示的错误
		return nil, errors.New("预留汉化组头像失败")
	}

	// 在数据库中预填充头像对象 Key
	if err := a.teamRepo.PreFillAvatarOSSKey(args.TeamID, avatarOSSKey); err != nil {
		// 记录预写失败
		lgr.Error(
			"预留汉化组头像失败：写入头像 OSS Key 失败",
			zap.String("team_id", args.TeamID),
			zap.Error(err),
		)

		// 返回客户端可展示的错误
		return nil, errors.New("预留汉化组头像失败")
	}

	// 返回预留结果
	return &val.ReserveTeamAvatarRes{PutURL: putURL}, nil
}

func (a *teamAppImpl) ConfirmAvatarUploaded(
	cx context.Context,
	currUserID string,
	teamID string,
) error {
	// 获取上下文中的日志记录器
	lgr := retrieveLgr(cx)

	// 查询当前用户在目标汉化组中的成员信息，用于鉴权
	currMember, err := a.memberRepo.Get(model.MemberQueryOpt{
		UserID: &currUserID,
		TeamID: &teamID,
	})
	if err != nil || !currMember.HasAnyRole(model.RoleAdmin) {
		// 记录权限校验失败
		lgr.Warn(
			"确认汉化组头像上传失败：权限不足",
			zap.String("curr_user_id", currUserID),
			zap.String("team_id", teamID),
		)

		// 返回客户端可展示的错误
		return errors.New("权限不足")
	}

	// 将头像状态标记为已上传
	if err := a.teamRepo.ConfirmAvatarUploaded(teamID); err != nil {
		// 记录确认失败
		lgr.Error(
			"确认汉化组头像上传失败",
			zap.String("team_id", teamID),
			zap.Error(err),
		)

		// 返回客户端可展示的错误
		return errors.New("确认汉化组头像上传失败")
	}

	// 返回确认成功
	return nil
}

// assembleTeamInfo 将领域层汉化组信息转换为 app 层值对象
func assembleTeamInfo(
	info *model.TeamInfo,
	ossClient oss.Client,
) (*val.TeamInfo, error) {
	// 默认头像地址为空字符串
	avatarURL := ""

	// 仅在头像已上传且存在对象 Key 时生成访问地址
	if info.IsAvatarUploaded && info.AvatarOSSKey != "" {
		// 生成头像下载链接
		url, err := ossClient.GenerateGetPresignedURL(info.AvatarOSSKey)
		if err != nil {
			// 将错误返回给调用方处理
			return nil, err
		}

		// 保存生成后的头像地址
		avatarURL = url
	}

	// 返回组装后的 app 层值对象
	return &val.TeamInfo{
		ID:               info.ID,
		Name:             info.Name,
		Description:      info.Description,
		AvatarURL:        avatarURL,
		IsAvatarUploaded: info.IsAvatarUploaded,
		CreatedAt:        info.CreatedAt.UnixMilli(),
		UpdatedAt:        info.UpdatedAt.UnixMilli(),
	}, nil
}

// logTeamAppImpl 是 TeamApp 的日志包装实现
type logTeamAppImpl struct {
	app TeamApp
}

func NewLogTeamApp(
	app TeamApp,
) TeamApp {
	// 校验构造函数依赖
	if app == nil {
		zap.L().Panic(
			"NewLogTeamApp: 依赖项不能为空",
			zap.Bool("app_nil", app == nil),
		)
	}

	// 返回日志包装实现
	return &logTeamAppImpl{app: app}
}

func (a *logTeamAppImpl) Create(
	cx context.Context,
	currUserID string,
	args *val.CreateTeamArgs,
) (*val.CreateTeamRes, error) {
	// 校验包装器实例本身是否合法
	if a == nil || a.app == nil {
		return nil, errors.New("TeamApp 不可用")
	}

	// 校验参数
	if args == nil || args.Name == "" {
		return nil, errors.New("汉化组名称不能为空")
	}

	// 为当前调用构造带上下文的日志记录器
	lgr := retrieveLgr(cx).With(
		zap.String("method", "Create"),
		zap.String("curr_user_id", currUserID),
	)

	// 将日志记录器注入上下文
	cx = injectLgr(cx, lgr)

	// 记录方法调用日志
	lgr.Info("[logTeamAppImpl.Create] CALL")

	// 转发调用到真实实现
	return a.app.Create(cx, currUserID, args)
}

func (a *logTeamAppImpl) List(
	cx context.Context,
	currUserID string,
	args *val.ListTeamArgs,
) ([]*val.TeamInfo, error) {
	// 校验包装器实例本身是否合法
	if a == nil || a.app == nil {
		return nil, errors.New("TeamApp 不可用")
	}

	// 为当前调用构造带上下文的日志记录器
	lgr := retrieveLgr(cx).With(
		zap.String("method", "List"),
		zap.String("curr_user_id", currUserID),
	)

	// 将日志记录器注入上下文
	cx = injectLgr(cx, lgr)

	// 记录方法调用日志
	lgr.Info("[logTeamAppImpl.List] CALL")

	// 转发调用到真实实现
	return a.app.List(cx, currUserID, args)
}

func (a *logTeamAppImpl) ListMy(
	cx context.Context,
	currUserID string,
	args *val.ListMyTeamArgs,
) ([]*val.TeamInfo, error) {
	// 校验包装器实例本身是否合法
	if a == nil || a.app == nil {
		return nil, errors.New("TeamApp 不可用")
	}

	// 为当前调用构造带上下文的日志记录器
	lgr := retrieveLgr(cx).With(
		zap.String("method", "ListMy"),
		zap.String("curr_user_id", currUserID),
	)

	// 将日志记录器注入上下文
	cx = injectLgr(cx, lgr)

	// 记录方法调用日志
	lgr.Info("[logTeamAppImpl.ListMy] CALL")

	// 转发调用到真实实现
	return a.app.ListMy(cx, currUserID, args)
}

func (a *logTeamAppImpl) Update(
	cx context.Context,
	currUserID string,
	args *val.UpdateTeamArgs,
) error {
	// 校验包装器实例本身是否合法
	if a == nil || a.app == nil {
		return errors.New("TeamApp 不可用")
	}

	// 校验参数
	if args == nil || args.ID == "" || args.Name == "" {
		return errors.New("参数不合法")
	}

	// 为当前调用构造带上下文的日志记录器
	lgr := retrieveLgr(cx).With(
		zap.String("method", "Update"),
		zap.String("curr_user_id", currUserID),
		zap.String("team_id", args.ID),
	)

	// 将日志记录器注入上下文
	cx = injectLgr(cx, lgr)

	// 记录方法调用日志
	lgr.Info("[logTeamAppImpl.Update] CALL")

	// 转发调用到真实实现
	return a.app.Update(cx, currUserID, args)
}

func (a *logTeamAppImpl) Remove(
	cx context.Context,
	currUserID string,
	teamID string,
) error {
	// 校验包装器实例本身是否合法
	if a == nil || a.app == nil {
		return errors.New("TeamApp 不可用")
	}

	// 校验汉化组 ID
	if teamID == "" {
		return errors.New("汉化组 ID 不能为空")
	}

	// 为当前调用构造带上下文的日志记录器
	lgr := retrieveLgr(cx).With(
		zap.String("method", "Remove"),
		zap.String("curr_user_id", currUserID),
		zap.String("team_id", teamID),
	)

	// 将日志记录器注入上下文
	cx = injectLgr(cx, lgr)

	// 记录方法调用日志
	lgr.Info("[logTeamAppImpl.Remove] CALL")

	// 转发调用到真实实现
	return a.app.Remove(cx, currUserID, teamID)
}

func (a *logTeamAppImpl) ReserveAvatar(
	cx context.Context,
	currUserID string,
	args *val.ReserveTeamAvatarArgs,
) (*val.ReserveTeamAvatarRes, error) {
	// 校验包装器实例本身是否合法
	if a == nil || a.app == nil {
		return nil, errors.New("TeamApp 不可用")
	}

	// 校验上传参数
	if args == nil || args.TeamID == "" {
		return nil, errors.New("汉化组 ID 不能为空")
	}

	if args.FileName == "" {
		return nil, errors.New("文件名不能为空")
	}

	// 为当前调用构造带上下文的日志记录器
	lgr := retrieveLgr(cx).With(
		zap.String("method", "ReserveAvatar"),
		zap.String("curr_user_id", currUserID),
		zap.String("team_id", args.TeamID),
	)

	// 将日志记录器注入上下文
	cx = injectLgr(cx, lgr)

	// 记录方法调用日志
	lgr.Info("[logTeamAppImpl.ReserveAvatar] CALL")

	// 转发调用到真实实现
	return a.app.ReserveAvatar(cx, currUserID, args)
}

func (a *logTeamAppImpl) ConfirmAvatarUploaded(
	cx context.Context,
	currUserID string,
	teamID string,
) error {
	// 校验包装器实例本身是否合法
	if a == nil || a.app == nil {
		return errors.New("TeamApp 不可用")
	}

	// 校验汉化组 ID
	if teamID == "" {
		return errors.New("汉化组 ID 不能为空")
	}

	// 为当前调用构造带上下文的日志记录器
	lgr := retrieveLgr(cx).With(
		zap.String("method", "ConfirmAvatarUploaded"),
		zap.String("curr_user_id", currUserID),
		zap.String("team_id", teamID),
	)

	// 将日志记录器注入上下文
	cx = injectLgr(cx, lgr)

	// 记录方法调用日志
	lgr.Info("[logTeamAppImpl.ConfirmAvatarUploaded] CALL")

	// 转发调用到真实实现
	return a.app.ConfirmAvatarUploaded(cx, currUserID, teamID)
}
