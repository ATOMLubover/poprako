package assembler

import (
	"labelplus-next-web-be/internal/domain/model"
	"labelplus-next-web-be/internal/util"
	"labelplus-next-web-be/internal/value"
)

// AssembleAssignmentInfo 将 model.AssignmentInfo 转换为 value.AssignmentInfo。
// User 和 Chapter 字段仅在 includes 指定时填充，nil 表示未请求该关联。
func AssembleAssignmentInfo(assignmentInfo model.AssignmentInfo, onLoadURL OnLoadURL) value.AssignmentInfo {
	result := value.AssignmentInfo{
		ID:                    assignmentInfo.ID,
		ChapterID:             assignmentInfo.ChapterID,
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

	if assignmentInfo.User != nil {
		user := AssembleUserInfo(*assignmentInfo.User, onLoadURL)
		result.User = &user
	}

	if assignmentInfo.Chapter != nil {
		chapter := AssembleChapterInfo(*assignmentInfo.Chapter, onLoadURL)
		result.Chapter = &chapter
	}

	return result
}
