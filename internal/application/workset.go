package application

import (
	"errors"

	"labelplus-next-web-be/internal/application/adapter"
	"labelplus-next-web-be/internal/application/assembler"
	"labelplus-next-web-be/internal/domain/external"
	"labelplus-next-web-be/internal/domain/model"
	"labelplus-next-web-be/internal/domain/repository"
	domain_service "labelplus-next-web-be/internal/domain/service"
	repository_infra "labelplus-next-web-be/internal/infrastructure/repository"
	"labelplus-next-web-be/internal/infrastructure/repository/query_option"
	"labelplus-next-web-be/internal/util"
	"labelplus-next-web-be/internal/value"

	"go.uber.org/zap"
)

type WorksetApplication interface {
	ListWorksets(
		scope util.TraceScope,
		currentUserID string,
		args value.ListWorksetArgs,
	) ([]value.WorksetInfo, error)
	CreateWorkset(
		scope util.TraceScope,
		currentUserID string,
		args value.CreateWorksetArgs,
	) (value.CreateWorksetResult, error)
	UpdateWorkset(
		scope util.TraceScope,
		currentUserID string,
		args value.UpdateWorksetArgs,
	) error
	DeleteWorkset(
		scope util.TraceScope,
		currentUserID string,
		worksetID string,
	) error
}

type worksetApplication struct {
	ossClient         external.OSSClient
	memberRepository  repository.MemberRepository
	worksetRepository repository.WorksetRepository
}

func NewWorksetApplication(
	ossClient external.OSSClient,
	memberRepository repository.MemberRepository,
	worksetRepository repository.WorksetRepository,
) WorksetApplication {
	if ossClient == nil ||
		memberRepository == nil ||
		worksetRepository == nil {
		zap.L().Panic(
			"NewWorksetApplication: 依赖项不能为空",
			zap.Bool("ossClient_nil", ossClient == nil),
			zap.Bool("memberRepository_nil", memberRepository == nil),
			zap.Bool("worksetRepository_nil", worksetRepository == nil),
		)
	}

	return &worksetApplication{
		ossClient:         ossClient,
		memberRepository:  memberRepository,
		worksetRepository: worksetRepository,
	}
}

func (wa *worksetApplication) ListWorksets(
	scope util.TraceScope,
	currentUserID string,
	args value.ListWorksetArgs,
) ([]value.WorksetInfo, error) {
	const fn = "WorksetApplication.ListWorksets"

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

	// 鉴权：检查当前用户在指定的汉化组是否有权限查看工作集列表
	if !model.PermWorksetList().Check(
		currentUserID,
		args.TeamID,
		adapter.HandleLoadMemberInfo(wa.memberRepository),
	) {
		scope.Logger().Warn(fn + ": 权限检查失败")
		return nil, errors.New("权限不足")
	}

	includeSpec := domain_service.ResolveWorksetListIncludeSpec(args.Includes)
	queryOptions := []repository.QueryOption{
		query_option.WorksetQuery().FilterByTeamID(args.TeamID),
		query_option.WorksetQuery().OrderByIndexAsc(),
	}

	if includeSpec.NeedTeam {
		queryOptions = append(queryOptions, query_option.WorksetQuery().IncludeTeamInfo())
	}

	queryOptions = append(queryOptions, query_option.Paginate(args.Offset, args.Limit))

	worksetList, err := wa.worksetRepository.List(nil, queryOptions...)
	if err != nil {
		scope.Logger().Error(fn+": 获取工作集列表失败", zap.Error(err))
		return nil, errors.New("无法获取工作集列表")
	}

	result := make([]value.WorksetInfo, len(worksetList))
	for i, workset := range worksetList {
		result[i] = assembler.AssembleWorksetInfo(workset, wa.ossClient.GenerateGetPresignedURL)
	}

	return result, nil
}

func (wa *worksetApplication) CreateWorkset(
	scope util.TraceScope,
	currentUserID string,
	args value.CreateWorksetArgs,
) (value.CreateWorksetResult, error) {
	const fn = "WorksetApplication.CreateWorkset"

	if err := args.Validate(); err != nil {
		scope.Logger().Warn(fn+": 参数验证失败", zap.Error(err))
		return value.CreateWorksetResult{}, errors.New("参数错误: " + err.Error())
	}

	scope.
		WithFields(
			zap.String("current_user_id", currentUserID),
			zap.Any("args", args),
		).
		Logger().
		Debug(fn + ": 被调用")

	// 鉴权：检查当前用户在指定汉化组是否有创建工作集权限
	if !model.PermWorksetCreate().Check(
		currentUserID,
		args.TeamID,
		adapter.HandleLoadMemberInfo(wa.memberRepository),
	) {
		scope.Logger().Warn(fn + ": 权限检查失败")
		return value.CreateWorksetResult{}, errors.New("权限不足")
	}

	transactionExecutor := wa.worksetRepository.BeginTransaction()
	if transactionExecutor.Error != nil {
		scope.Logger().Error(fn+": 开启事务失败", zap.Error(transactionExecutor.Error))
		return value.CreateWorksetResult{}, errors.New("创建工作集失败")
	}

	var transactionErr error

	defer func() {
		if transactionErr != nil {
			if rollbackErr := transactionExecutor.Rollback().Error; rollbackErr != nil {
				scope.Logger().Error(fn+": 回滚事务失败", zap.Error(rollbackErr))
			}
		}
	}()

	transactionErr = wa.worksetRepository.LockByTeamID(transactionExecutor, args.TeamID)
	if transactionErr != nil {
		scope.Logger().Error(fn+": 锁定工作集记录失败", zap.Error(transactionErr))
		return value.CreateWorksetResult{}, errors.New("创建工作集失败")
	}

	worksetCount, transactionErr := wa.worksetRepository.Count(
		transactionExecutor,
		query_option.WorksetQuery().FilterByTeamID(args.TeamID),
	)
	if transactionErr != nil {
		scope.Logger().Error(fn+": 统计工作集数量失败", zap.Error(transactionErr))
		return value.CreateWorksetResult{}, errors.New("创建工作集失败")
	}

	description := ""
	if args.Description != nil {
		description = *args.Description
	}

	worksetCreation := model.NewWorksetCreation(
		args.TeamID,
		int(worksetCount),
		args.Name,
		description,
	)

	worksetID, transactionErr := wa.worksetRepository.Create(transactionExecutor, *worksetCreation)
	if transactionErr != nil {
		scope.Logger().Error(fn+": 创建工作集失败", zap.Error(transactionErr))
		return value.CreateWorksetResult{}, errors.New("创建工作集失败")
	}

	if commitErr := transactionExecutor.Commit().Error; commitErr != nil {
		scope.Logger().Error(fn+": 提交事务失败", zap.Error(commitErr))
		return value.CreateWorksetResult{}, errors.New("创建工作集失败")
	}

	return value.CreateWorksetResult{ID: worksetID}, nil
}

func (wa *worksetApplication) UpdateWorkset(
	scope util.TraceScope,
	currentUserID string,
	args value.UpdateWorksetArgs,
) error {
	const fn = "WorksetApplication.UpdateWorkset"

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

	targetWorkset, err := wa.worksetRepository.Get(
		nil,
		query_option.FilterByID(repository_infra.WorksetTable, args.ID),
	)
	if err != nil {
		scope.Logger().Error(fn+": 获取目标工作集信息失败", zap.Error(err))
		return errors.New("无法获取工作集信息")
	}

	// 鉴权：检查当前用户在目标工作集所属汉化组是否有更新权限
	if !model.PermWorksetUpdate().Check(
		currentUserID,
		targetWorkset.TeamID,
		adapter.HandleLoadMemberInfo(wa.memberRepository),
	) {
		scope.Logger().Warn(fn + ": 权限检查失败")
		return errors.New("权限不足")
	}

	worksetUpdate := model.NewWorksetUpdate(args.ID, args.Name, args.Description)

	if err := wa.worksetRepository.Update(nil, worksetUpdate); err != nil {
		scope.Logger().Error(fn+": 更新工作集失败", zap.Error(err))
		return errors.New("更新工作集失败")
	}

	return nil
}

func (wa *worksetApplication) DeleteWorkset(
	scope util.TraceScope,
	currentUserID string,
	worksetID string,
) error {
	const fn = "WorksetApplication.DeleteWorkset"

	if worksetID == "" {
		scope.Logger().Warn(fn + ": worksetID 为空")
		return errors.New(ErrInternalError)
	}

	scope.
		WithFields(
			zap.String("current_user_id", currentUserID),
			zap.String("workset_id", worksetID),
		).
		Logger().
		Debug(fn + ": 被调用")

	targetWorkset, err := wa.worksetRepository.Get(
		nil,
		query_option.FilterByID(repository_infra.WorksetTable, worksetID),
	)
	if err != nil {
		scope.Logger().Error(fn+": 获取目标工作集信息失败", zap.Error(err))
		return errors.New("无法获取工作集信息")
	}

	// 鉴权：检查当前用户在目标工作集所属汉化组是否有删除权限
	if !model.PermWorksetDelete().Check(
		currentUserID,
		targetWorkset.TeamID,
		adapter.HandleLoadMemberInfo(wa.memberRepository),
	) {
		scope.Logger().Warn(fn + ": 权限检查失败")
		return errors.New("权限不足")
	}

	if err := wa.worksetRepository.Delete(nil, worksetID); err != nil {
		scope.Logger().Error(fn+": 删除工作集失败", zap.Error(err))
		return errors.New("删除工作集失败")
	}

	return nil
}
