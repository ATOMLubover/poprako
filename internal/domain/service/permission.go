package service

import (
	"labelplus-next-web-be/internal/domain/model"

	"go.uber.org/zap"
)

func CheckInvitationPermission(
	targetTeamID string,
	currentUserMemberships []model.MemberProfile,
	permission model.Permission,
) bool {
	var targetMemberInfo *model.MemberProfile

	for i := range currentUserMemberships {
		if currentUserMemberships[i].TeamID == targetTeamID {
			targetMemberInfo = &currentUserMemberships[i]
			break
		}
	}

	if targetMemberInfo == nil {
		zap.L().Warn(
			"checkInvitationPermission: 未找到目标汉化组下的成员信息",
			zap.String("targetTeamID", targetTeamID),
		)
		return false
	}

	switch permission {
	case model.PermissionInvitationList,
		model.PermissionInvitationCreate,
		model.PermissionInvitationDelete,
		model.PermissionInvitationUpdate:
		// 目前仅管理员有权限管理邀请
		return targetMemberInfo.HasAnyRole(model.RoleAdmin)

	default:
		zap.L().Warn(
			"CheckInvitationPermission: 无法识别的权限",
			zap.String("permission", string(permission)),
		)
		return false
	}
}

func CheckUserPermission(
	currentUser *model.UserInfo,
	targetUserID string,
	permission model.Permission,
) bool {
	if currentUser == nil {
		zap.L().Warn(
			"checkUserPermission: currentUser 为空",
			zap.String("targetUserID", targetUserID),
			zap.String("permission", string(permission)),
		)

		return false
	}

	switch permission {
	case model.PermissionUserRemove:
		// 目前仅超级管理员有权限删除用户，且不能删除自己
		return currentUser.ID != targetUserID &&
			currentUser.IsSuperAdmin

	case model.PermissionUserView:
		// 目前用户可以查看自己的信息，超级管理员可以查看所有用户的信息
		return currentUser.ID == targetUserID ||
			currentUser.IsSuperAdmin

	case model.PermissionUserList:
		// 目前仅超级管理员有权限查看用户列表
		return currentUser.IsSuperAdmin

	default:
		zap.L().Warn(
			"CheckUserPermission: 无法识别的权限",
			zap.String("permission", string(permission)),
		)

		return false
	}
}

func CheckTeamPermission(
	targetTeamID string,
	currentUser *model.UserInfo,
	currentUserMemberships []model.MemberProfile,
	permission model.Permission,
) bool {
	if currentUser == nil {
		zap.L().Warn(
			"checkTeamPermission: currentUser 为空",
			zap.String("targetTeamID", targetTeamID),
			zap.String("permission", string(permission)),
		)

		return false
	}

	switch permission {
	case model.PermissionTeamCreate,
		model.PermissionTeamListAll:
		// 只有超级管理员可以创建汉化组和查看所有汉化组
		if !currentUser.IsSuperAdmin {
			return false
		}

		return true

	case model.PermissionTeamUpdate,
		model.PermissionTeamDelete:
		// 目前仅超级管理员、汉化组管理员有权限更新或删除汉化组
		if currentUser.IsSuperAdmin {
			return true
		}

		var targetMemberInfo *model.MemberProfile

		for i := range currentUserMemberships {
			if currentUserMemberships[i].TeamID == targetTeamID {
				targetMemberInfo = &currentUserMemberships[i]
				break
			}
		}

		if targetMemberInfo == nil {
			zap.L().Warn(
				"checkTeamPermission: 未找到目标汉化组下的成员信息",
				zap.String("targetTeamID", targetTeamID),
			)
			return false
		}

		return targetMemberInfo.HasAnyRole(model.RoleAdmin)

	default:
		zap.L().Warn(
			"CheckTeamPermission: 无法识别的权限",
			zap.String("permission", string(permission)),
		)
		return false
	}
}

func CheckMemberPermission(
	targetTeamID string,
	currentUser *model.UserInfo,
	currentUserMemberships []model.MemberProfile,
	permission model.Permission,
) bool {
	switch permission {
	case model.PermissionMemberCreate:
		// 仅有超级管理员才可以直接添加成员
		if currentUser == nil {
			zap.L().Warn("checkMemberPermission: currentUser 为空")
			return false
		}

		return currentUser.IsSuperAdmin
	}

	var targetMemberInfo *model.MemberProfile

	for i := range currentUserMemberships {
		if currentUserMemberships[i].TeamID == targetTeamID {
			targetMemberInfo = &currentUserMemberships[i]
			break
		}
	}

	if targetMemberInfo == nil {
		zap.L().Warn(
			"checkMemberPermission: 未找到目标汉化组下的成员信息",
			zap.String("targetTeamID", targetTeamID),
		)
		return false
	}

	switch permission {
	case model.PermissionMemberList,
		model.PermissionMemberUpdate,
		model.PermissionMemberDelete:
		// 目前仅管理员有权限管理成员
		return targetMemberInfo.HasAnyRole(model.RoleAdmin)

	default:
		zap.L().Warn(
			"CheckMemberPermission: 无法识别的权限",
			zap.String("permission", string(permission)),
		)
		return false
	}
}
