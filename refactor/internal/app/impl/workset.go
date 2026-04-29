package app_impl

import (
	"context"

	app_iface "poprako-s/internal/app"
	app_res "poprako-s/internal/app/res"
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
	errClsf      repo_iface.ErrClsf
}

// `NewWorksetApp` creates a `WorksetApp` implementation.
func NewWorksetApp(
	txnCtrl repo_iface.TxnCtrl,
	worksetSvc svc.WorksetSvc,
	memberRepo repo_iface.MemberRepo,
	worksetRepo repo_iface.WorksetRepo,
	errClsf repo_iface.ErrClsf,
) app_iface.WorksetApp {
	if txnCtrl == nil ||
		memberRepo == nil ||
		worksetRepo == nil ||
		errClsf == nil {
		zap.L().Panic(
			"[NewWorksetApp] nil dependency",
			zap.Bool("txnCtrl", txnCtrl == nil),
			zap.Bool("memberRepo", memberRepo == nil),
			zap.Bool("worksetRepo", worksetRepo == nil),
			zap.Bool("errClsf", errClsf == nil),
		)
	}

	return &worksetAppImpl{
		txnCtrl:     txnCtrl,
		worksetSvc:  worksetSvc,
		memberRepo:  memberRepo,
		worksetRepo: worksetRepo,
		errClsf:      errClsf,
	}
}

// `List` returns all active worksets for the given team.
func (a *worksetAppImpl) List(cx context.Context, currUid string, args *val.ListWorksetArgs) app_res.AppRes[[]val.WorksetVal] {
	lgr := app_util.TakeLgr(cx)

	if re := vfyListWorksetArgs(args); re.IsReject() {
		return app_res.Reject[[]val.WorksetVal](re.Code(), re.Msg())
	}

	/// Permission check: only team members can list the team's worksets.
	if re := a.worksetSvc.CanListWorkset(currUid, args.TeamId, a.memberRepo, a.errClsf); re.IsReject() {
		return app_res.Reject[[]val.WorksetVal](app_res.ErrCode(re.Code()), re.Msg())
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
		enum.WorksetInclTeam,
	)
	if err != nil {
		lgr.Error(
			"[worksetAppImpl.List] failed to list worksets",
			zap.Error(err),
		)

		return app_res.Reject[[]val.WorksetVal](app_res.ServerError, "获取作品集列表失败")
	}

	// Assemble the value-object slice.
	vals := make([]val.WorksetVal, len(worksets))
	for i, ws := range worksets {
		vals[i] = asmWorksetVal(ws)
	}

	return app_res.Accept(&vals)
}

// `Create` creates a new workset inside a team.
func (a *worksetAppImpl) Create(cx context.Context, currUid string, args *val.CreateWorksetArgs) app_res.AppRes[val.WorksetCreatedRes] {
	lgr := app_util.TakeLgr(cx)

	// Open a transaction to atomically verify permission, count active rows,
	// and create the new workset.
	re, err := repo_iface.RunWithTxn(a.txnCtrl, func(prov repo_iface.Prov) (app_res.AppRes[val.WorksetCreatedRes], error) {
		memberRepo := prov.MemberRepo()
		worksetRepo := prov.WorksetRepo()

		// Verify admin role via domain service.
		if re := a.worksetSvc.CanAdminWorkset(currUid, args.TeamId, memberRepo, a.errClsf); re.IsReject() {
			return app_res.Reject[val.WorksetCreatedRes](app_res.ErrCode(re.Code()), re.Msg()), app_res.DefErr()
		}

		// Count active worksets to determine the next index.
		count, err := worksetRepo.Count(&query.ListWorksetOpt{TeamId: &args.TeamId})
		if err != nil {
			return app_res.Reject[val.WorksetCreatedRes](app_res.ServerError, "创建作品集失败"), err
		}

		// Build the creation input via the domain service.
		cre := a.worksetSvc.NewWorksetCre(args.TeamId, int(count), args.Name, args.Desc)

		// Persist the new workset.
		ws, err := worksetRepo.Create(cre)
		if err != nil {
			return app_res.Reject[val.WorksetCreatedRes](app_res.ServerError, "创建作品集失败"), err
		}

		return app_res.Accept(&val.WorksetCreatedRes{Id: ws.Id}), nil
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
func (a *worksetAppImpl) Update(cx context.Context, currUid string, args *val.WorksetUpdArgs) app_res.AppRes[app_res.None] {
	lgr := app_util.TakeLgr(cx)

	// Load the target workset first.
	ws, err := a.worksetRepo.GetById(args.Id)
	if err != nil {
		lgr.Error("[worksetAppImpl.Update] failed to get workset", zap.Error(err))

		return app_res.Reject[app_res.None](app_res.BadRequest, "作品集不存在")
	}

	// Verify admin role via domain service.
	if re := a.worksetSvc.CanAdminWorkset(currUid, ws.TeamId, a.memberRepo, a.errClsf); re.IsReject() {
		return app_res.Reject[app_res.None](app_res.ErrCode(re.Code()), re.Msg())
	}

	// Apply the `PUT` update.
	upd := &aggr.WorksetUpd{Id: args.Id, Name: args.Name, Desc: args.Desc}

	if err := a.worksetRepo.Update(upd); err != nil {
		lgr.Error("[worksetAppImpl.Update] failed to update workset", zap.Error(err))

		return app_res.Reject[app_res.None](app_res.ServerError, "更新作品集失败")
	}

	return app_res.Accept(&app_res.None{})
}

// `Remove` soft-deletes a workset by id.
func (a *worksetAppImpl) Remove(cx context.Context, currUid string, worksetId string) app_res.AppRes[app_res.None] {
	lgr := app_util.TakeLgr(cx)

	// Load the target workset first.
	ws, err := a.worksetRepo.GetById(worksetId)
	if err != nil {
		lgr.Error("[worksetAppImpl.Remove] failed to get workset", zap.Error(err))

		return app_res.Reject[app_res.None](app_res.BadRequest, "作品集不存在")
	}

	// Verify admin role via domain service.
	if re := a.worksetSvc.CanAdminWorkset(currUid, ws.TeamId, a.memberRepo, a.errClsf); re.IsReject() {
		return app_res.Reject[app_res.None](app_res.ErrCode(re.Code()), re.Msg())
	}

	// Soft-delete the workset.
	if err := a.worksetRepo.Remove(worksetId); err != nil {
		lgr.Error("[worksetAppImpl.Remove] failed to soft-delete workset", zap.Error(err))

		return app_res.Reject[app_res.None](app_res.ServerError, "删除作品集失败")
	}

	return app_res.Accept(&app_res.None{})
}
