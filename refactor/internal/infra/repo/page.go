package repo_infra

import (
	"time"

	"poprako-s/internal/domain/model/aggr"
	"poprako-s/internal/domain/model/query"
	repo_iface "poprako-s/internal/domain/repo"
	"poprako-s/internal/infra/repo/entity"
)

import "gorm.io/gorm"

// `pageRepoImpl` is the GORM-backed implementation of `PageRepo`.
type pageRepoImpl struct {
	// `gdb` is the underlying GORM handle.
	gdb *gorm.DB
}

// `NewPageRepo` creates a non-transaction-scoped `PageRepo`.
func NewPageRepo(gdb *gorm.DB) repo_iface.PageRepo {
	return &pageRepoImpl{gdb: gdb}
}

// `GetById` retrieves one page by primary key.
func (r *pageRepoImpl) GetById(id string) (*aggr.Page, repo_iface.RepoErr) {
	var row entity.PageRow

	err := r.gdb.Table(entity.PAGE_TABLE).Where("id = ?", id).First(&row).Error
	if err != nil {
		return nil, err
	}

	return row.ToPageAggr(), nil
}

// `List` returns pages matching query options.
func (r *pageRepoImpl) List(opt *query.ListPageOpt) ([]*aggr.Page, repo_iface.RepoErr) {
	var rows []entity.PageRow

	q := r.gdb.Table(entity.PAGE_TABLE)
	if opt != nil && opt.ChapterId != nil {
		q = q.Where("chapter_id = ?", *opt.ChapterId)
	}
	if opt != nil && opt.Pagi.Offset > 0 {
		q = q.Offset(opt.Pagi.Offset)
	}
	if opt != nil && opt.Pagi.Limit > 0 {
		q = q.Limit(opt.Pagi.Limit)
	}

	err := q.Order("index ASC").Find(&rows).Error
	if err != nil {
		return nil, err
	}

	result := make([]*aggr.Page, len(rows))
	for i := range rows {
		result[i] = rows[i].ToPageAggr()
	}

	return result, nil
}

// `Count` returns page count matching query options.
func (r *pageRepoImpl) Count(opt *query.ListPageOpt) (int64, repo_iface.RepoErr) {
	var count int64

	q := r.gdb.Table(entity.PAGE_TABLE)
	if opt != nil && opt.ChapterId != nil {
		q = q.Where("chapter_id = ?", *opt.ChapterId)
	}

	err := q.Count(&count).Error
	if err != nil {
		return 0, err
	}

	return count, nil
}

// `CreateBatch` inserts many pages in one batch.
func (r *pageRepoImpl) CreateBatch(cre []*aggr.PageCre) repo_iface.RepoErr {
	if len(cre) == 0 {
		return nil
	}

	rows := make([]map[string]any, 0, len(cre))
	for i := range cre {
		row := entity.NewPageCreRowFromAggr(cre[i])
		rows = append(rows, map[string]any{
			"id":                    row.Id,
			"chapter_id":            row.ChapterId,
			"index":                 row.Index,
			"image_key":             row.ImageKey,
			"image_uploaded":        row.ImageUploaded,
			"total_unit_count":      row.TotalUnitCount,
			"translated_unit_count": row.TranslatedUnitCount,
			"proofread_unit_count":  row.ProofreadUnitCount,
		})
	}

	return r.gdb.Table(entity.PAGE_TABLE).Create(rows).Error
}

// `MarkImageUploaded` marks one page image as uploaded.
func (r *pageRepoImpl) MarkImageUploaded(id string) repo_iface.RepoErr {
	upRow := &entity.PageImageUploadedUpdRow{
		ImageUploaded: true,
		UpdatedAt:     time.Now(),
	}

	return r.gdb.Table(entity.PAGE_TABLE).Where("id = ?", id).Select("image_uploaded", "updated_at").Updates(upRow).Error
}

// `DeleteByChapterId` hard-deletes all pages under one chapter.
func (r *pageRepoImpl) DeleteByChapterId(chapterId string) repo_iface.RepoErr {
	return r.gdb.Table(entity.PAGE_TABLE).Where("chapter_id = ?", chapterId).Delete(nil).Error
}
