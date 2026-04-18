package app

import (
	"errors"
	"testing"
	"time"

	"poprako-s/internal/app/val"
	"poprako-s/internal/domain/event"
	"poprako-s/internal/domain/model"
	"poprako-s/internal/domain/service"
	mock_repo "poprako-s/internal/infra/repo/mock"
)

func TestAssignmentAppListMy(t *testing.T) {
	now := time.Now()
	repo := &mock_repo.AssignmentRepo{Infos: map[string]model.AssignmentInfo{
		"assign-1": {ID: "assign-1", ChapterID: "chapter-1", UserID: "user-1", AssignedReviewerAt: &now},
	}}

	app := NewAssignmentApp(service.NewAssignmentService(), repo, mock_repo.NewMockChapterInvitationRepo(), &mock_repo.ChapterRepo{}, mock_repo.NewMockComicRepo(), mock_repo.NewMockUserRepo(), &mock_repo.TxnMgr{}, newMockEventBus(), newMockOSSClient())

	got, err := app.ListMy(background(), "user-1", &val.ListMyAssignmentArgs{})
	requireNoErr(t, err)

	if len(got) != 1 || got[0].ID != "assign-1" || got[0].Roles != model.RoleMask(model.RoleReviewer) {
		t.Fatalf("unexpected result: %#v", got)
	}
}

func TestAssignmentAppListByChapterIncludesNestedData(t *testing.T) {
	now := time.Now()
	reviewer := reviewerAssignment()
	reviewer.CreatedAt = now
	reviewer.UpdatedAt = now

	repo := mock_repo.NewMockAssignmentRepo()
	repo.Infos[reviewer.ID] = *reviewer
	repo.Infos["assignment-2"] = model.AssignmentInfo{
		ID:                   "assignment-2",
		ChapterID:            "chapter-1",
		UserID:               "user-2",
		AssignedTranslatorAt: &now,
		Chapter: &model.ChapterInfo{
			ID:        "chapter-1",
			ComicID:   "comic-1",
			Subtitle:  "Chapter 1",
			CreatedAt: now,
			UpdatedAt: now,
		},
		User: &model.UserInfo{
			ID:               "user-2",
			Name:             "User 2",
			QQ:               "10002",
			AvatarKey:        "avatar-user-2",
			IsAvatarUploaded: true,
			LastLoginAt:      now,
			CreatedAt:        now,
			UpdatedAt:        now,
		},
		CreatedAt: now.Add(time.Second),
		UpdatedAt: now.Add(time.Second),
	}
	ossClient := newMockOSSClient()
	ossClient.SetGetURLs(map[string]string{"avatar-user-2": "https://cdn.example/avatar-user-2"})

	app := NewAssignmentApp(
		service.NewAssignmentService(),
		repo,
		mock_repo.NewMockChapterInvitationRepo(),
		mock_repo.NewMockChapterRepo(),
		mock_repo.NewMockComicRepo(),
		mock_repo.NewMockUserRepo(),
		mock_repo.NewMockTxnMgr(nil),
		newMockEventBus(),
		ossClient,
	)

	got, err := app.ListByChapter(background(), "user-1", &val.ListChapterAssignmentArgs{ChapterID: "chapter-1"})
	requireNoErr(t, err)

	if len(got) != 2 {
		t.Fatalf("unexpected assignments: %#v", got)
	}
	if got[1].User == nil || got[1].User.AvatarURL != "https://cdn.example/avatar-user-2" {
		t.Fatalf("expected nested user avatar url, got %#v", got[1].User)
	}
	if got[1].Chapter == nil || got[1].Chapter.Subtitle != "Chapter 1" {
		t.Fatalf("expected nested chapter info, got %#v", got[1].Chapter)
	}
	if got[1].Roles != model.RoleMask(model.RoleTranslator) {
		t.Fatalf("unexpected roles: %#v", got[1].Roles)
	}
}

func TestAssignmentAppCreateUsesMockTxnRepos(t *testing.T) {
	assignmentRepo := mock_repo.NewMockAssignmentRepo()
	assignmentRepo.Infos["assignment-reviewer"] = *reviewerAssignment()
	userRepo := mock_repo.NewMockUserRepo()
	eventBus := newMockEventBus()
	txnMgr := mock_repo.NewMockTxnMgr(newMockTxnContext(mockTxnRepos{assignment: assignmentRepo, user: userRepo}))

	app := NewAssignmentApp(
		service.NewAssignmentService(),
		assignmentRepo,
		mock_repo.NewMockChapterInvitationRepo(),
		mock_repo.NewMockChapterRepo(),
		mock_repo.NewMockComicRepo(),
		userRepo,
		txnMgr,
		eventBus,
		newMockOSSClient(),
	)

	res, err := app.Create(background(), "user-1", &val.CreateAssignmentArgs{
		ChapterID: "chapter-1",
		UserID:    "user-2",
		Roles:     model.RoleMask(model.RoleTranslator) | model.RoleMask(model.RolePublisher),
	})
	requireNoErr(t, err)

	created, ok := assignmentRepo.Infos[res.ID]
	if !ok {
		t.Fatalf("expected created assignment in repo, got %#v", assignmentRepo.Infos)
	}
	if created.UserID != "user-2" || created.AssignedTranslatorAt == nil || created.AssignedPublisherAt == nil {
		t.Fatalf("unexpected created assignment: %#v", created)
	}
	pubCalls := eventBus.PubCalls()
	if len(pubCalls) != 1 {
		t.Fatalf("expected one sync event, got %#v", pubCalls)
	}
	createdEvent, ok := pubCalls[0][0].(*event.AssignmentCreatedEvent)
	if !ok || createdEvent.UserID != "user-2" || createdEvent.ChapterID != "chapter-1" {
		t.Fatalf("unexpected event: %#v", pubCalls)
	}
}

func TestAssignmentAppUpdateReplacesRoles(t *testing.T) {
	now := time.Now()
	assignmentRepo := mock_repo.NewMockAssignmentRepo()
	assignmentRepo.Infos["assignment-reviewer"] = *reviewerAssignment()
	assignmentRepo.Infos["assignment-target"] = model.AssignmentInfo{
		ID:                   "assignment-target",
		ChapterID:            "chapter-1",
		UserID:               "user-2",
		AssignedTranslatorAt: &now,
		CreatedAt:            now,
		UpdatedAt:            now,
	}

	app := NewAssignmentApp(
		service.NewAssignmentService(),
		assignmentRepo,
		mock_repo.NewMockChapterInvitationRepo(),
		mock_repo.NewMockChapterRepo(),
		mock_repo.NewMockComicRepo(),
		mock_repo.NewMockUserRepo(),
		mock_repo.NewMockTxnMgr(nil),
		newMockEventBus(),
		newMockOSSClient(),
	)

	err := app.Update(background(), "user-1", &val.UpdateAssignmentArgs{ID: "assignment-target", Roles: model.RoleMask(model.RolePublisher)})
	requireNoErr(t, err)

	updated := assignmentRepo.Infos["assignment-target"]
	if updated.AssignedTranslatorAt != nil || updated.AssignedPublisherAt == nil {
		t.Fatalf("unexpected updated assignment: %#v", updated)
	}
}

func TestAssignmentAppRemoveUsesMockTxnRepos(t *testing.T) {
	assignmentRepo := mock_repo.NewMockAssignmentRepo()
	assignmentRepo.Infos["assignment-reviewer"] = *reviewerAssignment()
	assignmentRepo.Infos["assignment-target"] = model.AssignmentInfo{ID: "assignment-target", ChapterID: "chapter-1", UserID: "user-2"}
	chapterRepo := mock_repo.NewMockChapterRepo()
	chapterRepo.Infos["chapter-1"] = model.ChapterInfo{ID: "chapter-1", ComicID: "comic-1"}
	userRepo := mock_repo.NewMockUserRepo()
	eventBus := newMockEventBus()
	txnMgr := mock_repo.NewMockTxnMgr(newMockTxnContext(mockTxnRepos{assignment: assignmentRepo, chapter: chapterRepo, user: userRepo}))

	app := NewAssignmentApp(
		service.NewAssignmentService(),
		assignmentRepo,
		mock_repo.NewMockChapterInvitationRepo(),
		chapterRepo,
		mock_repo.NewMockComicRepo(),
		userRepo,
		txnMgr,
		eventBus,
		newMockOSSClient(),
	)

	err := app.Remove(background(), "user-1", "assignment-target")
	requireNoErr(t, err)

	if _, ok := assignmentRepo.Infos["assignment-target"]; ok {
		t.Fatalf("expected assignment deleted, got %#v", assignmentRepo.Infos)
	}
	pubCalls := eventBus.PubCalls()
	removedEvent, ok := pubCalls[0][0].(*event.AssignmentRemovedEvent)
	if !ok || removedEvent.UserID != "user-2" || removedEvent.WasPublished {
		t.Fatalf("unexpected remove event: %#v", pubCalls)
	}
}

func TestAssignmentAppListByChapterForbidden(t *testing.T) {
	repo := mock_repo.NewMockAssignmentRepo()
	repo.Infos["assignment-2"] = model.AssignmentInfo{ID: "assignment-2", ChapterID: "chapter-1", UserID: "user-2"}

	app := NewAssignmentApp(
		service.NewAssignmentService(),
		repo,
		mock_repo.NewMockChapterInvitationRepo(),
		mock_repo.NewMockChapterRepo(),
		mock_repo.NewMockComicRepo(),
		mock_repo.NewMockUserRepo(),
		mock_repo.NewMockTxnMgr(nil),
		newMockEventBus(),
		newMockOSSClient(),
	)

	if _, err := app.ListByChapter(background(), "user-1", &val.ListChapterAssignmentArgs{ChapterID: "chapter-1"}); err == nil {
		t.Fatal("expected forbidden error")
	}
}

func TestAssignmentAppCreateUpdateAndRemoveErrorPaths(t *testing.T) {
	t.Run("create wraps transaction errors", func(t *testing.T) {
		app := NewAssignmentApp(
			service.NewAssignmentService(),
			mock_repo.NewMockAssignmentRepo(),
			mock_repo.NewMockChapterInvitationRepo(),
			mock_repo.NewMockChapterRepo(),
			mock_repo.NewMockComicRepo(),
			mock_repo.NewMockUserRepo(),
			mock_repo.NewMockTxnMgr(nil),
			newMockEventBus(),
			newMockOSSClient(),
		)

		if _, err := app.Create(background(), "user-1", &val.CreateAssignmentArgs{ChapterID: "chapter-1", UserID: "user-2", Roles: model.RoleMask(model.RoleTranslator)}); err == nil {
			t.Fatal("expected create failure")
		}
	})

	t.Run("update rejects missing target", func(t *testing.T) {
		app := NewAssignmentApp(
			service.NewAssignmentService(),
			mock_repo.NewMockAssignmentRepo(),
			mock_repo.NewMockChapterInvitationRepo(),
			mock_repo.NewMockChapterRepo(),
			mock_repo.NewMockComicRepo(),
			mock_repo.NewMockUserRepo(),
			mock_repo.NewMockTxnMgr(nil),
			newMockEventBus(),
			newMockOSSClient(),
		)

		if err := app.Update(background(), "user-1", &val.UpdateAssignmentArgs{ID: "missing", Roles: model.RoleMask(model.RoleReviewer)}); err == nil {
			t.Fatal("expected missing assignment error")
		}
	})

	t.Run("update rejects unauthorized editor", func(t *testing.T) {
		assignmentRepo := mock_repo.NewMockAssignmentRepo()
		assignmentRepo.Infos["assignment-target"] = model.AssignmentInfo{ID: "assignment-target", ChapterID: "chapter-1", UserID: "user-2"}
		app := NewAssignmentApp(
			service.NewAssignmentService(),
			assignmentRepo,
			mock_repo.NewMockChapterInvitationRepo(),
			mock_repo.NewMockChapterRepo(),
			mock_repo.NewMockComicRepo(),
			mock_repo.NewMockUserRepo(),
			mock_repo.NewMockTxnMgr(nil),
			newMockEventBus(),
			newMockOSSClient(),
		)

		if err := app.Update(background(), "user-1", &val.UpdateAssignmentArgs{ID: "assignment-target", Roles: model.RoleMask(model.RolePublisher)}); err == nil {
			t.Fatal("expected forbidden update error")
		}
	})

	t.Run("remove returns permission error", func(t *testing.T) {
		assignmentRepo := mock_repo.NewMockAssignmentRepo()
		assignmentRepo.Infos["assignment-target"] = model.AssignmentInfo{ID: "assignment-target", ChapterID: "chapter-1", UserID: "user-2"}
		chapterRepo := mock_repo.NewMockChapterRepo()
		chapterRepo.Infos["chapter-1"] = model.ChapterInfo{ID: "chapter-1", ComicID: "comic-1"}
		userRepo := mock_repo.NewMockUserRepo()
		txnMgr := mock_repo.NewMockTxnMgr(newMockTxnContext(mockTxnRepos{assignment: assignmentRepo, chapter: chapterRepo, user: userRepo}))
		app := NewAssignmentApp(service.NewAssignmentService(), assignmentRepo, mock_repo.NewMockChapterInvitationRepo(), chapterRepo, mock_repo.NewMockComicRepo(), userRepo, txnMgr, newMockEventBus(), newMockOSSClient())

		err := app.Remove(background(), "user-1", "assignment-target")
		if err == nil || !errors.Is(err, errors.New("权限不足")) && err.Error() != "权限不足" {
			t.Fatalf("expected permission error, got %v", err)
		}
	})
}

func TestAssignmentAppJoinInvitorChapterCreatesAssignmentAndInvalidatesInvitation(t *testing.T) {
	assignmentRepo := mock_repo.NewMockAssignmentRepo()
	chapterInvRepo := mock_repo.NewMockChapterInvitationRepo()
	chapterInvRepo.Infos["cinv-1"] = model.ChapterInvitationInfo{
		ID:              "cinv-1",
		ChapterID:       "chapter-1",
		InviterID:       "user-2",
		InviteeQQ:       "10001",
		InvitationCode:  "654321",
		Pending:         true,
		ToBeTranslator:  true,
		ToBeReviewer:    true,
		ToBeRawProvider: true,
	}
	userRepo := mock_repo.NewMockUserRepo()
	userRepo.Infos["user-1"] = *normalUser()
	eventBus := newMockEventBus()
	txnMgr := mock_repo.NewMockTxnMgr(newMockTxnContext(mockTxnRepos{assignment: assignmentRepo, chapterInvitation: chapterInvRepo, user: userRepo}))

	app := NewAssignmentApp(service.NewAssignmentService(), assignmentRepo, chapterInvRepo, mock_repo.NewMockChapterRepo(), mock_repo.NewMockComicRepo(), userRepo, txnMgr, eventBus, newMockOSSClient())

	err := app.JoinInvitorChapter(background(), "user-1", &val.JoinInvitorChapterArgs{InvitationCode: "654321"})
	requireNoErr(t, err)

	if len(assignmentRepo.Infos) != 1 {
		t.Fatalf("expected one assignment created, got %#v", assignmentRepo.Infos)
	}

	for _, info := range assignmentRepo.Infos {
		if info.ChapterID != "chapter-1" || info.UserID != "user-1" {
			t.Fatalf("unexpected created assignment: %#v", info)
		}
		if info.AssignedTranslatorAt == nil || info.AssignedReviewerAt == nil || info.AssignedRawProviderAt == nil {
			t.Fatalf("expected invited roles assigned: %#v", info)
		}
	}

	if chapterInvRepo.Infos["cinv-1"].Pending {
		t.Fatalf("expected invitation invalidated: %#v", chapterInvRepo.Infos["cinv-1"])
	}

	pubCalls := eventBus.PubCalls()
	if len(pubCalls) != 1 || len(pubCalls[0]) != 1 {
		t.Fatalf("expected one assignment-created event publish, got %#v", pubCalls)
	}

	createdEvent, ok := pubCalls[0][0].(*event.AssignmentCreatedEvent)
	if !ok || createdEvent.UserID != "user-1" || createdEvent.ChapterID != "chapter-1" {
		t.Fatalf("unexpected event: %#v", pubCalls)
	}
}

func TestAssignmentAppJoinInvitorChapterMergesExistingAssignment(t *testing.T) {
	now := time.Now()
	assignmentRepo := mock_repo.NewMockAssignmentRepo()
	assignmentRepo.Infos["assignment-1"] = model.AssignmentInfo{
		ID:                   "assignment-1",
		ChapterID:            "chapter-1",
		UserID:               "user-1",
		AssignedTranslatorAt: &now,
		CreatedAt:            now,
		UpdatedAt:            now,
	}
	chapterInvRepo := mock_repo.NewMockChapterInvitationRepo()
	chapterInvRepo.Infos["cinv-1"] = model.ChapterInvitationInfo{
		ID:             "cinv-1",
		ChapterID:      "chapter-1",
		InviteeQQ:      "10001",
		InvitationCode: "654321",
		Pending:        true,
		ToBeReviewer:   true,
	}
	userRepo := mock_repo.NewMockUserRepo()
	userRepo.Infos["user-1"] = *normalUser()
	eventBus := newMockEventBus()
	txnMgr := mock_repo.NewMockTxnMgr(newMockTxnContext(mockTxnRepos{assignment: assignmentRepo, chapterInvitation: chapterInvRepo, user: userRepo}))

	app := NewAssignmentApp(service.NewAssignmentService(), assignmentRepo, chapterInvRepo, mock_repo.NewMockChapterRepo(), mock_repo.NewMockComicRepo(), userRepo, txnMgr, eventBus, newMockOSSClient())

	err := app.JoinInvitorChapter(background(), "user-1", &val.JoinInvitorChapterArgs{InvitationCode: "654321"})
	requireNoErr(t, err)

	updated := assignmentRepo.Infos["assignment-1"]
	if updated.AssignedTranslatorAt == nil || updated.AssignedReviewerAt == nil {
		t.Fatalf("expected existing and invited roles both present: %#v", updated)
	}

	if chapterInvRepo.Infos["cinv-1"].Pending {
		t.Fatalf("expected invitation invalidated: %#v", chapterInvRepo.Infos["cinv-1"])
	}

	if len(eventBus.PubCalls()) != 0 {
		t.Fatalf("expected no create event when assignment already exists, got %#v", eventBus.PubCalls())
	}
}

func TestAssignmentAppJoinInvitorChapterRejectsInvalidCode(t *testing.T) {
	assignmentRepo := mock_repo.NewMockAssignmentRepo()
	chapterInvRepo := mock_repo.NewMockChapterInvitationRepo()
	chapterInvRepo.Infos["cinv-1"] = model.ChapterInvitationInfo{
		ID:             "cinv-1",
		ChapterID:      "chapter-1",
		InviteeQQ:      "10001",
		InvitationCode: "654321",
		Pending:        true,
		ToBeReviewer:   true,
	}
	userRepo := mock_repo.NewMockUserRepo()
	userRepo.Infos["user-1"] = *normalUser()

	app := NewAssignmentApp(service.NewAssignmentService(), assignmentRepo, chapterInvRepo, mock_repo.NewMockChapterRepo(), mock_repo.NewMockComicRepo(), userRepo, mock_repo.NewMockTxnMgr(newMockTxnContext(mockTxnRepos{assignment: assignmentRepo, chapterInvitation: chapterInvRepo, user: userRepo})), newMockEventBus(), newMockOSSClient())

	err := app.JoinInvitorChapter(background(), "user-1", &val.JoinInvitorChapterArgs{InvitationCode: "000000"})
	if err == nil {
		t.Fatal("expected invalid invitation code error")
	}

	if len(assignmentRepo.Infos) != 0 {
		t.Fatalf("expected no assignment created, got %#v", assignmentRepo.Infos)
	}
}
