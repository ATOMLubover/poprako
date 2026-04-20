package model

import "time"

// MemberInfo 表示用户在某个汉化组中的成员信息，包含角色分配时间
type MemberInfo struct {
	ID string

	UserID string
	TeamID string

	// User 仅在 includes 指定时填充
	User *UserInfo
	// Team 仅在 includes 指定时填充
	Team *TeamInfo

	AssignedRawProviderAt *time.Time
	AssignedTranslatorAt  *time.Time
	AssignedProofreaderAt *time.Time
	AssignedTypesetterAt  *time.Time
	AssignedRedrawerAt    *time.Time
	AssignedReviewerAt    *time.Time
	AssignedPublisherAt   *time.Time
	AssignedAdminAt       *time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
}

// HasAnyRole 检查当前成员是否拥有任意给定的角色
func (mi *MemberInfo) HasAnyRole(roles ...Role) bool {
	for _, role := range roles {
		switch role {
		case RoleRawProvider:
			if mi.AssignedRawProviderAt != nil {
				return true
			}
		case RoleTranslator:
			if mi.AssignedTranslatorAt != nil {
				return true
			}
		case RoleProofreader:
			if mi.AssignedProofreaderAt != nil {
				return true
			}
		case RoleTypesetter:
			if mi.AssignedTypesetterAt != nil {
				return true
			}
		case RoleRedrawer:
			if mi.AssignedRedrawerAt != nil {
				return true
			}
		case RoleReviewer:
			if mi.AssignedReviewerAt != nil {
				return true
			}
		case RolePublisher:
			if mi.AssignedPublisherAt != nil {
				return true
			}
		case RoleAdmin:
			if mi.AssignedAdminAt != nil {
				return true
			}
		}
	}

	return false
}

// Roles 返回当前成员所拥有的所有角色列表
func (mi *MemberInfo) Roles() []Role {
	roles := make([]Role, 0, 8)

	if mi.AssignedRawProviderAt != nil {
		roles = append(roles, RoleRawProvider)
	}
	if mi.AssignedTranslatorAt != nil {
		roles = append(roles, RoleTranslator)
	}
	if mi.AssignedProofreaderAt != nil {
		roles = append(roles, RoleProofreader)
	}
	if mi.AssignedTypesetterAt != nil {
		roles = append(roles, RoleTypesetter)
	}
	if mi.AssignedRedrawerAt != nil {
		roles = append(roles, RoleRedrawer)
	}
	if mi.AssignedReviewerAt != nil {
		roles = append(roles, RoleReviewer)
	}
	if mi.AssignedPublisherAt != nil {
		roles = append(roles, RolePublisher)
	}
	if mi.AssignedAdminAt != nil {
		roles = append(roles, RoleAdmin)
	}

	return roles
}

// MemberCreation 是创建成员时的载荷
type MemberCreation struct {
	// ID 由 domain service 生成，外部不提供
	ID string

	UserID string
	TeamID string

	// 以下布尔字段表示待分配的角色
	ToBeRawProvider bool
	ToBeTranslator  bool
	ToBeProofreader bool
	ToBeTypesetter  bool
	ToBeRedrawer    bool
	ToBeReviewer    bool
	ToBePublisher   bool
	ToBeAdmin       bool
}

// MemberUpdate 用于 PUT 语义的角色全量替换，保留已有角色的时间戳
type MemberUpdate struct {
	ID string

	AssignedRawProviderAt *time.Time
	AssignedTranslatorAt  *time.Time
	AssignedProofreaderAt *time.Time
	AssignedTypesetterAt  *time.Time
	AssignedRedrawerAt    *time.Time
	AssignedReviewerAt    *time.Time
	AssignedPublisherAt   *time.Time
	AssignedAdminAt       *time.Time
}

// MemberQueryOpt 指定成员查询的可选筛选条件
// 所有字段均为可空，nil 表示不参与筛选
type MemberQueryOpt struct {
	// ID 按成员 ID 筛选
	ID *string
	// UserID 按用户 ID 筛选
	UserID *string
	// TeamID 按汉化组 ID 筛选
	TeamID *string
}
