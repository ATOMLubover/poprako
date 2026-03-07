package service

import (
	"labelplus-next-web-be/internal/domain/model"

	"go.uber.org/zap"
)

type OnLoadUserInfo func(userID string) (model.UserInfo, error)

type OnLoadMemberInfo func(teamID string, userID string) (model.MemberInfo, error)

func CheckInvitationPermission(
	userID string,
	teamID string,
	onLoadMemberInfo OnLoadMemberInfo,
	permission model.Permission,
) bool {
	memberInfo, err := onLoadMemberInfo(teamID, userID)
	if err != nil {
		zap.L().Error(
			"CheckInvitationPermission: 获取成员信息失败",
			zap.String("teamID", teamID),
			zap.String("userID", userID),
			zap.Error(err),
		)

		return false
	}

	switch permission {
	case model.PermissionInvitationList,
		model.PermissionInvitationCreate,
		model.PermissionInvitationDelete,
		model.PermissionInvitationUpdate:
		// 目前仅管理员有权限管理邀请
		return memberInfo.HasAnyRole(model.RoleAdmin)

	default:
		zap.L().Warn(
			"CheckInvitationPermission: 无法识别的权限",
			zap.String("permission", string(permission)),
		)
		return false
	}
}

func CheckUserPermission(
	userID string,
	targetUserID string,
	onLoadUserInfo OnLoadUserInfo,
	permission model.Permission,
) bool {
	userInfo, err := onLoadUserInfo(userID)
	if err != nil {
		zap.L().Error(
			"CheckUserPermission: 获取用户信息失败",
			zap.String("currentUserID", userID),
			zap.Error(err),
		)

		return false
	}

	switch permission {
	case model.PermissionUserRemove:
		// 目前仅超级管理员有权限删除用户，且不能删除自己
		return userInfo.IsSuperAdmin && userID != targetUserID

	case model.PermissionUserView:
		// 目前用户可以查看自己的信息，超级管理员可以查看所有用户的信息
		return userID == targetUserID || userInfo.IsSuperAdmin

	case model.PermissionUserList:
		// 目前仅超级管理员有权限查看用户列表
		return userInfo.IsSuperAdmin

	default:
		zap.L().Warn(
			"CheckUserPermission: 无法识别的权限",
			zap.String("permission", string(permission)),
		)

		return false
	}
}

func CheckTeamPermission(
	userID string,
	teamID string,
	onLoadUserInfo OnLoadUserInfo,
	onLoadMemberInfo OnLoadMemberInfo,
	permission model.Permission,
) bool {
	userInfo, err := onLoadUserInfo(userID)
	if err != nil {
		zap.L().Error(
			"CheckTeamPermission: 获取用户信息失败",
			zap.String("userID", userID),
			zap.Error(err),
		)

		return false
	}

	switch permission {
	case model.PermissionTeamCreate,
		model.PermissionTeamListAll:
		// 只有超级管理员可以创建汉化组和查看所有汉化组
		if !userInfo.IsSuperAdmin {
			return false
		}

		return true

	case model.PermissionTeamUpdate,
		model.PermissionTeamDelete:
		// 仅超级管理员、汉化组管理员有权限更新或删除汉化组
		if userInfo.IsSuperAdmin {
			return true
		}

		memberInfo, err := onLoadMemberInfo(teamID, userID)
		if err != nil {
			zap.L().Error(
				"CheckTeamPermission: 获取成员信息失败",
				zap.String("teamID", teamID),
				zap.String("userID", userID),
				zap.Error(err),
			)

			return false
		}

		return memberInfo.HasAnyRole(model.RoleAdmin)

	default:
		zap.L().Warn(
			"CheckTeamPermission: 无法识别的权限",
			zap.String("permission", string(permission)),
		)

		return false
	}
}

func CheckMemberPermission(
	userID string,
	teamID string,
	onLoadUserInfo OnLoadUserInfo,
	onLoadMemberInfo OnLoadMemberInfo,
	permission model.Permission,
) bool {
	userInfo, err := onLoadUserInfo(userID)
	if err != nil {
		zap.L().Error(
			"CheckMemberPermission: 获取用户信息失败",
			zap.String("userID", userID),
			zap.Error(err),
		)

		return false
	}

	memberInfo, err := onLoadMemberInfo(teamID, userID)
	if err != nil {
		zap.L().Error(
			"CheckMemberPermission: 获取成员信息失败",
			zap.String("teamID", teamID),
			zap.String("userID", userID),
			zap.Error(err),
		)

		return false
	}

	switch permission {
	case model.PermissionMemberCreate:
		// 仅超级管理员有权限 **直接** 添加成员
		return userInfo.IsSuperAdmin
	//
	case model.PermissionMemberList,
		model.PermissionMemberUpdate,
		model.PermissionMemberDelete:
		// 仅管理员有权限管理成员
		return memberInfo.HasAnyRole(model.RoleAdmin)

	default:
		zap.L().Warn(
			"CheckMemberPermission: 无法识别的权限",
			zap.String("permission", string(permission)),
		)

		return false
	}
}

func CheckComicPermission(
	userID string,
	teamID string,
	onLoadMemberInfo OnLoadMemberInfo,
	permission model.Permission,
) bool {
	memberInfo, err := onLoadMemberInfo(teamID, userID)
	if err != nil {
		zap.L().Error(
			"CheckComicPermission: 获取成员信息失败",
			zap.String("teamID", teamID),
			zap.String("userID", userID),
			zap.Error(err),
		)

		return false
	}

	switch permission {
	case model.PermissionComicList:
		// 所有成员都可以查看漫画列表
		return true

	case model.PermissionComicCreate,
		model.PermissionComicUpdate,
		model.PermissionComicDelete:
		// 目前仅管理员有权限管理漫画
		return memberInfo.HasAnyRole(model.RoleAdmin)

	default:
		zap.L().Warn(
			"CheckComicPermission: 无法识别的权限",
			zap.String("permission", string(permission)),
		)

		return false
	}
}

func GetPermissionType(permission model.Permission) model.Permission {
	switch {
	case hasPrefix(permission, model.PrefixPermissionInvitation):
		return model.PrefixPermissionInvitation

	case hasPrefix(permission, model.PrefixPermissionMember):
		return model.PrefixPermissionMember

	case hasPrefix(permission, model.PrefixPermissionComic):
		return model.PrefixPermissionComic

	case hasPrefix(permission, model.PrefixPermissionTeam):
		return model.PrefixPermissionTeam

	default:
		zap.L().Warn(
			"GetPermissionType: 无法识别的权限",
			zap.String("permission", string(permission)),
		)
		return ""
	}
}

func hasPrefix(permission model.Permission, prefix model.Permission) bool {
	return len(permission) >= len(prefix) && permission[:len(prefix)] == prefix
}
