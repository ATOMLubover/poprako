package val

import "poprako-s/internal/domain/model/enum"

// `ChapterVal` is the app-facing chapter value object.
type ChapterVal struct {
	Id string `json:"id"`

	ComicId string `json:"comic_id"`

	IsPinned bool `json:"is_pinned"`

	Index    int    `json:"index"`
	Subtitle string `json:"subtitle"`

	PageCount           int `json:"page_count"`
	TotalUnitCount      int `json:"total_unit_count"`
	TranslatedUnitCount int `json:"translated_unit_count"`
	ProofreadUnitCount  int `json:"proofread_unit_count"`

	CreatorId string `json:"creator_id"`

	UploadedAt *int64 `json:"uploaded_at"`

	TransalatingAt *int64 `json:"transalating_at"`
	TranslatedAt   *int64 `json:"translated_at"`

	ProofreadingAt *int64 `json:"proofreading_at"`
	ProofreadAt    *int64 `json:"proofread_at"`

	TypesettingAt *int64 `json:"typesetting_at"`
	TypesetAt     *int64 `json:"typeset_at"`

	ReviewedAt  *int64 `json:"reviewed_at"`
	PublishedAt *int64 `json:"published_at"`

	CreatedAt int64 `json:"created_at"`
	UpdatedAt int64 `json:"updated_at"`
}

// `ListChapterArgs` holds list query arguments for chapter listing.
type ListChapterArgs struct {
	// `ComicId` is target comic identifier.
	ComicId string `url:"comic_id"`

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

// `ChapterUpdArgs` holds mutable chapter fields.
type ChapterUpdArgs struct {
	// `Id` identifies target chapter.
	Id string `json:"id"`

	// `Subtitle` is optional subtitle update.
	Subtitle *string `json:"subtitle"`
	// `IsPinned` optionally updates pinned status.
	IsPinned *bool `json:"is_pinned"`

	// `WorkflowTransition` drives workflow timestamp mutation.
	WorkflowTransition *enum.WorkflowTransition `json:"workflow_transition"`
}
