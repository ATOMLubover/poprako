package app_impl

import (
	"context"

	app_iface "poprako-s/internal/app"
	"poprako-s/internal/app/res"
	app_util "poprako-s/internal/app/util"
	"poprako-s/internal/app/val"
	"poprako-s/internal/domain/model/aggr"
	"poprako-s/internal/domain/model/enum"
	"poprako-s/internal/domain/model/query"
	repo_iface "poprako-s/internal/domain/repo"
	"poprako-s/internal/domain/svc"

	"go.uber.org/zap"
)

// `worksetAppImpl` is the default implementation of `WorksetApp`.
type worksetAppImpl struct {
	txnCtrl     repo_iface.TxnCtrl
	worksetSvc  svc.WorksetSvc
	memberRepo  repo_iface.MemberRepo
	worksetRepo repo_iface.WorksetRepo
}

// `NewWorksetApp` creates a `WorksetApp` implementation.
func NewWorksetApp(
	txnCtrl repo_iface.TxnCtrl,
	worksetSvc svc.WorksetSvc,
	memberRepo repo_iface.MemberRepo,
	worksetRepo repo_iface.WorksetRepo,
) app_iface.WorksetApp {
	if txnCtrl == nil ||
		memberRepo == nil ||
		worksetRepo == nil {
		zap.L().Panic(
			"[NewWorksetApp] nil dependency",
			zap.Bool("txnCtrl", txnCtrl == nil),
			zap.Bool("memberRepo", memberRepo == nil),
			zap.Bool("worksetRepo", worksetRepo == nil),
		)
	}

	return &worksetAppImpl{
		txnCtrl:     txnCtrl,
		worksetSvc:  worksetSvc,
		memberRepo:  memberRepo,
		worksetRepo: worksetRepo,
	}
}

// `List` returns all active worksets for the given team.
func (a *worksetAppImpl) List(cx context.Context, currUid string, args *val.ListWorksetArgs) res.AppRes[[]val.WorksetVal] {
	lgr := app_util.TakeLgr(cx)

	if args == nil {
		return res.Reject[[]val.WorksetVal](res.BadRequest, "分页参数不能为空")
	}

	if args.Offset < 0 {
		return res.Reject[[]val.WorksetVal](res.BadRequest, "offset 不能小于 0")
	}

	if args.Limit <= 0 {
		args.Limit = 20
	}

	if args.Limit > 200 {
		args.Limit = 200
	}

	// Verify that the caller is a member of the team.
	ok, err := a.memberRepo.ExistByUserTeamId(currUid, args.TeamId)
	if err != nil {
		lgr.Error(
			"[worksetAppImpl.List] failed to verify team membership",
			zap.Error(err),
		)

		return res.Reject[[]val.WorksetVal](res.ServerError, "获取作品集列表失败")
	}

	if !ok {
		return res.Reject[[]val.WorksetVal](res.Forbidden, "无权访问该汉化组的作品集")
	}

	// Retrieve all active worksets for the team.
	worksets, err := a.worksetRepo.List(
		&query.ListWorksetOpt{
			TeamId: &args.TeamId,
			Pagi: query.PagiOpt{
				Offset: args.Offset,
				Limit:  args.Limit,
			},
		},
		enum.WorksetInclTeam)
	if err != nil {
		lgr.Error(
			"[worksetAppImpl.List] failed to list worksets",
			zap.Error(err),
		)

		return res.Reject[[]val.WorksetVal](res.ServerError, "获取作品集列表失败")
	}

	// Assemble the value-object slice.
	vals := make([]val.WorksetVal, len(worksets))

	for i, ws := range worksets {
		vals[i] = asmWorksetVal(ws)
	}

	return res.Accept(&vals)
}

// `Create` creates a new workset inside a team.
func (a *worksetAppImpl) Create(cx context.Context, currUid string, args *val.CreateWorksetArgs) res.AppRes[val.WorksetCreatedRes] {
	lgr := app_util.TakeLgr(cx)

	// Open a transaction to atomically verify permission, count active rows,
	// and create the new workset.
	re, err := repo_iface.RunWithTxn[res.AppRes[val.WorksetCreatedRes]](a.txnCtrl, func(prov repo_iface.Prov) (res.AppRes[val.WorksetCreatedRes], error) {
		memberRepo := prov.MemberRepo()
		worksetRepo := prov.WorksetRepo()

		// Verify admin role in the target team.
		member, err := memberRepo.GetByUserTeamId(currUid, args.TeamId)
		if err != nil {
			return res.Reject[val.WorksetCreatedRes](res.Forbidden, "仅汉化组管理员可创建作品集"), res.DefErr()
		}
		if member == nil || !member.HasAnyRole(enum.RoleAdmin) {
			return res.Reject[val.WorksetCreatedRes](res.Forbidden, "仅汉化组管理员可创建作品集"), res.DefErr()
		}

		// Count active worksets to determine the next index.
		count, err := worksetRepo.Count(&query.ListWorksetOpt{TeamId: &args.TeamId})
		if err != nil {
			return res.Reject[val.WorksetCreatedRes](res.ServerError, "创建作品集失败"), err
		}

		// Build the creation input via the domain service.
		cre := a.worksetSvc.NewWorksetCre(args.TeamId, int(count), args.Name, args.Desc)

		// Persist the new workset.
		ws, err := worksetRepo.Create(cre)
		if err != nil {
			return res.Reject[val.WorksetCreatedRes](res.ServerError, "创建作品集失败"), err
		}

		return res.Accept(&val.WorksetCreatedRes{Id: ws.Id}), nil
	})

	if err != nil {
		lgr.Error(
			"[worksetAppImpl.Create] failed to run create workset transaction",
			zap.Error(err),
		)

		return re
	}

	return re
}

// `Update` updates the name and/or description of an existing workset.
func (a *worksetAppImpl) Update(cx context.Context, currUid string, args *val.WorksetUpdArgs) res.AppRes[res.None] {
	lgr := app_util.TakeLgr(cx)

	// Load the target workset first.
	ws, err := a.worksetRepo.GetById(args.Id)
	if err != nil {
		lgr.Error("[worksetAppImpl.Update] failed to get workset", zap.Error(err))

		return res.Reject[res.None](res.BadRequest, "作品集不存在")
	}

	// Verify the caller membership and admin role in the owning team.
	member, err := a.memberRepo.GetByUserTeamId(currUid, ws.TeamId)
	if err != nil {
		return res.Reject[res.None](res.BadRequest, "无权限访问该作品集")
	}

	if member == nil || !member.HasAnyRole(enum.RoleAdmin) {
		return res.Reject[res.None](res.Forbidden, "仅汉化组管理员可更新作品集")
	}

	// Apply the `PUT` update.
	upd := &aggr.WorksetUpd{Id: args.Id, Name: args.Name, Desc: args.Desc}

	if err := a.worksetRepo.Update(upd); err != nil {
		lgr.Error("[worksetAppImpl.Update] failed to update workset", zap.Error(err))

		return res.Reject[res.None](res.ServerError, "更新作品集失败")
	}

	return res.Accept(&res.None{})
}

// `Remove` soft-deletes a workset by id.
func (a *worksetAppImpl) Remove(cx context.Context, currUid string, worksetId string) res.AppRes[res.None] {
	lgr := app_util.TakeLgr(cx)

	// Load the target workset first.
	ws, err := a.worksetRepo.GetById(worksetId)
	if err != nil {
		lgr.Error("[worksetAppImpl.Remove] failed to get workset", zap.Error(err))

		return res.Reject[res.None](res.BadRequest, "作品集不存在")
	}

	// Verify the caller membership and admin role in the owning team.
	member, err := a.memberRepo.GetByUserTeamId(currUid, ws.TeamId)
	if err != nil {
		return res.Reject[res.None](res.BadRequest, "无权限访问该作品集")
	}

	if member == nil || !member.HasAnyRole(enum.RoleAdmin) {
		return res.Reject[res.None](res.Forbidden, "仅汉化组管理员可删除作品集")
	}

	// Soft-delete the workset.
	if err := a.worksetRepo.Remove(worksetId); err != nil {
		lgr.Error("[worksetAppImpl.Remove] failed to soft-delete workset", zap.Error(err))

		return res.Reject[res.None](res.ServerError, "删除作品集失败")
	}

	return res.Accept(&res.None{})
}
