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
	// Comic 仅在 includes 指定时填充
	Comic *ComicInfo

	// 是否顶置，通常一个 chapter 被刚创建出来时
	// 会自动变成该 comic 中的顶置章节
	// 这是为了能够快速查询到对应章节的数据
	IsPinned bool

	Index    int
	Subtitle string

	PageCount           int
	TotalUnitCount      int
	TranslatedUnitCount int
	ProofreadUnitCount  int

	CreatorID string
	// Creator 仅在 includes 指定时填充
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

type ChapterCreation struct {
	// ID 由 domain service 生成，外部不提供
	ID string

	// ComicID 是必须的，因为一个章节必须属于一个漫画
	ComicID string
	// Index 是必须的，因为章节在漫画中有一个明确的位置
	// 这是在事务中锁住 ComicID 所属漫画的所有章节 COUNT 得到的
	Index int
	// Subtitle 是可选的，因为有些章节可能没有副标题
	// 默认为 "Ch.{Index}"，但用户也可以提供一个更具描述性的副标题
	Subtitle *string
	// CreatorID 是必须的，因为我们需要知道是谁创建了这个章节
	CreatorID string
	// IsPinned 是可选的，因为默认情况下新创建的章节会被 **自动** 顶置
	// 但用户也可以选择不顶置
	IsPinned *bool
}
