package repo_infra

import (
	"context"
	"errors"
	"fmt"
	"time"

	"poprako-s/internal/domain/model"
	iface "poprako-s/internal/domain/repo"
	entity "poprako-s/internal/infra/repo/entity"

	"gorm.io/gorm"
)

type chapterRepoImpl struct {
	gdb *gorm.DB
}

func NewChapterRepo(gdb *gorm.DB) iface.ChapterRepo {
	return &chapterRepoImpl{gdb: gdb}
}

func NewChapterRepoFromCx(cx context.Context) (iface.ChapterRepo, error) {
	gdb, ok := cx.Value(txnKey).(*gorm.DB)
	if !ok {
		return nil, errors.New("[NewChapterRepoFromCx]: 无法从上下文中获取事务数据库连接")
	}

	return &chapterRepoImpl{gdb: gdb}, nil
}

func (r *chapterRepoImpl) FromTxnCx(cx context.Context) (iface.ChapterRepo, error) {
	return NewChapterRepoFromCx(cx)
}

func (r *chapterRepoImpl) GetByID(id string) (*model.ChapterInfo, error) {
	var row entity.ChapterInfoRow
	err := r.gdb.Table(entity.ChapterTable).
		Where("id = ? AND deleted_at IS NULL", id).
		First(&row).Error
	if err != nil {
		return nil, err
	}

	info := entity.ToChapterInfo(row)
	return &info, nil
}

func (r *chapterRepoImpl) FindPinnedByComicID(comicID string) (*model.ChapterInfo, error) {
	var row entity.ChapterInfoRow
	err := r.gdb.Table(entity.ChapterTable).
		Where("comic_id = ? AND pinned = TRUE AND deleted_at IS NULL", comicID).
		Order("index DESC").
		First(&row).Error
	if err != nil {
		return nil, err
	}

	info := entity.ToChapterInfo(row)
	return &info, nil
}

func (r *chapterRepoImpl) List(opt model.ChapterQueryOpt) ([]model.ChapterInfo, error) {
	db := r.gdb.Table(entity.ChapterTable).Where("deleted_at IS NULL")

	if opt.ID != nil {
		db = db.Where("id = ?", *opt.ID)
	}
	if opt.ComicID != nil {
		db = db.Where("comic_id = ?", *opt.ComicID)
	}
	if opt.CreatorID != nil {
		db = db.Where("creator_id = ?", *opt.CreatorID)
	}

	var rows []entity.ChapterInfoRow
	if err := db.Order("index DESC").Find(&rows).Error; err != nil {
		return nil, err
	}

	items := make([]model.ChapterInfo, 0, len(rows))
	for _, row := range rows {
		items = append(items, entity.ToChapterInfo(row))
	}

	return items, nil
}

func (r *chapterRepoImpl) Count(opt model.ChapterQueryOpt) (int64, error) {
	db := r.gdb.Table(entity.ChapterTable).Where("deleted_at IS NULL")

	if opt.ID != nil {
		db = db.Where("id = ?", *opt.ID)
	}
	if opt.ComicID != nil {
		db = db.Where("comic_id = ?", *opt.ComicID)
	}
	if opt.CreatorID != nil {
		db = db.Where("creator_id = ?", *opt.CreatorID)
	}

	var n int64
	err := db.Count(&n).Error
	return n, err
}

func (r *chapterRepoImpl) Create(c *model.ChapterCreation) (*model.ChapterInfo, error) {
	now := time.Now()
	subtitle := fmt.Sprintf("Ch.%d", c.Index)
	if c.Subtitle != nil {
		subtitle = *c.Subtitle
	}

	err := r.gdb.Transaction(func(tx *gorm.DB) error {
		if err := tx.Table(entity.ChapterTable).
			Where("comic_id = ? AND deleted_at IS NULL AND pinned = TRUE", c.ComicID).
			Updates(map[string]any{
				"pinned":     false,
				"updated_at": now,
			}).Error; err != nil {
			return err
		}

		row := map[string]any{
			"id":                    c.ID,
			"comic_id":              c.ComicID,
			"pinned":                true,
			"index":                 c.Index,
			"subtitle":              subtitle,
			"page_count":            0,
			"total_unit_count":      0,
			"translated_unit_count": 0,
			"proofread_unit_count":  0,
			"creator_id":            c.CreatorID,
			"created_at":            now,
			"updated_at":            now,
		}

		return tx.Table(entity.ChapterTable).Create(row).Error
	})
	if err != nil {
		return nil, err
	}

	return r.GetByID(c.ID)
}

func (r *chapterRepoImpl) Update(u *model.ChapterUpdate) error {
	return r.gdb.Transaction(func(tx *gorm.DB) error {
		updates := map[string]any{
			"updated_at":      time.Now(),
			"uploaded_at":     u.UploadedAt,
			"transalating_at": u.TransalatingAt,
			"translated_at":   u.TranslatedAt,
			"proofreading_at": u.ProofreadingAt,
			"proofread_at":    u.ProofreadAt,
			"typesetting_at":  u.TypesettingAt,
			"typeset_at":      u.TypesetAt,
			"reviewed_at":     u.ReviewedAt,
			"published_at":    u.PublishedAt,
		}
		if u.Subtitle != nil {
			updates["subtitle"] = *u.Subtitle
		}

		if u.IsPinned != nil {
			if *u.IsPinned {
				var current entity.ChapterInfoRow
				err := tx.Table(entity.ChapterTable).
					Select("id", "comic_id").
					Where("id = ? AND deleted_at IS NULL", u.ID).
					First(&current).Error
				if err != nil {
					return err
				}

				if err := tx.Table(entity.ChapterTable).
					Where("comic_id = ? AND id <> ? AND deleted_at IS NULL AND pinned = TRUE", current.ComicID, u.ID).
					Updates(map[string]any{"pinned": false, "updated_at": time.Now()}).Error; err != nil {
					return err
				}
			}

			updates["pinned"] = *u.IsPinned
		}

		return tx.Table(entity.ChapterTable).
			Where("id = ? AND deleted_at IS NULL", u.ID).
			Updates(updates).Error
	})
}

func (r *chapterRepoImpl) Remove(id string) error {
	return r.gdb.Table(entity.ChapterTable).
		Where("id = ? AND deleted_at IS NULL", id).
		Updates(map[string]any{
			"deleted_at": time.Now(),
			"updated_at": time.Now(),
		}).Error
}

func (r *chapterRepoImpl) UpdateStats(stats *model.ChapterStats) error {
	return r.gdb.Table(entity.ChapterTable).
		Where("id = ? AND deleted_at IS NULL", stats.ChapterID).
		Updates(map[string]any{
			"total_unit_count":      stats.TotalUnitCount,
			"translated_unit_count": stats.TranslatedUnitCount,
			"proofread_unit_count":  stats.ProofreadUnitCount,
			"updated_at":            time.Now(),
		}).Error
}
