package app

import (
	"poprako-s/internal/app/val"
	"poprako-s/internal/domain/model"
	"testing"
)

func TestNewLogMemberAppRejectsNilDependencies(t *testing.T) {
	assertPanics(t, func() { NewLogMemberApp(nil) })
}

func TestLogMemberAppForwardsAllMethods(t *testing.T) {
	stub := &memberAppStub{}
	app := NewLogMemberApp(stub)
	_, _ = app.Create(background(), "user-1", &val.CreateMemberArgs{UserID: "user-2", TeamID: "team-1"})
	_, _ = app.ListByTeam(background(), "user-1", &val.ListTeamMemberArgs{TeamID: "team-1"})
	_, _ = app.ListMy(background(), "user-1", &val.ListMyMemberArgs{})
	_ = app.UpdateRole(background(), "user-1", &val.UpdateMemberRoleArgs{ID: "member-1", Roles: model.RoleMask(model.RoleAdmin)})
	_ = app.Remove(background(), "user-1", "member-1")
	_ = app.JoinTeam(background(), "user-1", &val.JoinTeamArgs{InvitationCode: "654321"})
	if !stub.createCalled || !stub.listByTeamCalled || !stub.listMyCalled || !stub.updateCalled || !stub.removeCalled || !stub.joinTeamCalled {
		t.Fatalf("expected all member wrapper calls to forward: %#v", stub)
	}
}

func TestLogMemberAppRejectsInvalidArgs(t *testing.T) {
	app := NewLogMemberApp(&memberAppStub{})
	if _, err := app.Create(background(), "user-1", &val.CreateMemberArgs{}); err == nil {
		t.Fatal("expected create validation error")
	}
	if _, err := app.ListByTeam(background(), "user-1", nil); err == nil {
		t.Fatal("expected list by team validation error")
	}
	if err := app.UpdateRole(background(), "user-1", &val.UpdateMemberRoleArgs{}); err == nil {
		t.Fatal("expected update validation error")
	}
	if err := app.Remove(background(), "user-1", ""); err == nil {
		t.Fatal("expected remove validation error")
	}
	if err := app.JoinTeam(background(), "user-1", &val.JoinTeamArgs{}); err == nil {
		t.Fatal("expected join validation error")
	}
}
