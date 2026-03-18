package value

import "labelplus-next-web-be/internal/util"

type UnitInfo struct {
	ID string `json:"id"`

	PageID string `json:"page_id"`
	Index  int    `json:"index"`

	XCoord int `json:"x_coord"`
	YCoord int `json:"y_coord"`

	IsBubble bool `json:"is_bubble"`

	TranslatedText    *string `json:"translated_text,omitempty"`
	TranslatorID      *string `json:"translator_id,omitempty"`
	TranslatorComment *string `json:"translator_comment,omitempty"`

	IsProofread        bool    `json:"is_proofread"`
	ProofreadText      *string `json:"proofread_text,omitempty"`
	ProofreaderID      *string `json:"proofreader_id,omitempty"`
	ProofreaderComment *string `json:"proofreader_comment,omitempty"`
}

type UnitPatch struct {
	ID string `json:"id"`

	PageID string `json:"page_id"`
	Index  int    `json:"index"`

	XCoord int `json:"x_coord"`
	YCoord int `json:"y_coord"`

	IsBubble bool `json:"is_bubble"`

	TranslatedText    util.Option[*string] `json:"translated_text,omitempty"`
	TranslatorID      util.Option[*string] `json:"translator_id,omitempty"`
	TranslatorComment util.Option[*string] `json:"translator_comment,omitempty"`

	IsProofread        util.Option[bool]    `json:"is_proofread"`
	ProofreadText      util.Option[*string] `json:"proofread_text,omitempty"`
	ProofreaderID      util.Option[*string] `json:"proofreader_id,omitempty"`
	ProofreaderComment util.Option[*string] `json:"proofreader_comment,omitempty"`
}

type UnitCreation = UnitInfo

type UnitDiff struct {
	Insert []UnitCreation `json:"insert"`
	Patch  []UnitPatch    `json:"patch"`
	Delete []string       `json:"delete"` // 仅包含 ID 列表
}

type SavePageUnitArgs struct {
	PageID   string   `json:"page_id"`
	UnitDiff UnitDiff `json:"unit_diff"`
}
