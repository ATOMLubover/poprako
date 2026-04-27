package app_impl

import (
	"context"

	"go.uber.org/zap"
)

const LgrKey = "logger"

func takeLgr(cx context.Context) *zap.Logger {
	if cx == nil {
		return zap.L()
	}

	lgr, ok := cx.Value(LgrKey).(*zap.Logger)
	if !ok || lgr == nil {
		return zap.L()
	}

	return lgr
}

func saveLgr(cx context.Context, lgr *zap.Logger) context.Context {
	return context.WithValue(cx, LgrKey, lgr)
}
