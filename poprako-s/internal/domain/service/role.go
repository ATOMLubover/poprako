package service

import "poprako-s/internal/domain/model"

type RoleService interface {
	// MaskRoles 接受一个角色列表，返回对应的掩码
	MaskRoles(roles []model.Role) model.RoleMask
	// UnmaskRoles 接受一个角色掩码，返回对应的角色列表
	UnmaskRoles(mask model.RoleMask) []model.Role
}

type roleServiceImpl struct{}

func NewRoleService() RoleService {
	return &roleServiceImpl{}
}

func (*roleServiceImpl) MaskRoles(roles []model.Role) model.RoleMask {
	mask := model.RoleMask(0)

	for _, role := range roles {
		mask |= model.RoleMask(role)
	}

	return mask
}

func (*roleServiceImpl) UnmaskRoles(mask model.RoleMask) []model.Role {
	roles := make([]model.Role, 0, 6)

	if mask&model.RoleMask(model.RoleRawProvider) != 0 {
		roles = append(roles, model.RoleRawProvider)
	}
	if mask&model.RoleMask(model.RoleTranslator) != 0 {
		roles = append(roles, model.RoleTranslator)
	}
	if mask&model.RoleMask(model.RoleProofreader) != 0 {
		roles = append(roles, model.RoleProofreader)
	}
	if mask&model.RoleMask(model.RoleTypesetter) != 0 {
		roles = append(roles, model.RoleTypesetter)
	}
	if mask&model.RoleMask(model.RoleReviewer) != 0 {
		roles = append(roles, model.RoleReviewer)
	}
	if mask&model.RoleMask(model.RolePublisher) != 0 {
		roles = append(roles, model.RolePublisher)
	}
	if mask&model.RoleMask(model.RoleAdmin) != 0 {
		roles = append(roles, model.RoleAdmin)
	}

	return roles
}
