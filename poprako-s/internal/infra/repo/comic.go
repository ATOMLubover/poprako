package repo_infra

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"poprako-s/internal/domain/model"
	iface "poprako-s/internal/domain/repo"
	entity "poprako-s/internal/infra/repo/entity"

	"gorm.io/gorm"
)

type comicRepoImpl struct {
	gdb *gorm.DB
}

func NewComicRepo(gdb *gorm.DB) iface.ComicRepo {
	return &comicRepoImpl{gdb: gdb}
}

func NewComicRepoFromCx(cx context.Context) (iface.ComicRepo, error) {
	gdb, ok := cx.Value(txnKey).(*gorm.DB)
	if !ok {
		return nil, errors.New("[NewComicRepoFromCx]: 无法从上下文中获取事务数据库连接")
	}

	return &comicRepoImpl{gdb: gdb}, nil
}

func (r *comicRepoImpl) FromTxnCx(cx context.Context) (iface.ComicRepo, error) {
	return NewComicRepoFromCx(cx)
}

func (r *comicRepoImpl) GetByID(id string) (*model.ComicInfo, error) {
	var row entity.ComicInfoRow
	err := r.gdb.Table(entity.ComicTable).
		Where("id = ? AND deleted_at IS NULL", id).
		First(&row).Error
	if err != nil {
		return nil, err
	}

	info := entity.ToComicInfo(row)
	return &info, nil
}

func (r *comicRepoImpl) List(opt model.ComicQueryOpt) ([]model.ComicInfo, error) {
	db := r.gdb.Table(entity.ComicTable).
		Where("workset_id = ? AND deleted_at IS NULL", opt.WorksetID)

	if opt.ID != nil {
		db = db.Where("id = ?", *opt.ID)
	}
	if opt.FuzzyTitle != nil {
		term := strings.TrimSpace(*opt.FuzzyTitle)
		if term != "" {
			db = db.Where("composed_title ILIKE ?", "%"+term+"%")
		}
	}

	db = applyComicWorkflowFilter(db, opt.UploadStatus, "pinned_uploaded_at", "")
	db = applyComicWorkflowFilter(db, opt.TranslateStatus, "pinned_transalating_at", "pinned_translated_at")
	db = applyComicWorkflowFilter(db, opt.ProofreadStatus, "pinned_proofreading_at", "pinned_proofread_at")
	db = applyComicWorkflowFilter(db, opt.TypesetStatus, "pinned_typesetting_at", "pinned_typeset_at")
	db = applyComicWorkflowFilter(db, opt.ReviewStatus, "pinned_reviewed_at", "")
	db = applyComicWorkflowFilter(db, opt.PublishStatus, "pinned_published_at", "")

	if opt.Limit > 0 {
		db = db.Offset(opt.Offset).Limit(opt.Limit)
	}

	var rows []entity.ComicInfoRow
	if err := db.Order("index ASC").Find(&rows).Error; err != nil {
		return nil, err
	}

	items := make([]model.ComicInfo, 0, len(rows))
	for _, row := range rows {
		items = append(items, entity.ToComicInfo(row))
	}

	return items, nil
}

func (r *comicRepoImpl) Count(opt model.ComicQueryOpt) (int64, error) {
	db := r.gdb.Table(entity.ComicTable).
		Where("workset_id = ? AND deleted_at IS NULL", opt.WorksetID)

	if opt.ID != nil {
		db = db.Where("id = ?", *opt.ID)
	}
	if opt.FuzzyTitle != nil {
		term := strings.TrimSpace(*opt.FuzzyTitle)
		if term != "" {
			db = db.Where("composed_title ILIKE ?", "%"+term+"%")
		}
	}

	db = applyComicWorkflowFilter(db, opt.UploadStatus, "pinned_uploaded_at", "")
	db = applyComicWorkflowFilter(db, opt.TranslateStatus, "pinned_transalating_at", "pinned_translated_at")
	db = applyComicWorkflowFilter(db, opt.ProofreadStatus, "pinned_proofreading_at", "pinned_proofread_at")
	db = applyComicWorkflowFilter(db, opt.TypesetStatus, "pinned_typesetting_at", "pinned_typeset_at")
	db = applyComicWorkflowFilter(db, opt.ReviewStatus, "pinned_reviewed_at", "")
	db = applyComicWorkflowFilter(db, opt.PublishStatus, "pinned_published_at", "")

	var n int64
	err := db.Count(&n).Error
	return n, err
}

func (r *comicRepoImpl) Create(c *model.ComicCreation) (*model.ComicInfo, error) {
	now := time.Now()
	row := map[string]any{
		"id":                 c.ID,
		"workset_id":         c.WorksetID,
		"index":              c.Index,
		"title":              c.Title,
		"author":             c.Author,
		"composed_title":     fmt.Sprintf("【%d】[%s] %s", c.Index, c.Author, c.Title),
		"description":        c.Description,
		"chapter_count":      0,
		"has_pinned_chapter": false,
		"creator_id":         c.CreatorID,
		"last_active_at":     now,
		"created_at":         now,
		"updated_at":         now,
	}

	if err := r.gdb.Table(entity.ComicTable).Create(row).Error; err != nil {
		return nil, err
	}

	return r.GetByID(c.ID)
}

func (r *comicRepoImpl) Update(u *model.ComicUpdate) error {
	var row entity.ComicInfoRow
	err := r.gdb.Table(entity.ComicTable).
		Select("id", "index").
		Where("id = ? AND deleted_at IS NULL", u.ID).
		First(&row).Error
	if err != nil {
		return err
	}

	return r.gdb.Table(entity.ComicTable).
		Where("id = ? AND deleted_at IS NULL", u.ID).
		Updates(map[string]any{
			"title":          u.Title,
			"author":         u.Author,
			"description":    u.Description,
			"composed_title": fmt.Sprintf("【%d】[%s] %s", row.Index, u.Author, u.Title),
			"updated_at":     time.Now(),
		}).Error
}

func (r *comicRepoImpl) UpdateChapterCount(id string, delta int) error {
	return r.gdb.Table(entity.ComicTable).
		Where("id = ? AND deleted_at IS NULL", id).
		Updates(map[string]any{
			"chapter_count":  gorm.Expr("chapter_count + ?", delta),
			"last_active_at": time.Now(),
			"updated_at":     time.Now(),
		}).Error
}

func (r *comicRepoImpl) Delete(id string) error {
	return r.gdb.Table(entity.ComicTable).
		Where("id = ? AND deleted_at IS NULL", id).
		Updates(map[string]any{
			"deleted_at": time.Now(),
			"updated_at": time.Now(),
		}).Error
}

func applyComicWorkflowFilter(db *gorm.DB, phase *model.WorkflowPhase, startedColumn string, completedColumn string) *gorm.DB {
	if phase == nil {
		return db
	}

	switch *phase {
	case model.WorkflowPending:
		if completedColumn == "" {
			return db.Where("has_pinned_chapter = FALSE OR " + startedColumn + " IS NULL")
		}
		return db.Where("has_pinned_chapter = FALSE OR " + startedColumn + " IS NULL")
	case model.WorkflowOngoing:
		if completedColumn == "" {
			return db.Where("1 = 0")
		}
		return db.Where(startedColumn + " IS NOT NULL").Where(completedColumn + " IS NULL")
	case model.WorkflowCompleted:
		if completedColumn == "" {
			return db.Where(startedColumn + " IS NOT NULL")
		}
		return db.Where(completedColumn + " IS NOT NULL")
	default:
		return db
	}
}
