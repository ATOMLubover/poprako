package entity

import (
	"time"

	"labelplus-next-web-be/internal/domain/model"
)

const MemberTable = "member_table"

// MemberProfileRow 用于 List、GetByID，映射成员信息。
type MemberProfileRow struct {
	ID     string `gorm:"column:id"`
	UserID string `gorm:"column:user_id"`
	TeamID string `gorm:"column:team_id"`

	AssignedRawProviderAt *time.Time `gorm:"column:assigned_raw_provider_at"`
	AssignedTranslatorAt  *time.Time `gorm:"column:assigned_translator_at"`
	AssignedProofreaderAt *time.Time `gorm:"column:assigned_proofreader_at"`
	AssignedTypesetterAt  *time.Time `gorm:"column:assigned_typesetter_at"`
	AssignedReviewerAt    *time.Time `gorm:"column:assigned_reviewer_at"`
	AssignedPublisherAt   *time.Time `gorm:"column:assigned_publisher_at"`
	AssignedAdminAt       *time.Time `gorm:"column:assigned_admin_at"`

	CreatedAt time.Time  `gorm:"column:created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at"`
}

func (MemberProfileRow) TableName() string { return MemberTable }

// MemberInsertRow 用于 Create，仅包含写入所需字段。
type MemberInsertRow struct {
	ID     string `gorm:"column:id"`
	UserID string `gorm:"column:user_id"`
	TeamID string `gorm:"column:team_id"`

	AssignedRawProviderAt *time.Time `gorm:"column:assigned_raw_provider_at"`
	AssignedTranslatorAt  *time.Time `gorm:"column:assigned_translator_at"`
	AssignedProofreaderAt *time.Time `gorm:"column:assigned_proofreader_at"`
	AssignedTypesetterAt  *time.Time `gorm:"column:assigned_typesetter_at"`
	AssignedReviewerAt    *time.Time `gorm:"column:assigned_reviewer_at"`
	AssignedPublisherAt   *time.Time `gorm:"column:assigned_publisher_at"`
	AssignedAdminAt       *time.Time `gorm:"column:assigned_admin_at"`
}

func (MemberInsertRow) TableName() string { return MemberTable }

type MemberWithUserRow struct {
	MemberProfileRow
	UserName             string    `gorm:"column:user_name"`
	UserQQ               string    `gorm:"column:user_qq"`
	UserAvatarOSSKey     string    `gorm:"column:user_avatar_oss_key"`
	UserIsAvatarUploaded bool      `gorm:"column:user_is_avatar_uploaded"`
	UserIsSuperAdmin     bool      `gorm:"column:user_is_super_admin"`
	UserCreatedAt        time.Time `gorm:"column:user_created_at"`
	UserUpdatedAt        time.Time `gorm:"column:user_updated_at"`
}

type MemberWithTeamRow struct {
	MemberProfileRow
	TeamName             string    `gorm:"column:team_name"`
	TeamDescription      string    `gorm:"column:team_description"`
	TeamAvatarOSSKey     string    `gorm:"column:team_avatar_oss_key"`
	TeamIsAvatarUploaded bool      `gorm:"column:team_is_avatar_uploaded"`
	TeamCreatedAt        time.Time `gorm:"column:team_created_at"`
	TeamUpdatedAt        time.Time `gorm:"column:team_updated_at"`
}

func ToMemberProfile(row MemberProfileRow, userInfo model.UserInfo) model.MemberWithUserInfo {
	return model.MemberWithUserInfo{
		ID:                    row.ID,
		UserInfo:              userInfo,
		TeamID:                row.TeamID,
		AssignedRawProviderAt: row.AssignedRawProviderAt,
		AssignedTranslatorAt:  row.AssignedTranslatorAt,
		AssignedProofreaderAt: row.AssignedProofreaderAt,
		AssignedTypesetterAt:  row.AssignedTypesetterAt,
		AssignedReviewerAt:    row.AssignedReviewerAt,
		AssignedPublishererAt: row.AssignedPublisherAt,
		AssignedAdminAt:       row.AssignedAdminAt,
		CreatedAt:             row.CreatedAt,
		UpdatedAt:             row.UpdatedAt,
	}
}

func ToMemberWithTeamInfo(row MemberWithTeamRow) model.MemberWithTeamInfo {
	team := model.TeamInfo{
		ID:               row.TeamID,
		Name:             row.TeamName,
		Description:      row.TeamDescription,
		AvatarOSSKey:     row.TeamAvatarOSSKey,
		IsAvatarUploaded: row.TeamIsAvatarUploaded,
		CreatedAt:        row.TeamCreatedAt,
		UpdatedAt:        row.TeamUpdatedAt,
	}

	return model.NewMemberWithTeamInfo(
		row.ID,
		row.UserID,
		team,
		row.AssignedRawProviderAt,
		row.AssignedTranslatorAt,
		row.AssignedProofreaderAt,
		row.AssignedTypesetterAt,
		row.AssignedReviewerAt,
		row.AssignedPublisherAt,
		row.AssignedAdminAt,
	)
}
