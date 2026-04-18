package event

const (
	EventTypeAssignmentCreated EventType = "AssignmentCreated"
	EventTypeAssignmentRemoved EventType = "AssignmentRemoved"
)

// AssignmentCreatedEvent 代表创建 assignment 后的同步统计事件
type AssignmentCreatedEvent struct {
	UserID    string
	ChapterID string
}

func (e *AssignmentCreatedEvent) EventType() EventType {
	return EventTypeAssignmentCreated
}

func (e *AssignmentCreatedEvent) Payload() any {
	return e
}

// AssignmentRemovedEvent 代表删除 assignment 后的同步统计事件
type AssignmentRemovedEvent struct {
	UserID       string
	WasPublished bool
}

func (e *AssignmentRemovedEvent) EventType() EventType {
	return EventTypeAssignmentRemoved
}

func (e *AssignmentRemovedEvent) Payload() any {
	return e
}
