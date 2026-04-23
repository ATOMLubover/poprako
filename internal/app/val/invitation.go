package val

import "poprako-s/internal/domain/model"

// InvitationInfo 表示用于在应用层和接口层传递的邀请信息 VO
type InvitationInfo struct {
	// ID 是邀请记录的唯一标识
	ID string `json:"id"`

	// InvitorID 是邀请发起者的用户 ID
	InvitorID string `json:"invitor_id"`
	// InviteeQQ 是被邀请者的 QQ 号
	InviteeQQ string `json:"invitee_qq"`
	// TeamID 是目标汉化组 ID
	TeamID string `json:"team_id"`
	// InvitationCode 是邀请码
	InvitationCode string `json:"invitation_code"`

	// Pending 表示邀请是否仍有效
	Pending bool `json:"pending"`
	// Roles 是邀请中指定的角色位掩码，各位含义如下：
	//   bit 0 (1)  = RawProvider（图源）
	//   bit 1 (2)  = Translator（翻译）
	//   bit 2 (4)  = Proofreader（校对）
	//   bit 3 (8)  = Typesetter（嵌字）
	//   bit 4 (16) = Redrawer（美工）
	//   bit 5 (32) = Reviewer（监修）
	//   bit 6 (64) = Publisher（发布）
	//   bit 7 (128)= Admin（管理）
	Roles model.RoleMask `json:"roles"`

	// CreatedAt 是记录创建时间的 Unix 毫秒时间戳
	CreatedAt int64 `json:"created_at"`
}

// ListTeamInvitationArgs 表示列出指定汉化组邀请请求的参数
type ListTeamInvitationArgs struct {
	// TeamID 是目标汉化组 ID
	TeamID string `json:"team_id" url:"team_id" validate:"required"`
	// Pending 是可选的邀请有效性筛选条件
	Pending *bool `json:"pending,omitempty" url:"pending"`
	Offset  int   `json:"offset" url:"offset"`
	Limit   int   `json:"limit" url:"limit"`
}

// CreateInvitationArgs 表示创建邀请请求的参数
type CreateInvitationArgs struct {
	// TeamID 是目标汉化组 ID
	TeamID string `json:"team_id" validate:"required"`
	// InviteeQQ 是被邀请者的 QQ 号
	InviteeQQ string `json:"invitee_qq" validate:"required"`
	// Roles 是邀请中指定的角色位掩码，各位含义如下：
	//   bit 0 (1)  = RawProvider（图源）
	//   bit 1 (2)  = Translator（翻译）
	//   bit 2 (4)  = Proofreader（校对）
	//   bit 3 (8)  = Typesetter（嵌字）
	//   bit 4 (16) = Redrawer（美工）
	//   bit 5 (32) = Reviewer（监修）
	//   bit 6 (64) = Publisher（发布）
	//   bit 7 (128)= Admin（管理）
	Roles model.RoleMask `json:"roles" validate:"required"`
}

// UpdateInvitationArgs 表示更新邀请请求的参数
type UpdateInvitationArgs struct {
	// ID 是要更新的邀请记录 ID
	ID string `json:"id" validate:"required"`
	// TeamID 是目标汉化组 ID，用于鉴权
	TeamID string `json:"team_id" validate:"required"`
	// Roles 是更新后的角色位掩码，各位含义如下：
	//   bit 0 (1)  = RawProvider（图源）
	//   bit 1 (2)  = Translator（翻译）
	//   bit 2 (4)  = Proofreader（校对）
	//   bit 3 (8)  = Typesetter（嵌字）
	//   bit 4 (16) = Redrawer（美工）
	//   bit 5 (32) = Reviewer（监修）
	//   bit 6 (64) = Publisher（发布）
	//   bit 7 (128)= Admin（管理）
	Roles model.RoleMask `json:"roles" validate:"required"`
}
