package main

import (
	"fmt"

	event_infra "poprako-s/internal/infra/event"
)

func main() {
	eventBus, err := event_infra.NewEventBus()
	if err != nil {
		panic(fmt.Sprintf("初始化事件总线失败: %v", err))
	}
}
