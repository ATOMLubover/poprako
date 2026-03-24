package repository

import (
	"testing"

	"labelplus-next-web-be/internal/domain/model"
	intf "labelplus-next-web-be/internal/domain/repository"
	"labelplus-next-web-be/internal/infrastructure/repository/query_option"
	"labelplus-next-web-be/internal/util"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnitRepository_BatchLifecycleAndIdempotency(t *testing.T) {
	runRepositoryTest(t, func(executor intf.Executor) {
		userRepo := NewUserRepository(executor)
		teamRepo := NewTeamRepository(executor)
		worksetRepo := NewWorksetRepository(executor)
		comicRepo := NewComicRepository(executor)
		unitRepo := NewUnitRepository(executor)

		users := createUsers(t, userRepo, []model.UserCreation{
			model.NewUserRegistration("Unit User One", "550001", "pw-1"),
			model.NewUserRegistration("Unit User Two", "550002", "pw-2"),
		})
		teams := createTeams(t, teamRepo, []model.TeamCreation{
			*model.NewTeamCreation("Unit Team", "Repository unit tests"),
		})
		worksets := createWorksets(t, worksetRepo, []model.WorksetCreation{
			*model.NewWorksetCreation(teams[0].ID, 0, "Unit Shelf", "unit workset"),
		})
		comics := createComics(t, comicRepo, []model.ComicCreation{
			*model.NewComicCreation(worksets[0].ID, 0, "Unit Comic", "Unit Author", "unit desc", users[0].ID),
		})

		chapterID := createChapterRow(t, executor, comics[0].ID, 0, "Unit Chapter", users[0].ID)
		pageID := createPageRow(t, executor, chapterID, 0, "oss://unit/page-0.png", users[0].ID)
		otherPageID := createPageRow(t, executor, chapterID, 1, "oss://unit/page-1.png", users[1].ID)

		transactionExecutor := unitRepo.BeginTransaction()
		require.NoError(t, transactionExecutor.Error)
		err := unitRepo.LockByPageID(transactionExecutor, pageID)
		require.NoError(t, err)
		require.NoError(t, transactionExecutor.Rollback().Error)

		translatedText := "translated text"
		translatorComment := "translator comment"
		proofreadText := "proofread text"
		proofreaderComment := "proofreader comment"

		unitCreations := []model.UnitCreation{
			model.NewUnitCreation("unit-two", 2, 20, 30, true, nil, nil, nil, false, nil, nil, nil),
			model.NewUnitCreation("unit-zero", 0, 5, 10, true, nil, nil, nil, false, nil, nil, nil),
			model.NewUnitCreation("unit-one", 1, 12, 18, false, &translatedText, &users[0].ID, &translatorComment, false, nil, nil, nil),
			model.NewUnitCreation("unit-other-page", 0, 99, 77, true, &proofreadText, &users[1].ID, nil, true, &proofreadText, &users[1].ID, &proofreaderComment),
		}
		unitCreations[0].PageID = pageID
		unitCreations[1].PageID = pageID
		unitCreations[2].PageID = pageID
		unitCreations[3].PageID = otherPageID

		err = unitRepo.CreateBatch(nil, unitCreations)
		require.NoError(t, err)

		listed, err := unitRepo.List(nil,
			query_option.UnitQuery().FilterByPageID(pageID),
			query_option.UnitQuery().OrderByIndexAsc(),
		)
		require.NoError(t, err)
		require.Len(t, listed, 3)
		assert.Equal(t, []string{"unit-zero", "unit-one", "unit-two"}, []string{listed[0].ID, listed[1].ID, listed[2].ID})
		assert.Equal(t, []int{0, 1, 2}, []int{listed[0].Index, listed[1].Index, listed[2].Index})

		var clearedText *string
		var clearedUserID *string
		var clearedComment *string

		err = unitRepo.PatchBatch(nil, []model.UnitPatch{
			model.NewUnitPatch(
				"unit-zero",
				util.NewSomeOption(7),
				util.NewSomeOption(100),
				util.NewSomeOption(200),
				util.NewSomeOption(false),
				util.NewSomeOption(&translatedText),
				util.NewSomeOption(&users[0].ID),
				util.NewSomeOption(&translatorComment),
				util.NewSomeOption(true),
				util.NewSomeOption(&proofreadText),
				util.NewSomeOption(&users[1].ID),
				util.NewSomeOption(&proofreaderComment),
			),
			model.NewUnitPatch(
				"unit-one",
				util.NewNoneOption[int](),
				util.NewNoneOption[int](),
				util.NewNoneOption[int](),
				util.NewNoneOption[bool](),
				util.NewSomeOption(clearedText),
				util.NewSomeOption(clearedUserID),
				util.NewSomeOption(clearedComment),
				util.NewSomeOption(false),
				util.NewSomeOption(clearedText),
				util.NewSomeOption(clearedUserID),
				util.NewSomeOption(clearedComment),
			),
			model.NewUnitPatch(
				"missing-unit",
				util.NewSomeOption(9),
				util.NewNoneOption[int](),
				util.NewNoneOption[int](),
				util.NewNoneOption[bool](),
				util.NewNoneOption[*string](),
				util.NewNoneOption[*string](),
				util.NewNoneOption[*string](),
				util.NewNoneOption[bool](),
				util.NewNoneOption[*string](),
				util.NewNoneOption[*string](),
				util.NewNoneOption[*string](),
			),
		})
		require.NoError(t, err)

		listed, err = unitRepo.List(nil,
			query_option.UnitQuery().FilterByPageID(pageID),
			query_option.UnitQuery().OrderByIndexAsc(),
		)
		require.NoError(t, err)
		require.Len(t, listed, 3)
		assert.Equal(t, []string{"unit-one", "unit-two", "unit-zero"}, []string{listed[0].ID, listed[1].ID, listed[2].ID})
		assert.Nil(t, listed[0].TranslatedText)
		assert.False(t, listed[0].IsProofread)
		assert.Equal(t, 7, listed[2].Index)
		assert.Equal(t, 100, listed[2].XCoord)
		assert.Equal(t, 200, listed[2].YCoord)
		assert.False(t, listed[2].IsBubble)
		require.NotNil(t, listed[2].TranslatedText)
		assert.Equal(t, translatedText, *listed[2].TranslatedText)
		require.NotNil(t, listed[2].ProofreadText)
		assert.Equal(t, proofreadText, *listed[2].ProofreadText)
		assert.True(t, listed[2].IsProofread)

		err = unitRepo.DeleteBatch(nil, []string{"unit-two", "missing-unit"})
		require.NoError(t, err)

		remainingPageUnits, err := unitRepo.List(nil,
			query_option.UnitQuery().FilterByPageID(pageID),
			query_option.UnitQuery().OrderByIndexAsc(),
		)
		require.NoError(t, err)
		require.Len(t, remainingPageUnits, 2)
		assert.Equal(t, []string{"unit-one", "unit-zero"}, []string{remainingPageUnits[0].ID, remainingPageUnits[1].ID})

		otherPageUnits, err := unitRepo.List(nil,
			query_option.UnitQuery().FilterByPageID(otherPageID),
			query_option.UnitQuery().OrderByIndexAsc(),
		)
		require.NoError(t, err)
		require.Len(t, otherPageUnits, 1)
		assert.Equal(t, "unit-other-page", otherPageUnits[0].ID)
	})
}
