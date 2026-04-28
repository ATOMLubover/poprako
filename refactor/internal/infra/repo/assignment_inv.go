package repo_infra

import (
	"poprako-s/internal/domain/model/aggr"
	"poprako-s/internal/domain/model/query"
	repo_iface "poprako-s/internal/domain/repo"
	"poprako-s/internal/infra/repo/entity"

	"gorm.io/gorm"
)

// `assignmentInvRepoImpl` is gorm implementation of `AssignmentInvRepo`.
type assignmentInvRepoImpl struct {
	gdb *gorm.DB
}

// `NewAssignmentInvRepo` creates a non transaction-scoped `AssignmentInvRepo`.
func NewAssignmentInvRepo(gdb *gorm.DB) repo_iface.AssignmentInvRepo {
	return &assignmentInvRepoImpl{gdb: gdb}
}

// `GetById` returns one invitation by id.
func (r *assignmentInvRepoImpl) GetById(id string) (*aggr.AssignmentInv, repo_iface.RepoErr) {
	var row entity.AssignmentInvRow

	err := r.gdb.
		Table(entity.ASSIGNMENT_INV_TABLE).
		Where("id = ?", id).
		First(&row).Error
	if err != nil {
		return nil, err
	}

	return row.ToAssignmentInvAggr(), nil
}

// `List` returns invitation list by chapter and optional pending.
func (r *assignmentInvRepoImpl) List(opt query.ListAssignmentInvOpt) ([]aggr.AssignmentInv, repo_iface.RepoErr) {
	var rows []entity.AssignmentInvRow

	qry := r.gdb.
		Table(entity.ASSIGNMENT_INV_TABLE).
		Where("chapter_id = ?", opt.ChapterId)

	if opt.Pending != nil {
		qry = qry.Where("pending = ?", *opt.Pending)
	}

	if opt.Pagi.Offset > 0 {
		qry = qry.Offset(opt.Pagi.Offset)
	}

	if opt.Pagi.Limit > 0 {
		qry = qry.Limit(opt.Pagi.Limit)
	}

	err := qry.
		Order("created_at DESC").
		Find(&rows).Error
	if err != nil {
		return nil, err
	}

	items := make([]aggr.AssignmentInv, len(rows))
	for i := range rows {
		inv := rows[i].ToAssignmentInvAggr()
		if inv != nil {
			items[i] = *inv
		}
	}

	return items, nil
}

// `ListPendingByInviteeQid` returns pending invitations by invitee qid.
func (r *assignmentInvRepoImpl) ListPendingByInviteeQid(qid string) ([]aggr.AssignmentInv, repo_iface.RepoErr) {
	var rows []entity.AssignmentInvRow

	err := r.gdb.
		Table(entity.ASSIGNMENT_INV_TABLE).
		Where("invitee_qid = ? AND pending = TRUE", qid).
		Order("created_at DESC").
		Find(&rows).Error
	if err != nil {
		return nil, err
	}

	items := make([]aggr.AssignmentInv, len(rows))
	for i := range rows {
		inv := rows[i].ToAssignmentInvAggr()
		if inv != nil {
			items[i] = *inv
		}
	}

	return items, nil
}

// `Create` inserts one invitation.
func (r *assignmentInvRepoImpl) Create(cre *aggr.AssignmentInvCre) (*aggr.AssignmentInv, repo_iface.RepoErr) {
	row := entity.NewAssignmentInvCreRowFromAggr(cre)

	err := r.gdb.
		Table(entity.ASSIGNMENT_INV_TABLE).
		Create(row).Error
	if err != nil {
		return nil, err
	}

	var created entity.AssignmentInvRow
	err = r.gdb.
		Table(entity.ASSIGNMENT_INV_TABLE).
		Where("id = ?", row.Id).
		First(&created).Error
	if err != nil {
		return nil, err
	}

	return created.ToAssignmentInvAggr(), nil
}

// `Delete` hard deletes one invitation.
func (r *assignmentInvRepoImpl) Delete(id string) repo_iface.RepoErr {
	return r.gdb.
		Table(entity.ASSIGNMENT_INV_TABLE).
		Where("id = ?", id).
		Delete(&entity.AssignmentInvRow{}).Error
}

// `MarkCompleted` marks one invitation as completed.
func (r *assignmentInvRepoImpl) MarkCompleted(id string) repo_iface.RepoErr {
	upd := entity.NewAssignmentInvMarkCompletedUpdRow()

	return r.gdb.
		Table(entity.ASSIGNMENT_INV_TABLE).
		Where("id = ?", id).
		Select("pending", "updated_at").
		Updates(upd).Error
}
