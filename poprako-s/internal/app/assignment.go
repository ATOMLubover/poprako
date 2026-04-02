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

type AssignmentApp interface {
	// ListByChapter 获取指定章节的分配列表
	ListByChapter(
		cx context.Context,
		currUserID string,
		args *val.ListChapterAssignmentArgs,
	) ([]*val.AssignmentInfo, error)

	// ListMy 获取当前用户的所有分配列表
	ListMy(
		cx context.Context,
		currUserID string,
		args *val.ListMyAssignmentArgs,
	) ([]*val.AssignmentInfo, error)

	// Create 创建分配记录
	Create(
		cx context.Context,
		currUserID string,
		args *val.CreateAssignmentArgs,
	) (*val.CreateAssignmentRes, error)

	// Update 更新分配记录的角色
	Update(
		cx context.Context,
		currUserID string,
		args *val.UpdateAssignmentArgs,
	) error

	// Remove 删除分配记录
	Remove(
		cx context.Context,
		currUserID string,
		assignmentID string,
	) error
}

type assignmentAppImpl struct {
	assignmentSvc service.AssignmentService

	assignmentRepo repo.AssignmentRepo
	ossClient      oss.Client
}

func NewAssignmentApp(
	assignmentSvc service.AssignmentService,
	assignmentRepo repo.AssignmentRepo,
	ossClient oss.Client,
) AssignmentApp {
	// 校验构造函数依赖
	if assignmentSvc == nil ||
		assignmentRepo == nil ||
		ossClient == nil {
		zap.L().Panic(
			"NewAssignmentApp: 依赖项不能为空",
			zap.Bool("assignmentSvc_nil", assignmentSvc == nil),
			zap.Bool("assignmentRepo_nil", assignmentRepo == nil),
			zap.Bool("ossClient_nil", ossClient == nil),
		)
	}

	// 返回真实业务实现
	return &assignmentAppImpl{
		assignmentSvc:  assignmentSvc,
		assignmentRepo: assignmentRepo,
		ossClient:      ossClient,
	}
}

func (a *assignmentAppImpl) ListByChapter(
	cx context.Context,
	currUserID string,
	args *val.ListChapterAssignmentArgs,
) ([]*val.AssignmentInfo, error) {
	// 获取上下文中的日志记录器
	lgr := retrieveLgr(cx)

	// 鉴权：检查当前用户是否为该章节的分配人员
	_, err := a.assignmentRepo.Get(model.AssignmentQueryOpt{
		ChapterID: &args.ChapterID,
		UserID:    &currUserID,
	})
	if err != nil {
		// 记录权限校验失败
		lgr.Warn(
			"获取分配列表失败：权限不足",
			zap.String("curr_user_id", currUserID),
			zap.String("chapter_id", args.ChapterID),
		)

		// 返回客户端可展示的错误
		return nil, errors.New("权限不足")
	}

	// 查询分配列表
	assignments, err := a.assignmentRepo.List(model.AssignmentQueryOpt{
		ChapterID: &args.ChapterID,
	})
	if err != nil {
		// 记录查询失败
		lgr.Error(
			"获取分配列表失败",
			zap.String("chapter_id", args.ChapterID),
			zap.Error(err),
		)

		// 返回客户端可展示的错误
		return nil, errors.New("获取分配列表失败")
	}

	// 组装为 app 层值对象列表
	result := make([]*val.AssignmentInfo, len(assignments))

	for i, assignment := range assignments {
		result[i] = assembleAssignmentInfo(&assignment, a.ossClient)
	}

	// 返回分配列表
	return result, nil
}

func (a *assignmentAppImpl) ListMy(
	cx context.Context,
	currUserID string,
	args *val.ListMyAssignmentArgs,
) ([]*val.AssignmentInfo, error) {
	// 获取上下文中的日志记录器
	lgr := retrieveLgr(cx)

	// 查询当前用户的所有分配列表
	assignments, err := a.assignmentRepo.List(model.AssignmentQueryOpt{
		UserID: &currUserID,
	})
	if err != nil {
		// 记录查询失败
		lgr.Error(
			"获取我的分配列表失败",
			zap.String("curr_user_id", currUserID),
			zap.Error(err),
		)

		// 返回客户端可展示的错误
		return nil, errors.New("获取分配列表失败")
	}

	// 组装为 app 层值对象列表
	result := make([]*val.AssignmentInfo, len(assignments))

	for i, assignment := range assignments {
		result[i] = assembleAssignmentInfo(&assignment, a.ossClient)
	}

	// 返回分配列表
	return result, nil
}

func (a *assignmentAppImpl) Create(
	cx context.Context,
	currUserID string,
	args *val.CreateAssignmentArgs,
) (*val.CreateAssignmentRes, error) {
	// 获取上下文中的日志记录器
	lgr := retrieveLgr(cx)

	// 通过领域服务构造创建载荷（含权限校验）
	creation, err := a.assignmentSvc.NewCreation(
		a.assignmentRepo,
		currUserID,
		args.ChapterID,
		args.UserID,
		args.Roles,
	)
	if err != nil {
		// 记录权限校验失败
		lgr.Warn(
			"创建分配失败：权限不足",
			zap.String("curr_user_id", currUserID),
			zap.Error(err),
		)

		// 返回领域服务返回的错误
		return nil, err
	}

	// 持久化分配记录
	assignInfo, err := a.assignmentRepo.Create(creation)
	if err != nil {
		// 记录创建失败
		lgr.Error(
			"创建分配失败：持久化失败",
			zap.Error(err),
		)

		// 返回客户端可展示的错误
		return nil, errors.New("创建分配失败")
	}

	// 返回创建结果
	return &val.CreateAssignmentRes{ID: assignInfo.ID}, nil
}

func (a *assignmentAppImpl) Update(
	cx context.Context,
	currUserID string,
	args *val.UpdateAssignmentArgs,
) error {
	// 获取上下文中的日志记录器
	lgr := retrieveLgr(cx)

	// 查询当前分配信息
	currAssignInfo, err := a.assignmentRepo.GetByID(args.ID)
	if err != nil {
		// 记录查询失败
		lgr.Error(
			"更新分配失败：获取当前分配信息失败",
			zap.String("assignment_id", args.ID),
			zap.Error(err),
		)

		// 返回客户端可展示的错误
		return errors.New("无法获取分配信息")
	}

	// 通过领域服务构造更新载荷（含权限校验）
	update, err := a.assignmentSvc.NewUpdate(
		a.assignmentRepo,
		currUserID,
		args.ID,
		currAssignInfo,
		args.Roles,
	)
	if err != nil {
		// 记录权限校验失败
		lgr.Warn(
			"更新分配失败：权限不足",
			zap.String("curr_user_id", currUserID),
			zap.Error(err),
		)

		// 返回领域服务返回的错误
		return err
	}

	// 持久化更新
	if err := a.assignmentRepo.Update(update); err != nil {
		// 记录更新失败
		lgr.Error(
			"更新分配失败",
			zap.String("assignment_id", args.ID),
			zap.Error(err),
		)

		// 返回客户端可展示的错误
		return errors.New("更新分配失败")
	}

	// 返回更新成功
	return nil
}

func (a *assignmentAppImpl) Remove(
	cx context.Context,
	currUserID string,
	assignmentID string,
) error {
	// 获取上下文中的日志记录器
	lgr := retrieveLgr(cx)

	// 查询目标分配信息以获取所属章节
	targetAssignment, err := a.assignmentRepo.GetByID(assignmentID)
	if err != nil {
		// 记录查询失败
		lgr.Error(
			"删除分配失败：获取目标分配信息失败",
			zap.String("assignment_id", assignmentID),
			zap.Error(err),
		)

		// 返回客户端可展示的错误
		return errors.New("无法获取分配信息")
	}

	// 鉴权：检查当前用户是否为该章节的监修
	currAssignment, err := a.assignmentRepo.Get(model.AssignmentQueryOpt{
		ChapterID: &targetAssignment.ChapterID,
		UserID:    &currUserID,
	})
	if err != nil || !currAssignment.HasAnyRole(model.RoleReviewer) {
		// 记录权限校验失败
		lgr.Warn(
			"删除分配失败：权限不足",
			zap.String("curr_user_id", currUserID),
		)

		// 返回客户端可展示的错误
		return errors.New("权限不足")
	}

	// 执行删除
	if err := a.assignmentRepo.Delete(assignmentID); err != nil {
		// 记录删除失败
		lgr.Error(
			"删除分配失败",
			zap.String("assignment_id", assignmentID),
			zap.Error(err),
		)

		// 返回客户端可展示的错误
		return errors.New("删除分配失败")
	}

	// 返回删除成功
	return nil
}

// assembleAssignmentInfo 将领域层分配信息转换为 app 层值对象
func assembleAssignmentInfo(
	info *model.AssignmentInfo,
	ossClient oss.Client,
) *val.AssignmentInfo {
	result := &val.AssignmentInfo{
		ID:        info.ID,
		ChapterID: info.ChapterID,
		UserID:    info.UserID,
		Roles:     info.AssignedRoleMask(),
		CreatedAt: info.CreatedAt.UnixMilli(),
		UpdatedAt: info.UpdatedAt.UnixMilli(),
	}

	// 若包含章节信息则一并组装
	if info.Chapter != nil {
		result.Chapter = assembleChapterInfo(info.Chapter, ossClient)
	}

	// 若包含用户信息则一并组装
	if info.User != nil {
		userInfo, _ := assembleUserInfo(info.User, ossClient)
		result.User = userInfo
	}

	return result
}

// logAssignmentAppImpl 是 AssignmentApp 的日志包装实现
type logAssignmentAppImpl struct {
	lgr *zap.Logger
	app AssignmentApp
}

func NewLogAssignmentApp(
	lgr *zap.Logger,
	app AssignmentApp,
) AssignmentApp {
	if lgr == nil || app == nil {
		zap.L().Panic(
			"NewLogAssignmentApp: 依赖项不能为空",
			zap.Bool("lgr_nil", lgr == nil),
			zap.Bool("app_nil", app == nil),
		)
	}

	return &logAssignmentAppImpl{lgr: lgr, app: app}
}

func (a *logAssignmentAppImpl) ListByChapter(
	cx context.Context,
	currUserID string,
	args *val.ListChapterAssignmentArgs,
) ([]*val.AssignmentInfo, error) {
	if a == nil || a.app == nil || a.lgr == nil {
		return nil, errors.New("AssignmentApp 不可用")
	}

	if args == nil || args.ChapterID == "" {
		return nil, errors.New("参数不合法")
	}

	lgr := a.lgr.With(zap.String("method", "ListByChapter"), zap.String("curr_user_id", currUserID), zap.String("chapter_id", args.ChapterID))

	cx = injectLgr(cx, lgr)

	lgr.Info("[logAssignmentAppImpl.ListByChapter] CALL")

	return a.app.ListByChapter(cx, currUserID, args)
}

func (a *logAssignmentAppImpl) ListMy(
	cx context.Context,
	currUserID string,
	args *val.ListMyAssignmentArgs,
) ([]*val.AssignmentInfo, error) {
	if a == nil || a.app == nil || a.lgr == nil {
		return nil, errors.New("AssignmentApp 不可用")
	}

	if args == nil {
		return nil, errors.New("参数不合法")
	}

	lgr := a.lgr.With(zap.String("method", "ListMy"), zap.String("curr_user_id", currUserID))

	cx = injectLgr(cx, lgr)

	lgr.Info("[logAssignmentAppImpl.ListMy] CALL")

	return a.app.ListMy(cx, currUserID, args)
}

func (a *logAssignmentAppImpl) Create(
	cx context.Context,
	currUserID string,
	args *val.CreateAssignmentArgs,
) (*val.CreateAssignmentRes, error) {
	if a == nil || a.app == nil || a.lgr == nil {
		return nil, errors.New("AssignmentApp 不可用")
	}

	if args == nil || args.ChapterID == "" || args.UserID == "" {
		return nil, errors.New("参数不合法")
	}

	lgr := a.lgr.With(zap.String("method", "Create"), zap.String("curr_user_id", currUserID))

	cx = injectLgr(cx, lgr)

	lgr.Info("[logAssignmentAppImpl.Create] CALL")

	return a.app.Create(cx, currUserID, args)
}

func (a *logAssignmentAppImpl) Update(
	cx context.Context,
	currUserID string,
	args *val.UpdateAssignmentArgs,
) error {
	if a == nil || a.app == nil || a.lgr == nil {
		return errors.New("AssignmentApp 不可用")
	}

	if args == nil || args.ID == "" {
		return errors.New("参数不合法")
	}

	lgr := a.lgr.With(zap.String("method", "Update"), zap.String("curr_user_id", currUserID), zap.String("assignment_id", args.ID))

	cx = injectLgr(cx, lgr)

	lgr.Info("[logAssignmentAppImpl.Update] CALL")

	return a.app.Update(cx, currUserID, args)
}

func (a *logAssignmentAppImpl) Remove(
	cx context.Context,
	currUserID string,
	assignmentID string,
) error {
	if a == nil || a.app == nil || a.lgr == nil {
		return errors.New("AssignmentApp 不可用")
	}

	if assignmentID == "" {
		return errors.New("分配 ID 不能为空")
	}

	lgr := a.lgr.With(zap.String("method", "Remove"), zap.String("curr_user_id", currUserID), zap.String("assignment_id", assignmentID))

	cx = injectLgr(cx, lgr)

	lgr.Info("[logAssignmentAppImpl.Remove] CALL")

	return a.app.Remove(cx, currUserID, assignmentID)
}
