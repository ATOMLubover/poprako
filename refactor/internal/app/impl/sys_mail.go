package app_impl

import (
	"context"

	app_iface "poprako-s/internal/app"
	"poprako-s/internal/app/res"
	app_util "poprako-s/internal/app/util"
	"poprako-s/internal/app/val"
	"poprako-s/internal/domain/model/query"
	repo_iface "poprako-s/internal/domain/repo"

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
func (a *sysMailAppImpl) List(cx context.Context, currUid string, args *val.ListSysMailArgs) res.AppRes[[]val.SysMailVal] {
	lgr := app_util.TakeLgr(cx)

	if code, msg, reject := vfyListSysMailArgs(args); reject {
		return res.Reject[[]val.SysMailVal](code, msg)
	}

	items, err := a.sysMailRepo.ListUnreadByRcvId(
		currUid,
		query.PagiOpt{Offset: args.Offset, Limit: args.Limit},
	)
	if err != nil {
		code, msg, _ := app_util.ClassifyRepoErr(err, 0, "", "获取系统消息超时", "系统消息服务暂不可用", "获取系统消息失败")
		lgr.Error("[sysMailAppImpl.List] failed to list unread system mails", zap.Error(err))
		return res.Reject[[]val.SysMailVal](code, msg)
	}

	vals := make([]val.SysMailVal, len(items))

	for i := range items {
		vals[i] = asmSysMailVal(&items[i])
	}

	return res.Accept(&vals)
}

// `MarkRead` marks one system mail as read for current user.
func (a *sysMailAppImpl) MarkRead(cx context.Context, currUid string, id string) res.AppRes[res.None] {
	lgr := app_util.TakeLgr(cx)

	if code, msg, reject := vfyMarkReadSysMailId(id); reject {
		return res.Reject[res.None](code, msg)
	}

	err := a.sysMailRepo.MarkReadByRcvId(id, currUid)
	if err != nil {
		code, msg, _ := app_util.ClassifyRepoErr(err, res.NotFound, "系统消息不存在", "标记系统消息已读超时", "系统消息服务暂不可用", "标记系统消息已读失败")
		lgr.Error("[sysMailAppImpl.MarkRead] failed to mark system mail as read", zap.String("sys_mail_id", id), zap.Error(err))
		return res.Reject[res.None](code, msg)
	}

	return res.Accept(&res.None{})
}
