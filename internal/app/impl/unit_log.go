package app_impl

import (
	"context"

	app_iface "poprako-s/internal/app"
	app_res "poprako-s/internal/app/res"
	app_util "poprako-s/internal/app/util"
	"poprako-s/internal/app/val"

	"go.uber.org/zap"
)

// `unitLogAppImpl` is a logging decorator for `UnitApp`.
type unitLogAppImpl struct {
	// `inner` is the wrapped unit application.
	inner app_iface.UnitApp
}

// `NewUnitLogApp` creates the logging decorator for `UnitApp`.
func NewUnitLogApp(inner app_iface.UnitApp) app_iface.UnitApp {
	if inner == nil {
		zap.L().Panic("[NewUnitLogApp] nil dependency", zap.Bool("inner", inner == nil))
	}

	return &unitLogAppImpl{inner: inner}
}

// `ListByPage` enriches logger context then forwards the call.
func (a *unitLogAppImpl) ListByPage(cx context.Context, currUid string, args *val.ListPageUnitsArgs) app_res.AppRes[val.ListPageUnitsRes] {
	if cx == nil {
		cx = context.Background()
	}

	if args == nil {
		return app_res.Reject[val.ListPageUnitsRes](app_res.BadRequest, "page_id 不能为空")
	}

	lgr := app_util.TakeLgr(cx)
	lgr = lgr.With(zap.String("curr_uid", currUid), zap.String("page_id", args.PageId))
	cx = app_util.SaveLgr(cx, lgr)

	return a.inner.ListByPage(cx, currUid, args)
}

// `SaveByPage` enriches logger context then forwards the call.
func (a *unitLogAppImpl) SaveByPage(cx context.Context, currUid string, args *val.SavePageUnitsArgs) app_res.AppRes[val.SavePageUnitsRes] {
	if cx == nil {
		cx = context.Background()
	}

	if args == nil {
		return app_res.Reject[val.SavePageUnitsRes](app_res.BadRequest, "page_id 不能为空")
	}

	lgr := app_util.TakeLgr(cx)
	lgr = lgr.With(zap.String("curr_uid", currUid), zap.String("page_id", args.PageId))

	cx = app_util.SaveLgr(cx, lgr)

	lgr.Debug("[unitLogAppImpl.SaveByPage] Saving page units")

	return a.inner.SaveByPage(cx, currUid, args)
}
