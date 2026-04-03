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

func (*UserLoginEvent) PubType() PubType {
	return PubTypeAsync
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

func (*UserCreatedEvent) PubType() PubType {
	// 用户创建事件需要被同步处理，以保证在用户创建后立即执行相关的后续操作
	// 如给对应的邀请者发送通知等
	// 这都要求必须在事务中处理完用户创建事件后才能执行，因此需要同步处理
	return PubTypeSync
}

func (e *UserCreatedEvent) Payload() any {
	return e
}
