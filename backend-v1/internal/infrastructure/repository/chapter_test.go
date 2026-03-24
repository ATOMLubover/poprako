package repository

import (
	"testing"

	"labelplus-next-web-be/internal/domain/model"
	intf "labelplus-next-web-be/internal/domain/repository"
	"labelplus-next-web-be/internal/infrastructure/repository/entity"
	"labelplus-next-web-be/internal/infrastructure/repository/query_option"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type createdChapter struct {
	ID       string
	Subtitle string
}

func createChapters(t *testing.T, repo intf.ChapterRepository, creations []model.ChapterCreation) []createdChapter {
	t.Helper()

	chapters := make([]createdChapter, 0, len(creations))
	for _, creation := range creations {
		id, err := repo.Create(nil, creation)
		require.NoError(t, err)
		chapters = append(chapters, createdChapter{ID: id, Subtitle: creation.Subtitle})
	}

	return chapters
}

func fetchChapterRow(t *testing.T, executor intf.Executor, id string) entity.ChapterInfoRow {
	t.Helper()

	var row entity.ChapterInfoRow
	err := executor.Table(entity.ChapterTable).Where("id = ?", id).First(&row).Error
	require.NoError(t, err)

	return row
}

func TestChapterRepository_LifecycleIncludeAndStats(t *testing.T) {
	runRepositoryTest(t, func(executor intf.Executor) {
		userRepo := NewUserRepository(executor)
		teamRepo := NewTeamRepository(executor)
		worksetRepo := NewWorksetRepository(executor)
		comicRepo := NewComicRepository(executor)
		chapterRepo := NewChapterRepository(executor)

		users := createUsers(t, userRepo, []model.UserCreation{
			model.NewUserRegistration("Chapter Creator One", "520001", "pw-1"),
			model.NewUserRegistration("Chapter Creator Two", "520002", "pw-2"),
		})
		teams := createTeams(t, teamRepo, []model.TeamCreation{
			*model.NewTeamCreation("Chapter Team", "Repository chapter tests"),
		})
		worksets := createWorksets(t, worksetRepo, []model.WorksetCreation{
			*model.NewWorksetCreation(teams[0].ID, 0, "Chapter Shelf", "chapter workset"),
		})
		comics := createComics(t, comicRepo, []model.ComicCreation{
			*model.NewComicCreation(worksets[0].ID, 0, "Chapter Comic", "Author One", "desc one", users[0].ID),
			*model.NewComicCreation(worksets[0].ID, 1, "Chapter Side Comic", "Author Two", "desc two", users[1].ID),
		})

		transactionExecutor := chapterRepo.BeginTransaction()
		require.NoError(t, transactionExecutor.Error)

		err := chapterRepo.LockByComicID(transactionExecutor, comics[0].ID)
		require.NoError(t, err)

		count, err := chapterRepo.Count(transactionExecutor, query_option.ChapterQuery().FilterByComicID(comics[0].ID))
		require.NoError(t, err)
		assert.Zero(t, count)
		require.NoError(t, transactionExecutor.Rollback().Error)

		chapters := createChapters(t, chapterRepo, []model.ChapterCreation{
			model.NewChapterCreation(comics[0].ID, 0, "Opening", users[0].ID),
			model.NewChapterCreation(comics[0].ID, 1, "Middle", users[1].ID),
			model.NewChapterCreation(comics[1].ID, 0, "Side Story", users[1].ID),
		})

		listed, err := chapterRepo.List(nil,
			query_option.ChapterQuery().FilterByComicID(comics[0].ID),
			query_option.ChapterQuery().OrderByIndexDesc(),
			query_option.ChapterQuery().IncludeCreatorInfo(),
		)
		require.NoError(t, err)
		require.Len(t, listed, 2)
		assert.Equal(t, []string{"Middle", "Opening"}, []string{listed[0].Subtitle, listed[1].Subtitle})
		require.NotNil(t, listed[0].Creator)
		assert.Equal(t, users[1].ID, listed[0].Creator.ID)
		assert.Equal(t, "Chapter Creator Two", listed[0].Creator.Name)

		count, err = chapterRepo.Count(nil, query_option.ChapterQuery().FilterByComicID(comics[0].ID))
		require.NoError(t, err)
		assert.EqualValues(t, 2, count)

		plainChapter, err := chapterRepo.Get(nil, query_option.FilterByID(entity.ChapterTable, chapters[0].ID))
		require.NoError(t, err)
		assert.Equal(t, comics[0].ID, plainChapter.ComicID)
		assert.Equal(t, "Opening", plainChapter.Subtitle)
		assert.Nil(t, plainChapter.Creator)

		currentChapter, err := chapterRepo.Get(nil, query_option.FilterByID(entity.ChapterTable, chapters[1].ID))
		require.NoError(t, err)

		updatedSubtitle := "Middle Revised"
		uploadCompleted := model.WorkflowCompleted
		translateInProgress := model.WorkflowInProgress
		err = chapterRepo.Update(nil, model.NewChapterUpdate(
			chapters[1].ID,
			&updatedSubtitle,
			currentChapter,
			&uploadCompleted,
			&translateInProgress,
			nil,
			nil,
			nil,
			nil,
		))
		require.NoError(t, err)

		updatedChapter, err := chapterRepo.Get(nil, query_option.FilterByID(entity.ChapterTable, chapters[1].ID))
		require.NoError(t, err)
		assert.Equal(t, "Middle Revised", updatedChapter.Subtitle)
		assert.NotNil(t, updatedChapter.UploadedAt)
		assert.NotNil(t, updatedChapter.TransalatingAt)
		assert.Nil(t, updatedChapter.TranslatedAt)

		err = chapterRepo.UpdateStats(nil, model.NewChapterStats(chapters[1].ID, 9, 6, 4))
		require.NoError(t, err)

		stats, err := chapterRepo.GetStatsByID(nil, chapters[1].ID)
		require.NoError(t, err)
		assert.Equal(t, chapters[1].ID, stats.ChapterID)
		assert.Equal(t, 9, stats.TotalUnitCount)
		assert.Equal(t, 6, stats.TranslatedUnitCount)
		assert.Equal(t, 4, stats.ProofreadUnitCount)

		transactionExecutor = chapterRepo.BeginTransaction()
		require.NoError(t, transactionExecutor.Error)
		err = chapterRepo.LockByID(transactionExecutor, chapters[1].ID)
		require.NoError(t, err)
		require.NoError(t, transactionExecutor.Rollback().Error)

		err = chapterRepo.Delete(nil, chapters[0].ID)
		require.NoError(t, err)

		deletedRow := fetchChapterRow(t, executor, chapters[0].ID)
		require.NotNil(t, deletedRow.DeletedAt)

		_, err = chapterRepo.Get(nil, query_option.FilterByID(entity.ChapterTable, chapters[0].ID))
		assert.ErrorIs(t, err, ErrRecordNotFound)

		remaining, err := chapterRepo.List(nil,
			query_option.ChapterQuery().FilterByComicID(comics[0].ID),
			query_option.ChapterQuery().OrderByIndexDesc(),
		)
		require.NoError(t, err)
		require.Len(t, remaining, 1)
		assert.Equal(t, chapters[1].ID, remaining[0].ID)
	})
}
