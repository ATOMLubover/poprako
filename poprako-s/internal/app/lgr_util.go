package app

import (
	"context"

	"go.uber.org/zap"
)

func injectLgr(cx context.Context, lgr *zap.Logger) context.Context {
	return context.WithValue(cx, "lgr", lgr)
}

func retrieveLgr(cx context.Context) *zap.Logger {
	if cx == nil {
		return zap.L()
	}

	lgr, ok := cx.Value("lgr").(*zap.Logger)
	if !ok {
		return zap.L()
	}

	return lgr
}
