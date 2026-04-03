package service

import "poprako-s/internal/domain/model"

// RoleService 提供与角色掩码相关的转换功能
type RoleService interface {
	// MaskRoles 接受一个角色列表，返回对应的位掩码
	MaskRoles(
		roles []model.Role,
	) model.RoleMask

	// UnmaskRoles 接受一个位掩码，返回对应的角色列表
	UnmaskRoles(
		mask model.RoleMask,
	) []model.Role
}

// roleServiceImpl 是 RoleService 的具体实现，无内禀状态
type roleServiceImpl struct{}

// NewRoleService 返回 RoleService 的默认实现
func NewRoleService() RoleService {
	// 返回无状态实现
	return &roleServiceImpl{}
}

// MaskRoles 将角色切片编码为位掩码
func (*roleServiceImpl) MaskRoles(
	roles []model.Role,
) model.RoleMask {
	// 初始化空掩码
	mask := model.RoleMask(0)

	// 逐个将角色写入掩码
	for _, role := range roles {
		mask |= model.RoleMask(role)
	}

	// 返回最终掩码
	return mask
}

// UnmaskRoles 将位掩码解码为角色列表
func (*roleServiceImpl) UnmaskRoles(
	mask model.RoleMask,
) []model.Role {
	// 初始化结果切片
	roles := make([]model.Role, 0, 7)

	// 逐个校验每个角色位是否存在于掩码中
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

	// 返回解码后的角色列表
	return roles
}
