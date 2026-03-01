package service

import (
	"strings"

	"labelplus-next-web-be/internal/domain/model"

	"go.uber.org/zap"
)

type PermissionService interface {
	CheckPermission(userInfo *model.UserInfo, permission model.Permission) bool
}

type permissionService struct{}

func NewPermissionService() PermissionService {
	return &permissionService{}
}

func (*permissionService) CheckPermission(userInfo *model.UserInfo, permission model.Permission) bool {
	if userInfo == nil {
		return false
	}

	switch extractPermissionPrefix(permission) {
	case model.PrefixPermissionInvitation:
		return checkInvitationPermission(userInfo, permission)
	case model.PrefixPermissionUsers:
		return checkUsersPermission(userInfo, permission)
	default:
		zap.L().Warn(
			"checkPermission: 无法识别的权限前缀",
			zap.String("permission", string(permission)),
		)
		return false
	}
}

func extractPermissionPrefix(permission model.Permission) model.Permission {
	if permission == "" {
		zap.L().Warn(
			"extractPermissionPrefix: permission 长度小于 prefix 长度，无法比较",
			zap.String("permission", string(permission)),
		)
		return ""
	}

	if strings.HasPrefix(string(permission), string(model.PrefixPermissionInvitation)) {
		return model.PrefixPermissionInvitation
	}

	if strings.HasPrefix(string(permission), string(model.PrefixPermissionUsers)) {
		return model.PrefixPermissionUsers
	}

	// 无匹配，返回空字符串
	return ""
}

func checkInvitationPermission(userInfo *model.UserInfo, permission model.Permission) bool {
	if userInfo == nil {
		zap.L().Warn("checkInvitationPermission: userInfo 为空")
		return false
	}

	if !strings.HasPrefix(string(permission), string(model.PrefixPermissionInvitation)) {
		zap.L().Warn(
			"checkInvitationPermission: permission 前缀不匹配",
			zap.String("permission", string(permission)),
		)
		return false
	}

	// 查看、创建和删除邀请，均需要超级管理员权限
	return userInfo.HasRole(model.RoleSuperAdmin)
}

func checkUsersPermission(userInfo *model.UserInfo, permission model.Permission) bool {
	if userInfo == nil {
		zap.L().Warn("checkUsersPermission: userInfo 为空")
		return false
	}

	switch permission {
	case model.PermissionUsersRemove:
		// 删除用户，仅超级管理员可操作
		return userInfo.HasRole(model.RoleSuperAdmin)
	default:
		zap.L().Warn(
			"checkUsersPermission: 无法识别的权限",
			zap.String("permission", string(permission)),
		)
		return false
	}
}
