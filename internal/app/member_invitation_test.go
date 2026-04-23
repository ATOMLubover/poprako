package app

import (
	"testing"
	"time"

	"poprako-s/internal/app/val"
	"poprako-s/internal/domain/model"
	"poprako-s/internal/domain/service"
	mock_repo "poprako-s/internal/infra/repo/mock"
)

func boolPtr(v bool) *bool {
	return &v
}

func TestInvitationAppList(t *testing.T) {
	memberRepo := &mock_repo.MemberRepo{Infos: map[string]model.MemberInfo{
		"member-1": {ID: "member-1", UserID: "user-1", TeamID: "team-1"},
	}}
	now := time.Now()
	invRepo := &mock_repo.InvitationRepo{Infos: map[string]model.MemberInvitationInfo{
		"inv-1": {ID: "inv-1", TeamID: "team-1", InviteeQQ: "123", Pending: true, CreatedAt: now},
		"inv-2": {ID: "inv-2", TeamID: "team-1", InviteeQQ: "456", Pending: false, CreatedAt: now},
	}}

	app := NewInvitationApp(service.NewMemberInvitationService(), memberRepo, invRepo)

	got, err := app.List(background(), "user-1", &val.ListTeamInvitationArgs{TeamID: "team-1"})
	requireNoErr(t, err)

	if len(got) != 2 {
		t.Fatalf("unexpected result: %#v", got)
	}

	got, err = app.List(background(), "user-1", &val.ListTeamInvitationArgs{TeamID: "team-1", Pending: boolPtr(true)})
	requireNoErr(t, err)

	if len(got) != 1 || got[0].ID != "inv-1" {
		t.Fatalf("unexpected filtered result: %#v", got)
	}
}

func TestInvitationAppCreatePersistsRoles(t *testing.T) {
	memberRepo := mock_repo.NewMockMemberRepo()
	memberRepo.Infos["member-1"] = *adminMember()
	invRepo := mock_repo.NewMockInvitationRepo()

	app := NewInvitationApp(service.NewMemberInvitationService(), memberRepo, invRepo)

	got, err := app.Create(background(), "user-1", &val.CreateInvitationArgs{
		TeamID:    "team-1",
		InviteeQQ: "10002",
		Roles:     model.RoleMask(model.RoleTranslator) | model.RoleMask(model.RoleProofreader),
	})
	requireNoErr(t, err)

	stored := invRepo.Infos[got.ID]
	if !stored.ToBeTranslator || !stored.ToBeProofreader || stored.ToBeAdmin {
		t.Fatalf("unexpected invitation: %#v", stored)
	}
}

func TestInvitationAppUpdatePersistsRoles(t *testing.T) {
	memberRepo := mock_repo.NewMockMemberRepo()
	memberRepo.Infos["member-1"] = *adminMember()
	invRepo := mock_repo.NewMockInvitationRepo()
	invRepo.Infos["inv-1"] = model.MemberInvitationInfo{ID: "inv-1", TeamID: "team-1", Pending: true}

	app := NewInvitationApp(service.NewMemberInvitationService(), memberRepo, invRepo)

	err := app.Update(background(), "user-1", &val.UpdateInvitationArgs{ID: "inv-1", TeamID: "team-1", Roles: model.RoleMask(model.RoleAdmin) | model.RoleMask(model.RolePublisher)})
	requireNoErr(t, err)

	updated := invRepo.Infos["inv-1"]
	if !updated.ToBeAdmin || !updated.ToBePublisher {
		t.Fatalf("unexpected invitation update: %#v", updated)
	}
}

func TestInvitationAppRejectsRedrawerRole(t *testing.T) {
	memberRepo := mock_repo.NewMockMemberRepo()
	memberRepo.Infos["member-1"] = *adminMember()
	invRepo := mock_repo.NewMockInvitationRepo()
	invRepo.Infos["inv-1"] = model.MemberInvitationInfo{ID: "inv-1", TeamID: "team-1", Pending: true}

	app := NewInvitationApp(service.NewMemberInvitationService(), memberRepo, invRepo)

	if _, err := app.Create(background(), "user-1", &val.CreateInvitationArgs{TeamID: "team-1", InviteeQQ: "10002", Roles: model.RoleMask(model.RoleRedrawer)}); err == nil {
		t.Fatal("expected redrawer create rejection")
	}

	if err := app.Update(background(), "user-1", &val.UpdateInvitationArgs{ID: "inv-1", TeamID: "team-1", Roles: model.RoleMask(model.RoleRedrawer)}); err == nil {
		t.Fatal("expected redrawer update rejection")
	}
}

func TestInvitationAppRemoveDeletesInvitation(t *testing.T) {
	memberRepo := mock_repo.NewMockMemberRepo()
	memberRepo.Infos["member-1"] = *adminMember()
	invRepo := mock_repo.NewMockInvitationRepo()
	invRepo.Infos["inv-1"] = model.MemberInvitationInfo{ID: "inv-1", TeamID: "team-1", Pending: true}

	app := NewInvitationApp(service.NewMemberInvitationService(), memberRepo, invRepo)

	err := app.Remove(background(), "user-1", "inv-1")
	requireNoErr(t, err)

	if _, ok := invRepo.Infos["inv-1"]; ok {
		t.Fatalf("expected invitation deleted, got %#v", invRepo.Infos)
	}
}

func TestInvitationAppRemoveMissingInvitation(t *testing.T) {
	app := NewInvitationApp(service.NewMemberInvitationService(), mock_repo.NewMockMemberRepo(), mock_repo.NewMockInvitationRepo())
	if err := app.Remove(background(), "user-1", "missing"); err == nil {
		t.Fatal("expected missing invitation error")
	}
}

func TestInvitationAppPermissionErrors(t *testing.T) {
	app := NewInvitationApp(service.NewMemberInvitationService(), mock_repo.NewMockMemberRepo(), mock_repo.NewMockInvitationRepo())
	if _, err := app.List(background(), "user-1", &val.ListTeamInvitationArgs{TeamID: "team-1"}); err == nil {
		t.Fatal("expected list forbidden error")
	}
	if _, err := app.Create(background(), "user-1", &val.CreateInvitationArgs{TeamID: "team-1", InviteeQQ: "10002", Roles: model.RoleMask(model.RoleTranslator)}); err == nil {
		t.Fatal("expected create forbidden error")
	}
	if err := app.Update(background(), "user-1", &val.UpdateInvitationArgs{ID: "inv-1", TeamID: "team-1", Roles: model.RoleMask(model.RoleAdmin)}); err == nil {
		t.Fatal("expected update forbidden error")
	}
	if err := app.Remove(background(), "user-1", "inv-1"); err == nil {
		t.Fatal("expected remove forbidden error")
	}
}
