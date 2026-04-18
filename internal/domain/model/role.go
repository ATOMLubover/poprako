package model

// Role 表示单个成员角色的标识，使用位标志表示具体角色
type Role uint32

// RoleMask 表示多个 Role 的位掩码，用于高效地编码/解码角色集合
// **仅用于向外序列化，在 domain 内禁止使用 RoleMask 进行业务逻辑处理！！**
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

// MaskRoles 将角色切片编码为对应的 RoleMask
func MaskRoles(roles []Role) RoleMask {
	mask := RoleMask(0)

	for _, role := range roles {
		mask |= RoleMask(role)
	}

	return mask
}

// UnmaskRoles 从 RoleMask 解码得到角色切片
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
