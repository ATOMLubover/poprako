package application

import (
	"errors"

	"labelplus-next-web-be/internal/application/adapter"
	"labelplus-next-web-be/internal/domain/model"
	"labelplus-next-web-be/internal/domain/repository"
	"labelplus-next-web-be/internal/domain/service"
	repository_infra "labelplus-next-web-be/internal/repository"
	"labelplus-next-web-be/internal/repository/query_option"
	"labelplus-next-web-be/internal/util"
	"labelplus-next-web-be/internal/value"

	"go.uber.org/zap"
)

type InvitationApplication interface {
	ListInvitations(
		scope util.TraceScope,
		currentUserID string,
		args value.ListTeamInvitationArgs,
	) ([]value.InvitationInfo, error)
	CreateInvitation(
		scope util.TraceScope,
		currentUserID string,
		args value.CreateInvitationArgs,
	) (value.InvitationInfo, error)
	UpdateInvitation(
		scope util.TraceScope,
		currentUserID string,
		args value.UpdateInvitationArgs,
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

func (ia invitationApplication) ListInvitations(
	scope util.TraceScope,
	currentUserID string,
	args value.ListTeamInvitationArgs,
) ([]value.InvitationInfo, error) {
	const fn = "InvitationApplication.ListInvitations"

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

	// 鉴权：检查当前用户在指定的汉化组是否有权限查看邀请信息
	if !service.CheckInvitationPermission(
		currentUserID,
		args.TeamID,
		adapter.HandleLoadMemberInfo(ia.memberRepository),
		model.PermissionInvitationList,
	) {
		scope.Logger().Warn(fn + ": 权限检查不通过")
		return nil, errors.New("没有权限查看邀请信息")
	}

	// 获取邀请信息列表
	invitationList, err := ia.invitationRepository.List(
		nil,
		query_option.InvitationQuery().FilterByTeamID(args.TeamID),
		query_option.CreatedAtDesc(repository_infra.InvitationTable),
		query_option.Paginate(args.Offset, args.Limit),
	)
	if err != nil {
		scope.Logger().Error(fn+": 获取指定汉化组邀请信息列表失败", zap.Error(err))
		return nil, errors.New("无法获取邀请信息列表")
	}

	// 转换为应用层的值对
	result := make([]value.InvitationInfo, len(invitationList))
	for i, info := range invitationList {
		result[i] = value.NewInvitationInfoFromModel(info)
	}

	return result, nil
}

func (ia invitationApplication) CreateInvitation(
	scope util.TraceScope,
	currentUserID string,
	args value.CreateInvitationArgs,
) (value.InvitationInfo, error) {
	const fn = "InvitationApplication.CreateInvitation"

	if err := args.Validate(); err != nil {
		scope.Logger().Warn(fn+": 参数验证失败", zap.Error(err))
		return value.InvitationInfo{}, errors.New("参数错误: " + err.Error())
	}

	scope.
		WithFields(
			zap.String("current_user_id", currentUserID),
			zap.Any("args", args),
		).
		Logger().
		Debug(fn + ": 被调用")

	// 鉴权：检查当前用户是否有权限创建邀请
	if !service.CheckInvitationPermission(
		currentUserID,
		args.TeamID,
		adapter.HandleLoadMemberInfo(ia.memberRepository),
		model.PermissionInvitationCreate,
	) {
		scope.Logger().Warn(fn + ": 权限检查不通过")
		return value.InvitationInfo{}, errors.New("没有权限创建邀请")
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
		return value.InvitationInfo{}, errors.New("无法检查成员信息")
	}

	if isMemberExisting {
		return value.InvitationInfo{}, errors.New("该用户已经加入该汉化组")
	}

	// 创建邀请信息
	invitationCode, err := service.GenerateInvitationCode()
	if err != nil {
		scope.Logger().Error(fn+": 生成邀请码失败", zap.Error(err))
		return value.InvitationInfo{}, errors.New("生成邀请码失败")
	}

	invitationCreation := model.NewInvitationCreation(
		currentUserID,
		args.TeamID,
		args.InviteeQQ,
		invitationCode,
		model.UnmaskRoles(args.Roles)...,
	)

	// 此处应当使用乐观锁进行更新，不采用事务
	invitationID, err := ia.invitationRepository.Create(nil, *invitationCreation)
	if err != nil {
		scope.Logger().Error(fn+": 创建邀请信息失败", zap.Error(err))
		return value.InvitationInfo{}, errors.New("创建邀请失败")
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

func (ia invitationApplication) UpdateInvitation(
	scope util.TraceScope,
	currentUserID string,
	args value.UpdateInvitationArgs,
) error {
	const fn = "InvitationApplication.UpdateInvitation"

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

	// 鉴权：检查当前用户是否有权限修改邀请
	if !service.CheckInvitationPermission(
		currentUserID,
		args.TeamID,
		adapter.HandleLoadMemberInfo(ia.memberRepository),
		model.PermissionInvitationUpdate,
	) {
		scope.Logger().Warn(fn + ": 权限检查不通过")
		return errors.New("没有权限修改邀请")
	}

	invitationUpdate := model.NewInvitationUpdate(
		args.ID,
		model.UnmaskRoles(args.Roles)...,
	)

	if err := ia.invitationRepository.Update(
		nil,
		*invitationUpdate,
	); err != nil {
		scope.Logger().Error(fn+": 更新邀请信息失败", zap.Error(err))
		return errors.New("更新邀请失败")
	}

	return nil
}

func (ia invitationApplication) DeleteInvitation(
	scope util.TraceScope,
	currentUserID string,
	invitationID string,
) error {
	const fn = "InvitationApplication.DeleteInvitation"

	scope.
		WithFields(
			zap.String("current_user_id", currentUserID),
			zap.String("invitation_id", invitationID),
		).
		Logger().
		Debug(fn + ": 被调用")

	invitation, err := ia.invitationRepository.Get(
		nil,
		query_option.FilterByID(repository_infra.InvitationTable, invitationID),
	)
	if err != nil {
		scope.Logger().Error(fn+": 获取邀请信息失败", zap.Error(err))
		return errors.New("无法获取邀请信息")
	}

	// 鉴权：检查当前用户在邀请所属汉化组是否有删除邀请权限
	if !service.CheckInvitationPermission(
		currentUserID,
		invitation.TeamID,
		adapter.HandleLoadMemberInfo(ia.memberRepository),
		model.PermissionInvitationDelete,
	) {
		scope.Logger().Warn(fn + ": 权限检查不通过")
		return errors.New("没有权限删除邀请")
	}

	// 删除邀请信息
	err = ia.invitationRepository.Delete(nil, invitationID)
	if err != nil {
		scope.Logger().Error(fn+": 删除邀请信息失败", zap.Error(err))
		return errors.New("删除邀请失败")
	}

	return nil
}
