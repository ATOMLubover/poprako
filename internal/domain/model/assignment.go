package model

import "time"

type AssignmentInfo struct {
	ID string

	ComicID string
	UserID  string

	AssignedRawProviderAt *time.Time
	AssignedTranslatorAt  *time.Time
	AssignedProofreaderAt *time.Time
	AssignedTypesetterAt  *time.Time
	AssignedReviewerAt    *time.Time
	AssignedPublisherAt   *time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewAssignmentInfo(
	id string,
	comicID string,
	userID string,
	assignedRawProviderAt *time.Time,
	assignedTranslatorAt *time.Time,
	assignedProofreaderAt *time.Time,
	assignedTypesetterAt *time.Time,
	assignedReviewerAt *time.Time,
	assignedPublisherAt *time.Time,
	createdAt time.Time,
	updatedAt time.Time,
) AssignmentInfo {
	return AssignmentInfo{
		ID:                    id,
		ComicID:               comicID,
		UserID:                userID,
		AssignedRawProviderAt: assignedRawProviderAt,
		AssignedTranslatorAt:  assignedTranslatorAt,
		AssignedProofreaderAt: assignedProofreaderAt,
		AssignedTypesetterAt:  assignedTypesetterAt,
		AssignedReviewerAt:    assignedReviewerAt,
		AssignedPublisherAt:   assignedPublisherAt,
		CreatedAt:             createdAt,
		UpdatedAt:             updatedAt,
	}
}

func (a *AssignmentInfo) HasAnyRole(roles ...RoleFlag) bool {
	for _, role := range roles {
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
