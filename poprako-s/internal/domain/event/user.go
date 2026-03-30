package event

type UserLoginEvent struct {
	UserQQ string
}

const EventTypeUserLogin EventType = "UserLoginEvent"

func (*UserLoginEvent) EventType() EventType {
	return EventTypeUserLogin
}

// 返回 string 类型的 UserID，表示登录的用户 ID
func (e *UserLoginEvent) Payload() any {
	return e.UserQQ
}
