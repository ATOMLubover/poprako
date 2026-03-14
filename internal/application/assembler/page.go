package assembler

import (
	"labelplus-next-web-be/internal/domain/model"
	"labelplus-next-web-be/internal/value"
)

func AssemblePageInfo(pageInfo model.PageInfo, onLoadURL OnLoadURL) value.PageInfo {
	imageURL, err := onLoadURL(pageInfo.OSSKey)
	if err != nil {
		imageURL = ""
	}

	result := value.PageInfo{
		ID:                  pageInfo.ID,
		ChapterID:           pageInfo.ChapterID,
		Index:               pageInfo.Index,
		CreatorID:           pageInfo.CreatorID,
		ImageURL:            imageURL,
		TotalUnitCount:      pageInfo.TotalUnitCount,
		TranslatedUnitCount: pageInfo.TranslatedUnitCount,
		ProofreadUnitCount:  pageInfo.ProofreadUnitCount,
		CreatedAt:           pageInfo.CreatedAt.Unix(),
		UpdatedAt:           pageInfo.UpdatedAt.Unix(),
	}

	if pageInfo.Creator != nil {
		creator := AssembleUserInfo(*pageInfo.Creator, onLoadURL)
		result.Creator = &creator
	}

	return result
}
