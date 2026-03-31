package model

import "time"

// Workflow 表示工作流阶段的位移常量集合，用于在掩码中编码各阶段状态
type Workflow uint32

// 之所以使用 2 * iota 步进，是因为需要在生成掩码时
// 为每个阶段保留两个位，即 0x0~0x3 用于表示阶段的状态（Pending、Ongoing、Completed）
// 剩下的位用于保证扩展性
const (
	WorkflowUpload Workflow = iota * 2
	WorkflowTranslate
	WorkflowProofread
	WorkflowTypesett
	WorkflowReview
	WorkflowPublish
)

// WorkflowPhase 表示单个阶段的状态（待处理/进行中/已完成）
type WorkflowPhase uint8

const (
	WorkflowPending WorkflowPhase = iota
	WorkflowOngoing
	WorkflowCompleted
)

// UploadStatus 表示上传阶段的状态和完成时间
type UploadStatus struct {
	Status      WorkflowPhase
	CompletedAt *time.Time
}

// TranslateStatus 表示翻译阶段的状态和时间点
type TranslateStatus struct {
	Status      WorkflowPhase
	StartedAt   *time.Time
	CompletedAt *time.Time
}

// ProofreadStatus 表示校对阶段的状态和时间点
type ProofreadStatus struct {
	Status      WorkflowPhase
	StartedAt   *time.Time
	CompletedAt *time.Time
}

// TypesettStatus 表示嵌字阶段的状态和时间点
type TypesettStatus struct {
	Status      WorkflowPhase
	StartedAt   *time.Time
	CompletedAt *time.Time
}

// ReviewStatus 表示监修阶段的状态
type ReviewStatus struct {
	Status      WorkflowPhase
	CompletedAt *time.Time
}

// PublishStatus 表示发布阶段的状态
type PublishStatus struct {
	Status      WorkflowPhase
	CompletedAt *time.Time
}

// WorkflowTransition 描述了一个章节工作流改变事件
// 它以扁平的方式定义了可能的工作流事件类型，例如上传完成、翻译开始等
type WorkflowTransition string

const (
	WorkflowUploadComplete    WorkflowTransition = "upload_complete"
	WorkflowTranslateStart    WorkflowTransition = "translate_start"
	WorkflowTranslateComplete WorkflowTransition = "translate_complete"
	WorkflowProofreadStart    WorkflowTransition = "proofread_start"
	WorkflowProofreadComplete WorkflowTransition = "proofread_complete"
	WorkflowTypesetStart      WorkflowTransition = "typeset_start"
	WorkflowTypesetComplete   WorkflowTransition = "typeset_complete"
	WorkflowReviewComplete    WorkflowTransition = "review_complete"
	WorkflowPublishComplete   WorkflowTransition = "publish_complete"
)
