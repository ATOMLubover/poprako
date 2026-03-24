package repository

import (
	"testing"

	"labelplus-next-web-be/internal/domain/model"
	intf "labelplus-next-web-be/internal/domain/repository"
	"labelplus-next-web-be/internal/infrastructure/repository/entity"
	"labelplus-next-web-be/internal/infrastructure/repository/query_option"
	"labelplus-next-web-be/internal/util"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func fetchPageRow(t *testing.T, executor intf.Executor, id string) entity.PageInfoRow {
	t.Helper()

	var row entity.PageInfoRow
	err := executor.Table(entity.PageTable).Where("id = ?", id).First(&row).Error
	require.NoError(t, err)

	return row
}

func TestPageRepository_LifecycleIncludeAndStats(t *testing.T) {
	runRepositoryTest(t, func(executor intf.Executor) {
		userRepo := NewUserRepository(executor)
		teamRepo := NewTeamRepository(executor)
		worksetRepo := NewWorksetRepository(executor)
		comicRepo := NewComicRepository(executor)
		pageRepo := NewPageRepository(executor)

		users := createUsers(t, userRepo, []model.UserCreation{
			model.NewUserRegistration("Page Creator One", "530001", "pw-1"),
			model.NewUserRegistration("Page Creator Two", "530002", "pw-2"),
		})
		teams := createTeams(t, teamRepo, []model.TeamCreation{
			*model.NewTeamCreation("Page Team", "Repository page tests"),
		})
		worksets := createWorksets(t, worksetRepo, []model.WorksetCreation{
			*model.NewWorksetCreation(teams[0].ID, 0, "Page Shelf", "page workset"),
		})
		comics := createComics(t, comicRepo, []model.ComicCreation{
			*model.NewComicCreation(worksets[0].ID, 0, "Page Comic", "Page Author", "page desc", users[0].ID),
		})

		firstChapterID := createChapterRow(t, executor, comics[0].ID, 0, "Page Chapter A", users[0].ID)
		secondChapterID := createChapterRow(t, executor, comics[0].ID, 1, "Page Chapter B", users[1].ID)

		transactionExecutor := pageRepo.BeginTransaction()
		require.NoError(t, transactionExecutor.Error)
		err := pageRepo.LockByChapterID(transactionExecutor, firstChapterID)
		require.NoError(t, err)
		require.NoError(t, transactionExecutor.Rollback().Error)

		pageCreations := []model.PageCreation{
			model.NewPageCreation(util.GenerateUUID(), firstChapterID, 1, "oss://chapter-a/page-1.png", users[0].ID),
			model.NewPageCreation(util.GenerateUUID(), firstChapterID, 0, "oss://chapter-a/page-0.png", users[1].ID),
			model.NewPageCreation(util.GenerateUUID(), firstChapterID, 2, "oss://chapter-a/page-2.png", users[0].ID),
			model.NewPageCreation(util.GenerateUUID(), secondChapterID, 0, "oss://chapter-b/page-0.png", users[1].ID),
		}
		err = pageRepo.CreateBatch(nil, pageCreations)
		require.NoError(t, err)

		listed, err := pageRepo.List(nil,
			query_option.PageQuery().FilterByChapterID(firstChapterID),
			query_option.PageQuery().OrderByIndexAsc(),
			query_option.PageQuery().IncludeCreatorInfo(),
		)
		require.NoError(t, err)
		require.Len(t, listed, 3)
		assert.Equal(t, []int{0, 1, 2}, []int{listed[0].Index, listed[1].Index, listed[2].Index})
		require.NotNil(t, listed[0].Creator)
		assert.Equal(t, users[1].ID, listed[0].Creator.ID)
		assert.Equal(t, "Page Creator Two", listed[0].Creator.Name)

		plainPage, err := pageRepo.Get(nil, query_option.FilterByID(entity.PageTable, pageCreations[0].ID))
		require.NoError(t, err)
		assert.Equal(t, firstChapterID, plainPage.ChapterID)
		assert.Equal(t, 1, plainPage.Index)
		assert.Nil(t, plainPage.Creator)

		err = pageRepo.Update(nil, model.NewPageUpdate(pageCreations[0].ID, 5, "oss://chapter-a/page-5.png", true, 8, 6, 4))
		require.NoError(t, err)

		updatedPage, err := pageRepo.Get(nil, query_option.FilterByID(entity.PageTable, pageCreations[0].ID))
		require.NoError(t, err)
		assert.Equal(t, 5, updatedPage.Index)
		assert.Equal(t, "oss://chapter-a/page-5.png", updatedPage.OSSKey)
		assert.True(t, updatedPage.IsUploaded)
		assert.Equal(t, 8, updatedPage.TotalUnitCount)
		assert.Equal(t, 6, updatedPage.TranslatedUnitCount)
		assert.Equal(t, 4, updatedPage.ProofreadUnitCount)

		pageStats, err := pageRepo.GetStatsByID(nil, pageCreations[0].ID)
		require.NoError(t, err)
		assert.Equal(t, 8, pageStats.TotalUnitCount)
		assert.Equal(t, 6, pageStats.TranslatedUnitCount)
		assert.Equal(t, 4, pageStats.ProofreadUnitCount)

		err = pageRepo.UpdateStats(nil, model.NewPageStats(pageCreations[1].ID, 3, 2, 1))
		require.NoError(t, err)

		statsUpdatedPage, err := pageRepo.Get(nil, query_option.FilterByID(entity.PageTable, pageCreations[1].ID))
		require.NoError(t, err)
		assert.Equal(t, 3, statsUpdatedPage.TotalUnitCount)
		assert.Equal(t, 2, statsUpdatedPage.TranslatedUnitCount)
		assert.Equal(t, 1, statsUpdatedPage.ProofreadUnitCount)

		updatedRow := fetchPageRow(t, executor, pageCreations[0].ID)
		assert.Equal(t, 5, updatedRow.Index)

		transactionExecutor = pageRepo.BeginTransaction()
		require.NoError(t, transactionExecutor.Error)
		err = pageRepo.LockByID(transactionExecutor, pageCreations[0].ID)
		require.NoError(t, err)
		require.NoError(t, transactionExecutor.Rollback().Error)

		err = pageRepo.Delete(nil, pageCreations[2].ID)
		require.NoError(t, err)

		_, err = pageRepo.Get(nil, query_option.FilterByID(entity.PageTable, pageCreations[2].ID))
		assert.ErrorIs(t, err, ErrRecordNotFound)

		err = pageRepo.DeleteBatch(nil, []string{pageCreations[0].ID, pageCreations[1].ID, "missing-page"})
		require.NoError(t, err)

		remainingFirstChapterPages, err := pageRepo.List(nil,
			query_option.PageQuery().FilterByChapterID(firstChapterID),
			query_option.PageQuery().OrderByIndexAsc(),
		)
		require.NoError(t, err)
		assert.Empty(t, remainingFirstChapterPages)

		remainingSecondChapterPages, err := pageRepo.List(nil,
			query_option.PageQuery().FilterByChapterID(secondChapterID),
			query_option.PageQuery().OrderByIndexAsc(),
		)
		require.NoError(t, err)
		require.Len(t, remainingSecondChapterPages, 1)
		assert.Equal(t, pageCreations[3].ID, remainingSecondChapterPages[0].ID)
	})
}
