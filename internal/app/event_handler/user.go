package event_handler

import (
	"context"
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

func (h *UserLoginHandler) PubType() event.PubType {
	return event.PubTypeAsync
}

func (h *UserLoginHandler) Handle(_ context.Context, ev event.Event) error {
	payload, ok := ev.(*event.UserLoginEvent)
	if !ok {
		return errors.New("[UserLoginHandler] 事件载荷类型不是 UserLoginEvent")
	}

	return h.userRepo.RefreshLastLogin(payload.UserQQ, time.Now())
}
