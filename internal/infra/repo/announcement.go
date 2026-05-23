package repo_infra

import (
	"poprako-s/internal/domain/model/aggr"
	"poprako-s/internal/domain/model/enum"
	"poprako-s/internal/domain/model/query"
	repo_iface "poprako-s/internal/domain/repo"
	"poprako-s/internal/infra/repo/entity"

	"gorm.io/gorm"
)

// `announcementRepoImpl` is the gorm implementation of `AnnouncementRepo`.
type announcementRepoImpl struct {
	gdb *gorm.DB
}

// `NewAnnouncementRepo` creates a non transaction-scoped `AnnouncementRepo`.
func NewAnnouncementRepo(gdb *gorm.DB) repo_iface.AnnouncementRepo {
	return &announcementRepoImpl{gdb: gdb}
}

// `List` returns team announcements by filter options.
func (r *announcementRepoImpl) List(opt query.ListAnnouncementOpt, inc ...enum.AnnouncementIncl) ([]aggr.Announcement, repo_iface.RepoErr) {
	var rows []entity.AnnouncementRow

	qry := r.gdb.
		Table(entity.ANNOUNCEMENT_TABLE).
		Where("team_id = ?", opt.TeamId)

	if opt.Pagi.Offset > 0 {
		qry = qry.Offset(opt.Pagi.Offset)
	}

	if opt.Pagi.Limit > 0 {
		qry = qry.Limit(opt.Pagi.Limit)
	}

	qry = withAnnouncementIncl(qry, inc...)

	err := qry.
		Order("created_at DESC").
		Find(&rows).Error
	if err != nil {
		return nil, err
	}

	items := make([]aggr.Announcement, len(rows))
	for i := range rows {
		item := rows[i].ToAnnouncementAggr()
		if item != nil {
			items[i] = *item
		}
	}

	return items, nil
}

// `Create` inserts one announcement record and returns a fully loaded aggregate.
func (r *announcementRepoImpl) Create(cre *aggr.AnnouncementCre) (*aggr.Announcement, repo_iface.RepoErr) {
	creRow := entity.NewAnnouncementCreRowFromAggr(cre)

	err := r.gdb.
		Table(creRow.TableName()).
		Create(creRow).Error
	if err != nil {
		return nil, err
	}

	var row entity.AnnouncementRow

	qry := r.gdb.
		Table(entity.ANNOUNCEMENT_TABLE).
		Where("id = ?", creRow.Id)

	qry = withAnnouncementIncl(qry, enum.AnnouncementInclUser)

	err = qry.First(&row).Error
	if err != nil {
		return nil, err
	}

	return row.ToAnnouncementAggr(), nil
}

// `withAnnouncementIncl` maps typed include options to preloads.
func withAnnouncementIncl(query *gorm.DB, inc ...enum.AnnouncementIncl) *gorm.DB {
	for _, i := range inc {
		switch i {
		case enum.AnnouncementInclUser:
			query = query.Preload("User")
		}
	}

	return query
}
