package app_impl

import (
	"context"

	app_iface "poprako-s/internal/app"
	"poprako-s/internal/app/res"
	app_util "poprako-s/internal/app/util"
	"poprako-s/internal/app/val"

	"go.uber.org/zap"
)

// `userStatsLogAppImpl` is logging decorator for `UserStatsApp`.
type userStatsLogAppImpl struct {
	inner app_iface.UserStatsApp
}

// `NewUserStatsLogApp` creates logging decorator for `UserStatsApp`.
func NewUserStatsLogApp(inner app_iface.UserStatsApp) app_iface.UserStatsApp {
	if inner == nil {
		zap.L().Panic("[NewUserStatsLogApp] nil dependency", zap.Bool("inner", inner == nil))
	}

	return &userStatsLogAppImpl{inner: inner}
}

// `GetStats` enriches logger context and forwards call.
func (a *userStatsLogAppImpl) GetStats(cx context.Context, currUid string) app_res.AppRes[val.UserStats] {
	if cx == nil {
		cx = context.Background()
	}

	lgr := app_util.TakeLgr(cx).With(zap.String("curr_uid", currUid))
	cx = app_util.SaveLgr(cx, lgr)

	return a.inner.GetStats(cx, currUid)
}
