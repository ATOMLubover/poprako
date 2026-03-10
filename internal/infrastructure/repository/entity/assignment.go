package entity

import (
	"time"

	"labelplus-next-web-be/internal/domain/model"
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
	AssignedReviewerAt    *time.Time `gorm:"column:assigned_reviewer_at"`
	AssignedPublisherAt   *time.Time `gorm:"column:assigned_publisher_at"`

	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (AssignmentInfoRow) TableName() string { return AssignmentTable }

func ToAssignmentInfo(row AssignmentInfoRow) model.AssignmentInfo {
	return model.NewAssignmentInfo(
		row.ID,
		row.ChapterID,
		row.UserID,
		row.AssignedRawProviderAt,
		row.AssignedTranslatorAt,
		row.AssignedProofreaderAt,
		row.AssignedTypesetterAt,
		row.AssignedReviewerAt,
		row.AssignedPublisherAt,
		row.CreatedAt,
		row.UpdatedAt,
	)
}
