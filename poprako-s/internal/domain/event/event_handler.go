package event

// EventHandler 是一个接口，表示一个事件处理器
// EventDispatcher 会根据 EventType 将事件分发给对应的 EventHandler 来处理
type EventHandler interface {
	// EventType 返回该处理器能够处理的事件类型
	EventType() EventType
	// Handle 接收一个事件并处理它，返回处理结果
	Handle(ev Event) error
}
