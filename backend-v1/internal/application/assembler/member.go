package assembler

import (
	"labelplus-next-web-be/internal/domain/model"
	"labelplus-next-web-be/internal/util"
	"labelplus-next-web-be/internal/value"
)

func AssembleMemberInfo(memberInfo model.MemberInfo, onLoadURL OnLoadURL) value.MemberInfo {
	result := value.MemberInfo{
		ID:                    memberInfo.ID,
		UserID:                memberInfo.UserID,
		TeamID:                memberInfo.TeamID,
		Roles:                 model.MaskRoles(memberInfo.Roles()),
		AssignedRawProviderAt: util.ToUnixPtr(memberInfo.AssignedRawProviderAt),
		AssignedTranslatorAt:  util.ToUnixPtr(memberInfo.AssignedTranslatorAt),
		AssignedProofreaderAt: util.ToUnixPtr(memberInfo.AssignedProofreaderAt),
		AssignedTypesetterAt:  util.ToUnixPtr(memberInfo.AssignedTypesetterAt),
		AssignedReviewerAt:    util.ToUnixPtr(memberInfo.AssignedReviewerAt),
		AssignedPublisherAt:   util.ToUnixPtr(memberInfo.AssignedPublisherAt),
		AssignedAdminAt:       util.ToUnixPtr(memberInfo.AssignedAdminAt),
		CreatedAt:             memberInfo.CreatedAt.UnixMilli(),
		UpdatedAt:             memberInfo.UpdatedAt.UnixMilli(),
	}

	if memberInfo.User != nil {
		user := AssembleUserInfo(*memberInfo.User, onLoadURL)
		result.User = &user
	}

	if memberInfo.Team != nil {
		team := AssembleTeamInfo(*memberInfo.Team, onLoadURL)
		result.Team = &team
	}

	return result
}
