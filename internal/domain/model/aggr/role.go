package aggr

import (
	"time"

	"poprako-s/internal/domain/model/enum"
)

// `RoleMask` is a bitmask that encodes a set of roles as a single uint32 value
type RoleMask uint32

// `HasAnyRole` reports whether the mask contains at least one of the given roles
func (m RoleMask) HasAnyRole(r ...enum.Role) bool {
	for _, role := range r {
		if m&RoleMask(role) != 0 {
			return true
		}
	}

	return false
}

// `ToRoleArr` expands the mask into a slice of individual roles in ascending bit order
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

// `ToRoleMask` returns the receiver as-is, satisfying the `WithRoles` interface
func (m RoleMask) ToRoleMask() RoleMask {
	return m
}

// `FromRoleArr` sets the mask by ORing together all roles in the given slice
func (m *RoleMask) FromRoleArr(arr []enum.Role) {
	var mask RoleMask

	for _, role := range arr {
		mask |= RoleMask(role)
	}

	*m = mask
}

// `TimedRoles` records when each role was assigned by storing a nullable timestamp per role.
// A nil timestamp means the role is not currently held
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

// `HasAnyRole` reports whether at least one of the given roles has been assigned
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

// `ToRoleArr` returns the slice of roles for which an assignment timestamp exists
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

// `ToRoleMask` converts the timed role set into a `RoleMask` bitmask
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

// `FromRoleArr` assigns roles from the given slice by setting each role's timestamp to now,
// only if the role has not been previously assigned
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

// `FromRoleMask` assigns roles from a `RoleMask` by setting each role's timestamp to now,
// only if the role has not been previously assigned
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
