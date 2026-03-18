package model

type Workflow string

const (
	WorkflowUploading    Workflow = "uploading"
	WorkflowTranslating  Workflow = "translating"
	WorkflowProofreading Workflow = "proofreading"
	WorkflowTypesetting  Workflow = "typesetting"
	WorkflowReviewing    Workflow = "reviewing"
	WorkflowPublishing   Workflow = "publishing"
)

type WorkflowStatus string

const (
	WorkflowPending    WorkflowStatus = "pending"
	WorkflowInProgress WorkflowStatus = "in_progress"
	WorkflowCompleted  WorkflowStatus = "completed"
	// 不参与筛选，或者不需要更新
	WorkflowUnset WorkflowStatus = "unset"
)

func IsValidWorkflowCombination(workflow Workflow, status WorkflowStatus) bool {
	switch workflow {
	case WorkflowUploading, WorkflowReviewing, WorkflowPublishing:
		// 这三个状态不区分进行中和未开始，只有待处理和已完成两种状态
		return status == WorkflowPending || status == WorkflowCompleted
	case WorkflowTranslating, WorkflowProofreading, WorkflowTypesetting:
		// 其他状态区分待处理、进行中和已完成三种状态
		return status == WorkflowPending || status == WorkflowInProgress || status == WorkflowCompleted
	}

	return false
}
