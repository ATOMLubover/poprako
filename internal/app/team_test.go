package app

import (
	"errors"
	"testing"
	"time"

	"poprako-s/internal/app/val"
	"poprako-s/internal/domain/model"
	"poprako-s/internal/domain/service"
	mock_repo "poprako-s/internal/infra/repo/mock"
)

func TestTeamAppReserveAvatar(t *testing.T) {
	member := adminMember()
	memberRepo := &mock_repo.MemberRepo{Infos: map[string]model.MemberInfo{
		member.ID: *member,
	}}
	teamRepo := &mock_repo.TeamRepo{Infos: map[string]model.TeamInfo{
		"team-1": {ID: "team-1"},
	}}
	ossClient := newMockOSSClient()
	ossClient.SetGeneratePutPresignedURLFunc(func(objectKey string) (string, error) {
		return "https://upload.example/" + objectKey, nil
	})

	app := NewTeamApp(service.NewTeamService(), service.NewMemberService(), &mock_repo.UserRepo{}, teamRepo, memberRepo, ossClient)

	got, err := app.ReserveAvatar(background(), "user-1", &val.ReserveTeamAvatarArgs{TeamID: "team-1", FileName: "avatar.png"})
	requireNoErr(t, err)

	if got.PutURL != "https://upload.example/team-avatar_team-1.png" {
		t.Fatalf("unexpected put url: %#v", got)
	}
	team := teamRepo.Infos["team-1"]
	if team.AvatarOSSKey != "team-avatar_team-1.png" {
		t.Fatalf("unexpected avatar key: %#v", team)
	}
}

func TestTeamAppAdminFlows(t *testing.T) {
	now := time.Now()
	userRepo := mock_repo.NewMockUserRepo()
	userRepo.Infos["user-1"] = *superAdminUser()
	memberRepo := mock_repo.NewMockMemberRepo()
	memberRepo.Infos["member-1"] = *adminMember()
	teamRepo := mock_repo.NewMockTeamRepo()
	teamRepo.Infos["team-1"] = model.TeamInfo{ID: "team-1", Name: "Existing", AvatarOSSKey: "team-avatar_team-1", IsAvatarUploaded: true, CreatedAt: now, UpdatedAt: now}
	ossClient := newMockOSSClient()
	ossClient.SetGetURLs(map[string]string{"team-avatar_team-1": "https://cdn.example/team-1", "team-avatar_team-1.png": "https://cdn.example/team-1"})
	ossClient.SetPutURLs(map[string]string{"team-avatar_team-1.png": "https://upload.example/team-1"})

	app := NewTeamApp(service.NewTeamService(), service.NewMemberService(), userRepo, teamRepo, memberRepo, ossClient)

	createRes, err := app.Create(background(), "user-1", &val.CreateTeamArgs{Name: "Created", Description: "desc"})
	requireNoErr(t, err)
	if _, ok := teamRepo.Infos[createRes.ID]; !ok {
		t.Fatalf("expected created team, got %#v", teamRepo.Infos)
	}

	listRes, err := app.List(background(), "user-1", &val.ListTeamArgs{})
	requireNoErr(t, err)
	foundAvatar := false
	for _, team := range listRes {
		if team.ID == "team-1" && team.AvatarURL == "https://cdn.example/team-1" {
			foundAvatar = true
			break
		}
	}
	if !foundAvatar {
		t.Fatalf("expected avatar url in team list: %#v", listRes)
	}

	listMyRes, err := app.ListMy(background(), "user-1", &val.ListMyTeamArgs{})
	requireNoErr(t, err)
	if len(listMyRes) == 0 {
		t.Fatalf("expected teams in my list: %#v", listMyRes)
	}

	err = app.Update(background(), "user-1", &val.UpdateTeamArgs{ID: "team-1", Name: "Updated", Description: "new desc"})
	requireNoErr(t, err)
	if teamRepo.Infos["team-1"].Name != "Updated" {
		t.Fatalf("unexpected team update: %#v", teamRepo.Infos["team-1"])
	}

	reserveRes, err := app.ReserveAvatar(background(), "user-1", &val.ReserveTeamAvatarArgs{TeamID: "team-1", FileName: "avatar.png"})
	requireNoErr(t, err)
	if reserveRes.PutURL != "https://upload.example/team-1" {
		t.Fatalf("unexpected reserve avatar result: %#v", reserveRes)
	}

	err = app.ConfirmAvatarUploaded(background(), "user-1", "team-1")
	requireNoErr(t, err)
	if !teamRepo.Infos["team-1"].IsAvatarUploaded {
		t.Fatalf("expected uploaded avatar flag: %#v", teamRepo.Infos["team-1"])
	}

	err = app.Remove(background(), "user-1", "team-1")
	requireNoErr(t, err)
	if _, ok := teamRepo.Infos["team-1"]; ok {
		t.Fatalf("expected deleted team, got %#v", teamRepo.Infos)
	}
	deleted := ossClient.Deleted()
	if len(deleted) == 0 || deleted[0] == "" {
		t.Fatalf("expected avatar oss cleanup, got %#v", deleted)
	}
}

func TestTeamAppRemoveFailsWhenAvatarCleanupFails(t *testing.T) {
	now := time.Now()
	userRepo := mock_repo.NewMockUserRepo()
	userRepo.Infos["user-1"] = model.UserInfo{ID: "user-1", QQ: "100001", Name: "Super", IsSuperAdmin: true, LastLoginAt: now, CreatedAt: now, UpdatedAt: now}
	teamRepo := mock_repo.NewMockTeamRepo()
	teamRepo.Infos["team-1"] = model.TeamInfo{ID: "team-1", Name: "Team", AvatarOSSKey: "team-avatar-1", IsAvatarUploaded: true, CreatedAt: now, UpdatedAt: now}
	memberRepo := mock_repo.NewMockMemberRepo()
	ossClient := newMockOSSClient()
	ossClient.SetDeleteErr(errors.New("boom"))

	app := NewTeamApp(service.NewTeamService(), service.NewMemberService(), userRepo, teamRepo, memberRepo, ossClient)

	err := app.Remove(background(), "user-1", "team-1")
	if err == nil {
		t.Fatal("expected remove failure when avatar cleanup fails")
	}

	if _, ok := teamRepo.Infos["team-1"]; !ok {
		t.Fatalf("expected team to remain, got %#v", teamRepo.Infos)
	}

	if len(ossClient.Deleted()) != 3 {
		t.Fatalf("expected 3 delete attempts, got %#v", ossClient.Deleted())
	}
}

func TestTeamAppPermissionErrors(t *testing.T) {
	now := time.Now()
	userRepo := mock_repo.NewMockUserRepo()
	userRepo.Infos["user-1"] = model.UserInfo{ID: "user-1", QQ: "100001", Name: "Normal", LastLoginAt: now, CreatedAt: now, UpdatedAt: now}
	teamRepo := mock_repo.NewMockTeamRepo()
	teamRepo.Infos["team-1"] = model.TeamInfo{ID: "team-1", Name: "Team"}

	app := NewTeamApp(service.NewTeamService(), service.NewMemberService(), userRepo, teamRepo, mock_repo.NewMockMemberRepo(), newMockOSSClient())
	if _, err := app.List(background(), "user-1", &val.ListTeamArgs{}); err == nil {
		t.Fatal("expected list forbidden error")
	}
	if err := app.Remove(background(), "user-1", "team-1"); err == nil {
		t.Fatal("expected remove forbidden error")
	}
}

func TestTeamAppAdditionalErrorPaths(t *testing.T) {
	t.Run("create rejects non super admin", func(t *testing.T) {
		now := time.Now()
		userRepo := mock_repo.NewMockUserRepo()
		userRepo.Infos["user-1"] = model.UserInfo{ID: "user-1", QQ: "100001", Name: "Normal", LastLoginAt: now, CreatedAt: now, UpdatedAt: now}
		app := NewTeamApp(service.NewTeamService(), service.NewMemberService(), userRepo, mock_repo.NewMockTeamRepo(), mock_repo.NewMockMemberRepo(), newMockOSSClient())
		if _, err := app.Create(background(), "user-1", &val.CreateTeamArgs{Name: "Team"}); err == nil {
			t.Fatal("expected create forbidden error")
		}
	})

	t.Run("list my returns nil for no membership", func(t *testing.T) {
		app := NewTeamApp(service.NewTeamService(), service.NewMemberService(), mock_repo.NewMockUserRepo(), mock_repo.NewMockTeamRepo(), mock_repo.NewMockMemberRepo(), newMockOSSClient())
		got, err := app.ListMy(background(), "user-1", &val.ListMyTeamArgs{})
		requireNoErr(t, err)
		if got != nil {
			t.Fatalf("expected nil team list, got %#v", got)
		}
	})

	t.Run("reserve avatar wraps oss failures", func(t *testing.T) {
		memberRepo := mock_repo.NewMockMemberRepo()
		memberRepo.Infos["member-1"] = *adminMember()
		teamRepo := mock_repo.NewMockTeamRepo()
		teamRepo.Infos["team-1"] = model.TeamInfo{ID: "team-1"}
		ossClient := newMockOSSClient()
		ossClient.SetPutErr(errors.New("boom"))
		app := NewTeamApp(service.NewTeamService(), service.NewMemberService(), mock_repo.NewMockUserRepo(), teamRepo, memberRepo, ossClient)
		if _, err := app.ReserveAvatar(background(), "user-1", &val.ReserveTeamAvatarArgs{TeamID: "team-1", FileName: "avatar.png"}); err == nil {
			t.Fatal("expected reserve avatar failure")
		}
	})

	t.Run("confirm avatar rejects non admin", func(t *testing.T) {
		memberRepo := mock_repo.NewMockMemberRepo()
		memberRepo.Infos["member-1"] = model.MemberInfo{ID: "member-1", UserID: "user-1", TeamID: "team-1"}
		teamRepo := mock_repo.NewMockTeamRepo()
		teamRepo.Infos["team-1"] = model.TeamInfo{ID: "team-1"}
		app := NewTeamApp(service.NewTeamService(), service.NewMemberService(), mock_repo.NewMockUserRepo(), teamRepo, memberRepo, newMockOSSClient())
		if err := app.ConfirmAvatarUploaded(background(), "user-1", "team-1"); err == nil {
			t.Fatal("expected confirm forbidden error")
		}
	})
}
