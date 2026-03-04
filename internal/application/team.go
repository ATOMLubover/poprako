package application

import (
	"errors"
	"time"

	"labelplus-next-web-be/internal/domain/model"
	"labelplus-next-web-be/internal/domain/repository"
	"labelplus-next-web-be/internal/domain/service"
	"labelplus-next-web-be/internal/repository/query_option"
	"labelplus-next-web-be/internal/util"
	"labelplus-next-web-be/internal/value"

	"go.uber.org/zap"
)

type TeamApplication interface {
	CreateTeam(
		scope util.TraceScope,
		currentUserID string,
		args *value.CreateTeamArgs,
	) (*value.TeamInfo, error)
	ListAllTeams(
		scope util.TraceScope,
		currentUserID string,
	) ([]*value.TeamInfo, error)
	ListMyTeams(
		scope util.TraceScope,
		currentUserID string,
	) ([]*value.TeamInfo, error)
	UpdateTeam(
		scope util.TraceScope,
		currentUserID string,
		args *value.UpdateTeamArgs,
	) error
	DeleteTeam(
		scope util.TraceScope,
		currentUserID string,
		teamID string,
	) error
}

type teamApplication struct {
	userRepository   repository.UserRepository
	teamRepository   repository.TeamRepository
	memberRepository repository.MemberRepository
}

func NewTeamApplication(
	userRepository repository.UserRepository,
	teamRepository repository.TeamRepository,
	memberRepository repository.MemberRepository,
) TeamApplication {
	if userRepository == nil ||
		teamRepository == nil ||
		memberRepository == nil {
		zap.L().Panic(
			"NewTeamApplication: 依赖项不能为空",
			zap.Bool("userRepository_nil", userRepository == nil),
			zap.Bool("teamRepository_nil", teamRepository == nil),
			zap.Bool("memberRepository_nil", memberRepository == nil),
		)
	}

	return &teamApplication{
		userRepository:   userRepository,
		teamRepository:   teamRepository,
		memberRepository: memberRepository,
	}
}

func (ta *teamApplication) CreateTeam(
	scope util.TraceScope,
	currentUserID string,
	args *value.CreateTeamArgs,
) (*value.TeamInfo, error) {
	const fn = "TeamApplication.CreateTeam"

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

	// 鉴权：获取当前用户信息
	currentUser, err := ta.userRepository.GetInfoByID(nil, currentUserID)
	if err != nil {
		scope.Logger().Error(fn+": 获取当前用户信息失败", zap.Error(err))
		return nil, errors.New("无法获取用户信息")
	}

	// 鉴权：检查当前用户是否有权限创建汉化组（仅超级管理员）
	if !service.CheckTeamPermission(
		"",
		currentUser,
		nil,
		model.PermissionTeamCreate,
	) {
		return nil, errors.New("没有权限创建汉化组")
	}

	// 创建汉化组
	teamCreation := model.NewTeamCreation(args.Name, args.Description)

	teamID, err := ta.teamRepository.Create(nil, teamCreation)
	if err != nil {
		scope.Logger().Error(fn+": 创建汉化组失败", zap.Error(err))
		return nil, errors.New("创建汉化组失败")
	}

	// 构造返回值，不再查询数据库
	now := time.Now()
	return value.NewTeamInfoFromModel(&model.TeamInfo{
		ID:          teamID,
		Name:        args.Name,
		Description: args.Description,
		CreatedAt:   now,
		UpdatedAt:   now,
	}), nil
}

func (ta *teamApplication) ListAllTeams(
	scope util.TraceScope,
	currentUserID string,
) ([]*value.TeamInfo, error) {
	const fn = "TeamApplication.ListAllTeams"

	if currentUserID == "" {
		scope.Logger().Warn(fn + ": currentUserID 为空")
		return nil, errors.New(ErrInternalError)
	}

	scope.WithFields(
		zap.String("current_user_id", currentUserID),
	)

	scope.Logger().Debug(fn + ": 被调用")

	// 鉴权：获取当前用户信息
	currentUser, err := ta.userRepository.GetInfoByID(nil, currentUserID)
	if err != nil {
		scope.Logger().Error(fn+": 获取当前用户信息失败", zap.Error(err))
		return nil, errors.New("无法获取用户信息")
	}

	// 鉴权：检查当前用户是否有权限查看所有汉化组（仅超级管理员）
	if !service.CheckTeamPermission(
		"",
		currentUser,
		nil,
		model.PermissionTeamListAll,
	) {
		return nil, errors.New("没有权限查看所有汉化组")
	}

	// 获取所有汉化组列表
	teamList, err := ta.teamRepository.List(nil)
	if err != nil {
		scope.Logger().Error(fn+": 获取汉化组列表失败", zap.Error(err))
		return nil, errors.New("无法获取汉化组列表")
	}

	result := make([]*value.TeamInfo, len(teamList))

	for i, team := range teamList {
		result[i] = value.NewTeamInfoFromModel(team)
	}

	return result, nil
}

func (ta *teamApplication) ListMyTeams(
	scope util.TraceScope,
	currentUserID string,
) ([]*value.TeamInfo, error) {
	const fn = "TeamApplication.ListMyTeams"

	if currentUserID == "" {
		scope.Logger().Warn(fn + ": currentUserID 为空")
		return nil, errors.New(ErrInternalError)
	}

	scope.WithFields(
		zap.String("current_user_id", currentUserID),
	)

	scope.Logger().Debug(fn + ": 被调用")

	// 鉴权：检查当前用户是否有权限查看所在汉化组列表（所有成员均可）
	if !service.CheckTeamPermission(
		"",
		nil,
		nil,
		model.PermissionTeamListMine,
	) {
		return nil, errors.New("没有权限查看所在汉化组列表")
	}

	// 获取当前用户的所有成员记录
	currentUserMemberships, err := ta.memberRepository.List(
		nil,
		query_option.MemberQuery().FilterByUserID(currentUserID),
	)
	if err != nil {
		scope.Logger().Error(fn+": 获取当前用户的成员信息失败", zap.Error(err))
		return nil, errors.New("无法获取成员信息")
	}

	if len(currentUserMemberships) == 0 {
		return nil, nil
	}

	// 构建用户所在汉化组 ID 的数组
	teamIDs := make([]string, len(currentUserMemberships))
	for i, membership := range currentUserMemberships {
		teamIDs[i] = membership.TeamID
	}

	// 由于用户所在汉化组数量较少，直接采用 N + 1 查询方式获取汉化组信息
	teams, err := ta.teamRepository.List(
		nil,
		query_option.TeamQuery().FilterByIDs(teamIDs),
	)
	if err != nil {
		scope.Logger().Error(fn+": 获取汉化组列表失败", zap.Error(err))
		return nil, errors.New("无法获取汉化组列表")
	}

	result := make([]*value.TeamInfo, len(teams))

	for i, team := range teams {
		result[i] = value.NewTeamInfoFromModel(team)
	}

	return result, nil
}

func (ta *teamApplication) UpdateTeam(
	scope util.TraceScope,
	currentUserID string,
	args *value.UpdateTeamArgs,
) error {
	const fn = "TeamApplication.UpdateTeam"

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

	// 鉴权：获取当前用户信息和成员信息
	currentUser, err := ta.userRepository.GetInfoByID(nil, currentUserID)
	if err != nil {
		scope.Logger().Error(fn+": 获取当前用户信息失败", zap.Error(err))
		return errors.New("无法获取用户信息")
	}

	currentUserMemberships, err := ta.memberRepository.List(
		nil,
		query_option.MemberQuery().FilterByUserID(currentUserID),
	)
	if err != nil {
		scope.Logger().Error(fn+": 获取当前用户的成员信息失败", zap.Error(err))
		return errors.New("无法获取成员信息")
	}

	// 鉴权：检查当前用户是否有权限更新汉化组（超级管理员或团队管理员）
	if !service.CheckTeamPermission(
		args.ID,
		currentUser,
		currentUserMemberships,
		model.PermissionTeamUpdate,
	) {
		return errors.New("没有权限更新汉化组")
	}

	// 构建更新对象，仅传入有值的字段
	teamUpdate := model.NewTeamUpdate(args.ID)
	teamUpdate.Name = args.Name
	teamUpdate.Description = args.Description

	if err := ta.teamRepository.Update(nil, teamUpdate); err != nil {
		scope.Logger().Error(fn+": 更新汉化组信息失败", zap.Error(err))
		return errors.New("更新汉化组失败")
	}

	return nil
}

func (ta *teamApplication) DeleteTeam(
	scope util.TraceScope,
	currentUserID string,
	teamID string,
) error {
	const fn = "TeamApplication.DeleteTeam"

	if currentUserID == "" {
		scope.Logger().Warn(fn + ": currentUserID 为空")
		return errors.New(ErrInternalError)
	}

	scope.WithFields(
		zap.String("current_user_id", currentUserID),
		zap.String("team_id", teamID),
	)

	scope.Logger().Debug(fn + ": 被调用")

	// 鉴权：获取当前用户信息和成员信息
	currentUser, err := ta.userRepository.GetInfoByID(nil, currentUserID)
	if err != nil {
		scope.Logger().Error(fn+": 获取当前用户信息失败", zap.Error(err))
		return errors.New("无法获取用户信息")
	}

	currentUserMemberships, err := ta.memberRepository.List(
		nil,
		query_option.MemberQuery().FilterByUserID(currentUserID),
	)
	if err != nil {
		scope.Logger().Error(fn+": 获取当前用户的成员信息失败", zap.Error(err))
		return errors.New("无法获取成员信息")
	}

	// 鉴权：检查当前用户是否有权限删除汉化组（超级管理员或团队管理员）
	if !service.CheckTeamPermission(
		teamID,
		currentUser,
		currentUserMemberships,
		model.PermissionTeamDelete,
	) {
		return errors.New("没有权限删除汉化组")
	}

	// 删除汉化组
	if err := ta.teamRepository.DeleteByID(nil, teamID); err != nil {
		scope.Logger().Error(fn+": 删除汉化组失败", zap.Error(err))
		return errors.New("删除汉化组失败")
	}

	return nil
}
