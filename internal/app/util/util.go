package app_util

import (
	"context"
	"time"

	app_res "poprako-s/internal/app/res"

	"go.uber.org/zap"
)

const (
	lgrKey       = "logger"
	listDefLimit = 20
	listMaxLimit = 200
)

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

func ToUnixMilliPtr(t *time.Time) *int64 {
	if t == nil {
		return nil
	}

	m := t.UnixMilli()

	return &m
}

// `ClampOffsetLimit` validates offset and normalizes limit to [defLimit, maxLimit].
func ClampOffsetLimit(offset int, limit *int) app_res.AppRes[app_res.None] {
	if offset < 0 {
		return app_res.Reject[app_res.None](app_res.BadRequest, "offset 不能小于 0")
	}

	if *limit <= 0 {
		*limit = listDefLimit
	}

	if *limit > listMaxLimit {
		*limit = listMaxLimit
	}

	return app_res.Accept(&app_res.None{})
}
