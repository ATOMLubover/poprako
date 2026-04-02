package app

import (
	"context"
	"errors"

	"poprako-s/internal/app/val"
	"poprako-s/internal/domain/model"
	"poprako-s/internal/domain/repo"
	"poprako-s/internal/domain/service"

	"go.uber.org/zap"
)

type InvitationApp interface {
	// List 获取指定团队的邀请列表
	List(
		cx context.Context,
		currUserID string,
		args *val.ListTeamInvitationArgs,
	) ([]*val.InvitationInfo, error)

	// Create 创建一条新的邀请
	Create(
		cx context.Context,
		currUserID string,
		args *val.CreateInvitationArgs,
	) (*val.InvitationInfo, error)

	// Update 更新邀请的角色信息
	Update(
		cx context.Context,
		currUserID string,
		args *val.UpdateInvitationArgs,
	) error

	// Remove 删除邀请
	Remove(
		cx context.Context,
		currUserID string,
		invitationID string,
	) error
}

type invitationAppImpl struct {
	invSvc service.InvitationService

	memberRepo repo.MemberRepo
	invRepo    repo.InvitationRepo
}

func NewInvitationApp(
	invSvc service.InvitationService,
	memberRepo repo.MemberRepo,
	invRepo repo.InvitationRepo,
) InvitationApp {
	// 校验构造函数依赖
	if invSvc == nil ||
		memberRepo == nil ||
		invRepo == nil {
		zap.L().Panic(
			"NewInvitationApp: 依赖项不能为空",
			zap.Bool("invSvc_nil", invSvc == nil),
			zap.Bool("memberRepo_nil", memberRepo == nil),
			zap.Bool("invRepo_nil", invRepo == nil),
		)
	}

	// 返回真实业务实现
	return &invitationAppImpl{
		invSvc:     invSvc,
		memberRepo: memberRepo,
		invRepo:    invRepo,
	}
}

func (a *invitationAppImpl) List(
	cx context.Context,
	currUserID string,
	args *val.ListTeamInvitationArgs,
) ([]*val.InvitationInfo, error) {
	// 获取上下文中的日志记录器
	lgr := retrieveLgr(cx)

	// 鉴权：检查当前用户是否为该团队成员
	_, err := a.memberRepo.Get(model.MemberQueryOpt{
		UserID: &currUserID,
		TeamID: &args.TeamID,
	})
	if err != nil {
		// 记录权限校验失败
		lgr.Warn(
			"获取邀请列表失败：权限不足",
			zap.String("curr_user_id", currUserID),
			zap.String("team_id", args.TeamID),
		)

		// 返回客户端可展示的错误
		return nil, errors.New("权限不足")
	}

	// 查询邀请列表
	invitations, err := a.invRepo.List(model.InvitationQueryOpt{
		TeamID: &args.TeamID,
	})
	if err != nil {
		// 记录查询失败
		lgr.Error(
			"获取邀请列表失败",
			zap.String("team_id", args.TeamID),
			zap.Error(err),
		)

		// 返回客户端可展示的错误
		return nil, errors.New("获取邀请列表失败")
	}

	// 组装为 app 层值对象列表
	result := make([]*val.InvitationInfo, len(invitations))

	for i, inv := range invitations {
		result[i] = assembleInvitationInfo(&inv)
	}

	// 返回邀请列表
	return result, nil
}

func (a *invitationAppImpl) Create(
	cx context.Context,
	currUserID string,
	args *val.CreateInvitationArgs,
) (*val.InvitationInfo, error) {
	// 获取上下文中的日志记录器
	lgr := retrieveLgr(cx)

	// 通过领域服务构造邀请创建载荷（含权限校验）
	creation, err := a.invSvc.NewCreation(
		a.memberRepo,
		currUserID,
		args.TeamID,
		args.InviteeQQ,
		model.UnmaskRoles(args.Roles)...,
	)
	if err != nil {
		// 记录权限校验失败
		lgr.Warn(
			"创建邀请失败：权限不足",
			zap.String("curr_user_id", currUserID),
			zap.Error(err),
		)

		// 返回领域服务返回的错误
		return nil, err
	}

	// 持久化邀请信息
	invInfo, err := a.invRepo.Create(creation)
	if err != nil {
		// 记录创建失败
		lgr.Error(
			"创建邀请失败：持久化失败",
			zap.Error(err),
		)

		// 返回客户端可展示的错误
		return nil, errors.New("创建邀请失败")
	}

	// 组装并返回创建结果
	return assembleInvitationInfo(invInfo), nil
}

func (a *invitationAppImpl) Update(
	cx context.Context,
	currUserID string,
	args *val.UpdateInvitationArgs,
) error {
	// 获取上下文中的日志记录器
	lgr := retrieveLgr(cx)

	// 鉴权：检查当前用户在目标团队中是否为管理员
	currMember, err := a.memberRepo.Get(model.MemberQueryOpt{
		UserID: &currUserID,
		TeamID: &args.TeamID,
	})
	if err != nil || !currMember.HasAnyRole(model.RoleAdmin) {
		// 记录权限校验失败
		lgr.Warn(
			"更新邀请失败：权限不足",
			zap.String("curr_user_id", currUserID),
			zap.String("team_id", args.TeamID),
		)

		// 返回客户端可展示的错误
		return errors.New("权限不足")
	}

	// 解码角色并构造更新载荷
	roles := model.UnmaskRoles(args.Roles)

	update := &model.InvitationUpdate{
		ID: args.ID,
	}

	for _, role := range roles {
		switch role {
		case model.RoleRawProvider:
			update.ToBeRawProvider = true
		case model.RoleTranslator:
			update.ToBeTranslator = true
		case model.RoleProofreader:
			update.ToBeProofreader = true
		case model.RoleTypesetter:
			update.ToBeTypesetter = true
		case model.RoleReviewer:
			update.ToBeReviewer = true
		case model.RolePublisher:
			update.ToBePublisher = true
		case model.RoleAdmin:
			update.ToBeAdmin = true
		}
	}

	// 持久化更新
	if err := a.invRepo.Update(update); err != nil {
		// 记录更新失败
		lgr.Error(
			"更新邀请失败",
			zap.String("invitation_id", args.ID),
			zap.Error(err),
		)

		// 返回客户端可展示的错误
		return errors.New("更新邀请失败")
	}

	// 返回更新成功
	return nil
}

func (a *invitationAppImpl) Remove(
	cx context.Context,
	currUserID string,
	invitationID string,
) error {
	// 获取上下文中的日志记录器
	lgr := retrieveLgr(cx)

	// 查询邀请信息以获取所属团队 ID
	invitations, err := a.invRepo.List(model.InvitationQueryOpt{})
	if err != nil {
		// 记录查询失败
		lgr.Error(
			"删除邀请失败：获取邀请信息失败",
			zap.String("invitation_id", invitationID),
			zap.Error(err),
		)

		// 返回客户端可展示的错误
		return errors.New("删除邀请失败：无法获取邀请信息")
	}

	// 查找目标邀请
	var targetInv *model.InvitationInfo

	for i, inv := range invitations {
		if inv.ID == invitationID {
			targetInv = &invitations[i]

			break
		}
	}

	// 若邀请不存在
	if targetInv == nil {
		return errors.New("邀请不存在")
	}

	// 鉴权：检查当前用户在邀请所属团队中是否为管理员
	currMember, err := a.memberRepo.Get(model.MemberQueryOpt{
		UserID: &currUserID,
		TeamID: &targetInv.TeamID,
	})
	if err != nil || !currMember.HasAnyRole(model.RoleAdmin) {
		// 记录权限校验失败
		lgr.Warn(
			"删除邀请失败：权限不足",
			zap.String("curr_user_id", currUserID),
		)

		// 返回客户端可展示的错误
		return errors.New("权限不足")
	}

	// 执行删除
	if err := a.invRepo.Delete(invitationID); err != nil {
		// 记录删除失败
		lgr.Error(
			"删除邀请失败",
			zap.String("invitation_id", invitationID),
			zap.Error(err),
		)

		// 返回客户端可展示的错误
		return errors.New("删除邀请失败")
	}

	// 返回删除成功
	return nil
}

// assembleInvitationInfo 将领域层邀请信息转换为 app 层值对象
func assembleInvitationInfo(info *model.InvitationInfo) *val.InvitationInfo {
	return &val.InvitationInfo{
		ID:             info.ID,
		InvitorID:      info.InvitorID,
		InviteeQQ:      info.InviteeQQ,
		TeamID:         info.TeamID,
		InvitationCode: info.InvitationCode,
		Pending:        info.Pending,
		Roles:          info.InvitedRoleMask(),
		CreatedAt:      info.CreatedAt.UnixMilli(),
	}
}

// logInvitationAppImpl 是 InvitationApp 的日志包装实现
type logInvitationAppImpl struct {
	lgr *zap.Logger
	app InvitationApp
}

func NewLogInvitationApp(
	lgr *zap.Logger,
	app InvitationApp,
) InvitationApp {
	if lgr == nil || app == nil {
		zap.L().Panic(
			"NewLogInvitationApp: 依赖项不能为空",
			zap.Bool("lgr_nil", lgr == nil),
			zap.Bool("app_nil", app == nil),
		)
	}

	return &logInvitationAppImpl{lgr: lgr, app: app}
}

func (a *logInvitationAppImpl) List(
	cx context.Context,
	currUserID string,
	args *val.ListTeamInvitationArgs,
) ([]*val.InvitationInfo, error) {
	if a == nil || a.app == nil || a.lgr == nil {
		return nil, errors.New("InvitationApp 不可用")
	}

	if args == nil || args.TeamID == "" {
		return nil, errors.New("参数不合法")
	}

	lgr := a.lgr.With(zap.String("method", "List"), zap.String("curr_user_id", currUserID), zap.String("team_id", args.TeamID))

	cx = injectLgr(cx, lgr)

	lgr.Info("[logInvitationAppImpl.List] CALL")

	return a.app.List(cx, currUserID, args)
}

func (a *logInvitationAppImpl) Create(
	cx context.Context,
	currUserID string,
	args *val.CreateInvitationArgs,
) (*val.InvitationInfo, error) {
	if a == nil || a.app == nil || a.lgr == nil {
		return nil, errors.New("InvitationApp 不可用")
	}

	if args == nil || args.TeamID == "" || args.InviteeQQ == "" {
		return nil, errors.New("参数不合法")
	}

	lgr := a.lgr.With(zap.String("method", "Create"), zap.String("curr_user_id", currUserID))

	cx = injectLgr(cx, lgr)

	lgr.Info("[logInvitationAppImpl.Create] CALL")

	return a.app.Create(cx, currUserID, args)
}

func (a *logInvitationAppImpl) Update(
	cx context.Context,
	currUserID string,
	args *val.UpdateInvitationArgs,
) error {
	if a == nil || a.app == nil || a.lgr == nil {
		return errors.New("InvitationApp 不可用")
	}

	if args == nil || args.ID == "" || args.TeamID == "" {
		return errors.New("参数不合法")
	}

	lgr := a.lgr.With(zap.String("method", "Update"), zap.String("curr_user_id", currUserID), zap.String("invitation_id", args.ID))

	cx = injectLgr(cx, lgr)

	lgr.Info("[logInvitationAppImpl.Update] CALL")

	return a.app.Update(cx, currUserID, args)
}

func (a *logInvitationAppImpl) Remove(
	cx context.Context,
	currUserID string,
	invitationID string,
) error {
	if a == nil || a.app == nil || a.lgr == nil {
		return errors.New("InvitationApp 不可用")
	}

	if invitationID == "" {
		return errors.New("邀请 ID 不能为空")
	}

	lgr := a.lgr.With(zap.String("method", "Remove"), zap.String("curr_user_id", currUserID), zap.String("invitation_id", invitationID))

	cx = injectLgr(cx, lgr)

	lgr.Info("[logInvitationAppImpl.Remove] CALL")

	return a.app.Remove(cx, currUserID, invitationID)
}
