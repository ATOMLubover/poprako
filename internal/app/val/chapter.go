package val

import (
	"poprako-s/internal/domain/model/aggr"
	"poprako-s/internal/domain/model/enum"
)

// `ChapterVal` is the app-facing chapter value object.
type ChapterVal struct {
	Id string `json:"id"`

	ComicId string    `json:"comic_id"`
	Comic   *ComicVal `json:"comic,omitempty"`

	IsPinned bool `json:"is_pinned"`

	Index    int    `json:"index"`
	Subtitle string `json:"subtitle"`

	PageCount           int `json:"page_count"`
	TotalUnitCount      int `json:"total_unit_count"`
	TranslatedUnitCount int `json:"translated_unit_count"`
	ProofreadUnitCount  int `json:"proofread_unit_count"`

	CreatorId string   `json:"creator_id"`
	Creator   *UserVal `json:"creator,omitempty"`

	UploadedAt *int64 `json:"uploaded_at,omitempty"`

	TranslatingAt *int64 `json:"translating_at,omitempty"`
	TranslatedAt  *int64 `json:"translated_at,omitempty"`

	ProofreadingAt *int64 `json:"proofreading_at,omitempty"`
	ProofreadAt    *int64 `json:"proofread_at,omitempty"`

	TypesettingAt *int64 `json:"typesetting_at,omitempty"`
	TypesetAt     *int64 `json:"typeset_at,omitempty"`

	ReviewedAt  *int64 `json:"reviewed_at,omitempty"`
	PublishedAt *int64 `json:"published_at,omitempty"`

	CreatedAt int64 `json:"created_at"`
	UpdatedAt int64 `json:"updated_at"`
}

// `ListChapterArgs` holds list query arguments for chapter listing.
type ListChapterArgs struct {
	// `ComicId` is target comic identifier.
	ComicId string `url:"comic_id"`

	Includes []enum.ChapterIncl `url:"includes"`

	Offset int `url:"offset"`
	Limit  int `url:"limit"`
}

// `CreateChapterArgs` holds chapter creation input.
type CreateChapterArgs struct {
	// `ComicId` is target comic identifier.
	ComicId string `json:"comic_id"`

	// `Subtitle` is optional subtitle.
	Subtitle *string `json:"subtitle"`
}

// `ChapterCreatedRes` is returned after chapter creation.
type ChapterCreatedRes struct {
	// `Id` is created chapter identifier.
	Id string `json:"id"`
}

// `GetChapterByIdArgs` carries query args for chapter detail API.
type GetChapterByIdArgs struct {
	// `ChapterId` identifies the target chapter.
	ChapterId string `url:"chapter_id"`

	// `Includes` controls relation assembly fields.
	Includes []enum.ChapterIncl `url:"includes"`
}

// `ChapterUpdArgs` holds mutable chapter fields.
// `WorkflowTransition` and `RevertTransition` are mutually exclusive; set at most one
type ChapterUpdArgs struct {
	// `Id` identifies target chapter.
	Id string `json:"id"`

	// `Subtitle` is optional subtitle update.
	Subtitle *string `json:"subtitle"`
	// `IsPinned` optionally updates pinned status.
	IsPinned *bool `json:"is_pinned"`

	// `WorkflowTransition` drives forward workflow timestamp mutation
	WorkflowTransition *enum.WorkflowTransition `json:"workflow_transition"`
	// `RevertTransition` drives a revert of one workflow timestamp back to NULL
	// Mutually exclusive with `WorkflowTransition`
	// Publish-complete cannot be reverted
	RevertTransition *enum.WorkflowTransition `json:"revert_transition"`
}

// `JoinChapterArgs` holds input for chapter joining.
type JoinChapterArgs struct {
	// `ChapterId` identifies target chapter.
	ChapterId string `json:"chapter_id"`

	// `RoleMask` specifies assignment roles to be added.
	RoleMask aggr.RoleMask `json:"role_mask"`
}
