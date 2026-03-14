package assembler

import (
	"labelplus-next-web-be/internal/domain/model"
	"labelplus-next-web-be/internal/value"

	"go.uber.org/zap"
)

func AssembleTeamInfo(teamInfo model.TeamInfo, onLoadURL OnLoadURL) value.TeamInfo {
	avatarURL, err := onLoadURL(teamInfo.AvatarOSSKey)
	if err != nil {
		zap.L().Error(
			"加载汉化组头像 URL 失败",
			zap.String("team_id", teamInfo.ID),
			zap.String("avatar_oss_key", teamInfo.AvatarOSSKey),
			zap.Error(err),
		)
		avatarURL = ""
	}

	return value.TeamInfo{
		ID:               teamInfo.ID,
		Name:             teamInfo.Name,
		Description:      teamInfo.Description,
		AvatarURL:        avatarURL,
		IsAvatarUploaded: teamInfo.IsAvatarUploaded,
		CreatedAt:        teamInfo.CreatedAt.UnixMilli(),
		UpdatedAt:        teamInfo.UpdatedAt.UnixMilli(),
	}
}
