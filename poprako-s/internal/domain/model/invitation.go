package model

import "time"

// InvitationInfo 表示一条团队邀请记录
type InvitationInfo struct {
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
func (ii *InvitationInfo) InvitedRoleMask() RoleMask {
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
	if ii.ToBePublisher {
		mask |= RoleMask(RolePublisher)
	}
	if ii.ToBeAdmin {
		mask |= RoleMask(RoleAdmin)
	}

	return mask
}

// InvitedRoles 返回邀请中记录的角色列表
func (ii *InvitationInfo) InvitedRoles() []Role {
	roles := make([]Role, 0)

	if ii.ToBeRawProvider {
		roles = append(roles, RoleRawProvider)
	}
	if ii.ToBeTranslator {
		roles = append(roles, RoleTranslator)
	}
	if ii.ToBeProofreader {
		roles = append(roles, RoleProofreader)
	}
	if ii.ToBeTypesetter {
		roles = append(roles, RoleTypesetter)
	}
	if ii.ToBeReviewer {
		roles = append(roles, RoleReviewer)
	}
	if ii.ToBePublisher {
		roles = append(roles, RolePublisher)
	}
	if ii.ToBeAdmin {
		roles = append(roles, RoleAdmin)
	}

	return roles
}

// InvitationCreation 是创建邀请时的载荷
type InvitationCreation struct {
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

// InvitationUpdate 是邀请更新载荷，用于修改待分配角色
type InvitationUpdate struct {
	ID string

	ToBeRawProvider bool
	ToBeTranslator  bool
	ToBeProofreader bool
	ToBeTypesetter  bool
	ToBeReviewer    bool
	ToBePublisher   bool
	ToBeAdmin       bool
}

// InvitationQueryOpt 指定邀请查询的可选筛选条件
// 所有字段均为可空，nil 表示不参与筛选
type InvitationQueryOpt struct {
	// TeamID 按目标团队 ID 筛选
	TeamID *string
	// InvitationCode 按邀请码筛选
	InvitationCode *string
	// Pending 按是否有效筛选
	Pending bool
}
