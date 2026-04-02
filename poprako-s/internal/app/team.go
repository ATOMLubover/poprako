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

type TeamApp interface {
	// Create 创建一个新的团队，仅超级管理员可操作
	Create(
		cx context.Context,
		currUserID string,
		args *val.CreateTeamArgs,
	) (*val.CreateTeamRes, error)

	// List 获取所有团队列表（超级管理员专用）
	List(
		cx context.Context,
		currUserID string,
		args *val.ListTeamArgs,
	) ([]*val.TeamInfo, error)

	// ListMy 获取当前用户所属的团队列表
	ListMy(
		cx context.Context,
		currUserID string,
		args *val.ListMyTeamArgs,
	) ([]*val.TeamInfo, error)

	// Update 更新团队信息
	Update(
		cx context.Context,
		currUserID string,
		args *val.UpdateTeamArgs,
	) error

	// Remove 删除团队
	Remove(
		cx context.Context,
		currUserID string,
		teamID string,
	) error

	// ReserveAvatar 预留团队头像上传所需的预签名 URL
	ReserveAvatar(
		cx context.Context,
		currUserID string,
		teamID string,
	) (*val.ReserveTeamAvatarRes, error)

	// ConfirmAvatarUploaded 确认团队头像已经完成上传
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
			"创建团队失败：无法获取当前用户信息",
			zap.String("curr_user_id", currUserID),
			zap.Error(err),
		)

		// 返回客户端可展示的错误
		return nil, errors.New("创建团队失败：无法获取用户信息")
	}

	// 通过领域服务构造团队创建载荷（含权限校验）
	creation, err := a.teamSvc.NewCreation(currUser, args.Name, args.Description)
	if err != nil {
		// 记录权限校验失败
		lgr.Warn(
			"创建团队失败：权限不足",
			zap.String("curr_user_id", currUserID),
			zap.Error(err),
		)

		// 返回领域服务返回的错误
		return nil, err
	}

	// 持久化团队信息
	teamInfo, err := a.teamRepo.Create(creation)
	if err != nil {
		// 记录创建失败
		lgr.Error(
			"创建团队失败：持久化失败",
			zap.String("curr_user_id", currUserID),
			zap.Error(err),
		)

		// 返回客户端可展示的错误
		return nil, errors.New("创建团队失败")
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
			"获取团队列表失败：无法获取当前用户信息",
			zap.String("curr_user_id", currUserID),
			zap.Error(err),
		)

		// 返回客户端可展示的错误
		return nil, errors.New("获取团队列表失败")
	}

	// 校验超级管理员权限
	if !currUser.IsSuperAdmin {
		// 记录权限校验失败
		lgr.Warn(
			"获取团队列表失败：仅超级管理员可以查看所有团队",
			zap.String("curr_user_id", currUserID),
		)

		// 返回客户端可展示的错误
		return nil, errors.New("权限不足")
	}

	// 查询所有团队列表
	teams, err := a.teamRepo.List(model.TeamQueryOpt{})
	if err != nil {
		// 记录查询失败
		lgr.Error(
			"获取团队列表失败",
			zap.Error(err),
		)

		// 返回客户端可展示的错误
		return nil, errors.New("获取团队列表失败")
	}

	// 组装为 app 层值对象列表
	result := make([]*val.TeamInfo, len(teams))

	for i, team := range teams {
		// 组装单个团队信息
		info, err := assembleTeamInfo(&team, a.ossClient)
		if err != nil {
			// 记录组装失败
			lgr.Error(
				"获取团队列表失败：组装团队信息失败",
				zap.String("team_id", team.ID),
				zap.Error(err),
			)

			// 返回客户端可展示的错误
			return nil, errors.New("获取团队列表失败")
		}

		result[i] = info
	}

	// 返回团队列表
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
			"获取我的团队列表失败：查询成员记录失败",
			zap.String("curr_user_id", currUserID),
			zap.Error(err),
		)

		// 返回客户端可展示的错误
		return nil, errors.New("获取我的团队列表失败")
	}

	// 若用户未加入任何团队，直接返回空列表
	if len(members) == 0 {
		return nil, nil
	}

	// 逐个查询成员所属的团队信息
	result := make([]*val.TeamInfo, 0, len(members))

	for _, member := range members {
		// 根据团队 ID 查询团队信息
		team, err := a.teamRepo.GetByID(member.TeamID)
		if err != nil {
			// 记录查询失败
			lgr.Error(
				"获取我的团队列表失败：查询团队信息失败",
				zap.String("team_id", member.TeamID),
				zap.Error(err),
			)

			// 返回客户端可展示的错误
			return nil, errors.New("获取我的团队列表失败")
		}

		// 组装单个团队信息
		info, err := assembleTeamInfo(team, a.ossClient)
		if err != nil {
			// 记录组装失败
			lgr.Error(
				"获取我的团队列表失败：组装团队信息失败",
				zap.String("team_id", team.ID),
				zap.Error(err),
			)

			// 返回客户端可展示的错误
			return nil, errors.New("获取我的团队列表失败")
		}

		result = append(result, info)
	}

	// 返回团队列表
	return result, nil
}

func (a *teamAppImpl) Update(
	cx context.Context,
	currUserID string,
	args *val.UpdateTeamArgs,
) error {
	// 获取上下文中的日志记录器
	lgr := retrieveLgr(cx)

	// 查询当前用户在目标团队中的成员信息，用于鉴权
	currMember, err := a.memberRepo.Get(model.MemberQueryOpt{
		UserID: &currUserID,
		TeamID: &args.ID,
	})
	if err != nil || !currMember.HasAnyRole(model.RoleAdmin) {
		// 记录权限校验失败
		lgr.Warn(
			"更新团队失败：权限不足",
			zap.String("curr_user_id", currUserID),
			zap.String("team_id", args.ID),
		)

		// 返回客户端可展示的错误
		return errors.New("权限不足")
	}

	// 组装团队更新载荷
	update := &model.TeamUpdate{
		ID:          args.ID,
		Name:        args.Name,
		Description: args.Description,
	}

	// 持久化更新
	if err := a.teamRepo.Update(update); err != nil {
		// 记录更新失败
		lgr.Error(
			"更新团队失败",
			zap.String("team_id", args.ID),
			zap.Error(err),
		)

		// 返回客户端可展示的错误
		return errors.New("更新团队失败")
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
			"删除团队失败：无法获取当前用户信息",
			zap.String("curr_user_id", currUserID),
			zap.Error(err),
		)

		// 返回客户端可展示的错误
		return errors.New("删除团队失败")
	}

	// 校验超级管理员权限
	if !currUser.IsSuperAdmin {
		// 记录权限校验失败
		lgr.Warn(
			"删除团队失败：权限不足",
			zap.String("curr_user_id", currUserID),
		)

		// 返回客户端可展示的错误
		return errors.New("权限不足")
	}

	// 执行删除
	if err := a.teamRepo.Delete(teamID); err != nil {
		// 记录删除失败
		lgr.Error(
			"删除团队失败",
			zap.String("team_id", teamID),
			zap.Error(err),
		)

		// 返回客户端可展示的错误
		return errors.New("删除团队失败")
	}

	// 返回删除成功
	return nil
}

func (a *teamAppImpl) ReserveAvatar(
	cx context.Context,
	currUserID string,
	teamID string,
) (*val.ReserveTeamAvatarRes, error) {
	// 获取上下文中的日志记录器
	lgr := retrieveLgr(cx)

	// 查询当前用户在目标团队中的成员信息，用于鉴权
	currMember, err := a.memberRepo.Get(model.MemberQueryOpt{
		UserID: &currUserID,
		TeamID: &teamID,
	})
	if err != nil || !currMember.HasAnyRole(model.RoleAdmin) {
		// 记录权限校验失败
		lgr.Warn(
			"预留团队头像失败：权限不足",
			zap.String("curr_user_id", currUserID),
			zap.String("team_id", teamID),
		)

		// 返回客户端可展示的错误
		return nil, errors.New("权限不足")
	}

	// 基于团队 ID 生成头像对象 Key
	avatarOSSKey := a.teamSvc.GenAvatarOSSKey(teamID)

	// 为客户端生成预签名上传链接
	putURL, err := a.ossClient.GeneratePutPresignedURL(avatarOSSKey)
	if err != nil {
		// 记录上传链接生成失败
		lgr.Error(
			"预留团队头像失败：生成上传链接失败",
			zap.String("team_id", teamID),
			zap.Error(err),
		)

		// 返回客户端可展示的错误
		return nil, errors.New("预留团队头像失败")
	}

	// 在数据库中预填充头像对象 Key
	if err := a.teamRepo.PreFillAvatarOSSKey(teamID, avatarOSSKey); err != nil {
		// 记录预写失败
		lgr.Error(
			"预留团队头像失败：写入头像 OSS Key 失败",
			zap.String("team_id", teamID),
			zap.Error(err),
		)

		// 返回客户端可展示的错误
		return nil, errors.New("预留团队头像失败")
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

	// 查询当前用户在目标团队中的成员信息，用于鉴权
	currMember, err := a.memberRepo.Get(model.MemberQueryOpt{
		UserID: &currUserID,
		TeamID: &teamID,
	})
	if err != nil || !currMember.HasAnyRole(model.RoleAdmin) {
		// 记录权限校验失败
		lgr.Warn(
			"确认团队头像上传失败：权限不足",
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
			"确认团队头像上传失败",
			zap.String("team_id", teamID),
			zap.Error(err),
		)

		// 返回客户端可展示的错误
		return errors.New("确认团队头像上传失败")
	}

	// 返回确认成功
	return nil
}

// assembleTeamInfo 将领域层团队信息转换为 app 层值对象
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
	lgr *zap.Logger
	app TeamApp
}

func NewLogTeamApp(
	lgr *zap.Logger,
	app TeamApp,
) TeamApp {
	// 校验构造函数依赖
	if lgr == nil || app == nil {
		zap.L().Panic(
			"NewLogTeamApp: 依赖项不能为空",
			zap.Bool("lgr_nil", lgr == nil),
			zap.Bool("app_nil", app == nil),
		)
	}

	// 返回日志包装实现
	return &logTeamAppImpl{lgr: lgr, app: app}
}

func (a *logTeamAppImpl) Create(
	cx context.Context,
	currUserID string,
	args *val.CreateTeamArgs,
) (*val.CreateTeamRes, error) {
	// 校验包装器实例本身是否合法
	if a == nil || a.app == nil || a.lgr == nil {
		return nil, errors.New("TeamApp 不可用")
	}

	// 校验参数
	if args == nil || args.Name == "" {
		return nil, errors.New("团队名称不能为空")
	}

	// 为当前调用构造带上下文的日志记录器
	lgr := a.lgr.With(
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
	if a == nil || a.app == nil || a.lgr == nil {
		return nil, errors.New("TeamApp 不可用")
	}

	// 为当前调用构造带上下文的日志记录器
	lgr := a.lgr.With(
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
	if a == nil || a.app == nil || a.lgr == nil {
		return nil, errors.New("TeamApp 不可用")
	}

	// 为当前调用构造带上下文的日志记录器
	lgr := a.lgr.With(
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
	if a == nil || a.app == nil || a.lgr == nil {
		return errors.New("TeamApp 不可用")
	}

	// 校验参数
	if args == nil || args.ID == "" || args.Name == "" {
		return errors.New("参数不合法")
	}

	// 为当前调用构造带上下文的日志记录器
	lgr := a.lgr.With(
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
	if a == nil || a.app == nil || a.lgr == nil {
		return errors.New("TeamApp 不可用")
	}

	// 校验团队 ID
	if teamID == "" {
		return errors.New("团队 ID 不能为空")
	}

	// 为当前调用构造带上下文的日志记录器
	lgr := a.lgr.With(
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
	teamID string,
) (*val.ReserveTeamAvatarRes, error) {
	// 校验包装器实例本身是否合法
	if a == nil || a.app == nil || a.lgr == nil {
		return nil, errors.New("TeamApp 不可用")
	}

	// 校验团队 ID
	if teamID == "" {
		return nil, errors.New("团队 ID 不能为空")
	}

	// 为当前调用构造带上下文的日志记录器
	lgr := a.lgr.With(
		zap.String("method", "ReserveAvatar"),
		zap.String("curr_user_id", currUserID),
		zap.String("team_id", teamID),
	)

	// 将日志记录器注入上下文
	cx = injectLgr(cx, lgr)

	// 记录方法调用日志
	lgr.Info("[logTeamAppImpl.ReserveAvatar] CALL")

	// 转发调用到真实实现
	return a.app.ReserveAvatar(cx, currUserID, teamID)
}

func (a *logTeamAppImpl) ConfirmAvatarUploaded(
	cx context.Context,
	currUserID string,
	teamID string,
) error {
	// 校验包装器实例本身是否合法
	if a == nil || a.app == nil || a.lgr == nil {
		return errors.New("TeamApp 不可用")
	}

	// 校验团队 ID
	if teamID == "" {
		return errors.New("团队 ID 不能为空")
	}

	// 为当前调用构造带上下文的日志记录器
	lgr := a.lgr.With(
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
