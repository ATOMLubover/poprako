package model

type Role uint32

type RoleMask uint32

// 语义化的成员分工角色，兼容掩码
const (
	RoleRawProvider Role = 1 << iota
	RoleTranslator
	RoleProofreader
	RoleTypesetter
	RoleReviewer
	RolePublisher
	RoleAdmin
)

func MaskRoles(roles []Role) RoleMask {
	mask := RoleMask(0)

	for _, role := range roles {
		mask |= RoleMask(role)
	}

	return mask
}

func UnmaskRoles(mask RoleMask) []Role {
	roles := make([]Role, 0, 6)

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
	if mask&RoleMask(RolePublisher) != 0 {
		roles = append(roles, RolePublisher)
	}
	if mask&RoleMask(RoleAdmin) != 0 {
		roles = append(roles, RoleAdmin)
	}

	return roles
}
