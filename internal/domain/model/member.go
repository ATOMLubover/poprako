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

// MemberWithInfo 是统一的成员详情模型，替代分裂的 MemberWithUserInfo 和 MemberWithTeamInfo。
// User 和 Team 字段均为可选，仅在 includes 指定时填充。
type MemberWithInfo struct {
	ID     string
	UserID string
	TeamID string

	User *UserInfo
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

func NewMemberWithInfo(
	id string,
	userID string,
	teamID string,
	user *UserInfo,
	team *TeamInfo,
	assignedRawProviderAt *time.Time,
	assignedTranslatorAt *time.Time,
	assignedProofreaderAt *time.Time,
	assignedTypesetterAt *time.Time,
	assignedReviewerAt *time.Time,
	assignedPublisherAt *time.Time,
	assignedAdminAt *time.Time,
	createdAt time.Time,
	updatedAt time.Time,
) MemberWithInfo {
	return MemberWithInfo{
		ID:                    id,
		UserID:                userID,
		TeamID:                teamID,
		User:                  user,
		Team:                  team,
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

func (m *MemberWithInfo) HasAnyRole(roles ...RoleFlag) bool {
	for _, role := range roles {
		switch role {
		case RoleRawProvider:
			if m.AssignedRawProviderAt != nil {
				return true
			}
		case RoleTranslator:
			if m.AssignedTranslatorAt != nil {
				return true
			}
		case RoleProofreader:
			if m.AssignedProofreaderAt != nil {
				return true
			}
		case RoleTypesetter:
			if m.AssignedTypesetterAt != nil {
				return true
			}
		case RoleReviewer:
			if m.AssignedReviewerAt != nil {
				return true
			}
		case RolePublisher:
			if m.AssignedPublisherAt != nil {
				return true
			}
		case RoleAdmin:
			if m.AssignedAdminAt != nil {
				return true
			}
		}
	}

	return false
}

func (m *MemberWithInfo) Roles() []RoleFlag {
	roles := make([]RoleFlag, 0)

	if m.AssignedRawProviderAt != nil {
		roles = append(roles, RoleRawProvider)
	}
	if m.AssignedTranslatorAt != nil {
		roles = append(roles, RoleTranslator)
	}
	if m.AssignedProofreaderAt != nil {
		roles = append(roles, RoleProofreader)
	}
	if m.AssignedTypesetterAt != nil {
		roles = append(roles, RoleTypesetter)
	}
	if m.AssignedReviewerAt != nil {
		roles = append(roles, RoleReviewer)
	}
	if m.AssignedPublisherAt != nil {
		roles = append(roles, RolePublisher)
	}
	if m.AssignedAdminAt != nil {
		roles = append(roles, RoleAdmin)
	}

	return roles
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

func NewMemberUpdate(id string, current MemberWithInfo, targetRoles RoleMask) MemberUpdate {
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
		AssignPublisher:   resolveRoleAssignedAt(current.AssignedPublisherAt, RolePublisher),
		AssignAdmin:       resolveRoleAssignedAt(current.AssignedAdminAt, RoleAdmin),
	}
}
