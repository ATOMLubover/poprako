package assembler

import (
	"labelplus-next-web-be/internal/domain/model"
	"labelplus-next-web-be/internal/value"

	"go.uber.org/zap"
)

func AssembleUserInfo(userInfo model.UserInfo, onLoadURL OnLoadURL) value.UserInfo {
	avatarURL, err := onLoadURL(userInfo.AvatarOSSKey)
	if err != nil {
		// 这里不应该因为加载头像失败而导致整个用户信息加载失败，所以我们捕获错误并记录日志，但不返回错误
		zap.L().Error(
			"加载用户头像 URL 失败",
			zap.String("user_id", userInfo.ID),
			zap.String("avatar_oss_key", userInfo.AvatarOSSKey),
			zap.Error(err),
		)
		avatarURL = "" // 如果加载失败，我们可以选择返回一个空字符串或者一个默认的头像 URL
	}

	return value.UserInfo{
		ID:               userInfo.ID,
		Name:             userInfo.Name,
		QQ:               userInfo.QQ,
		AvatarURL:        avatarURL,
		IsAvatarUploaded: userInfo.IsAvatarUploaded,
		IsSuperAdmin:     userInfo.IsSuperAdmin,
		CreatedAt:        userInfo.CreatedAt.UnixMilli(),
		UpdatedAt:        userInfo.UpdatedAt.UnixMilli(),
	}
}

func AssembleUserStatsInfo(userStats model.UserStats) value.UserStatsInfo {
	return value.UserStatsInfo{
		UserID:                  userStats.UserID,
		TotalAssignmentCount:    userStats.TotalAssignmentCount,
		ActiveAssignmentCount:   userStats.ActiveAssignmentCount,
		FinishedAssignmentCount: userStats.FinishedAssignmentCount,
	}
}
