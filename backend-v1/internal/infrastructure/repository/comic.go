package repository

import (
	"errors"
	"fmt"
	"time"

	"labelplus-next-web-be/internal/domain/model"
	intf "labelplus-next-web-be/internal/domain/repository"
	"labelplus-next-web-be/internal/infrastructure/repository/entity"
	"labelplus-next-web-be/internal/util"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type comicRepository struct {
	executor intf.Executor
}

func NewComicRepository(executor intf.Executor) intf.ComicRepository {
	return &comicRepository{executor: executor}
}

func composeComicTitle(index int, author string, title string) string {
	return fmt.Sprintf("【%d】[%s]%s", index, author, title)
}

func (r *comicRepository) withTransaction(executor intf.Executor) intf.Executor {
	if executor != nil {
		return executor
	}

	return r.executor
}

func (r *comicRepository) BeginTransaction() intf.Executor {
	return r.executor.Begin()
}

func (r *comicRepository) List(executor intf.Executor, options ...intf.QueryOption) ([]model.ComicInfo, error) {
	executor = r.withTransaction(executor)
	executor = executor.Table(entity.ComicTable).Where("comic_table.deleted_at IS NULL")

	for _, opt := range options {
		executor = opt(executor)
	}

	var rows []entity.ComicInfoRow

	if err := executor.Find(&rows).Error; err != nil {
		return nil, err
	}

	result := make([]model.ComicInfo, len(rows))
	for i, row := range rows {
		result[i] = entity.ToComicInfo(row)
	}

	return result, nil
}

func (r *comicRepository) Get(executor intf.Executor, options ...intf.QueryOption) (model.ComicInfo, error) {
	executor = r.withTransaction(executor)
	executor = executor.Table(entity.ComicTable).Where("comic_table.deleted_at IS NULL")

	for _, opt := range options {
		executor = opt(executor)
	}

	var row entity.ComicInfoRow
	if err := executor.First(&row).Error; err != nil {
		return model.ComicInfo{}, err
	}

	return entity.ToComicInfo(row), nil
}

func (r *comicRepository) GetLatestChapterFirstPageOSSKey(executor intf.Executor, comicID string) (*string, error) {
	executor = r.withTransaction(executor)

	type pageOSSKeyRow struct {
		OSSKey string `gorm:"column:oss_key"`
	}

	var row pageOSSKeyRow
	err := executor.
		Table(entity.PageTable+" AS page_table").
		Select("page_table.oss_key").
		Joins("JOIN chapter_table ON chapter_table.id = page_table.chapter_id AND chapter_table.deleted_at IS NULL").
		Where("chapter_table.comic_id = ?", comicID).
		Order("chapter_table.index DESC").
		Order("page_table.index ASC").
		Limit(1).
		Take(&row).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return &row.OSSKey, nil
}

func (r *comicRepository) LockByWorksetID(executor intf.Executor, worksetID string) error {
	executor = r.withTransaction(executor)

	var lockedIDs []string

	return executor.
		Table(entity.ComicTable).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("workset_id = ? AND deleted_at IS NULL", worksetID).
		Pluck("id", &lockedIDs).Error
}

func (r *comicRepository) Count(executor intf.Executor, options ...intf.QueryOption) (int64, error) {
	executor = r.withTransaction(executor)
	executor = executor.Table(entity.ComicTable).Where("deleted_at IS NULL")
	for _, opt := range options {
		executor = opt(executor)
	}

	var count int64
	if err := executor.Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (r *comicRepository) SyncLatestChapterReplica(executor intf.Executor, comicID string) error {
	executor = r.withTransaction(executor)

	return executor.Exec(`
UPDATE comic_table AS comic
SET
	chapter_count = chapter_stats.chapter_count,
	latest_uploaded_at = latest_chapter.uploaded_at,
	latest_transalating_at = latest_chapter.transalating_at,
	latest_translated_at = latest_chapter.translated_at,
	latest_proofreading_at = latest_chapter.proofreading_at,
	latest_proofread_at = latest_chapter.proofread_at,
	latest_typesetting_at = latest_chapter.typesetting_at,
	latest_typeset_at = latest_chapter.typeset_at,
	latest_reviewed_at = latest_chapter.reviewed_at,
	latest_published_at = latest_chapter.published_at
FROM (
	SELECT COUNT(*)::INTEGER AS chapter_count
	FROM chapter_table
	WHERE chapter_table.comic_id = ?
		AND chapter_table.deleted_at IS NULL
) AS chapter_stats
LEFT JOIN LATERAL (
	SELECT
		uploaded_at,
		transalating_at,
		translated_at,
		proofreading_at,
		proofread_at,
		typesetting_at,
		typeset_at,
		reviewed_at,
		published_at
	FROM chapter_table
	WHERE chapter_table.comic_id = ?
		AND chapter_table.deleted_at IS NULL
	ORDER BY chapter_table.index DESC
	LIMIT 1
) AS latest_chapter ON TRUE
WHERE comic.id = ?
	AND comic.deleted_at IS NULL
`, comicID, comicID, comicID).Error
}

func (r *comicRepository) Create(executor intf.Executor, creation model.ComicCreation) (string, error) {
	executor = r.withTransaction(executor)

	row := entity.ComicInsertRow{
		ID:            util.GenerateUUID(),
		WorksetID:     creation.WorksetID,
		Index:         creation.Index,
		Title:         creation.Title,
		Author:        creation.Author,
		ComposedTitle: composeComicTitle(creation.Index, creation.Author, creation.Title),
		Description:   creation.Description,
		CreatorID:     creation.CreatorID,
	}

	if err := executor.Create(&row).Error; err != nil {
		return "", err
	}

	return row.ID, nil
}

func (r *comicRepository) Update(executor intf.Executor, update model.ComicUpdate) error {
	executor = r.withTransaction(executor)

	return executor.
		Table(entity.ComicTable).
		Where("id = ? AND deleted_at IS NULL", update.ID).
		Updates(map[string]any{
			"title":          update.Title,
			"author":         update.Author,
			"composed_title": gorm.Expr("CONCAT('【', \"index\", '】[', ?, ']', ?)", update.Author, update.Title),
			"description":    update.Description,
		}).Error
}

func (r *comicRepository) Delete(executor intf.Executor, comicID string) error {
	executor = r.withTransaction(executor)

	return executor.
		Table(entity.ComicTable).
		Where("id = ?", comicID).
		Update("deleted_at", time.Now()).Error
}
