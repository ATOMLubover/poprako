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

type MemberApp interface {
	// Create 创建成员记录（超级管理员专用）
	Create(
		cx context.Context,
		currUserID string,
		args *val.CreateMemberArgs,
	) (*val.CreateMemberRes, error)

	// ListByTeam 获取指定汉化组的成员列表
	ListByTeam(
		cx context.Context,
		currUserID string,
		args *val.ListTeamMemberArgs,
	) ([]*val.MemberInfo, error)

	// ListMy 获取当前用户的所有成员记录
	ListMy(
		cx context.Context,
		currUserID string,
		args *val.ListMyMemberArgs,
	) ([]*val.MemberInfo, error)

	// UpdateRole 更新成员角色（PUT 语义全量替换）
	UpdateRole(
		cx context.Context,
		currUserID string,
		args *val.UpdateMemberRoleArgs,
	) error

	// Remove 删除成员记录
	Remove(
		cx context.Context,
		currUserID string,
		memberID string,
	) error

	// JoinTeam 通过邀请码加入汉化组
	JoinTeam(
		cx context.Context,
		currUserID string,
		args *val.JoinTeamArgs,
	) error
}

type memberAppImpl struct {
	memberSvc service.MemberService

	userRepo   repo.UserRepo
	memberRepo repo.MemberRepo
	invRepo    repo.InvitationRepo
	txnMgr     repo.TxnMgr

	ossClient oss.Client
}

func NewMemberApp(
	memberSvc service.MemberService,
	userRepo repo.UserRepo,
	memberRepo repo.MemberRepo,
	invRepo repo.InvitationRepo,
	txnMgr repo.TxnMgr,
	ossClient oss.Client,
) MemberApp {
	// 校验构造函数依赖
	if memberSvc == nil ||
		userRepo == nil ||
		memberRepo == nil ||
		invRepo == nil ||
		txnMgr == nil ||
		ossClient == nil {
		zap.L().Panic(
			"NewMemberApp: 依赖项不能为空",
			zap.Bool("memberSvc_nil", memberSvc == nil),
			zap.Bool("userRepo_nil", userRepo == nil),
			zap.Bool("memberRepo_nil", memberRepo == nil),
			zap.Bool("invRepo_nil", invRepo == nil),
			zap.Bool("txnMgr_nil", txnMgr == nil),
			zap.Bool("ossClient_nil", ossClient == nil),
		)
	}

	// 返回真实业务实现
	return &memberAppImpl{
		memberSvc:  memberSvc,
		userRepo:   userRepo,
		memberRepo: memberRepo,
		invRepo:    invRepo,
		txnMgr:     txnMgr,
		ossClient:  ossClient,
	}
}

func (a *memberAppImpl) Create(
	cx context.Context,
	currUserID string,
	args *val.CreateMemberArgs,
) (*val.CreateMemberRes, error) {
	// 获取上下文中的日志记录器
	lgr := retrieveLgr(cx)

	// 查询当前用户信息，用于权限校验
	currUser, err := a.userRepo.GetByID(currUserID)
	if err != nil {
		// 记录查询失败
		lgr.Warn(
			"创建成员失败：无法获取当前用户信息",
			zap.String("curr_user_id", currUserID),
			zap.Error(err),
		)

		// 返回客户端可展示的错误
		return nil, errors.New("创建成员失败：无法获取用户信息")
	}

	// 检查目标用户是否已是该汉化组成员
	isExisting, err := a.memberRepo.Exist(model.MemberQueryOpt{
		UserID: &args.UserID,
		TeamID: &args.TeamID,
	})
	if err != nil {
		// 记录查询失败
		lgr.Error(
			"创建成员失败：检查成员信息失败",
			zap.Error(err),
		)

		// 返回客户端可展示的错误
		return nil, errors.New("创建成员失败")
	}

	// 若已存在，返回重复错误
	if isExisting {
		return nil, errors.New("该用户已经加入该汉化组")
	}

	// 通过领域服务构造成员创建载荷（含权限校验）
	creation, err := a.memberSvc.NewCreation(
		currUser,
		args.UserID,
		args.TeamID,
		model.UnmaskRoles(args.Roles)...,
	)
	if err != nil {
		// 记录权限校验失败
		lgr.Warn(
			"创建成员失败：权限不足",
			zap.String("curr_user_id", currUserID),
			zap.Error(err),
		)

		// 返回领域服务返回的错误
		return nil, err
	}

	// 持久化成员信息
	memberInfo, err := a.memberRepo.Create(creation)
	if err != nil {
		// 记录创建失败
		lgr.Error(
			"创建成员失败：持久化失败",
			zap.Error(err),
		)

		// 返回客户端可展示的错误
		return nil, errors.New("创建成员失败")
	}

	// 返回创建结果
	return &val.CreateMemberRes{ID: memberInfo.ID}, nil
}

func (a *memberAppImpl) ListByTeam(
	cx context.Context,
	currUserID string,
	args *val.ListTeamMemberArgs,
) ([]*val.MemberInfo, error) {
	// 获取上下文中的日志记录器
	lgr := retrieveLgr(cx)

	// 鉴权：检查当前用户是否为该汉化组成员
	_, err := a.memberRepo.Get(model.MemberQueryOpt{
		UserID: &currUserID,
		TeamID: &args.TeamID,
	})
	if err != nil {
		// 记录权限校验失败
		lgr.Warn(
			"获取汉化组成员列表失败：权限不足",
			zap.String("curr_user_id", currUserID),
			zap.String("team_id", args.TeamID),
		)

		// 返回客户端可展示的错误
		return nil, errors.New("权限不足")
	}

	// 查询汉化组成员列表
	members, err := a.memberRepo.List(model.MemberQueryOpt{
		TeamID: &args.TeamID,
	})
	if err != nil {
		// 记录查询失败
		lgr.Error(
			"获取汉化组成员列表失败",
			zap.String("team_id", args.TeamID),
			zap.Error(err),
		)

		// 返回客户端可展示的错误
		return nil, errors.New("获取汉化组成员列表失败")
	}

	// 组装为 app 层值对象列表
	result := make([]*val.MemberInfo, len(members))

	for i, member := range members {
		result[i] = assembleMemberInfo(&member, a.ossClient)
	}

	// 返回成员列表
	return result, nil
}

func (a *memberAppImpl) ListMy(
	cx context.Context,
	currUserID string,
	args *val.ListMyMemberArgs,
) ([]*val.MemberInfo, error) {
	// 获取上下文中的日志记录器
	lgr := retrieveLgr(cx)

	// 查询当前用户的所有成员记录
	members, err := a.memberRepo.List(model.MemberQueryOpt{
		UserID: &currUserID,
	})
	if err != nil {
		// 记录查询失败
		lgr.Error(
			"获取我的成员列表失败",
			zap.String("curr_user_id", currUserID),
			zap.Error(err),
		)

		// 返回客户端可展示的错误
		return nil, errors.New("获取我的成员列表失败")
	}

	// 组装为 app 层值对象列表
	result := make([]*val.MemberInfo, len(members))

	for i, member := range members {
		result[i] = assembleMemberInfo(&member, a.ossClient)
	}

	// 返回成员列表
	return result, nil
}

func (a *memberAppImpl) UpdateRole(
	cx context.Context,
	currUserID string,
	args *val.UpdateMemberRoleArgs,
) error {
	// 获取上下文中的日志记录器
	lgr := retrieveLgr(cx)

	// 查询目标成员信息，获取可信的 TeamID
	targetMember, err := a.memberRepo.GetByID(args.ID)
	if err != nil {
		// 记录查询失败
		lgr.Error(
			"更新成员角色失败：获取目标成员信息失败",
			zap.String("member_id", args.ID),
			zap.Error(err),
		)

		// 返回客户端可展示的错误
		return errors.New("更新成员角色失败：无法获取成员信息")
	}

	// 通过领域服务构造更新载荷（含权限校验）
	update, err := a.memberSvc.NewUpdate(
		a.memberRepo,
		currUserID,
		targetMember.TeamID,
		args.Roles,
	)
	if err != nil {
		// 记录权限校验失败
		lgr.Warn(
			"更新成员角色失败：权限不足",
			zap.String("curr_user_id", currUserID),
			zap.Error(err),
		)

		// 返回领域服务返回的错误
		return err
	}

	// 设置正确的成员 ID
	update.ID = args.ID

	// 持久化更新
	if err := a.memberRepo.Update(update); err != nil {
		// 记录更新失败
		lgr.Error(
			"更新成员角色失败",
			zap.String("member_id", args.ID),
			zap.Error(err),
		)

		// 返回客户端可展示的错误
		return errors.New("更新成员角色失败")
	}

	// 返回更新成功
	return nil
}

func (a *memberAppImpl) Remove(
	cx context.Context,
	currUserID string,
	memberID string,
) error {
	// 获取上下文中的日志记录器
	lgr := retrieveLgr(cx)

	// 查询目标成员信息，获取可信的 TeamID
	targetMember, err := a.memberRepo.GetByID(memberID)
	if err != nil {
		// 记录查询失败
		lgr.Error(
			"删除成员失败：获取目标成员信息失败",
			zap.String("member_id", memberID),
			zap.Error(err),
		)

		// 返回客户端可展示的错误
		return errors.New("删除成员失败：无法获取成员信息")
	}

	// 鉴权：检查当前用户在目标汉化组中是否为管理员
	currMember, err := a.memberRepo.Get(model.MemberQueryOpt{
		UserID: &currUserID,
		TeamID: &targetMember.TeamID,
	})
	if err != nil || !currMember.HasAnyRole(model.RoleAdmin) {
		// 记录权限校验失败
		lgr.Warn(
			"删除成员失败：权限不足",
			zap.String("curr_user_id", currUserID),
			zap.String("team_id", targetMember.TeamID),
		)

		// 返回客户端可展示的错误
		return errors.New("权限不足")
	}

	// 执行删除
	if err := a.memberRepo.Delete(memberID); err != nil {
		// 记录删除失败
		lgr.Error(
			"删除成员失败",
			zap.String("member_id", memberID),
			zap.Error(err),
		)

		// 返回客户端可展示的错误
		return errors.New("删除成员失败")
	}

	// 返回删除成功
	return nil
}

func (a *memberAppImpl) JoinTeam(
	cx context.Context,
	currUserID string,
	args *val.JoinTeamArgs,
) error {
	// 获取上下文中的日志记录器
	lgr := retrieveLgr(cx)

	// 查询当前用户信息
	currUser, err := a.userRepo.GetByID(currUserID)
	if err != nil {
		// 记录查询失败
		lgr.Error(
			"加入汉化组失败：无法获取用户信息",
			zap.String("curr_user_id", currUserID),
			zap.Error(err),
		)

		// 返回客户端可展示的错误
		return errors.New("加入汉化组失败：无法获取用户信息")
	}

	// 根据 QQ 查找待消耗的邀请
	inv, err := a.invRepo.GetByInviteeQQ(currUser.QQ)
	if err != nil {
		// 记录查询失败
		lgr.Warn(
			"加入汉化组失败：无效的邀请码",
			zap.String("qq", currUser.QQ),
			zap.Error(err),
		)

		// 返回客户端可展示的错误
		return errors.New("邀请码无效或已被使用")
	}

	// 校验邀请码
	if inv.InvitationCode != args.InvitationCode || !inv.Pending {
		// 返回客户端可展示的错误
		return errors.New("邀请码无效或已被使用")
	}

	// 检查用户是否已经是该汉化组成员
	isExisting, err := a.memberRepo.Exist(model.MemberQueryOpt{
		UserID: &currUserID,
		TeamID: &inv.TeamID,
	})
	if err != nil {
		// 记录查询失败
		lgr.Error(
			"加入汉化组失败：检查成员信息失败",
			zap.Error(err),
		)

		// 返回客户端可展示的错误
		return errors.New("加入汉化组失败")
	}

	// 若已是成员，返回错误
	if isExisting {
		return errors.New("您已经是该汉化组的成员")
	}

	// 通过领域服务构造成员创建载荷
	creation, err := a.memberSvc.NewCreationFromInvitation(currUser, inv)
	if err != nil {
		// 记录构造失败
		lgr.Error(
			"加入汉化组失败：构造成员载荷失败",
			zap.Error(err),
		)

		// 返回客户端可展示的错误
		return errors.New("加入汉化组失败")
	}

	// 在事务中创建成员并使邀请失效
	if err := a.txnMgr.RunInTxn(func(cx context.Context) error {
		memberRepoTxn, err := a.memberRepo.FromTxnCx(cx)
		if err != nil {
			return err
		}

		invRepoTxn, err := a.invRepo.FromTxnCx(cx)
		if err != nil {
			return err
		}

		// 持久化成员信息
		if _, err := memberRepoTxn.Create(creation); err != nil {
			return err
		}

		// 使邀请失效
		if err := invRepoTxn.Invalidate(inv.ID); err != nil {
			return err
		}

		// 事务执行成功
		return nil
	}); err != nil {
		// 记录事务执行失败
		lgr.Error(
			"加入汉化组失败：事务执行失败",
			zap.Error(err),
		)

		// 返回客户端可展示的错误
		return errors.New("加入汉化组失败")
	}

	// 返回加入成功
	return nil
}

// assembleMemberInfo 将领域层成员信息转换为 app 层值对象
func assembleMemberInfo(
	info *model.MemberInfo,
	ossClient oss.Client,
) *val.MemberInfo {
	// 组装基础成员信息
	result := &val.MemberInfo{
		ID:        info.ID,
		UserID:    info.UserID,
		TeamID:    info.TeamID,
		Roles:     model.MaskRoles(info.Roles()),
		CreatedAt: info.CreatedAt.UnixMilli(),
		UpdatedAt: info.UpdatedAt.UnixMilli(),
	}

	// 如果包含了用户信息，则组装用户子对象
	if info.User != nil {
		userInfo, _ := assembleUserInfo(info.User, ossClient)

		result.User = userInfo
	}

	// 如果包含了汉化组信息，则组装汉化组子对象
	if info.Team != nil {
		teamInfo, _ := assembleTeamInfo(info.Team, ossClient)

		result.Team = teamInfo
	}

	// 返回组装后的成员信息
	return result
}

// logMemberAppImpl 是 MemberApp 的日志包装实现
type logMemberAppImpl struct {
	app MemberApp
}

func NewLogMemberApp(
	app MemberApp,
) MemberApp {
	// 校验构造函数依赖
	if app == nil {
		zap.L().Panic(
			"NewLogMemberApp: 依赖项不能为空",
			zap.Bool("app_nil", app == nil),
		)
	}

	// 返回日志包装实现
	return &logMemberAppImpl{app: app}
}

func (a *logMemberAppImpl) Create(
	cx context.Context,
	currUserID string,
	args *val.CreateMemberArgs,
) (*val.CreateMemberRes, error) {
	if a == nil || a.app == nil {
		return nil, errors.New("MemberApp 不可用")
	}

	if args == nil || args.UserID == "" || args.TeamID == "" {
		return nil, errors.New("参数不合法")
	}

	lgr := retrieveLgr(cx).With(zap.String("method", "Create"), zap.String("curr_user_id", currUserID))

	cx = injectLgr(cx, lgr)

	lgr.Info("[logMemberAppImpl.Create] CALL")

	return a.app.Create(cx, currUserID, args)
}

func (a *logMemberAppImpl) ListByTeam(
	cx context.Context,
	currUserID string,
	args *val.ListTeamMemberArgs,
) ([]*val.MemberInfo, error) {
	if a == nil || a.app == nil {
		return nil, errors.New("MemberApp 不可用")
	}

	if args == nil || args.TeamID == "" {
		return nil, errors.New("参数不合法")
	}

	lgr := retrieveLgr(cx).With(zap.String("method", "ListByTeam"), zap.String("curr_user_id", currUserID), zap.String("team_id", args.TeamID))

	cx = injectLgr(cx, lgr)

	lgr.Info("[logMemberAppImpl.ListByTeam] CALL")

	return a.app.ListByTeam(cx, currUserID, args)
}

func (a *logMemberAppImpl) ListMy(
	cx context.Context,
	currUserID string,
	args *val.ListMyMemberArgs,
) ([]*val.MemberInfo, error) {
	if a == nil || a.app == nil {
		return nil, errors.New("MemberApp 不可用")
	}

	lgr := retrieveLgr(cx).With(zap.String("method", "ListMy"), zap.String("curr_user_id", currUserID))

	cx = injectLgr(cx, lgr)

	lgr.Info("[logMemberAppImpl.ListMy] CALL")

	return a.app.ListMy(cx, currUserID, args)
}

func (a *logMemberAppImpl) UpdateRole(
	cx context.Context,
	currUserID string,
	args *val.UpdateMemberRoleArgs,
) error {
	if a == nil || a.app == nil {
		return errors.New("MemberApp 不可用")
	}

	if args == nil || args.ID == "" {
		return errors.New("参数不合法")
	}

	lgr := retrieveLgr(cx).With(zap.String("method", "UpdateRole"), zap.String("curr_user_id", currUserID), zap.String("member_id", args.ID))

	cx = injectLgr(cx, lgr)

	lgr.Info("[logMemberAppImpl.UpdateRole] CALL")

	return a.app.UpdateRole(cx, currUserID, args)
}

func (a *logMemberAppImpl) Remove(
	cx context.Context,
	currUserID string,
	memberID string,
) error {
	if a == nil || a.app == nil {
		return errors.New("MemberApp 不可用")
	}

	if memberID == "" {
		return errors.New("成员 ID 不能为空")
	}

	lgr := retrieveLgr(cx).With(zap.String("method", "Remove"), zap.String("curr_user_id", currUserID), zap.String("member_id", memberID))

	cx = injectLgr(cx, lgr)

	lgr.Info("[logMemberAppImpl.Remove] CALL")

	return a.app.Remove(cx, currUserID, memberID)
}

func (a *logMemberAppImpl) JoinTeam(
	cx context.Context,
	currUserID string,
	args *val.JoinTeamArgs,
) error {
	if a == nil || a.app == nil {
		return errors.New("MemberApp 不可用")
	}

	if args == nil || args.InvitationCode == "" {
		return errors.New("邀请码不能为空")
	}

	lgr := retrieveLgr(cx).With(zap.String("method", "JoinTeam"), zap.String("curr_user_id", currUserID))

	cx = injectLgr(cx, lgr)

	lgr.Info("[logMemberAppImpl.JoinTeam] CALL")

	return a.app.JoinTeam(cx, currUserID, args)
}
