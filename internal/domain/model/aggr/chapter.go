package aggr

import (
	"errors"
	"fmt"
	"time"

	"poprako-s/internal/domain/model/enum"
	"poprako-s/internal/domain/model/event"
	event_iface "poprako-s/internal/event"
)

// `Chapter` represents one chapter inside a comic
// It contains workflow timestamps and pinned state
// Mutable workflow fields are advanced by `TransiteWorkflow`
// Relation fields are include-driven
// Time pointers keep nullable database semantics
type Chapter struct {
	// Embedded for domain events:
	// - `ChapterPublishedEv`: one chapter reaches publish-complete state
	// - `ChapterWorkflowCompletedEv`: one workflow-complete transition with next phase succeeds
	event_iface.EvBase

	// `Id` is chapter identifier
	Id string

	// `ComicId` is owner comic identifier
	ComicId string
	// `Comic` is included only when includes contain `comic`
	Comic *Comic

	// `IsPinned` indicates whether this chapter is pinned in comic
	IsPinned bool

	// `Index` is comic-scoped chapter order
	Index int
	// `Subtitle` is chapter subtitle
	Subtitle string

	// `PageCount` is chapter page count
	PageCount int
	// `TotalUnitCount` is total translation unit count
	TotalUnitCount int
	// `TranslatedUnitCount` is translated unit count
	TranslatedUnitCount int
	// `ProofreadUnitCount` is proofread unit count
	ProofreadUnitCount int

	// `CreatorId` is chapter creator user identifier
	CreatorId string
	// `Creator` is included only when includes contain `creator`.
	Creator *User

	// `CreatedAt` is creation timestamp
	CreatedAt time.Time
	// `UpdatedAt` is update timestamp
	UpdatedAt time.Time

	// Workflow timestamps.
	UploadedAt *time.Time

	TranslatingAt *time.Time
	TranslatedAt  *time.Time

	ProofreadingAt *time.Time
	ProofreadAt    *time.Time

	TypesettingAt *time.Time
	TypesetAt     *time.Time

	ReviewedAt  *time.Time
	PublishedAt *time.Time
}

// `TransiteWorkflow` applies a workflow transition to chapter state.
func (c *Chapter) TransiteWorkflow(t enum.WorkflowTransition) error {
	now := time.Now()

	markStarted := func(
		startedAt **time.Time,
		completedAt *time.Time,
		ongoingMsg string,
		completedMsg string,
	) error {
		if completedAt != nil {
			return errors.New(completedMsg)
		}

		if *startedAt != nil {
			return errors.New(ongoingMsg)
		}

		*startedAt = &now

		return nil
	}

	markCompleted := func(
		startedAt *time.Time,
		completedAt **time.Time,
		pendingMsg string,
		completedMsg string,
	) error {
		if *completedAt != nil {
			return errors.New(completedMsg)
		}

		if startedAt == nil {
			return errors.New(pendingMsg)
		}

		*completedAt = &now

		return nil
	}

	markOnce := func(completedAt **time.Time, completedMsg string) error {
		if *completedAt != nil {
			return errors.New(completedMsg)
		}

		*completedAt = &now

		return nil
	}

	switch t {
	case enum.WorkflowUploadComplete:
		if err := markOnce(&c.UploadedAt, "上传进度已经标记为完成"); err != nil {
			return err
		}

		c.PushEv(event.NewChapterWorkflowCompletedEv(c.Id, t))

	case enum.WorkflowTranslateStart:
		if err := markStarted(
			&c.TranslatingAt,
			c.TranslatedAt,
			"翻译进度已经处于进行中",
			"翻译进度已经标记为完成",
		); err != nil {
			return err
		}

	case enum.WorkflowTranslateComplete:
		if err := markCompleted(
			c.TranslatingAt,
			&c.TranslatedAt,
			"翻译进度尚未开始",
			"翻译进度已经标记为完成",
		); err != nil {
			return err
		}

		c.PushEv(event.NewChapterWorkflowCompletedEv(c.Id, t))

	case enum.WorkflowProofreadStart:
		if err := markStarted(
			&c.ProofreadingAt,
			c.ProofreadAt,
			"校对进度已经处于进行中",
			"校对进度已经标记为完成",
		); err != nil {
			return err
		}

	case enum.WorkflowProofreadComplete:
		if err := markCompleted(
			c.ProofreadingAt,
			&c.ProofreadAt,
			"校对进度尚未开始",
			"校对进度已经标记为完成",
		); err != nil {
			return err
		}

		c.PushEv(event.NewChapterWorkflowCompletedEv(c.Id, t))

	case enum.WorkflowTypesetStart:
		if err := markStarted(
			&c.TypesettingAt,
			c.TypesetAt,
			"嵌字进度已经处于进行中",
			"嵌字进度已经标记为完成",
		); err != nil {
			return err
		}

	case enum.WorkflowTypesetComplete:
		if err := markCompleted(
			c.TypesettingAt,
			&c.TypesetAt,
			"嵌字进度尚未开始",
			"嵌字进度已经标记为完成",
		); err != nil {
			return err
		}

		c.PushEv(event.NewChapterWorkflowCompletedEv(c.Id, t))

	case enum.WorkflowReviewComplete:
		if err := markOnce(&c.ReviewedAt, "监修进度已经标记为完成"); err != nil {
			return err
		}

		c.PushEv(event.NewChapterWorkflowCompletedEv(c.Id, t))

	case enum.WorkflowPublishComplete:
		if err := markOnce(&c.PublishedAt, "发布进度已经标记为完成"); err != nil {
			return err
		}

		c.PushEv(event.NewChapterPublishedEv(c.Id))

	default:
		return fmt.Errorf("unknown workflow transition: %s", t)
	}

	return nil
}

// `RevertWorkflow` clears a single workflow timestamp back to NULL
// Publish-complete has no revert — call with a publish revert constant returns an error
// If the target timestamp is already NULL, the method is a no-op and no event is emitted
// On a successful clear, `ChapterWorkflowRevertedEv` is pushed to the aggregate
func (c *Chapter) RevertWorkflow(t enum.WorkflowTransition) error {
	// Helper: clear a timestamp field; returns true when an actual change was made
	clear := func(field **time.Time) bool {
		if *field == nil {
			return false
		}

		*field = nil

		return true
	}

	var changed bool

	switch t {
	case enum.WorkflowUploadRevert:
		changed = clear(&c.UploadedAt)
	case enum.WorkflowTranslateStartRevert:
		changed = clear(&c.TranslatingAt)
	case enum.WorkflowTranslateRevert:
		changed = clear(&c.TranslatedAt)
	case enum.WorkflowProofreadStartRevert:
		changed = clear(&c.ProofreadingAt)
	case enum.WorkflowProofreadRevert:
		changed = clear(&c.ProofreadAt)
	case enum.WorkflowTypesetStartRevert:
		changed = clear(&c.TypesettingAt)
	case enum.WorkflowTypesetRevert:
		changed = clear(&c.TypesetAt)
	case enum.WorkflowReviewRevert:
		changed = clear(&c.ReviewedAt)
	default:
		return fmt.Errorf("unknown or disallowed revert transition: %s", t)
	}

	// Emit event only when an actual change was made
	if changed {
		c.PushEv(event.NewChapterWorkflowRevertedEv(c.Id, t))
	}

	return nil
}

// `ChapterCre` holds creation payload for chapter insert.
type ChapterCre struct {
	// `Id` is generated by chapter service
	Id string

	// `ComicId` is owner comic identifier
	ComicId string
	// `Index` is comic-scoped chapter index
	Index int
	// `Subtitle` is optional subtitle and defaults to `Ch.<index>` when nil
	Subtitle *string

	// `CreatorId` is chapter creator user identifier
	CreatorId string
}

// `ChapterUpd` holds mutable fields for chapter update.
// Optional fields are patch-like and applied only when selected.
// Workflow fields use pointer-to-pointer to distinguish skip (nil) from
// an explicit nullable write (&nil means write NULL; &ptr means write *ptr).
// `WorkflowTransition` and `RevertTransition` are mutually exclusive; set at most one
// `Id` identifies the target chapter
type ChapterUpd struct {
	// `Id` identifies target chapter.
	Id string

	// `Subtitle` is optional subtitle update.
	Subtitle *string
	// `IsPinned` optionally updates pinned status.
	IsPinned *bool
	// `WorkflowTransition` identifies the forward workflow transition driving this update
	WorkflowTransition *enum.WorkflowTransition
	// `RevertTransition` identifies the revert transition driving this update
	// Mutually exclusive with `WorkflowTransition`
	RevertTransition *enum.WorkflowTransition

	// `UploadedAt` is the upload completion timestamp update.
	UploadedAt **time.Time
	// `TranslatingAt` is the translate-start timestamp update.
	TranslatingAt **time.Time
	// `TranslatedAt` is the translate-complete timestamp update.
	TranslatedAt **time.Time
	// `ProofreadingAt` is the proofread-start timestamp update.
	ProofreadingAt **time.Time
	// `ProofreadAt` is the proofread-complete timestamp update.
	ProofreadAt **time.Time
	// `TypesettingAt` is the typeset-start timestamp update.
	TypesettingAt **time.Time
	// `TypesetAt` is the typeset-complete timestamp update.
	TypesetAt **time.Time
	// `ReviewedAt` is the review-complete timestamp update.
	ReviewedAt **time.Time
	// `PublishedAt` is the publish-complete timestamp update.
	PublishedAt **time.Time
}
