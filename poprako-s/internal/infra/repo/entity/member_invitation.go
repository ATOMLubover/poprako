package repo_entity

import (
	"time"

	"poprako-s/internal/domain/model"
)

const MemberInvitationTable = "member_invitation_table"

type MemberInvitationInfoRow struct {
	ID string `gorm:"column:id"`

	InvitorID string `gorm:"column:invitor_id"`
	TeamID    string `gorm:"column:team_id"`

	InviteeQQ      string `gorm:"column:invitee_qq"`
	InvitationCode string `gorm:"column:invitation_code"`

	ToBeRawProvider bool `gorm:"column:to_be_raw_provider"`
	ToBeTranslator  bool `gorm:"column:to_be_translator"`
	ToBeProofreader bool `gorm:"column:to_be_proofreader"`
	ToBeTypesetter  bool `gorm:"column:to_be_typesetter"`
	ToBeReviewer    bool `gorm:"column:to_be_reviewer"`
	ToBePublisher   bool `gorm:"column:to_be_publisher"`
	ToBeAdmin       bool `gorm:"column:to_be_admin"`

	Pending bool `gorm:"column:pending"`

	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func ToMemberInvitationInfo(row MemberInvitationInfoRow) model.MemberInvitationInfo {
	return model.MemberInvitationInfo{
		ID:              row.ID,
		InvitorID:       row.InvitorID,
		InviteeQQ:       row.InviteeQQ,
		TeamID:          row.TeamID,
		InvitationCode:  row.InvitationCode,
		Pending:         row.Pending,
		ToBeRawProvider: row.ToBeRawProvider,
		ToBeTranslator:  row.ToBeTranslator,
		ToBeProofreader: row.ToBeProofreader,
		ToBeTypesetter:  row.ToBeTypesetter,
		ToBeReviewer:    row.ToBeReviewer,
		ToBePublisher:   row.ToBePublisher,
		ToBeAdmin:       row.ToBeAdmin,
		CreatedAt:       row.CreatedAt,
	}
}
