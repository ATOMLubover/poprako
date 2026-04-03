package app

import (
	"poprako-s/internal/app/val"
	"testing"
)

func TestNewLogWorksetAppRejectsNilDependencies(t *testing.T) {
	assertPanics(t, func() { NewLogWorksetApp(nil) })
}

func TestLogWorksetAppForwardsAllMethods(t *testing.T) {
	stub := &worksetAppStub{}
	app := NewLogWorksetApp(stub)
	_, _ = app.List(background(), "user-1", &val.ListWorksetArgs{TeamID: "team-1"})
	_, _ = app.Create(background(), "user-1", &val.CreateWorksetArgs{TeamID: "team-1", Name: "Workset"})
	_ = app.Update(background(), "user-1", &val.UpdateWorksetArgs{ID: "workset-1", Name: "Workset"})
	_ = app.Remove(background(), "user-1", "workset-1")
	if !stub.listCalled || !stub.createCalled || !stub.updateCalled || !stub.removeCalled {
		t.Fatalf("expected all workset wrapper calls to forward: %#v", stub)
	}
}

func TestLogWorksetAppRejectsInvalidArgs(t *testing.T) {
	app := NewLogWorksetApp(&worksetAppStub{})
	if _, err := app.List(background(), "user-1", nil); err == nil {
		t.Fatal("expected list validation error")
	}
	if _, err := app.Create(background(), "user-1", &val.CreateWorksetArgs{}); err == nil {
		t.Fatal("expected create validation error")
	}
	if err := app.Update(background(), "user-1", &val.UpdateWorksetArgs{}); err == nil {
		t.Fatal("expected update validation error")
	}
	if err := app.Remove(background(), "user-1", ""); err == nil {
		t.Fatal("expected remove validation error")
	}
}
