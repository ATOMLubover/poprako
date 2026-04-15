package event

// EventType 用于类型安全地标识事件类型
type EventType string

// PubType 用于标识一个事件是需要被同步处理还是异步处理
type PubType int

const (
	PubTypeSync PubType = iota
	PubTypeAsync
)

// EventBase 是带有 events 的类型的基础结构体
// 它提供统一的事件操作方法，简化 model 的实现
type EventBase struct {
	events []Event
}

func (b *EventBase) PushEvent(e Event) {
	b.events = append(b.events, e)
}

func (b *EventBase) ClearEvents() {
	b.events = nil
}

// 实现 EventSource 接口
func (b *EventBase) PullEvents() []Event {
	// 保证每次拉取事件后都清空事件列表，避免重复拉取同一批事件
	events := b.events
	b.ClearEvents()

	return events
}
