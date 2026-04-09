package event_handler

import (
	"context"
	"errors"

	"poprako-s/internal/domain/event"
	"poprako-s/internal/domain/model"
	"poprako-s/internal/domain/service"
)

type ChapterCreateHandler struct{}

func NewChapterCreateHandler() event.EventHandler {
	return &ChapterCreateHandler{}
}

func (h *ChapterCreateHandler) EventType() event.EventType {
	return event.EventTypeChapterCreated
}

func (h *ChapterCreateHandler) PubType() event.PubType {
	return event.PubTypeAsync
}

func (h *ChapterCreateHandler) Handle(cx context.Context, ev event.Event) error {
	e, ok := ev.(*event.ChapterCreatedEvent)
	if !ok {
		return errors.New("[ChapterCreateHandler] 事件载荷类型不是 ChapterCreatedEvent")
	}

	comicRepo, err := comicRepoFromCx(cx)
	if err != nil {
		return errors.New("[ChapterCreateHandler] 无法创建事务版 ComicRepo")
	}

	return comicRepo.UpdateChapterCount(e.ComicID, 1)
}

type ChapterRemoveHandler struct{}

func NewChapterRemoveHandler() event.EventHandler {
	return &ChapterRemoveHandler{}
}

func (h *ChapterRemoveHandler) EventType() event.EventType {
	return event.EventTypeChapterRemoved
}

func (h *ChapterRemoveHandler) PubType() event.PubType {
	return event.PubTypeAsync
}

func (h *ChapterRemoveHandler) Handle(cx context.Context, ev event.Event) error {
	e, ok := ev.(*event.ChapterRemovedEvent)
	if !ok {
		return errors.New("[ChapterRemoveHandler] 事件载荷类型不是 ChapterRemovedEvent")
	}

	comicRepo, err := comicRepoFromCx(cx)
	if err != nil {
		return errors.New("[ChapterRemoveHandler] 无法创建事务版 ComicRepo")
	}

	if err := comicRepo.UpdateChapterCount(e.ComicID, -1); err != nil {
		return err
	}

	if e.WasPublished || len(e.AssignedUserIDs) == 0 {
		return nil
	}

	userRepo, err := userRepoFromCx(cx)
	if err != nil {
		return errors.New("[ChapterRemoveHandler] 无法创建事务版 UserRepo")
	}

	for _, userID := range e.AssignedUserIDs {
		if err := userRepo.PatchStats(&model.UserStatsPatch{
			UserID:                     userID,
			ActiveAssignmentCountDelta: -1,
		}); err != nil {
			return err
		}
	}

	return nil
}

type ChapterPublishedHandler struct{}

func NewChapterPublishedHandler() event.EventHandler {
	return &ChapterPublishedHandler{}
}

func (h *ChapterPublishedHandler) EventType() event.EventType {
	return event.EventTypeChapterPublished
}

func (h *ChapterPublishedHandler) PubType() event.PubType {
	return event.PubTypeAsync
}

func (h *ChapterPublishedHandler) Handle(cx context.Context, ev event.Event) error {
	e, ok := ev.(*event.ChapterPublishedEvent)
	if !ok {
		return errors.New("[ChapterPublishedHandler] 事件载荷类型不是 ChapterPublishedEvent")
	}

	assignmentRepo, err := assignmentRepoFromCx(cx)
	if err != nil {
		return errors.New("[ChapterPublishedHandler] 无法创建事务版 AssignmentRepo")
	}

	userRepo, err := userRepoFromCx(cx)
	if err != nil {
		return errors.New("[ChapterPublishedHandler] 无法创建事务版 UserRepo")
	}

	assignments, err := assignmentRepo.List(model.AssignmentQueryOpt{
		ChapterID: &e.ChapterID,
	})
	if err != nil {
		return err
	}

	for _, assignment := range assignments {
		if err := userRepo.PatchStats(&model.UserStatsPatch{
			UserID:                       assignment.UserID,
			ActiveAssignmentCountDelta:   -1,
			FinishedAssignmentCountDelta: 1,
		}); err != nil {
			return err
		}
	}

	return nil
}

// ChapterCreatorAssignedHandler 在章节创建时将创建者指定为监修
type ChapterCreatorAssignedHandler struct{}

// NewChapterCreatorAssignedHandler 返回 ChapterCreatorAssignedHandler 的默认实现
func NewChapterCreatorAssignedHandler() event.EventHandler {
	return &ChapterCreatorAssignedHandler{}
}

// EventType 返回该处理器订阅的事件类型
func (h *ChapterCreatorAssignedHandler) EventType() event.EventType {
	return event.EventTypeChapterCreatorAssigned
}

// PubType 返回该处理器的发布类型
func (h *ChapterCreatorAssignedHandler) PubType() event.PubType {
	return event.PubTypeAsync
}

// Handle 处理章节创建者分配监修事件
func (h *ChapterCreatorAssignedHandler) Handle(cx context.Context, ev event.Event) error {
	e, ok := ev.(*event.ChapterCreatorAssignedEvent)
	if !ok {
		return errors.New("[ChapterCreatorAssignedHandler] 事件载荷类型不是 ChapterCreatorAssignedEvent")
	}

	assignmentRepo, err := assignmentRepoFromCx(cx)
	if err != nil {
		return errors.New("[ChapterCreatorAssignedHandler] 无法创建事务版 AssignmentRepo")
	}

	// 由领域服务构造初始监修分配载荷
	svc := service.NewAssignmentService()

	creation := svc.NewInitReviewerCreation(e.ChapterID, e.CreatorID)

	// 持久化初始监修分配记录
	_, err = assignmentRepo.Create(creation)

	return err
}
