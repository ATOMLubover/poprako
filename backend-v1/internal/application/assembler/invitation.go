package assembler

import (
	"labelplus-next-web-be/internal/domain/model"
	"labelplus-next-web-be/internal/value"
)

func AssembleInvitationInfo(invitationInfo model.InvitationInfo, onLoadURL OnLoadURL) value.InvitationInfo {
	result := value.InvitationInfo{
		ID:        invitationInfo.ID,
		InvitorID: invitationInfo.InvitorID,
		InviteeQQ: invitationInfo.InviteeQQ,
		Pending:   invitationInfo.Pending,
		Roles:     invitationInfo.RoleMask(),
		CreatedAt: invitationInfo.CreatedAt.UnixMilli(),
	}

	if invitationInfo.Invitor != nil {
		invitor := AssembleUserInfo(*invitationInfo.Invitor, onLoadURL)
		result.Invitor = &invitor
	}

	return result
}
