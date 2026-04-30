package event_handler

import (
	"context"
	"errors"

	"poprako-s/internal/domain/event"
)

type UnitSaveHandler struct{}

func NewUnitSaveHandler() event.EventHandler {
	return &UnitSaveHandler{}
}

func (h *UnitSaveHandler) EventType() event.EventType {
	return event.EventTypeUnitSave
}

func (h *UnitSaveHandler) PubType() event.PubType {
	return event.PubTypeSync
}

func (h *UnitSaveHandler) Handle(cx context.Context, ev event.Event) error {
	e, ok := ev.(*event.UnitSaveEvent)
	if !ok {
		return errors.New("[UnitSaveHandler] 事件载荷类型不是 UnitSaveEvent")
	}

	// 从事务上下文中创建绑定事务的 repo
	pageRepo, err := pageRepoFromCx(cx)
	if err != nil {
		return errors.New("[UnitSaveHandler] 无法创建事务版 PageRepo")
	}

	chapterRepo, err := chapterRepoFromCx(cx)
	if err != nil {
		return errors.New("[UnitSaveHandler] 无法创建事务版 ChapterRepo")
	}

	// 锁定 page 和 chapter 行，确保强一致性
	if err := pageRepo.LockByID(e.PageID); err != nil {
		return errors.New("[UnitSaveHandler] 锁定页面行失败")
	}

	if err := chapterRepo.LockByID(e.ChapterID); err != nil {
		return errors.New("[UnitSaveHandler] 锁定章节行失败")
	}

	// 使用原子增量更新 Page 统计
	if err := pageRepo.UpdateStats(e.PageID, e.TotalDelta, e.TranslatedDelta, e.ProofreadDelta); err != nil {
		return errors.New("[UnitSaveHandler] 更新页面统计失败")
	}

	// 使用原子增量更新 Chapter 统计
	if err := chapterRepo.UpdateStats(e.ChapterID, e.TotalDelta, e.TranslatedDelta, e.ProofreadDelta); err != nil {
		return errors.New("[UnitSaveHandler] 更新章节统计失败")
	}

	return nil
}
