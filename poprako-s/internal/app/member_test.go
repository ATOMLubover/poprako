package app

import (
	"testing"
	"time"

	"poprako-s/internal/app/val"
	"poprako-s/internal/domain/model"
	"poprako-s/internal/domain/service"
	mock_repo "poprako-s/internal/infra/repo/mock"
)

func TestMemberAppListMy(t *testing.T) {
	now := time.Now()
	memberRepo := &mock_repo.MemberRepo{Infos: map[string]model.MemberInfo{
		"member-1": {ID: "member-1", UserID: "user-1", TeamID: "team-1", AssignedAdminAt: &now},
	}}

	app := NewMemberApp(service.NewMemberService(), &mock_repo.UserRepo{}, memberRepo, &mock_repo.InvitationRepo{}, &mock_repo.TxnMgr{}, newMockOSSClient())

	got, err := app.ListMy(background(), "user-1", &val.ListMyMemberArgs{})
	requireNoErr(t, err)

	if len(got) != 1 || got[0].ID != "member-1" {
		t.Fatalf("unexpected result: %#v", got)
	}
}

func TestMemberAppCreateForSuperAdmin(t *testing.T) {
	userRepo := mock_repo.NewMockUserRepo()
	userRepo.Infos["user-1"] = *superAdminUser()
	memberRepo := mock_repo.NewMockMemberRepo()

	app := NewMemberApp(service.NewMemberService(), userRepo, memberRepo, mock_repo.NewMockInvitationRepo(), mock_repo.NewMockTxnMgr(nil), newMockOSSClient())

	res, err := app.Create(background(), "user-1", &val.CreateMemberArgs{UserID: "user-2", TeamID: "team-1", Roles: model.RoleMask(model.RoleAdmin)})
	requireNoErr(t, err)

	stored := memberRepo.Infos[res.ID]
	if stored.UserID != "user-2" || stored.AssignedAdminAt == nil {
		t.Fatalf("unexpected member: %#v", stored)
	}
}

func TestMemberAppListByTeamAssemblesNestedObjects(t *testing.T) {
	now := time.Now()
	memberRepo := mock_repo.NewMockMemberRepo()
	memberRepo.Infos["member-self"] = model.MemberInfo{ID: "member-self", UserID: "user-1", TeamID: "team-1", CreatedAt: now, UpdatedAt: now}
	memberRepo.Infos["member-2"] = model.MemberInfo{
		ID:              "member-2",
		UserID:          "user-2",
		TeamID:          "team-1",
		AssignedAdminAt: &now,
		CreatedAt:       now.Add(time.Second),
		UpdatedAt:       now.Add(time.Second),
		User:            &model.UserInfo{ID: "user-2", Name: "User2", QQ: "10002", AvatarKey: "avatar-user-2", IsAvatarUploaded: true, LastLoginAt: now, CreatedAt: now, UpdatedAt: now},
		Team:            &model.TeamInfo{ID: "team-1", Name: "Team1", AvatarOSSKey: "team-avatar-1", IsAvatarUploaded: true, CreatedAt: now, UpdatedAt: now},
	}
	ossClient := newMockOSSClient()
	ossClient.SetGetURLs(map[string]string{
		"avatar-user-2": "https://cdn.example/user-2",
		"team-avatar-1": "https://cdn.example/team-1",
	})

	app := NewMemberApp(service.NewMemberService(), mock_repo.NewMockUserRepo(), memberRepo, mock_repo.NewMockInvitationRepo(), mock_repo.NewMockTxnMgr(nil), ossClient)

	got, err := app.ListByTeam(background(), "user-1", &val.ListTeamMemberArgs{TeamID: "team-1"})
	requireNoErr(t, err)

	if len(got) != 2 || got[1].User == nil || got[1].User.AvatarURL != "https://cdn.example/user-2" || got[1].Team == nil || got[1].Team.AvatarURL != "https://cdn.example/team-1" {
		t.Fatalf("unexpected members: %#v", got)
	}
}

func TestMemberAppUpdateRolePersistsRoleChange(t *testing.T) {
	memberRepo := mock_repo.NewMockMemberRepo()
	memberRepo.Infos["member-admin"] = *adminMember()
	memberRepo.Infos["member-target"] = model.MemberInfo{ID: "member-target", UserID: "user-2", TeamID: "team-1"}

	app := NewMemberApp(service.NewMemberService(), mock_repo.NewMockUserRepo(), memberRepo, mock_repo.NewMockInvitationRepo(), mock_repo.NewMockTxnMgr(nil), newMockOSSClient())

	err := app.UpdateRole(background(), "user-1", &val.UpdateMemberRoleArgs{ID: "member-target", Roles: model.RoleMask(model.RoleAdmin)})
	requireNoErr(t, err)

	updated := memberRepo.Infos["member-target"]
	if updated.AssignedAdminAt == nil {
		t.Fatalf("expected admin role assigned: %#v", updated)
	}
}

func TestMemberAppJoinTeamUsesMockTxnRepos(t *testing.T) {
	userRepo := mock_repo.NewMockUserRepo()
	userRepo.Infos["user-1"] = *normalUser()
	memberRepo := mock_repo.NewMockMemberRepo()
	invRepo := mock_repo.NewMockInvitationRepo()
	invRepo.Infos["inv-1"] = model.InvitationInfo{ID: "inv-1", TeamID: "team-1", InviteeQQ: "10001", InvitationCode: "654321", Pending: true, ToBeTranslator: true}
	txnMgr := mock_repo.NewMockTxnMgr(newMockTxnContext(mockTxnRepos{member: memberRepo, invitation: invRepo}))

	app := NewMemberApp(service.NewMemberService(), userRepo, memberRepo, invRepo, txnMgr, newMockOSSClient())

	err := app.JoinTeam(background(), "user-1", &val.JoinTeamArgs{InvitationCode: "654321"})
	requireNoErr(t, err)

	if len(memberRepo.Infos) != 1 {
		t.Fatalf("expected created membership, got %#v", memberRepo.Infos)
	}
	if invRepo.Infos["inv-1"].Pending {
		t.Fatalf("expected invitation invalidated: %#v", invRepo.Infos["inv-1"])
	}
}

func TestMemberAppJoinTeamRejectsBadCode(t *testing.T) {
	userRepo := mock_repo.NewMockUserRepo()
	userRepo.Infos["user-1"] = *normalUser()
	invRepo := mock_repo.NewMockInvitationRepo()
	invRepo.Infos["inv-1"] = model.InvitationInfo{ID: "inv-1", TeamID: "team-1", InviteeQQ: "10001", InvitationCode: "654321", Pending: true}

	app := NewMemberApp(service.NewMemberService(), userRepo, mock_repo.NewMockMemberRepo(), invRepo, mock_repo.NewMockTxnMgr(nil), newMockOSSClient())
	if err := app.JoinTeam(background(), "user-1", &val.JoinTeamArgs{InvitationCode: "000000"}); err == nil {
		t.Fatal("expected invalid invitation code error")
	}
}

func TestMemberAppRemoveSuccess(t *testing.T) {
	memberRepo := mock_repo.NewMockMemberRepo()
	memberRepo.Infos["member-admin"] = *adminMember()
	memberRepo.Infos["member-target"] = model.MemberInfo{ID: "member-target", UserID: "user-2", TeamID: "team-1"}

	app := NewMemberApp(service.NewMemberService(), mock_repo.NewMockUserRepo(), memberRepo, mock_repo.NewMockInvitationRepo(), mock_repo.NewMockTxnMgr(nil), newMockOSSClient())
	err := app.Remove(background(), "user-1", "member-target")
	requireNoErr(t, err)

	if _, ok := memberRepo.Infos["member-target"]; ok {
		t.Fatalf("expected member deleted, got %#v", memberRepo.Infos)
	}
}

func TestMemberAppErrorPaths(t *testing.T) {
	t.Run("create rejects duplicate member", func(t *testing.T) {
		userRepo := mock_repo.NewMockUserRepo()
		userRepo.Infos["user-1"] = *superAdminUser()
		memberRepo := mock_repo.NewMockMemberRepo()
		memberRepo.Infos["member-1"] = model.MemberInfo{ID: "member-1", UserID: "user-2", TeamID: "team-1"}
		app := NewMemberApp(service.NewMemberService(), userRepo, memberRepo, mock_repo.NewMockInvitationRepo(), mock_repo.NewMockTxnMgr(nil), newMockOSSClient())
		if _, err := app.Create(background(), "user-1", &val.CreateMemberArgs{UserID: "user-2", TeamID: "team-1", Roles: model.RoleMask(model.RoleAdmin)}); err == nil {
			t.Fatal("expected duplicate member error")
		}
	})

	t.Run("list by team rejects non member", func(t *testing.T) {
		app := NewMemberApp(service.NewMemberService(), mock_repo.NewMockUserRepo(), mock_repo.NewMockMemberRepo(), mock_repo.NewMockInvitationRepo(), mock_repo.NewMockTxnMgr(nil), newMockOSSClient())
		if _, err := app.ListByTeam(background(), "user-1", &val.ListTeamMemberArgs{TeamID: "team-1"}); err == nil {
			t.Fatal("expected list forbidden error")
		}
	})

	t.Run("update rejects missing target", func(t *testing.T) {
		app := NewMemberApp(service.NewMemberService(), mock_repo.NewMockUserRepo(), mock_repo.NewMockMemberRepo(), mock_repo.NewMockInvitationRepo(), mock_repo.NewMockTxnMgr(nil), newMockOSSClient())
		if err := app.UpdateRole(background(), "user-1", &val.UpdateMemberRoleArgs{ID: "missing", Roles: model.RoleMask(model.RoleAdmin)}); err == nil {
			t.Fatal("expected missing member error")
		}
	})

	t.Run("remove rejects non admin", func(t *testing.T) {
		memberRepo := mock_repo.NewMockMemberRepo()
		memberRepo.Infos["member-target"] = model.MemberInfo{ID: "member-target", UserID: "user-2", TeamID: "team-1"}
		memberRepo.Infos["member-1"] = model.MemberInfo{ID: "member-1", UserID: "user-1", TeamID: "team-1"}
		app := NewMemberApp(service.NewMemberService(), mock_repo.NewMockUserRepo(), memberRepo, mock_repo.NewMockInvitationRepo(), mock_repo.NewMockTxnMgr(nil), newMockOSSClient())
		if err := app.Remove(background(), "user-1", "member-target"); err == nil {
			t.Fatal("expected remove forbidden error")
		}
	})

	t.Run("join team rejects existing member", func(t *testing.T) {
		userRepo := mock_repo.NewMockUserRepo()
		userRepo.Infos["user-1"] = *normalUser()
		memberRepo := mock_repo.NewMockMemberRepo()
		memberRepo.Infos["member-1"] = model.MemberInfo{ID: "member-1", UserID: "user-1", TeamID: "team-1"}
		invRepo := mock_repo.NewMockInvitationRepo()
		invRepo.Infos["inv-1"] = model.InvitationInfo{ID: "inv-1", TeamID: "team-1", InviteeQQ: "10001", InvitationCode: "654321", Pending: true}
		app := NewMemberApp(service.NewMemberService(), userRepo, memberRepo, invRepo, mock_repo.NewMockTxnMgr(nil), newMockOSSClient())
		if err := app.JoinTeam(background(), "user-1", &val.JoinTeamArgs{InvitationCode: "654321"}); err == nil {
			t.Fatal("expected existing member error")
		}
	})
}
