package repo_infra

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"poprako-s/internal/domain/model/aggr"
	"poprako-s/internal/domain/model/enum"
	"poprako-s/internal/domain/model/query"
	repo_iface "poprako-s/internal/domain/repo"
	"poprako-s/internal/infra/repo/entity"

	"gorm.io/gorm"
)

const chapterPinAlias = "pch"

// `comicRepoImpl` is the GORM-backed implementation of `ComicRepo`
type comicRepoImpl struct {
	// `gdb` is the underlying GORM handle
	gdb *gorm.DB
}

// `TxnComicRepo` creates a transaction-scoped `ComicRepo` extracted from `cx`
func TxnComicRepo(cx context.Context) (repo_iface.ComicRepo, repo_iface.RepoErr) {
	gdb := takeTxnGdb(cx)
	if gdb == nil {
		return nil, errors.New("[TxnComicRepo] no transaction context found for ComicRepo")
	}

	return &comicRepoImpl{gdb: gdb}, nil
}

// `NewComicRepo` creates a non-transaction-scoped `ComicRepo`
func NewComicRepo(gdb *gorm.DB) repo_iface.ComicRepo {
	return &comicRepoImpl{gdb: gdb}
}

// `GetById` retrieves one active comic by primary key
func (r *comicRepoImpl) GetById(id string, inc ...enum.ComicIncl) (*aggr.Comic, repo_iface.RepoErr) {
	var row entity.ComicRow

	q := r.gdb.
		Table(entity.COMIC_TABLE).
		Where("t_comic.id = ? AND t_comic.deleted_at IS NULL", id)

	q = withComicIncl(q, inc...)

	err := q.First(&row).Error
	if err != nil {
		return nil, err
	}

	return row.ToComicAggr(), nil
}

// `List` returns active comics matching query options and includes
func (r *comicRepoImpl) List(opt *query.ListComicOpt, inc ...enum.ComicIncl) ([]*aggr.Comic, repo_iface.RepoErr) {
	var rows []entity.ComicRow

	q := r.gdb.
		Table(entity.COMIC_TABLE).
		Where("t_comic.deleted_at IS NULL")

	if hasComicSearchFilter(opt) {
		q = q.Where("t_comic.is_completed = FALSE")
	}

	if opt != nil && opt.WorksetId != nil {
		q = q.Where("workset_id = ?", *opt.WorksetId)
	}

	if opt != nil && opt.FuzzyTitle != nil {
		fuzzyTitle := strings.TrimSpace(*opt.FuzzyTitle)

		if fuzzyTitle != "" {
			q = q.Where("t_comic.fuzzy_title ILIKE ?", "%"+fuzzyTitle+"%")
		}
	}

	if hasComicWorkflowFilter(opt) {
		q = withPinnedChapterJoin(q)

		q = applyComicWorkflowFilter(q, opt.UploadPhase, "uploaded_at", "")
		q = applyComicWorkflowFilter(q, opt.TranslatePhase, "transalating_at", "translated_at")
		q = applyComicWorkflowFilter(q, opt.ProofreadPhase, "proofreading_at", "proofread_at")
		q = applyComicWorkflowFilter(q, opt.TypesetPhase, "typesetting_at", "typeset_at")
		q = applyComicWorkflowFilter(q, opt.ReviewPhase, "reviewed_at", "")
		q = applyComicWorkflowFilter(q, opt.PublishPhase, "published_at", "")
	}

	if opt != nil && opt.Pagi.Offset > 0 {
		q = q.Offset(opt.Pagi.Offset)
	}

	if opt != nil && opt.Pagi.Limit > 0 {
		q = q.Limit(opt.Pagi.Limit)
	}

	q = withComicIncl(q, inc...)

	err := q.Order("last_active_at DESC").Find(&rows).Error
	if err != nil {
		return nil, err
	}

	result := make([]*aggr.Comic, len(rows))

	for i := range rows {
		result[i] = rows[i].ToComicAggr()
	}

	return result, nil
}

// `Count` returns active comic count matching query options
func (r *comicRepoImpl) Count(opt *query.ListComicOpt) (int64, repo_iface.RepoErr) {
	var count int64

	q := r.gdb.
		Table(entity.COMIC_TABLE).
		Where("t_comic.deleted_at IS NULL")

	if hasComicSearchFilter(opt) {
		q = q.Where("t_comic.is_completed = FALSE")
	}

	if opt != nil && opt.WorksetId != nil {
		q = q.Where("workset_id = ?", *opt.WorksetId)
	}

	if opt != nil && opt.FuzzyTitle != nil {
		fuzzyTitle := strings.TrimSpace(*opt.FuzzyTitle)

		if fuzzyTitle != "" {
			q = q.Where("t_comic.fuzzy_title ILIKE ?", "%"+fuzzyTitle+"%")
		}
	}

	if hasComicWorkflowFilter(opt) {
		q = withPinnedChapterJoin(q)

		q = applyComicWorkflowFilter(q, opt.UploadPhase, "uploaded_at", "")
		q = applyComicWorkflowFilter(q, opt.TranslatePhase, "transalating_at", "translated_at")
		q = applyComicWorkflowFilter(q, opt.ProofreadPhase, "proofreading_at", "proofread_at")
		q = applyComicWorkflowFilter(q, opt.TypesetPhase, "typesetting_at", "typeset_at")
		q = applyComicWorkflowFilter(q, opt.ReviewPhase, "reviewed_at", "")
		q = applyComicWorkflowFilter(q, opt.PublishPhase, "published_at", "")
	}

	err := q.Count(&count).Error
	if err != nil {
		return 0, err
	}

	return count, nil
}

// `Create` inserts one comic row then reloads the created aggregate
func (r *comicRepoImpl) Create(cre *aggr.ComicCre) (*aggr.Comic, repo_iface.RepoErr) {
	fuzzyTitle := buildComicFuzzyTitle(cre.Index, cre.Author, cre.Title)
	row := entity.NewComicCreRowFromAggr(cre, fuzzyTitle)

	err := r.gdb.
		Table(row.TableName()).
		Create(row).Error
	if err != nil {
		return nil, err
	}

	return r.GetById(row.Id)
}

// `Update` applies put-style mutable fields to matching comic row
func (r *comicRepoImpl) Update(upd *aggr.ComicUpd) repo_iface.RepoErr {
	now := time.Now()

	updRow := &entity.ComicUpdRow{
		Title:      upd.Title,
		Author:     upd.Author,
		FuzzyTitle: "",
		Desc:       upd.Desc,
		UpdatedAt:  now,
	}

	// Resolve current comic index so fuzzy text always carries index info.
	cm, err := r.GetById(upd.Id)
	if err != nil {
		return err
	}

	updRow.FuzzyTitle = buildComicFuzzyTitle(cm.Index, upd.Author, upd.Title)

	err = r.gdb.
		Table(entity.COMIC_TABLE).
		Where("t_comic.id = ? AND t_comic.deleted_at IS NULL", upd.Id).
		Select("title", "author", "fuzzy_title", "description", "updated_at").
		Updates(updRow).Error

	return err
}

// `UpdateChapterCount` applies delta to comic chapter counter.
func (r *comicRepoImpl) UpdateChapterCount(id string, delta int) repo_iface.RepoErr {
	now := time.Now()

	return r.gdb.
		Table(entity.COMIC_TABLE).
		Where("t_comic.id = ? AND t_comic.deleted_at IS NULL", id).
		Updates(map[string]any{
			"chapter_count": gorm.Expr("GREATEST(chapter_count + ?, 0)", delta),
			"updated_at":    now,
		}).Error
}

// `TouchLastActive` refreshes comic activity timestamp.
func (r *comicRepoImpl) TouchLastActive(id string) repo_iface.RepoErr {
	now := time.Now()

	updRow := &entity.ComicLastActiveUpdRow{
		LastActiveAt: now,
		UpdatedAt:    now,
	}

	err := r.gdb.
		Table(entity.COMIC_TABLE).
		Where("t_comic.id = ? AND t_comic.deleted_at IS NULL", id).
		Select("last_active_at", "updated_at").
		Updates(updRow).Error

	return err
}

// `Remove` soft-deletes one comic by setting `deleted_at`
func (r *comicRepoImpl) Remove(id string) repo_iface.RepoErr {
	now := time.Now()

	err := r.gdb.
		Table(entity.COMIC_TABLE).
		Where("t_comic.id = ? AND t_comic.deleted_at IS NULL", id).
		Updates(map[string]any{
			"deleted_at": now,
			"updated_at": now,
		}).Error

	return err
}

// `withComicIncl` maps typed include options to preload actions
func withComicIncl(q *gorm.DB, inc ...enum.ComicIncl) *gorm.DB {
	for _, i := range inc {
		switch i {
		case enum.ComicInclWorkset:
			q = q.Preload("Workset")
		}
	}

	return q
}

// `withPinnedChapterJoin` joins pinned chapter row for workflow filtering.
func withPinnedChapterJoin(q *gorm.DB) *gorm.DB {
	return q.Joins(
		"LEFT JOIN t_chapter AS " + chapterPinAlias +
			" ON " + chapterPinAlias + ".comic_id = t_comic.id" +
			" AND " + chapterPinAlias + ".pinned = TRUE" +
			" AND " + chapterPinAlias + ".deleted_at IS NULL",
	)
}

// `applyComicWorkflowFilter` applies one workflow phase filter against pinned chapter fields.
func applyComicWorkflowFilter(
	q *gorm.DB,
	phase *enum.WorkflowPhase,
	startedCol string,
	completedCol string,
) *gorm.DB {
	if phase == nil {
		return q
	}

	startedExpr := chapterPinAlias + "." + startedCol
	completedExpr := chapterPinAlias + "." + completedCol

	switch *phase {
	case enum.WorkflowPending:
		return q.Where(
			"(" + chapterPinAlias + ".id IS NULL OR " + startedExpr + " IS NULL)",
		)

	case enum.WorkflowOngoing:
		if completedCol == "" {
			return q.Where("1 = 0")
		}

		return q.Where(
			startedExpr + " IS NOT NULL AND " + completedExpr + " IS NULL",
		)

	case enum.WorkflowCompleted:
		if completedCol == "" {
			return q.Where(startedExpr + " IS NOT NULL")
		}

		return q.Where(completedExpr + " IS NOT NULL")

	default:
		return q
	}
}

// `hasComicWorkflowFilter` reports whether comic query contains workflow constraints.
func hasComicWorkflowFilter(opt *query.ListComicOpt) bool {
	if opt == nil {
		return false
	}

	return opt.UploadPhase != nil ||
		opt.TranslatePhase != nil ||
		opt.ProofreadPhase != nil ||
		opt.TypesetPhase != nil ||
		opt.ReviewPhase != nil ||
		opt.PublishPhase != nil
}

// `hasComicSearchFilter` reports whether comic query contains any business filter.
func hasComicSearchFilter(opt *query.ListComicOpt) bool {
	return hasComicFuzzyTitleFilter(opt) || hasComicWorkflowFilter(opt)
}

// `hasComicFuzzyTitleFilter` reports whether comic query contains fuzzy title filter.
func hasComicFuzzyTitleFilter(opt *query.ListComicOpt) bool {
	if opt == nil || opt.FuzzyTitle == nil {
		return false
	}

	return strings.TrimSpace(*opt.FuzzyTitle) != ""
}

// `buildComicFuzzyTitle` builds dedicated fuzzy text containing index and title fields.
func buildComicFuzzyTitle(index int, author string, title string) string {
	return fmt.Sprintf("【%d】[%s] %s", index+1, author, title)
}
