package repo_infra

import (
	"poprako-s/internal/domain/model/aggr"
	"poprako-s/internal/domain/model/enum"
	"poprako-s/internal/domain/model/query"
	repo_iface "poprako-s/internal/domain/repo"
	"poprako-s/internal/infra/repo/entity"

	"gorm.io/gorm"
)

// `memberInvRepoImpl` is the gorm implementation of `MemberInvRepo`
type memberInvRepoImpl struct {
	gdb *gorm.DB
}

// `NewMemberInvRepo` creates a non transaction-scoped `MemberInvRepo`
func NewMemberInvRepo(gdb *gorm.DB) repo_iface.MemberInvRepo {
	return &memberInvRepoImpl{gdb: gdb}
}

// `GetPendingByInviteeQid` returns the latest invitation by invitee qid
func (r *memberInvRepoImpl) GetPendingByInviteeQid(qid string, inc ...enum.MemberInvIncl) (*aggr.MemberInv, repo_iface.RepoErr) {
	var row entity.MemberInvRow

	query := r.gdb.
		Table(entity.MEMBER_INV_TABLE).
		Where("invitee_qid = ? AND pending = true", qid).
		Order("created_at DESC")

	query = withMemberInvIncl(query, inc...)

	err := query.First(&row).Error
	if err != nil {
		return nil, err
	}

	return row.ToMemberInvAggr(), nil
}

// `GetById` returns one invitation by id.
func (r *memberInvRepoImpl) GetById(id string, inc ...enum.MemberInvIncl) (*aggr.MemberInv, repo_iface.RepoErr) {
	var row entity.MemberInvRow

	query := r.gdb.
		Table(entity.MEMBER_INV_TABLE).
		Where("id = ?", id)

	query = withMemberInvIncl(query, inc...)

	err := query.First(&row).Error
	if err != nil {
		return nil, err
	}

	return row.ToMemberInvAggr(), nil
}

// `List` returns invitations by filter options
func (r *memberInvRepoImpl) List(opt query.ListMemberInvOpt, inc ...enum.MemberInvIncl) ([]aggr.MemberInv, repo_iface.RepoErr) {
	var rows []entity.MemberInvRow

	qry := r.gdb.
		Table(entity.MEMBER_INV_TABLE).
		Where("team_id = ?", opt.TeamId)

	if opt.Pending != nil {
		qry = qry.Where("pending = ?", *opt.Pending)
	}

	if opt.Pagi.Offset > 0 {
		qry = qry.Offset(opt.Pagi.Offset)
	}

	if opt.Pagi.Limit > 0 {
		qry = qry.Limit(opt.Pagi.Limit)
	}

	qry = withMemberInvIncl(qry, inc...)

	err := qry.
		Order("created_at DESC").
		Find(&rows).Error
	if err != nil {
		return nil, err
	}

	items := make([]aggr.MemberInv, len(rows))
	for i := range rows {
		item := rows[i].ToMemberInvAggr()
		if item != nil {
			items[i] = *item
		}
	}

	return items, nil
}

// `withMemberInvIncl` maps typed include options to preloads.
func withMemberInvIncl(query *gorm.DB, inc ...enum.MemberInvIncl) *gorm.DB {
	for _, i := range inc {
		switch i {
		case enum.MemberInvInclInvitor:
			query = query.Preload("Invitor")

		case enum.MemberInvInclInvitee:
			query = query.Preload("Invitee")
		}
	}

	return query
}

// `Create` inserts a member invitation record
func (r *memberInvRepoImpl) Create(cre *aggr.MemberInvCre) (*aggr.MemberInv, repo_iface.RepoErr) {
	creRow := entity.NewMemberInvCreRowFromAggr(cre)

	err := r.gdb.
		Table(creRow.TableName()).
		Create(creRow).Error
	if err != nil {
		return nil, err
	}

	var row entity.MemberInvRow
	err = r.gdb.
		Table(entity.MEMBER_INV_TABLE).
		Where("id = ?", creRow.Id).
		First(&row).Error
	if err != nil {
		return nil, err
	}

	return row.ToMemberInvAggr(), nil
}

// `Update` applies put-style role update to one invitation row.
func (r *memberInvRepoImpl) Update(upd *aggr.MemberInvUpd) repo_iface.RepoErr {
	updRow := entity.NewMemberInvUpdRowFromAggr(upd)

	updRe := r.gdb.
		Table(entity.MEMBER_INV_TABLE).
		Where("id = ?", upd.Id).
		Select("role_mask", "updated_at").
		Updates(updRow)
	if updRe.Error != nil {
		return updRe.Error
	}

	if updRe.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// `Delete` hard deletes one invitation record
func (r *memberInvRepoImpl) Delete(id string) repo_iface.RepoErr {
	queryRe := r.gdb.
		Table(entity.MEMBER_INV_TABLE).
		Where("id = ?", id).
		Delete(&entity.MemberInvRow{})
	if queryRe.Error != nil {
		return queryRe.Error
	}

	if queryRe.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// `MarkCompleted` marks one invitation record as completed
func (r *memberInvRepoImpl) MarkCompleted(id string) repo_iface.RepoErr {
	upd := entity.NewMemberInvMarkCompletedUpdRow()

	updRe := r.gdb.
		Table(entity.MEMBER_INV_TABLE).
		Where("id = ?", id).
		Select("pending", "updated_at").
		Updates(upd)
	if updRe.Error != nil {
		return updRe.Error
	}

	if updRe.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
