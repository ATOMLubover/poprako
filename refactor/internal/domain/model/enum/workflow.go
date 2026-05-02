package enum

// `WorkflowPhase` represents one workflow phase status in chapter progress.
// It is reused by comic list filters that depend on pinned chapter progress.
// `WorkflowPending` means the phase has not started.
// `WorkflowOngoing` means started but not completed.
// `WorkflowCompleted` means the phase has reached its terminal state.
// Note: some phases do not support `Ongoing`; the repo filter returns no rows for those.
type WorkflowPhase uint8

const (
	// `WorkflowPending` means workflow phase has not started.
	WorkflowPending WorkflowPhase = iota
	// `WorkflowOngoing` means workflow phase started but not completed.
	WorkflowOngoing
	// `WorkflowCompleted` means workflow phase is completed.
	WorkflowCompleted
)

// `WorkflowTransition` describes a chapter workflow transition event.
// Each value maps to one deterministic timestamp mutation in the chapter aggregate.
// `Start` transitions set the phase-ing timestamp; `Complete` transitions set the phase-done timestamp.
// Upload, review, and publish only have complete transitions.
type WorkflowTransition string

const (
	// `WorkflowUploadComplete` marks upload phase as completed
	WorkflowUploadComplete WorkflowTransition = "upload_complete"
	// `WorkflowTranslateStart` marks translate phase as started
	WorkflowTranslateStart WorkflowTransition = "translate_start"
	// `WorkflowTranslateComplete` marks translate phase as completed
	WorkflowTranslateComplete WorkflowTransition = "translate_complete"
	// `WorkflowProofreadStart` marks proofread phase as started
	WorkflowProofreadStart WorkflowTransition = "proofread_start"
	// `WorkflowProofreadComplete` marks proofread phase as completed
	WorkflowProofreadComplete WorkflowTransition = "proofread_complete"
	// `WorkflowTypesetStart` marks typeset phase as started
	WorkflowTypesetStart WorkflowTransition = "typeset_start"
	// `WorkflowTypesetComplete` marks typeset phase as completed
	WorkflowTypesetComplete WorkflowTransition = "typeset_complete"
	// `WorkflowReviewComplete` marks review phase as completed
	WorkflowReviewComplete WorkflowTransition = "review_complete"
	// `WorkflowPublishComplete` marks publish phase as completed
	WorkflowPublishComplete WorkflowTransition = "publish_complete"
)
