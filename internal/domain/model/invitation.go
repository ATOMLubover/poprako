package model

import (
	"time"
)

type InvitationCreation struct {
	InvitorID    string
	TargetTeamID string
	InviteeQQ    string

	InvitationCode string

	ToBeRawProvider bool
	ToBeTranslator  bool
	ToBeProofreader bool
	ToBeTypesetter  bool
	ToBeReviewer    bool
	ToBeUploader    bool
	ToBeAdmin       bool
}

func NewInvitationCreation(
	invitorID,
	targetTeamID,
	inviteeQQ string,
	invitationCode string,
	roles ...RoleFlag,
) *InvitationCreation {
	creation := &InvitationCreation{
		InvitorID:      invitorID,
		TargetTeamID:   targetTeamID,
		InviteeQQ:      inviteeQQ,
		InvitationCode: invitationCode,
	}

	creation.setRoles(roles...)

	return creation
}

func (ic *InvitationCreation) setRoles(roles ...RoleFlag) {
	for _, role := range roles {
		switch role {
		case RoleRawProvider:
			ic.ToBeRawProvider = true
		case RoleTranslator:
			ic.ToBeTranslator = true
		case RoleProofreader:
			ic.ToBeProofreader = true
		case RoleTypesetter:
			ic.ToBeTypesetter = true
		case RoleReviewer:
			ic.ToBeReviewer = true
		case RoleUploader:
			ic.ToBeUploader = true
		case RoleAdmin:
			ic.ToBeAdmin = true
		}
	}
}

type InvitationInfo struct {
	ID string

	InvitorID      string
	InviteeQQ      string
	TeamID         string
	InvitationCode string

	Pending bool

	ToBeRawProvider bool
	ToBeTranslator  bool
	ToBeProofreader bool
	ToBeTypesetter  bool
	ToBeReviewer    bool
	ToBeUploader    bool
	ToBeAdmin       bool

	CreatedAt time.Time
}

func (ii *InvitationInfo) RoleMask() RoleMask {
	mask := RoleMask(0)

	if ii.ToBeRawProvider {
		mask |= RoleMask(RoleRawProvider)
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
	if ii.ToBeUploader {
		mask |= RoleMask(RoleUploader)
	}
	if ii.ToBeAdmin {
		mask |= RoleMask(RoleAdmin)
	}

	return mask
}

type InvitationUpdate struct {
	ID              string
	ToBeRawProvider bool
	ToBeTranslator  bool
	ToBeProofreader bool
	ToBeTypesetter  bool
	ToBeReviewer    bool
	ToBeUploader    bool
	ToBeAdmin       bool
}

func NewInvitationUpdate(
	id string,
	roles ...RoleFlag,
) *InvitationUpdate {
	patch := &InvitationUpdate{
		ID: id,
	}

	patch.setRoles(roles...)

	return patch
}

func (ip *InvitationUpdate) setRoles(roles ...RoleFlag) {
	for _, role := range roles {
		switch role {
		case RoleRawProvider:
			ip.ToBeRawProvider = true
		case RoleTranslator:
			ip.ToBeTranslator = true
		case RoleProofreader:
			ip.ToBeProofreader = true
		case RoleTypesetter:
			ip.ToBeTypesetter = true
		case RoleReviewer:
			ip.ToBeReviewer = true
		case RoleUploader:
			ip.ToBeUploader = true
		case RoleAdmin:
			ip.ToBeAdmin = true
		}
	}
}
