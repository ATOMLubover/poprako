package main

import (
	"fmt"

	"poprako-s/internal/domain/event"
)

func main() {
	eventBus, err := event.NewEventBus()
	if err != nil {
		panic(fmt.Sprintf("初始化事件总线失败: %v", err))
	}
}
