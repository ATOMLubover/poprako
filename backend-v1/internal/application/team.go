package application

import (
	"errors"

	"labelplus-next-web-be/internal/application/adapter"
	"labelplus-next-web-be/internal/application/assembler"
	"labelplus-next-web-be/internal/domain/external"
	"labelplus-next-web-be/internal/domain/model"
	"labelplus-next-web-be/internal/domain/repository"
	"labelplus-next-web-be/internal/domain/service"
	repository_infra "labelplus-next-web-be/internal/infrastructure/repository"
	"labelplus-next-web-be/internal/infrastructure/repository/query_option"
	"labelplus-next-web-be/internal/util"
	"labelplus-next-web-be/internal/value"

	"go.uber.org/zap"
)

type TeamApplication interface {
	CreateTeam(
		scope util.TraceScope,
		currentUserID string,
		args value.CreateTeamArgs,
	) (value.CreateTeamResult, error)
	ListTeams(
		scope util.TraceScope,
		currentUserID string,
		args value.ListTeamArgs,
	) ([]value.TeamInfo, error)
	ListMyTeams(
		scope util.TraceScope,
		currentUserID string,
		args value.ListMyTeamArgs,
	) ([]value.TeamInfo, error)
	ReserveTeamAvatar(
		scope util.TraceScope,
		currentUserID string,
		teamID string,
	) (value.ReserveTeamAvatarResult, error)
	ConfirmTeamAvatarUploaded(
		scope util.TraceScope,
		currentUserID string,
		teamID string,
	) error
	UpdateTeam(
		scope util.TraceScope,
		currentUserID string,
		args value.UpdateTeamArgs,
	) error
	RemoveTeam(
		scope util.TraceScope,
		currentUserID string,
		teamID string,
	) error
}

type teamApplication struct {
	ossClient        external.OSSClient
	userRepository   repository.UserRepository
	teamRepository   repository.TeamRepository
	memberRepository repository.MemberRepository
}

func NewTeamApplication(
	ossClient external.OSSClient,
	userRepository repository.UserRepository,
	teamRepository repository.TeamRepository,
	memberRepository repository.MemberRepository,
) TeamApplication {
	if ossClient == nil ||
		userRepository == nil ||
		teamRepository == nil ||
		memberRepository == nil {
		zap.L().Panic(
			"NewTeamApplication: 依赖项不能为空",
			zap.Bool("ossClient_nil", ossClient == nil),
			zap.Bool("userRepository_nil", userRepository == nil),
			zap.Bool("teamRepository_nil", teamRepository == nil),
			zap.Bool("memberRepository_nil", memberRepository == nil),
		)
	}

	return &teamApplication{
		ossClient:        ossClient,
		userRepository:   userRepository,
		teamRepository:   teamRepository,
		memberRepository: memberRepository,
	}
}

func (ta *teamApplication) CreateTeam(
	scope util.TraceScope,
	currentUserID string,
	args value.CreateTeamArgs,
) (value.CreateTeamResult, error) {
	const fn = "TeamApplication.CreateTeam"

	if err := args.Validate(); err != nil {
		scope.Logger().Warn(fn+": 参数验证失败", zap.Error(err))
		return value.CreateTeamResult{}, errors.New("参数错误: " + err.Error())
	}

	scope.
		WithFields(
			zap.String("current_user_id", currentUserID),
			zap.Any("args", args),
		).
		Logger().
		Debug(fn + ": 被调用")

	if !model.PermTeamCreate().Check(
		currentUserID,
		adapter.HandleLoadUserInfo(ta.userRepository),
	) {
		scope.Logger().Warn(fn + ": 权限检查失败")
		return value.CreateTeamResult{}, errors.New("权限不足")
	}

	// 创建汉化组
	teamCreation := model.NewTeamCreation(args.Name, args.Description)

	teamID, err := ta.teamRepository.Create(nil, *teamCreation)
	if err != nil {
		scope.Logger().Error(fn+": 创建汉化组失败", zap.Error(err))
		return value.CreateTeamResult{}, errors.New("创建汉化组失败")
	}

	return value.CreateTeamResult{ID: teamID}, nil
}

func (ta *teamApplication) ListTeams(
	scope util.TraceScope,
	currentUserID string,
	args value.ListTeamArgs,
) ([]value.TeamInfo, error) {
	const fn = "TeamApplication.ListTeams"

	if err := args.Validate(); err != nil {
		scope.Logger().Warn(fn+": 参数验证失败", zap.Error(err))
		return nil, errors.New("参数错误: " + err.Error())
	}

	scope.
		WithFields(
			zap.String("current_user_id", currentUserID),
		).
		Logger().
		Debug(fn + ": 被调用")

	if !model.PermTeamListAll().Check(
		currentUserID,
		adapter.HandleLoadUserInfo(ta.userRepository),
	) {
		scope.Logger().Warn(fn + ": 权限检查失败")
		return nil, errors.New("权限不足")
	}

	// 获取所有汉化组列表
	teamList, err := ta.teamRepository.List(
		nil,
		query_option.UpdatedAtDesc(repository_infra.TeamTable),
		query_option.Paginate(args.Offset, args.Limit),
	)
	if err != nil {
		scope.Logger().Error(fn+": 获取汉化组列表失败", zap.Error(err))
		return nil, errors.New("无法获取汉化组列表")
	}

	result := make([]value.TeamInfo, len(teamList))

	for i, team := range teamList {
		avatarURL, err := ta.ossClient.GenerateGetPresignedURL(team.AvatarOSSKey)
		if err != nil {
			scope.Logger().Error(fn+": 生成头像访问链接失败", zap.Error(err))
			return nil, errors.New("无法获取汉化组列表")
		}

		_ = avatarURL
		result[i] = assembler.AssembleTeamInfo(team, ta.ossClient.GenerateGetPresignedURL)
	}

	return result, nil
}

func (ta *teamApplication) ReserveTeamAvatar(
	scope util.TraceScope,
	currentUserID string,
	teamID string,
) (value.ReserveTeamAvatarResult, error) {
	const fn = "TeamApplication.ReserveTeamAvatar"

	scope.
		WithFields(
			zap.String("current_user_id", currentUserID),
			zap.String("team_id", teamID),
		).
		Logger().
		Debug(fn + ": 被调用")

	if teamID == "" {
		return value.ReserveTeamAvatarResult{}, errors.New("汉化组 ID 不能为空")
	}

	if !model.PermTeamUpdate().Check(
		currentUserID,
		teamID,
		adapter.HandleLoadUserInfo(ta.userRepository),
		adapter.HandleLoadMemberInfo(ta.memberRepository),
	) {
		scope.Logger().Warn(fn + ": 权限检查失败")
		return value.ReserveTeamAvatarResult{}, errors.New("权限不足")
	}

	avatarOSSKey := service.GenerateTeamAvatarOSSKey(teamID)

	putURL, err := ta.ossClient.GeneratePutPresignedURL(avatarOSSKey)
	if err != nil {
		scope.Logger().Error(fn+": 生成头像访问链接失败", zap.Error(err))
		return value.ReserveTeamAvatarResult{}, errors.New("预留汉化组头像失败")
	}

	if err := ta.teamRepository.ReserveAvatar(nil, teamID, avatarOSSKey); err != nil {
		scope.Logger().Error(fn+": 写入头像 OSS Key 失败", zap.Error(err))
		return value.ReserveTeamAvatarResult{}, errors.New("预留汉化组头像失败")
	}

	return value.ReserveTeamAvatarResult{
		AvatarOSSKey: avatarOSSKey,
		PutURL:       putURL,
	}, nil
}

func (ta *teamApplication) ListMyTeams(
	scope util.TraceScope,
	currentUserID string,
	args value.ListMyTeamArgs,
) ([]value.TeamInfo, error) {
	const fn = "TeamApplication.ListMyTeams"

	if err := args.Validate(); err != nil {
		scope.Logger().Warn(fn+": 参数验证失败", zap.Error(err))
		return nil, errors.New("参数错误: " + err.Error())
	}

	scope.
		WithFields(
			zap.String("current_user_id", currentUserID),
		).
		Logger().
		Debug(fn + ": 被调用")

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
		query_option.UpdatedAtDesc(repository_infra.TeamTable),
		query_option.FilterByIDs(repository_infra.TeamTable, teamIDs),
		query_option.Paginate(args.Offset, args.Limit),
	)
	if err != nil {
		scope.Logger().Error(fn+": 获取汉化组列表失败", zap.Error(err))
		return nil, errors.New("无法获取汉化组列表")
	}

	result := make([]value.TeamInfo, len(teams))

	for i, team := range teams {
		avatarURL, err := ta.ossClient.GenerateGetPresignedURL(team.AvatarOSSKey)
		if err != nil {
			scope.Logger().Error(fn+": 生成头像访问链接失败", zap.Error(err))
			return nil, errors.New("无法获取汉化组列表")
		}

		_ = avatarURL
		result[i] = assembler.AssembleTeamInfo(team, ta.ossClient.GenerateGetPresignedURL)
	}

	return result, nil
}

func (ta *teamApplication) UpdateTeam(
	scope util.TraceScope,
	currentUserID string,
	args value.UpdateTeamArgs,
) error {
	const fn = "TeamApplication.UpdateTeam"

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

	if !model.PermTeamUpdate().Check(
		currentUserID,
		args.ID,
		adapter.HandleLoadUserInfo(ta.userRepository),
		adapter.HandleLoadMemberInfo(ta.memberRepository),
	) {
		scope.Logger().Warn(fn + ": 权限检查失败")
		return errors.New("权限不足")
	}

	teamUpdate := model.NewTeamUpdate(
		args.ID,
		args.Name,
		args.Description,
	)

	if err := ta.teamRepository.Update(nil, teamUpdate); err != nil {
		scope.Logger().Error(fn+": 更新汉化组信息失败", zap.Error(err))
		return errors.New("更新汉化组失败")
	}

	return nil
}

func (ta *teamApplication) ConfirmTeamAvatarUploaded(
	scope util.TraceScope,
	currentUserID string,
	teamID string,
) error {
	const fn = "TeamApplication.ConfirmTeamAvatarUploaded"

	scope.
		WithFields(
			zap.String("current_user_id", currentUserID),
			zap.String("team_id", teamID),
		).
		Logger().
		Debug(fn + ": 被调用")

	if teamID == "" {
		return errors.New("汉化组 ID 不能为空")
	}

	if !model.PermTeamUpdate().Check(
		currentUserID,
		teamID,
		adapter.HandleLoadUserInfo(ta.userRepository),
		adapter.HandleLoadMemberInfo(ta.memberRepository),
	) {
		scope.Logger().Warn(fn + ": 权限检查失败")
		return errors.New("权限不足")
	}

	if err := ta.teamRepository.ConfirmAvatarUploaded(nil, teamID); err != nil {
		scope.Logger().Error(fn+": 确认汉化组头像上传失败", zap.Error(err))
		return errors.New("确认汉化组头像上传失败")
	}

	return nil
}

func (ta *teamApplication) RemoveTeam(
	scope util.TraceScope,
	currentUserID string,
	teamID string,
) error {
	const fn = "TeamApplication.DeleteTeam"

	scope.
		WithFields(
			zap.String("current_user_id", currentUserID),
			zap.String("team_id", teamID),
		).
		Logger().
		Debug(fn + ": 被调用")

	if !model.PermTeamRemove().Check(
		currentUserID,
		teamID,
		adapter.HandleLoadUserInfo(ta.userRepository),
		adapter.HandleLoadMemberInfo(ta.memberRepository),
	) {
		scope.Logger().Warn(fn + ": 权限检查失败")
		return errors.New("权限不足")
	}

	// 删除汉化组
	if err := ta.teamRepository.Delete(nil, teamID); err != nil {
		scope.Logger().Error(fn+": 删除汉化组失败", zap.Error(err))
		return errors.New("删除汉化组失败")
	}

	return nil
}
