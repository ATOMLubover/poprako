package app

import (
	"errors"
	"testing"

	"poprako-s/internal/app/val"
	"poprako-s/internal/domain/event"
	"poprako-s/internal/domain/model"
	"poprako-s/internal/domain/service"
	mock_repo "poprako-s/internal/infra/repo/mock"
)

func TestComicAppList(t *testing.T) {
	worksetRepo := &mock_repo.WorksetRepo{Infos: map[string]model.WorksetInfo{
		"workset-1": {ID: "workset-1", TeamID: "team-1"},
	}}
	memberRepo := &mock_repo.MemberRepo{Infos: map[string]model.MemberInfo{
		"member-1": {ID: "member-1", UserID: "user-1", TeamID: "team-1"},
	}}
	comicRepo := &mock_repo.ComicRepo{Infos: map[string]model.ComicInfo{
		"comic-1": {ID: "comic-1", WorksetID: "workset-1", Title: "title"},
	}}

	app := NewComicApp(service.NewComicService(), memberRepo, worksetRepo, comicRepo, &mock_repo.TxnMgr{}, mock_repo.NewMockOSSMessageRepo(), newMockEventBus(), newMockOSSClient())

	got, err := app.List(background(), "user-1", &val.ListComicArgs{WorksetID: "workset-1"})
	requireNoErr(t, err)

	if len(got) != 1 || got[0].ID != "comic-1" {
		t.Fatalf("unexpected result: %#v", got)
	}
}

func TestComicAppListFiltersByFuzzyTitle(t *testing.T) {
	worksetRepo := mock_repo.NewMockWorksetRepo()
	worksetRepo.Infos["workset-1"] = model.WorksetInfo{ID: "workset-1", TeamID: "team-1"}
	memberRepo := mock_repo.NewMockMemberRepo()
	memberRepo.Infos["member-1"] = model.MemberInfo{ID: "member-1", UserID: "user-1", TeamID: "team-1"}
	comicRepo := mock_repo.NewMockComicRepo()
	comicRepo.Infos["comic-1"] = model.ComicInfo{ID: "comic-1", WorksetID: "workset-1", Index: 0, Title: "Target Story", Author: "A"}
	comicRepo.Infos["comic-2"] = model.ComicInfo{ID: "comic-2", WorksetID: "workset-1", Index: 1, Title: "Another", Author: "B"}

	app := NewComicApp(service.NewComicService(), memberRepo, worksetRepo, comicRepo, mock_repo.NewMockTxnMgr(nil), mock_repo.NewMockOSSMessageRepo(), newMockEventBus(), newMockOSSClient())

	got, err := app.List(background(), "user-1", &val.ListComicArgs{WorksetID: "workset-1", FuzzyTitle: "Target"})
	requireNoErr(t, err)

	if len(got) != 1 || got[0].ID != "comic-1" {
		t.Fatalf("unexpected comics: %#v", got)
	}
}

func TestComicAppCreateUsesMockTxnRepos(t *testing.T) {
	comicRepo := mock_repo.NewMockComicRepo()
	memberRepo := mock_repo.NewMockMemberRepo()
	memberRepo.Infos["member-1"] = *adminMember()
	worksetRepo := mock_repo.NewMockWorksetRepo()
	worksetRepo.Infos["workset-1"] = model.WorksetInfo{ID: "workset-1", TeamID: "team-1"}
	eventBus := newMockEventBus()
	txnMgr := mock_repo.NewMockTxnMgr(newMockTxnContext(mockTxnRepos{comic: comicRepo, member: memberRepo, workset: worksetRepo}))

	app := NewComicApp(service.NewComicService(), memberRepo, worksetRepo, comicRepo, txnMgr, mock_repo.NewMockOSSMessageRepo(), eventBus, newMockOSSClient())

	res, err := app.Create(background(), "user-1", &val.CreateComicArgs{WorksetID: "workset-1", Title: "Title", Author: "Author", Desc: "Desc"})
	requireNoErr(t, err)

	created := comicRepo.Infos[res.ID]
	if created.Title != "Title" || created.Index != 0 {
		t.Fatalf("unexpected comic: %#v", created)
	}
	pubCalls := eventBus.PubCalls()
	createdEvent, ok := pubCalls[0][0].(*event.ComicCreatedEvent)
	if !ok || createdEvent.WorksetID != "workset-1" {
		t.Fatalf("unexpected event: %#v", pubCalls)
	}
}

func TestComicAppRemoveUsesMockTxnRepos(t *testing.T) {
	comicRepo := mock_repo.NewMockComicRepo()
	comicRepo.Infos["comic-1"] = model.ComicInfo{ID: "comic-1", WorksetID: "workset-1", CoverOSSKey: "comic-cover-1", IsCoverUploaded: true}
	memberRepo := mock_repo.NewMockMemberRepo()
	memberRepo.Infos["member-1"] = *adminMember()
	worksetRepo := mock_repo.NewMockWorksetRepo()
	worksetRepo.Infos["workset-1"] = model.WorksetInfo{ID: "workset-1", TeamID: "team-1"}
	eventBus := newMockEventBus()
	msgRepo := mock_repo.NewMockOSSMessageRepo()
	txnMgr := mock_repo.NewMockTxnMgr(newMockTxnContext(mockTxnRepos{comic: comicRepo, workset: worksetRepo, ossMessage: msgRepo}))
	ossClient := newMockOSSClient()

	app := NewComicApp(service.NewComicService(), memberRepo, worksetRepo, comicRepo, txnMgr, msgRepo, eventBus, ossClient)

	err := app.Remove(background(), "user-1", "comic-1")
	requireNoErr(t, err)

	if _, ok := comicRepo.Infos["comic-1"]; ok {
		t.Fatalf("expected comic deleted: %#v", comicRepo.Infos)
	}
	pubCalls := eventBus.PubCalls()
	removedEvent, ok := pubCalls[0][0].(*event.ComicRemovedEvent)
	if !ok || removedEvent.WorksetID != "workset-1" {
		t.Fatalf("unexpected event: %#v", pubCalls)
	}
	if len(msgRepo.Messages) != 1 || msgRepo.Messages[0].ObjectKey != "comic-cover-1" {
		t.Fatalf("expected cover oss cleanup message, got %#v", msgRepo.Messages)
	}
}

func TestComicAppRemoveFailsWhenCoverCleanupFails(t *testing.T) {
	comicRepo := mock_repo.NewMockComicRepo()
	comicRepo.Infos["comic-1"] = model.ComicInfo{ID: "comic-1", WorksetID: "workset-1", CoverOSSKey: "comic-cover-1", IsCoverUploaded: true}
	memberRepo := mock_repo.NewMockMemberRepo()
	memberRepo.Infos["member-1"] = *adminMember()
	worksetRepo := mock_repo.NewMockWorksetRepo()
	worksetRepo.Infos["workset-1"] = model.WorksetInfo{ID: "workset-1", TeamID: "team-1"}
	msgRepo := mock_repo.NewMockOSSMessageRepo()
	msgRepo.InsertErr = errors.New("boom")
	txnMgr := mock_repo.NewMockTxnMgr(newMockTxnContext(mockTxnRepos{comic: comicRepo, workset: worksetRepo, ossMessage: msgRepo}))

	app := NewComicApp(service.NewComicService(), memberRepo, worksetRepo, comicRepo, txnMgr, msgRepo, newMockEventBus(), newMockOSSClient())

	err := app.Remove(background(), "user-1", "comic-1")
	if err == nil {
		t.Fatal("expected remove failure when cover cleanup fails")
	}

	if len(msgRepo.Messages) != 0 {
		t.Fatalf("expected no queued messages, got %#v", msgRepo.Messages)
	}
}

func TestComicAppListForbidden(t *testing.T) {
	worksetRepo := mock_repo.NewMockWorksetRepo()
	worksetRepo.Infos["workset-1"] = model.WorksetInfo{ID: "workset-1", TeamID: "team-1"}

	app := NewComicApp(service.NewComicService(), mock_repo.NewMockMemberRepo(), worksetRepo, mock_repo.NewMockComicRepo(), mock_repo.NewMockTxnMgr(nil), mock_repo.NewMockOSSMessageRepo(), newMockEventBus(), newMockOSSClient())

	if _, err := app.List(background(), "user-1", &val.ListComicArgs{WorksetID: "workset-1"}); err == nil {
		t.Fatal("expected forbidden error")
	}
}

func TestComicAppUpdateSuccess(t *testing.T) {
	memberRepo := mock_repo.NewMockMemberRepo()
	memberRepo.Infos["member-1"] = *adminMember()
	worksetRepo := mock_repo.NewMockWorksetRepo()
	worksetRepo.Infos["workset-1"] = model.WorksetInfo{ID: "workset-1", TeamID: "team-1"}
	comicRepo := mock_repo.NewMockComicRepo()
	comicRepo.Infos["comic-1"] = model.ComicInfo{ID: "comic-1", WorksetID: "workset-1", Title: "Old", Author: "A"}

	app := NewComicApp(service.NewComicService(), memberRepo, worksetRepo, comicRepo, mock_repo.NewMockTxnMgr(nil), mock_repo.NewMockOSSMessageRepo(), newMockEventBus(), newMockOSSClient())
	err := app.Update(background(), "user-1", &val.UpdateComicArgs{ID: "comic-1", Title: "New", Author: "B", Desc: "Desc"})
	requireNoErr(t, err)

	updated := comicRepo.Infos["comic-1"]
	if updated.Title != "New" || updated.Author != "B" || updated.Desc != "Desc" {
		t.Fatalf("unexpected comic update: %#v", updated)
	}
}

func TestComicAppErrorPaths(t *testing.T) {
	t.Run("create rejects non admin", func(t *testing.T) {
		memberRepo := mock_repo.NewMockMemberRepo()
		memberRepo.Infos["member-1"] = model.MemberInfo{ID: "member-1", UserID: "user-1", TeamID: "team-1"}
		worksetRepo := mock_repo.NewMockWorksetRepo()
		worksetRepo.Infos["workset-1"] = model.WorksetInfo{ID: "workset-1", TeamID: "team-1"}
		comicRepo := mock_repo.NewMockComicRepo()
		app := NewComicApp(service.NewComicService(), memberRepo, worksetRepo, comicRepo, mock_repo.NewMockTxnMgr(newMockTxnContext(mockTxnRepos{member: memberRepo, workset: worksetRepo, comic: comicRepo})), mock_repo.NewMockOSSMessageRepo(), newMockEventBus(), newMockOSSClient())
		if _, err := app.Create(background(), "user-1", &val.CreateComicArgs{WorksetID: "workset-1", Title: "T", Author: "A"}); err == nil {
			t.Fatal("expected forbidden create error")
		}
	})

	t.Run("update rejects missing comic", func(t *testing.T) {
		app := NewComicApp(service.NewComicService(), mock_repo.NewMockMemberRepo(), mock_repo.NewMockWorksetRepo(), mock_repo.NewMockComicRepo(), mock_repo.NewMockTxnMgr(nil), mock_repo.NewMockOSSMessageRepo(), newMockEventBus(), newMockOSSClient())
		if err := app.Update(background(), "user-1", &val.UpdateComicArgs{ID: "missing", Title: "T", Author: "A"}); err == nil {
			t.Fatal("expected missing comic error")
		}
	})

	t.Run("remove rejects non admin", func(t *testing.T) {
		memberRepo := mock_repo.NewMockMemberRepo()
		memberRepo.Infos["member-1"] = model.MemberInfo{ID: "member-1", UserID: "user-1", TeamID: "team-1"}
		worksetRepo := mock_repo.NewMockWorksetRepo()
		worksetRepo.Infos["workset-1"] = model.WorksetInfo{ID: "workset-1", TeamID: "team-1"}
		comicRepo := mock_repo.NewMockComicRepo()
		comicRepo.Infos["comic-1"] = model.ComicInfo{ID: "comic-1", WorksetID: "workset-1"}
		app := NewComicApp(service.NewComicService(), memberRepo, worksetRepo, comicRepo, mock_repo.NewMockTxnMgr(nil), mock_repo.NewMockOSSMessageRepo(), newMockEventBus(), newMockOSSClient())
		if err := app.Remove(background(), "user-1", "comic-1"); err == nil {
			t.Fatal("expected forbidden remove error")
		}
	})
}
