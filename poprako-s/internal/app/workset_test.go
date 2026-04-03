package app

import (
	"testing"

	"poprako-s/internal/app/val"
	"poprako-s/internal/domain/model"
	"poprako-s/internal/domain/service"
	mock_repo "poprako-s/internal/infra/repo/mock"
)

func TestWorksetAppList(t *testing.T) {
	memberRepo := &mock_repo.MemberRepo{Infos: map[string]model.MemberInfo{
		"member-1": {ID: "member-1", UserID: "user-1", TeamID: "team-1"},
	}}
	worksetRepo := &mock_repo.WorksetRepo{Infos: map[string]model.WorksetInfo{
		"workset-1": {ID: "workset-1", TeamID: "team-1", Name: "ws"},
	}}

	app := NewWorksetApp(service.NewWorksetService(), memberRepo, worksetRepo, &mock_repo.TxnMgr{})

	got, err := app.List(background(), "user-1", &val.ListWorksetArgs{TeamID: "team-1"})
	requireNoErr(t, err)

	if len(got) != 1 || got[0].ID != "workset-1" {
		t.Fatalf("unexpected result: %#v", got)
	}
}

func TestWorksetAppCreateUsesMockTxnRepos(t *testing.T) {
	memberRepo := mock_repo.NewMockMemberRepo()
	memberRepo.Infos["member-1"] = *adminMember()
	worksetRepo := mock_repo.NewMockWorksetRepo()
	description := "desc"
	txnMgr := mock_repo.NewMockTxnMgr(newMockTxnContext(mockTxnRepos{member: memberRepo, workset: worksetRepo}))

	app := NewWorksetApp(service.NewWorksetService(), memberRepo, worksetRepo, txnMgr)

	res, err := app.Create(background(), "user-1", &val.CreateWorksetArgs{TeamID: "team-1", Name: "WS", Description: &description})
	requireNoErr(t, err)

	stored := worksetRepo.Infos[res.ID]
	if stored.Index != 0 || stored.Description != "desc" {
		t.Fatalf("unexpected workset: %#v", stored)
	}
}

func TestWorksetAppUpdateAndRemove(t *testing.T) {
	memberRepo := mock_repo.NewMockMemberRepo()
	memberRepo.Infos["member-1"] = *adminMember()
	worksetRepo := mock_repo.NewMockWorksetRepo()
	worksetRepo.Infos["workset-1"] = model.WorksetInfo{ID: "workset-1", TeamID: "team-1", Name: "Old"}
	app := NewWorksetApp(service.NewWorksetService(), memberRepo, worksetRepo, mock_repo.NewMockTxnMgr(nil))

	desc := "new desc"
	err := app.Update(background(), "user-1", &val.UpdateWorksetArgs{ID: "workset-1", Name: "New", Description: &desc})
	requireNoErr(t, err)
	if worksetRepo.Infos["workset-1"].Name != "New" {
		t.Fatalf("unexpected updated workset: %#v", worksetRepo.Infos["workset-1"])
	}

	err = app.Remove(background(), "user-1", "workset-1")
	requireNoErr(t, err)
	if _, ok := worksetRepo.Infos["workset-1"]; ok {
		t.Fatalf("expected deleted workset, got %#v", worksetRepo.Infos)
	}
}

func TestWorksetAppListForbidden(t *testing.T) {
	worksetRepo := mock_repo.NewMockWorksetRepo()
	worksetRepo.Infos["workset-1"] = model.WorksetInfo{ID: "workset-1", TeamID: "team-1"}

	app := NewWorksetApp(service.NewWorksetService(), mock_repo.NewMockMemberRepo(), worksetRepo, mock_repo.NewMockTxnMgr(nil))
	if _, err := app.List(background(), "user-1", &val.ListWorksetArgs{TeamID: "team-1"}); err == nil {
		t.Fatal("expected forbidden error")
	}
}

func TestWorksetAppAdditionalErrorPaths(t *testing.T) {
	t.Run("create rejects non admin", func(t *testing.T) {
		memberRepo := mock_repo.NewMockMemberRepo()
		memberRepo.Infos["member-1"] = model.MemberInfo{ID: "member-1", UserID: "user-1", TeamID: "team-1"}
		worksetRepo := mock_repo.NewMockWorksetRepo()
		txnMgr := mock_repo.NewMockTxnMgr(newMockTxnContext(mockTxnRepos{member: memberRepo, workset: worksetRepo}))
		app := NewWorksetApp(service.NewWorksetService(), memberRepo, worksetRepo, txnMgr)
		if _, err := app.Create(background(), "user-1", &val.CreateWorksetArgs{TeamID: "team-1", Name: "WS"}); err == nil {
			t.Fatal("expected forbidden create error")
		}
	})

	t.Run("update rejects missing workset", func(t *testing.T) {
		app := NewWorksetApp(service.NewWorksetService(), mock_repo.NewMockMemberRepo(), mock_repo.NewMockWorksetRepo(), mock_repo.NewMockTxnMgr(nil))
		if err := app.Update(background(), "user-1", &val.UpdateWorksetArgs{ID: "missing", Name: "WS"}); err == nil {
			t.Fatal("expected missing workset error")
		}
	})

	t.Run("remove rejects non admin", func(t *testing.T) {
		memberRepo := mock_repo.NewMockMemberRepo()
		memberRepo.Infos["member-1"] = model.MemberInfo{ID: "member-1", UserID: "user-1", TeamID: "team-1"}
		worksetRepo := mock_repo.NewMockWorksetRepo()
		worksetRepo.Infos["workset-1"] = model.WorksetInfo{ID: "workset-1", TeamID: "team-1"}
		app := NewWorksetApp(service.NewWorksetService(), memberRepo, worksetRepo, mock_repo.NewMockTxnMgr(nil))
		if err := app.Remove(background(), "user-1", "workset-1"); err == nil {
			t.Fatal("expected forbidden remove error")
		}
	})
}
