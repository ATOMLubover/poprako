package entity

import (
	"time"

	"poprako-s/internal/domain/model/aggr"
)

const MEMBER_TABLE = "t_member"

// `MemberUserNicknameCol` is the `user_nickname` column name in `t_member`.
const MemberUserNicknameCol = "user_nickname"

type MemberRow struct {
	Id string `gorm:"column:id;primaryKey"`

	UserId       string   `gorm:"column:user_id"`
	UserNickname string   `gorm:"column:user_nickname"`
	User         *UserRow `gorm:"foreignKey:UserId"`

	TeamId string   `gorm:"column:team_id"`
	Team   *TeamRow `gorm:"foreignKey:TeamId"`

	AssignedRawProviderAt *time.Time `gorm:"column:assigned_raw_provider_at"`
	AssignedTranslatorAt  *time.Time `gorm:"column:assigned_translator_at"`
	AssignedProofreaderAt *time.Time `gorm:"column:assigned_proofreader_at"`
	AssignedTypesetterAt  *time.Time `gorm:"column:assigned_typesetter_at"`
	AssignedRedrawerAt    *time.Time `gorm:"column:assigned_redrawer_at"`
	AssignedReviewerAt    *time.Time `gorm:"column:assigned_reviewer_at"`
	AssignedPublisherAt   *time.Time `gorm:"column:assigned_publisher_at"`
	AssignedAdminAt       *time.Time `gorm:"column:assigned_admin_at"`

	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (*MemberRow) TableName() string {
	return MEMBER_TABLE
}

func (r *MemberRow) ToMemberAggr() *aggr.Member {
	if r == nil {
		// Do not treat nil as an error, as it may be used in includes.
		return nil
	}

	return &aggr.Member{
		TimedRoles: aggr.TimedRoles{
			AssignedRawProviderAt: r.AssignedRawProviderAt,
			AssignedTranslatorAt:  r.AssignedTranslatorAt,
			AssignedProofreaderAt: r.AssignedProofreaderAt,
			AssignedTypesetterAt:  r.AssignedTypesetterAt,
			AssignedRedrawerAt:    r.AssignedRedrawerAt,
			AssignedReviewerAt:    r.AssignedReviewerAt,
			AssignedPublisherAt:   r.AssignedPublisherAt,
			AssignedAdminAt:       r.AssignedAdminAt,
		},
		Id:           r.Id,
		UserId:       r.UserId,
		UserNickname: r.UserNickname,
		User:         r.User.ToUserAggr(),
		TeamId:       r.TeamId,
		Team:         r.Team.ToTeamAggr(),
		CreatedAt:    r.CreatedAt,
		UpdatedAt:    r.UpdatedAt,
	}
}

type MemberCreRow struct {
	Id string `gorm:"column:id;primaryKey"`

	UserId       string `gorm:"column:user_id"`
	UserNickname string `gorm:"column:user_nickname"`
	TeamId       string `gorm:"column:team_id"`

	AssignedRawProviderAt *time.Time `gorm:"column:assigned_raw_provider_at"`
	AssignedTranslatorAt  *time.Time `gorm:"column:assigned_translator_at"`
	AssignedProofreaderAt *time.Time `gorm:"column:assigned_proofreader_at"`
	AssignedTypesetterAt  *time.Time `gorm:"column:assigned_typesetter_at"`
	AssignedRedrawerAt    *time.Time `gorm:"column:assigned_redrawer_at"`
	AssignedReviewerAt    *time.Time `gorm:"column:assigned_reviewer_at"`
	AssignedPublisherAt   *time.Time `gorm:"column:assigned_publisher_at"`
	AssignedAdminAt       *time.Time `gorm:"column:assigned_admin_at"`

	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func NewMemberCreRowFromAggr(cre *aggr.MemberCre) *MemberCreRow {
	roles := aggr.TimedRoles{}
	now := time.Now()

	roles.FromRoleMask(cre.RoleMask)

	return &MemberCreRow{
		Id: cre.Id,

		UserId:       cre.UserId,
		UserNickname: cre.UserNickname,
		TeamId:       cre.TeamId,

		AssignedRawProviderAt: roles.AssignedRawProviderAt,
		AssignedTranslatorAt:  roles.AssignedTranslatorAt,
		AssignedProofreaderAt: roles.AssignedProofreaderAt,
		AssignedTypesetterAt:  roles.AssignedTypesetterAt,
		AssignedRedrawerAt:    roles.AssignedRedrawerAt,
		AssignedReviewerAt:    roles.AssignedReviewerAt,
		AssignedPublisherAt:   roles.AssignedPublisherAt,
		AssignedAdminAt:       roles.AssignedAdminAt,

		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (*MemberCreRow) TableName() string {
	return MEMBER_TABLE
}

type memberRoleUpdRow struct {
	Id string `gorm:"column:id;primaryKey"`

	AssignedRawProviderAt *time.Time `gorm:"column:assigned_raw_provider_at"`
	AssignedTranslatorAt  *time.Time `gorm:"column:assigned_translator_at"`
	AssignedProofreaderAt *time.Time `gorm:"column:assigned_proofreader_at"`
	AssignedTypesetterAt  *time.Time `gorm:"column:assigned_typesetter_at"`
	AssignedRedrawerAt    *time.Time `gorm:"column:assigned_redrawer_at"`
	AssignedReviewerAt    *time.Time `gorm:"column:assigned_reviewer_at"`
	AssignedPublisherAt   *time.Time `gorm:"column:assigned_publisher_at"`
	AssignedAdminAt       *time.Time `gorm:"column:assigned_admin_at"`

	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func NewMemberRoleUpdRowFromAggr(upd *aggr.MemberRoleUpd) *memberRoleUpdRow {
	roles := aggr.TimedRoles{}
	now := time.Now()

	roles.FromRoleMask(upd.RoleMask)

	return &memberRoleUpdRow{
		Id: upd.Id,

		AssignedRawProviderAt: roles.AssignedRawProviderAt,
		AssignedTranslatorAt:  roles.AssignedTranslatorAt,
		AssignedProofreaderAt: roles.AssignedProofreaderAt,
		AssignedTypesetterAt:  roles.AssignedTypesetterAt,
		AssignedRedrawerAt:    roles.AssignedRedrawerAt,
		AssignedReviewerAt:    roles.AssignedReviewerAt,
		AssignedPublisherAt:   roles.AssignedPublisherAt,
		AssignedAdminAt:       roles.AssignedAdminAt,

		UpdatedAt: now,
	}
}

func (*memberRoleUpdRow) TableName() string {
	return MEMBER_TABLE
}
