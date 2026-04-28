package event

import event_iface "poprako-s/internal/event"

// `ChapterPublishedEv` represents chapter-published domain event payload.
type ChapterPublishedEv struct {
	ChapterId       string
	AssignedUserIds []string
}

// `NewChapterPublishedEv` creates chapter-published event.
func NewChapterPublishedEv(chapterId string, assignedUserIds []string) event_iface.Event {
	return &ChapterPublishedEv{ChapterId: chapterId, AssignedUserIds: assignedUserIds}
}

// `EvTyp` returns event type identifier.
func (e *ChapterPublishedEv) EvTyp() event_iface.EvTyp {
	return EvChapterPublished
}

// `Payload` returns typed payload.
func (e *ChapterPublishedEv) Payload() any {
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
