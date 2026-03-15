package model

import (
	"time"
)

type MemberCreation struct {
	UserID string
	TeamID string

	ToBeRawProvider bool
	ToBeTranslator  bool
	ToBeProofreader bool
	ToBeTypesetter  bool
	ToBeReviewer    bool
	ToBePublisher   bool
	ToBeAdmin       bool
}

func NewMemberCreation(userID, teamID string, roles ...RoleFlag) MemberCreation {
	mc := MemberCreation{
		UserID: userID,
		TeamID: teamID,
	}

	for _, role := range roles {
		switch role {
		case RoleRawProvider:
			mc.ToBeRawProvider = true
		case RoleTranslator:
			mc.ToBeTranslator = true
		case RoleProofreader:
			mc.ToBeProofreader = true
		case RoleTypesetter:
			mc.ToBeTypesetter = true
		case RoleReviewer:
			mc.ToBeReviewer = true
		case RolePublisher:
			mc.ToBePublisher = true
		case RoleAdmin:
			mc.ToBeAdmin = true
		}
	}

	return mc
}

type MemberInfo struct {
	ID string

	UserID string
	TeamID string

	// User 仅在 includes 指定时填充。
	User *UserInfo
	// Team 仅在 includes 指定时填充。
	Team *TeamInfo

	AssignedRawProviderAt *time.Time
	AssignedTranslatorAt  *time.Time
	AssignedProofreaderAt *time.Time
	AssignedTypesetterAt  *time.Time
	AssignedReviewerAt    *time.Time
	AssignedPublisherAt   *time.Time
	AssignedAdminAt       *time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewMemberInfo(
	id string,
	userID string,
	teamID string,
	assignedRawProviderAt *time.Time,
	assignedTranslatorAt *time.Time,
	assignedProofreaderAt *time.Time,
	assignedTypesetterAt *time.Time,
	assignedReviewerAt *time.Time,
	assignedPublisherAt *time.Time,
	assignedAdminAt *time.Time,
	createdAt time.Time,
	updatedAt time.Time,
) MemberInfo {
	return MemberInfo{
		ID:                    id,
		UserID:                userID,
		TeamID:                teamID,
		AssignedRawProviderAt: assignedRawProviderAt,
		AssignedTranslatorAt:  assignedTranslatorAt,
		AssignedProofreaderAt: assignedProofreaderAt,
		AssignedTypesetterAt:  assignedTypesetterAt,
		AssignedReviewerAt:    assignedReviewerAt,
		AssignedPublisherAt:   assignedPublisherAt,
		AssignedAdminAt:       assignedAdminAt,
		CreatedAt:             createdAt,
		UpdatedAt:             updatedAt,
	}
}

func (mi *MemberInfo) HasAnyRole(roles ...RoleFlag) bool {
	for _, role := range roles {
		switch role {
		case RoleRawProvider:
			if mi.AssignedRawProviderAt != nil {
				return true
			}
		case RoleTranslator:
			if mi.AssignedTranslatorAt != nil {
				return true
			}
		case RoleProofreader:
			if mi.AssignedProofreaderAt != nil {
				return true
			}
		case RoleTypesetter:
			if mi.AssignedTypesetterAt != nil {
				return true
			}
		case RoleReviewer:
			if mi.AssignedReviewerAt != nil {
				return true
			}
		case RolePublisher:
			if mi.AssignedPublisherAt != nil {
				return true
			}
		case RoleAdmin:
			if mi.AssignedAdminAt != nil {
				return true
			}
		}
	}

	return false
}
func (mi *MemberInfo) Roles() []RoleFlag {
	roles := make([]RoleFlag, 0)

	if mi.AssignedRawProviderAt != nil {
		roles = append(roles, RoleRawProvider)
	}
	if mi.AssignedTranslatorAt != nil {
		roles = append(roles, RoleTranslator)
	}
	if mi.AssignedProofreaderAt != nil {
		roles = append(roles, RoleProofreader)
	}
	if mi.AssignedTypesetterAt != nil {
		roles = append(roles, RoleTypesetter)
	}
	if mi.AssignedReviewerAt != nil {
		roles = append(roles, RoleReviewer)
	}
	if mi.AssignedPublisherAt != nil {
		roles = append(roles, RolePublisher)
	}
	if mi.AssignedAdminAt != nil {
		roles = append(roles, RoleAdmin)
	}

	return roles
}

type MemberUpdate struct {
	ID                    string
	AssignedRawProviderAt *time.Time
	AssignedTranslatorAt  *time.Time
	AssignedProofreaderAt *time.Time
	AssignedTypesetterAt  *time.Time
	AssignedReviewerAt    *time.Time
	AssignedPublisherAt   *time.Time
	AssignedAdminAt       *time.Time
}

func NewMemberUpdate(id string, current MemberInfo, targetRoles RoleMask) MemberUpdate {
	now := time.Now()

	resolveRoleAssignedAt := func(currentAssignedAt *time.Time, role RoleFlag) *time.Time {
		if targetRoles&RoleMask(role) == 0 {
			return nil
		}

		if currentAssignedAt != nil {
			t := *currentAssignedAt
			return &t
		}

		t := now
		return &t
	}

	return MemberUpdate{
		ID:                    id,
		AssignedRawProviderAt: resolveRoleAssignedAt(current.AssignedRawProviderAt, RoleRawProvider),
		AssignedTranslatorAt:  resolveRoleAssignedAt(current.AssignedTranslatorAt, RoleTranslator),
		AssignedProofreaderAt: resolveRoleAssignedAt(current.AssignedProofreaderAt, RoleProofreader),
		AssignedTypesetterAt:  resolveRoleAssignedAt(current.AssignedTypesetterAt, RoleTypesetter),
		AssignedReviewerAt:    resolveRoleAssignedAt(current.AssignedReviewerAt, RoleReviewer),
		AssignedPublisherAt:   resolveRoleAssignedAt(current.AssignedPublisherAt, RolePublisher),
		AssignedAdminAt:       resolveRoleAssignedAt(current.AssignedAdminAt, RoleAdmin),
	}
}
