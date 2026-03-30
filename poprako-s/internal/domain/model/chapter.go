package model

import (
	"time"

	"poprako-s/internal/domain/event"
)

// ChapterInfo 包含了一个漫画章节的核心内容
// 它提供的能力表现为状态机
type ChapterInfo struct {
	ID string

	ComicID string
	// Comic 仅在 includes 指定时填充。
	Comic *ComicInfo

	Index    int
	Subtitle string

	PageCount           int
	TotalUnitCount      int
	TranslatedUnitCount int
	ProofreadUnitCount  int

	CreatorID string
	// Creator 仅在 includes 指定时填充。
	Creator *UserInfo

	CreatedAt time.Time
	UpdatedAt time.Time

	// 所有进度相关数据的核心内聚在 state 中
	state workflowState

	// 在工作流转换过程中产生的事件
	events []event.Event
}

// Transite 接受一个工作流转换事件，根据事件类型和当前状态执行相应的状态转换
func (c *ChapterInfo) TransiteWorkflow(t WorkflowTransition) error {
	if err := c.state.transite(t); err != nil {
		return err
	}

	switch t {
	case WorkflowUploadComplete:
		c.events = append(c.events, &event.WorkflowUploadCompletedEvent{ChapterID: c.ID})
	case WorkflowTranslateComplete:
		c.events = append(c.events, &event.WorkflowTranslateCompletedEvent{ChapterID: c.ID})
	case WorkflowProofreadComplete:
		c.events = append(c.events, &event.WorkflowProofreadCompletedEvent{ChapterID: c.ID})
	case WorkflowTypesetComplete:
		c.events = append(c.events, &event.WorkflowTypesettCompletedEvent{ChapterID: c.ID})
	case WorkflowReviewComplete:
		c.events = append(c.events, &event.WorkflowReviewCompletedEvent{ChapterID: c.ID})
	case WorkflowPublishComplete:
		c.events = append(c.events, &event.WorkflowPublishCompletedEvent{ChapterID: c.ID})
	}

	return nil
}

// 实现 EventSource 接口
func (c *ChapterInfo) Events() []event.Event {
	return c.events
}
