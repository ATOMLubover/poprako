package event

import (
	"poprako-s/internal/domain/model/enum"
	event_iface "poprako-s/internal/event"
)

// `ChapterPublishedEv` represents chapter-published domain event payload.
type ChapterPublishedEv struct {
	ChapterId string
}

// `NewChapterPublishedEv` creates chapter-published event.
func NewChapterPublishedEv(chapterId string) event_iface.Event {
	return &ChapterPublishedEv{ChapterId: chapterId}
}

// `EvTyp` returns event type identifier.
func (e *ChapterPublishedEv) EvTyp() event_iface.EvTyp {
	return EvChapterPublished
}

// `Payload` returns typed payload.
func (e *ChapterPublishedEv) Payload() any {
	return e
}

// `ChapterWorkflowCompletedEv` represents chapter workflow-completed domain event payload.
type ChapterWorkflowCompletedEv struct {
	ChapterId           string
	CompletedTransition enum.WorkflowTransition
}

// `NewChapterWorkflowCompletedEv` creates chapter workflow-completed event.
func NewChapterWorkflowCompletedEv(chapterId string, completedTransition enum.WorkflowTransition) event_iface.Event {
	return &ChapterWorkflowCompletedEv{
		ChapterId:           chapterId,
		CompletedTransition: completedTransition,
	}
}

// `EvTyp` returns event type identifier.
func (e *ChapterWorkflowCompletedEv) EvTyp() event_iface.EvTyp {
	return EvChapterWorkflowCompleted
}

// `Payload` returns typed payload.
func (e *ChapterWorkflowCompletedEv) Payload() any {
	return e
}

// `ChapterWorkflowRevertedEv` represents chapter workflow-reverted domain event payload
type ChapterWorkflowRevertedEv struct {
	ChapterId          string
	RevertedTransition enum.WorkflowTransition
}

// `NewChapterWorkflowRevertedEv` creates chapter workflow-reverted event
func NewChapterWorkflowRevertedEv(chapterId string, revertedTransition enum.WorkflowTransition) event_iface.Event {
	return &ChapterWorkflowRevertedEv{
		ChapterId:          chapterId,
		RevertedTransition: revertedTransition,
	}
}

// `EvTyp` returns event type identifier
func (e *ChapterWorkflowRevertedEv) EvTyp() event_iface.EvTyp {
	return EvChapterWorkflowReverted
}

// `Payload` returns typed payload
func (e *ChapterWorkflowRevertedEv) Payload() any {
	return e
}

// `ChapterRemovedEv` represents chapter-removed domain event payload.
type ChapterRemovedEv struct {
	ChapterId       string
	WasPublished    bool
	AssignedUserIds []string
}

// `NewChapterRemovedEv` creates chapter-removed event.
func NewChapterRemovedEv(chapterId string, wasPublished bool, assignedUserIds []string) event_iface.Event {
	return &ChapterRemovedEv{ChapterId: chapterId, WasPublished: wasPublished, AssignedUserIds: assignedUserIds}
}

// `EvTyp` returns event type identifier.
func (e *ChapterRemovedEv) EvTyp() event_iface.EvTyp {
	return EvChapterRemoved
}

// `Payload` returns typed payload.
func (e *ChapterRemovedEv) Payload() any {
	return e
}
