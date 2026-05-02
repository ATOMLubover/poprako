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

type WorksetApp interface {
	// List 获取指定汉化组的作品集列表
	List(
		cx context.Context,
		currUserID string,
		args *val.ListWorksetArgs,
	) ([]*val.WorksetInfo, error)

	// Create 创建一个新的作品集
	Create(
		cx context.Context,
		currUserID string,
		args *val.CreateWorksetArgs,
	) (*val.CreateWorksetRes, error)

	// Update 更新作品集信息
	Update(
		cx context.Context,
		currUserID string,
		args *val.UpdateWorksetArgs,
	) error

	// Remove 删除作品集
	Remove(
		cx context.Context,
		currUserID string,
		worksetID string,
	) error
}

type worksetAppImpl struct {
	worksetSvc service.WorksetService

	memberRepo  repo.MemberRepo
	worksetRepo repo.WorksetRepo
	txnMgr      repo.TxnMgr
}

func NewWorksetApp(
	worksetSvc service.WorksetService,
	memberRepo repo.MemberRepo,
	worksetRepo repo.WorksetRepo,
	txnMgr repo.TxnMgr,
) WorksetApp {
	// 校验构造函数依赖
	if worksetSvc == nil ||
		memberRepo == nil ||
		worksetRepo == nil ||
		txnMgr == nil {
		zap.L().Panic(
			"NewWorksetApp: 依赖项不能为空",
			zap.Bool("worksetSvc_nil", worksetSvc == nil),
			zap.Bool("memberRepo_nil", memberRepo == nil),
			zap.Bool("worksetRepo_nil", worksetRepo == nil),
			zap.Bool("txnMgr_nil", txnMgr == nil),
		)
	}

	// 返回真实业务实现
	return &worksetAppImpl{
		worksetSvc:  worksetSvc,
		memberRepo:  memberRepo,
		worksetRepo: worksetRepo,
		txnMgr:      txnMgr,
	}
}

func (a *worksetAppImpl) List(
	cx context.Context,
	currUserID string,
	args *val.ListWorksetArgs,
) ([]*val.WorksetInfo, error) {
	// 获取上下文中的日志记录器
	lgr := retrieveLgr(cx)

	// 鉴权：检查当前用户是否为该汉化组成员
	_, err := a.memberRepo.Get(model.MemberQueryOpt{
		UserID: &currUserID,
		TeamID: &args.TeamID,
	})
	if err != nil {
		// 记录权限校验失败
		lgr.Warn(
			"获取作品集列表失败：权限不足",
			zap.String("curr_user_id", currUserID),
			zap.String("team_id", args.TeamID),
		)

		// 返回客户端可展示的错误
		return nil, errors.New("权限不足")
	}

	// 查询作品集列表
	worksets, err := a.worksetRepo.List(model.WorksetQueryOpt{
		TeamID: &args.TeamID,
	})
	if err != nil {
		// 记录查询失败
		lgr.Error(
			"获取作品集列表失败",
			zap.String("team_id", args.TeamID),
			zap.Error(err),
		)

		// 返回客户端可展示的错误
		return nil, errors.New("获取作品集列表失败")
	}

	// 组装为 app 层值对象列表
	result := make([]*val.WorksetInfo, len(worksets))

	for i, ws := range worksets {
		result[i] = assembleWorksetInfo(&ws)
	}

	// 返回作品集列表
	return result, nil
}

func (a *worksetAppImpl) Create(
	cx context.Context,
	currUserID string,
	args *val.CreateWorksetArgs,
) (*val.CreateWorksetRes, error) {
	// 获取上下文中的日志记录器
	lgr := retrieveLgr(cx)

	// 默认描述为空字符串
	desc := ""

	if args.Desc != nil {
		desc = *args.Desc
	}

	// 在事务中创建作品集（需要 count 获取 index）

	var createdID string

	if err := a.txnMgr.RunInTxn(func(cx context.Context) error {
		worksetRepoTxn, err := a.worksetRepo.FromTxnCx(cx)
		if err != nil {
			return err
		}

		memberRepoTxn, err := a.memberRepo.FromTxnCx(cx)
		if err != nil {
			return err
		}

		// 统计当前汉化组下的作品集数量以确定 index
		count, err := worksetRepoTxn.Count(model.WorksetQueryOpt{
			TeamID: &args.TeamID,
		})
		if err != nil {
			return err
		}

		// 通过领域服务构造创建载荷（含权限校验）
		creation, err := a.worksetSvc.NewCreation(
			memberRepoTxn,
			currUserID,
			args.TeamID,
			int(count),
			args.Name,
			desc,
		)
		if err != nil {
			return err
		}

		// 持久化作品集
		wsInfo, err := worksetRepoTxn.Create(creation)
		if err != nil {
			return err
		}

		createdID = wsInfo.ID

		return nil
	}); err != nil {
		// 记录创建失败
		lgr.Error(
			"创建作品集失败",
			zap.String("curr_user_id", currUserID),
			zap.String("team_id", args.TeamID),
			zap.Error(err),
		)

		// 返回客户端可展示的错误
		return nil, errors.New("创建作品集失败")
	}

	// 返回创建结果
	return &val.CreateWorksetRes{ID: createdID}, nil
}

func (a *worksetAppImpl) Update(
	cx context.Context,
	currUserID string,
	args *val.UpdateWorksetArgs,
) error {
	// 获取上下文中的日志记录器
	lgr := retrieveLgr(cx)

	// 查询目标作品集信息以获取所属汉化组 ID
	targetWorkset, err := a.worksetRepo.GetByID(args.ID)
	if err != nil {
		// 记录查询失败
		lgr.Error(
			"更新作品集失败：获取目标作品集信息失败",
			zap.String("workset_id", args.ID),
			zap.Error(err),
		)

		// 返回客户端可展示的错误
		return errors.New("无法获取作品集信息")
	}

	// 鉴权：检查当前用户在作品集所属汉化组中是否为管理员
	currMember, err := a.memberRepo.Get(model.MemberQueryOpt{
		UserID: &currUserID,
		TeamID: &targetWorkset.TeamID,
	})
	if err != nil || !currMember.HasAnyRole(model.RoleAdmin) {
		// 记录权限校验失败
		lgr.Warn(
			"更新作品集失败：权限不足",
			zap.String("curr_user_id", currUserID),
		)

		// 返回客户端可展示的错误
		return errors.New("权限不足")
	}

	// 构造更新载荷
	update := &model.WorksetUpdate{
		ID:   args.ID,
		Name: args.Name,
		Desc: args.Desc,
	}

	// 持久化更新
	if err := a.worksetRepo.Update(update); err != nil {
		// 记录更新失败
		lgr.Error(
			"更新作品集失败",
			zap.String("workset_id", args.ID),
			zap.Error(err),
		)

		// 返回客户端可展示的错误
		return errors.New("更新作品集失败")
	}

	// 返回更新成功
	return nil
}

func (a *worksetAppImpl) Remove(
	cx context.Context,
	currUserID string,
	worksetID string,
) error {
	// 获取上下文中的日志记录器
	lgr := retrieveLgr(cx)

	// 查询目标作品集信息以获取所属汉化组 ID
	targetWorkset, err := a.worksetRepo.GetByID(worksetID)
	if err != nil {
		// 记录查询失败
		lgr.Error(
			"删除作品集失败：获取目标作品集信息失败",
			zap.String("workset_id", worksetID),
			zap.Error(err),
		)

		// 返回客户端可展示的错误
		return errors.New("无法获取作品集信息")
	}

	// 鉴权：检查当前用户在作品集所属汉化组中是否为管理员
	currMember, err := a.memberRepo.Get(model.MemberQueryOpt{
		UserID: &currUserID,
		TeamID: &targetWorkset.TeamID,
	})
	if err != nil || !currMember.HasAnyRole(model.RoleAdmin) {
		// 记录权限校验失败
		lgr.Warn(
			"删除作品集失败：权限不足",
			zap.String("curr_user_id", currUserID),
		)

		// 返回客户端可展示的错误
		return errors.New("权限不足")
	}

	// 执行删除
	if err := a.worksetRepo.Delete(worksetID); err != nil {
		// 记录删除失败
		lgr.Error(
			"删除作品集失败",
			zap.String("workset_id", worksetID),
			zap.Error(err),
		)

		// 返回客户端可展示的错误
		return errors.New("删除作品集失败")
	}

	// 返回删除成功
	return nil
}

// assembleWorksetInfo 将领域层作品集信息转换为 app 层值对象
func assembleWorksetInfo(info *model.WorksetInfo) *val.WorksetInfo {
	result := &val.WorksetInfo{
		ID:         info.ID,
		TeamID:     info.TeamID,
		Index:      info.Index,
		Name:       info.Name,
		Desc:       info.Desc,
		ComicCount: info.ComicCount,
		CreatedAt:  info.CreatedAt.UnixMilli(),
		UpdatedAt:  info.UpdatedAt.UnixMilli(),
	}

	return result
}

// logWorksetAppImpl 是 WorksetApp 的日志包装实现
type logWorksetAppImpl struct {
	app WorksetApp
}

func NewLogWorksetApp(
	app WorksetApp,
) WorksetApp {
	if app == nil {
		zap.L().Panic(
			"NewLogWorksetApp: 依赖项不能为空",
			zap.Bool("app_nil", app == nil),
		)
	}

	return &logWorksetAppImpl{app: app}
}

func (a *logWorksetAppImpl) List(
	cx context.Context,
	currUserID string,
	args *val.ListWorksetArgs,
) ([]*val.WorksetInfo, error) {
	if a == nil || a.app == nil {
		return nil, errors.New("WorksetApp 不可用")
	}

	if args == nil || args.TeamID == "" {
		return nil, errors.New("参数不合法")
	}

	lgr := retrieveLgr(cx).With(zap.String("method", "List"), zap.String("curr_user_id", currUserID), zap.String("team_id", args.TeamID))

	cx = injectLgr(cx, lgr)

	lgr.Info("[logWorksetAppImpl.List] CALL")

	return a.app.List(cx, currUserID, args)
}

func (a *logWorksetAppImpl) Create(
	cx context.Context,
	currUserID string,
	args *val.CreateWorksetArgs,
) (*val.CreateWorksetRes, error) {
	if a == nil || a.app == nil {
		return nil, errors.New("WorksetApp 不可用")
	}

	if args == nil || args.TeamID == "" || args.Name == "" {
		return nil, errors.New("参数不合法")
	}

	lgr := retrieveLgr(cx).With(zap.String("method", "Create"), zap.String("curr_user_id", currUserID))

	cx = injectLgr(cx, lgr)

	lgr.Info("[logWorksetAppImpl.Create] CALL")

	return a.app.Create(cx, currUserID, args)
}

func (a *logWorksetAppImpl) Update(
	cx context.Context,
	currUserID string,
	args *val.UpdateWorksetArgs,
) error {
	if a == nil || a.app == nil {
		return errors.New("WorksetApp 不可用")
	}

	if args == nil || args.ID == "" {
		return errors.New("参数不合法")
	}

	lgr := retrieveLgr(cx).With(zap.String("method", "Update"), zap.String("curr_user_id", currUserID), zap.String("workset_id", args.ID))

	cx = injectLgr(cx, lgr)

	lgr.Info("[logWorksetAppImpl.Update] CALL")

	return a.app.Update(cx, currUserID, args)
}

func (a *logWorksetAppImpl) Remove(
	cx context.Context,
	currUserID string,
	worksetID string,
) error {
	if a == nil || a.app == nil {
		return errors.New("WorksetApp 不可用")
	}

	if worksetID == "" {
		return errors.New("作品集 ID 不能为空")
	}

	lgr := retrieveLgr(cx).With(zap.String("method", "Remove"), zap.String("curr_user_id", currUserID), zap.String("workset_id", worksetID))

	cx = injectLgr(cx, lgr)

	lgr.Info("[logWorksetAppImpl.Remove] CALL")

	return a.app.Remove(cx, currUserID, worksetID)
}
