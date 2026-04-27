package query

import "poprako-s/internal/domain/model/enum"

// `ListComicOpt` holds optional filters for listing comics
// All fields are pointer types nil means the filter is not applied
// Pagination is carried by `Pagi` to keep list interfaces explicit
// Query options are domain-level contracts not infra implementation details
type ListComicOpt struct {
	// `WorksetId` restricts result to comics that belong to the target workset
	WorksetId *string

	// `FuzzyTitle` applies case-insensitive fuzzy matching on comic title
	FuzzyTitle *string

	// Workflow phase filters based on pinned chapter progress.
	UploadPhase *enum.WorkflowPhase

	TranslatePhase *enum.WorkflowPhase

	ProofreadPhase *enum.WorkflowPhase

	TypesetPhase *enum.WorkflowPhase

	ReviewPhase *enum.WorkflowPhase

	PublishPhase *enum.WorkflowPhase

	// `Pagi` specifies pagination options for list query
	Pagi PagiOpt
}
