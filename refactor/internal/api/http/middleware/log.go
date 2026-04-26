package middleware

import (
	"time"

	"poprako-s/internal/cfg"

	"github.com/kataras/iris/v12"
	"go.uber.org/zap"
)

func LogLatency(cf *cfg.AppCfg) iris.Handler {
	if cf.Env == cfg.EnvProd {
		return func(cx iris.Context) {
			cx.Next()
		}
	}

	// Only record latency in non-prod environment, to avoid performance overhead.
	return func(cx iris.Context) {
		start := time.Now()
		cx.Next()

		latency := time.Since(start)
		method := cx.Method()
		path := cx.Path()
		status := cx.GetStatusCode()

		zap.L().Debug(
			"[LogLatency]",
			zap.String("method", method),
			zap.String("path", path),
			zap.Int("status", status),
			zap.Duration("latency", latency),
		)
	}
}
