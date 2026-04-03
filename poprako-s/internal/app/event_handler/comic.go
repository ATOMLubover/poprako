package event_handler

import (
	"errors"

	"poprako-s/internal/domain/event"
)

type ComicCreateHandler struct{}

func NewComicCreateHandler() event.EventHandler {
	return &ComicCreateHandler{}
}

func (h *ComicCreateHandler) EventType() event.EventType {
	return event.EventTypeComicCreated
}

func (h *ComicCreateHandler) Handle(ev event.Event) error {
	e, ok := ev.(*event.ComicCreatedEvent)
	if !ok {
		return errors.New("[ComicCreateHandler] 事件载荷类型不是 ComicCreatedEvent")
	}

	worksetRepo, err := worksetRepoFromCx(e.Cx)
	if err != nil {
		return errors.New("[ComicCreateHandler] 无法创建事务版 WorksetRepo")
	}

	return worksetRepo.UpdateComicCount(e.WorksetID, 1)
}

type ComicRemoveHandler struct{}

func NewComicRemoveHandler() event.EventHandler {
	return &ComicRemoveHandler{}
}

func (h *ComicRemoveHandler) EventType() event.EventType {
	return event.EventTypeComicRemoved
}

func (h *ComicRemoveHandler) Handle(ev event.Event) error {
	e, ok := ev.(*event.ComicRemovedEvent)
	if !ok {
		return errors.New("[ComicRemoveHandler] 事件载荷类型不是 ComicRemovedEvent")
	}

	worksetRepo, err := worksetRepoFromCx(e.Cx)
	if err != nil {
		return errors.New("[ComicRemoveHandler] 无法创建事务版 WorksetRepo")
	}

	return worksetRepo.UpdateComicCount(e.WorksetID, -1)
}
