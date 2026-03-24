package repository

import (
	"testing"
	"time"

	"labelplus-next-web-be/internal/domain/model"
	intf "labelplus-next-web-be/internal/domain/repository"
	"labelplus-next-web-be/internal/infrastructure/repository/entity"
	"labelplus-next-web-be/internal/infrastructure/repository/query_option"
	"labelplus-next-web-be/internal/util"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type createdComic struct {
	ID    string
	Title string
}

type comicPersistedRow struct {
	ID            string     `gorm:"column:id"`
	Title         string     `gorm:"column:title"`
	Author        string     `gorm:"column:author"`
	ComposedTitle string     `gorm:"column:composed_title"`
	Description   string     `gorm:"column:description"`
	DeletedAt     *time.Time `gorm:"column:deleted_at"`
}

func createComics(t *testing.T, repo intf.ComicRepository, creations []model.ComicCreation) []createdComic {
	t.Helper()

	comics := make([]createdComic, 0, len(creations))
	for _, creation := range creations {
		id, err := repo.Create(nil, creation)
		require.NoError(t, err)
		comics = append(comics, createdComic{ID: id, Title: creation.Title})
	}

	return comics
}

func fetchComicRow(t *testing.T, executor intf.Executor, id string) comicPersistedRow {
	t.Helper()

	var row comicPersistedRow
	err := executor.Table(entity.ComicTable).Where("id = ?", id).First(&row).Error
	require.NoError(t, err)

	return row
}

func createChapterRow(t *testing.T, executor intf.Executor, comicID string, index int, subtitle string, creatorID string) string {
	t.Helper()

	id := util.GenerateUUID()
	err := executor.Create(&entity.ChapterInsertRow{
		ID:        id,
		ComicID:   comicID,
		Index:     index,
		Subtitle:  subtitle,
		CreatorID: creatorID,
	}).Error
	require.NoError(t, err)

	return id
}

func createPageRow(t *testing.T, executor intf.Executor, chapterID string, index int, ossKey string, creatorID string) string {
	t.Helper()

	id := util.GenerateUUID()
	err := executor.Create(&entity.PageInsertRow{
		ID:        id,
		ChapterID: chapterID,
		Index:     index,
		OSSKey:    ossKey,
		CreatorID: creatorID,
	}).Error
	require.NoError(t, err)

	return id
}

func TestComicRepository_LifecycleFiltersAndReplicaQueries(t *testing.T) {
	runRepositoryTest(t, func(executor intf.Executor) {
		userRepo := NewUserRepository(executor)
		teamRepo := NewTeamRepository(executor)
		worksetRepo := NewWorksetRepository(executor)
		comicRepo := NewComicRepository(executor)

		users := createUsers(t, userRepo, []model.UserCreation{
			model.NewUserRegistration("Comic Creator One", "510001", "pw-1"),
			model.NewUserRegistration("Comic Creator Two", "510002", "pw-2"),
		})
		teams := createTeams(t, teamRepo, []model.TeamCreation{
			*model.NewTeamCreation("Comic Team", "Repository test team"),
		})
		worksets := createWorksets(t, worksetRepo, []model.WorksetCreation{
			*model.NewWorksetCreation(teams[0].ID, 0, "Main Shelf", "Primary workset"),
			*model.NewWorksetCreation(teams[0].ID, 1, "Side Shelf", "Secondary workset"),
		})

		transactionExecutor := comicRepo.BeginTransaction()
		require.NoError(t, transactionExecutor.Error)

		err := comicRepo.LockByWorksetID(transactionExecutor, worksets[0].ID)
		require.NoError(t, err)

		count, err := comicRepo.Count(transactionExecutor, query_option.ComicQuery().FilterByWorksetID(worksets[0].ID))
		require.NoError(t, err)
		assert.Zero(t, count)
		require.NoError(t, transactionExecutor.Rollback().Error)

		comics := createComics(t, comicRepo, []model.ComicCreation{
			*model.NewComicCreation(worksets[0].ID, 0, "Alpha Mission", "Aki", "alpha desc", users[0].ID),
			*model.NewComicCreation(worksets[0].ID, 1, "Beta Mission", "Beni", "beta desc", users[0].ID),
			*model.NewComicCreation(worksets[0].ID, 2, "Gamma Mission", "Cara", "gamma desc", users[1].ID),
			*model.NewComicCreation(worksets[1].ID, 0, "Delta Mission", "Dino", "delta desc", users[1].ID),
		})

		baseTime := time.Date(2026, time.March, 1, 8, 0, 0, 0, time.UTC)
		require.NoError(t, executor.Table(entity.ComicTable).Where("id = ?", comics[0].ID).Update("last_active_at", baseTime.Add(1*time.Hour)).Error)
		require.NoError(t, executor.Table(entity.ComicTable).Where("id = ?", comics[1].ID).Update("last_active_at", baseTime.Add(3*time.Hour)).Error)
		require.NoError(t, executor.Table(entity.ComicTable).Where("id = ?", comics[2].ID).Update("last_active_at", baseTime.Add(2*time.Hour)).Error)
		require.NoError(t, executor.Table(entity.ComicTable).Where("id = ?", comics[3].ID).Update("last_active_at", baseTime.Add(4*time.Hour)).Error)

		listed, err := comicRepo.List(nil,
			query_option.ComicQuery().FilterByWorksetID(worksets[0].ID),
			query_option.ComicQuery().OrderByLastActiveAtDesc(),
			query_option.ComicQuery().IncludeWorksetAndCreatorInfo(),
		)
		require.NoError(t, err)
		require.Len(t, listed, 3)
		assert.Equal(t, []string{"Beta Mission", "Gamma Mission", "Alpha Mission"}, []string{listed[0].Title, listed[1].Title, listed[2].Title})
		require.NotNil(t, listed[0].Workset)
		require.NotNil(t, listed[0].Creator)
		assert.Equal(t, worksets[0].ID, listed[0].Workset.ID)
		assert.Equal(t, "Main Shelf", listed[0].Workset.Name)
		assert.Equal(t, users[0].ID, listed[0].Creator.ID)
		assert.Equal(t, "Comic Creator One", listed[0].Creator.Name)

		creatorIncluded, err := comicRepo.List(nil,
			query_option.ComicQuery().FilterByWorksetID(worksets[0].ID),
			query_option.ComicQuery().FilterByCreatorID(users[1].ID),
			query_option.ComicQuery().IncludeCreatorInfo(),
		)
		require.NoError(t, err)
		require.Len(t, creatorIncluded, 1)
		assert.Nil(t, creatorIncluded[0].Workset)
		require.NotNil(t, creatorIncluded[0].Creator)
		assert.Equal(t, "Comic Creator Two", creatorIncluded[0].Creator.Name)

		worksetIncluded, err := comicRepo.List(nil,
			query_option.ComicQuery().FilterByWorksetID(worksets[0].ID),
			query_option.ComicQuery().FilterByCreatorID(users[0].ID),
			query_option.ComicQuery().IncludeWorksetInfo(),
		)
		require.NoError(t, err)
		require.Len(t, worksetIncluded, 2)
		require.NotNil(t, worksetIncluded[0].Workset)
		assert.Nil(t, worksetIncluded[0].Creator)
		assert.Equal(t, "Main Shelf", worksetIncluded[0].Workset.Name)

		fuzzyMatched, err := comicRepo.List(nil,
			query_option.ComicQuery().FilterByWorksetID(worksets[0].ID),
			query_option.ComicQuery().FuzzyFilterByTitle("[Beni] Beta"),
		)
		require.NoError(t, err)
		require.Len(t, fuzzyMatched, 1)
		assert.Equal(t, comics[1].ID, fuzzyMatched[0].ID)

		firstChapterID := createChapterRow(t, executor, comics[0].ID, 0, "v1", users[0].ID)
		latestChapterID := createChapterRow(t, executor, comics[0].ID, 1, "v2", users[0].ID)
		createPageRow(t, executor, firstChapterID, 0, "oss://comic-alpha/chapter-0-page-0.png", users[0].ID)
		createPageRow(t, executor, latestChapterID, 1, "oss://comic-alpha/chapter-1-page-1.png", users[0].ID)
		createPageRow(t, executor, latestChapterID, 0, "oss://comic-alpha/chapter-1-page-0.png", users[0].ID)

		require.NoError(t, executor.Table(entity.ChapterTable).Where("id = ?", firstChapterID).Updates(map[string]any{
			"uploaded_at": baseTime.Add(10 * time.Minute),
		}).Error)
		require.NoError(t, executor.Table(entity.ChapterTable).Where("id = ?", latestChapterID).Updates(map[string]any{
			"uploaded_at":     baseTime.Add(20 * time.Minute),
			"transalating_at": baseTime.Add(30 * time.Minute),
		}).Error)

		completedChapterID := createChapterRow(t, executor, comics[1].ID, 0, "final", users[0].ID)
		createPageRow(t, executor, completedChapterID, 0, "oss://comic-beta/chapter-0-page-0.png", users[0].ID)
		require.NoError(t, executor.Table(entity.ChapterTable).Where("id = ?", completedChapterID).Updates(map[string]any{
			"uploaded_at":     baseTime.Add(40 * time.Minute),
			"transalating_at": baseTime.Add(50 * time.Minute),
			"translated_at":   baseTime.Add(60 * time.Minute),
			"proofreading_at": baseTime.Add(70 * time.Minute),
			"proofread_at":    baseTime.Add(80 * time.Minute),
			"typesetting_at":  baseTime.Add(90 * time.Minute),
			"typeset_at":      baseTime.Add(100 * time.Minute),
			"reviewed_at":     baseTime.Add(110 * time.Minute),
			"published_at":    baseTime.Add(120 * time.Minute),
		}).Error)

		err = comicRepo.SyncLatestChapterReplica(nil, comics[0].ID)
		require.NoError(t, err)
		err = comicRepo.SyncLatestChapterReplica(nil, comics[1].ID)
		require.NoError(t, err)
		err = comicRepo.SyncLatestChapterReplica(nil, comics[2].ID)
		require.NoError(t, err)

		translateInProgress := model.WorkflowInProgress
		translateMatched, err := comicRepo.List(nil,
			query_option.ComicQuery().FilterByWorksetID(worksets[0].ID),
			query_option.ComicQuery().FilterByLatestChapterStatus(nil, &translateInProgress, nil, nil, nil, nil),
		)
		require.NoError(t, err)
		require.Len(t, translateMatched, 1)
		assert.Equal(t, comics[0].ID, translateMatched[0].ID)

		publishedCompleted := model.WorkflowCompleted
		publishedMatched, err := comicRepo.List(nil,
			query_option.ComicQuery().FilterByWorksetID(worksets[0].ID),
			query_option.ComicQuery().FilterByLatestChapterStatus(nil, nil, nil, nil, nil, &publishedCompleted),
		)
		require.NoError(t, err)
		require.Len(t, publishedMatched, 1)
		assert.Equal(t, comics[1].ID, publishedMatched[0].ID)

		uploadPending := model.WorkflowPending
		pendingMatched, err := comicRepo.List(nil,
			query_option.ComicQuery().FilterByWorksetID(worksets[0].ID),
			query_option.ComicQuery().FilterByLatestChapterStatus(&uploadPending, nil, nil, nil, nil, nil),
		)
		require.NoError(t, err)
		require.Len(t, pendingMatched, 1)
		assert.Equal(t, comics[2].ID, pendingMatched[0].ID)

		count, err = comicRepo.Count(nil, query_option.ComicQuery().FilterByWorksetID(worksets[0].ID))
		require.NoError(t, err)
		assert.EqualValues(t, 3, count)

		err = comicRepo.Update(nil, model.NewComicUpdate(comics[1].ID, "Beta Mission Revised", "Beni Revised", "beta updated desc", "【1】[Beni Revised] Beta Mission Revised"))
		require.NoError(t, err)

		updatedRow := fetchComicRow(t, executor, comics[1].ID)
		assert.Equal(t, "Beta Mission Revised", updatedRow.Title)
		assert.Equal(t, "Beni Revised", updatedRow.Author)
		assert.Equal(t, "【1】[Beni Revised] Beta Mission Revised", updatedRow.ComposedTitle)
		assert.Equal(t, "beta updated desc", updatedRow.Description)

		plainComic, err := comicRepo.Get(nil, query_option.FilterByID(entity.ComicTable, comics[1].ID))
		require.NoError(t, err)
		assert.Nil(t, plainComic.Workset)
		assert.Nil(t, plainComic.Creator)

		coverOSSKey, err := comicRepo.GetLatestChapterFirstPageOSSKey(nil, comics[0].ID)
		require.NoError(t, err)
		require.NotNil(t, coverOSSKey)
		assert.Equal(t, "oss://comic-alpha/chapter-1-page-0.png", *coverOSSKey)

		coverOSSKey, err = comicRepo.GetLatestChapterFirstPageOSSKey(nil, comics[2].ID)
		require.NoError(t, err)
		assert.Nil(t, coverOSSKey)

		err = comicRepo.Delete(nil, comics[0].ID)
		require.NoError(t, err)

		deletedRow := fetchComicRow(t, executor, comics[0].ID)
		require.NotNil(t, deletedRow.DeletedAt)

		_, err = comicRepo.Get(nil, query_option.FilterByID(entity.ComicTable, comics[0].ID))
		assert.ErrorIs(t, err, ErrRecordNotFound)

		remaining, err := comicRepo.List(nil,
			query_option.ComicQuery().FilterByWorksetID(worksets[0].ID),
			query_option.ComicQuery().OrderByLastActiveAtDesc(),
		)
		require.NoError(t, err)
		require.Len(t, remaining, 2)

		count, err = comicRepo.Count(nil, query_option.ComicQuery().FilterByWorksetID(worksets[0].ID))
		require.NoError(t, err)
		assert.EqualValues(t, 2, count)
	})
}
