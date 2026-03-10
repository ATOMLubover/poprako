package entity

import (
	"time"

	"labelplus-next-web-be/internal/domain/model"
)

const InvitationTable = "invitation_table"

// InvitationInfoRow 用于 List、GetByID、GetByInviteeQQAndCode。
type InvitationInfoRow struct {
	ID             string `gorm:"column:id"`
	InvitorID      string `gorm:"column:invitor_id"`
	TargetTeamID   string `gorm:"column:target_team_id"`
	InviteeQQ      string `gorm:"column:invitee_qq"`
	InvitationCode string `gorm:"column:invitation_code"`

	ToBeRawProvider bool `gorm:"column:to_be_raw_provider"`
	ToBeTranslator  bool `gorm:"column:to_be_translator"`
	ToBeProofreader bool `gorm:"column:to_be_proofreader"`
	ToBeTypesetter  bool `gorm:"column:to_be_typesetter"`
	ToBeReviewer    bool `gorm:"column:to_be_reviewer"`
	ToBeUploader    bool `gorm:"column:to_be_publisher"`
	ToBeAdmin       bool `gorm:"column:to_be_admin"`

	Pending bool `gorm:"column:pending"`

	CreatedAt time.Time `gorm:"column:created_at"`
}

func (InvitationInfoRow) TableName() string { return InvitationTable }

// InvitationInsertRow 用于 Create，仅包含写入所需字段。
type InvitationInsertRow struct {
	ID             string `gorm:"column:id"`
	InvitorID      string `gorm:"column:invitor_id"`
	TargetTeamID   string `gorm:"column:target_team_id"`
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
}

func (InvitationInsertRow) TableName() string { return InvitationTable }

func ToInvitationInfo(row InvitationInfoRow) model.InvitationInfo {
	return model.InvitationInfo{
		ID:              row.ID,
		InvitorID:       row.InvitorID,
		InviteeQQ:       row.InviteeQQ,
		TeamID:          row.TargetTeamID,
		InvitationCode:  row.InvitationCode,
		Pending:         row.Pending,
		ToBeRawProvider: row.ToBeRawProvider,
		ToBeTranslator:  row.ToBeTranslator,
		ToBeProofreader: row.ToBeProofreader,
		ToBeTypesetter:  row.ToBeTypesetter,
		ToBeReviewer:    row.ToBeReviewer,
		ToBeUploader:    row.ToBeUploader,
		ToBeAdmin:       row.ToBeAdmin,
		CreatedAt:       row.CreatedAt,
	}
}
