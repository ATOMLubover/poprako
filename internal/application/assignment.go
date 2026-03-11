package application

import (
	"errors"

	"labelplus-next-web-be/internal/application/adapter"
	"labelplus-next-web-be/internal/domain/external"
	"labelplus-next-web-be/internal/domain/model"
	"labelplus-next-web-be/internal/domain/repository"
	repository_infra "labelplus-next-web-be/internal/infrastructure/repository"
	"labelplus-next-web-be/internal/infrastructure/repository/query_option"
	"labelplus-next-web-be/internal/util"
	"labelplus-next-web-be/internal/value"

	"go.uber.org/zap"
)

type AssignmentApplication interface {
	ListChapterAssignments(
		scope util.TraceScope,
		currentUserID string,
		args value.ListChapterAssignmentArgs,
	) ([]value.AssignmentWithUserInfo, error)
	ListMyAssignments(
		scope util.TraceScope,
		currentUserID string,
		args value.ListUserAssignmentArgs,
	) ([]value.AssignmentWithChapterInfo, error)
	CreateChapterAssignment(
		scope util.TraceScope,
		currentUserID string,
		args value.CreateChapterAssignmentArgs,
	) (value.CreateChapterAssignmentResult, error)
	UpdateAssignment(
		scope util.TraceScope,
		currentUserID string,
		args value.UpdateAssignmentArgs,
	) error
	RemoveAssignment(
		scope util.TraceScope,
		currentUserID string,
		assignmentID string,
	) error
}

type assignmentApplication struct {
	ossClient            external.OSSClient
	memberRepository     repository.MemberRepository
	comicRepository      repository.ComicRepository
	chapterRepository    repository.ChapterRepository
	assignmentRepository repository.AssignmentRepository
}

func NewAssignmentApplication(
	ossClient external.OSSClient,
	memberRepository repository.MemberRepository,
	comicRepository repository.ComicRepository,
	chapterRepository repository.ChapterRepository,
	assignmentRepository repository.AssignmentRepository,
) AssignmentApplication {
	if ossClient == nil ||
		memberRepository == nil ||
		comicRepository == nil ||
		chapterRepository == nil ||
		assignmentRepository == nil {
		zap.L().Panic(
			"NewAssignmentApplication: 依赖项不能为空",
			zap.Bool("ossClient_nil", ossClient == nil),
			zap.Bool("memberRepository_nil", memberRepository == nil),
			zap.Bool("comicRepository_nil", comicRepository == nil),
			zap.Bool("chapterRepository_nil", chapterRepository == nil),
			zap.Bool("assignmentRepository_nil", assignmentRepository == nil),
		)
	}

	return &assignmentApplication{
		ossClient:            ossClient,
		memberRepository:     memberRepository,
		comicRepository:      comicRepository,
		chapterRepository:    chapterRepository,
		assignmentRepository: assignmentRepository,
	}
}

func (aa *assignmentApplication) ListChapterAssignments(
	scope util.TraceScope,
	currentUserID string,
	args value.ListChapterAssignmentArgs,
) ([]value.AssignmentWithUserInfo, error) {
	const fn = "AssignmentApplication.ListChapterAssignments"

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

	if !model.PermAssignmentList().Check(
		currentUserID,
		args.ChapterID,
		adapter.HandleLoadChapterInfo(aa.chapterRepository),
		adapter.HandleLoadComicInfo(aa.comicRepository),
		adapter.HandleLoadMemberInfo(aa.memberRepository),
	) {
		scope.Logger().Warn(fn + ": 权限检查失败")
		return nil, errors.New("权限不足")
	}

	assignments, err := aa.assignmentRepository.ListWithUserInfo(
		nil,
		query_option.UpdatedAtDesc(repository_infra.AssignmentTable),
		query_option.AssignmentQuery().FilterByChapterID(args.ChapterID),
	)
	if err != nil {
		scope.Logger().Error(fn+": 获取分配列表失败", zap.Error(err))
		return nil, errors.New("获取分配列表失败")
	}

	result := make([]value.AssignmentWithUserInfo, len(assignments))
	for i, a := range assignments {
		avatarURL, err := aa.ossClient.GenerateGetPresignedURL(a.User.AvatarOSSKey)
		if err != nil {
			scope.Logger().Error(fn+": 生成头像链接失败", zap.Error(err))
			return nil, errors.New("获取分配列表失败")
		}
		result[i] = value.NewAssignmentWithUserInfo(a, avatarURL)
	}

	return result, nil
}

func (aa *assignmentApplication) ListMyAssignments(
	scope util.TraceScope,
	currentUserID string,
	args value.ListUserAssignmentArgs,
) ([]value.AssignmentWithChapterInfo, error) {
	const fn = "AssignmentApplication.ListMyAssignments"

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

	assignments, err := aa.assignmentRepository.ListWithChapterInfo(
		nil,
		query_option.UpdatedAtDesc(repository_infra.AssignmentTable),
		query_option.AssignmentQuery().FilterByUserID(currentUserID),
		query_option.Paginate(args.Offset, args.Limit),
	)
	if err != nil {
		scope.Logger().Error(fn+": 获取用户分配列表失败", zap.Error(err))
		return nil, errors.New("获取用户分配列表失败")
	}

	result := make([]value.AssignmentWithChapterInfo, len(assignments))
	for i, a := range assignments {
		result[i] = value.NewAssignmentWithChapterInfo(a, a.Chapter.CoverURL)
	}

	return result, nil
}

func (aa *assignmentApplication) CreateChapterAssignment(
	scope util.TraceScope,
	currentUserID string,
	args value.CreateChapterAssignmentArgs,
) (value.CreateChapterAssignmentResult, error) {
	const fn = "AssignmentApplication.CreateChapterAssignment"

	if err := args.Validate(); err != nil {
		scope.Logger().Warn(fn+": 参数验证失败", zap.Error(err))
		return value.CreateChapterAssignmentResult{}, errors.New("参数错误: " + err.Error())
	}

	scope.
		WithFields(
			zap.String("current_user_id", currentUserID),
			zap.Any("args", args),
		).
		Logger().
		Debug(fn + ": 被调用")

	if !model.PermAssignmentCreate().Check(
		currentUserID,
		args.ChapterID,
		adapter.HandleLoadAssignmentInfo(aa.assignmentRepository),
	) {
		scope.Logger().Warn(fn + ": 权限检查失败")
		return value.CreateChapterAssignmentResult{}, errors.New("权限不足")
	}

	isExisting, err := aa.assignmentRepository.Exist(
		nil,
		query_option.AssignmentQuery().FilterByChapterID(args.ChapterID),
		query_option.AssignmentQuery().FilterByUserID(args.UserID),
	)
	if err != nil {
		scope.Logger().Error(fn+": 检查分配记录失败", zap.Error(err))
		return value.CreateChapterAssignmentResult{}, errors.New("创建分配失败")
	}
	if isExisting {
		return value.CreateChapterAssignmentResult{}, errors.New("该用户已在此章节中存在分配记录")
	}

	creation := model.NewAssignmentCreation(args.ChapterID, args.UserID, args.Role)

	assignmentID, err := aa.assignmentRepository.Create(nil, creation)
	if err != nil {
		scope.Logger().Error(fn+": 创建分配失败", zap.Error(err))
		return value.CreateChapterAssignmentResult{}, errors.New("创建分配失败")
	}

	return value.NewCreateChapterAssignmentResult(assignmentID), nil
}

func (aa *assignmentApplication) UpdateAssignment(
	scope util.TraceScope,
	currentUserID string,
	args value.UpdateAssignmentArgs,
) error {
	const fn = "AssignmentApplication.UpdateAssignment"

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

	// 先加载目标分配，获取可信的 ChapterID
	targetAssignment, err := aa.assignmentRepository.Get(
		nil,
		query_option.FilterByID(repository_infra.AssignmentTable, args.ID),
	)
	if err != nil {
		scope.Logger().Error(fn+": 获取分配信息失败", zap.Error(err))
		return errors.New("无法获取分配信息")
	}

	if !model.PermAssignmentUpdate().Check(
		currentUserID,
		targetAssignment.ChapterID,
		adapter.HandleLoadAssignmentInfo(aa.assignmentRepository),
	) {
		scope.Logger().Warn(fn + ": 权限检查失败")
		return errors.New("权限不足")
	}

	// PUT 语义：保留已有角色的时间戳，新增角色写入当前时间，移除角色置 null
	assignmentUpdate := model.NewAssignmentUpdate(args.ID, targetAssignment, args.Role)

	if err := aa.assignmentRepository.Update(nil, assignmentUpdate); err != nil {
		scope.Logger().Error(fn+": 更新分配失败", zap.Error(err))
		return errors.New("更新分配失败")
	}

	return nil
}

func (aa *assignmentApplication) RemoveAssignment(
	scope util.TraceScope,
	currentUserID string,
	assignmentID string,
) error {
	const fn = "AssignmentApplication.RemoveAssignment"

	scope.
		WithFields(
			zap.String("current_user_id", currentUserID),
			zap.String("assignment_id", assignmentID),
		).
		Logger().
		Debug(fn + ": 被调用")

	// 先加载目标分配，获取可信的 ChapterID
	targetAssignment, err := aa.assignmentRepository.Get(
		nil,
		query_option.FilterByID(repository_infra.AssignmentTable, assignmentID),
	)
	if err != nil {
		scope.Logger().Error(fn+": 获取分配信息失败", zap.Error(err))
		return errors.New("无法获取分配信息")
	}

	if !model.PermAssignmentDelete().Check(
		currentUserID,
		targetAssignment.ChapterID,
		adapter.HandleLoadAssignmentInfo(aa.assignmentRepository),
	) {
		scope.Logger().Warn(fn + ": 权限检查失败")
		return errors.New("权限不足")
	}

	if err := aa.assignmentRepository.Delete(nil, assignmentID); err != nil {
		scope.Logger().Error(fn+": 删除分配失败", zap.Error(err))
		return errors.New("删除分配失败")
	}

	return nil
}
