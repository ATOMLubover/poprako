package repo_entity

import (
	"time"

	"poprako-s/internal/domain/model"
)

const ChapterInvitationTable = "chapter_invitation"

type ChapterInvitationInfoRow struct {
	ID string `gorm:"column:id"`

	ChapterID string `gorm:"column:chapter_id"`
	InviterID string `gorm:"column:inviter_id"`
	InviteeQQ string `gorm:"column:invitee_qq"`

	InvitationCode string `gorm:"column:invitation_code"`
	Pending        bool   `gorm:"column:pending"`

	ToBeRawProvider bool `gorm:"column:to_be_raw_provider"`
	ToBeTranslator  bool `gorm:"column:to_be_translator"`
	ToBeProofreader bool `gorm:"column:to_be_proofreader"`
	ToBeTypesetter  bool `gorm:"column:to_be_typesetter"`
	ToBeRedrawer    bool `gorm:"column:to_be_redrawer"`
	ToBeReviewer    bool `gorm:"column:to_be_reviewer"`
	ToBePublisher   bool `gorm:"column:to_be_publisher"`

	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func ToChapterInvitationInfo(row ChapterInvitationInfoRow) model.ChapterInvitationInfo {
	return model.ChapterInvitationInfo{
		ID:              row.ID,
		ChapterID:       row.ChapterID,
		InviterID:       row.InviterID,
		InviteeQQ:       row.InviteeQQ,
		InvitationCode:  row.InvitationCode,
		Pending:         row.Pending,
		ToBeRawProvider: row.ToBeRawProvider,
		ToBeTranslator:  row.ToBeTranslator,
		ToBeProofreader: row.ToBeProofreader,
		ToBeTypesetter:  row.ToBeTypesetter,
		ToBeRedrawer:    row.ToBeRedrawer,
		ToBeReviewer:    row.ToBeReviewer,
		ToBePublisher:   row.ToBePublisher,
		CreatedAt:       row.CreatedAt,
		UpdatedAt:       row.UpdatedAt,
	}
}
