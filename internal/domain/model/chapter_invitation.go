package model

import "time"

// ChapterInvitationInfo 表示一条章节邀请记录
type ChapterInvitationInfo struct {
	ID string

	ChapterID string
	InviterID string
	InviteeQQ string

	InvitationCode string
	Pending        bool

	ToBeRawProvider bool
	ToBeTranslator  bool
	ToBeProofreader bool
	ToBeTypesetter  bool
	ToBeRedrawer    bool
	ToBeReviewer    bool
	ToBePublisher   bool

	CreatedAt time.Time
	UpdatedAt time.Time
}

// InvitedRoleMask 根据邀请中记录的角色信息计算 RoleMask
func (i *ChapterInvitationInfo) InvitedRoleMask() RoleMask {
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
	if i.ToBeRedrawer {
		mask |= RoleMask(RoleRedrawer)
	}
	if i.ToBeReviewer {
		mask |= RoleMask(RoleReviewer)
	}
	if i.ToBePublisher {
		mask |= RoleMask(RolePublisher)
	}

	return mask
}

// ChapterInvitationCreation 是创建章节邀请时的载荷
type ChapterInvitationCreation struct {
	ID string

	ChapterID string
	InviterID string
	InviteeQQ string

	InvitationCode string

	ToBeRawProvider bool
	ToBeTranslator  bool
	ToBeProofreader bool
	ToBeTypesetter  bool
	ToBeRedrawer    bool
	ToBeReviewer    bool
	ToBePublisher   bool
}

// ChapterInvitationQueryOpt 指定章节邀请查询的可选筛选条件
type ChapterInvitationQueryOpt struct {
	ChapterID       *string
	InvitationCode  *string
	InviteeQQ       *string
	OnlyPendingTrue bool
}
