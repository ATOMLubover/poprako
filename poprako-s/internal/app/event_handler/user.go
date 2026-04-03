package event_handler

import (
	"errors"
	"time"

	"poprako-s/internal/domain/event"
	"poprako-s/internal/domain/repo"
)

type UserLoginHandler struct {
	userRepo repo.UserRepo
}

func (h *UserLoginHandler) EventType() event.EventType {
	return event.EventTypeUserLogin
}

func (h *UserLoginHandler) Handle(ev event.Event) error {
	payload, ok := ev.(*event.UserLoginEvent)
	if !ok {
		return errors.New("[UserLoginHandler] 事件载荷类型不是 UserLoginEvent")
	}

	return h.userRepo.RefreshLastLogin(payload.UserQQ, time.Now())
}
