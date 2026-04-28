package app_util

import (
	"context"
	"time"

	"poprako-s/internal/app/res"
	repo_iface "poprako-s/internal/domain/repo"
	repo_infra "poprako-s/internal/infra/repo"

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
func ClampOffsetLimit(offset int, limit *int) (res.ErrCode, string, bool) {
	if offset < 0 {
		return res.BadRequest, "offset 不能小于 0", true
	}

	if *limit <= 0 {
		*limit = listDefLimit
	}

	if *limit > listMaxLimit {
		*limit = listMaxLimit
	}

	return 0, "", false
}

// `ClassifyRepoErr` maps a repo error to the standard (code, msg, reject) triple.
func ClassifyRepoErr(err repo_iface.RepoErr, notFoundCode res.ErrCode, notFoundMsg, timeoutMsg, unavailableMsg, serverMsg string) (res.ErrCode, string, bool) {
	if repo_infra.IsNotFound(err) {
		return notFoundCode, notFoundMsg, true
	}

	if repo_infra.IsTimeout(err) {
		return res.Timeout, timeoutMsg, true
	}

	if repo_infra.IsUnavailable(err) {
		return res.Unavailable, unavailableMsg, true
	}

	return res.ServerError, serverMsg, true
}
