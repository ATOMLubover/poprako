package val

import "poprako-s/internal/domain/model"

// MemberInfo 表示用于在应用层和接口层传递的成员信息 VO
type MemberInfo struct {
	// ID 是成员记录的唯一标识
	ID string `json:"id"`

	// UserID 是该成员对应的用户 ID
	UserID string `json:"user_id"`
	// TeamID 是该成员所属的团队 ID
	TeamID string `json:"team_id"`

	// Roles 是该成员所拥有的角色掩码
	Roles model.RoleMask `json:"roles"`

	// User 是可选的用户信息（仅在 includes 时填充）
	User *UserInfo `json:"user,omitempty"`
	// Team 是可选的团队信息（仅在 includes 时填充）
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
	// TeamID 是目标团队 ID
	TeamID string `json:"team_id" validate:"required"`
	// Roles 是分配的角色掩码
	Roles model.RoleMask `json:"roles" validate:"required"`
}

// CreateMemberRes 表示创建成员成功后的响应数据
type CreateMemberRes struct {
	// ID 是新创建成员记录的标识
	ID string `json:"id"`
}

// ListTeamMemberArgs 表示列出指定团队成员请求的参数
type ListTeamMemberArgs struct {
	// TeamID 是目标团队 ID
	TeamID string `json:"team_id" validate:"required"`
	Offset int    `json:"offset"`
	Limit  int    `json:"limit"`
}

// ListMyMemberArgs 表示列出当前用户所有成员记录请求的参数
type ListMyMemberArgs struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

// UpdateMemberRoleArgs 表示更新成员角色请求的参数
type UpdateMemberRoleArgs struct {
	// ID 是目标成员记录 ID
	ID string `json:"id" validate:"required"`
	// Roles 是目标角色掩码（PUT 语义全量替换）
	Roles model.RoleMask `json:"roles" validate:"required"`
}

// JoinTeamArgs 表示用户通过邀请码加入团队请求的参数
type JoinTeamArgs struct {
	// InvitationCode 是邀请码
	InvitationCode string `json:"invitation_code" validate:"required"`
}
