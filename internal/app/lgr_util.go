package app

import (
	"context"

	"go.uber.org/zap"
)

const lgrKey = "lgr"

func injectLgr(cx context.Context, lgr *zap.Logger) context.Context {
	return context.WithValue(cx, lgrKey, lgr)
}

func retrieveLgr(cx context.Context) *zap.Logger {
	if cx == nil {
		return zap.L()
	}

	lgr, ok := cx.Value(lgrKey).(*zap.Logger)
	if !ok {
		return zap.L()
	}

	return lgr
}
