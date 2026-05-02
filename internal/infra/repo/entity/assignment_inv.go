package entity

import (
	"time"

	"poprako-s/internal/domain/model/aggr"
)

// `ASSIGNMENT_INV_TABLE` is table name for assignment invitation.
const ASSIGNMENT_INV_TABLE = "t_assignment_invitation"

// `AssignmentInvRow` maps one assignment invitation record for read queries.
type AssignmentInvRow struct {
	Id string `gorm:"column:id;primaryKey"`

	ChapterId string `gorm:"column:chapter_id"`
	InviterId string `gorm:"column:inviter_id"`

	InviteeQid string `gorm:"column:invitee_qid"`
	InvCode    string `gorm:"column:invitation_code"`

	Pending bool `gorm:"column:pending"`

	RoleMask int64 `gorm:"column:role_mask"`

	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

// `TableName` returns table name for `AssignmentInvRow`.
func (*AssignmentInvRow) TableName() string {
	return ASSIGNMENT_INV_TABLE
}

// `ToAssignmentInvAggr` converts row to assignment invitation aggregate.
func (r *AssignmentInvRow) ToAssignmentInvAggr() *aggr.AssignmentInv {
	if r == nil {
		return nil
	}

	return &aggr.AssignmentInv{
		Id:         r.Id,
		ChapterId:  r.ChapterId,
		InviterId:  r.InviterId,
		InviteeQid: r.InviteeQid,
		InvCode:    r.InvCode,
		Pending:    r.Pending,
		RoleMask:   aggr.RoleMask(r.RoleMask),
		CreatedAt:  r.CreatedAt,
		UpdatedAt:  r.UpdatedAt,
	}
}

// `AssignmentInvCreRow` maps assignment invitation columns for create queries.
type AssignmentInvCreRow struct {
	Id string `gorm:"column:id;primaryKey"`

	ChapterId string `gorm:"column:chapter_id"`
	InviterId string `gorm:"column:inviter_id"`

	InviteeQid string `gorm:"column:invitee_qid"`
	InvCode    string `gorm:"column:invitation_code"`

	Pending bool `gorm:"column:pending"`

	RoleMask int64 `gorm:"column:role_mask"`

	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

// `TableName` returns table name for `AssignmentInvCreRow`.
func (*AssignmentInvCreRow) TableName() string {
	return ASSIGNMENT_INV_TABLE
}

// `NewAssignmentInvCreRowFromAggr` builds create row from assignment invitation aggregate.
func NewAssignmentInvCreRowFromAggr(cre *aggr.AssignmentInvCre) *AssignmentInvCreRow {
	now := time.Now()

	return &AssignmentInvCreRow{
		Id:         cre.Id,
		ChapterId:  cre.ChapterId,
		InviterId:  cre.InviterId,
		InviteeQid: cre.InviteeQid,
		InvCode:    cre.InvCode,
		Pending:    true,
		RoleMask:   int64(cre.RoleMask),
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}

// `assignmentInvMarkCompletedUpdRow` maps update columns for invitation completion.
type assignmentInvMarkCompletedUpdRow struct {
	Pending   bool      `gorm:"column:pending"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

// `NewAssignmentInvMarkCompletedUpdRow` builds update row for invitation completion.
func NewAssignmentInvMarkCompletedUpdRow() *assignmentInvMarkCompletedUpdRow {
	return &assignmentInvMarkCompletedUpdRow{
		Pending:   false,
		UpdatedAt: time.Now(),
	}
}
