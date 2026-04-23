package app

import (
	"errors"
	"testing"
	"time"

	"poprako-s/internal/app/val"
	"poprako-s/internal/cfg"
	"poprako-s/internal/domain/model"
	"poprako-s/internal/domain/service"
	mock_repo "poprako-s/internal/infra/repo/mock"
)

func TestUserAppGetMyStats(t *testing.T) {
	userRepo := &mock_repo.UserRepo{Stats: map[string]model.UserStats{
		"user-1": {UserID: "user-1", TotalAssignmentCount: 3, ActiveAssignmentCount: 2, FinishedAssignmentCount: 1},
	}}

	app := NewUserApp(service.NewUserService(), service.NewMemberService(), userRepo, &mock_repo.InvitationRepo{}, &mock_repo.MemberRepo{}, &mock_repo.TxnMgr{}, mock_repo.NewMockOSSMessageRepo(), newMockEventBus(), newMockOSSClient(), &cfg.AuthCfg{})

	got, err := app.GetMyStats(background(), "user-1")
	requireNoErr(t, err)

	if got.UserID != "user-1" || got.TotalAssignmentCount != 3 || got.ActiveAssignmentCount != 2 || got.FinishedAssignmentCount != 1 {
		t.Fatalf("unexpected stats: %#v", got)
	}
}

func TestUserAppLoginSuccess(t *testing.T) {
	userSvc := service.NewUserService()
	hash, err := userSvc.HashPwd("secret123")
	requireNoErr(t, err)

	now := time.Now()
	userRepo := mock_repo.NewMockUserRepo()
	userRepo.Creds["100001"] = model.UserCreds{QQ: "100001", PwdHash: hash}
	userRepo.Infos["user-1"] = model.UserInfo{ID: "user-1", QQ: "100001", Name: "Tester", LastLoginAt: now, CreatedAt: now, UpdatedAt: now}
	eventBus := newMockEventBus()

	app := NewUserApp(userSvc, service.NewMemberService(), userRepo, mock_repo.NewMockInvitationRepo(), mock_repo.NewMockMemberRepo(), mock_repo.NewMockTxnMgr(nil), mock_repo.NewMockOSSMessageRepo(), eventBus, newMockOSSClient(), &cfg.AuthCfg{SecretKey: "secret-key", ExpHrs: 1})

	res, err := app.Login(background(), &val.LoginUserArgs{QQ: "100001", Pwd: "secret123"})
	requireNoErr(t, err)

	if res.UserID != "user-1" || res.AccessToken == "" {
		t.Fatalf("unexpected login result: %#v", res)
	}
	pubCalls := eventBus.PubCalls()
	if len(pubCalls) != 1 {
		t.Fatalf("expected login event publication, got %#v", pubCalls)
	}
}

func TestUserAppRegUsesMockTxnRepos(t *testing.T) {
	userRepo := mock_repo.NewMockUserRepo()
	invRepo := mock_repo.NewMockInvitationRepo()
	invRepo.Infos["inv-1"] = model.MemberInvitationInfo{ID: "inv-1", InvitorID: "admin-1", TeamID: "team-1", InviteeQQ: "100001", Pending: true}
	memberRepo := mock_repo.NewMockMemberRepo()
	eventBus := newMockEventBus()
	txnMgr := mock_repo.NewMockTxnMgr(newMockTxnContext(mockTxnRepos{user: userRepo, invitation: invRepo, member: memberRepo}))

	app := NewUserApp(service.NewUserService(), service.NewMemberService(), userRepo, invRepo, memberRepo, txnMgr, mock_repo.NewMockOSSMessageRepo(), eventBus, newMockOSSClient(), &cfg.AuthCfg{SecretKey: "secret-key", ExpHrs: 1})

	res, err := app.Reg(background(), &val.RegUserArgs{QQ: "100001", Pwd: "secret123", Name: "New User", InvCode: "whatever"})
	requireNoErr(t, err)

	if _, ok := userRepo.Infos[res.UserID]; !ok {
		t.Fatalf("expected created user, got %#v", userRepo.Infos)
	}
	if len(memberRepo.Infos) != 1 {
		t.Fatalf("expected member record created, got %#v", memberRepo.Infos)
	}
	pubCalls := eventBus.PubCalls()
	if len(pubCalls) != 1 {
		t.Fatalf("expected user created event publication, got %#v", pubCalls)
	}
}

func TestUserAppProfileAvatarAndRemoveFlows(t *testing.T) {
	now := time.Now()
	userRepo := mock_repo.NewMockUserRepo()
	userRepo.Infos["user-1"] = model.UserInfo{ID: "user-1", QQ: "100001", Name: "Old Name", AvatarKey: "user-avatar_user-1", IsAvatarUploaded: true, LastLoginAt: now, CreatedAt: now, UpdatedAt: now}
	userRepo.Infos["user-2"] = model.UserInfo{ID: "user-2", QQ: "100002", Name: "To Delete", LastLoginAt: now, CreatedAt: now, UpdatedAt: now}
	ossClient := newMockOSSClient()
	ossClient.SetGetURLs(map[string]string{"user-avatar_user-1": "https://cdn.example/user-1", "user-avatar_user-1.png": "https://cdn.example/user-1"})
	ossClient.SetPutURLs(map[string]string{"user_user-1/avatar.png": "https://upload.example/user-1"})
	msgRepo := mock_repo.NewMockOSSMessageRepo()
	txCx := newMockTxnContext(mockTxnRepos{user: userRepo, ossMessage: msgRepo})
	txnMgr := mock_repo.NewMockTxnMgr(txCx)

	app := NewUserApp(service.NewUserService(), service.NewMemberService(), userRepo, mock_repo.NewMockInvitationRepo(), mock_repo.NewMockMemberRepo(), txnMgr, msgRepo, newMockEventBus(), ossClient, &cfg.AuthCfg{SecretKey: "secret-key", ExpHrs: 1})

	info, err := app.GetInfo(background(), "user-1")
	requireNoErr(t, err)
	if info.AvatarURL != "https://cdn.example/user-1" {
		t.Fatalf("unexpected user info: %#v", info)
	}

	err = app.UpdateMyInfo(background(), "user-1", &val.UpdateUserArgs{Name: "New Name", QQ: "100003"})
	requireNoErr(t, err)
	if userRepo.Infos["user-1"].Name != "New Name" {
		t.Fatalf("unexpected user update: %#v", userRepo.Infos["user-1"])
	}

	reserveRes, err := app.ReserveMyAvatar(background(), "user-1", &val.ReserveUserAvatarArgs{FileName: "avatar.png"})
	requireNoErr(t, err)
	if reserveRes.PutURL != "https://upload.example/user-1" {
		t.Fatalf("unexpected reserve avatar result: %#v", reserveRes)
	}

	err = app.ConfirmMyAvatarUploaded(background(), "user-1")
	requireNoErr(t, err)
	if !userRepo.Infos["user-1"].IsAvatarUploaded {
		t.Fatalf("expected uploaded avatar flag: %#v", userRepo.Infos["user-1"])
	}

	err = app.Remove(background(), "user-1", "user-1")
	if err == nil {
		t.Fatal("expected self remove to fail")
	}

	err = app.Remove(background(), "user-1", "user-2")
	requireNoErr(t, err)
	if _, ok := userRepo.Infos["user-2"]; ok {
		t.Fatalf("expected removed user, got %#v", userRepo.Infos)
	}
}

func TestUserAppRemoveFailsWhenAvatarCleanupFails(t *testing.T) {
	now := time.Now()
	userRepo := mock_repo.NewMockUserRepo()
	userRepo.Infos["user-1"] = model.UserInfo{ID: "user-1", QQ: "100001", Name: "Admin", LastLoginAt: now, CreatedAt: now, UpdatedAt: now}
	userRepo.Infos["user-2"] = model.UserInfo{ID: "user-2", QQ: "100002", Name: "Target", AvatarKey: "avatar-user-2", LastLoginAt: now, CreatedAt: now, UpdatedAt: now}
	msgRepo := mock_repo.NewMockOSSMessageRepo()
	msgRepo.InsertErr = errors.New("boom")

	app := NewUserApp(service.NewUserService(), service.NewMemberService(), userRepo, mock_repo.NewMockInvitationRepo(), mock_repo.NewMockMemberRepo(), mock_repo.NewMockTxnMgr(newMockTxnContext(mockTxnRepos{user: userRepo, ossMessage: msgRepo})), msgRepo, newMockEventBus(), newMockOSSClient(), &cfg.AuthCfg{SecretKey: "secret-key", ExpHrs: 1})

	err := app.Remove(background(), "user-1", "user-2")
	if err == nil {
		t.Fatal("expected remove failure when avatar cleanup fails")
	}

	if len(msgRepo.Messages) != 0 {
		t.Fatalf("expected no queued messages, got %#v", msgRepo.Messages)
	}
}

func TestUserAppLoginRejectsWrongPassword(t *testing.T) {
	userSvc := service.NewUserService()
	hash, err := userSvc.HashPwd("secret123")
	requireNoErr(t, err)
	now := time.Now()
	userRepo := mock_repo.NewMockUserRepo()
	userRepo.Creds["100001"] = model.UserCreds{QQ: "100001", PwdHash: hash}
	userRepo.Infos["user-1"] = model.UserInfo{ID: "user-1", QQ: "100001", Name: "Tester", LastLoginAt: now, CreatedAt: now, UpdatedAt: now}

	app := NewUserApp(userSvc, service.NewMemberService(), userRepo, mock_repo.NewMockInvitationRepo(), mock_repo.NewMockMemberRepo(), mock_repo.NewMockTxnMgr(nil), mock_repo.NewMockOSSMessageRepo(), newMockEventBus(), newMockOSSClient(), &cfg.AuthCfg{SecretKey: "secret-key", ExpHrs: 1})
	if _, err := app.Login(background(), &val.LoginUserArgs{QQ: "100001", Pwd: "wrongpwd"}); err == nil {
		t.Fatal("expected wrong password error")
	}
}

func TestUserAppGetMyInfo(t *testing.T) {
	now := time.Now()
	userRepo := mock_repo.NewMockUserRepo()
	userRepo.Infos["user-1"] = model.UserInfo{ID: "user-1", QQ: "100001", Name: "Tester", AvatarKey: "avatar-1", IsAvatarUploaded: true, LastLoginAt: now, CreatedAt: now, UpdatedAt: now}
	ossClient := newMockOSSClient()
	ossClient.SetGetURLs(map[string]string{"avatar-1": "https://cdn.example/avatar-1"})

	app := NewUserApp(service.NewUserService(), service.NewMemberService(), userRepo, mock_repo.NewMockInvitationRepo(), mock_repo.NewMockMemberRepo(), mock_repo.NewMockTxnMgr(nil), mock_repo.NewMockOSSMessageRepo(), newMockEventBus(), ossClient, &cfg.AuthCfg{SecretKey: "secret", ExpHrs: 1})

	info, err := app.GetMyInfo(background(), "user-1")
	requireNoErr(t, err)
	if info.ID != "user-1" || info.AvatarURL != "https://cdn.example/avatar-1" {
		t.Fatalf("unexpected user info: %#v", info)
	}
}

func TestUserValidationHelpers(t *testing.T) {
	app := &logUserAppImpl{}

	if err := app.validateLoginArgs(nil); err == nil {
		t.Fatal("expected nil login args error")
	}
	if err := app.validateLoginArgs(&val.LoginUserArgs{QQ: "1", Pwd: "123456"}); err == nil {
		t.Fatal("expected qq validation error")
	}
	if err := app.validateLoginArgs(&val.LoginUserArgs{QQ: "100001", Pwd: "123"}); err == nil {
		t.Fatal("expected password validation error")
	}
	if err := app.validateLoginArgs(&val.LoginUserArgs{QQ: "100001", Pwd: "123456"}); err != nil {
		t.Fatalf("unexpected login validation error: %v", err)
	}

	if err := app.validateRegArgs(nil); err == nil {
		t.Fatal("expected nil reg args error")
	}
	if err := app.validateRegArgs(&val.RegUserArgs{QQ: "100001", Pwd: "123456", Name: "", InvCode: "x"}); err == nil {
		t.Fatal("expected missing name error")
	}
	if err := app.validateRegArgs(&val.RegUserArgs{QQ: "100001", Pwd: "123456", Name: "Name", InvCode: ""}); err == nil {
		t.Fatal("expected missing invcode error")
	}
	if err := app.validateRegArgs(&val.RegUserArgs{QQ: "100001", Pwd: "123456", Name: "Name", InvCode: "inv"}); err != nil {
		t.Fatalf("unexpected reg validation error: %v", err)
	}

	if err := app.validateUpdateUserArgs(nil); err == nil {
		t.Fatal("expected nil update args error")
	}
	if err := app.validateUpdateUserArgs(&val.UpdateUserArgs{Name: "", QQ: "100001"}); err == nil {
		t.Fatal("expected missing name error")
	}
	if err := app.validateUpdateUserArgs(&val.UpdateUserArgs{Name: "Name", QQ: "1"}); err == nil {
		t.Fatal("expected qq validation error")
	}
	if err := app.validateUpdateUserArgs(&val.UpdateUserArgs{Name: "Name", QQ: "100001"}); err != nil {
		t.Fatalf("unexpected update validation error: %v", err)
	}
}

func TestUserAppErrorPaths(t *testing.T) {
	t.Run("login wraps event bus failures", func(t *testing.T) {
		userSvc := service.NewUserService()
		hash, err := userSvc.HashPwd("secret123")
		requireNoErr(t, err)
		now := time.Now()
		userRepo := mock_repo.NewMockUserRepo()
		userRepo.Creds["100001"] = model.UserCreds{QQ: "100001", PwdHash: hash}
		userRepo.Infos["user-1"] = model.UserInfo{ID: "user-1", QQ: "100001", Name: "Tester", LastLoginAt: now, CreatedAt: now, UpdatedAt: now}
		eventBus := newMockEventBus()
		eventBus.SetPubErr(errors.New("boom"))
		app := NewUserApp(userSvc, service.NewMemberService(), userRepo, mock_repo.NewMockInvitationRepo(), mock_repo.NewMockMemberRepo(), mock_repo.NewMockTxnMgr(nil), mock_repo.NewMockOSSMessageRepo(), eventBus, newMockOSSClient(), &cfg.AuthCfg{SecretKey: "secret", ExpHrs: 1})
		if _, err := app.Login(background(), &val.LoginUserArgs{QQ: "100001", Pwd: "secret123"}); err == nil {
			t.Fatal("expected login event bus error")
		}
	})

	t.Run("reg rejects invalid invitation", func(t *testing.T) {
		userRepo := mock_repo.NewMockUserRepo()
		invRepo := mock_repo.NewMockInvitationRepo()
		memberRepo := mock_repo.NewMockMemberRepo()
		txnMgr := mock_repo.NewMockTxnMgr(newMockTxnContext(mockTxnRepos{user: userRepo, invitation: invRepo, member: memberRepo}))
		app := NewUserApp(service.NewUserService(), service.NewMemberService(), userRepo, invRepo, memberRepo, txnMgr, mock_repo.NewMockOSSMessageRepo(), newMockEventBus(), newMockOSSClient(), &cfg.AuthCfg{SecretKey: "secret", ExpHrs: 1})
		if _, err := app.Reg(background(), &val.RegUserArgs{QQ: "100001", Pwd: "secret123", Name: "Name", InvCode: "bad"}); err == nil {
			t.Fatal("expected invalid invitation error")
		}
	})

	t.Run("get info rejects missing user", func(t *testing.T) {
		app := NewUserApp(service.NewUserService(), service.NewMemberService(), mock_repo.NewMockUserRepo(), mock_repo.NewMockInvitationRepo(), mock_repo.NewMockMemberRepo(), mock_repo.NewMockTxnMgr(nil), mock_repo.NewMockOSSMessageRepo(), newMockEventBus(), newMockOSSClient(), &cfg.AuthCfg{SecretKey: "secret", ExpHrs: 1})
		if _, err := app.GetInfo(background(), "missing"); err == nil {
			t.Fatal("expected missing user error")
		}
	})

	t.Run("reserve avatar wraps put failures", func(t *testing.T) {
		userRepo := mock_repo.NewMockUserRepo()
		userRepo.Infos["user-1"] = *normalUser()
		ossClient := newMockOSSClient()
		ossClient.SetPutErr(errors.New("boom"))
		app := NewUserApp(service.NewUserService(), service.NewMemberService(), userRepo, mock_repo.NewMockInvitationRepo(), mock_repo.NewMockMemberRepo(), mock_repo.NewMockTxnMgr(nil), mock_repo.NewMockOSSMessageRepo(), newMockEventBus(), ossClient, &cfg.AuthCfg{SecretKey: "secret", ExpHrs: 1})
		if _, err := app.ReserveMyAvatar(background(), "user-1", &val.ReserveUserAvatarArgs{FileName: "avatar.png"}); err == nil {
			t.Fatal("expected reserve avatar failure")
		}
	})

	t.Run("confirm avatar rejects missing user", func(t *testing.T) {
		app := NewUserApp(service.NewUserService(), service.NewMemberService(), mock_repo.NewMockUserRepo(), mock_repo.NewMockInvitationRepo(), mock_repo.NewMockMemberRepo(), mock_repo.NewMockTxnMgr(nil), mock_repo.NewMockOSSMessageRepo(), newMockEventBus(), newMockOSSClient(), &cfg.AuthCfg{SecretKey: "secret", ExpHrs: 1})
		if err := app.ConfirmMyAvatarUploaded(background(), "user-1"); err == nil {
			t.Fatal("expected confirm avatar failure")
		}
	})
}
