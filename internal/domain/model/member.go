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
	ToBeUploader    bool
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
		case RoleUploader:
			mc.ToBeUploader = true
		case RoleAdmin:
			mc.ToBeAdmin = true
		}
	}

	return mc
}

type MemberProfile struct {
	ID string

	UserInfo *UserInfo

	TeamID string

	AssignedRawProviderAt *time.Time
	AssignedTranslatorAt  *time.Time
	AssignedProofreaderAt *time.Time
	AssignedTypesetterAt  *time.Time
	AssignedReviewerAt    *time.Time
	AssignedUploaderAt    *time.Time
	AssignedAdminAt       *time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (mp *MemberProfile) HasAnyRole(roles ...RoleFlag) bool {
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
		case RoleAdmin:
			if mp.AssignedAdminAt != nil {
				return true
			}
		}
	}

	return false
}

func (mp *MemberProfile) Roles() []RoleFlag {
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
	if mp.AssignedUploaderAt != nil {
		roles = append(roles, RoleUploader)
	}
	if mp.AssignedAdminAt != nil {
		roles = append(roles, RoleAdmin)
	}

	return roles
}

type RoleWithTime struct {
	Role       RoleFlag
	AssignedAt time.Time
}

type MemberInfo struct {
	ID string

	UserID string

	AssignRawProvider *time.Time
	AssignTranslator  *time.Time
	AssignProofreader *time.Time
	AssignTypesetter  *time.Time
	AssignReviewer    *time.Time
	AssignUploader    *time.Time
	AssignAdmin       *time.Time
}

func NewMemberInfo(id string, userID string, roles ...RoleWithTime) MemberInfo {
	memberInfo := MemberInfo{
		ID:     id,
		UserID: userID,
	}

	for _, role := range roles {
		switch role.Role {
		case RoleRawProvider:
			t := role.AssignedAt
			memberInfo.AssignRawProvider = &t
		case RoleTranslator:
			t := role.AssignedAt
			memberInfo.AssignTranslator = &t
		case RoleProofreader:
			t := role.AssignedAt
			memberInfo.AssignProofreader = &t
		case RoleTypesetter:
			t := role.AssignedAt
			memberInfo.AssignTypesetter = &t
		case RoleReviewer:
			t := role.AssignedAt
			memberInfo.AssignReviewer = &t
		case RoleUploader:
			t := role.AssignedAt
			memberInfo.AssignUploader = &t
		case RoleAdmin:
			t := role.AssignedAt
			memberInfo.AssignAdmin = &t
		}
	}

	return memberInfo
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
		case RoleUploader:
			if mi.AssignUploader != nil {
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

type MemberUpdate struct {
	ID                string
	AssignRawProvider *time.Time
	AssignTranslator  *time.Time
	AssignProofreader *time.Time
	AssignTypesetter  *time.Time
	AssignReviewer    *time.Time
	AssignUploader    *time.Time
	AssignAdmin       *time.Time
}

func NewMemberUpdate(id string, roles ...RoleWithTime) MemberUpdate {
	mu := MemberUpdate{
		ID: id,
	}

	for _, role := range roles {
		switch role.Role {
		case RoleRawProvider:
			t := role.AssignedAt
			mu.AssignRawProvider = &t
		case RoleTranslator:
			t := role.AssignedAt
			mu.AssignTranslator = &t
		case RoleProofreader:
			t := role.AssignedAt
			mu.AssignProofreader = &t
		case RoleTypesetter:
			t := role.AssignedAt
			mu.AssignTypesetter = &t
		case RoleReviewer:
			t := role.AssignedAt
			mu.AssignReviewer = &t
		case RoleUploader:
			t := role.AssignedAt
			mu.AssignUploader = &t
		case RoleAdmin:
			t := role.AssignedAt
			mu.AssignAdmin = &t
		}
	}

	return mu
}
