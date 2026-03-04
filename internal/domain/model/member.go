package model

import (
	"time"

	"labelplus-next-web-be/internal/util"
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

func NewMemberCreation(userID, teamID string, roles ...RoleFlag) *MemberCreation {
	mc := &MemberCreation{
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

/* func (mi *MemberInfo) RoleMask() RoleMask {
	mask := RoleMask(0)

	if mi.AssignedRawProviderAt != nil {
		mask |= RoleMask(RoleRawProvider)
	}
	if mi.AssignedTranslatorAt != nil {
		mask |= RoleMask(RoleTranslator)
	}
	if mi.AssignedProofreaderAt != nil {
		mask |= RoleMask(RoleProofreader)
	}
	if mi.AssignedTypesetterAt != nil {
		mask |= RoleMask(RoleTypesetter)
	}
	if mi.AssignedReviewerAt != nil {
		mask |= RoleMask(RoleReviewer)
	}
	if mi.AssignedUploaderAt != nil {
		mask |= RoleMask(RoleUploader)
	}
	if mi.AssignedAdminAt != nil {
		mask |= RoleMask(RoleAdmin)
	}

	return mask
} */

type MemberUpdate struct {
	ID string

	AssignRawProvider util.Option[time.Time]
	AssignTranslator  util.Option[time.Time]
	AssignProofreader util.Option[time.Time]
	AssignTypesetter  util.Option[time.Time]
	AssignReviewer    util.Option[time.Time]
	AssignUploader    util.Option[time.Time]
	AssignAdmin       util.Option[time.Time]
}

func NewMemberUpdate(id string, roles ...RoleFlag) *MemberUpdate {
	mu := &MemberUpdate{
		ID: id,
	}

	for _, role := range roles {
		switch role {
		case RoleRawProvider:
			mu.AssignRawProvider = util.NewSomeOption(time.Now())
		case RoleTranslator:
			mu.AssignTranslator = util.NewSomeOption(time.Now())
		case RoleProofreader:
			mu.AssignProofreader = util.NewSomeOption(time.Now())
		case RoleTypesetter:
			mu.AssignTypesetter = util.NewSomeOption(time.Now())
		case RoleReviewer:
			mu.AssignReviewer = util.NewSomeOption(time.Now())
		case RoleUploader:
			mu.AssignUploader = util.NewSomeOption(time.Now())
		case RoleAdmin:
			mu.AssignAdmin = util.NewSomeOption(time.Now())
		}
	}

	return mu
}
