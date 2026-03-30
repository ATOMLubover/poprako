package model

import "time"

type AssignmentInfo struct {
	ID string

	ChapterID string
	// Chapter 仅在 includes 指定时填充。
	Chapter *ChapterInfo

	UserID string
	// User 仅在 includes 指定时填充。
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
