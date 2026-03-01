package model

import (
	"time"
)

type UserInfo struct {
	ID string

	Name string
	QQ   string

	AvatarURL string

	AssignedPictureSource *time.Time
	AssignedTranslatorAt  *time.Time
	AssignedProofreaderAt *time.Time
	AssignedTypesetterAt  *time.Time
	AssignedReviewerAt    *time.Time

	AssignedAdminAt      *time.Time
	AssignedSuperAdminAt *time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (u *UserInfo) HasRole(role RoleFlag) bool {
	switch role {
	case RolePictureSource:
		return u.AssignedPictureSource != nil
	case RoleTranslator:
		return u.AssignedTranslatorAt != nil
	case RoleProofreader:
		return u.AssignedProofreaderAt != nil
	case RoleTypesetter:
		return u.AssignedTypesetterAt != nil
	case RoleReviewer:
		return u.AssignedReviewerAt != nil
	case RoleAdmin:
		return u.AssignedAdminAt != nil
	case RoleSuperAdmin:
		return u.AssignedSuperAdminAt != nil
	default:
		return false
	}
}

func (u *UserInfo) HasAnyRole(roles ...RoleFlag) bool {
	for _, role := range roles {
		if u.HasRole(role) {
			return true
		}
	}

	return false
}

func (u *UserInfo) HasAllRoles(roles ...RoleFlag) bool {
	for _, role := range roles {
		if !u.HasRole(role) {
			return false
		}
	}

	return true
}

func (u *UserInfo) MaskRoles() RoleMask {
	mask := RoleMask(0)

	if u.HasRole(RolePictureSource) {
		mask |= RoleMask(RolePictureSource)
	}
	if u.HasRole(RoleTranslator) {
		mask |= RoleMask(RoleTranslator)
	}
	if u.HasRole(RoleProofreader) {
		mask |= RoleMask(RoleProofreader)
	}
	if u.HasRole(RoleTypesetter) {
		mask |= RoleMask(RoleTypesetter)
	}
	if u.HasRole(RoleReviewer) {
		mask |= RoleMask(RoleReviewer)
	}
	if u.HasRole(RoleAdmin) {
		mask |= RoleMask(RoleAdmin)
	}
	if u.HasRole(RoleSuperAdmin) {
		mask |= RoleMask(RoleSuperAdmin)
	}

	return mask
}

// FIXME: 是否需要这个函数？
func (u *UserInfo) UnmaskRoles(mask RoleMask) []RoleFlag {
	roles := make([]RoleFlag, 0, 6)

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
	if mask&RoleMask(RoleAdmin) != 0 {
		roles = append(roles, RoleAdmin)
	}
	if mask&RoleMask(RoleSuperAdmin) != 0 {
		roles = append(roles, RoleSuperAdmin)
	}

	return roles
}

// UserCredentials 仅用于登录时校验，不对外暴露完整用户信息
type UserCredentials struct {
	UserID       string
	PasswordHash string
}

type UserRegistration struct {
	Name              string
	QQ                string
	Password          string
	ToBePictureSource bool
	ToBeTranslator    bool
	ToBeProofreader   bool
	ToBeTypesetter    bool
	ToBeReviewer      bool
	ToBeAdmin         bool
	ToBeSuperAdmin    bool
}

func NewUserRegistration(
	name string,
	qq string,
	password string,
	roles ...RoleFlag,
) *UserRegistration {
	registration := &UserRegistration{
		Name:           name,
		QQ:             qq,
		Password:       password,
	}

	registration.setRoles(roles...)

	return registration
}

func (ur *UserRegistration) setRoles(roles ...RoleFlag) {
	for _, role := range roles {
		switch role {
		case RolePictureSource:
			ur.ToBePictureSource = true
		case RoleTranslator:
			ur.ToBeTranslator = true
		case RoleProofreader:
			ur.ToBeProofreader = true
		case RoleTypesetter:
			ur.ToBeTypesetter = true
		case RoleReviewer:
			ur.ToBeReviewer = true
		case RoleAdmin:
			ur.ToBeAdmin = true
		case RoleSuperAdmin:
			ur.ToBeSuperAdmin = true
		}
	}
}
