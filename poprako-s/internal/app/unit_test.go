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

func TestUnitAppList(t *testing.T) {
	pageRepo := &mock_repo.PageRepo{Infos: map[string]model.PageInfo{
		"page-1": {ID: "page-1", ChapterID: "chapter-1"},
	}}
	assignment := reviewerAssignment()
	assignmentRepo := &mock_repo.AssignmentRepo{Infos: map[string]model.AssignmentInfo{
		assignment.ID: *assignment,
	}}
	unitRepo := &mock_repo.UnitRepo{Infos: map[string]model.UnitInfo{
		"unit-1": {ID: "unit-1", PageID: "page-1", Index: 1},
	}}

	app := NewUnitApp(service.NewUnitService(), newMockEventBus(), &mock_repo.UserRepo{}, pageRepo, &mock_repo.ChapterRepo{}, assignmentRepo, unitRepo, &mock_repo.TxnMgr{})

	got, err := app.List(background(), "user-1", "page-1")
	requireNoErr(t, err)

	if len(got) != 1 || got[0].ID != "unit-1" {
		t.Fatalf("unexpected result: %#v", got)
	}
}

func TestUnitAppSaveUsesMockTxnRepos(t *testing.T) {
	userRepo := mock_repo.NewMockUserRepo()
	userRepo.Infos["user-1"] = *normalUser()
	pageRepo := mock_repo.NewMockPageRepo()
	pageRepo.Infos["page-1"] = model.PageInfo{ID: "page-1", ChapterID: "chapter-1"}
	chapterRepo := mock_repo.NewMockChapterRepo()
	assignmentRepo := mock_repo.NewMockAssignmentRepo()
	assignmentRepo.Infos["assignment-1"] = *proofreaderAssignment()
	unitRepo := mock_repo.NewMockUnitRepo()
	translated := "translated"
	proofread := "proofread"
	unitRepo.Infos["unit-old-1"] = model.UnitInfo{ID: "unit-old-1", PageID: "page-1", Index: 0, TranslatedText: &translated, IsProofread: true, ProofreadText: &proofread}
	unitRepo.Infos["unit-old-2"] = model.UnitInfo{ID: "unit-old-2", PageID: "page-1", Index: 1}
	eventBus := newMockEventBus()
	txnMgr := mock_repo.NewMockTxnMgr(newMockTxnContext(mockTxnRepos{unit: unitRepo, page: pageRepo, chapter: chapterRepo}))

	app := NewUnitApp(service.NewUnitService(), eventBus, userRepo, pageRepo, chapterRepo, assignmentRepo, unitRepo, txnMgr)

	err := app.Save(background(), "user-1", &val.SavePageUnitArgs{
		PageID: "page-1",
		UnitDiff: val.UnitDiff{
			Insert: []val.UnitCreation{{ID: "unit-new", Index: 2, TranslatedText: ptr("new text"), IsProofread: true, ProofreadText: ptr("new proof")}},
			Patch:  []val.UnitPatch{{ID: "unit-old-2", TranslatedText: ptrptr("patched text"), IsProofread: ptr(true), ProofreadText: ptrptr("patched proof")}},
			Delete: []string{"unit-old-1"},
		},
	})
	requireNoErr(t, err)

	if _, ok := unitRepo.Infos["unit-old-1"]; ok {
		t.Fatalf("expected deleted unit, got %#v", unitRepo.Infos)
	}
	if unitRepo.Infos["unit-old-2"].TranslatedText == nil || !unitRepo.Infos["unit-old-2"].IsProofread {
		t.Fatalf("expected patched unit, got %#v", unitRepo.Infos["unit-old-2"])
	}
	if _, ok := unitRepo.Infos["unit-new"]; !ok {
		t.Fatalf("expected inserted unit, got %#v", unitRepo.Infos)
	}

	pubCalls := eventBus.PubCalls()
	unitSaveEvent, ok := pubCalls[0][0].(*event.UnitSaveEvent)
	if !ok || unitSaveEvent.TotalDelta != 0 || unitSaveEvent.TranslatedDelta != 1 || unitSaveEvent.ProofreadDelta != 1 {
		t.Fatalf("unexpected unit save event: %#v", pubCalls)
	}
}

func TestUnitAppListForbidden(t *testing.T) {
	pageRepo := mock_repo.NewMockPageRepo()
	pageRepo.Infos["page-1"] = model.PageInfo{ID: "page-1", ChapterID: "chapter-1"}

	app := NewUnitApp(service.NewUnitService(), newMockEventBus(), mock_repo.NewMockUserRepo(), pageRepo, mock_repo.NewMockChapterRepo(), mock_repo.NewMockAssignmentRepo(), mock_repo.NewMockUnitRepo(), mock_repo.NewMockTxnMgr(nil))
	if _, err := app.List(background(), "user-1", "page-1"); err == nil {
		t.Fatal("expected forbidden error")
	}
}

func TestUnitAppSaveErrorPaths(t *testing.T) {
	t.Run("save rejects missing user", func(t *testing.T) {
		app := NewUnitApp(service.NewUnitService(), newMockEventBus(), mock_repo.NewMockUserRepo(), mock_repo.NewMockPageRepo(), mock_repo.NewMockChapterRepo(), mock_repo.NewMockAssignmentRepo(), mock_repo.NewMockUnitRepo(), mock_repo.NewMockTxnMgr(nil))
		if err := app.Save(background(), "user-1", &val.SavePageUnitArgs{PageID: "page-1"}); err == nil {
			t.Fatal("expected missing user error")
		}
	})

	t.Run("save rejects missing page", func(t *testing.T) {
		userRepo := mock_repo.NewMockUserRepo()
		userRepo.Infos["user-1"] = *normalUser()
		app := NewUnitApp(service.NewUnitService(), newMockEventBus(), userRepo, mock_repo.NewMockPageRepo(), mock_repo.NewMockChapterRepo(), mock_repo.NewMockAssignmentRepo(), mock_repo.NewMockUnitRepo(), mock_repo.NewMockTxnMgr(nil))
		if err := app.Save(background(), "user-1", &val.SavePageUnitArgs{PageID: "page-1"}); err == nil {
			t.Fatal("expected missing page error")
		}
	})

	t.Run("save wraps event publish failures", func(t *testing.T) {
		userRepo := mock_repo.NewMockUserRepo()
		userRepo.Infos["user-1"] = *normalUser()
		pageRepo := mock_repo.NewMockPageRepo()
		pageRepo.Infos["page-1"] = model.PageInfo{ID: "page-1", ChapterID: "chapter-1"}
		chapterRepo := mock_repo.NewMockChapterRepo()
		assignmentRepo := mock_repo.NewMockAssignmentRepo()
		assignmentRepo.Infos["assignment-1"] = *proofreaderAssignment()
		unitRepo := mock_repo.NewMockUnitRepo()
		eventBus := newMockEventBus()
		eventBus.SetPubErr(errors.New("boom"))
		txnMgr := mock_repo.NewMockTxnMgr(newMockTxnContext(mockTxnRepos{unit: unitRepo, page: pageRepo, chapter: chapterRepo}))
		app := NewUnitApp(service.NewUnitService(), eventBus, userRepo, pageRepo, chapterRepo, assignmentRepo, unitRepo, txnMgr)
		if err := app.Save(background(), "user-1", &val.SavePageUnitArgs{PageID: "page-1"}); err == nil {
			t.Fatal("expected publish failure")
		}
	})
}
