package event

import "context"

// EventHandler 是一个接口，表示一个事件处理器
// EventDispatcher 会根据 EventType 将事件分发给对应的 EventHandler 来处理
type EventHandler interface {
	// EventType 返回该处理器能够处理的事件类型
	// 这种自声明决定它如何被注册到事件总线上
	EventType() EventType
	// PubType 返回该处理器的发布类型，决定了它是同步处理事件还是异步处理事件
	PubType() PubType
	// Handle 接收一个事件并处理它，返回处理结果
	// cx 是必须的，因为它总是需要在一个隐式连贯的上下文中处理信息
	Handle(cx context.Context, ev Event) error
}
