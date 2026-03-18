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
	ToBePublisher   bool `gorm:"column:to_be_publisher"`
	ToBeAdmin       bool `gorm:"column:to_be_admin"`

	Pending bool `gorm:"column:pending"`

	// Invitor 别名列（IncludeInvitorInfo() 时填充）
	InvitorName             string    `gorm:"column:invitor_name"`
	InvitorQQ               string    `gorm:"column:invitor_qq"`
	InvitorAvatarOSSKey     string    `gorm:"column:invitor_avatar_oss_key"`
	InvitorIsAvatarUploaded bool      `gorm:"column:invitor_is_avatar_uploaded"`
	InvitorIsSuperAdmin     bool      `gorm:"column:invitor_is_super_admin"`
	InvitorCreatedAt        time.Time `gorm:"column:invitor_created_at"`
	InvitorUpdatedAt        time.Time `gorm:"column:invitor_updated_at"`

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
	result := model.InvitationInfo{
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
		ToBePublisher:   row.ToBePublisher,
		ToBeAdmin:       row.ToBeAdmin,
		CreatedAt:       row.CreatedAt,
	}

	if !row.InvitorCreatedAt.IsZero() {
		invitorInfo := model.NewUserInfo(
			row.InvitorID,
			row.InvitorName,
			row.InvitorQQ,
			row.InvitorAvatarOSSKey,
			row.InvitorIsAvatarUploaded,
			row.InvitorIsSuperAdmin,
			row.InvitorCreatedAt,
			row.InvitorUpdatedAt,
		)
		result.Invitor = &invitorInfo
	}

	return result
}
