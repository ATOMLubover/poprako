package event

// UserLoginEvent 表示用户通过 QQ 登录的事件
type UserLoginEvent struct {
	UserQQ string
}

const EventTypeUserLogin EventType = "UserLoginEvent"

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
