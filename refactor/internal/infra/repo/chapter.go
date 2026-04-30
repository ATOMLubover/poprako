package repo_infra

import (
	"errors"
	"strconv"
	"time"

	"poprako-s/internal/domain/model/aggr"
	"poprako-s/internal/domain/model/enum"
	"poprako-s/internal/domain/model/query"
	repo_iface "poprako-s/internal/domain/repo"
	"poprako-s/internal/infra/repo/entity"

	"gorm.io/gorm"
)

// `chapterRepoImpl` is GORM-backed implementation of `ChapterRepo`.
type chapterRepoImpl struct {
	// `gdb` is underlying GORM database handle.
	gdb *gorm.DB
}

// `NewChapterRepo` creates non-transaction chapter repo.
func NewChapterRepo(gdb *gorm.DB) repo_iface.ChapterRepo {
	return &chapterRepoImpl{gdb: gdb}
}

// `GetById` retrieves one active chapter by id.
func (r *chapterRepoImpl) GetById(id string, inc ...enum.ChapterIncl) (*aggr.Chapter, repo_iface.RepoErr) {
	var row entity.ChapterRow

	q := r.gdb.Table(entity.CHAPTER_TABLE).Where("id = ? AND deleted_at IS NULL", id)
	q = withChapterIncl(q, inc...)

	err := q.First(&row).Error
	if err != nil {
		return nil, err
	}

	return row.ToChapterAggr(), nil
}

// `FindPinnedByComicId` retrieves pinned chapter under comic.
func (r *chapterRepoImpl) FindPinnedByComicId(comicId string, inc ...enum.ChapterIncl) (*aggr.Chapter, repo_iface.RepoErr) {
	var row entity.ChapterRow

	q := r.gdb.
		Table(entity.CHAPTER_TABLE).
		Where("comic_id = ? AND pinned = TRUE AND deleted_at IS NULL", comicId)
	q = withChapterIncl(q, inc...)

	err := q.First(&row).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return row.ToChapterAggr(), nil
}

// `List` lists active chapters matching options.
func (r *chapterRepoImpl) List(opt *query.ListChapterOpt, inc ...enum.ChapterIncl) ([]*aggr.Chapter, repo_iface.RepoErr) {
	var rows []entity.ChapterRow

	q := r.gdb.Table(entity.CHAPTER_TABLE).Where("deleted_at IS NULL")

	if opt != nil && opt.ComicId != nil {
		q = q.Where("comic_id = ?", *opt.ComicId)
	}

	if opt != nil && opt.Pagi.Offset > 0 {
		q = q.Offset(opt.Pagi.Offset)
	}

	if opt != nil && opt.Pagi.Limit > 0 {
		q = q.Limit(opt.Pagi.Limit)
	}

	q = withChapterIncl(q, inc...)

	err := q.Order("index DESC").Find(&rows).Error
	if err != nil {
		return nil, err
	}

	result := make([]*aggr.Chapter, len(rows))
	for i := range rows {
		result[i] = rows[i].ToChapterAggr()
	}

	return result, nil
}

// `Count` counts active chapters matching options.
func (r *chapterRepoImpl) Count(opt *query.ListChapterOpt) (int64, repo_iface.RepoErr) {
	var count int64

	q := r.gdb.Table(entity.CHAPTER_TABLE).Where("deleted_at IS NULL")

	if opt != nil && opt.ComicId != nil {
		q = q.Where("comic_id = ?", *opt.ComicId)
	}

	err := q.Count(&count).Error
	if err != nil {
		return 0, err
	}

	return count, nil
}

// `Create` inserts one chapter and keeps pinned uniqueness in comic.
func (r *chapterRepoImpl) Create(cre *aggr.ChapterCre) (*aggr.Chapter, repo_iface.RepoErr) {
	now := time.Now()

	subtitle := ""
	if cre.Subtitle != nil && *cre.Subtitle != "" {
		subtitle = *cre.Subtitle
	} else {
		subtitle = "Ch." + fmtInt(cre.Index)
	}

	row := entity.NewChapterCreRowFromAggr(cre, subtitle)

	err := r.gdb.Transaction(func(tx *gorm.DB) error {
		pinClrRow := &entity.ChapterPinUpdRow{IsPinned: false, UpdatedAt: now}

		if err := tx.
			Table(entity.CHAPTER_TABLE).
			Where("comic_id = ? AND deleted_at IS NULL AND pinned = TRUE", cre.ComicId).
			Select("pinned", "updated_at").
			Updates(pinClrRow).Error; err != nil {
			return err
		}

		if err := tx.Table(entity.CHAPTER_TABLE).Create(row).Error; err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return r.GetById(row.Id)
}

// `Update` applies mutable fields to chapter row.
func (r *chapterRepoImpl) Update(upd *aggr.ChapterUpd) repo_iface.RepoErr {
	now := time.Now()
	updRow := &entity.ChapterUpdRow{UpdatedAt: now}

	selectCols := make([]string, 0, 16)

	if upd.Subtitle != nil {
		updRow.Subtitle = upd.Subtitle
		selectCols = append(selectCols, "subtitle")
	}

	if upd.IsPinned != nil {
		if *upd.IsPinned {
			if err := r.gdb.Table(entity.CHAPTER_TABLE).
				Where("comic_id = (SELECT comic_id FROM t_chapter WHERE id = ? AND deleted_at IS NULL) AND id <> ? AND deleted_at IS NULL AND pinned = TRUE", upd.Id, upd.Id).
				Select("pinned", "updated_at").
				Updates(&entity.ChapterPinUpdRow{IsPinned: false, UpdatedAt: now}).Error; err != nil {
				return err
			}
		}

		updRow.IsPinned = upd.IsPinned
		selectCols = append(selectCols, "pinned")
	}

	if upd.UploadedAt != nil {
		updRow.UploadedAt = *upd.UploadedAt
		selectCols = append(selectCols, "uploaded_at")
	}

	if upd.TransalatingAt != nil {
		updRow.TransalatingAt = *upd.TransalatingAt
		selectCols = append(selectCols, "transalating_at")
	}

	if upd.TranslatedAt != nil {
		updRow.TranslatedAt = *upd.TranslatedAt
		selectCols = append(selectCols, "translated_at")
	}

	if upd.ProofreadingAt != nil {
		updRow.ProofreadingAt = *upd.ProofreadingAt
		selectCols = append(selectCols, "proofreading_at")
	}

	if upd.ProofreadAt != nil {
		updRow.ProofreadAt = *upd.ProofreadAt
		selectCols = append(selectCols, "proofread_at")
	}

	if upd.TypesettingAt != nil {
		updRow.TypesettingAt = *upd.TypesettingAt
		selectCols = append(selectCols, "typesetting_at")
	}

	if upd.TypesetAt != nil {
		updRow.TypesetAt = *upd.TypesetAt
		selectCols = append(selectCols, "typeset_at")
	}

	if upd.ReviewedAt != nil {
		updRow.ReviewedAt = *upd.ReviewedAt
		selectCols = append(selectCols, "reviewed_at")
	}

	if upd.PublishedAt != nil {
		updRow.PublishedAt = *upd.PublishedAt
		selectCols = append(selectCols, "published_at")
	}

	if len(selectCols) == 0 {
		return nil
	}

	selectCols = append(selectCols, "updated_at")

	return r.gdb.
		Table(entity.CHAPTER_TABLE).
		Where("id = ? AND deleted_at IS NULL", upd.Id).
		Select(selectCols).
		Updates(updRow).Error
}

// `SetPageCount` overwrites the page count of one chapter.
func (r *chapterRepoImpl) SetPageCount(id string, count int) repo_iface.RepoErr {
	updRow := &entity.ChapterPageCountUpdRow{
		PageCount: count,
		UpdatedAt: time.Now(),
	}

	return r.gdb.
		Table(entity.CHAPTER_TABLE).
		Where("id = ? AND deleted_at IS NULL", id).
		Select("page_count", "updated_at").
		Updates(updRow).Error
}

// `Remove` soft-deletes one chapter.
func (r *chapterRepoImpl) Remove(id string) repo_iface.RepoErr {
	now := time.Now()
	updRow := &entity.ChapterRemoveUpdRow{DeletedAt: now, UpdatedAt: now}

	return r.gdb.
		Table(entity.CHAPTER_TABLE).
		Where("id = ? AND deleted_at IS NULL", id).
		Select("deleted_at", "updated_at").
		Updates(updRow).Error
}

// `withChapterIncl` maps typed include options to preloads.
func withChapterIncl(q *gorm.DB, inc ...enum.ChapterIncl) *gorm.DB {
	for _, i := range inc {
		switch i {
		case enum.ChapterInclComic:
			q = q.Preload("Comic")
		}
	}

	return q
}

// `fmtInt` formats index without importing strconv in multiple places.
func fmtInt(v int) string {
	return strconv.FormatInt(int64(v), 10)
}
