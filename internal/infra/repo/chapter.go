package repo_infra

import (
	"errors"
	"fmt"
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

// `GetById` retrieves one chapter by id.
func (r *chapterRepoImpl) GetById(id string, inc ...enum.ChapterIncl) (*aggr.Chapter, repo_iface.RepoErr) {
	var row entity.ChapterRow

	q := r.gdb.Table(entity.CHAPTER_TABLE).Where("id = ?", id)
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
		Where("comic_id = ? AND pinned = TRUE", comicId)
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

// `FindPinnedByComics` retrieves pinned chapters for many comics.
func (r *chapterRepoImpl) FindPinnedByComics(comicIds []string) ([]*aggr.Chapter, repo_iface.RepoErr) {
	if len(comicIds) == 0 {
		return nil, nil
	}

	var rows []entity.ChapterRow

	err := r.gdb.
		Table(entity.CHAPTER_TABLE).
		Where("comic_id IN ? AND pinned = TRUE", comicIds).
		Find(&rows).Error
	if err != nil {
		return nil, err
	}

	chapterByComicId := make(map[string]*aggr.Chapter, len(rows))
	for i := range rows {
		ch := rows[i].ToChapterAggr()

		chapterByComicId[ch.ComicId] = ch
	}

	re := make([]*aggr.Chapter, len(comicIds))
	for i := range comicIds {
		re[i] = chapterByComicId[comicIds[i]]
	}

	return re, nil
}

// `List` lists chapters matching options.
func (r *chapterRepoImpl) List(opt *query.ListChapterOpt, inc ...enum.ChapterIncl) ([]*aggr.Chapter, repo_iface.RepoErr) {
	var rows []entity.ChapterRow

	q := r.gdb.Table(entity.CHAPTER_TABLE)

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

// `Count` counts chapters matching options.
func (r *chapterRepoImpl) Count(opt *query.ListChapterOpt) (int64, repo_iface.RepoErr) {
	var count int64

	q := r.gdb.Table(entity.CHAPTER_TABLE)

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
		subtitle = fmt.Sprintf("第%d话", cre.Index+1)
	}

	row := entity.NewChapterCreRowFromAggr(cre, subtitle)

	err := r.gdb.Transaction(func(tx *gorm.DB) error {
		pinClrRow := &entity.ChapterPinUpdRow{IsPinned: false, UpdatedAt: now}

		if err := tx.
			Table(entity.CHAPTER_TABLE).
			Where("comic_id = ? AND pinned = TRUE", cre.ComicId).
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
				Where("comic_id = (SELECT comic_id FROM t_chapter WHERE id = ?) AND id <> ? AND pinned = TRUE", upd.Id, upd.Id).
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

	if upd.TranslatingAt != nil {
		updRow.TranslatingAt = *upd.TranslatingAt
		selectCols = append(selectCols, "translating_at")
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

	q := r.gdb.
		Table(entity.CHAPTER_TABLE).
		Where("id = ?", upd.Id)

	if upd.WorkflowTransition != nil {
		switch *upd.WorkflowTransition {
		case enum.WorkflowUploadComplete:
			q = q.Where("uploaded_at IS NULL")
		case enum.WorkflowTranslateStart:
			q = q.Where("translating_at IS NULL").Where("translated_at IS NULL")
		case enum.WorkflowTranslateComplete:
			q = q.Where("translating_at IS NOT NULL").Where("translated_at IS NULL")
		case enum.WorkflowProofreadStart:
			q = q.Where("proofreading_at IS NULL").Where("proofread_at IS NULL")
		case enum.WorkflowProofreadComplete:
			q = q.Where("proofreading_at IS NOT NULL").Where("proofread_at IS NULL")
		case enum.WorkflowTypesetStart:
			q = q.Where("typesetting_at IS NULL").Where("typeset_at IS NULL")
		case enum.WorkflowTypesetComplete:
			q = q.Where("typesetting_at IS NOT NULL").Where("typeset_at IS NULL")
		case enum.WorkflowReviewComplete:
			q = q.Where("reviewed_at IS NULL")
		case enum.WorkflowPublishComplete:
			q = q.Where("published_at IS NULL")
		}
	}

	tx := q.Select(selectCols).Updates(updRow)
	if tx.Error != nil {
		return tx.Error
	}

	// Revert transitions are no-op when the field is already NULL; skip the rows-affected check
	if upd.WorkflowTransition != nil && tx.RowsAffected == 0 {
		return errConditionalUpdateFailed
	}

	return nil
}

// `SetPageCount` overwrites the page count of one chapter.
func (r *chapterRepoImpl) SetPageCount(id string, count int) repo_iface.RepoErr {
	updRow := &entity.ChapterPageCountUpdRow{
		PageCount: count,
		UpdatedAt: time.Now(),
	}

	return r.gdb.
		Table(entity.CHAPTER_TABLE).
		Where("id = ?", id).
		Select("page_count", "updated_at").
		Updates(updRow).Error
}

// `AdjustUnitCounts` atomically applies unit count delta to one chapter.
func (r *chapterRepoImpl) AdjustUnitCounts(id string, deltaTotal int, deltaTranslated int, deltaProofread int) repo_iface.RepoErr {
	return r.gdb.
		Table(entity.CHAPTER_TABLE).
		Where("id = ?", id).
		Updates(map[string]any{
			"total_unit_count":      gorm.Expr("total_unit_count + ?", deltaTotal),
			"translated_unit_count": gorm.Expr("translated_unit_count + ?", deltaTranslated),
			"proofread_unit_count":  gorm.Expr("proofread_unit_count + ?", deltaProofread),
			"updated_at":            time.Now(),
		}).Error
}

// `Delete` hard-deletes one chapter.
func (r *chapterRepoImpl) Delete(id string) repo_iface.RepoErr {
	return r.gdb.
		Table(entity.CHAPTER_TABLE).
		Where("id = ?", id).
		Delete(&entity.ChapterRow{}).Error
}

// `withChapterIncl` maps typed include options to preloads.
func withChapterIncl(q *gorm.DB, inc ...enum.ChapterIncl) *gorm.DB {
	for _, i := range inc {
		switch i {
		case enum.ChapterInclComic:
			q = q.Preload("Comic")

		case enum.ChapterInclComicWorkset:
			q = q.Preload("Comic.Workset")

		case enum.ChapterInclComicWorksetTeam:
			q = q.Preload("Comic.Workset.Team")

		case enum.ChapterInclComicCreator:
			q = q.Preload("Comic.Creator")

		case enum.ChapterInclCreator:
			q = q.Preload("Creator")
		}
	}

	return q
}
