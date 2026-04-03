package event

import "context"

// 工作流发生变化时产生的事件类型
// 我们通常只记录工作流完成的事件，因为它们是最重要的里程碑事件
const (
	EventTypeWorkflowUploadCompleted    EventType = "WorkflowUploadCompleted"
	EventTypeWorkflowTranslateCompleted EventType = "WorkflowTranslateCompleted"
	EventTypeWorkflowProofreadCompleted EventType = "WorkflowProofreadCompleted"
	EventTypeWorkflowTypesettCompleted  EventType = "WorkflowTypesettCompleted"
	EventTypeWorkflowReviewCompleted    EventType = "WorkflowReviewCompleted"
	EventTypeWorkflowPublishCompleted   EventType = "WorkflowPublishCompleted"
	EventTypeChapterCreated             EventType = "ChapterCreated"
	EventTypeChapterRemoved             EventType = "ChapterRemoved"
	EventTypeChapterPublished           EventType = "ChapterPublished"
)

// ChapterCreatedEvent 代表章节创建后的同步统计事件
type ChapterCreatedEvent struct {
	ComicID string
	Cx      context.Context
}

func (e *ChapterCreatedEvent) EventType() EventType {
	return EventTypeChapterCreated
}

func (e *ChapterCreatedEvent) PubType() PubType {
	return PubTypeSync
}

func (e *ChapterCreatedEvent) Payload() any {
	return e
}

// ChapterRemovedEvent 代表章节删除后的同步统计事件
type ChapterRemovedEvent struct {
	ComicID         string
	WasPublished    bool
	AssignedUserIDs []string
	Cx              context.Context
}

func (e *ChapterRemovedEvent) EventType() EventType {
	return EventTypeChapterRemoved
}

func (e *ChapterRemovedEvent) PubType() PubType {
	return PubTypeSync
}

func (e *ChapterRemovedEvent) Payload() any {
	return e
}

// ChapterPublishedEvent 代表章节发布完成后的同步统计事件
type ChapterPublishedEvent struct {
	ChapterID string
	Cx        context.Context
}

func (e *ChapterPublishedEvent) EventType() EventType {
	return EventTypeChapterPublished
}

func (e *ChapterPublishedEvent) PubType() PubType {
	return PubTypeSync
}

func (e *ChapterPublishedEvent) Payload() any {
	return e
}

// WorkflowUploadCompletedEvent 代表章节上传流程完成事件
type WorkflowUploadCompletedEvent struct {
	ChapterID string
}

// EventType 返回事件的标识符
func (e *WorkflowUploadCompletedEvent) EventType() EventType {
	return EventTypeWorkflowUploadCompleted
}

// PubType 返回事件的发布类型
func (e *WorkflowUploadCompletedEvent) PubType() PubType {
	return PubTypeAsync
}

// Payload 返回事件载荷
func (e *WorkflowUploadCompletedEvent) Payload() any {
	return e
}

// WorkflowTranslateCompletedEvent 代表章节翻译流程完成事件
type WorkflowTranslateCompletedEvent struct {
	ChapterID string
}

// EventType 返回事件的标识符
func (e *WorkflowTranslateCompletedEvent) EventType() EventType {
	return EventTypeWorkflowTranslateCompleted
}

// PubType 返回事件的发布类型
func (e *WorkflowTranslateCompletedEvent) PubType() PubType {
	return PubTypeAsync
}

// Payload 返回事件载荷
func (e *WorkflowTranslateCompletedEvent) Payload() any {
	return e
}

// WorkflowProofreadCompletedEvent 代表章节校对流程完成事件
type WorkflowProofreadCompletedEvent struct {
	ChapterID string
}

// EventType 返回事件的标识符
func (e *WorkflowProofreadCompletedEvent) EventType() EventType {
	return EventTypeWorkflowProofreadCompleted
}

// PubType 返回事件的发布类型
func (e *WorkflowProofreadCompletedEvent) PubType() PubType {
	return PubTypeAsync
}

// Payload 返回事件载荷
func (e *WorkflowProofreadCompletedEvent) Payload() any {
	return e
}

// WorkflowTypesettCompletedEvent 代表章节嵌字流程完成事件
type WorkflowTypesettCompletedEvent struct {
	ChapterID string
}

// EventType 返回事件的标识符
func (e *WorkflowTypesettCompletedEvent) EventType() EventType {
	return EventTypeWorkflowTypesettCompleted
}

// PubType 返回事件的发布类型
func (e *WorkflowTypesettCompletedEvent) PubType() PubType {
	return PubTypeAsync
}

// Payload 返回事件载荷
func (e *WorkflowTypesettCompletedEvent) Payload() any {
	return e
}

// WorkflowReviewCompletedEvent 代表章节监修流程完成事件
type WorkflowReviewCompletedEvent struct {
	ChapterID string
}

// EventType 返回事件的标识符
func (e *WorkflowReviewCompletedEvent) EventType() EventType {
	return EventTypeWorkflowReviewCompleted
}

// PubType 返回事件的发布类型
func (e *WorkflowReviewCompletedEvent) PubType() PubType {
	return PubTypeAsync
}

// Payload 返回事件载荷
func (e *WorkflowReviewCompletedEvent) Payload() any {
	return e
}

// WorkflowPublishCompletedEvent 代表章节发布流程完成事件
type WorkflowPublishCompletedEvent struct {
	ChapterID string
}

// EventType 返回事件的标识符
func (e *WorkflowPublishCompletedEvent) EventType() EventType {
	return EventTypeWorkflowPublishCompleted
}

// PubType 返回事件的发布类型
func (e *WorkflowPublishCompletedEvent) PubType() PubType {
	return PubTypeAsync
}

// Payload 返回事件载荷
func (e *WorkflowPublishCompletedEvent) Payload() any {
	return e
}
