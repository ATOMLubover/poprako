package app

import (
	"poprako-s/internal/app/val"
	"testing"
)

func TestNewLogUserAppRejectsNilDependencies(t *testing.T) {
	assertPanics(t, func() { NewLogUserApp(nil) })
}

func TestLogUserAppForwardsAllMethods(t *testing.T) {
	stub := &userAppStub{}
	app := NewLogUserApp(stub)
	_, _ = app.ParseToken(background(), "token")
	_, _ = app.Login(background(), &val.LoginUserArgs{QQ: "100001", Pwd: "secret123"})
	_, _ = app.Reg(background(), &val.RegUserArgs{QQ: "100001", Pwd: "secret123", Name: "Name", InvCode: "654321"})
	_, _ = app.GetInfo(background(), "user-1")
	_, _ = app.GetMyInfo(background(), "user-1")
	_ = app.UpdateMyInfo(background(), "user-1", &val.UpdateUserArgs{Name: "Name", QQ: "100001"})
	_, _ = app.GetMyStats(background(), "user-1")
	_, _ = app.ReserveMyAvatar(background(), "user-1", &val.ReserveUserAvatarArgs{FileName: "avatar.png"})
	_ = app.ConfirmMyAvatarUploaded(background(), "user-1")
	_ = app.Remove(background(), "user-1", "user-2")
	if !stub.parseTokenCalled || !stub.loginCalled || !stub.regCalled || !stub.getCalled || !stub.getMyCalled || !stub.updateCalled || !stub.statsCalled || !stub.reserveCalled || !stub.confirmCalled || !stub.removeCalled {
		t.Fatalf("expected all user wrapper calls to forward: %#v", stub)
	}
}

func TestLogUserAppRejectsInvalidArgs(t *testing.T) {
	app := NewLogUserApp(&userAppStub{})
	if _, err := app.Login(background(), nil); err == nil {
		t.Fatal("expected login validation error")
	}
	if _, err := app.Reg(background(), &val.RegUserArgs{QQ: "100001", Pwd: "secret123", Name: "Name"}); err == nil {
		t.Fatal("expected reg validation error")
	}
	if _, err := app.GetInfo(background(), ""); err == nil {
		t.Fatal("expected get info validation error")
	}
	if _, err := app.GetMyInfo(background(), ""); err == nil {
		t.Fatal("expected get my info validation error")
	}
	if err := app.UpdateMyInfo(background(), "user-1", &val.UpdateUserArgs{Name: "", QQ: "1"}); err == nil {
		t.Fatal("expected update validation error")
	}
	if _, err := app.GetMyStats(background(), ""); err == nil {
		t.Fatal("expected stats validation error")
	}
	if _, err := app.ReserveMyAvatar(background(), "user-1", &val.ReserveUserAvatarArgs{}); err == nil {
		t.Fatal("expected reserve avatar validation error")
	}
	if err := app.ConfirmMyAvatarUploaded(background(), ""); err == nil {
		t.Fatal("expected confirm avatar validation error")
	}
	if err := app.Remove(background(), "user-1", ""); err == nil {
		t.Fatal("expected remove validation error")
	}
}
