package app_impl

import (
	"context"
	"fmt"

	app_iface "poprako-s/internal/app"
	"poprako-s/internal/app/res"
	app_util "poprako-s/internal/app/util"
	"poprako-s/internal/app/val"
	"poprako-s/internal/domain/model/aggr"
	"poprako-s/internal/domain/model/enum"
	"poprako-s/internal/domain/model/query"
	repo_iface "poprako-s/internal/domain/repo"
	"poprako-s/internal/domain/svc"
	repo_infra "poprako-s/internal/infra/repo"

	"go.uber.org/zap"
)

// `comicAppImpl` is the default implementation of `ComicApp`
type comicAppImpl struct {
	txnCtrl repo_iface.TxnCtrl

	comicSvc svc.ComicSvc

	memberRepo  repo_iface.MemberRepo
	worksetRepo repo_iface.WorksetRepo
	comicRepo   repo_iface.ComicRepo
}

// `NewComicApp` creates a `ComicApp` implementation
func NewComicApp(
	txnCtrl repo_iface.TxnCtrl,
	comicSvc svc.ComicSvc,
	memberRepo repo_iface.MemberRepo,
	worksetRepo repo_iface.WorksetRepo,
	comicRepo repo_iface.ComicRepo,
) app_iface.ComicApp {
	if txnCtrl == nil ||
		memberRepo == nil ||
		worksetRepo == nil ||
		comicRepo == nil {
		zap.L().Panic(
			"[NewComicApp] nil dependency",
			zap.Bool("txnCtrl", txnCtrl == nil),
			zap.Bool("memberRepo", memberRepo == nil),
			zap.Bool("worksetRepo", worksetRepo == nil),
			zap.Bool("comicRepo", comicRepo == nil),
		)
	}

	return &comicAppImpl{
		txnCtrl:     txnCtrl,
		comicSvc:    comicSvc,
		memberRepo:  memberRepo,
		worksetRepo: worksetRepo,
		comicRepo:   comicRepo,
	}
}

// `List` returns all active comics for one workset
func (a *comicAppImpl) List(cx context.Context, currUid string, args *val.ListComicArgs) res.AppRes[[]val.ComicVal] {
	lgr := app_util.TakeLgr(cx)

	if code, msg, reject := vfyListComicArgs(args); reject {
		return res.Reject[[]val.ComicVal](code, msg)
	}

	ws, err := a.worksetRepo.GetById(args.WorksetId)
	if err != nil {
		return res.Reject[[]val.ComicVal](res.BadRequest, "作品集不存在")
	}

	ok, err := a.memberRepo.ExistByUserTeamId(currUid, ws.TeamId)
	if err != nil {
		lgr.Error(
			"[comicAppImpl.List] failed to verify team membership",
			zap.Error(err),
		)

		return res.Reject[[]val.ComicVal](res.ServerError, "获取漫画列表失败")
	}
	if !ok {
		return res.Reject[[]val.ComicVal](res.Forbidden, "无权访问该作品集的漫画")
	}

	listOpt := &query.ListComicOpt{
		WorksetId:      &args.WorksetId,
		UploadPhase:    args.UploadPhase,
		TranslatePhase: args.TranslatePhase,
		ProofreadPhase: args.ProofreadPhase,
		TypesetPhase:   args.TypesetPhase,
		ReviewPhase:    args.ReviewPhase,
		PublishPhase:   args.PublishPhase,
		Pagi: query.PagiOpt{
			Offset: args.Offset,
			Limit:  args.Limit,
		},
	}

	if args.FuzzyTitle != "" {
		listOpt.FuzzyTitle = &args.FuzzyTitle
	}

	comics, err := a.comicRepo.List(listOpt)
	if err != nil {
		lgr.Error(
			"[comicAppImpl.List] failed to list comics",
			zap.Error(err),
		)

		return res.Reject[[]val.ComicVal](res.ServerError, "获取漫画列表失败")
	}

	vals := make([]val.ComicVal, len(comics))
	for i, cm := range comics {
		vals[i] = asmComicVal(cm)
	}

	return res.Accept(&vals)
}

// `Create` creates a comic in target workset
func (a *comicAppImpl) Create(cx context.Context, currUid string, args *val.CreateComicArgs) res.AppRes[val.ComicCreatedRes] {
	lgr := app_util.TakeLgr(cx)

	if code, msg, reject := vfyCreateComicArgs(args); reject {
		return res.Reject[val.ComicCreatedRes](code, msg)
	}

	var (
		createdId string
		errCode   = res.BadRequest
	)

	if err := a.txnCtrl.RunWithTxn(func(cx context.Context) error {
		memberRepo, err := repo_infra.TxnMemberRepo(cx)
		if err != nil {
			errCode = res.ServerError
			return err
		}

		worksetRepo, err := repo_infra.TxnWorksetRepo(cx)
		if err != nil {
			errCode = res.ServerError
			return err
		}

		comicRepo, err := repo_infra.TxnComicRepo(cx)
		if err != nil {
			errCode = res.ServerError
			return err
		}

		ws, err := worksetRepo.GetById(args.WorksetId)
		if err != nil {
			return err
		}

		member, err := memberRepo.GetByUserTeamId(currUid, ws.TeamId)
		if err != nil {
			return err
		}

		if member == nil || !member.HasAnyRole(enum.RoleAdmin) {
			return fmt.Errorf("only team admin can create comic")
		}

		count, err := comicRepo.Count(&query.ListComicOpt{WorksetId: &args.WorksetId})
		if err != nil {
			errCode = res.ServerError
			return err
		}

		cre := a.comicSvc.NewComicCre(
			args.WorksetId,
			int(count),
			args.Title,
			args.Author,
			args.Desc,
			currUid,
		)

		cm, err := comicRepo.Create(cre)
		if err != nil {
			return err
		}

		if err := worksetRepo.UpdateComicCount(args.WorksetId, 1); err != nil {
			errCode = res.ServerError
			return err
		}

		createdId = cm.Id

		return nil
	}); err != nil {
		lgr.Error(
			"[comicAppImpl.Create] failed to run create comic transaction",
			zap.Error(err),
		)

		switch errCode {
		case res.ServerError:
			return res.Reject[val.ComicCreatedRes](res.ServerError, "创建漫画失败")
		default:
			return res.Reject[val.ComicCreatedRes](res.Forbidden, "仅汉化组管理员可创建漫画")
		}
	}

	return res.Accept(&val.ComicCreatedRes{Id: createdId})
}

// `Update` updates title author and description of a comic
func (a *comicAppImpl) Update(cx context.Context, currUid string, args *val.ComicUpdArgs) res.AppRes[res.None] {
	lgr := app_util.TakeLgr(cx)

	if code, msg, reject := vfyUpdateComicArgs(args); reject {
		return res.Reject[res.None](code, msg)
	}

	cm, err := a.comicRepo.GetById(args.Id)
	if err != nil {
		lgr.Error("[comicAppImpl.Update] failed to get comic", zap.Error(err))

		return res.Reject[res.None](res.BadRequest, "漫画不存在")
	}

	ws, err := a.worksetRepo.GetById(cm.WorksetId)
	if err != nil {
		lgr.Error("[comicAppImpl.Update] failed to get workset", zap.Error(err))

		return res.Reject[res.None](res.ServerError, "更新漫画失败")
	}

	member, err := a.memberRepo.GetByUserTeamId(currUid, ws.TeamId)
	if err != nil {
		lgr.Error("[comicAppImpl.Update] failed to get member", zap.Error(err))

		return res.Reject[res.None](res.ServerError, "更新漫画失败")
	}

	if member == nil || !member.HasAnyRole(enum.RoleAdmin) {
		return res.Reject[res.None](res.Forbidden, "仅汉化组管理员可更新漫画")
	}

	upd := &aggr.ComicUpd{
		Id:     args.Id,
		Title:  args.Title,
		Author: args.Author,
		Desc:   args.Desc,
	}

	if err := a.comicRepo.Update(upd); err != nil {
		lgr.Error("[comicAppImpl.Update] failed to update comic", zap.Error(err))

		return res.Reject[res.None](res.ServerError, "更新漫画失败")
	}

	return res.Accept(&res.None{})
}

// `Remove` soft-deletes a comic and updates workset comic counter
func (a *comicAppImpl) Remove(cx context.Context, currUid string, comicId string) res.AppRes[res.None] {
	lgr := app_util.TakeLgr(cx)

	if code, msg, reject := vfyRemoveComicId(comicId); reject {
		return res.Reject[res.None](code, msg)
	}

	var errCode = res.BadRequest

	if err := a.txnCtrl.RunWithTxn(func(cx context.Context) error {
		memberRepo, err := repo_infra.TxnMemberRepo(cx)
		if err != nil {
			errCode = res.ServerError
			return err
		}

		worksetRepo, err := repo_infra.TxnWorksetRepo(cx)
		if err != nil {
			errCode = res.ServerError
			return err
		}

		comicRepo, err := repo_infra.TxnComicRepo(cx)
		if err != nil {
			errCode = res.ServerError
			return err
		}

		cm, err := comicRepo.GetById(comicId)
		if err != nil {
			return err
		}

		ws, err := worksetRepo.GetById(cm.WorksetId)
		if err != nil {
			return err
		}

		member, err := memberRepo.GetByUserTeamId(currUid, ws.TeamId)
		if err != nil {
			return err
		}

		if member == nil || !member.HasAnyRole(enum.RoleAdmin) {
			return fmt.Errorf("only team admin can remove comic")
		}

		if err := comicRepo.Remove(comicId); err != nil {
			errCode = res.ServerError
			return err
		}

		if err := worksetRepo.UpdateComicCount(ws.Id, -1); err != nil {
			errCode = res.ServerError
			return err
		}

		return nil
	}); err != nil {
		lgr.Error(
			"[comicAppImpl.Remove] failed to run remove comic transaction",
			zap.Error(err),
		)

		switch errCode {
		case res.ServerError:
			return res.Reject[res.None](res.ServerError, "删除漫画失败")
		default:
			return res.Reject[res.None](res.Forbidden, "仅汉化组管理员可删除漫画")
		}
	}

	return res.Accept(&res.None{})
}
