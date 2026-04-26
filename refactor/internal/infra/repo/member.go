package repo_infra

import (
	"context"
	"errors"

	"poprako-s/internal/domain/model/aggr"
	"poprako-s/internal/domain/model/enum"
	"poprako-s/internal/domain/model/query"
	repo_iface "poprako-s/internal/domain/repo"
	"poprako-s/internal/infra/repo/entity"

	"gorm.io/gorm"
)

type memberRepoImpl struct {
	gdb *gorm.DB
}

func TxnMemberRepo(cx context.Context) (repo_iface.MemberRepo, error) {
	gdb := takeTxnGdb(cx)
	if gdb == nil {
		return nil, errors.New("[TxnMemberRepo] no transaction context found for MemberRepo")
	}

	return &memberRepoImpl{gdb: gdb}, nil
}

func NewMemberRepo(gdb *gorm.DB) repo_iface.MemberRepo {
	return &memberRepoImpl{gdb: gdb}
}

func (r *memberRepoImpl) GetById(id string, inc ...enum.MemberIncl) (*aggr.Member, repo_iface.RepoErr) {
	var row entity.MemberRow

	query := r.gdb.
		Table(entity.MEMBER_TABLE).
		Where("id = ?", id)

	for _, i := range inc {
		switch i {
		case enum.MemberInclUser:
			query = query.Preload("User")
		case enum.MemberInclTeam:
			query = query.Preload("Team")
		}
	}

	err := query.First(&row).Error
	if err != nil {
		return nil, err
	}

	return row.ToMemberAggr(), nil
}

func (r *memberRepoImpl) List(opt *query.ListMemberOpt, inc ...enum.MemberIncl) ([]*aggr.Member, repo_iface.RepoErr) {
	var rows []entity.MemberRow

	query := r.gdb.
		Table(entity.MEMBER_TABLE)

	if opt.UserId != nil {
		query = query.Where("user_id = ?", *opt.UserId)
	}
	if opt.TeamId != nil {
		query = query.Where("team_id = ?", *opt.TeamId)
	}

	for _, i := range inc {
		switch i {
		case enum.MemberInclUser:
			query = query.Preload("User")
		case enum.MemberInclTeam:
			query = query.Preload("Team")
		}
	}

	err := query.Find(&rows).Error
	if err != nil {
		return nil, err
	}

	members := make([]*aggr.Member, len(rows))
	for i, row := range rows {
		members[i] = row.ToMemberAggr()
	}

	return members, nil
}

func (r *memberRepoImpl) ExistByUserTeamId(userId string, teamId string) (bool, repo_iface.RepoErr) {
	var count int64

	err := r.gdb.
		Table(entity.MEMBER_TABLE).
		Where("user_id = ? AND team_id = ?", userId, teamId).
		Count(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *memberRepoImpl) Create(cre *aggr.MemberCre) (*aggr.Member, repo_iface.RepoErr) {
	creRow := entity.NewMemberCreRowFromAggr(cre)

	err := r.gdb.
		Table(creRow.TableName()).
		Create(creRow).Error
	if err != nil {
		return nil, err
	}

	return r.GetById(creRow.Id)
}

func (r *memberRepoImpl) UpdateRoles(upd *aggr.MemberRoleUpd) repo_iface.RepoErr {
	updRow := entity.NewMemberRoleUpdRowFromAggr(upd)

	err := r.gdb.
		Table(entity.MEMBER_TABLE).
		Where("id = ?", updRow.Id).
		Updates(updRow).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *memberRepoImpl) Delete(id string) repo_iface.RepoErr {
	err := r.gdb.
		Table(entity.MEMBER_TABLE).
		Where("id = ?", id).
		Delete(&entity.MemberRow{}).Error
	if err != nil {
		return err
	}

	return nil
}
