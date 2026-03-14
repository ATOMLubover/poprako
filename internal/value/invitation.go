package value

import (
	"errors"

	"labelplus-next-web-be/internal/domain/model"
)

type CreateInvitationArgs struct {
	TeamID    string         `json:"team_id"`
	InviteeQQ string         `json:"invitee_qq"`
	Roles     model.RoleMask `json:"roles"`
}

type ListTeamInvitationArgs struct {
	TeamID   string   `url:"team_id"`
	Includes []string `url:"includes[]"`
	PaginationParams
}

func (ltia *ListTeamInvitationArgs) Validate() error {
	if ltia == nil {
		return errors.New("参数不能为空")
	}

	if ltia.TeamID == "" {
		return errors.New("汉化组 ID 不能为空")
	}

	if err := ltia.PaginationParams.Validate(); err != nil {
		return errors.New("分页参数无效: " + err.Error())
	}

	return nil
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

	InvitorID      string    `json:"invitor_id"`
	InvitorInfo    *UserInfo `json:"invitor_info,omitempty"`
	InviteeQQ      string    `json:"invitee_qq"`
	InvitationCode string    `json:"invitation_code"`

	Pending bool `json:"pending"`

	Roles model.RoleMask `json:"roles"`

	CreatedAt int64 `json:"created_at"`
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
		return errors.New("汉化组 ID 不能为空")
	}

	if pia.Roles == 0 {
		return errors.New("分工角色不能为空")
	}

	return nil
}
