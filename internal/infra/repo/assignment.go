package repo_infra

import (
	"poprako-s/internal/domain/model/aggr"
	"poprako-s/internal/domain/model/enum"
	"poprako-s/internal/domain/model/query"
	repo_iface "poprako-s/internal/domain/repo"
	"poprako-s/internal/infra/repo/entity"

	"gorm.io/gorm"
)

// `assignmentRepoImpl` is gorm implementation of `AssignmentRepo`.
type assignmentRepoImpl struct {
	gdb *gorm.DB
}

// `NewAssignmentRepo` creates a non transaction-scoped `AssignmentRepo`.
func NewAssignmentRepo(gdb *gorm.DB) repo_iface.AssignmentRepo {
	return &assignmentRepoImpl{gdb: gdb}
}

// `GetById` returns one assignment by id.
func (r *assignmentRepoImpl) GetById(id string) (*aggr.Assignment, repo_iface.RepoErr) {
	var row entity.AssignmentRow

	err := r.gdb.
		Table(entity.ASSIGNMENT_TABLE).
		Where("id = ?", id).
		First(&row).Error
	if err != nil {
		return nil, err
	}

	return row.ToAssignmentAggr(), nil
}

// `GetByChapterUserId` returns one assignment by chapter and user.
func (r *assignmentRepoImpl) GetByChapterUserId(chapterId string, userId string) (*aggr.Assignment, repo_iface.RepoErr) {
	var row entity.AssignmentRow

	err := r.gdb.
		Table(entity.ASSIGNMENT_TABLE).
		Where("chapter_id = ? AND user_id = ?", chapterId, userId).
		First(&row).Error
	if err != nil {
		return nil, err
	}

	return row.ToAssignmentAggr(), nil
}

// `List` returns assignment list by query options.
func (r *assignmentRepoImpl) List(opt *query.ListAssignmentOpt, inc ...enum.AssignmentIncl) ([]*aggr.Assignment, repo_iface.RepoErr) {
	var rows []entity.AssignmentRow

	qry := r.gdb.Table(entity.ASSIGNMENT_TABLE)

	qry = withAssignmentIncl(qry, inc...)

	if opt != nil && opt.ChapterId != nil {
		qry = qry.Where("chapter_id = ?", *opt.ChapterId)
	}

	if opt != nil && opt.UserId != nil {
		qry = qry.Where("user_id = ?", *opt.UserId)
	}

	if opt != nil && opt.Pagi.Offset > 0 {
		qry = qry.Offset(opt.Pagi.Offset)
	}

	if opt != nil && opt.Pagi.Limit > 0 {
		qry = qry.Limit(opt.Pagi.Limit)
	}

	err := qry.Order("created_at DESC").Find(&rows).Error
	if err != nil {
		return nil, err
	}

	items := make([]*aggr.Assignment, len(rows))
	for i := range rows {
		items[i] = rows[i].ToAssignmentAggr()
	}

	return items, nil
}

// `withAssignmentIncl` maps typed include options to preloads.
func withAssignmentIncl(q *gorm.DB, inc ...enum.AssignmentIncl) *gorm.DB {
	for _, i := range inc {
		switch i {
		case enum.AssignmentInclUser:
			q = q.Preload("User")

		case enum.AssignmentInclChapter:
			q = q.Preload("Chapter")

		case enum.AssignmentInclChapterComic:
			q = q.Preload("Chapter.Comic")

		case enum.AssignmentInclChapterComicWorkset:
			q = q.Preload("Chapter.Comic.Workset")

		case enum.AssignmentInclChapterComicWorksetTeam:
			q = q.Preload("Chapter.Comic.Workset.Team")

		case enum.AssignmentInclChapterCreator:
			q = q.Preload("Chapter.Creator")

		case enum.AssignmentInclChapterComicCreator:
			q = q.Preload("Chapter.Comic.Creator")
		}
	}

	return q
}

// `Create` inserts one assignment.
func (r *assignmentRepoImpl) Create(cre *aggr.AssignmentCre) (*aggr.Assignment, repo_iface.RepoErr) {
	row := entity.NewAssignmentCreRowFromAggr(cre)

	err := r.gdb.Table(entity.ASSIGNMENT_TABLE).Create(row).Error
	if err != nil {
		return nil, err
	}

	var created entity.AssignmentRow
	err = r.gdb.Table(entity.ASSIGNMENT_TABLE).Where("id = ?", row.Id).First(&created).Error
	if err != nil {
		return nil, err
	}

	return created.ToAssignmentAggr(), nil
}

// `Put` overwrites all role timestamp fields.
func (r *assignmentRepoImpl) Put(put *aggr.AssignmentPut) repo_iface.RepoErr {
	row := entity.NewAssignmentPutRowFromAggr(put)

	return r.gdb.
		Table(entity.ASSIGNMENT_TABLE).
		Where("id = ?", put.Id).
		Select(
			"assigned_raw_provider_at",
			"assigned_translator_at",
			"assigned_proofreader_at",
			"assigned_typesetter_at",
			"assigned_redrawer_at",
			"assigned_reviewer_at",
			"assigned_publisher_at",
			"updated_at",
		).
		Updates(row).Error
}

// `Delete` hard deletes one assignment by id.
func (r *assignmentRepoImpl) Delete(id string) repo_iface.RepoErr {
	return r.gdb.
		Table(entity.ASSIGNMENT_TABLE).
		Where("id = ?", id).
		Delete(&entity.AssignmentRow{}).Error
}

// `DeleteByChapterUserId` hard deletes one assignment by chapter and user.
func (r *assignmentRepoImpl) DeleteByChapterUserId(chapterId string, userId string) repo_iface.RepoErr {
	return r.gdb.
		Table(entity.ASSIGNMENT_TABLE).
		Where("chapter_id = ? AND user_id = ?", chapterId, userId).
		Delete(&entity.AssignmentRow{}).Error
}

// `DeleteByChapterId` hard deletes all assignments by chapter.
func (r *assignmentRepoImpl) DeleteByChapterId(chapterId string) repo_iface.RepoErr {
	return r.gdb.
		Table(entity.ASSIGNMENT_TABLE).
		Where("chapter_id = ?", chapterId).
		Delete(&entity.AssignmentRow{}).Error
}
