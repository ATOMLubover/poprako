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

// Revert transition identifiers used by `Chapter.RevertWorkflow`
// Each clears the corresponding timestamp back to NULL
// `publish_complete` has no revert variant — publish is irreversible
const (
	// `WorkflowUploadRevert` clears the upload-complete timestamp
	WorkflowUploadRevert WorkflowTransition = "upload_revert"
	// `WorkflowTranslateStartRevert` clears the translate-start timestamp
	WorkflowTranslateStartRevert WorkflowTransition = "translate_start_revert"
	// `WorkflowTranslateRevert` clears the translate-complete timestamp
	WorkflowTranslateRevert WorkflowTransition = "translate_revert"
	// `WorkflowProofreadStartRevert` clears the proofread-start timestamp
	WorkflowProofreadStartRevert WorkflowTransition = "proofread_start_revert"
	// `WorkflowProofreadRevert` clears the proofread-complete timestamp
	WorkflowProofreadRevert WorkflowTransition = "proofread_revert"
	// `WorkflowTypesetStartRevert` clears the typeset-start timestamp
	WorkflowTypesetStartRevert WorkflowTransition = "typeset_start_revert"
	// `WorkflowTypesetRevert` clears the typeset-complete timestamp
	WorkflowTypesetRevert WorkflowTransition = "typeset_revert"
	// `WorkflowReviewRevert` clears the review-complete timestamp
	WorkflowReviewRevert WorkflowTransition = "review_revert"
)
