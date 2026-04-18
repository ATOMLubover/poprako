package event

const (
	EventTypeComicCreated EventType = "ComicCreated"
	EventTypeComicRemoved EventType = "ComicRemoved"
)

// ComicCreatedEvent 代表漫画创建后的同步统计事件
type ComicCreatedEvent struct {
	WorksetID string
}

func (e *ComicCreatedEvent) EventType() EventType {
	return EventTypeComicCreated
}

func (e *ComicCreatedEvent) Payload() any {
	return e
}

// ComicRemovedEvent 代表漫画删除后的同步统计事件
type ComicRemovedEvent struct {
	WorksetID string
}

func (e *ComicRemovedEvent) EventType() EventType {
	return EventTypeComicRemoved
}

func (e *ComicRemovedEvent) Payload() any {
	return e
}
