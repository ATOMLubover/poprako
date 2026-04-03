package app

import (
	"poprako-s/internal/app/val"
	"testing"
)

func TestNewLogPageAppRejectsNilDependencies(t *testing.T) {
	assertPanics(t, func() { NewLogPageApp(nil) })
}

func TestLogPageAppForwardsAllMethods(t *testing.T) {
	stub := &pageAppStub{}
	app := NewLogPageApp(stub)
	_, _ = app.Reserve(background(), "user-1", &val.ReserveChapterPagesArgs{ChapterID: "chapter-1", PageCount: 1})
	_, _ = app.List(background(), "user-1", &val.ListChapterPageArgs{ChapterID: "chapter-1"})
	_ = app.Update(background(), "user-1", &val.UpdatePageArgs{ID: "page-1"})
	_ = app.Remove(background(), "user-1", "page-1")
	if !stub.reserveCalled || !stub.listCalled || !stub.updateCalled || !stub.removeCalled {
		t.Fatalf("expected all page wrapper calls to forward: %#v", stub)
	}
}

func TestLogPageAppRejectsInvalidArgs(t *testing.T) {
	app := NewLogPageApp(&pageAppStub{})
	if _, err := app.Reserve(background(), "user-1", &val.ReserveChapterPagesArgs{}); err == nil {
		t.Fatal("expected reserve validation error")
	}
	if _, err := app.List(background(), "user-1", nil); err == nil {
		t.Fatal("expected list validation error")
	}
	if err := app.Update(background(), "user-1", &val.UpdatePageArgs{}); err == nil {
		t.Fatal("expected update validation error")
	}
	if err := app.Remove(background(), "user-1", ""); err == nil {
		t.Fatal("expected remove validation error")
	}
}
