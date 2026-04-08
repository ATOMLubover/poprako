package model

import "time"

// MemberInvitationInfo 表示一条汉化组邀请记录
type MemberInvitationInfo struct {
	ID string

	InvitorID string
	// Invitor 仅在 includes 指定时填充
	Invitor *UserInfo

	InviteeQQ      string
	TeamID         string
	InvitationCode string

	// Pending 表示邀请是否仍有效
	Pending bool

	ToBeRawProvider bool
	ToBeTranslator  bool
	ToBeProofreader bool
	ToBeTypesetter  bool
	ToBeReviewer    bool
	ToBePublisher   bool
	ToBeAdmin       bool

	CreatedAt time.Time
}

// InvitedRoleMask 根据邀请中记录的角色信息计算 RoleMask
// 在非 app 层不应该使用这个方法！！
func (i *MemberInvitationInfo) InvitedRoleMask() RoleMask {
	mask := RoleMask(0)

	if i.ToBeRawProvider {
		mask |= RoleMask(RoleRawProvider)
	}
	if i.ToBeTranslator {
		mask |= RoleMask(RoleTranslator)
	}
	if i.ToBeProofreader {
		mask |= RoleMask(RoleProofreader)
	}
	if i.ToBeTypesetter {
		mask |= RoleMask(RoleTypesetter)
	}
	if i.ToBeReviewer {
		mask |= RoleMask(RoleReviewer)
	}
	if i.ToBePublisher {
		mask |= RoleMask(RolePublisher)
	}
	if i.ToBeAdmin {
		mask |= RoleMask(RoleAdmin)
	}

	return mask
}

// InvitedRoles 返回邀请中记录的角色列表
func (i *MemberInvitationInfo) InvitedRoles() []Role {
	roles := make([]Role, 0)

	if i.ToBeRawProvider {
		roles = append(roles, RoleRawProvider)
	}
	if i.ToBeTranslator {
		roles = append(roles, RoleTranslator)
	}
	if i.ToBeProofreader {
		roles = append(roles, RoleProofreader)
	}
	if i.ToBeTypesetter {
		roles = append(roles, RoleTypesetter)
	}
	if i.ToBeReviewer {
		roles = append(roles, RoleReviewer)
	}
	if i.ToBePublisher {
		roles = append(roles, RolePublisher)
	}
	if i.ToBeAdmin {
		roles = append(roles, RoleAdmin)
	}

	return roles
}

// MemberInvitationCreation 是创建邀请时的载荷
type MemberInvitationCreation struct {
	// ID 由 domain service 生成，外部不提供
	ID string

	InvitorID    string
	TargetTeamID string
	InviteeQQ    string

	InvitationCode string

	ToBeRawProvider bool
	ToBeTranslator  bool
	ToBeProofreader bool
	ToBeTypesetter  bool
	ToBeReviewer    bool
	ToBePublisher   bool
	ToBeAdmin       bool
}

// MemberInvitationUpdate 是邀请更新载荷，用于修改待分配角色
type MemberInvitationUpdate struct {
	ID string

	ToBeRawProvider bool
	ToBeTranslator  bool
	ToBeProofreader bool
	ToBeTypesetter  bool
	ToBeReviewer    bool
	ToBePublisher   bool
	ToBeAdmin       bool
}

// MemberInvitationQueryOpt 指定邀请查询的可选筛选条件
// 所有字段均为可空，nil 表示不参与筛选
type MemberInvitationQueryOpt struct {
	// TeamID 按目标汉化组 ID 筛选
	TeamID *string
	// InvitationCode 按邀请码筛选
	InvitationCode *string
	// Pending 按是否有效筛选
	Pending bool
}
