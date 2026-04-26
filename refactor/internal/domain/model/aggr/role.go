package aggr

import (
	"time"

	"poprako-s/internal/domain/model/enum"
)

type RoleMask uint32

func (m RoleMask) HasAnyRole(r ...enum.Role) bool {
	for _, role := range r {
		if m&RoleMask(role) != 0 {
			return true
		}
	}

	return false
}

func (m RoleMask) ToRoleArr() []enum.Role {
	var arr []enum.Role

	for i := int(enum.RoleInf); i <= int(enum.RoleSup); i = i << 1 {
		r := enum.Role(i)
		if m&RoleMask(r) != 0 {
			arr = append(arr, r)
		}
	}

	return arr
}

func (m RoleMask) ToRoleMask() RoleMask {
	return m
}

func (m *RoleMask) FromRoleArr(arr []enum.Role) {
	var mask RoleMask

	for _, role := range arr {
		mask |= RoleMask(role)
	}

	*m = mask
}

type TimedRoles struct {
	AssignedRawProviderAt *time.Time
	AssignedTranslatorAt  *time.Time
	AssignedProofreaderAt *time.Time
	AssignedTypesetterAt  *time.Time
	AssignedRedrawerAt    *time.Time
	AssignedReviewerAt    *time.Time
	AssignedPublisherAt   *time.Time
	AssignedAdminAt       *time.Time
}

func (t *TimedRoles) HasAnyRole(r ...enum.Role) bool {
	for _, role := range r {
		switch role {
		case enum.RoleRawProvider:
			if t.AssignedRawProviderAt != nil {
				return true
			}
		case enum.RoleTranslator:
			if t.AssignedTranslatorAt != nil {
				return true
			}
		case enum.RoleProofreader:
			if t.AssignedProofreaderAt != nil {
				return true
			}
		case enum.RoleTypesetter:
			if t.AssignedTypesetterAt != nil {
				return true
			}
		case enum.RoleRedrawer:
			if t.AssignedRedrawerAt != nil {
				return true
			}
		case enum.RoleReviewer:
			if t.AssignedReviewerAt != nil {
				return true
			}
		case enum.RolePublisher:
			if t.AssignedPublisherAt != nil {
				return true
			}
		case enum.RoleAdmin:
			if t.AssignedAdminAt != nil {
				return true
			}
		}
	}

	return false
}

func (t *TimedRoles) ToRoleArr() []enum.Role {
	var arr []enum.Role

	if t.AssignedRawProviderAt != nil {
		arr = append(arr, enum.RoleRawProvider)
	}
	if t.AssignedTranslatorAt != nil {
		arr = append(arr, enum.RoleTranslator)
	}
	if t.AssignedProofreaderAt != nil {
		arr = append(arr, enum.RoleProofreader)
	}
	if t.AssignedTypesetterAt != nil {
		arr = append(arr, enum.RoleTypesetter)
	}
	if t.AssignedRedrawerAt != nil {
		arr = append(arr, enum.RoleRedrawer)
	}
	if t.AssignedReviewerAt != nil {
		arr = append(arr, enum.RoleReviewer)
	}
	if t.AssignedPublisherAt != nil {
		arr = append(arr, enum.RolePublisher)
	}
	if t.AssignedAdminAt != nil {
		arr = append(arr, enum.RoleAdmin)
	}

	return arr
}

func (t *TimedRoles) ToRoleMask() RoleMask {
	var m RoleMask

	if t.AssignedRawProviderAt != nil {
		m |= RoleMask(enum.RoleRawProvider)
	}
	if t.AssignedTranslatorAt != nil {
		m |= RoleMask(enum.RoleTranslator)
	}
	if t.AssignedProofreaderAt != nil {
		m |= RoleMask(enum.RoleProofreader)
	}
	if t.AssignedTypesetterAt != nil {
		m |= RoleMask(enum.RoleTypesetter)
	}
	if t.AssignedRedrawerAt != nil {
		m |= RoleMask(enum.RoleRedrawer)
	}
	if t.AssignedReviewerAt != nil {
		m |= RoleMask(enum.RoleReviewer)
	}
	if t.AssignedPublisherAt != nil {
		m |= RoleMask(enum.RolePublisher)
	}
	if t.AssignedAdminAt != nil {
		m |= RoleMask(enum.RoleAdmin)
	}

	return m
}

func (t *TimedRoles) FromRoleArr(arr []enum.Role) {
	now := time.Now()

	for _, role := range arr {
		switch role {
		case enum.RoleRawProvider:
			if t.AssignedRawProviderAt == nil {
				t.AssignedRawProviderAt = &now
			}
		case enum.RoleTranslator:
			if t.AssignedTranslatorAt == nil {
				t.AssignedTranslatorAt = &now
			}
		case enum.RoleProofreader:
			if t.AssignedProofreaderAt == nil {
				t.AssignedProofreaderAt = &now
			}
		case enum.RoleTypesetter:
			if t.AssignedTypesetterAt == nil {
				t.AssignedTypesetterAt = &now
			}
		case enum.RoleRedrawer:
			if t.AssignedRedrawerAt == nil {
				t.AssignedRedrawerAt = &now
			}
		case enum.RoleReviewer:
			if t.AssignedReviewerAt == nil {
				t.AssignedReviewerAt = &now
			}
		case enum.RolePublisher:
			if t.AssignedPublisherAt == nil {
				t.AssignedPublisherAt = &now
			}
		case enum.RoleAdmin:
			if t.AssignedAdminAt == nil {
				t.AssignedAdminAt = &now
			}
		}
	}
}

func (t *TimedRoles) FromRoleMask(m RoleMask) {
	now := time.Now()

	if m&RoleMask(enum.RoleRawProvider) != 0 && t.AssignedRawProviderAt == nil {
		t.AssignedRawProviderAt = &now
	}
	if m&RoleMask(enum.RoleTranslator) != 0 && t.AssignedTranslatorAt == nil {
		t.AssignedTranslatorAt = &now
	}
	if m&RoleMask(enum.RoleProofreader) != 0 && t.AssignedProofreaderAt == nil {
		t.AssignedProofreaderAt = &now
	}
	if m&RoleMask(enum.RoleTypesetter) != 0 && t.AssignedTypesetterAt == nil {
		t.AssignedTypesetterAt = &now
	}
	if m&RoleMask(enum.RoleRedrawer) != 0 && t.AssignedRedrawerAt == nil {
		t.AssignedRedrawerAt = &now
	}
	if m&RoleMask(enum.RoleReviewer) != 0 && t.AssignedReviewerAt == nil {
		t.AssignedReviewerAt = &now
	}
	if m&RoleMask(enum.RolePublisher) != 0 && t.AssignedPublisherAt == nil {
		t.AssignedPublisherAt = &now
	}
	if m&RoleMask(enum.RoleAdmin) != 0 && t.AssignedAdminAt == nil {
		t.AssignedAdminAt = &now
	}
}
