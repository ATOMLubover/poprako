package event_infra

import (
	"context"
	"fmt"

	"poprako-s/internal/domain/model/aggr"
	"poprako-s/internal/domain/model/enum"
	"poprako-s/internal/domain/model/event"
	"poprako-s/internal/domain/model/query"
	repo_iface "poprako-s/internal/domain/repo"
	event_iface "poprako-s/internal/event"
	"poprako-s/pkg/util"

	"go.uber.org/zap"
)

// `NotifyChapterProgressHandler` sends workflow progress sys-mails to both
// next-phase assignees and reviewer assignees after a workflow-complete
// transition. It merges the two receiver groups and deduplicates by `RcvId`
// so that a user who holds both the next-phase role and `RoleReviewer`
// receives only one mail.
type NotifyChapterProgressHandler struct {
	chapterRepo    repo_iface.ChapterRepo
	assignmentRepo repo_iface.AssignmentRepo
	sysMailRepo    repo_iface.SysMailRepo
}

// `NewNotifyChapterProgressHandler` creates one `NotifyChapterProgressHandler`.
func NewNotifyChapterProgressHandler(
	chapterRepo repo_iface.ChapterRepo,
	assignmentRepo repo_iface.AssignmentRepo,
	sysMailRepo repo_iface.SysMailRepo,
) *NotifyChapterProgressHandler {
	return &NotifyChapterProgressHandler{
		chapterRepo:    chapterRepo,
		assignmentRepo: assignmentRepo,
		sysMailRepo:    sysMailRepo,
	}
}

// `EvTyp` returns the target event type of `NotifyChapterProgressHandler`.
func (h *NotifyChapterProgressHandler) EvTyp() event_iface.EvTyp {
	return event.EvChapterWorkflowCompleted
}

// `Handle` creates and sends progress sys-mails for next-phase and reviewer
// assignees, deduplicated by receiver id.
func (h *NotifyChapterProgressHandler) Handle(cx context.Context, ev event_iface.Event) {
	_ = cx

	// Parse payload and validate event type.
	payload, ok := ev.Payload().(*event.ChapterWorkflowCompletedEv)
	if !ok || payload == nil {
		zap.L().Error(
			"[NotifyChapterProgressHandler.Handle] invalid event payload",
			zap.Any("payload", ev.Payload()),
		)

		return
	}

	// Resolve next-phase role, reviewer label, and chinese workflow name from
	// completed transition. Returns false for unknown transitions.
	nextRole, reviewerLabel, workflowCn, ok := resolveProgressMailCfg(payload.CompletedTransition)
	if !ok {
		return
	}

	// Query chapter with include chain to fetch `comic`, `workset` and `team` info.
	chapter, err := h.chapterRepo.GetById(
		payload.ChapterId,
		enum.ChapterInclComicWorksetTeam,
	)
	if err != nil {
		zap.L().Error(
			"[NotifyChapterProgressHandler.Handle] failed to get chapter info",
			zap.String("chapter_id", payload.ChapterId),
			zap.Error(err),
		)

		return
	}

	if chapter.Comic == nil || chapter.Comic.Workset == nil || chapter.Comic.Workset.Team == nil {
		zap.L().Error(
			"[NotifyChapterProgressHandler.Handle] missing include chain data",
			zap.String("chapter_id", payload.ChapterId),
		)

		return
	}

	// Query all assignments under the chapter.
	assignmentOpt := &query.ListAssignmentOpt{ChapterId: &chapter.Id}
	assignments, err := h.assignmentRepo.List(assignmentOpt)
	if err != nil {
		zap.L().Error(
			"[NotifyChapterProgressHandler.Handle] failed to list chapter assignments",
			zap.String("chapter_id", payload.ChapterId),
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

	// Build mail content templates for next-phase and reviewer groups.
	// Both use the same `workflowCn` label since the labels are identical.
	mailBody := fmt.Sprintf(
		"「%s」-「%s」漫画 #%d『%s』章节 #%d 「%s」已完成。",
		chapter.Comic.Workset.Team.Name,
		chapter.Comic.Workset.Name,
		chapter.Comic.Index+1,
		shortTitle,
		chapter.Index+1,
		workflowCn,
	)

	// Collect mails with deduplication by `RcvId`.
	// Next-phase assignees are added first; reviewer assignees are skipped
	// when the same user already appears as a next-phase receiver.
	seen := make(map[string]bool, len(assignments))
	mailCres := make([]*aggr.SysMailCre, 0, len(assignments))

	// Collect next-phase assignees.
	for i := range assignments {
		if assignments[i] == nil || !assignments[i].HasAnyRole(nextRole) {
			continue
		}

		seen[assignments[i].UserId] = true

		mailCres = append(mailCres, &aggr.SysMailCre{
			Id:      util.GenId("sys_mail"),
			RcvId:   assignments[i].UserId,
			Title:   mailTitle,
			Content: mailBody,
		})
	}

	// Collect reviewer assignees, skipping already-notified users.
	// `reviewerLabel` is empty for `WorkflowTypesetComplete` because reviewers
	// are already covered by the next-phase collection above.
	if reviewerLabel != "" {
		for i := range assignments {
			if assignments[i] == nil || !assignments[i].HasAnyRole(enum.RoleReviewer) {
				continue
			}

			if seen[assignments[i].UserId] {
				continue
			}

			seen[assignments[i].UserId] = true

			mailCres = append(mailCres, &aggr.SysMailCre{
				Id:      util.GenId("sys_mail"),
				RcvId:   assignments[i].UserId,
				Title:   mailTitle,
				Content: mailBody,
			})
		}
	}

	// Early return when no receivers are found.
	if len(mailCres) == 0 {
		return
	}

	if err := h.sysMailRepo.SendBatch(mailCres); err != nil {
		zap.L().Error(
			"[NotifyChapterProgressHandler.Handle] failed to send chapter progress sys-mails",
			zap.String("chapter_id", payload.ChapterId),
			zap.String("workflow_transition", string(payload.CompletedTransition)),
			zap.Error(err),
		)

		return
	}
}

// `resolveProgressMailCfg` maps a completed transition to the next-phase role,
// reviewer label, and chinese workflow name used in notification mails.
// `reviewerLabel` is empty when the transition's next-phase is already
// `RoleReviewer` (i.e. `WorkflowTypesetComplete`), avoiding redundant collection.
func resolveProgressMailCfg(t enum.WorkflowTransition) (enum.Role, string, string, bool) {
	switch t {
	case enum.WorkflowUploadComplete:
		return enum.RoleTranslator, "上传", "上传", true

	case enum.WorkflowTranslateComplete:
		return enum.RoleProofreader, "翻译", "翻译", true

	case enum.WorkflowProofreadComplete:
		return enum.RoleTypesetter, "校对", "校对", true

	case enum.WorkflowTypesetComplete:
		// Next-phase is already `RoleReviewer`; no separate reviewer mail needed.
		return enum.RoleReviewer, "", "嵌字", true

	case enum.WorkflowReviewComplete:
		return enum.RolePublisher, "监修", "监修", true

	default:
		return 0, "", "", false
	}
}

// `truncateTitle` truncates title to at most `maxRune` runes and appends `...` when truncated.
func truncateTitle(title string, maxRune int) string {
	titleRunes := []rune(title)
	if len(titleRunes) <= maxRune {
		return title
	}

	return string(titleRunes[:maxRune]) + "..."
}
