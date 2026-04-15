package event

const (
	// EventTypeUserLogin 是 UserLoginEvent 的事件类型
	EventTypeUserLogin EventType = "UserLogin"
	// EventTypeUserCreated 是 UserCreatedEvent 的事件类型
	EventTypeUserCreated EventType = "UserCreated"
)

// UserLoginEvent 表示用户通过 QQ 登录的事件
type UserLoginEvent struct {
	UserQQ string
}

func (*UserLoginEvent) EventType() EventType {
	return EventTypeUserLogin
}

// Payload 返回 *UserLoginEvent 本身，供事件处理器使用
func (e *UserLoginEvent) Payload() any {
	return e
}

type UserCreatedEvent struct {
	InvitorID     string
	CreatedUserID string
}

func (*UserCreatedEvent) EventType() EventType {
	return EventTypeUserCreated
}

func (e *UserCreatedEvent) Payload() any {
	return e
}
