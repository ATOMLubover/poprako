package model

import (
	"errors"
	"fmt"
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

	UploadedAt     *time.Time
	TransalatingAt *time.Time
	TranslatedAt   *time.Time
	ProofreadingAt *time.Time
	ProofreadAt    *time.Time
	TypesettingAt  *time.Time
	TypesetAt      *time.Time
	ReviewedAt     *time.Time
	PublishedAt    *time.Time

	// 在工作流转换过程中产生的事件
	events []event.Event
}

// TransiteWorkflow 接受一个工作流转换事件，根据事件类型和当前状态执行相应的状态转换
func (c *ChapterInfo) TransiteWorkflow(t WorkflowTransition) error {
	now := time.Now()

	markStarted := func(startedAt **time.Time, completedAt *time.Time, ongoingMsg, completedMsg string) error {
		if completedAt != nil {
			return errors.New(completedMsg)
		}
		if *startedAt != nil {
			return errors.New(ongoingMsg)
		}

		*startedAt = &now
		return nil
	}

	markCompleted := func(startedAt *time.Time, completedAt **time.Time, pendingMsg, completedMsg string) error {
		if *completedAt != nil {
			return errors.New(completedMsg)
		}
		if startedAt == nil {
			return errors.New(pendingMsg)
		}

		*completedAt = &now
		return nil
	}

	markOnce := func(completedAt **time.Time, completedMsg string) error {
		if *completedAt != nil {
			return errors.New(completedMsg)
		}

		*completedAt = &now
		return nil
	}

	switch t {
	case WorkflowUploadComplete:
		if err := markOnce(&c.UploadedAt, "上传进度已经标记为完成，无法重复标记"); err != nil {
			return err
		}

		c.events = append(c.events, &event.WorkflowUploadCompletedEvent{ChapterID: c.ID})

	case WorkflowTranslateStart:
		if err := markStarted(&c.TransalatingAt, c.TranslatedAt, "翻译进度已经处于进行中，无法重复开始", "翻译进度已经标记为完成，无法重新开始"); err != nil {
			return err
		}

	case WorkflowTranslateComplete:
		if err := markCompleted(c.TransalatingAt, &c.TranslatedAt, "翻译进度尚未开始，无法标记为完成", "翻译进度已经标记为完成，无法重复标记"); err != nil {
			return err
		}

		c.events = append(c.events, &event.WorkflowTranslateCompletedEvent{ChapterID: c.ID})

	case WorkflowProofreadStart:
		if err := markStarted(&c.ProofreadingAt, c.ProofreadAt, "校对进度已经处于进行中，无法重复开始", "校对进度已经标记为完成，无法重新开始"); err != nil {
			return err
		}

	case WorkflowProofreadComplete:
		if err := markCompleted(c.ProofreadingAt, &c.ProofreadAt, "校对进度尚未开始，无法标记为完成", "校对进度已经标记为完成，无法重复标记"); err != nil {
			return err
		}

		c.events = append(c.events, &event.WorkflowProofreadCompletedEvent{ChapterID: c.ID})

	case WorkflowTypesetStart:
		if err := markStarted(&c.TypesettingAt, c.TypesetAt, "嵌字进度已经处于进行中，无法重复开始", "嵌字进度已经标记为完成，无法重新开始"); err != nil {
			return err
		}

	case WorkflowTypesetComplete:
		if err := markCompleted(c.TypesettingAt, &c.TypesetAt, "嵌字进度尚未开始，无法标记为完成", "嵌字进度已经标记为完成，无法重复标记"); err != nil {
			return err
		}

		c.events = append(c.events, &event.WorkflowTypesettCompletedEvent{ChapterID: c.ID})

	case WorkflowReviewComplete:
		if err := markOnce(&c.ReviewedAt, "监修进度已经标记为完成，无法重复标记"); err != nil {
			return err
		}

		c.events = append(c.events, &event.WorkflowReviewCompletedEvent{ChapterID: c.ID})

	case WorkflowPublishComplete:
		if err := markOnce(&c.PublishedAt, "发布进度已经标记为完成，无法重复标记"); err != nil {
			return err
		}

		c.events = append(c.events, &event.ChapterPublishedEvent{ChapterID: c.ID})

		c.events = append(c.events, &event.WorkflowPublishCompletedEvent{ChapterID: c.ID})

	default:
		return fmt.Errorf("未知的工作流转换：%s", t)
	}

	return nil
}

// PullEvents 实现 EventSource 接口，返回并清除所有待发布的事件
func (c *ChapterInfo) PullEvents() []event.Event {
	evs := c.events
	c.events = nil
	return evs
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

	event.EventBase
}

// ChapterStats 聚合一个章节的统计数据
type ChapterStats struct {
	ChapterID string

	TotalUnitCount      int
	TranslatedUnitCount int
	ProofreadUnitCount  int
}

// ChapterUpdate 用于 PUT 语义的章节全量替换，保留已有字段的值
type ChapterUpdate struct {
	ID string

	// ComicID 不允许修改，因为章节与漫画的关系一旦建立就不应该改变
	// Index 不允许修改，因为章节在漫画中的位置由创建时确定，修改可能导致数据不一致

	// Subtitle 是可选的，因为用户可能只想修改副标题而不改变其他字段
	Subtitle *string
	// IsPinned 是可选的，因为用户可能只想修改是否顶置而不改变其他字段
	IsPinned *bool

	UploadedAt     *time.Time
	TransalatingAt *time.Time
	TranslatedAt   *time.Time
	ProofreadingAt *time.Time
	ProofreadAt    *time.Time
	TypesettingAt  *time.Time
	TypesetAt      *time.Time
	ReviewedAt     *time.Time
	PublishedAt    *time.Time
}

// ChapterQueryOpt 指定章节查询的可选筛选条件
// 所有字段均为可空，nil 表示不参与筛选
type ChapterQueryOpt struct {
	// ID 按章节 ID 筛选
	ID *string
	// ComicID 按所属漫画 ID 筛选
	ComicID *string
	// CreatorID 按创建者 ID 筛选
	CreatorID *string
}
