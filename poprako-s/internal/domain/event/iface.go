package event

// Event 是领域事件的抽象，它类型擦除了具体的类型
// 保证可以在总线上携带数据地传输
type Event interface {
	EventType() EventType
	// Payload 返回事件的载荷数据，类型为 any，可以是任意类型
	// 应当使用类型断言来获取具体的载荷类型，以最大可能地保障类型安全
	Payload() any
}

// EventSource 是一个接口，表示一个事件源
// 通常一个聚合根会实现这个接口，以便在其生命周期内产生事件
type EventSource interface {
	Events() []Event
}

// EventHandler 是一个接口，表示一个事件处理器
// EventDispatcher 会根据 EventType 将事件分发给对应的 EventHandler 来处理
type EventHandler interface {
	Handle(event Event) error
}
