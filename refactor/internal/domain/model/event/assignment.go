package event

import event_iface "poprako-s/internal/event"

// `AssignmentCreatedEv` represents assignment-created domain event payload.
type AssignmentCreatedEv struct {
	UserId    string
	ChapterId string
}

// `NewAssignmentCreatedEv` creates assignment-created event.
func NewAssignmentCreatedEv(userId string, chapterId string) event_iface.Event {
	return &AssignmentCreatedEv{UserId: userId, ChapterId: chapterId}
}

// `EvTyp` returns event type identifier.
func (e *AssignmentCreatedEv) EvTyp() event_iface.EvTyp {
	return EvAssignmentCreated
}

// `Payload` returns typed payload.
func (e *AssignmentCreatedEv) Payload() any {
	return e
}

// `AssignmentRemovedEv` represents assignment-removed domain event payload.
type AssignmentRemovedEv struct {
	UserId       string
	ChapterId    string
	WasPublished bool
}

// `NewAssignmentRemovedEv` creates assignment-removed event.
func NewAssignmentRemovedEv(userId string, chapterId string, wasPublished bool) event_iface.Event {
	return &AssignmentRemovedEv{
		UserId:       userId,
		ChapterId:    chapterId,
		WasPublished: wasPublished,
	}
}

// `EvTyp` returns event type identifier.
func (e *AssignmentRemovedEv) EvTyp() event_iface.EvTyp {
	return EvAssignmentRemoved
}

// `Payload` returns typed payload.
func (e *AssignmentRemovedEv) Payload() any {
	return e
}
