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
	Units []UnitVal `json:"units"`

	TotalUnitCount      int `json:"total_unit_count"`
	TranslatedUnitCount int `json:"translated_unit_count"`
	ProofreadUnitCount  int `json:"proofread_unit_count"`
}

// `UnitDiffVal` is the transport-safe representation of `UnitDiff`.
type UnitDiffVal struct {
	PageId string `json:"page_id"`

	Ops []UnitOpVal `json:"ops"`

	CandOrder []string `json:"cand_order"`
}

// `UnitOpVal` is the transport-safe representation of one unit op.
type UnitOpVal struct {
	Id      string `json:"id,omitempty"`
	LocalId string `json:"local_id,omitempty"`

	IsBubble    *bool `json:"is_bubble,omitempty"`
	IsProofread *bool `json:"is_proofread,omitempty"`

	XCoord *float64 `json:"x_coord,omitempty"`
	YCoord *float64 `json:"y_coord,omitempty"`

	TranslatedText    *string `json:"translated_text,omitempty"`
	TranslatorComment *string `json:"translator_comment,omitempty"`
	LastTranslatorId  *string `json:"last_translator_id,omitempty"`

	ProofreadText      *string `json:"proofread_text,omitempty"`
	ProofreaderComment *string `json:"proofreader_comment,omitempty"`
	LastProofreaderId  *string `json:"last_proofreader_id,omitempty"`
}

// `SavePageUnitsArgs` holds input for saving page units.
type SavePageUnitsArgs struct {
	PageId string `json:"page_id"`

	Diff *UnitDiffVal `json:"diff"`
}

// `SavePageUnitsRes` holds synchronized unit counts after save.
type SavePageUnitsRes struct {
	TotalUnitCount      int `json:"total_unit_count"`
	TranslatedUnitCount int `json:"translated_unit_count"`
	ProofreadUnitCount  int `json:"proofread_unit_count"`
}
