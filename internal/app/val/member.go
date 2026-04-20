package val

import "poprako-s/internal/domain/model"

// MemberInfo 表示用于在应用层和接口层传递的成员信息 VO
type MemberInfo struct {
	// ID 是成员记录的唯一标识
	ID string `json:"id"`

	// UserID 是该成员对应的用户 ID
	UserID string `json:"user_id"`
	// TeamID 是该成员所属的汉化组 ID
	TeamID string `json:"team_id"`

	// Roles 是该成员所拥有的角色位掩码，各位含义如下：
	//   bit 0 (1)  = RawProvider（图源）
	//   bit 1 (2)  = Translator（翻译）
	//   bit 2 (4)  = Proofreader（校对）
	//   bit 3 (8)  = Typesetter（嵌字）
	//   bit 4 (16) = Redrawer（美工）
	//   bit 5 (32) = Reviewer（监修）
	//   bit 6 (64) = Publisher（发布）
	//   bit 7 (128)= Admin（管理）
	Roles model.RoleMask `json:"roles"`

	// User 是可选的用户信息（仅在 includes 时填充）
	User *UserInfo `json:"user,omitempty"`
	// Team 是可选的汉化组信息（仅在 includes 时填充）
	Team *TeamInfo `json:"team,omitempty"`

	// CreatedAt 是记录创建时间的 Unix 毫秒时间戳
	CreatedAt int64 `json:"created_at"`
	// UpdatedAt 是记录最近一次更新时间的 Unix 毫秒时间戳
	UpdatedAt int64 `json:"updated_at"`
}

// CreateMemberArgs 表示创建成员请求的参数（超级管理员专用）
type CreateMemberArgs struct {
	// UserID 是要创建成员的用户 ID
	UserID string `json:"user_id" validate:"required"`
	// TeamID 是目标汉化组 ID
	TeamID string `json:"team_id" validate:"required"`
	// Roles 是分配的角色位掩码，各位含义如下：
	//   bit 0 (1)  = RawProvider（图源）
	//   bit 1 (2)  = Translator（翻译）
	//   bit 2 (4)  = Proofreader（校对）
	//   bit 3 (8)  = Typesetter（嵌字）
	//   bit 4 (16) = Reviewer（监修）
	//   bit 5 (32) = Publisher（发布）
	//   bit 6 (64) = Admin（管理员）
	//   bit 7 (128)= Redrawer（美工）
	Roles model.RoleMask `json:"roles" validate:"required"`
}

// CreateMemberRes 表示创建成员成功后的响应数据
type CreateMemberRes struct {
	// ID 是新创建成员记录的标识
	ID string `json:"id"`
}

// ListTeamMemberArgs 表示列出指定汉化组成员请求的参数
type ListTeamMemberArgs struct {
	// TeamID 是目标汉化组 ID
	TeamID   string   `json:"team_id" url:"team_id" validate:"required"`
	Includes []string `json:"includes" url:"includes"`
	Offset   int      `json:"offset" url:"offset"`
	Limit    int      `json:"limit" url:"limit"`
}

// ListMyMemberArgs 表示列出当前用户所有成员记录请求的参数
type ListMyMemberArgs struct {
	Includes []string `json:"includes" url:"includes"`
	Offset   int      `json:"offset" url:"offset"`
	Limit    int      `json:"limit" url:"limit"`
}

// UpdateMemberRoleArgs 表示更新成员角色请求的参数
type UpdateMemberRoleArgs struct {
	// ID 是目标成员记录 ID
	ID string `json:"id" validate:"required"`
	// Roles 是目标角色位掩码（PUT 语义全量替换），各位含义如下：
	//   bit 0 (1)  = RawProvider（图源）
	//   bit 1 (2)  = Translator（翻译）
	//   bit 2 (4)  = Proofreader（校对）
	//   bit 3 (8)  = Typesetter（嵌字）
	//   bit 4 (16) = Reviewer（监修）
	//   bit 5 (32) = Publisher（发布）
	//   bit 6 (64) = Admin（管理员）
	//   bit 7 (128)= Redrawer（美工）
	Roles model.RoleMask `json:"roles" validate:"required"`
}

// JoinTeamArgs 表示用户通过邀请码加入汉化组请求的参数
type JoinTeamArgs struct {
	// InvitationCode 是邀请码
	InvitationCode string `json:"invitation_code" validate:"required"`
}
