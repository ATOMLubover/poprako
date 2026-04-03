package app

import (
	"poprako-s/internal/app/val"
	"testing"
)

func TestNewLogInvitationAppRejectsNilDependencies(t *testing.T) {
	assertPanics(t, func() { NewLogInvitationApp(nil) })
}

func TestLogInvitationAppForwardsAllMethods(t *testing.T) {
	stub := &invitationAppStub{}
	app := NewLogInvitationApp(stub)
	_, _ = app.List(background(), "user-1", &val.ListTeamInvitationArgs{TeamID: "team-1"})
	_, _ = app.Create(background(), "user-1", &val.CreateInvitationArgs{TeamID: "team-1", InviteeQQ: "100002"})
	_ = app.Update(background(), "user-1", &val.UpdateInvitationArgs{ID: "inv-1", TeamID: "team-1"})
	_ = app.Remove(background(), "user-1", "inv-1")
	if !stub.listCalled || !stub.createCalled || !stub.updateCalled || !stub.removeCalled {
		t.Fatalf("expected all invitation wrapper calls to forward: %#v", stub)
	}
}

func TestLogInvitationAppRejectsInvalidArgs(t *testing.T) {
	app := NewLogInvitationApp(&invitationAppStub{})
	if _, err := app.List(background(), "user-1", nil); err == nil {
		t.Fatal("expected list validation error")
	}
	if _, err := app.Create(background(), "user-1", &val.CreateInvitationArgs{}); err == nil {
		t.Fatal("expected create validation error")
	}
	if err := app.Update(background(), "user-1", &val.UpdateInvitationArgs{}); err == nil {
		t.Fatal("expected update validation error")
	}
	if err := app.Remove(background(), "user-1", ""); err == nil {
		t.Fatal("expected remove validation error")
	}
}
