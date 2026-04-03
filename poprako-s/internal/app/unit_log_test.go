package app

import (
	"poprako-s/internal/app/val"
	"testing"
)

func TestNewLogUnitAppRejectsNilDependencies(t *testing.T) {
	assertPanics(t, func() { NewLogUnitApp(nil) })
}

func TestLogUnitAppForwardsAllMethods(t *testing.T) {
	stub := &unitAppStub{}
	app := NewLogUnitApp(stub)
	_, _ = app.List(background(), "user-1", "page-1")
	_ = app.Save(background(), "user-1", &val.SavePageUnitArgs{PageID: "page-1"})
	if !stub.listCalled || !stub.saveCalled {
		t.Fatalf("expected all unit wrapper calls to forward: %#v", stub)
	}
}

func TestLogUnitAppRejectsInvalidArgs(t *testing.T) {
	app := NewLogUnitApp(&unitAppStub{})
	if _, err := app.List(background(), "user-1", ""); err == nil {
		t.Fatal("expected list validation error")
	}
	if err := app.Save(background(), "user-1", &val.SavePageUnitArgs{}); err == nil {
		t.Fatal("expected save validation error")
	}
}
