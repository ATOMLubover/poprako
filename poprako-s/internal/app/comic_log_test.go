package app

import (
	"poprako-s/internal/app/val"
	"testing"
)

func TestNewLogComicAppRejectsNilDependencies(t *testing.T) {
	assertPanics(t, func() { NewLogComicApp(nil) })
}

func TestLogComicAppForwardsAllMethods(t *testing.T) {
	stub := &comicAppStub{}
	app := NewLogComicApp(stub)
	_, _ = app.List(background(), "user-1", &val.ListComicArgs{WorksetID: "workset-1"})
	_, _ = app.Create(background(), "user-1", &val.CreateComicArgs{WorksetID: "workset-1", Title: "Comic"})
	_ = app.Update(background(), "user-1", &val.UpdateComicArgs{ID: "comic-1"})
	_ = app.Remove(background(), "user-1", "comic-1")
	if !stub.listCalled || !stub.createCalled || !stub.updateCalled || !stub.removeCalled {
		t.Fatalf("expected all comic wrapper calls to forward: %#v", stub)
	}
}

func TestLogComicAppRejectsInvalidArgs(t *testing.T) {
	app := NewLogComicApp(&comicAppStub{})
	if _, err := app.List(background(), "user-1", nil); err == nil {
		t.Fatal("expected list validation error")
	}
	if _, err := app.Create(background(), "user-1", &val.CreateComicArgs{}); err == nil {
		t.Fatal("expected create validation error")
	}
	if err := app.Update(background(), "user-1", &val.UpdateComicArgs{}); err == nil {
		t.Fatal("expected update validation error")
	}
	if err := app.Remove(background(), "user-1", ""); err == nil {
		t.Fatal("expected remove validation error")
	}
}
