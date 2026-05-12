package val

// `UnitVal` is the app-facing value object for unit data.
type UnitVal struct {
	Id string `json:"id"`

	PageId string `json:"page_id"`

	Index int `json:"index"`

	IsBubble bool `json:"is_bubble"`

	IsProofread bool `json:"is_proofread"`

	XCoord float64 `json:"x_coord"`
	YCoord float64 `json:"y_coord"`

	TranslatedText    *string `json:"translated_text,omitempty"`
	TranslatorComment *string `json:"translator_comment,omitempty"`
	LastTranslatorId  *string `json:"last_translator_id,omitempty"`

	ProofreadText      *string `json:"proofread_text,omitempty"`
	ProofreaderComment *string `json:"proofreader_comment,omitempty"`
	LastProofreaderId  *string `json:"last_proofreader_id,omitempty"`

	CreatedAt int64 `json:"created_at"`
	UpdatedAt int64 `json:"updated_at"`
}

// `ListPageUnitsArgs` carries query args for page unit list.
type ListPageUnitsArgs struct {
	PageId string `url:"page_id"`
}

// `ListPageUnitsRes` is the app response for one page unit list.
type ListPageUnitsRes struct {
	Units []UnitVal `json:"units,omitempty"`

	TotalUnitCount      int `json:"total_unit_count"`
	TranslatedUnitCount int `json:"translated_unit_count"`
	ProofreadUnitCount  int `json:"proofread_unit_count"`
}

// `UnitDiffVal` is the transport-safe representation of `UnitDiff`
// It describes a batch of unit mutations to apply on one page
//
// `Ops` is a list of unit operations (create, update or delete)
// applied in the given order before reindexing
//
// `CandOrder` is the client-suggested unit ordering expressed as an
// ordered list of unit identifiers  Use `local_id` values from create
// ops or real `id` values from save ops to place new and modified units
// relative to existing ones  Every `id` or `local_id` that appears in a
// create or save op must be listed in `CandOrder`; ids being deleted
// must be excluded  `CandOrder` must not contain duplicates
type UnitDiffVal struct {
	// `PageId` identifies the target page and must match the path parameter
	PageId string `json:"page_id"`

	// `Ops` is a list of unit operations applied before reindexing
	Ops []UnitOpVal `json:"operations,omitempty"`

	// `CandOrder` is the client-suggested unit ordering for reindexing
	// Use `local_id` for create ops or real `id` for save ops
	CandOrder []string `json:"candidate_order,omitempty"`
}

// `UnitOpVal` is the transport-safe representation of one unit op
// The op type is inferred from the presence of `id` and `local_id`:
//
//   - CREATE: `local_id` is set AND `id` is empty
//     Requires all of `is_bubble`, `is_proofread`, `x_coord`, `y_coord`
//   - SAVE:   `id` is set AND at least one mutable field is present
//     (is_bubble, is_proofread, x_coord, y_coord, translated_text,
//     translator_comment, last_translator_id, proofread_text,
//     proofreader_comment, last_proofreader_id)
//     Requires `is_bubble`, `is_proofread`, `x_coord`, `y_coord`
//   - DELETE: `id` is set AND NO mutable fields are present
//
// A single op must not combine `id` and `local_id`
type UnitOpVal struct {
	// `Id` is the server-assigned unit identifier (set for SAVE and DELETE ops)
	Id string `json:"id,omitempty"`
	// `LocalId` is a client-generated temporary identifier (set for CREATE ops)
	LocalId string `json:"local_id,omitempty"`

	// `IsBubble` indicates whether the unit is a speech bubble (required for CREATE and SAVE)
	IsBubble *bool `json:"is_bubble,omitempty"`
	// `IsProofread` indicates whether the unit has been proofread (required for CREATE and SAVE)
	IsProofread *bool `json:"is_proofread,omitempty"`

	// `XCoord` is the horizontal coordinate (required for CREATE and SAVE)
	XCoord *float64 `json:"x_coord,omitempty"`
	// `YCoord` is the vertical coordinate (required for CREATE and SAVE)
	YCoord *float64 `json:"y_coord,omitempty"`

	// `TranslatedText` is the translated content (mutable field for SAVE)
	TranslatedText *string `json:"translated_text,omitempty"`
	// `TranslatorComment` is the translator's note (mutable field for SAVE)
	TranslatorComment *string `json:"translator_comment,omitempty"`
	// `LastTranslatorId` is the user id of the last translator (mutable field for SAVE)
	LastTranslatorId *string `json:"last_translator_id,omitempty"`

	// `ProofreadText` is the proofread content (mutable field for SAVE)
	ProofreadText *string `json:"proofread_text,omitempty"`
	// `ProofreaderComment` is the proofreader's note (mutable field for SAVE)
	ProofreaderComment *string `json:"proofreader_comment,omitempty"`
	// `LastProofreaderId` is the user id of the last proofreader (mutable field for SAVE)
	LastProofreaderId *string `json:"last_proofreader_id,omitempty"`
}

// `SavePageUnitsArgs` holds input for saving page units
// `PageId` identifies the target page and must match the path parameter
// `Diff` contains the batch of unit mutations and the suggested ordering
// See `UnitDiffVal` and `UnitOpVal` for the detailed operation format
type SavePageUnitsArgs struct {
	PageId string `json:"page_id"`

	Diff *UnitDiffVal `json:"difference"`
}

// `SavePageUnitsRes` holds synchronized unit counts after save
// `TotalUnitCount` is the total number of units on the page after applying the diff
// `TranslatedUnitCount` is the number of units with non-empty `translated_text`
// `ProofreadUnitCount` is the number of units marked as proofread
type SavePageUnitsRes struct {
	TotalUnitCount      int `json:"total_unit_count"`
	TranslatedUnitCount int `json:"translated_unit_count"`
	ProofreadUnitCount  int `json:"proofread_unit_count"`
}
