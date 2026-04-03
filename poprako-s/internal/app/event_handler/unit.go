package event_handler

import (
	"errors"

	"poprako-s/internal/domain/event"
	"poprako-s/internal/domain/model"
)

type UnitSaveHandler struct{}

func NewUnitSaveHandler() event.EventHandler {
	return &UnitSaveHandler{}
}

func (h *UnitSaveHandler) EventType() event.EventType {
	return event.EventTypeUnitSave
}

func (h *UnitSaveHandler) Handle(ev event.Event) error {
	e, ok := ev.(*event.UnitSaveEvent)
	if !ok {
		return errors.New("[UnitSaveHandler] 事件载荷类型不是 UnitSaveEvent")
	}

	// 从事务上下文中创建绑定事务的 repo
	pageRepo, err := pageRepoFromCx(e.Cx)
	if err != nil {
		return errors.New("[UnitSaveHandler] 无法创建事务版 PageRepo")
	}

	chapterRepo, err := chapterRepoFromCx(e.Cx)
	if err != nil {
		return errors.New("[UnitSaveHandler] 无法创建事务版 ChapterRepo")
	}

	// 更新 Page 统计
	pageStats, err := pageRepo.GetStatsByID(e.PageID)
	if err != nil {
		return errors.New("[UnitSaveHandler] 获取页面统计数据失败")
	}

	updatedPage := &model.PageStats{
		PageID:              e.PageID,
		TotalUnitCount:      pageStats.TotalUnitCount + e.TotalDelta,
		TranslatedUnitCount: pageStats.TranslatedUnitCount + e.TranslatedDelta,
		ProofreadUnitCount:  pageStats.ProofreadUnitCount + e.ProofreadDelta,
	}

	if updatedPage.TotalUnitCount < 0 ||
		updatedPage.TranslatedUnitCount < 0 ||
		updatedPage.ProofreadUnitCount < 0 {
		return errors.New("[UnitSaveHandler] 页面统计值出现负数")
	}

	if err := pageRepo.UpdateStats(updatedPage); err != nil {
		return errors.New("[UnitSaveHandler] 更新页面统计失败")
	}

	// 更新 Chapter 统计
	chapterInfo, err := chapterRepo.GetByID(e.ChapterID)
	if err != nil {
		return errors.New("[UnitSaveHandler] 获取章节信息失败")
	}

	updatedChapter := &model.ChapterStats{
		ChapterID:           e.ChapterID,
		TotalUnitCount:      chapterInfo.TotalUnitCount + e.TotalDelta,
		TranslatedUnitCount: chapterInfo.TranslatedUnitCount + e.TranslatedDelta,
		ProofreadUnitCount:  chapterInfo.ProofreadUnitCount + e.ProofreadDelta,
	}

	if updatedChapter.TotalUnitCount < 0 ||
		updatedChapter.TranslatedUnitCount < 0 ||
		updatedChapter.ProofreadUnitCount < 0 {
		return errors.New("[UnitSaveHandler] 章节统计值出现负数")
	}

	if err := chapterRepo.UpdateStats(updatedChapter); err != nil {
		return errors.New("[UnitSaveHandler] 更新章节统计失败")
	}

	return nil
}
