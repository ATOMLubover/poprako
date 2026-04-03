package event

// EventBus 是一个事件总线接口，定义了发布和订阅事件的方法
// 支持同步地消费事件（Pub）和异步地消费事件（PubAsync），以及订阅（Sub）和取消订阅（Unsub）事件处理器的方法
type EventBus interface {
	// Pub 使用同步方式发布事件，等待所有处理器完成后返回
	// 注意，它会截取调用时的一个 snapshot 来保证在事件处理过程中分发映射的一致性
	// 这意味着在事件处理过程中对订阅的修改不会影响当前事件的分发，但会影响后续事件的分发
	Pub(ev []Event) error
	// PubAsync 使用异步方式发布事件，将事件发送到工作池中处理，立即返回
	// 注意，它不会截取调用时的一个 snapshot，而是直接使用当前的分发映射来处理事件
	// 而且它不保证所有的事件都原子地入队成功，如果事件总线过载，可能会返回错误，导致部分事件无法处理
	PubAsync(ev []Event) error
	// Sub 订阅事件处理器，使用线程安全的方式修改分发映射，适用于多线程环境
	Sub(h EventHandler) error
	// SubUnsafe 订阅事件处理器，直接修改当前的分发映射，不进行复制，适用于单线程环境或已知没有并发访问的情况
	SubUnsafe(h EventHandler) error
	// Unsub 取消订阅事件处理器，使用线程安全的方式修改分发映射，适用于多线程环境
	Unsub(h EventHandler) error
	// UnsubUnsafe 取消订阅事件处理器，直接修改当前的分发映射，不进行复制，适用于单线程环境或已知没有并发访问的情况
	UnsubUnsafe(h EventHandler) error

	Close()
}
