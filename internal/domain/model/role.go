package model

import "time"

type RoleFlag int32

type RoleMask int32

// 语义化的成员分工角色
// 兼容掩码
const (
	RoleRawProvider RoleFlag = 1 << iota
	RoleTranslator
	RoleProofreader
	RoleTypesetter
	RoleReviewer
	RoleUploader
	RoleAdmin
)

func MaskRoles(roles []RoleFlag) RoleMask {
	mask := RoleMask(0)

	for _, role := range roles {
		mask |= RoleMask(role)
	}

	return mask
}

func UnmaskRoles(mask RoleMask) []RoleFlag {
	roles := make([]RoleFlag, 0, 6)

	if mask&RoleMask(RoleRawProvider) != 0 {
		roles = append(roles, RoleRawProvider)
	}
	if mask&RoleMask(RoleTranslator) != 0 {
		roles = append(roles, RoleTranslator)
	}
	if mask&RoleMask(RoleProofreader) != 0 {
		roles = append(roles, RoleProofreader)
	}
	if mask&RoleMask(RoleTypesetter) != 0 {
		roles = append(roles, RoleTypesetter)
	}
	if mask&RoleMask(RoleReviewer) != 0 {
		roles = append(roles, RoleReviewer)
	}
	if mask&RoleMask(RoleUploader) != 0 {
		roles = append(roles, RoleUploader)
	}
	if mask&RoleMask(RoleAdmin) != 0 {
		roles = append(roles, RoleAdmin)
	}

	return roles
}

func UnmaskRolesWithTime(mask RoleMask) []RoleWithTime {
	flags := UnmaskRoles(mask)
	roles := make([]RoleWithTime, len(flags))
	now := time.Now()

	for i, flag := range flags {
		roles[i] = RoleWithTime{
			Role:       flag,
			AssignedAt: now,
		}
	}

	return roles
}
