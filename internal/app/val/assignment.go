package val

import (
	"poprako-s/internal/domain/model/aggr"
	"poprako-s/internal/domain/model/enum"
)

// `AssignmentVal` is app-facing value object for assignment.
type AssignmentVal struct {
	Id string `json:"id"`

	ChapterId string      `json:"chapter_id"`
	UserId    string      `json:"user_id"`
	Chapter   *ChapterVal `json:"chapter,omitempty"`
	User      *UserVal    `json:"user,omitempty"`

	RoleMask aggr.RoleMask `json:"role_mask"`

	AssignedRawProviderAt *int64 `json:"assigned_raw_provider_at,omitempty"`
	AssignedTranslatorAt  *int64 `json:"assigned_translator_at,omitempty"`
	AssignedProofreaderAt *int64 `json:"assigned_proofreader_at,omitempty"`
	AssignedTypesetterAt  *int64 `json:"assigned_typesetter_at,omitempty"`
	AssignedRedrawerAt    *int64 `json:"assigned_redrawer_at,omitempty"`
	AssignedReviewerAt    *int64 `json:"assigned_reviewer_at,omitempty"`
	AssignedPublisherAt   *int64 `json:"assigned_publisher_at,omitempty"`

	CreatedAt int64 `json:"created_at"`
	UpdatedAt int64 `json:"updated_at"`
}

// `ListAssignmentByChapterArgs` carries chapter list args for assignment.
type ListAssignmentByChapterArgs struct {
	ChapterId string `url:"chapter_id"`

	Includes []enum.AssignmentIncl `url:"includes"`

	Offset int `url:"offset"`
	Limit  int `url:"limit"`
}

// `ListAssignmentByUserArgs` carries my list args for assignment.
type ListAssignmentByUserArgs struct {
	Includes []enum.AssignmentIncl `url:"includes"`

	Offset int `url:"offset"`
	Limit  int `url:"limit"`
}

// `UpsertAssignmentArgs` carries upsert args for assignment.
type UpsertAssignmentArgs struct {
	ChapterId string        `json:"chapter_id"`
	UserId    string        `json:"user_id"`
	RoleMask  aggr.RoleMask `json:"role_mask"`
}
