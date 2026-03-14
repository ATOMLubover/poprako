package assembler

import (
	"labelplus-next-web-be/internal/domain/model"
	"labelplus-next-web-be/internal/util"
	"labelplus-next-web-be/internal/value"
)

func AssembleAssignmentInfoFromUser(assignmentInfo model.AssignmentWithUserInfo, onLoadURL OnLoadURL) value.AssignmentInfo {
	result := value.AssignmentInfo{
		ID:                    assignmentInfo.ID,
		ChapterID:             assignmentInfo.ChapterID,
		UserID:                assignmentInfo.User.ID,
		AssignedRawProviderAt: util.ToUnixPtr(assignmentInfo.AssignedRawProviderAt),
		AssignedTranslatorAt:  util.ToUnixPtr(assignmentInfo.AssignedTranslatorAt),
		AssignedProofreaderAt: util.ToUnixPtr(assignmentInfo.AssignedProofreaderAt),
		AssignedTypesetterAt:  util.ToUnixPtr(assignmentInfo.AssignedTypesetterAt),
		AssignedRedrawerAt:    util.ToUnixPtr(assignmentInfo.AssignedRedrawerAt),
		AssignedReviewerAt:    util.ToUnixPtr(assignmentInfo.AssignedReviewerAt),
		AssignedPublisherAt:   util.ToUnixPtr(assignmentInfo.AssignedPublisherAt),
		CreatedAt:             assignmentInfo.CreatedAt.UnixMilli(),
		UpdatedAt:             assignmentInfo.UpdatedAt.UnixMilli(),
	}

	user := AssembleUserInfo(assignmentInfo.User, onLoadURL)
	result.User = &user

	return result
}

func AssembleAssignmentInfoFromChapter(assignmentInfo model.AssignmentWithChapterInfo, onLoadURL OnLoadURL) value.AssignmentInfo {
	result := value.AssignmentInfo{
		ID:                    assignmentInfo.ID,
		ChapterID:             assignmentInfo.Chapter.ID,
		UserID:                assignmentInfo.UserID,
		AssignedRawProviderAt: util.ToUnixPtr(assignmentInfo.AssignedRawProviderAt),
		AssignedTranslatorAt:  util.ToUnixPtr(assignmentInfo.AssignedTranslatorAt),
		AssignedProofreaderAt: util.ToUnixPtr(assignmentInfo.AssignedProofreaderAt),
		AssignedTypesetterAt:  util.ToUnixPtr(assignmentInfo.AssignedTypesetterAt),
		AssignedRedrawerAt:    util.ToUnixPtr(assignmentInfo.AssignedRedrawerAt),
		AssignedReviewerAt:    util.ToUnixPtr(assignmentInfo.AssignedReviewerAt),
		AssignedPublisherAt:   util.ToUnixPtr(assignmentInfo.AssignedPublisherAt),
		CreatedAt:             assignmentInfo.CreatedAt.UnixMilli(),
		UpdatedAt:             assignmentInfo.UpdatedAt.UnixMilli(),
	}

	chapter := value.ChapterInfo{
		ID:                  assignmentInfo.Chapter.ID,
		ComicID:             assignmentInfo.Chapter.Comic.ID,
		Index:               assignmentInfo.Chapter.Index,
		ChapterNo:           assignmentInfo.Chapter.ChapterNo,
		PageCount:           assignmentInfo.Chapter.PageCount,
		TotalUnitCount:      assignmentInfo.Chapter.TotalUnitCount,
		TranslatedUnitCount: assignmentInfo.Chapter.TranslatedUnitCount,
		ProofreadUnitCount:  assignmentInfo.Chapter.ProofreadUnitCount,
		CreatedAt:           assignmentInfo.Chapter.CreatedAt.UnixMilli(),
		UpdatedAt:           assignmentInfo.Chapter.UpdatedAt.UnixMilli(),
	}

	comic := AssembleComicInfo(assignmentInfo.Chapter.Comic, onLoadURL)
	chapter.Comic = &comic

	result.Chapter = &chapter

	return result
}
