package app_impl

import (
	"context"

	app_iface "poprako-s/internal/app"
	app_res "poprako-s/internal/app/res"
	app_util "poprako-s/internal/app/util"
	"poprako-s/internal/app/val"
	oss_iface "poprako-s/internal/domain/ext/oss"
	"poprako-s/internal/domain/model/aggr"
	"poprako-s/internal/domain/model/enum"
	"poprako-s/internal/domain/model/query"
	repo_iface "poprako-s/internal/domain/repo"
	"poprako-s/internal/domain/svc"
	repo_infra "poprako-s/internal/infra/repo"

	"go.uber.org/zap"
)

// `teamAppImpl` is the default implementation of `TeamApp`.
type teamAppImpl struct {
	txnCtrl repo_iface.TxnCtrl

	teamSvc   svc.TeamSvc
	ossMsgSvc svc.OssMsgSvc

	userRepo   repo_iface.UserRepo
	teamRepo   repo_iface.TeamRepo
	memberRepo repo_iface.MemberRepo
	ossMsgRepo repo_iface.OssMsgRepo

	ossSigner oss_iface.Signer
	errClsf   repo_iface.ErrClsf
}

// `NewTeamApp` creates an app implementation for team use-cases.
func NewTeamApp(
	txnCtrl repo_iface.TxnCtrl,
	teamRepo repo_iface.TeamRepo,
	userRepo repo_iface.UserRepo,
	memberRepo repo_iface.MemberRepo,
	ossMsgRepo repo_iface.OssMsgRepo,
	teamSvc svc.TeamSvc,
	ossMsgSvc svc.OssMsgSvc,
	ossSigner oss_iface.Signer,
	errClsf repo_iface.ErrClsf,
) app_iface.TeamApp {
	if txnCtrl == nil || teamRepo == nil || userRepo == nil || memberRepo == nil || ossMsgRepo == nil || ossSigner == nil || errClsf == nil {
		zap.L().Panic(
			"[NewTeamApp] nil dependency",
			zap.Bool("txnCtrl", txnCtrl == nil),
			zap.Bool("teamRepo", teamRepo == nil),
			zap.Bool("userRepo", userRepo == nil),
			zap.Bool("memberRepo", memberRepo == nil),
			zap.Bool("ossMsgRepo", ossMsgRepo == nil),
			zap.Bool("ossSigner", ossSigner == nil),
			zap.Bool("errClsf", errClsf == nil),
		)
	}

	return &teamAppImpl{
		txnCtrl:    txnCtrl,
		teamSvc:    teamSvc,
		ossMsgSvc:  ossMsgSvc,
		userRepo:   userRepo,
		teamRepo:   teamRepo,
		memberRepo: memberRepo,
		ossMsgRepo: ossMsgRepo,
		ossSigner:  ossSigner,
		errClsf:    errClsf,
	}
}

// `Create` creates one team with super-admin permission.
func (a *teamAppImpl) Create(cx context.Context, currUid string, args *val.TeamCreArgs) app_res.AppRes[val.TeamCreRes] {
	lgr := app_util.TakeLgr(cx)

	if args == nil || args.Name == "" {
		return app_res.Reject[val.TeamCreRes](app_res.BadRequest, "name 不能为空")
	}

	if re := a.teamSvc.CanListTeam(currUid, a.userRepo, a.errClsf); re.IsReject() {
		return app_res.Reject[val.TeamCreRes](app_res.ErrCode(re.Code()), re.Msg())
	}

	cre, err := a.teamSvc.NewTeamCre(args.Name, args.Desc)
	if err != nil {
		return app_res.Reject[val.TeamCreRes](app_res.BadRequest, err.Error())
	}

	team, err := a.teamRepo.Create(cre)
	if err != nil {
		lgr.Error("[teamAppImpl.Create] failed to create team", zap.Error(err))

		if repo_infra.IsDupKey(err) {
			return app_res.Reject[val.TeamCreRes](app_res.BadRequest, "团队名称已存在")
		}

		return app_res.Reject[val.TeamCreRes](app_res.ServerError, "创建团队失败")
	}

	return app_res.Accept(&val.TeamCreRes{Id: team.Id})
}

// `GetInfo` gets team information by id.
func (a *teamAppImpl) GetInfo(cx context.Context, id string) app_res.AppRes[val.TeamVal] {
	lgr := app_util.TakeLgr(cx)

	team, err := a.teamRepo.GetById(id)
	if err != nil {
		lgr.Error(
			"[teamAppImpl.GetInfo] failed to get team by id",
			zap.String("team_id", id),
			zap.Error(err),
		)

		if repo_infra.IsNotFound(err) {
			return app_res.Reject[val.TeamVal](app_res.BadRequest, "团队不存在")
		}

		return app_res.Reject[val.TeamVal](app_res.ServerError, "获取团队信息失败")
	}

	teamVal, err := asmTeamVal(team, a.ossSigner)
	if err != nil {
		lgr.Error(
			"[teamAppImpl.GetInfo] failed to assemble team info",
			zap.String("team_id", id),
			zap.Error(err),
		)

		return app_res.Reject[val.TeamVal](app_res.ServerError, "获取团队信息失败")
	}

	return app_res.Accept(teamVal)
}

// `List` lists all teams.
func (a *teamAppImpl) List(cx context.Context, currUid string, args *val.ListTeamArgs) app_res.AppRes[[]val.TeamVal] {
	lgr := app_util.TakeLgr(cx)

	if re := vfyListTeamArgs(args); re.IsReject() {
		return app_res.Reject[[]val.TeamVal](re.Code(), re.Msg())
	}

	if re := a.teamSvc.CanListTeam(currUid, a.userRepo, a.errClsf); re.IsReject() {
		return app_res.Reject[[]val.TeamVal](app_res.ErrCode(re.Code()), re.Msg())
	}

	teams, err := a.teamRepo.List(&query.ListTeamOpt{Pagi: query.PagiOpt{Offset: args.Offset, Limit: args.Limit}})
	if err != nil {
		lgr.Error("[teamAppImpl.List] failed to list teams", zap.Error(err))

		return app_res.Reject[[]val.TeamVal](app_res.ServerError, "获取团队列表失败")
	}

	teamVals, err := asmTeamVals(teams, a.ossSigner)
	if err != nil {
		lgr.Error("[teamAppImpl.List] failed to assemble team values", zap.Error(err))

		return app_res.Reject[[]val.TeamVal](app_res.ServerError, "获取团队列表失败")
	}

	return app_res.Accept(&teamVals)
}

// `ListByUser` lists teams that a user belongs to.
func (a *teamAppImpl) ListByUser(cx context.Context, userId string, args *val.ListTeamArgs) app_res.AppRes[[]val.TeamVal] {
	lgr := app_util.TakeLgr(cx)

	if re := vfyListByUserArgs(args); re.IsReject() {
		return app_res.Reject[[]val.TeamVal](re.Code(), re.Msg())
	}

	teamIds, err := listUserTeamIds(a.memberRepo, userId)
	if err != nil {
		lgr.Error("[teamAppImpl.ListByUser] failed to list member rows", zap.Error(err))

		return app_res.Reject[[]val.TeamVal](app_res.ServerError, "获取团队列表失败")
	}

	offset := args.Offset
	if offset > len(teamIds) {
		offset = len(teamIds)
	}

	end := offset + args.Limit
	if end > len(teamIds) {
		end = len(teamIds)
	}

	teamIds = teamIds[offset:end]

	teams, err := listTeamsByIds(a.teamRepo, teamIds)
	if err != nil {
		lgr.Error("[teamAppImpl.ListByUser] failed to load teams", zap.Error(err))

		return app_res.Reject[[]val.TeamVal](app_res.ServerError, "获取团队列表失败")
	}

	teamVals, err := asmTeamVals(teams, a.ossSigner)
	if err != nil {
		lgr.Error("[teamAppImpl.ListByUser] failed to assemble team values", zap.Error(err))

		return app_res.Reject[[]val.TeamVal](app_res.ServerError, "获取团队列表失败")
	}

	return app_res.Accept(&teamVals)
}

// `Update` updates one team profile.
func (a *teamAppImpl) Update(cx context.Context, currUid string, args *val.TeamUpdArgs) app_res.AppRes[app_res.None] {
	lgr := app_util.TakeLgr(cx)

	if args == nil || args.Id == "" || args.Name == "" {
		return app_res.Reject[app_res.None](app_res.BadRequest, "id 和 name 不能为空")
	}

	if re := a.teamSvc.CanAdminTeam(currUid, args.Id, a.memberRepo, a.errClsf); re.IsReject() {
		return app_res.Reject[app_res.None](app_res.ErrCode(re.Code()), re.Msg())
	}

	if err := a.teamRepo.Update(&aggr.TeamUpd{Id: args.Id, Name: args.Name, Desc: args.Desc}); err != nil {
		if repo_infra.IsNotFound(err) {
			return app_res.Reject[app_res.None](app_res.BadRequest, "团队不存在")
		}

		lgr.Error("[teamAppImpl.Update] failed to update team", zap.Error(err))

		return app_res.Reject[app_res.None](app_res.ServerError, "更新团队失败")
	}

	return app_res.Accept(&app_res.None{})
}

// `ResvAvatar` reserves avatar upload for team.
func (a *teamAppImpl) ResvAvatar(cx context.Context, currUid string, args *val.ResvTeamAvatarArgs) app_res.AppRes[val.ResvTeamAvatarRes] {
	lgr := app_util.TakeLgr(cx)

	if args == nil || args.TeamId == "" || args.FileExt == "" {
		return app_res.Reject[val.ResvTeamAvatarRes](app_res.BadRequest, "team_id 和 file_extension 不能为空")
	}

	var key string

	re, err := repo_iface.RunWithTxn[app_res.AppRes[val.ResvTeamAvatarRes]](a.txnCtrl, func(prov repo_iface.Prov) (app_res.AppRes[val.ResvTeamAvatarRes], error) {
		memberRepo := prov.MemberRepo()
		teamRepo := prov.TeamRepo()
		ossMsgRepo := prov.OssMsgRepo()

		if re := a.teamSvc.CanAdminTeam(currUid, args.TeamId, memberRepo, a.errClsf); re.IsReject() {
			return app_res.Reject[val.ResvTeamAvatarRes](app_res.ErrCode(re.Code()), re.Msg()), app_res.DefErr()
		}

		team, err := teamRepo.GetById(args.TeamId)
		if err != nil {
			if repo_infra.IsNotFound(err) {
				return app_res.Reject[val.ResvTeamAvatarRes](app_res.BadRequest, "团队不存在"), app_res.DefErr()
			}

			return app_res.Reject[val.ResvTeamAvatarRes](app_res.ServerError, "生成团队头像上传信息失败"), err
		}

		key = team.GenAvatarKey(args.FileExt)
		oldKey := team.AvatarKey

		if err := teamRepo.PrefillAvatarKey(args.TeamId, key); err != nil {
			if repo_infra.IsNotFound(err) {
				return app_res.Reject[val.ResvTeamAvatarRes](app_res.BadRequest, "团队不存在"), app_res.DefErr()
			}

			return app_res.Reject[val.ResvTeamAvatarRes](app_res.ServerError, "生成团队头像上传信息失败"), err
		}

		if oldKey != "" && oldKey != key {
			if err := a.ossMsgSvc.SavePendingDel(ossMsgRepo, enum.OssResTeamAvatar, args.TeamId, []string{oldKey}); err != nil {
				return app_res.Reject[val.ResvTeamAvatarRes](app_res.ServerError, "生成团队头像上传信息失败"), err
			}
		}

		if err := a.ossMsgSvc.SavePendingCre(ossMsgRepo, enum.OssResTeamAvatar, args.TeamId, []string{key}); err != nil {
			return app_res.Reject[val.ResvTeamAvatarRes](app_res.ServerError, "生成团队头像上传信息失败"), err
		}

		return app_res.Accept(&val.ResvTeamAvatarRes{}), nil
	})
	if err != nil {
		lgr.Error("[teamAppImpl.ResvAvatar] failed to run transaction", zap.Error(err))

		return re
	}

	url, err := a.ossSigner.GenPutUrl(key)
	if err != nil {
		lgr.Error("[teamAppImpl.ResvAvatar] failed to generate avatar upload url", zap.Error(err))

		return app_res.Reject[val.ResvTeamAvatarRes](app_res.ServerError, "生成团队头像上传信息失败")
	}

	return app_res.Accept(&val.ResvTeamAvatarRes{PutUrl: url})
}

// `MarkAvatarUploaded` marks team avatar upload as completed.
func (a *teamAppImpl) MarkAvatarUploaded(cx context.Context, currUid string, teamId string) app_res.AppRes[app_res.None] {
	lgr := app_util.TakeLgr(cx)

	if teamId == "" {
		return app_res.Reject[app_res.None](app_res.BadRequest, "team_id 不能为空")
	}

	re, err := repo_iface.RunWithTxn[app_res.AppRes[app_res.None]](a.txnCtrl, func(prov repo_iface.Prov) (app_res.AppRes[app_res.None], error) {
		memberRepo := prov.MemberRepo()
		teamRepo := prov.TeamRepo()
		ossMsgRepo := prov.OssMsgRepo()

		if re := a.teamSvc.CanAdminTeam(currUid, teamId, memberRepo, a.errClsf); re.IsReject() {
			return app_res.Reject[app_res.None](app_res.ErrCode(re.Code()), re.Msg()), app_res.DefErr()
		}

		if err := teamRepo.MarkAvatarUploaded(teamId); err != nil {
			if repo_infra.IsNotFound(err) {
				return app_res.Reject[app_res.None](app_res.BadRequest, "团队不存在"), app_res.DefErr()
			}

			return app_res.Reject[app_res.None](app_res.ServerError, "标记团队头像上传状态失败"), err
		}

		if err := ossMsgRepo.MarkCompletedByRes(enum.OssResTeamAvatar, teamId); err != nil {
			if repo_infra.IsNotFound(err) {
				return app_res.Reject[app_res.None](app_res.BadRequest, "无效的团队头像上传状态"), app_res.DefErr()
			}

			return app_res.Reject[app_res.None](app_res.ServerError, "标记团队头像上传状态失败"), err
		}

		return app_res.Accept(&app_res.None{}), nil
	})
	if err != nil {
		lgr.Error("[teamAppImpl.MarkAvatarUploaded] failed to run transaction", zap.Error(err))

		return re
	}

	return re
}
