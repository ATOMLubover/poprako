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

// `NotifyNextPhaseHandler` sends chapter progress sys-mail to next-phase assignees.
type NotifyNextPhaseHandler struct {
	chapterRepo    repo_iface.ChapterRepo
	assignmentRepo repo_iface.AssignmentRepo
	sysMailRepo    repo_iface.SysMailRepo
}

// `NewNotifyNextPhaseHandler` creates one `NotifyNextPhaseHandler`.
func NewNotifyNextPhaseHandler(
	chapterRepo repo_iface.ChapterRepo,
	assignmentRepo repo_iface.AssignmentRepo,
	sysMailRepo repo_iface.SysMailRepo,
) *NotifyNextPhaseHandler {
	return &NotifyNextPhaseHandler{
		chapterRepo:    chapterRepo,
		assignmentRepo: assignmentRepo,
		sysMailRepo:    sysMailRepo,
	}
}

// `EvTyp` returns the target event type of `NotifyNextPhaseHandler`.
func (h *NotifyNextPhaseHandler) EvTyp() event_iface.EvTyp {
	return event.EvChapterWorkflowCompleted
}

// `Handle` creates and sends progress sys-mails for next-phase assignees.
func (h *NotifyNextPhaseHandler) Handle(cx context.Context, ev event_iface.Event) {
	_ = cx

	// Parse payload and validate event type.
	payload, ok := ev.Payload().(*event.ChapterWorkflowCompletedEv)
	if !ok || payload == nil {
		zap.L().Error(
			"[NotifyNextPhaseHandler.Handle] invalid event payload for NotifyNextPhaseHandler",
			zap.Any("payload", ev.Payload()),
		)

		return
	}

	// Resolve next phase role and chinese workflow name from completed transition.
	nextRole, workflowCn, ok := resolveWorkflowMailCfg(payload.CompletedTransition)
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
			"[NotifyNextPhaseHandler.Handle] failed to get chapter info",
			zap.String("chapter_id", payload.ChapterId),
			zap.Error(err),
		)

		return
	}

	if chapter.Comic == nil || chapter.Comic.Workset == nil || chapter.Comic.Workset.Team == nil {
		zap.L().Error(
			"[NotifyNextPhaseHandler.Handle] missing include chain data",
			zap.String("chapter_id", payload.ChapterId),
		)

		return
	}

	// Query all assignments under chapter and pick next-phase receivers.
	assignmentOpt := &query.ListAssignmentOpt{ChapterId: &chapter.Id}
	assignments, err := h.assignmentRepo.List(assignmentOpt)
	if err != nil {
		zap.L().Error(
			"[NotifyNextPhaseHandler.Handle] failed to list chapter assignments",
			zap.String("chapter_id", payload.ChapterId),
			zap.Error(err),
		)

		return
	}

	shortTitle := truncateTitle(chapter.Comic.Title, 15)
	mailTitle := fmt.Sprintf(
		"你参加的漫画『%s』#%d 章节有进度更新",
		shortTitle,
		chapter.Index,
	)
	mailContent := fmt.Sprintf(
		"「%s」-「%s」漫画 %d『%s』#%d 章节「%s」已完成。",
		chapter.Comic.Workset.Team.Name,
		chapter.Comic.Workset.Name,
		chapter.Comic.Index,
		shortTitle,
		chapter.Index,
		workflowCn,
	)

	mailCres := make([]*aggr.SysMailCre, 0, len(assignments))
	for i := range assignments {
		if assignments[i] == nil || !assignments[i].HasAnyRole(nextRole) {
			continue
		}

		mailCres = append(mailCres, &aggr.SysMailCre{
			Id:      util.GenId("sys_mail"),
			RcvId:   assignments[i].UserId,
			Title:   mailTitle,
			Content: mailContent,
		})
	}

	if err := h.sysMailRepo.SendBatch(mailCres); err != nil {
		zap.L().Error(
			"[NotifyNextPhaseHandler.Handle] failed to create chapter workflow progress sys-mails",
			zap.String("chapter_id", payload.ChapterId),
			zap.String("workflow_transition", string(payload.CompletedTransition)),
			zap.Error(err),
		)

		return
	}
}

// `resolveWorkflowMailCfg` maps completed transition to next role and chinese workflow name.
func resolveWorkflowMailCfg(t enum.WorkflowTransition) (enum.Role, string, bool) {
	switch t {
	case enum.WorkflowUploadComplete:
		return enum.RoleTranslator, "上传", true

	case enum.WorkflowTranslateComplete:
		return enum.RoleProofreader, "翻译", true

	case enum.WorkflowProofreadComplete:
		return enum.RoleTypesetter, "校对", true

	case enum.WorkflowTypesetComplete:
		return enum.RoleReviewer, "嵌字", true

	case enum.WorkflowReviewComplete:
		return enum.RolePublisher, "监修", true

	default:
		return 0, "", false
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
