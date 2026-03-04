package entity

import (
	"time"

	"labelplus-next-web-be/internal/domain/model"
)

type MemberRow struct {
	ID     string `gorm:"column:id"`
	UserID string `gorm:"column:user_id"`
	TeamID string `gorm:"column:team_id"`

	AssignedRawProviderAt *time.Time `gorm:"column:assigned_raw_provider_at"`
	AssignedTranslatorAt  *time.Time `gorm:"column:assigned_translator_at"`
	AssignedProofreaderAt *time.Time `gorm:"column:assigned_proofreader_at"`
	AssignedTypesetterAt  *time.Time `gorm:"column:assigned_typesetter_at"`
	AssignedReviewerAt    *time.Time `gorm:"column:assigned_reviewer_at"`
	AssignedUploaderAt    *time.Time `gorm:"column:assigned_uploader_at"`
	AssignedAdminAt       *time.Time `gorm:"column:assigned_admin_at"`

	CreatedAt time.Time  `gorm:"column:created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at"`
}

func (MemberRow) TableName() string { return "member_table" }

type MemberWithUserRow struct {
	MemberRow
	UserName         string    `gorm:"column:user_name"`
	UserQQ           string    `gorm:"column:user_qq"`
	UserAvatarURL    string    `gorm:"column:user_avatar_url"`
	UserIsSuperAdmin bool      `gorm:"column:user_is_super_admin"`
	UserCreatedAt    time.Time `gorm:"column:user_created_at"`
	UserUpdatedAt    time.Time `gorm:"column:user_updated_at"`
}

func ToMemberProfile(row MemberRow, userInfo *model.UserInfo) model.MemberProfile {
	return model.MemberProfile{
		ID:                    row.ID,
		UserInfo:              userInfo,
		TeamID:                row.TeamID,
		AssignedRawProviderAt: row.AssignedRawProviderAt,
		AssignedTranslatorAt:  row.AssignedTranslatorAt,
		AssignedProofreaderAt: row.AssignedProofreaderAt,
		AssignedTypesetterAt:  row.AssignedTypesetterAt,
		AssignedReviewerAt:    row.AssignedReviewerAt,
		AssignedUploaderAt:    row.AssignedUploaderAt,
		AssignedAdminAt:       row.AssignedAdminAt,
		CreatedAt:             row.CreatedAt,
		UpdatedAt:             row.UpdatedAt,
	}
}
