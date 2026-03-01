package model

import (
	"errors"
	"time"
	"unicode/utf8"
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

type UserRegistration struct {
	Name           string
	QQ             string
	Password       string
	InvitationCode string
}

func (r *UserRegistration) Validate() error {
	if err := r.validateUserName(r.Name); err != nil {
		return err
	}

	if err := r.validatePassword(r.Password); err != nil {
		return err
	}

	return nil
}

func (*UserRegistration) validateUserName(name string) error {
	if len := utf8.RuneCountInString(name); len < 2 || len > 20 {
		return errors.New("用户名长度必须在 2 到 20 个字符之间")
	}

	return nil
}

func (*UserRegistration) validatePassword(password string) error {
	if len := utf8.RuneCountInString(password); len < 6 || len > 20 {
		return errors.New("密码长度必须在 6 到 20 个字符之间")
	}

	for _, r := range password {
		if !(r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9') {
			return errors.New("密码只能包含字母和数字")
		}
	}

	return nil
}
