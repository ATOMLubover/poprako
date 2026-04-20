package repo_entity

import (
	"time"

	"poprako-s/internal/domain/model"
)

const MemberTable = "member_table"

type MemberInfoRow struct {
	ID string `gorm:"column:id"`

	UserID string `gorm:"column:user_id"`
	TeamID string `gorm:"column:team_id"`

	AssignedRawProviderAt *time.Time `gorm:"column:assigned_raw_provider_at"`
	AssignedTranslatorAt  *time.Time `gorm:"column:assigned_translator_at"`
	AssignedProofreaderAt *time.Time `gorm:"column:assigned_proofreader_at"`
	AssignedTypesetterAt  *time.Time `gorm:"column:assigned_typesetter_at"`
	AssignedRedrawerAt    *time.Time `gorm:"column:assigned_redrawer_at"`
	AssignedReviewerAt    *time.Time `gorm:"column:assigned_reviewer_at"`
	AssignedPublisherAt   *time.Time `gorm:"column:assigned_publisher_at"`
	AssignedAdminAt       *time.Time `gorm:"column:assigned_admin_at"`

	CreatedAt time.Time  `gorm:"column:created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at"`
}

func ToMemberInfo(row MemberInfoRow) model.MemberInfo {
	return model.MemberInfo{
		ID:                    row.ID,
		UserID:                row.UserID,
		TeamID:                row.TeamID,
		AssignedRawProviderAt: row.AssignedRawProviderAt,
		AssignedTranslatorAt:  row.AssignedTranslatorAt,
		AssignedProofreaderAt: row.AssignedProofreaderAt,
		AssignedTypesetterAt:  row.AssignedTypesetterAt,
		AssignedRedrawerAt:    row.AssignedRedrawerAt,
		AssignedReviewerAt:    row.AssignedReviewerAt,
		AssignedPublisherAt:   row.AssignedPublisherAt,
		AssignedAdminAt:       row.AssignedAdminAt,
		CreatedAt:             row.CreatedAt,
		UpdatedAt:             row.UpdatedAt,
	}
}
