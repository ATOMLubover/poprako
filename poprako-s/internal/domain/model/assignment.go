package model

import (
	"time"

	"poprako-s/internal/domain/event"
)

// AssignmentInfo 表示某个用户被分配到章节中的任务分配信息
type AssignmentInfo struct {
	ID string

	ChapterID string
	// Chapter 仅在 includes 指定时填充
	Chapter *ChapterInfo

	UserID string
	// User 仅在 includes 指定时填充
	User *UserInfo

	AssignedRawProviderAt *time.Time
	AssignedTranslatorAt  *time.Time
	AssignedProofreaderAt *time.Time
	AssignedTypesetterAt  *time.Time
	AssignedRedrawerAt    *time.Time
	AssignedReviewerAt    *time.Time
	AssignedPublisherAt   *time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
}

// HasAnyRole 检查当前分配信息是否包含任意给定的角色
func (a *AssignmentInfo) HasAnyRole(r ...Role) bool {
	for _, role := range r {
		switch role {
		case RoleRawProvider:
			if a.AssignedRawProviderAt != nil {
				return true
			}
		case RoleTranslator:
			if a.AssignedTranslatorAt != nil {
				return true
			}
		case RoleProofreader:
			if a.AssignedProofreaderAt != nil {
				return true
			}
		case RoleTypesetter:
			if a.AssignedTypesetterAt != nil {
				return true
			}
		case RoleReviewer:
			if a.AssignedReviewerAt != nil {
				return true
			}
		case RolePublisher:
			if a.AssignedPublisherAt != nil {
				return true
			}
		}
	}

	return false
}

// AssignedRoleMask 根据分配中记录的角色信息计算 RoleMask
func (a *AssignmentInfo) AssignedRoleMask() RoleMask {
	mask := RoleMask(0)

	if a.AssignedRawProviderAt != nil {
		mask |= RoleMask(RoleRawProvider)
	}
	if a.AssignedTranslatorAt != nil {
		mask |= RoleMask(RoleTranslator)
	}
	if a.AssignedProofreaderAt != nil {
		mask |= RoleMask(RoleProofreader)
	}
	if a.AssignedTypesetterAt != nil {
		mask |= RoleMask(RoleTypesetter)
	}
	if a.AssignedReviewerAt != nil {
		mask |= RoleMask(RoleReviewer)
	}
	if a.AssignedPublisherAt != nil {
		mask |= RoleMask(RolePublisher)
	}

	return mask
}

// AssignmentCreation 是创建分配记录时的载荷
type AssignmentCreation struct {
	// ID 由 domain service 生成，外部不提供
	ID string

	ChapterID string
	UserID    string

	AssignedRawProviderAt *time.Time
	AssignedTranslatorAt  *time.Time
	AssignedProofreaderAt *time.Time
	AssignedTypesetterAt  *time.Time
	AssignedRedrawerAt    *time.Time
	AssignedReviewerAt    *time.Time
	AssignedPublisherAt   *time.Time

	event.EventBase
}

// AssignmentUpdate 用于 PUT 语义的角色全量替换，保留已有角色的时间戳
type AssignmentUpdate struct {
	ID string

	AssignedRawProviderAt *time.Time
	AssignedTranslatorAt  *time.Time
	AssignedProofreaderAt *time.Time
	AssignedTypesetterAt  *time.Time
	AssignedRedrawerAt    *time.Time
	AssignedReviewerAt    *time.Time
	AssignedPublisherAt   *time.Time
}

// AssignmentQueryOpt 指定分配查询的可选筛选条件
// 所有字段均为可空，nil 表示不参与筛选
type AssignmentQueryOpt struct {
	// ChapterID 按所属章节 ID 筛选
	ChapterID *string
	// UserID 按用户 ID 筛选
	UserID *string
}
