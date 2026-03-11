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

type MemberWithUserInfo struct {
	ID string

	UserInfo UserInfo

	TeamID string

	AssignedRawProviderAt *time.Time
	AssignedTranslatorAt  *time.Time
	AssignedProofreaderAt *time.Time
	AssignedTypesetterAt  *time.Time
	AssignedReviewerAt    *time.Time
	AssignedPublishererAt *time.Time
	AssignedAdminAt       *time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (mp *MemberWithUserInfo) HasAnyRole(roles ...RoleFlag) bool {
	for _, role := range roles {
		switch role {
		case RoleRawProvider:
			if mp.AssignedRawProviderAt != nil {
				return true
			}
		case RoleTranslator:
			if mp.AssignedTranslatorAt != nil {
				return true
			}
		case RoleProofreader:
			if mp.AssignedProofreaderAt != nil {
				return true
			}
		case RoleTypesetter:
			if mp.AssignedTypesetterAt != nil {
				return true
			}
		case RoleReviewer:
			if mp.AssignedReviewerAt != nil {
				return true
			}
		case RolePublisher:
			if mp.AssignedPublishererAt != nil {
				return true
			}
		case RoleAdmin:
			if mp.AssignedAdminAt != nil {
				return true
			}
		}
	}

	return false
}

func (mp *MemberWithUserInfo) Roles() []RoleFlag {
	roles := make([]RoleFlag, 0)

	if mp.AssignedRawProviderAt != nil {
		roles = append(roles, RoleRawProvider)
	}
	if mp.AssignedTranslatorAt != nil {
		roles = append(roles, RoleTranslator)
	}
	if mp.AssignedProofreaderAt != nil {
		roles = append(roles, RoleProofreader)
	}
	if mp.AssignedTypesetterAt != nil {
		roles = append(roles, RoleTypesetter)
	}
	if mp.AssignedReviewerAt != nil {
		roles = append(roles, RoleReviewer)
	}
	if mp.AssignedPublishererAt != nil {
		roles = append(roles, RolePublisher)
	}
	if mp.AssignedAdminAt != nil {
		roles = append(roles, RoleAdmin)
	}

	return roles
}

type MemberInfo struct {
	ID string

	UserID string

	AssignRawProvider *time.Time
	AssignTranslator  *time.Time
	AssignProofreader *time.Time
	AssignTypesetter  *time.Time
	AssignReviewer    *time.Time
	AssignPublisher   *time.Time
	AssignAdmin       *time.Time
}

func NewMemberInfo(
	id string,
	userID string,
	assignRawProvider *time.Time,
	assignTranslator *time.Time,
	assignProofreader *time.Time,
	assignTypesetter *time.Time,
	assignReviewer *time.Time,
	assignPublisher *time.Time,
	assignAdmin *time.Time,
) MemberInfo {
	return MemberInfo{
		ID:                id,
		UserID:            userID,
		AssignRawProvider: assignRawProvider,
		AssignTranslator:  assignTranslator,
		AssignProofreader: assignProofreader,
		AssignTypesetter:  assignTypesetter,
		AssignReviewer:    assignReviewer,
		AssignPublisher:   assignPublisher,
		AssignAdmin:       assignAdmin,
	}
}

func (mi *MemberInfo) HasAnyRole(roles ...RoleFlag) bool {
	for _, role := range roles {
		switch role {
		case RoleRawProvider:
			if mi.AssignRawProvider != nil {
				return true
			}
		case RoleTranslator:
			if mi.AssignTranslator != nil {
				return true
			}
		case RoleProofreader:
			if mi.AssignProofreader != nil {
				return true
			}
		case RoleTypesetter:
			if mi.AssignTypesetter != nil {
				return true
			}
		case RoleReviewer:
			if mi.AssignReviewer != nil {
				return true
			}
		case RolePublisher:
			if mi.AssignPublisher != nil {
				return true
			}
		case RoleAdmin:
			if mi.AssignAdmin != nil {
				return true
			}
		}
	}

	return false
}

type MemberWithTeamInfo struct {
	ID string

	UserID string
	Team   TeamInfo

	AssignedRawProviderAt *time.Time
	AssignedTranslatorAt  *time.Time
	AssignedProofreaderAt *time.Time
	AssignedTypesetterAt  *time.Time
	AssignedReviewerAt    *time.Time
	AssignedPublisherAt   *time.Time
	AssignedAdminAt       *time.Time
}

func NewMemberWithTeamInfo(
	id string,
	userID string,
	team TeamInfo,
	assignedRawProviderAt *time.Time,
	assignedTranslatorAt *time.Time,
	assignedProofreaderAt *time.Time,
	assignedTypesetterAt *time.Time,
	assignedReviewerAt *time.Time,
	assignedPublisherAt *time.Time,
	assignedAdminAt *time.Time,
) MemberWithTeamInfo {
	return MemberWithTeamInfo{
		ID:                    id,
		UserID:                userID,
		Team:                  team,
		AssignedRawProviderAt: assignedRawProviderAt,
		AssignedTranslatorAt:  assignedTranslatorAt,
		AssignedProofreaderAt: assignedProofreaderAt,
		AssignedTypesetterAt:  assignedTypesetterAt,
		AssignedReviewerAt:    assignedReviewerAt,
		AssignedPublisherAt:   assignedPublisherAt,
		AssignedAdminAt:       assignedAdminAt,
	}
}

type MemberUpdate struct {
	ID                string
	AssignRawProvider *time.Time
	AssignTranslator  *time.Time
	AssignProofreader *time.Time
	AssignTypesetter  *time.Time
	AssignReviewer    *time.Time
	AssignPublisher   *time.Time
	AssignAdmin       *time.Time
}

func NewMemberUpdate(id string, current MemberWithUserInfo, targetRoles RoleMask) MemberUpdate {
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
		ID:                id,
		AssignRawProvider: resolveRoleAssignedAt(current.AssignedRawProviderAt, RoleRawProvider),
		AssignTranslator:  resolveRoleAssignedAt(current.AssignedTranslatorAt, RoleTranslator),
		AssignProofreader: resolveRoleAssignedAt(current.AssignedProofreaderAt, RoleProofreader),
		AssignTypesetter:  resolveRoleAssignedAt(current.AssignedTypesetterAt, RoleTypesetter),
		AssignReviewer:    resolveRoleAssignedAt(current.AssignedReviewerAt, RoleReviewer),
		AssignPublisher:   resolveRoleAssignedAt(current.AssignedPublishererAt, RolePublisher),
		AssignAdmin:       resolveRoleAssignedAt(current.AssignedAdminAt, RoleAdmin),
	}
}
