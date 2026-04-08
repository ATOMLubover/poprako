package event_handler

import (
	"context"
	"errors"

	"poprako-s/internal/domain/event"
	"poprako-s/internal/domain/model"
)

type AssignmentCreateHandler struct{}

func NewAssignmentCreateHandler() event.EventHandler {
	return &AssignmentCreateHandler{}
}

func (h *AssignmentCreateHandler) EventType() event.EventType {
	return event.EventTypeAssignmentCreated
}

func (h *AssignmentCreateHandler) Handle(cx context.Context, ev event.Event) error {
	e, ok := ev.(*event.AssignmentCreatedEvent)
	if !ok {
		return errors.New("[AssignmentCreateHandler] 事件载荷类型不是 AssignmentCreatedEvent")
	}

	userRepo, err := userRepoFromCx(cx)
	if err != nil {
		return errors.New("[AssignmentCreateHandler] 无法创建事务版 UserRepo")
	}

	return userRepo.PatchStats(&model.UserStatsPatch{
		UserID:                       e.UserID,
		TotalAssignmentCountDelta:    1,
		ActiveAssignmentCountDelta:   1,
		FinishedAssignmentCountDelta: 0,
	})
}

type AssignmentRemoveHandler struct{}

func NewAssignmentRemoveHandler() event.EventHandler {
	return &AssignmentRemoveHandler{}
}

func (h *AssignmentRemoveHandler) EventType() event.EventType {
	return event.EventTypeAssignmentRemoved
}

func (h *AssignmentRemoveHandler) Handle(cx context.Context, ev event.Event) error {
	e, ok := ev.(*event.AssignmentRemovedEvent)
	if !ok {
		return errors.New("[AssignmentRemoveHandler] 事件载荷类型不是 AssignmentRemovedEvent")
	}

	if e.WasPublished {
		return nil
	}

	userRepo, err := userRepoFromCx(cx)
	if err != nil {
		return errors.New("[AssignmentRemoveHandler] 无法创建事务版 UserRepo")
	}

	return userRepo.PatchStats(&model.UserStatsPatch{
		UserID:                     e.UserID,
		ActiveAssignmentCountDelta: -1,
	})
}
