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
	_, _ = app.Get(background(), "user-1", "comic-1")
	_, _ = app.List(background(), "user-1", &val.ListComicArgs{WorksetID: "workset-1"})
	_, _ = app.Create(background(), "user-1", &val.CreateComicArgs{WorksetID: "workset-1", Title: "Comic"})
	_ = app.Update(background(), "user-1", &val.UpdateComicArgs{ID: "comic-1"})
	_ = app.Remove(background(), "user-1", "comic-1")
	_, _ = app.ReserveCover(background(), "user-1", &val.ReserveComicCoverArgs{ComicID: "comic-1", FileName: "cover.jpg"})
	_ = app.ConfirmCoverUploaded(background(), "user-1", "comic-1")
	if !stub.getCalled || !stub.listCalled || !stub.createCalled || !stub.updateCalled || !stub.removeCalled || !stub.reserveCoverCalled || !stub.confirmCoverCalled {
		t.Fatalf("expected all comic wrapper calls to forward: %#v", stub)
	}
}

func TestLogComicAppRejectsInvalidArgs(t *testing.T) {
	app := NewLogComicApp(&comicAppStub{})
	if _, err := app.Get(background(), "user-1", ""); err == nil {
		t.Fatal("expected get validation error")
	}
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
	if _, err := app.ReserveCover(background(), "user-1", &val.ReserveComicCoverArgs{}); err == nil {
		t.Fatal("expected reserve cover validation error")
	}
	if _, err := app.ReserveCover(background(), "user-1", &val.ReserveComicCoverArgs{ComicID: "comic-1"}); err == nil {
		t.Fatal("expected reserve cover file_name validation error")
	}
	if err := app.ConfirmCoverUploaded(background(), "user-1", ""); err == nil {
		t.Fatal("expected confirm cover validation error")
	}
}
