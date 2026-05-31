package middleware

import (
	"time"

	"poprako-s/internal/cfg"

	"github.com/kataras/iris/v12"
	"go.uber.org/zap"
)

// `logLatencyWarnThreshold` is the latency threshold above which a warning is logged.
const logLatencyWarnThreshold = 50 * time.Millisecond

// `LogLatency` is an Iris middleware that records request latency
// In production (`EnvProd`) it is a no-op to avoid overhead
// In all other environments it logs method, path, status and duration at debug level
func LogLatency(cf *cfg.AppCfg) iris.Handler {
	if cf.Env() == cfg.EnvProd {
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

		// Warn when latency exceeds the threshold to surface slow requests.
		if latency > logLatencyWarnThreshold {
			zap.L().Warn(
				"[LogLatency] slow request",
				zap.String("method", method),
				zap.String("path", path),
				zap.Int("status", status),
				zap.Duration("latency", latency),
			)
		}
	}
}
