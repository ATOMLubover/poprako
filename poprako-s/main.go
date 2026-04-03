package main

import (
	"fmt"

	event_handler "poprako-s/internal/app/event_handler"
	"poprako-s/internal/domain/event"
	event_infra "poprako-s/internal/infra/event"
)

func main() {
	eventBus, err := event_infra.NewEventBus()
	if err != nil {
		panic(fmt.Sprintf("初始化事件总线失败: %v", err))
	}

	defer eventBus.Close()

	handlers := []event.EventHandler{
		event_handler.NewUnitSaveHandler(),
		event_handler.NewAssignmentCreateHandler(),
		event_handler.NewAssignmentRemoveHandler(),
		event_handler.NewComicCreateHandler(),
		event_handler.NewComicRemoveHandler(),
		event_handler.NewChapterCreateHandler(),
		event_handler.NewChapterRemoveHandler(),
		event_handler.NewChapterPublishedHandler(),
	}

	for _, handler := range handlers {
		if err := eventBus.SubUnsafe(handler); err != nil {
			panic(fmt.Sprintf("注册事件处理器失败: %v", err))
		}
	}
}
