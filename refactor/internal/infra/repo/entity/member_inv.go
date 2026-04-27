package entity

import (
	"time"

	"poprako-s/internal/domain/model/aggr"
)

// `MEMBER_INV_TABLE` is the table name for member invitations
const MEMBER_INV_TABLE = "t_member_invitation"

// `MemberInvRow` maps one invitation record for read queries
type MemberInvRow struct {
	Id string `gorm:"column:id;primaryKey"`

	InvitorId string `gorm:"column:inviter_id"`
	TeamId    string `gorm:"column:team_id"`

	InviteeQid string `gorm:"column:invitee_qid"`
	InvCode    string `gorm:"column:invitation_code"`

	Pending bool `gorm:"column:pending"`

	RoleMask int64 `gorm:"column:role_mask"`

	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

// `TableName` returns the table name of `MemberInvRow`
func (*MemberInvRow) TableName() string {
	return MEMBER_INV_TABLE
}

// `ToMemberInvAggr` converts a row to `MemberInv` aggregate
func (r *MemberInvRow) ToMemberInvAggr() *aggr.MemberInv {
	if r == nil {
		// Keep nil-safe conversion for optional includes
		return nil
	}

	return &aggr.MemberInv{
		Id:         r.Id,
		InvitorId:  r.InvitorId,
		TeamId:     r.TeamId,
		InviteeQid: r.InviteeQid,
		InvCode:    r.InvCode,
		Pending:    r.Pending,
		RoleMask:   aggr.RoleMask(r.RoleMask),
		CreatedAt:  r.CreatedAt,
	}
}

// `MemberInvCreRow` maps one invitation record for create queries
type MemberInvCreRow struct {
	Id string `gorm:"column:id;primaryKey"`

	InvitorId string `gorm:"column:inviter_id"`
	TeamId    string `gorm:"column:team_id"`

	InviteeQid string `gorm:"column:invitee_qid"`
	InvCode    string `gorm:"column:invitation_code"`

	Pending bool `gorm:"column:pending"`

	RoleMask int64 `gorm:"column:role_mask"`

	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

// `NewMemberInvCreRowFromAggr` builds a create row from `MemberInvCre`
func NewMemberInvCreRowFromAggr(cre *aggr.MemberInvCre) *MemberInvCreRow {
	now := time.Now()

	return &MemberInvCreRow{
		Id:         cre.Id,
		InvitorId:  cre.InvitorId,
		TeamId:     cre.TeamId,
		InviteeQid: cre.InviteeQid,
		InvCode:    cre.InvCode,
		Pending:    true,
		RoleMask:   int64(cre.RoleMask),
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}

// `TableName` returns the table name of `MemberInvCreRow`
func (*MemberInvCreRow) TableName() string {
	return MEMBER_INV_TABLE
}

// `memberInvMarkCmplUpdRow` maps update columns for completion mark
type memberInvMarkCmplUpdRow struct {
	Pending   bool      `gorm:"column:pending"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

// `NewMemberInvMarkCmplUpdRow` builds update row for `MarkCmpl`
func NewMemberInvMarkCmplUpdRow() *memberInvMarkCmplUpdRow {
	return &memberInvMarkCmplUpdRow{
		Pending:   false,
		UpdatedAt: time.Now(),
	}
}
