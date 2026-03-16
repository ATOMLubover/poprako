package assembler

import (
	"labelplus-next-web-be/internal/domain/model"
	"labelplus-next-web-be/internal/value"
)

// AssembleUnitInfo 将 model.UnitInfo 转换为 value.UnitInfo。
func AssembleUnitInfo(unitInfo model.UnitInfo) value.UnitInfo {
	return value.UnitInfo{
		ID:                 unitInfo.ID,
		PageID:             unitInfo.PageID,
		Index:              unitInfo.Index,
		XCoord:             unitInfo.XCoord,
		YCoord:             unitInfo.YCoord,
		IsBubble:           unitInfo.IsBubble,
		TranslatedText:     unitInfo.TranslatedText,
		TranslatorID:       unitInfo.TranslatorID,
		TranslatorComment:  unitInfo.TranslatorComment,
		IsProofread:        unitInfo.IsProofread,
		ProofreadText:      unitInfo.ProofreadText,
		ProofreaderID:      unitInfo.ProofreaderID,
		ProofreaderComment: unitInfo.ProofreaderComment,
	}
}
