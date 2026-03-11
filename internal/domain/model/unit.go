package model

import "time"

type UnitInfo struct {
	ID string

	Index int

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

	CreatedAt time.Time
	UpdatedAt time.Time
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
	createdAt time.Time,
	updatedAt time.Time,
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
		CreatedAt:          createdAt,
		UpdatedAt:          updatedAt,
	}
}
