package model

import (
	"labelplus-next-web-be/internal/util"
)

type UnitInfo struct {
	ID string

	PageID string
	Index  int

	XCoord int
	YCoord int

	IsBubble bool

	TranslatedText    *string
	TranslatorID      *string
	TranslatorComment *string

	IsProofread        bool
	ProofreadText      *string
	ProofreaderID      *string
	ProofreaderComment *string

	// CreatedAt time.Time --- IGNORE ---
	// UpdatedAt time.Time --- IGNORE ---
}

func NewUnitInfo(
	id string,
	index int,
	xCoord int,
	yCoord int,
	isBubble bool,
	translatedText *string,
	translatorID *string,
	translatorComment *string,
	isProofread bool,
	proofreadText *string,
	proofreaderID *string,
	proofreaderComment *string,
) UnitInfo {
	return UnitInfo{
		ID:                 id,
		Index:              index,
		XCoord:             xCoord,
		YCoord:             yCoord,
		IsBubble:           isBubble,
		TranslatedText:     translatedText,
		TranslatorID:       translatorID,
		TranslatorComment:  translatorComment,
		IsProofread:        isProofread,
		ProofreadText:      proofreadText,
		ProofreaderID:      proofreaderID,
		ProofreaderComment: proofreaderComment,
	}
}

// 纯粹 patch 语义的 struct，用于减少非必要字段修改
type UnitPatch struct {
	// 这个 ID 是必须的，因为我们需要知道修改哪个单元
	ID string

	// 注意，在下述的字段中，如果 Option 为 None，代表不修改该字段，Some 代表 **必须修改**
	// 如果模板参数是指针，则 Some(nil) 代表清空为 NULL

	Index util.Option[int]

	XCoord util.Option[int]
	YCoord util.Option[int]

	IsBubble util.Option[bool]

	TranslatedText    util.Option[*string]
	TranslatorID      util.Option[*string]
	TranslatorComment util.Option[*string]

	IsProofread        util.Option[bool]
	ProofreadText      util.Option[*string]
	ProofreaderID      util.Option[*string]
	ProofreaderComment util.Option[*string]
}

func NewUnitPatch(
	id string,
	index util.Option[int],
	xCoord util.Option[int],
	yCoord util.Option[int],
	isBubble util.Option[bool],
	translatedText util.Option[*string],
	translatorID util.Option[*string],
	translatorComment util.Option[*string],
	isProofread util.Option[bool],
	proofreadText util.Option[*string],
	proofreaderID util.Option[*string],
	proofreaderComment util.Option[*string],
) UnitPatch {
	return UnitPatch{
		ID:                 id,
		Index:              index,
		XCoord:             xCoord,
		YCoord:             yCoord,
		IsBubble:           isBubble,
		TranslatedText:     translatedText,
		TranslatorID:       translatorID,
		TranslatorComment:  translatorComment,
		IsProofread:        isProofread,
		ProofreadText:      proofreadText,
		ProofreaderID:      proofreaderID,
		ProofreaderComment: proofreaderComment,
	}
}

// 考虑到 UnitCreation 和 UnitInfo 的字段完全相同，且 UnitCreation 语义上就是创建一个 UnitInfo
// 所以直接 type alias，避免不必要的代码重复
type UnitCreation = UnitInfo

func NewUnitCreation(
	id string,
	index int,
	xCoord int,
	yCoord int,
	isBubble bool,
	translatedText *string,
	translatorID *string,
	translatorComment *string,
	isProofread bool,
	proofreadText *string,
	proofreaderID *string,
	proofreaderComment *string,
) UnitCreation {
	return UnitCreation{
		ID:                 id,
		Index:              index,
		XCoord:             xCoord,
		YCoord:             yCoord,
		IsBubble:           isBubble,
		TranslatedText:     translatedText,
		TranslatorID:       translatorID,
		TranslatorComment:  translatorComment,
		IsProofread:        isProofread,
		ProofreadText:      proofreadText,
		ProofreaderID:      proofreaderID,
		ProofreaderComment: proofreaderComment,
	}
}
