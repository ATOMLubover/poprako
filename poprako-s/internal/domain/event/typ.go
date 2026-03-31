package event

// EventType 用于类型安全地标识事件类型
type EventType string

// PubType 用于标识一个事件是需要被同步处理还是异步处理
type PubType int

const (
	PubTypeSync PubType = iota
	PubTypeAsync
)
