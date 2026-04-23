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

func TestChapterAppList(t *testing.T) {
	comicRepo := &mock_repo.ComicRepo{Infos: map[string]model.ComicInfo{
		"comic-1": {ID: "comic-1", WorksetID: "workset-1"},
	}}
	worksetRepo := &mock_repo.WorksetRepo{Infos: map[string]model.WorksetInfo{
		"workset-1": {ID: "workset-1", TeamID: "team-1"},
	}}
	memberRepo := &mock_repo.MemberRepo{Infos: map[string]model.MemberInfo{
		"member-1": {ID: "member-1", UserID: "user-1", TeamID: "team-1"},
	}}
	chapterRepo := &mock_repo.ChapterRepo{Infos: map[string]model.ChapterInfo{
		"chapter-1": {ID: "chapter-1", ComicID: "comic-1", Subtitle: "ch1"},
	}}

	app := NewChapterApp(service.NewChapterService(), service.NewChapterInvitationService(), memberRepo, worksetRepo, comicRepo, chapterRepo, &mock_repo.AssignmentRepo{}, mock_repo.NewMockUserRepo(), &mock_repo.PageRepo{}, mock_repo.NewMockChapterInvitationRepo(), &mock_repo.TxnMgr{}, mock_repo.NewMockOSSMessageRepo(), newMockEventBus(), newMockOSSClient())

	got, err := app.List(background(), "user-1", &val.ListChapterArgs{ComicID: "comic-1"})
	requireNoErr(t, err)

	if len(got) != 1 || got[0].ID != "chapter-1" {
		t.Fatalf("unexpected result: %#v", got)
	}
}

func TestChapterAppCreateUsesMockTxnRepos(t *testing.T) {
	memberRepo := mock_repo.NewMockMemberRepo()
	memberRepo.Infos["member-1"] = *adminMember()
	worksetRepo := mock_repo.NewMockWorksetRepo()
	worksetRepo.Infos["workset-1"] = model.WorksetInfo{ID: "workset-1", TeamID: "team-1"}
	comicRepo := mock_repo.NewMockComicRepo()
	comicRepo.Infos["comic-1"] = model.ComicInfo{ID: "comic-1", WorksetID: "workset-1"}
	chapterRepo := mock_repo.NewMockChapterRepo()
	assignmentRepo := mock_repo.NewMockAssignmentRepo()
	userRepo := mock_repo.NewMockUserRepo()
	msgRepo := mock_repo.NewMockOSSMessageRepo()
	eventBus := newMockEventBus()
	txnMgr := mock_repo.NewMockTxnMgr(newMockTxnContext(mockTxnRepos{chapter: chapterRepo, comic: comicRepo, assignment: assignmentRepo}))

	app := NewChapterApp(service.NewChapterService(), service.NewChapterInvitationService(), memberRepo, worksetRepo, comicRepo, chapterRepo, mock_repo.NewMockAssignmentRepo(), userRepo, mock_repo.NewMockPageRepo(), mock_repo.NewMockChapterInvitationRepo(), txnMgr, msgRepo, eventBus, newMockOSSClient())

	subtitle := "Intro"
	res, err := app.Create(background(), "user-1", &val.CreateChapterArgs{ComicID: "comic-1", Subtitle: &subtitle})
	requireNoErr(t, err)

	created := chapterRepo.Infos[res.ID]
	if created.Index != 0 || created.Subtitle != "Intro" {
		t.Fatalf("unexpected chapter: %#v", created)
	}

	assignment, err := assignmentRepo.Get(model.AssignmentQueryOpt{
		ChapterID: &res.ID,
		UserID:    ptr("user-1"),
	})
	requireNoErr(t, err)

	if !assignment.HasAnyRole(model.RoleReviewer) {
		t.Fatalf("expected creator reviewer assignment, got %#v", assignment)
	}

	pubCalls := eventBus.PubCalls()
	createdEvent, ok := pubCalls[0][0].(*event.ChapterCreatedEvent)
	if !ok || createdEvent.ComicID != "comic-1" {
		t.Fatalf("unexpected event: %#v", pubCalls)
	}
	if len(pubCalls) != 1 || len(pubCalls[0]) != 1 {
		t.Fatalf("unexpected publish calls: %#v", pubCalls)
	}
}

func TestChapterAppPublishUsesMockTxnRepos(t *testing.T) {
	chapterRepo := mock_repo.NewMockChapterRepo()
	chapterRepo.Infos["chapter-1"] = model.ChapterInfo{ID: "chapter-1", ComicID: "comic-1", Subtitle: "To Publish"}
	assignmentRepo := mock_repo.NewMockAssignmentRepo()
	assignmentRepo.Infos["assignment-1"] = *publisherAssignment()
	userRepo := mock_repo.NewMockUserRepo()
	msgRepo := mock_repo.NewMockOSSMessageRepo()
	eventBus := newMockEventBus()
	txnMgr := mock_repo.NewMockTxnMgr(newMockTxnContext(mockTxnRepos{chapter: chapterRepo, assignment: assignmentRepo, user: userRepo}))
	transition := model.WorkflowPublishComplete

	app := NewChapterApp(service.NewChapterService(), service.NewChapterInvitationService(), mock_repo.NewMockMemberRepo(), mock_repo.NewMockWorksetRepo(), mock_repo.NewMockComicRepo(), chapterRepo, assignmentRepo, userRepo, mock_repo.NewMockPageRepo(), mock_repo.NewMockChapterInvitationRepo(), txnMgr, msgRepo, eventBus, newMockOSSClient())

	err := app.Update(background(), "user-1", &val.UpdateChapterArgs{ChapterID: "chapter-1", WorkflowTransition: &transition})
	requireNoErr(t, err)

	updated := chapterRepo.Infos["chapter-1"]
	if updated.PublishedAt == nil {
		t.Fatalf("expected published chapter: %#v", updated)
	}
	pubCalls := eventBus.PubCalls()
	if len(pubCalls) != 1 || len(pubCalls[0]) != 2 {
		t.Fatalf("expected one Pub call with 2 events, got %#v", pubCalls)
	}
	if _, ok := pubCalls[0][0].(*event.ChapterPublishedEvent); !ok {
		t.Fatalf("unexpected sync event: %#v", pubCalls)
	}
	if _, ok := pubCalls[0][1].(*event.WorkflowPublishCompletedEvent); !ok {
		t.Fatalf("unexpected async event: %#v", pubCalls)
	}
}

func TestChapterAppRemoveUsesMockTxnRepos(t *testing.T) {
	chapterRepo := mock_repo.NewMockChapterRepo()
	chapterRepo.Infos["chapter-1"] = model.ChapterInfo{ID: "chapter-1", ComicID: "comic-1"}
	assignmentRepo := mock_repo.NewMockAssignmentRepo()
	assignmentRepo.Infos["assignment-1"] = model.AssignmentInfo{ID: "assignment-1", ChapterID: "chapter-1", UserID: "user-2"}
	assignmentRepo.Infos["assignment-2"] = model.AssignmentInfo{ID: "assignment-2", ChapterID: "chapter-1", UserID: "user-3"}
	pageRepo := mock_repo.NewMockPageRepo()
	pageRepo.Infos["page-1"] = model.PageInfo{ID: "page-1", ChapterID: "chapter-1", OSSKey: "page-oss-1"}
	pageRepo.Infos["page-2"] = model.PageInfo{ID: "page-2", ChapterID: "chapter-1", OSSKey: "page-oss-2"}
	comicRepo := mock_repo.NewMockComicRepo()
	comicRepo.Infos["comic-1"] = model.ComicInfo{ID: "comic-1", WorksetID: "workset-1"}
	worksetRepo := mock_repo.NewMockWorksetRepo()
	worksetRepo.Infos["workset-1"] = model.WorksetInfo{ID: "workset-1", TeamID: "team-1"}
	memberRepo := mock_repo.NewMockMemberRepo()
	memberRepo.Infos["member-1"] = *adminMember()
	userRepo := mock_repo.NewMockUserRepo()
	eventBus := newMockEventBus()
	msgRepo := mock_repo.NewMockOSSMessageRepo()
	txnMgr := mock_repo.NewMockTxnMgr(newMockTxnContext(mockTxnRepos{chapter: chapterRepo, assignment: assignmentRepo, comic: comicRepo, user: userRepo, ossMessage: msgRepo}))

	app := NewChapterApp(service.NewChapterService(), service.NewChapterInvitationService(), memberRepo, worksetRepo, comicRepo, chapterRepo, assignmentRepo, userRepo, pageRepo, mock_repo.NewMockChapterInvitationRepo(), txnMgr, msgRepo, eventBus, newMockOSSClient())

	err := app.Remove(background(), "user-1", "chapter-1")
	requireNoErr(t, err)

	if _, ok := chapterRepo.Infos["chapter-1"]; ok {
		t.Fatalf("expected deleted chapter, got %#v", chapterRepo.Infos)
	}
	pubCalls := eventBus.PubCalls()
	removedEvent, ok := pubCalls[0][0].(*event.ChapterRemovedEvent)
	if !ok || len(removedEvent.AssignedUserIDs) != 2 {
		t.Fatalf("unexpected remove event: %#v", pubCalls)
	}
	if len(msgRepo.Messages) != 1 {
		t.Fatalf("expected 1 delete_pending message, got %#v", msgRepo.Messages)
	}
}

func TestChapterAppInviteAssigneeUsesMockTxnRepos(t *testing.T) {
	chapterRepo := mock_repo.NewMockChapterRepo()
	chapterRepo.Infos["chapter-1"] = model.ChapterInfo{ID: "chapter-1", ComicID: "comic-1"}
	assignmentRepo := mock_repo.NewMockAssignmentRepo()
	assignmentRepo.Infos["assignment-1"] = *reviewerAssignment()
	comicRepo := mock_repo.NewMockComicRepo()
	worksetRepo := mock_repo.NewMockWorksetRepo()
	memberRepo := mock_repo.NewMockMemberRepo()
	chapterInvRepo := mock_repo.NewMockChapterInvitationRepo()
	txnMgr := mock_repo.NewMockTxnMgr(newMockTxnContext(mockTxnRepos{
		chapter:           chapterRepo,
		assignment:        assignmentRepo,
		chapterInvitation: chapterInvRepo,
	}))

	app := NewChapterApp(service.NewChapterService(), service.NewChapterInvitationService(), memberRepo, worksetRepo, comicRepo, chapterRepo, assignmentRepo, mock_repo.NewMockUserRepo(), mock_repo.NewMockPageRepo(), chapterInvRepo, txnMgr, mock_repo.NewMockOSSMessageRepo(), newMockEventBus(), newMockOSSClient())

	res, err := app.InviteAssignee(background(), "user-1", &val.InviteChapterAssigneeArgs{
		ChapterID: "chapter-1",
		InviteeQQ: "10001",
		Roles:     model.MaskRoles([]model.Role{model.RoleTranslator, model.RoleReviewer}),
	})
	requireNoErr(t, err)

	if res == nil || res.InvCode == "" {
		t.Fatalf("expected invitation code, got %#v", res)
	}

	if len(chapterInvRepo.Infos) != 1 {
		t.Fatalf("expected one invitation, got %#v", chapterInvRepo.Infos)
	}

	for _, info := range chapterInvRepo.Infos {
		if info.ChapterID != "chapter-1" || info.InviterID != "user-1" || info.InviteeQQ != "10001" {
			t.Fatalf("unexpected invitation info: %#v", info)
		}
		if !info.ToBeTranslator || !info.ToBeReviewer {
			t.Fatalf("unexpected invitation roles: %#v", info)
		}
	}
}

func TestChapterAppInviteAssigneeForbidden(t *testing.T) {
	chapterRepo := mock_repo.NewMockChapterRepo()
	chapterRepo.Infos["chapter-1"] = model.ChapterInfo{ID: "chapter-1", ComicID: "comic-1"}
	assignmentRepo := mock_repo.NewMockAssignmentRepo()
	assignmentRepo.Infos["assignment-1"] = *translatorAssignment()
	comicRepo := mock_repo.NewMockComicRepo()
	worksetRepo := mock_repo.NewMockWorksetRepo()
	memberRepo := mock_repo.NewMockMemberRepo()
	chapterInvRepo := mock_repo.NewMockChapterInvitationRepo()
	txnMgr := mock_repo.NewMockTxnMgr(newMockTxnContext(mockTxnRepos{
		chapter:           chapterRepo,
		assignment:        assignmentRepo,
		chapterInvitation: chapterInvRepo,
	}))

	app := NewChapterApp(service.NewChapterService(), service.NewChapterInvitationService(), memberRepo, worksetRepo, comicRepo, chapterRepo, assignmentRepo, mock_repo.NewMockUserRepo(), mock_repo.NewMockPageRepo(), chapterInvRepo, txnMgr, mock_repo.NewMockOSSMessageRepo(), newMockEventBus(), newMockOSSClient())

	_, err := app.InviteAssignee(background(), "user-1", &val.InviteChapterAssigneeArgs{
		ChapterID: "chapter-1",
		InviteeQQ: "10001",
		Roles:     model.MaskRoles([]model.Role{model.RoleTranslator}),
	})
	if err == nil {
		t.Fatal("expected forbidden invite error")
	}

	if len(chapterInvRepo.Infos) != 0 {
		t.Fatalf("expected no invitation on forbidden request, got %#v", chapterInvRepo.Infos)
	}
}

func TestChapterAppInviteAssigneeSupportsRedrawerRole(t *testing.T) {
	chapterRepo := mock_repo.NewMockChapterRepo()
	chapterRepo.Infos["chapter-1"] = model.ChapterInfo{ID: "chapter-1", ComicID: "comic-1"}
	assignmentRepo := mock_repo.NewMockAssignmentRepo()
	assignmentRepo.Infos["assignment-1"] = *reviewerAssignment()
	chapterInvRepo := mock_repo.NewMockChapterInvitationRepo()
	txnMgr := mock_repo.NewMockTxnMgr(newMockTxnContext(mockTxnRepos{
		chapter:           chapterRepo,
		assignment:        assignmentRepo,
		chapterInvitation: chapterInvRepo,
	}))

	app := NewChapterApp(service.NewChapterService(), service.NewChapterInvitationService(), mock_repo.NewMockMemberRepo(), mock_repo.NewMockWorksetRepo(), mock_repo.NewMockComicRepo(), chapterRepo, assignmentRepo, mock_repo.NewMockUserRepo(), mock_repo.NewMockPageRepo(), chapterInvRepo, txnMgr, mock_repo.NewMockOSSMessageRepo(), newMockEventBus(), newMockOSSClient())

	res, err := app.InviteAssignee(background(), "user-1", &val.InviteChapterAssigneeArgs{
		ChapterID: "chapter-1",
		InviteeQQ: "10001",
		Roles:     model.MaskRoles([]model.Role{model.RoleRedrawer}),
	})
	requireNoErr(t, err)

	if res == nil || res.InvCode == "" {
		t.Fatalf("expected invitation code, got %#v", res)
	}
	for _, info := range chapterInvRepo.Infos {
		if !info.ToBeRedrawer {
			t.Fatalf("expected redrawer invitation role: %#v", info)
		}
		if info.InvitedRoleMask() != model.RoleMask(model.RoleRedrawer) {
			t.Fatalf("unexpected redrawer invitation mask: %v", info.InvitedRoleMask())
		}
	}
}

func TestChapterAppListForbidden(t *testing.T) {
	comicRepo := mock_repo.NewMockComicRepo()
	comicRepo.Infos["comic-1"] = model.ComicInfo{ID: "comic-1", WorksetID: "workset-1"}
	worksetRepo := mock_repo.NewMockWorksetRepo()
	worksetRepo.Infos["workset-1"] = model.WorksetInfo{ID: "workset-1", TeamID: "team-1"}

	app := NewChapterApp(service.NewChapterService(), service.NewChapterInvitationService(), mock_repo.NewMockMemberRepo(), worksetRepo, comicRepo, mock_repo.NewMockChapterRepo(), mock_repo.NewMockAssignmentRepo(), mock_repo.NewMockUserRepo(), mock_repo.NewMockPageRepo(), mock_repo.NewMockChapterInvitationRepo(), mock_repo.NewMockTxnMgr(nil), mock_repo.NewMockOSSMessageRepo(), newMockEventBus(), newMockOSSClient())

	if _, err := app.List(background(), "user-1", &val.ListChapterArgs{ComicID: "comic-1"}); err == nil {
		t.Fatal("expected forbidden error")
	}
}

func TestChapterAppErrorPaths(t *testing.T) {
	t.Run("create rejects non admin", func(t *testing.T) {
		worksetRepo := mock_repo.NewMockWorksetRepo()
		worksetRepo.Infos["workset-1"] = model.WorksetInfo{ID: "workset-1", TeamID: "team-1"}
		comicRepo := mock_repo.NewMockComicRepo()
		comicRepo.Infos["comic-1"] = model.ComicInfo{ID: "comic-1", WorksetID: "workset-1"}
		memberRepo := mock_repo.NewMockMemberRepo()
		memberRepo.Infos["member-1"] = model.MemberInfo{ID: "member-1", UserID: "user-1", TeamID: "team-1"}

		app := NewChapterApp(service.NewChapterService(), service.NewChapterInvitationService(), memberRepo, worksetRepo, comicRepo, mock_repo.NewMockChapterRepo(), mock_repo.NewMockAssignmentRepo(), mock_repo.NewMockUserRepo(), mock_repo.NewMockPageRepo(), mock_repo.NewMockChapterInvitationRepo(), mock_repo.NewMockTxnMgr(nil), mock_repo.NewMockOSSMessageRepo(), newMockEventBus(), newMockOSSClient())
		if _, err := app.Create(background(), "user-1", &val.CreateChapterArgs{ComicID: "comic-1"}); err == nil {
			t.Fatal("expected forbidden create error")
		}
	})

	t.Run("update rejects missing chapter", func(t *testing.T) {
		app := NewChapterApp(service.NewChapterService(), service.NewChapterInvitationService(), mock_repo.NewMockMemberRepo(), mock_repo.NewMockWorksetRepo(), mock_repo.NewMockComicRepo(), mock_repo.NewMockChapterRepo(), mock_repo.NewMockAssignmentRepo(), mock_repo.NewMockUserRepo(), mock_repo.NewMockPageRepo(), mock_repo.NewMockChapterInvitationRepo(), mock_repo.NewMockTxnMgr(nil), mock_repo.NewMockOSSMessageRepo(), newMockEventBus(), newMockOSSClient())
		if err := app.Update(background(), "user-1", &val.UpdateChapterArgs{ChapterID: "missing"}); err == nil {
			t.Fatal("expected missing chapter error")
		}
	})

	t.Run("update rejects workflow transition without assignment", func(t *testing.T) {
		transition := model.WorkflowTranslateStart
		chapterRepo := mock_repo.NewMockChapterRepo()
		chapterRepo.Infos["chapter-1"] = model.ChapterInfo{ID: "chapter-1", ComicID: "comic-1"}
		app := NewChapterApp(service.NewChapterService(), service.NewChapterInvitationService(), mock_repo.NewMockMemberRepo(), mock_repo.NewMockWorksetRepo(), mock_repo.NewMockComicRepo(), chapterRepo, mock_repo.NewMockAssignmentRepo(), mock_repo.NewMockUserRepo(), mock_repo.NewMockPageRepo(), mock_repo.NewMockChapterInvitationRepo(), mock_repo.NewMockTxnMgr(nil), mock_repo.NewMockOSSMessageRepo(), newMockEventBus(), newMockOSSClient())
		if err := app.Update(background(), "user-1", &val.UpdateChapterArgs{ChapterID: "chapter-1", WorkflowTransition: &transition}); err == nil {
			t.Fatal("expected workflow permission error")
		}
	})

	t.Run("remove rejects non admin", func(t *testing.T) {
		chapterRepo := mock_repo.NewMockChapterRepo()
		chapterRepo.Infos["chapter-1"] = model.ChapterInfo{ID: "chapter-1", ComicID: "comic-1"}
		comicRepo := mock_repo.NewMockComicRepo()
		comicRepo.Infos["comic-1"] = model.ComicInfo{ID: "comic-1", WorksetID: "workset-1"}
		worksetRepo := mock_repo.NewMockWorksetRepo()
		worksetRepo.Infos["workset-1"] = model.WorksetInfo{ID: "workset-1", TeamID: "team-1"}
		memberRepo := mock_repo.NewMockMemberRepo()
		memberRepo.Infos["member-1"] = model.MemberInfo{ID: "member-1", UserID: "user-1", TeamID: "team-1"}

		app := NewChapterApp(service.NewChapterService(), service.NewChapterInvitationService(), memberRepo, worksetRepo, comicRepo, chapterRepo, mock_repo.NewMockAssignmentRepo(), mock_repo.NewMockUserRepo(), mock_repo.NewMockPageRepo(), mock_repo.NewMockChapterInvitationRepo(), mock_repo.NewMockTxnMgr(nil), mock_repo.NewMockOSSMessageRepo(), newMockEventBus(), newMockOSSClient())
		if err := app.Remove(background(), "user-1", "chapter-1"); err == nil {
			t.Fatal("expected forbidden remove error")
		}
	})

	t.Run("remove blocks db delete when oss enqueue fails", func(t *testing.T) {
		chapterRepo := mock_repo.NewMockChapterRepo()
		chapterRepo.Infos["chapter-1"] = model.ChapterInfo{ID: "chapter-1", ComicID: "comic-1"}
		comicRepo := mock_repo.NewMockComicRepo()
		comicRepo.Infos["comic-1"] = model.ComicInfo{ID: "comic-1", WorksetID: "workset-1"}
		worksetRepo := mock_repo.NewMockWorksetRepo()
		worksetRepo.Infos["workset-1"] = model.WorksetInfo{ID: "workset-1", TeamID: "team-1"}
		memberRepo := mock_repo.NewMockMemberRepo()
		memberRepo.Infos["member-1"] = *adminMember()
		assignmentRepo := mock_repo.NewMockAssignmentRepo()
		pageRepo := mock_repo.NewMockPageRepo()
		pageRepo.Infos["page-1"] = model.PageInfo{ID: "page-1", ChapterID: "chapter-1", OSSKey: "page-oss-1"}
		msgRepo := mock_repo.NewMockOSSMessageRepo()
		msgRepo.InsertErr = errors.New("boom")
		txnMgr := mock_repo.NewMockTxnMgr(newMockTxnContext(mockTxnRepos{chapter: chapterRepo, assignment: assignmentRepo, comic: comicRepo, user: mock_repo.NewMockUserRepo(), ossMessage: msgRepo}))

		app := NewChapterApp(service.NewChapterService(), service.NewChapterInvitationService(), memberRepo, worksetRepo, comicRepo, chapterRepo, assignmentRepo, mock_repo.NewMockUserRepo(), pageRepo, mock_repo.NewMockChapterInvitationRepo(), txnMgr, msgRepo, newMockEventBus(), newMockOSSClient())

		err := app.Remove(background(), "user-1", "chapter-1")
		if err == nil {
			t.Fatal("expected remove failure when oss enqueue fails")
		}

		if len(msgRepo.Messages) != 0 {
			t.Fatalf("expected no queued messages on failure, got %#v", msgRepo.Messages)
		}
	})
}
