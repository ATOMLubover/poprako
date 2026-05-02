package app_impl

import (
	"context"

	app_iface "poprako-s/internal/app"
	app_res "poprako-s/internal/app/res"
	app_util "poprako-s/internal/app/util"
	"poprako-s/internal/app/val"
	"poprako-s/internal/domain/model/query"
	repo_iface "poprako-s/internal/domain/repo"
	repo_infra "poprako-s/internal/infra/repo"

	"go.uber.org/zap"
)

// `sysMailAppImpl` is the default implementation of `SysMailApp`.
type sysMailAppImpl struct {
	sysMailRepo repo_iface.SysMailRepo
}

// `NewSysMailApp` creates a `SysMailApp` implementation.
func NewSysMailApp(sysMailRepo repo_iface.SysMailRepo) app_iface.SysMailApp {
	if sysMailRepo == nil {
		zap.L().Panic(
			"[NewSysMailApp] nil dependency",
			zap.Bool("sysMailRepo", sysMailRepo == nil),
		)
	}

	return &sysMailAppImpl{sysMailRepo: sysMailRepo}
}

// `List` returns unread system mails for current user with pagination.
func (a *sysMailAppImpl) List(cx context.Context, currUid string, args *val.ListSysMailArgs) app_res.AppRes[[]val.SysMailVal] {
	lgr := app_util.TakeLgr(cx)

	if re := vfyListSysMailArgs(args); re.IsReject() {
		return app_res.Reject[[]val.SysMailVal](re.Code(), re.Msg())
	}

	items, err := a.sysMailRepo.ListUnreadByRcvId(
		currUid,
		query.PagiOpt{Offset: args.Offset, Limit: args.Limit},
	)
	if err != nil {
		lgr.Error("[sysMailAppImpl.List] failed to list unread system mails", zap.Error(err))
		return app_res.Reject[[]val.SysMailVal](app_res.ServerError, "获取系统消息失败")
	}

	sysMailVals := make([]val.SysMailVal, len(items))

	for i := range items {
		sysMailVals[i] = asmSysMailVal(&items[i])
	}

	return app_res.Accept(&sysMailVals)
}

// `MarkRead` marks one system mail as read for current user.
func (a *sysMailAppImpl) MarkRead(cx context.Context, currUid string, id string) app_res.AppRes[app_res.None] {
	lgr := app_util.TakeLgr(cx)

	if re := vfyMarkReadSysMailId(id); re.IsReject() {
		return app_res.Reject[app_res.None](re.Code(), re.Msg())
	}

	err := a.sysMailRepo.MarkReadByRcvId(id, currUid)
	if err != nil {
		if repo_infra.IsNotFound(err) {
			return app_res.Reject[app_res.None](app_res.NotFound, "系统消息不存在")
		}
		lgr.Error("[sysMailAppImpl.MarkRead] failed to mark system mail as read", zap.String("sys_mail_id", id), zap.Error(err))
		return app_res.Reject[app_res.None](app_res.ServerError, "标记系统消息已读失败")
	}

	return app_res.Accept(&app_res.None{})
}
