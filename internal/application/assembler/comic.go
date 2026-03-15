package assembler

import (
	"labelplus-next-web-be/internal/domain/model"
	"labelplus-next-web-be/internal/value"
)

func AssembleComicInfo(comicInfo model.ComicInfo, onLoadURL OnLoadURL) value.ComicInfo {
	result := value.ComicInfo{
		ID:           comicInfo.ID,
		WorksetID:    comicInfo.WorksetID,
		Index:        comicInfo.Index,
		Title:        comicInfo.Title,
		Author:       comicInfo.Author,
		Description:  comicInfo.Description,
		CoverURL:     comicInfo.CoverURL,
		ChapterCount: comicInfo.ChapterCount,
		CreatorID:    comicInfo.CreatorID,
		LastActiveAt: comicInfo.LastActiveAt.UnixMilli(),
		CreatedAt:    comicInfo.CreatedAt.UnixMilli(),
		UpdatedAt:    comicInfo.UpdatedAt.UnixMilli(),
	}

	if comicInfo.Workset != nil {
		workset := AssembleWorksetInfo(*comicInfo.Workset, onLoadURL)
		result.Workset = &workset
	}

	if comicInfo.Creator != nil {
		creator := AssembleUserInfo(*comicInfo.Creator, onLoadURL)
		result.Creator = &creator
	}

	return result
}
