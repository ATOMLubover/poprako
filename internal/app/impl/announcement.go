package app_impl

import (
	"context"

	app_iface "poprako-s/internal/app"
	app_res "poprako-s/internal/app/res"
	app_util "poprako-s/internal/app/util"
	"poprako-s/internal/app/val"
	"poprako-s/internal/domain/model/enum"
	"poprako-s/internal/domain/model/query"
	repo_iface "poprako-s/internal/domain/repo"
	"poprako-s/internal/domain/svc"

	"go.uber.org/zap"
)

// `announcementAppImpl` is the default implementation of `AnnouncementApp`.
type announcementAppImpl struct {
	txnCtrl repo_iface.TxnCtrl

	announcementSvc svc.AnnouncementSvc

	memberRepo       repo_iface.MemberRepo
	announcementRepo repo_iface.AnnouncementRepo

	errClsf repo_iface.ErrClsf
}

// `NewAnnouncementApp` creates one `AnnouncementApp` implementation.
func NewAnnouncementApp(
	txnCtrl repo_iface.TxnCtrl,
	memberRepo repo_iface.MemberRepo,
	announcementRepo repo_iface.AnnouncementRepo,
	announcementSvc svc.AnnouncementSvc,
	errClsf repo_iface.ErrClsf,
) app_iface.AnnouncementApp {
	if txnCtrl == nil || memberRepo == nil || announcementRepo == nil || errClsf == nil {
		zap.L().Panic(
			"[NewAnnouncementApp] nil dependency",
			zap.Bool("txnCtrl", txnCtrl == nil),
			zap.Bool("memberRepo", memberRepo == nil),
			zap.Bool("announcementRepo", announcementRepo == nil),
			zap.Bool("errClsf", errClsf == nil),
		)
	}

	return &announcementAppImpl{
		txnCtrl:          txnCtrl,
		announcementSvc:  announcementSvc,
		memberRepo:       memberRepo,
		announcementRepo: announcementRepo,
		errClsf:          errClsf,
	}
}

// `List` lists announcements under one team.
func (a *announcementAppImpl) List(cx context.Context, currUid string, args *val.ListAnnouncementArgs) app_res.AppRes[[]val.AnnouncementVal] {
	lgr := app_util.TakeLgr(cx)

	if re := vfyListAnnouncementArgs(args); re.IsReject() {
		return app_res.Reject[[]val.AnnouncementVal](re.Code(), re.Msg())
	}

	if re := a.announcementSvc.CanListAnnouncement(currUid, args.TeamId, a.memberRepo, a.errClsf); re.IsReject() {
		return app_res.Reject[[]val.AnnouncementVal](app_res.ErrCode(re.Code()), re.Msg())
	}

	items, err := a.announcementRepo.List(
		query.ListAnnouncementOpt{
			TeamId: args.TeamId,
			Pagi: query.PagiOpt{
				Offset: args.Offset,
				Limit:  args.Limit,
			},
		},
		enum.AnnouncementInclUser,
	)
	if err != nil {
		lgr.Error("[announcementAppImpl.List] failed to list announcements", zap.Error(err))

		return app_res.Reject[[]val.AnnouncementVal](app_res.ServerError, "获取公告列表失败")
	}

	vals := make([]val.AnnouncementVal, len(items))
	for i := range items {
		vals[i] = asmAnnouncementVal(&items[i])
	}

	return app_res.Accept(&vals)
}

// `Create` creates one announcement under one team.
func (a *announcementAppImpl) Create(cx context.Context, currUid string, args *val.CreateAnnouncementArgs) app_res.AppRes[val.AnnouncementCreatedRes] {
	lgr := app_util.TakeLgr(cx)

	if re := vfyCreateAnnouncementArgs(args); re.IsReject() {
		return app_res.Reject[val.AnnouncementCreatedRes](re.Code(), re.Msg())
	}

	re, err := repo_iface.RunWithTxn[app_res.AppRes[val.AnnouncementCreatedRes]](a.txnCtrl, func(prov repo_iface.Prov) (app_res.AppRes[val.AnnouncementCreatedRes], error) {
		memberRepo := prov.MemberRepo()
		announcementRepo := prov.AnnouncementRepo()

		if re := a.announcementSvc.CanCreateAnnouncement(currUid, args.TeamId, memberRepo, a.errClsf); re.IsReject() {
			return app_res.Reject[val.AnnouncementCreatedRes](app_res.ErrCode(re.Code()), re.Msg()), app_res.DefErr()
		}

		announcementCre, err := a.announcementSvc.NewAnnouncementCre(currUid, args.TeamId, args.Title, args.Content)
		if err != nil {
			return app_res.Reject[val.AnnouncementCreatedRes](app_res.BadRequest, err.Error()), app_res.DefErr()
		}

		announcement, err := announcementRepo.Create(announcementCre)
		if err != nil {
			return app_res.Reject[val.AnnouncementCreatedRes](app_res.ServerError, "发布公告失败"), err
		}

		return app_res.Accept(&val.AnnouncementCreatedRes{Id: announcement.Id}), nil
	})
	if err != nil {
		lgr.Error("[announcementAppImpl.Create] failed to run transaction", zap.Error(err))

		return re
	}

	return re
}
