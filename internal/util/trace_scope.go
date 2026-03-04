package util

import (
	"go.uber.org/zap"
)

// 追踪链路的上下文
type TraceScope struct {
	logger *zap.Logger
}

func NewTraceScope(logger *zap.Logger) *TraceScope {
	return &TraceScope{
		logger: logger,
	}
}

func (ts *TraceScope) Logger() *zap.Logger {
	if ts.logger == nil {
		return zap.L()
	}

	return ts.logger
}

func (ts *TraceScope) WithFields(fields ...zap.Field) *TraceScope {
	ts.logger = ts.Logger().With(fields...)

	return ts
}
