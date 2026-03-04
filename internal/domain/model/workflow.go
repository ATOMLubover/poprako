package model

type Workflow string

const (
	WorkflowRawProviding Workflow = "raw_providing"
	WorkflowTranslating  Workflow = "translating"
	WorkflowProofreading Workflow = "proofreading"
	WorkflowTypesetting  Workflow = "typesetting"
	WorkflowReviewing    Workflow = "reviewing"
	WorkflowUploading    Workflow = "uploading"
)

type WorkflowStatus string

const (
	WorkflowPending    WorkflowStatus = "pending"
	WorkflowInProgress WorkflowStatus = "in_progress"
	WorkflowCompleted  WorkflowStatus = "completed"
)
