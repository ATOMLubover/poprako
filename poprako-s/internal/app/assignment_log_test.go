package app

import (
	"poprako-s/internal/app/val"
	"poprako-s/internal/domain/model"
	"testing"
)

func TestNewLogAssignmentAppRejectsNilDependencies(t *testing.T) {
	assertPanics(t, func() { NewLogAssignmentApp(nil) })
}

func TestLogAssignmentAppForwardsAllMethods(t *testing.T) {
	stub := &assignmentAppStub{}
	app := NewLogAssignmentApp(stub)
	_, _ = app.ListByChapter(background(), "user-1", &val.ListChapterAssignmentArgs{ChapterID: "chapter-1"})
	_, _ = app.ListMy(background(), "user-1", &val.ListMyAssignmentArgs{})
	_, _ = app.Create(background(), "user-1", &val.CreateAssignmentArgs{ChapterID: "chapter-1", UserID: "user-2", Roles: model.RoleMask(model.RoleTranslator)})
	_ = app.Update(background(), "user-1", &val.UpdateAssignmentArgs{ID: "assignment-1", Roles: model.RoleMask(model.RoleReviewer)})
	_ = app.Remove(background(), "user-1", "assignment-1")
	if !stub.listByChapterCalled || !stub.listMyCalled || !stub.createCalled || !stub.updateCalled || !stub.removeCalled {
		t.Fatalf("expected all assignment wrapper calls to forward: %#v", stub)
	}
}

func TestLogAssignmentAppRejectsInvalidArgs(t *testing.T) {
	app := NewLogAssignmentApp(&assignmentAppStub{})
	if _, err := app.ListByChapter(background(), "user-1", nil); err == nil {
		t.Fatal("expected list by chapter validation error")
	}
	if _, err := app.ListMy(background(), "user-1", nil); err == nil {
		t.Fatal("expected list my validation error")
	}
	if _, err := app.Create(background(), "user-1", &val.CreateAssignmentArgs{}); err == nil {
		t.Fatal("expected create validation error")
	}
	if err := app.Update(background(), "user-1", &val.UpdateAssignmentArgs{}); err == nil {
		t.Fatal("expected update validation error")
	}
	if err := app.Remove(background(), "user-1", ""); err == nil {
		t.Fatal("expected remove validation error")
	}
}
