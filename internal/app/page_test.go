package app

import (
	"errors"
	"testing"

	"poprako-s/internal/app/val"
	"poprako-s/internal/domain/model"
	"poprako-s/internal/domain/service"
	mock_repo "poprako-s/internal/infra/repo/mock"
)

func TestPageAppList(t *testing.T) {
	assignment := reviewerAssignment()
	assignmentRepo := &mock_repo.AssignmentRepo{Infos: map[string]model.AssignmentInfo{
		assignment.ID: *assignment,
	}}
	pageRepo := &mock_repo.PageRepo{Infos: map[string]model.PageInfo{
		"page-1": {ID: "page-1", ChapterID: "chapter-1", Index: 0},
	}}

	app := NewPageApp(service.NewPageService(), assignmentRepo, &mock_repo.ChapterRepo{}, pageRepo, mock_repo.NewMockTxnMgr(nil), mock_repo.NewMockOSSMessageRepo(), newMockOSSClient())

	got, err := app.List(background(), "user-1", &val.ListChapterPageArgs{ChapterID: "chapter-1"})
	requireNoErr(t, err)

	if len(got) != 1 || got[0].ID != "page-1" {
		t.Fatalf("unexpected result: %#v", got)
	}
}

func TestPageAppReserveUpdateAndRemove(t *testing.T) {
	assignment := rawProviderAssignment()
	assignmentRepo := mock_repo.NewMockAssignmentRepo()
	assignmentRepo.Infos[assignment.ID] = *assignment
	chapterRepo := mock_repo.NewMockChapterRepo()
	chapterRepo.Infos["chapter-1"] = model.ChapterInfo{ID: "chapter-1", ComicID: "comic-1", PageCount: 0}
	pageRepo := mock_repo.NewMockPageRepo()
	ossClient := newMockOSSClient()
	ossClient.SetPutURLs(map[string]string{"chapter_chapter-1/page_0": "https://upload.example/page_0", "chapter_chapter-1/page_1": "https://upload.example/page_1"})
	msgRepo := mock_repo.NewMockOSSMessageRepo()
	txCx := newMockTxnContext(mockTxnRepos{chapter: chapterRepo, page: pageRepo, ossMessage: msgRepo})
	txnMgr := mock_repo.NewMockTxnMgr(txCx)

	app := NewPageApp(service.NewPageService(), assignmentRepo, chapterRepo, pageRepo, txnMgr, msgRepo, ossClient)

	reserveRes, err := app.Reserve(background(), "user-1", &val.ReserveChapterPagesArgs{ChapterID: "chapter-1", PageCount: 2, Extension: "png"})
	requireNoErr(t, err)

	if len(reserveRes.Creations) != 2 || len(pageRepo.Infos) != 2 {
		t.Fatalf("unexpected reserve result: %#v %#v", reserveRes, pageRepo.Infos)
	}
	if chapterRepo.Infos["chapter-1"].PageCount != 2 {
		t.Fatalf("expected chapter page count to increase, got %#v", chapterRepo.Infos["chapter-1"])
	}

	pageID := reserveRes.Creations[0].PageID
	err = app.Update(background(), "user-1", &val.UpdatePageArgs{ID: pageID, IsUploaded: true})
	requireNoErr(t, err)
	if !pageRepo.Infos[pageID].IsUploaded {
		t.Fatalf("expected page uploaded: %#v", pageRepo.Infos[pageID])
	}

	err = app.Remove(background(), "user-1", pageID)
	requireNoErr(t, err)
	if _, ok := pageRepo.Infos[pageID]; ok {
		t.Fatalf("expected page deleted, got %#v", pageRepo.Infos)
	}
	if chapterRepo.Infos["chapter-1"].PageCount != 1 {
		t.Fatalf("expected chapter page count to decrease, got %#v", chapterRepo.Infos["chapter-1"])
	}
	if len(msgRepo.Messages) == 0 {
		t.Fatalf("expected oss enqueue message, got %#v", msgRepo.Messages)
	}
}

func TestPageAppListForbidden(t *testing.T) {
	pageRepo := mock_repo.NewMockPageRepo()
	pageRepo.Infos["page-1"] = model.PageInfo{ID: "page-1", ChapterID: "chapter-1"}

	app := NewPageApp(service.NewPageService(), mock_repo.NewMockAssignmentRepo(), mock_repo.NewMockChapterRepo(), pageRepo, mock_repo.NewMockTxnMgr(nil), mock_repo.NewMockOSSMessageRepo(), newMockOSSClient())
	if _, err := app.List(background(), "user-1", &val.ListChapterPageArgs{ChapterID: "chapter-1"}); err == nil {
		t.Fatal("expected forbidden error")
	}
}

func TestPageAppErrorPaths(t *testing.T) {
	t.Run("reserve wraps put url failures", func(t *testing.T) {
		assignmentRepo := mock_repo.NewMockAssignmentRepo()
		assignmentRepo.Infos["assignment-1"] = *rawProviderAssignment()
		ossClient := newMockOSSClient()
		ossClient.SetPutErr(errors.New("boom"))
		app := NewPageApp(service.NewPageService(), assignmentRepo, mock_repo.NewMockChapterRepo(), mock_repo.NewMockPageRepo(), mock_repo.NewMockTxnMgr(nil), mock_repo.NewMockOSSMessageRepo(), ossClient)
		if _, err := app.Reserve(background(), "user-1", &val.ReserveChapterPagesArgs{ChapterID: "chapter-1", PageCount: 1, Extension: "png"}); err == nil {
			t.Fatal("expected reserve failure")
		}
	})

	t.Run("update rejects missing page", func(t *testing.T) {
		app := NewPageApp(service.NewPageService(), mock_repo.NewMockAssignmentRepo(), mock_repo.NewMockChapterRepo(), mock_repo.NewMockPageRepo(), mock_repo.NewMockTxnMgr(nil), mock_repo.NewMockOSSMessageRepo(), newMockOSSClient())
		if err := app.Update(background(), "user-1", &val.UpdatePageArgs{ID: "missing", IsUploaded: true}); err == nil {
			t.Fatal("expected missing page error")
		}
	})

	t.Run("remove rejects missing page", func(t *testing.T) {
		app := NewPageApp(service.NewPageService(), mock_repo.NewMockAssignmentRepo(), mock_repo.NewMockChapterRepo(), mock_repo.NewMockPageRepo(), mock_repo.NewMockTxnMgr(nil), mock_repo.NewMockOSSMessageRepo(), newMockOSSClient())
		if err := app.Remove(background(), "user-1", "missing"); err == nil {
			t.Fatal("expected missing page error")
		}
	})

	t.Run("remove returns error when oss enqueue fails", func(t *testing.T) {
		assignmentRepo := mock_repo.NewMockAssignmentRepo()
		assignmentRepo.Infos["assignment-1"] = *rawProviderAssignment()
		pageRepo := mock_repo.NewMockPageRepo()
		pageRepo.Infos["page-1"] = model.PageInfo{ID: "page-1", ChapterID: "chapter-1", OSSKey: "page-oss-1"}
		chapterRepo := mock_repo.NewMockChapterRepo()
		chapterRepo.Infos["chapter-1"] = model.ChapterInfo{ID: "chapter-1", ComicID: "comic-1", PageCount: 1}
		msgRepo := mock_repo.NewMockOSSMessageRepo()
		msgRepo.InsertErr = errors.New("boom")
		txCx := newMockTxnContext(mockTxnRepos{page: pageRepo, chapter: chapterRepo, ossMessage: msgRepo})
		txnMgr := mock_repo.NewMockTxnMgr(txCx)

		app := NewPageApp(service.NewPageService(), assignmentRepo, chapterRepo, pageRepo, txnMgr, msgRepo, newMockOSSClient())

		err := app.Remove(background(), "user-1", "page-1")
		if err == nil {
			t.Fatal("expected remove failure")
		}

		if len(msgRepo.Messages) != 0 {
			t.Fatalf("expected no queued messages, got %#v", msgRepo.Messages)
		}
	})

	t.Run("assemble swallows image url errors", func(t *testing.T) {
		ossClient := newMockOSSClient()
		ossClient.SetGetErr(errors.New("boom"))
		info := assemblePageInfo(&model.PageInfo{ID: "page-1", OSSKey: "oss-key"}, ossClient)
		if info.ImageURL != "" {
			t.Fatalf("expected empty image url, got %#v", info)
		}
	})
}
