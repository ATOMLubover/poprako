package value

import (
	"labelplus-next-web-be/internal/domain/model"

	"go.uber.org/zap"
)

type CreateInvitationArgs struct {
	InviteeQQ string         `json:"invitee_qq"`
	Roles     model.RoleMask `json:"roles"`
}

type InvitationInfo struct {
	ID string `json:"id"`

	InvitorID string `json:"invitor_id"`
	InviteeQQ string `json:"invitee_qq"`

	Roles model.RoleMask `json:"roles"`

	CreatedAt int64 `json:"created_at"`
}

func NewInvitationInfo(
	ID string,
	InvitorID string,
	InviteeQQ string,
	Roles model.RoleMask,
	CreatedAt int64,
) *InvitationInfo {
	return &InvitationInfo{
		ID:        ID,
		InvitorID: InvitorID,
		InviteeQQ: InviteeQQ,
		Roles:     Roles,
		CreatedAt: CreatedAt,
	}
}

func NewInvitationInfoFromModel(invitation *model.InvitationInfo) *InvitationInfo {
	if invitation == nil {
		zap.L().Warn("NewInvitationInfoFromModel: invitation 为空")
		return nil
	}

	return &InvitationInfo{
		ID:        invitation.ID,
		InvitorID: invitation.InvitorID,
		InviteeQQ: invitation.InviteeQQ,
		Roles:     invitation.RoleMask(),
		CreatedAt: invitation.CreatedAt.UnixMilli(),
	}
}

type PatchInvitationArgs struct {
	ID    string         `json:"id"`
	Roles model.RoleMask `json:"roles"`
}
