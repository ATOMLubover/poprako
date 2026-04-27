package app_util

import (
	"context"

	"go.uber.org/zap"
)

const lgrKey = "logger"

func TakeLgr(cx context.Context) *zap.Logger {
	if cx == nil {
		return zap.L()
	}

	lgr, ok := cx.Value(lgrKey).(*zap.Logger)
	if !ok || lgr == nil {
		return zap.L()
	}

	return lgr
}

func SaveLgr(cx context.Context, lgr *zap.Logger) context.Context {
	return context.WithValue(cx, lgrKey, lgr)
}
