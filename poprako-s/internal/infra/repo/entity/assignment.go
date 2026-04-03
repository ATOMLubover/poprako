package repo_entity

import (
	"time"

	"poprako-s/internal/domain/model"
)

const AssignmentTable = "assignment_table"

type AssignmentInfoRow struct {
	ID string `gorm:"column:id"`

	ChapterID string `gorm:"column:chapter_id"`
	UserID    string `gorm:"column:user_id"`

	AssignedRawProviderAt *time.Time `gorm:"column:assigned_raw_provider_at"`
	AssignedTranslatorAt  *time.Time `gorm:"column:assigned_translator_at"`
	AssignedProofreaderAt *time.Time `gorm:"column:assigned_proofreader_at"`
	AssignedTypesetterAt  *time.Time `gorm:"column:assigned_typesetter_at"`
	AssignedRedrawerAt    *time.Time `gorm:"column:assigned_redrawer_at"`
	AssignedReviewerAt    *time.Time `gorm:"column:assigned_reviewer_at"`
	AssignedPublisherAt   *time.Time `gorm:"column:assigned_publisher_at"`

	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func ToAssignmentInfo(row AssignmentInfoRow) model.AssignmentInfo {
	return model.AssignmentInfo{
		ID:                    row.ID,
		ChapterID:             row.ChapterID,
		UserID:                row.UserID,
		AssignedRawProviderAt: row.AssignedRawProviderAt,
		AssignedTranslatorAt:  row.AssignedTranslatorAt,
		AssignedProofreaderAt: row.AssignedProofreaderAt,
		AssignedTypesetterAt:  row.AssignedTypesetterAt,
		AssignedRedrawerAt:    row.AssignedRedrawerAt,
		AssignedReviewerAt:    row.AssignedReviewerAt,
		AssignedPublisherAt:   row.AssignedPublisherAt,
		CreatedAt:             row.CreatedAt,
		UpdatedAt:             row.UpdatedAt,
	}
}
