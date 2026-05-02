package val

import "poprako-s/internal/domain/model/aggr"

// `AssignmentVal` is app-facing value object for assignment.
type AssignmentVal struct {
	Id string `json:"id"`

	ChapterId string `json:"chapter_id"`
	UserId    string `json:"user_id"`

	RoleMask aggr.RoleMask `json:"role_mask"`

	AssignedRawProviderAt *int64 `json:"assigned_raw_provider_at"`
	AssignedTranslatorAt  *int64 `json:"assigned_translator_at"`
	AssignedProofreaderAt *int64 `json:"assigned_proofreader_at"`
	AssignedTypesetterAt  *int64 `json:"assigned_typesetter_at"`
	AssignedRedrawerAt    *int64 `json:"assigned_redrawer_at"`
	AssignedReviewerAt    *int64 `json:"assigned_reviewer_at"`
	AssignedPublisherAt   *int64 `json:"assigned_publisher_at"`

	CreatedAt int64 `json:"created_at"`
	UpdatedAt int64 `json:"updated_at"`
}

// `ListAssignmentByChapterArgs` carries chapter list args for assignment.
type ListAssignmentByChapterArgs struct {
	ChapterId string `url:"chapter_id"`

	Offset int `url:"offset"`
	Limit  int `url:"limit"`
}

// `ListAssignmentByUserArgs` carries my list args for assignment.
type ListAssignmentByUserArgs struct {
	Offset int `url:"offset"`
	Limit  int `url:"limit"`
}

// `UpsertAssignmentArgs` carries upsert args for assignment.
type UpsertAssignmentArgs struct {
	ChapterId string        `json:"chapter_id"`
	UserId    string        `json:"user_id"`
	RoleMask  aggr.RoleMask `json:"role_mask"`
}
