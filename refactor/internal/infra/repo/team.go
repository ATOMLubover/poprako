package repo_infra

import (
	"poprako-s/internal/domain/model/aggr"
	repo_iface "poprako-s/internal/domain/repo"
	"poprako-s/internal/infra/repo/entity"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// `teamRepoImpl` implements `repo_iface.TeamRepo` with `gorm`.
type teamRepoImpl struct {
	gdb *gorm.DB
}

// `NewTeamRepo` creates a non-transaction-scoped `TeamRepo`.
func NewTeamRepo(gdb *gorm.DB) repo_iface.TeamRepo {
	return &teamRepoImpl{gdb: gdb}
}

// `GetById` retrieves one `Team` aggregate by its id.
func (r *teamRepoImpl) GetById(id string) (*aggr.Team, repo_iface.RepoErr) {
	var row entity.TeamRow

	// Query the team row by id from `t_team`.
	err := r.gdb.
		Table(entity.TEAM_TABLE).
		Where("id = ?", id).
		First(&row).Error
	if err != nil {
		return nil, err
	}

	// Convert the row into `aggr.Team`.
	return row.ToTeamAggr(), nil
}

// `IncrementWorksetNextIndex` allocates one next workset index from one team row.
func (r *teamRepoImpl) IncrementWorksetNextIndex(id string) (int, repo_iface.RepoErr) {
	type nextIndexRow struct {
		NextIndex int `gorm:"column:workset_next_index"`
	}

	var row nextIndexRow

	updRe := r.gdb.
		Table(entity.TEAM_TABLE).
		Where("id = ?", id).
		Select("workset_next_index").
		Clauses(clause.Returning{Columns: []clause.Column{{Name: "workset_next_index"}}}).
		Updates(map[string]any{"workset_next_index": gorm.Expr("workset_next_index + 1")}).
		Scan(&row)
	if updRe.Error != nil {
		return 0, updRe.Error
	}

	if updRe.RowsAffected == 0 {
		return 0, gorm.ErrRecordNotFound
	}

	return row.NextIndex - 1, nil
}
