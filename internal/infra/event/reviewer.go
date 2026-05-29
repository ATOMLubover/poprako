package event_infra

import (
	"context"
	"fmt"

	"poprako-s/internal/domain/model/aggr"
	"poprako-s/internal/domain/model/enum"
	event_impl "poprako-s/internal/domain/model/event"
	"poprako-s/internal/domain/model/query"
	repo_iface "poprako-s/internal/domain/repo"
	event_iface "poprako-s/internal/event"
	"poprako-s/pkg/util"

	"go.uber.org/zap"
)

// `NotifyReviewerOnProgressHandler` sends workflow progress sys-mail to all
// reviewer assignees of a chapter after each workflow-complete transition,
// except `WorkflowTypesetComplete` which is already covered by `NotifyNextPhaseHandler`.
type NotifyReviewerOnProgressHandler struct {
	chapterRepo    repo_iface.ChapterRepo
	assignmentRepo repo_iface.AssignmentRepo
	sysMailRepo    repo_iface.SysMailRepo
}

// `NewNotifyReviewerOnProgressHandler` creates one `NotifyReviewerOnProgressHandler`.
func NewNotifyReviewerOnProgressHandler(
	chapterRepo repo_iface.ChapterRepo,
	assignmentRepo repo_iface.AssignmentRepo,
	sysMailRepo repo_iface.SysMailRepo,
) *NotifyReviewerOnProgressHandler {
	return &NotifyReviewerOnProgressHandler{
		chapterRepo:    chapterRepo,
		assignmentRepo: assignmentRepo,
		sysMailRepo:    sysMailRepo,
	}
}

// `EvTyp` returns the target event type of `NotifyReviewerOnProgressHandler`.
func (h *NotifyReviewerOnProgressHandler) EvTyp() event_iface.EvTyp {
	return event_impl.EvChapterWorkflowCompleted
}

// `Handle` sends workflow progress sys-mails to all reviewer assignees.
// Skips `WorkflowTypesetComplete` to avoid duplicating the notification
// already sent by `NotifyNextPhaseHandler` for that transition.
func (h *NotifyReviewerOnProgressHandler) Handle(cx context.Context, ev event_iface.Event) {
	_ = cx

	// Parse payload and validate event type.
	payload, ok := ev.Payload().(*event_impl.ChapterWorkflowCompletedEv)
	if !ok || payload == nil {
		zap.L().Error(
			"[NotifyReviewerOnProgressHandler.Handle] invalid event payload",
			zap.Any("payload", ev.Payload()),
		)

		return
	}

	// Resolve chinese workflow label; returns false for transitions to skip.
	label, ok := resolveReviewerProgressLabel(payload.CompletedTransition)
	if !ok {
		return
	}

	// Delegate to shared helper.
	sendReviewerMails(
		"[NotifyReviewerOnProgressHandler.Handle]",
		payload.ChapterId,
		label,
		h.chapterRepo,
		h.assignmentRepo,
		h.sysMailRepo,
	)
}

// `NotifyReviewerOnPublishHandler` sends a publish sys-mail to all reviewer
// assignees of a chapter when a `ChapterPublishedEv` fires.
type NotifyReviewerOnPublishHandler struct {
	chapterRepo    repo_iface.ChapterRepo
	assignmentRepo repo_iface.AssignmentRepo
	sysMailRepo    repo_iface.SysMailRepo
}

// `NewNotifyReviewerOnPublishHandler` creates one `NotifyReviewerOnPublishHandler`.
func NewNotifyReviewerOnPublishHandler(
	chapterRepo repo_iface.ChapterRepo,
	assignmentRepo repo_iface.AssignmentRepo,
	sysMailRepo repo_iface.SysMailRepo,
) *NotifyReviewerOnPublishHandler {
	return &NotifyReviewerOnPublishHandler{
		chapterRepo:    chapterRepo,
		assignmentRepo: assignmentRepo,
		sysMailRepo:    sysMailRepo,
	}
}

// `EvTyp` returns the target event type of `NotifyReviewerOnPublishHandler`.
func (h *NotifyReviewerOnPublishHandler) EvTyp() event_iface.EvTyp {
	return event_impl.EvChapterPublished
}

// `Handle` sends publish sys-mails to all reviewer assignees of the chapter.
func (h *NotifyReviewerOnPublishHandler) Handle(cx context.Context, ev event_iface.Event) {
	_ = cx

	// Parse payload and validate event type.
	payload, ok := ev.Payload().(*event_impl.ChapterPublishedEv)
	if !ok || payload == nil {
		zap.L().Error(
			"[NotifyReviewerOnPublishHandler.Handle] invalid event payload",
			zap.Any("payload", ev.Payload()),
		)

		return
	}

	// Delegate to shared helper.
	sendReviewerMails(
		"[NotifyReviewerOnPublishHandler.Handle]",
		payload.ChapterId,
		"发布",
		h.chapterRepo,
		h.assignmentRepo,
		h.sysMailRepo,
	)
}

// `resolveReviewerProgressLabel` maps a completed workflow transition to a
// chinese label used in the reviewer notification mail. Returns `("", false)`
// for `WorkflowTypesetComplete` (already handled by `NotifyNextPhaseHandler`)
// and for any unknown transition.
func resolveReviewerProgressLabel(t enum.WorkflowTransition) (string, bool) {
	switch t {
	case enum.WorkflowUploadComplete:
		return "上传", true

	case enum.WorkflowTranslateComplete:
		return "翻译", true

	case enum.WorkflowProofreadComplete:
		return "校对", true

	case enum.WorkflowTypesetComplete:
		// `NotifyNextPhaseHandler` already notifies reviewers for typeset-complete;
		// skip here to avoid sending a duplicate mail.
		return "", false

	case enum.WorkflowReviewComplete:
		return "监修", true

	default:
		return "", false
	}
}

// `sendReviewerMails` queries the chapter and its `RoleReviewer` assignees,
// then sends one sys-mail per reviewer. `caller` is the bracketed caller label
// prepended to every error log line for easy triage.
func sendReviewerMails(
	caller string,
	chapterId string,
	workflowLabel string,
	chapterRepo repo_iface.ChapterRepo,
	assignmentRepo repo_iface.AssignmentRepo,
	sysMailRepo repo_iface.SysMailRepo,
) {
	// Query chapter with include chain to fetch `comic`, `workset` and `team` info.
	chapter, err := chapterRepo.GetById(chapterId, enum.ChapterInclComicWorksetTeam)
	if err != nil {
		zap.L().Error(
			caller+" failed to get chapter info",
			zap.String("chapter_id", chapterId),
			zap.Error(err),
		)

		return
	}

	if chapter.Comic == nil || chapter.Comic.Workset == nil || chapter.Comic.Workset.Team == nil {
		zap.L().Error(
			caller+" missing include chain data",
			zap.String("chapter_id", chapterId),
		)

		return
	}

	// Query all assignments under the chapter to locate reviewer receivers.
	assignmentOpt := &query.ListAssignmentOpt{ChapterId: &chapter.Id}
	assignments, err := assignmentRepo.List(assignmentOpt)
	if err != nil {
		zap.L().Error(
			caller+" failed to list chapter assignments",
			zap.String("chapter_id", chapterId),
			zap.Error(err),
		)

		return
	}

	shortTitle := truncateTitle(chapter.Comic.Title, 15)
	mailTitle := fmt.Sprintf(
		"你参加的漫画『%s』#%d 章节有进度更新",
		shortTitle,
		chapter.Index+1,
	)
	mailContent := fmt.Sprintf(
		"「%s」-「%s」漫画 #%d『%s』章节 #%d 「%s」已完成。",
		chapter.Comic.Workset.Team.Name,
		chapter.Comic.Workset.Name,
		chapter.Comic.Index+1,
		shortTitle,
		chapter.Index+1,
		workflowLabel,
	)

	// Collect one mail per reviewer assignee.
	mailCres := make([]*aggr.SysMailCre, 0, len(assignments))
	for i := range assignments {
		if assignments[i] == nil || !assignments[i].HasAnyRole(enum.RoleReviewer) {
			continue
		}

		mailCres = append(mailCres, &aggr.SysMailCre{
			Id:      util.GenId("sys_mail"),
			RcvId:   assignments[i].UserId,
			Title:   mailTitle,
			Content: mailContent,
		})
	}

	// Early return when no reviewer is assigned to avoid an empty batch insert.
	if len(mailCres) == 0 {
		return
	}

	if err := sysMailRepo.SendBatch(mailCres); err != nil {
		zap.L().Error(
			caller+" failed to send reviewer sys-mails",
			zap.String("chapter_id", chapterId),
			zap.Error(err),
		)

		return
	}
}
