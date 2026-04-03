package app

import (
	"context"
	"testing"

	"go.uber.org/zap"
)

func TestRetrieveLgrUsesInjectedLogger(t *testing.T) {
	defaultLgr := retrieveLgr(context.Background())
	if defaultLgr == nil {
		t.Fatal("expected default logger")
	}

	injected := zap.NewNop()
	if got := retrieveLgr(injectLgr(context.Background(), injected)); got != injected {
		t.Fatal("expected injected logger")
	}
}

func TestRetrieveLgrFallsBackToGlobalWhenContextNil(t *testing.T) {
	if got := retrieveLgr(nil); got == nil {
		t.Fatal("expected fallback logger for nil context")
	}
}
