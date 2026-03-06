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

type InvitationApplication interface {
	ListInvitations(
		scope util.TraceScope,
		currentUserID string,
		teamID string,
	) ([]*value.InvitationInfo, error)
	CreateInvitation(
		scope util.TraceScope,
		currentUserID string,
		args *value.CreateInvitationArgs,
	) (*value.InvitationInfo, error)
	UpdateInvitation(
		scope util.TraceScope,
		currentUserID string,
		args *value.UpdateInvitationArgs,
	) error
	DeleteInvitation(
		scope util.TraceScope,
		currentUserID string,
		invitationID string,
	) error
}

type invitationApplication struct {
	userRepository       repository.UserRepository
	memberRepository     repository.MemberRepository
	invitationRepository repository.InvitationRepository
}

func NewInvitationApplication(
	userRepository repository.UserRepository,
	memberRepository repository.MemberRepository,
	invitationRepository repository.InvitationRepository,
) InvitationApplication {
	if userRepository == nil ||
		invitationRepository == nil ||
		memberRepository == nil {
		zap.L().Panic(
			"NewInvitationApplication: 依赖项不能为空",
			zap.Bool("userRepository_nil", userRepository == nil),
			zap.Bool("memberRepository_nil", memberRepository == nil),
			zap.Bool("invitationRepository_nil", invitationRepository == nil),
		)
	}

	return &invitationApplication{
		userRepository:       userRepository,
		memberRepository:     memberRepository,
		invitationRepository: invitationRepository,
	}
}

func (ia *invitationApplication) ListInvitations(
	scope util.TraceScope,
	currentUserID string,
	targetTeamID string,
) ([]*value.InvitationInfo, error) {
	const fn = "InvitationApplication.ListInvitations"

	if currentUserID == "" {
		scope.Logger().Warn(fn + ": currentUserID 为空")
		return nil, errors.New(ErrInternalError)
	}

	scope.WithFields(
		zap.String("current_user_id", currentUserID),
		zap.String("target_team_id", targetTeamID),
	)

	scope.Logger().Debug(fn + ": 被调用")

	// 鉴权：检查当前用户在指定的汉化组是否有权限查看邀请信息
	if err := checkCurrentUserPermissionInTeam(
		scope,
		ia.memberRepository,
		currentUserID,
		targetTeamID,
		model.PermissionInvitationList,
	); err != nil {
		scope.Logger().Warn(fn+": 权限检查不通过", zap.Error(err))
		return nil, err
	}

	// 获取邀请信息列表
	invitationList, err := ia.invitationRepository.List(
		nil,
		query_option.InvitationQuery().FilterByTeamID(targetTeamID),
		query_option.CreatedAtDesc(),
	)
	if err != nil {
		scope.Logger().Error(fn+": 获取指定汉化组邀请信息列表失败", zap.Error(err))
		return nil, errors.New("无法获取邀请信息列表")
	}

	// 转换为应用层的值对
	result := make([]*value.InvitationInfo, len(invitationList))
	for i, info := range invitationList {
		result[i] = value.NewInvitationInfoFromModel(&info)
	}

	return result, nil
}

func (ia *invitationApplication) CreateInvitation(
	scope util.TraceScope,
	currentUserID string,
	args *value.CreateInvitationArgs,
) (*value.InvitationInfo, error) {
	const fn = "InvitationApplication.CreateInvitation"

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

	// 鉴权：检查当前用户是否有权限创建邀请
	if err := checkCurrentUserPermissionInTeam(
		scope,
		ia.memberRepository,
		currentUserID,
		args.TeamID,
		model.PermissionInvitationCreate,
	); err != nil {
		scope.Logger().Warn(fn+": 权限检查不通过", zap.Error(err))
		return nil, err
	}

	// 检查是否已经有相关 QQ 的成员存在
	isMemberExisting, err := ia.memberRepository.Exist(
		nil,
		query_option.MemberQuery().JoinUser(),
		query_option.MemberQuery().FilterByTeamID(args.TeamID),
		query_option.MemberQuery().FilterOnUserQQ(args.InviteeQQ),
	)
	if err != nil {
		scope.Logger().Error(fn+": 检查成员信息失败", zap.Error(err))
		return nil, errors.New("无法检查成员信息")
	}

	if isMemberExisting {
		return nil, errors.New("该用户已经加入该汉化组")
	}

	// 创建邀请信息
	invitationCode, err := service.GenerateInvitationCode()
	if err != nil {
		scope.Logger().Error(fn+": 生成邀请码失败", zap.Error(err))
		return nil, errors.New("生成邀请码失败")
	}

	invitationCreation := model.NewInvitationCreation(
		currentUserID,
		args.TeamID,
		args.InviteeQQ,
		invitationCode,
		model.UnmaskRoles(args.Roles)...,
	)

	// 此处应当使用乐观锁进行更新，不采用事务
	invitationID, err := ia.invitationRepository.Create(nil, invitationCreation)
	if err != nil {
		scope.Logger().Error(fn+": 创建邀请信息失败", zap.Error(err))
		return nil, errors.New("创建邀请失败")
	}

	// 构造创建成功的邀请信息值对象，不再查询数据库获取邀请信息
	invitationInfo := value.NewInvitationInfo(
		invitationID,
		currentUserID,
		args.InviteeQQ,
		invitationCode,
		true,
		args.Roles,
		util.NowMillis(),
	)

	return invitationInfo, nil
}

func (ia *invitationApplication) UpdateInvitation(
	scope util.TraceScope,
	currentUserID string,
	args *value.UpdateInvitationArgs,
) error {
	const fn = "InvitationApplication.UpdateInvitation"

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

	// 鉴权：检查当前用户是否有权限修改邀请
	if err := checkCurrentUserPermissionInTeam(
		scope,
		ia.memberRepository,
		currentUserID,
		args.TeamID,
		model.PermissionInvitationUpdate,
	); err != nil {
		scope.Logger().Warn(fn+": 权限检查不通过", zap.Error(err))
		return err
	}

	invitationUpdate := model.NewInvitationUpdate(
		args.ID,
		model.UnmaskRoles(args.Roles)...,
	)

	if err := ia.invitationRepository.Update(
		nil,
		invitationUpdate,
	); err != nil {
		scope.Logger().Error(fn+": 更新邀请信息失败", zap.Error(err))
		return errors.New("更新邀请失败")
	}

	return nil
}

func (ia *invitationApplication) DeleteInvitation(
	scope util.TraceScope,
	currentUserID string,
	invitationID string,
) error {
	const fn = "InvitationApplication.DeleteInvitation"

	scope.WithFields(
		zap.String("current_user_id", currentUserID),
		zap.String("invitation_id", invitationID),
	)

	scope.Logger().Debug(fn + ": 被调用")

	// 鉴权：检查当前用户是否有权限删除邀请
	currentUserMemberships, err := ia.memberRepository.List(
		nil,
		query_option.MemberQuery().FilterByUserID(currentUserID),
	)
	if err != nil {
		scope.Logger().Error(fn+": 获取当前用户的成员信息失败", zap.Error(err))
		return errors.New("无法获取成员信息")
	}

	invitation, err := ia.invitationRepository.GetByID(nil, invitationID)
	if invitation == nil || err != nil {
		scope.Logger().Error(fn+": 获取邀请信息失败", zap.Error(err))
		return errors.New("无法获取邀请信息")
	}

	// 只有 invitation 所对应的汉化组的管理员才有权限删除邀请
	if !service.CheckInvitationPermission(
		invitation.TeamID,
		currentUserMemberships,
		model.PermissionInvitationDelete,
	) {
		return errors.New("没有权限删除邀请")
	}

	// 删除邀请信息
	err = ia.invitationRepository.DeleteByID(nil, invitationID)
	if err != nil {
		scope.Logger().Error(fn+": 删除邀请信息失败", zap.Error(err))
		return errors.New("删除邀请失败")
	}

	return nil
}
