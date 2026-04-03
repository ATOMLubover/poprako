package app

import (
	"poprako-s/internal/app/val"
	"testing"
)

func TestNewLogTeamAppRejectsNilDependencies(t *testing.T) {
	assertPanics(t, func() { NewLogTeamApp(nil) })
}

func TestLogTeamAppForwardsAllMethods(t *testing.T) {
	stub := &teamAppStub{}
	app := NewLogTeamApp(stub)
	_, _ = app.Create(background(), "user-1", &val.CreateTeamArgs{Name: "Team"})
	_, _ = app.List(background(), "user-1", &val.ListTeamArgs{})
	_, _ = app.ListMy(background(), "user-1", &val.ListMyTeamArgs{})
	_ = app.Update(background(), "user-1", &val.UpdateTeamArgs{ID: "team-1", Name: "Team"})
	_ = app.Remove(background(), "user-1", "team-1")
	_, _ = app.ReserveAvatar(background(), "user-1", "team-1")
	_ = app.ConfirmAvatarUploaded(background(), "user-1", "team-1")
	if !stub.createCalled || !stub.listCalled || !stub.listMyCalled || !stub.updateCalled || !stub.removeCalled || !stub.reserveCalled || !stub.confirmCalled {
		t.Fatalf("expected all team wrapper calls to forward: %#v", stub)
	}
}

func TestLogTeamAppRejectsInvalidArgs(t *testing.T) {
	app := NewLogTeamApp(&teamAppStub{})
	if _, err := app.Create(background(), "user-1", &val.CreateTeamArgs{}); err == nil {
		t.Fatal("expected create validation error")
	}
	if err := app.Update(background(), "user-1", &val.UpdateTeamArgs{}); err == nil {
		t.Fatal("expected update validation error")
	}
	if err := app.Remove(background(), "user-1", ""); err == nil {
		t.Fatal("expected remove validation error")
	}
	if _, err := app.ReserveAvatar(background(), "user-1", ""); err == nil {
		t.Fatal("expected reserve avatar validation error")
	}
	if err := app.ConfirmAvatarUploaded(background(), "user-1", ""); err == nil {
		t.Fatal("expected confirm avatar validation error")
	}
}
