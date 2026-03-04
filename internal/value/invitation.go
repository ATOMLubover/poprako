package value

import (
	"errors"

	"labelplus-next-web-be/internal/domain/model"

	"go.uber.org/zap"
)

type CreateInvitationArgs struct {
	TeamID    string         `json:"team_id"`
	InviteeQQ string         `json:"invitee_qq"`
	Roles     model.RoleMask `json:"roles"`
}

func (cia *CreateInvitationArgs) Validate() error {
	if cia == nil {
		return errors.New("参数不能为空")
	}

	if cia.InviteeQQ == "" {
		return errors.New("被邀请人 QQ 不能为空")
	}

	if cia.Roles == 0 {
		return errors.New("分工角色不能为空")
	}

	return nil
}

type InvitationInfo struct {
	ID string `json:"id"`

	InvitorID string `json:"invitor_id"`
	InviteeQQ string `json:"invitee_qq"`

	Pending bool `json:"pending"`

	Roles model.RoleMask `json:"roles"`

	CreatedAt int64 `json:"created_at"`
}

func NewInvitationInfo(
	id string,
	invitorID string,
	inviteeQQ string,
	pengding bool,
	roles model.RoleMask,
	createdAt int64,
) *InvitationInfo {
	return &InvitationInfo{
		ID:        id,
		InvitorID: invitorID,
		InviteeQQ: inviteeQQ,
		Pending:   pengding,
		Roles:     roles,
		CreatedAt: createdAt,
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
		Pending:   invitation.Pending,
		Roles:     invitation.RoleMask(),
		CreatedAt: invitation.CreatedAt.UnixMilli(),
	}
}

type UpdateInvitationArgs struct {
	ID     string         `json:"id"`
	TeamID string         `json:"team_id"`
	Roles  model.RoleMask `json:"roles"`
}

func (pia *UpdateInvitationArgs) Validate() error {
	if pia == nil {
		return errors.New("参数不能为空")
	}

	if pia.ID == "" {
		return errors.New("邀请 ID 不能为空")
	}

	if pia.TeamID == "" {
		return errors.New("团队 ID 不能为空")
	}

	if pia.Roles == 0 {
		return errors.New("分工角色不能为空")
	}

	return nil
}
