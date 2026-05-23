package repo_infra

import (
	"poprako-s/internal/domain/model/aggr"
	"poprako-s/internal/domain/model/enum"
	"poprako-s/internal/domain/model/query"
	repo_iface "poprako-s/internal/domain/repo"
	"poprako-s/internal/infra/repo/entity"

	"gorm.io/gorm"
)

// `commentRepoImpl` is the gorm implementation of `CommentRepo`.
type commentRepoImpl struct {
	gdb *gorm.DB
}

// `NewCommentRepo` creates a non transaction-scoped `CommentRepo`.
func NewCommentRepo(gdb *gorm.DB) repo_iface.CommentRepo {
	return &commentRepoImpl{gdb: gdb}
}

// `List` returns team comments by filter options.
func (r *commentRepoImpl) List(opt query.ListCommentOpt, inc ...enum.CommentIncl) ([]aggr.Comment, repo_iface.RepoErr) {
	var rows []entity.CommentRow

	qry := r.gdb.
		Table(entity.COMMENT_TABLE).
		Where("team_id = ?", opt.TeamId)

	if opt.Pagi.Offset > 0 {
		qry = qry.Offset(opt.Pagi.Offset)
	}

	if opt.Pagi.Limit > 0 {
		qry = qry.Limit(opt.Pagi.Limit)
	}

	qry = withCommentIncl(qry, inc...)

	err := qry.
		Order("created_at DESC").
		Find(&rows).Error
	if err != nil {
		return nil, err
	}

	items := make([]aggr.Comment, len(rows))
	for i := range rows {
		item := rows[i].ToCommentAggr()
		if item != nil {
			items[i] = *item
		}
	}

	return items, nil
}

// `Create` inserts one team comment record and returns a fully loaded aggregate.
func (r *commentRepoImpl) Create(cre *aggr.CommentCre) (*aggr.Comment, repo_iface.RepoErr) {
	creRow := entity.NewCommentCreRowFromAggr(cre)

	err := r.gdb.
		Table(creRow.TableName()).
		Create(creRow).Error
	if err != nil {
		return nil, err
	}

	var row entity.CommentRow

	qry := r.gdb.
		Table(entity.COMMENT_TABLE).
		Where("id = ?", creRow.Id)

	qry = withCommentIncl(qry, enum.CommentInclUser)

	err = qry.First(&row).Error
	if err != nil {
		return nil, err
	}

	return row.ToCommentAggr(), nil
}

// `withCommentIncl` maps typed include options to preloads.
func withCommentIncl(query *gorm.DB, inc ...enum.CommentIncl) *gorm.DB {
	for _, i := range inc {
		switch i {
		case enum.CommentInclUser:
			query = query.Preload("User")
		}
	}

	return query
}
