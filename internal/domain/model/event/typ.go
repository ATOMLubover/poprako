package event

import event_iface "poprako-s/internal/event"

// Event type identifiers for user-domain events
const (
	// `EvUserReg` is emitted after a new user registration completes
	EvUserReg event_iface.EvTyp = "event:user_reg"
	// `EvAssignmentCreated` is emitted after one assignment is created
	EvAssignmentCreated event_iface.EvTyp = "event:assignment_created"
	// `EvAssignmentRemoved` is emitted after one assignment is removed
	EvAssignmentRemoved event_iface.EvTyp = "event:assignment_removed"
	// `EvChapterPublished` is emitted after one chapter reaches publish-complete state
	EvChapterPublished event_iface.EvTyp = "event:chapter_published"
	// `EvChapterRemoved` is emitted after one chapter is removed
	EvChapterRemoved event_iface.EvTyp = "event:chapter_removed"
)
