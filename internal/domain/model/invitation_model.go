package model

import (
	"time"
)

type InvitationCreation struct {
	InvitorID string
	InviteeQQ string

	InvitationCode string

	ToBePictureSource bool
	ToBeTranslator    bool
	ToBeProofreader   bool
	ToBeTypesetter    bool
	ToBeReviewer      bool
	ToBeAdmin         bool
	ToBeSuperAdmin    bool
}

func NewInvitationCreation(
	invitorID,
	inviteeQQ string,
	invitationCode string,
	roles ...RoleFlag,
) *InvitationCreation {
	creation := &InvitationCreation{
		InvitorID:      invitorID,
		InviteeQQ:      inviteeQQ,
		InvitationCode: invitationCode,
	}

	creation.setRoles(roles...)

	return creation
}

func (ic *InvitationCreation) setRoles(roles ...RoleFlag) {
	for _, role := range roles {
		switch role {
		case RolePictureSource:
			ic.ToBePictureSource = true
		case RoleTranslator:
			ic.ToBeTranslator = true
		case RoleProofreader:
			ic.ToBeProofreader = true
		case RoleTypesetter:
			ic.ToBeTypesetter = true
		case RoleReviewer:
			ic.ToBeReviewer = true
		case RoleAdmin:
			ic.ToBeAdmin = true
		case RoleSuperAdmin:
			ic.ToBeSuperAdmin = true
		}
	}
}

type InvitationInfo struct {
	ID string

	InvitorID      string
	InviteeQQ      string
	InvitationCode string

	Pending bool

	ToBePictureSource bool
	ToBeTranslator    bool
	ToBeProofreader   bool
	ToBeTypesetter    bool
	ToBeReviewer      bool
	ToBeAdmin         bool
	ToBeSuperAdmin    bool

	CreatedAt time.Time
}

func (ii *InvitationInfo) RoleMask() RoleMask {
	mask := RoleMask(0)

	if ii.ToBePictureSource {
		mask |= RoleMask(RolePictureSource)
	}
	if ii.ToBeTranslator {
		mask |= RoleMask(RoleTranslator)
	}
	if ii.ToBeProofreader {
		mask |= RoleMask(RoleProofreader)
	}
	if ii.ToBeTypesetter {
		mask |= RoleMask(RoleTypesetter)
	}
	if ii.ToBeReviewer {
		mask |= RoleMask(RoleReviewer)
	}
	if ii.ToBeAdmin {
		mask |= RoleMask(RoleAdmin)
	}
	if ii.ToBeSuperAdmin {
		mask |= RoleMask(RoleSuperAdmin)
	}

	return mask
}

type InvitationPatch struct {
	ID                string
	ToBePictureSource bool
	ToBeTranslator    bool
	ToBeProofreader   bool
	ToBeTypesetter    bool
	ToBeReviewer      bool
	ToBeAdmin         bool
	ToBeSuperAdmin    bool
}

func NewInvitationPatch(
	id string,
	roles ...RoleFlag,
) *InvitationPatch {
	patch := &InvitationPatch{
		ID: id,
	}

	patch.setRoles(roles...)

	return patch
}

func (ip *InvitationPatch) setRoles(roles ...RoleFlag) {
	for _, role := range roles {
		switch role {
		case RolePictureSource:
			ip.ToBePictureSource = true
		case RoleTranslator:
			ip.ToBeTranslator = true
		case RoleProofreader:
			ip.ToBeProofreader = true
		case RoleTypesetter:
			ip.ToBeTypesetter = true
		case RoleReviewer:
			ip.ToBeReviewer = true
		case RoleAdmin:
			ip.ToBeAdmin = true
		case RoleSuperAdmin:
			ip.ToBeSuperAdmin = true
		}
	}
}
