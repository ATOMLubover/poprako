package application

import (
	"errors"

	"labelplus-next-web-be/internal/application/adapter"
	"labelplus-next-web-be/internal/domain/model"
	"labelplus-next-web-be/internal/domain/repository"
	"labelplus-next-web-be/internal/domain/service"
	repository_infra "labelplus-next-web-be/internal/infrastructure/repository"
	"labelplus-next-web-be/internal/infrastructure/repository/query_option"
	"labelplus-next-web-be/internal/util"
	"labelplus-next-web-be/internal/value"

	"go.uber.org/zap"
)

type MemberApplication interface {
	CreateMember(
		scope util.TraceScope,
		currentUserID string,
		args value.CreateMemberArgs,
	) (value.CreateMemberResult, error)
	ListMembers(
		scope util.TraceScope,
		currentUserID string,
		args value.ListTeamMemberArgs,
	) ([]value.MemberProfile, error)
	UpdateMemberRole(
		scope util.TraceScope,
		currentUserID string,
		args value.UpdateMemberRoleArgs,
	) error
	RemoveMember(
		scope util.TraceScope,
		currentUserID string,
		memberID string,
	) error
	JoinTeam(
		scope util.TraceScope,
		currentUserID string,
		args value.JoinTeamArgs,
	) error
}

type memberApplication struct {
	userRepository       repository.UserRepository
	memberRepository     repository.MemberRepository
	invitationRepository repository.InvitationRepository
}

func NewMemberApplication(
	userRepository repository.UserRepository,
	memberRepository repository.MemberRepository,
	invitationRepository repository.InvitationRepository,
) MemberApplication {
	if userRepository == nil || memberRepository == nil || invitationRepository == nil {
		zap.L().Panic(
			"NewMemberApplication: 依赖项不能为空",
			zap.Bool("userRepository_nil", userRepository == nil),
			zap.Bool("memberRepository_nil", memberRepository == nil),
			zap.Bool("invitationRepository_nil", invitationRepository == nil),
		)
	}

	return &memberApplication{
		userRepository:       userRepository,
		memberRepository:     memberRepository,
		invitationRepository: invitationRepository,
	}
}

// 特供超级管理员使用的接口，允许直接创建成员记录
func (ma *memberApplication) CreateMember(
	scope util.TraceScope,
	currentUserID string,
	args value.CreateMemberArgs,
) (value.CreateMemberResult, error) {
	const fn = "MemberApplication.CreateMember"

	if err := args.Validate(); err != nil {
		scope.Logger().Warn(fn+": 参数验证失败", zap.Error(err))
		return value.CreateMemberResult{}, errors.New("参数错误: " + err.Error())
	}

	scope.
		WithFields(
			zap.String("current_user_id", currentUserID),
			zap.Any("args", args),
		).
		Logger().
		Debug(fn + ": 被调用")

	if !service.CheckMemberPermission(
		currentUserID,
		args.TeamID,
		adapter.HandleLoadUserInfo(ma.userRepository),
		adapter.HandleLoadMemberInfo(ma.memberRepository),
		model.PermissionMemberCreate,
	) {
		scope.Logger().Warn(fn + ": 权限检查失败")
		return value.CreateMemberResult{}, errors.New("权限不足")
	}

	isMemberExisting, err := ma.memberRepository.Exist(
		nil,
		query_option.MemberQuery().FilterByTeamID(args.TeamID),
		query_option.MemberQuery().FilterByUserID(args.UserID),
	)
	if err != nil {
		scope.Logger().Error(fn+": 检查成员信息失败", zap.Error(err))
		return value.CreateMemberResult{}, errors.New("无法检查成员信息")
	}

	if isMemberExisting {
		return value.CreateMemberResult{}, errors.New("该用户已经加入该汉化组")
	}

	memberID, err := ma.memberRepository.Create(
		nil,
		model.NewMemberCreation(args.UserID, args.TeamID, model.UnmaskRoles(args.Roles)...),
	)
	if err != nil {
		scope.Logger().Error(fn+": 创建成员失败", zap.Error(err))
		return value.CreateMemberResult{}, errors.New("创建成员失败")
	}

	result := value.NewCreateMemberResult(memberID)

	return result, nil
}

func (ma *memberApplication) ListMembers(
	scope util.TraceScope,
	currentUserID string,
	args value.ListTeamMemberArgs,
) ([]value.MemberProfile, error) {
	const fn = "MemberApplication.ListMembers"

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

	// 鉴权：检查当前用户在指定汉化组是否有查看成员列表权限
	if !service.CheckMemberPermission(
		currentUserID,
		args.TeamID,
		adapter.HandleLoadUserInfo(ma.userRepository),
		adapter.HandleLoadMemberInfo(ma.memberRepository),
		model.PermissionMemberList,
	) {
		scope.Logger().Warn(fn + ": 权限检查失败")
		return nil, errors.New("权限不足")
	}

	// 获取成员列表（含用户信息）
	memberList, err := ma.memberRepository.ListProfiles(
		nil,
		query_option.CreatedAtDesc(repository_infra.MemberTable),
		query_option.MemberQuery().FilterByTeamID(args.TeamID),
		query_option.Paginate(args.Offset, args.Limit),
	)
	if err != nil {
		scope.Logger().Error(fn+": 获取成员列表失败", zap.Error(err))
		return nil, errors.New("无法获取成员列表")
	}

	result := make([]value.MemberProfile, len(memberList))

	for i := range memberList {
		result[i] = value.NewMemberProfileFromModel(memberList[i])
	}

	return result, nil
}

func (ma *memberApplication) UpdateMemberRole(
	scope util.TraceScope,
	currentUserID string,
	args value.UpdateMemberRoleArgs,
) error {
	const fn = "MemberApplication.UpdateMemberRole"

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

	// 查询目标成员信息，获取可信的 TeamID
	targetMember, err := ma.memberRepository.GetProfile(
		nil,
		query_option.FilterByID(repository_infra.MemberTable, args.ID),
	)
	if err != nil {
		scope.Logger().Error(fn+": 获取目标成员信息失败", zap.Error(err))
		return errors.New("无法获取成员信息")
	}

	// 鉴权：使用从数据库查询到的可信 TeamID 进行权限检查
	if !service.CheckMemberPermission(
		currentUserID,
		targetMember.TeamID,
		adapter.HandleLoadUserInfo(ma.userRepository),
		adapter.HandleLoadMemberInfo(ma.memberRepository),
		model.PermissionMemberUpdate,
	) {
		scope.Logger().Warn(fn + ": 权限检查失败")
		return errors.New("权限不足")
	}

	// 构建更新对象，传入目标角色
	memberUpdate := model.NewMemberUpdate(args.ID, model.UnmaskRolesWithTime(args.Roles)...)

	if err := ma.memberRepository.Update(nil, memberUpdate); err != nil {
		scope.Logger().Error(fn+": 更新成员角色失败", zap.Error(err))
		return errors.New("更新成员角色失败")
	}

	return nil
}

func (ma *memberApplication) RemoveMember(
	scope util.TraceScope,
	currentUserID string,
	memberID string,
) error {
	const fn = "MemberApplication.RemoveMember"

	scope.
		WithFields(
			zap.String("current_user_id", currentUserID),
			zap.String("member_id", memberID),
		).
		Logger().
		Debug(fn + ": 被调用")

	// 查询目标成员信息，获取可信的 TeamID
	targetMember, err := ma.memberRepository.GetProfile(
		nil,
		query_option.FilterByID(repository_infra.MemberTable, memberID),
	)
	if err != nil {
		scope.Logger().Error(fn+": 获取目标成员信息失败", zap.Error(err))
		return errors.New("无法获取成员信息")
	}

	// 鉴权：使用从数据库查询到的可信 TeamID 进行权限检查
	if !service.CheckMemberPermission(
		currentUserID,
		targetMember.TeamID,
		adapter.HandleLoadUserInfo(ma.userRepository),
		adapter.HandleLoadMemberInfo(ma.memberRepository),
		model.PermissionMemberDelete,
	) {
		scope.Logger().Warn(fn + ": 权限检查失败")
		return errors.New("权限不足")
	}

	// 删除成员，不存在情况由 repo 层处理
	if err := ma.memberRepository.Delete(nil, memberID); err != nil {
		scope.Logger().Error(fn+": 删除成员失败", zap.Error(err))
		return errors.New("删除成员失败")
	}

	return nil
}

func (ma *memberApplication) JoinTeam(
	scope util.TraceScope,
	currentUserID string,
	args value.JoinTeamArgs,
) error {
	const fn = "MemberApplication.JoinTeam"

	if err := args.Validate(); err != nil {
		scope.Logger().Warn(fn+": 参数验证失败", zap.Error(err))
		return errors.New("参数错误: " + err.Error())
	}

	scope.
		WithFields(
			zap.String("current_user_id", currentUserID),
			zap.String("invitation_code", args.InvitationCode),
		).
		Logger().
		Debug(fn + ": 被调用")

	// 获取当前用户信息，提取 QQ
	currentUser, err := ma.userRepository.GetByID(
		nil,
		query_option.FilterByID(repository_infra.UserTable, currentUserID),
	)
	if err != nil || currentUser.ID == "" {
		scope.Logger().Error(fn+": 获取当前用户信息失败", zap.Error(err))
		return errors.New("无法获取用户信息")
	}

	// 根据 QQ 和邀请码查找待消耗的邀请
	invitationInfo, err := ma.invitationRepository.Get(
		nil,
		query_option.InvitationQuery().FilterByInviteeQQ(currentUser.QQ),
		query_option.InvitationQuery().FilterByCode(args.InvitationCode),
		query_option.InvitationQuery().FilterPending(true),
	)
	if err != nil && !errors.Is(err, repository_infra.ErrRecordNotFound) {
		scope.Logger().Error(fn+": 查询邀请信息失败", zap.Error(err))
		return errors.New("获取邀请信息失败")
	}

	if err != nil {
		scope.Logger().Warn(fn + ": 没有找到对应的邀请")
		return errors.New("邀请码无效或已被使用")
	}

	// 检查用户是否已经是该汉化组成员
	isMember, err := ma.memberRepository.Exist(
		nil,
		query_option.MemberQuery().FilterByTeamID(invitationInfo.TeamID),
		query_option.MemberQuery().FilterByUserID(currentUserID),
	)
	if err != nil {
		scope.Logger().Error(fn+": 检查成员信息失败", zap.Error(err))
		return errors.New("无法检查成员信息")
	}

	if isMember {
		return errors.New("您已经是该汉化组的成员")
	}

	// 在事务中创建成员并使邀请失效
	transactionExecutor := ma.memberRepository.BeginTransaction()
	if transactionExecutor.Error != nil {
		scope.Logger().Error(fn+": 创建事务失败", zap.Error(transactionExecutor.Error))
		return errors.New("加入汉化组失败")
	}

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

	memberCreation := model.NewMemberCreation(
		currentUserID,
		invitationInfo.TeamID,
		model.UnmaskRoles(invitationInfo.RoleMask())...,
	)

	_, transactionErr = ma.memberRepository.Create(transactionExecutor, memberCreation)
	if transactionErr != nil {
		scope.Logger().Error(fn+": 创建成员失败", zap.Error(transactionErr))
		return errors.New("加入汉化组失败")
	}

	transactionErr = ma.invitationRepository.Invalidate(transactionExecutor, invitationInfo.ID)
	if transactionErr != nil {
		scope.Logger().Error(fn+": 使邀请失效失败", zap.Error(transactionErr))
		return errors.New("加入汉化组失败")
	}

	if commitErr := transactionExecutor.Commit().Error; commitErr != nil {
		scope.Logger().Error(fn+": 事务提交失败", zap.Error(commitErr))
		return errors.New("加入汉化组失败")
	}

	return nil
}
