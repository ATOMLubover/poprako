package app

import (
	"poprako-s/internal/app/val"
	"testing"
)

func TestNewLogChapterAppRejectsNilDependencies(t *testing.T) {
	assertPanics(t, func() { NewLogChapterApp(nil) })
}

func TestLogChapterAppForwardsAllMethods(t *testing.T) {
	stub := &chapterAppStub{}
	app := NewLogChapterApp(stub)
	_, _ = app.List(background(), "user-1", &val.ListChapterArgs{ComicID: "comic-1"})
	_, _ = app.Create(background(), "user-1", &val.CreateChapterArgs{ComicID: "comic-1"})
	_ = app.Update(background(), "user-1", &val.UpdateChapterArgs{ChapterID: "chapter-1"})
	_ = app.Remove(background(), "user-1", "chapter-1")
	if !stub.listCalled || !stub.createCalled || !stub.updateCalled || !stub.removeCalled {
		t.Fatalf("expected all chapter wrapper calls to forward: %#v", stub)
	}
}

func TestLogChapterAppRejectsInvalidArgs(t *testing.T) {
	app := NewLogChapterApp(&chapterAppStub{})
	if _, err := app.List(background(), "user-1", nil); err == nil {
		t.Fatal("expected list validation error")
	}
	if _, err := app.Create(background(), "user-1", &val.CreateChapterArgs{}); err == nil {
		t.Fatal("expected create validation error")
	}
	if err := app.Update(background(), "user-1", &val.UpdateChapterArgs{}); err == nil {
		t.Fatal("expected update validation error")
	}
	if err := app.Remove(background(), "user-1", ""); err == nil {
		t.Fatal("expected remove validation error")
	}
}
