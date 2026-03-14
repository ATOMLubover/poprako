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

// MemberWithInfoRow 是统一的聚合行类型，可同时携带 user 和 team 别名列。
// User 列与 Team 列均可选，由 query option 决定是否 JOIN。
type MemberWithInfoRow struct {
	MemberProfileRow

	// User 别名列（IncludeUserInfo() 时填充）
	UserName             string    `gorm:"column:user_name"`
	UserQQ               string    `gorm:"column:user_qq"`
	UserAvatarOSSKey     string    `gorm:"column:user_avatar_oss_key"`
	UserIsAvatarUploaded bool      `gorm:"column:user_is_avatar_uploaded"`
	UserIsSuperAdmin     bool      `gorm:"column:user_is_super_admin"`
	UserCreatedAt        time.Time `gorm:"column:user_created_at"`
	UserUpdatedAt        time.Time `gorm:"column:user_updated_at"`

	// Team 别名列（IncludeTeamInfo() 时填充）
	TeamName             string    `gorm:"column:team_name"`
	TeamDescription      string    `gorm:"column:team_description"`
	TeamAvatarOSSKey     string    `gorm:"column:team_avatar_oss_key"`
	TeamIsAvatarUploaded bool      `gorm:"column:team_is_avatar_uploaded"`
	TeamCreatedAt        time.Time `gorm:"column:team_created_at"`
	TeamUpdatedAt        time.Time `gorm:"column:team_updated_at"`
}

func ToMemberWithInfo(row MemberWithInfoRow) model.MemberWithInfo {
	var userInfo *model.UserInfo
	if !row.UserCreatedAt.IsZero() {
		info := model.NewUserInfo(
			row.UserID,
			row.UserName,
			row.UserQQ,
			row.UserAvatarOSSKey,
			row.UserIsAvatarUploaded,
			row.UserIsSuperAdmin,
			row.UserCreatedAt,
			row.UserUpdatedAt,
		)
		userInfo = &info
	}

	var teamInfo *model.TeamInfo
	if !row.TeamCreatedAt.IsZero() {
		info := model.NewTeamInfo(
			row.TeamID,
			row.TeamName,
			row.TeamDescription,
			row.TeamAvatarOSSKey,
			row.TeamIsAvatarUploaded,
			row.TeamCreatedAt,
			row.TeamUpdatedAt,
		)
		teamInfo = &info
	}

	return model.NewMemberWithInfo(
		row.ID,
		row.UserID,
		row.TeamID,
		userInfo,
		teamInfo,
		row.AssignedRawProviderAt,
		row.AssignedTranslatorAt,
		row.AssignedProofreaderAt,
		row.AssignedTypesetterAt,
		row.AssignedReviewerAt,
		row.AssignedPublisherAt,
		row.AssignedAdminAt,
		row.CreatedAt,
		row.UpdatedAt,
	)
}
