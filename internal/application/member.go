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

type MemberApplication interface {
	ListMembers(
		scope util.TraceScope,
		currentUserID string,
		teamID string,
		paginationParams value.PaginationParams,
	) ([]*value.MemberProfile, error)
	UpdateMemberRole(
		scope util.TraceScope,
		currentUserID string,
		args *value.UpdateMemberRoleArgs,
	) error
	RemoveMember(
		scope util.TraceScope,
		currentUserID string,
		memberID string,
	) error
}

type memberApplication struct {
	memberRepository repository.MemberRepository
}

func NewMemberApplication(
	memberRepository repository.MemberRepository,
) MemberApplication {
	if memberRepository == nil {
		zap.L().Panic(
			"NewMemberApplication: 依赖项不能为空",
			zap.Bool("memberRepository_nil", memberRepository == nil),
		)
	}

	return &memberApplication{
		memberRepository: memberRepository,
	}
}

func (ma *memberApplication) ListMembers(
	scope util.TraceScope,
	currentUserID string,
	teamID string,
	paginationParams value.PaginationParams,
) ([]*value.MemberProfile, error) {
	const fn = "MemberApplication.ListMembers"

	if currentUserID == "" {
		scope.Logger().Warn(fn + ": currentUserID 为空")
		return nil, errors.New(ErrInternalError)
	}

	scope.WithFields(
		zap.String("current_user_id", currentUserID),
		zap.String("team_id", teamID),
	)

	scope.Logger().Debug(fn + ": 被调用")

	// 鉴权：获取当前用户在各汉化组的成员信息
	currentUserMemberships, err := ma.memberRepository.List(
		nil,
		query_option.MemberQuery().FilterByUserID(currentUserID),
	)
	if err != nil {
		scope.Logger().Error(fn+": 获取当前用户的成员信息失败", zap.Error(err))
		return nil, errors.New("无法获取成员信息")
	}

	// 鉴权：检查当前用户在指定汉化组是否有权限查看成员列表
	if !service.CheckMemberPermission(
		teamID,
		currentUserMemberships,
		model.PermissionMemberList,
	) {
		return nil, errors.New("没有权限查看成员列表")
	}

	// 获取成员列表（含用户信息）
	memberList, err := ma.memberRepository.ListWithUserInfo(
		nil,
		query_option.MemberQuery().FilterByTeamID(teamID),
		query_option.Paginate(paginationParams.Offset, paginationParams.Limit),
	)
	if err != nil {
		scope.Logger().Error(fn+": 获取成员列表失败", zap.Error(err))
		return nil, errors.New("无法获取成员列表")
	}

	result := make([]*value.MemberProfile, len(memberList))

	for i := range memberList {
		result[i] = value.NewMemberProfileFromModel(&memberList[i])
	}

	return result, nil
}

func (ma *memberApplication) UpdateMemberRole(
	scope util.TraceScope,
	currentUserID string,
	args *value.UpdateMemberRoleArgs,
) error {
	const fn = "MemberApplication.UpdateMemberRole"

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

	// 查询目标成员信息，获取可信的 TeamID
	targetMember, err := ma.memberRepository.GetByID(nil, args.ID)
	if err != nil {
		scope.Logger().Error(fn+": 获取目标成员信息失败", zap.Error(err))
		return errors.New("无法获取成员信息")
	}

	// 鉴权：获取当前用户在各汉化组的成员信息
	currentUserMemberships, err := ma.memberRepository.List(
		nil,
		query_option.MemberQuery().FilterByUserID(currentUserID),
	)
	if err != nil {
		scope.Logger().Error(fn+": 获取当前用户的成员信息失败", zap.Error(err))
		return errors.New("无法获取成员信息")
	}

	// 鉴权：使用从数据库查询到的可信 TeamID 进行权限检查
	if !service.CheckMemberPermission(
		targetMember.TeamID,
		currentUserMemberships,
		model.PermissionMemberUpdate,
	) {
		return errors.New("没有权限更新成员角色")
	}

	// 构建更新对象，传入目标角色
	memberUpdate := model.NewMemberUpdate(args.ID, model.UnmaskRoles(args.Roles)...)

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

	if currentUserID == "" {
		scope.Logger().Warn(fn + ": currentUserID 为空")
		return errors.New(ErrInternalError)
	}

	scope.WithFields(
		zap.String("current_user_id", currentUserID),
		zap.String("member_id", memberID),
	)

	scope.Logger().Debug(fn + ": 被调用")

	// 查询目标成员信息，获取可信的 TeamID
	targetMember, err := ma.memberRepository.GetByID(nil, memberID)
	if err != nil {
		scope.Logger().Error(fn+": 获取目标成员信息失败", zap.Error(err))
		return errors.New("无法获取成员信息")
	}

	// 鉴权：获取当前用户在各汉化组的成员信息
	currentUserMemberships, err := ma.memberRepository.List(
		nil,
		query_option.MemberQuery().FilterByUserID(currentUserID),
	)
	if err != nil {
		scope.Logger().Error(fn+": 获取当前用户的成员信息失败", zap.Error(err))
		return errors.New("无法获取成员信息")
	}

	// 鉴权：使用从数据库查询到的可信 TeamID 进行权限检查
	if !service.CheckMemberPermission(
		targetMember.TeamID,
		currentUserMemberships,
		model.PermissionMemberDelete,
	) {
		return errors.New("没有权限删除成员")
	}

	// 删除成员，不存在情况由 repo 层处理
	if err := ma.memberRepository.DeleteByID(nil, memberID); err != nil {
		scope.Logger().Error(fn+": 删除成员失败", zap.Error(err))
		return errors.New("删除成员失败")
	}

	return nil
}

