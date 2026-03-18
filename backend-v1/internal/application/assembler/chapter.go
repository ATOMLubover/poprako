package assembler

import (
	"labelplus-next-web-be/internal/domain/model"
	"labelplus-next-web-be/internal/util"
	"labelplus-next-web-be/internal/value"
)

func AssembleChapterInfo(chapterInfo model.ChapterInfo, onLoadURL OnLoadURL) value.ChapterInfo {
	result := value.ChapterInfo{
		ID:                  chapterInfo.ID,
		ComicID:             chapterInfo.ComicID,
		Index:               chapterInfo.Index,
		Subtitle:            chapterInfo.Subtitle,
		PageCount:           chapterInfo.PageCount,
		TotalUnitCount:      chapterInfo.TotalUnitCount,
		TranslatedUnitCount: chapterInfo.TranslatedUnitCount,
		ProofreadUnitCount:  chapterInfo.ProofreadUnitCount,
		UploadedAt:          util.ToUnixPtr(chapterInfo.UploadedAt),
		TransalatingAt:      util.ToUnixPtr(chapterInfo.TransalatingAt),
		TranslatedAt:        util.ToUnixPtr(chapterInfo.TranslatedAt),
		ProofreadingAt:      util.ToUnixPtr(chapterInfo.ProofreadingAt),
		ProofreadAt:         util.ToUnixPtr(chapterInfo.ProofreadAt),
		TypesettingAt:       util.ToUnixPtr(chapterInfo.TypesettingAt),
		TypesetAt:           util.ToUnixPtr(chapterInfo.TypesetAt),
		ReviewedAt:          util.ToUnixPtr(chapterInfo.ReviewedAt),
		PublishedAt:         util.ToUnixPtr(chapterInfo.PublishedAt),
		CreatorID:           chapterInfo.CreatorID,
		CreatedAt:           chapterInfo.CreatedAt.UnixMilli(),
		UpdatedAt:           chapterInfo.UpdatedAt.UnixMilli(),
	}

	if chapterInfo.Creator != nil {
		creator := AssembleUserInfo(*chapterInfo.Creator, onLoadURL)
		result.Creator = &creator
	}

	return result
}
