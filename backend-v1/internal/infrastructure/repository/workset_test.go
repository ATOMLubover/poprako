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

type createdWorkset struct {
	ID   string
	Name string
}

func createWorksets(t *testing.T, repo intf.WorksetRepository, creations []model.WorksetCreation) []createdWorkset {
	t.Helper()

	worksets := make([]createdWorkset, 0, len(creations))
	for _, creation := range creations {
		id, err := repo.Create(nil, creation)
		require.NoError(t, err)
		worksets = append(worksets, createdWorkset{ID: id, Name: creation.Name})
	}

	return worksets
}

func fetchWorksetRow(t *testing.T, executor intf.Executor, id string) entity.WorksetInfoRow {
	t.Helper()

	var row entity.WorksetInfoRow
	err := executor.Table(entity.WorksetTable).Where("id = ?", id).First(&row).Error
	require.NoError(t, err)

	return row
}

func TestWorksetRepository_LifecycleAndIncludeTeam(t *testing.T) {
	runRepositoryTest(t, func(executor intf.Executor) {
		teamRepo := NewTeamRepository(executor)
		worksetRepo := NewWorksetRepository(executor)

		teams := createTeams(t, teamRepo, []model.TeamCreation{
			*model.NewTeamCreation("Workset Team Alpha", "Alpha description"),
			*model.NewTeamCreation("Workset Team Beta", "Beta description"),
		})

		transactionExecutor := worksetRepo.BeginTransaction()
		require.NoError(t, transactionExecutor.Error)

		err := worksetRepo.LockByTeamID(transactionExecutor, teams[0].ID)
		require.NoError(t, err)

		count, err := worksetRepo.Count(transactionExecutor, query_option.WorksetQuery().FilterByTeamID(teams[0].ID))
		require.NoError(t, err)
		assert.Zero(t, count)
		require.NoError(t, transactionExecutor.Rollback().Error)

		worksets := createWorksets(t, worksetRepo, []model.WorksetCreation{
			*model.NewWorksetCreation(teams[0].ID, 1, "Volume Two", "Second workset"),
			*model.NewWorksetCreation(teams[0].ID, 0, "Volume One", "First workset"),
			*model.NewWorksetCreation(teams[1].ID, 0, "Beta Shelf", "Other team workset"),
		})

		listed, err := worksetRepo.List(nil,
			query_option.WorksetQuery().FilterByTeamID(teams[0].ID),
			query_option.WorksetQuery().OrderByIndexAsc(),
		)
		require.NoError(t, err)
		require.Len(t, listed, 2)
		assert.Equal(t, []string{"Volume One", "Volume Two"}, []string{listed[0].Name, listed[1].Name})
		assert.Equal(t, []int{0, 1}, []int{listed[0].Index, listed[1].Index})
		assert.Nil(t, listed[0].Team)

		count, err = worksetRepo.Count(nil, query_option.WorksetQuery().FilterByTeamID(teams[0].ID))
		require.NoError(t, err)
		assert.EqualValues(t, 2, count)

		included, err := worksetRepo.List(nil,
			query_option.WorksetQuery().FilterByTeamID(teams[0].ID),
			query_option.WorksetQuery().OrderByIndexAsc(),
			query_option.WorksetQuery().IncludeTeamInfo(),
		)
		require.NoError(t, err)
		require.Len(t, included, 2)
		require.NotNil(t, included[0].Team)
		assert.Equal(t, teams[0].ID, included[0].Team.ID)
		assert.Equal(t, "Workset Team Alpha", included[0].Team.Name)
		assert.Equal(t, "Alpha description", included[0].Team.Description)

		selected, err := worksetRepo.Get(nil, query_option.FilterByID(entity.WorksetTable, worksets[0].ID))
		require.NoError(t, err)
		assert.Equal(t, teams[0].ID, selected.TeamID)
		assert.Equal(t, "Volume Two", selected.Name)
		assert.Equal(t, "Second workset", selected.Description)

		err = worksetRepo.Update(nil, model.NewWorksetUpdate(worksets[1].ID, "Volume One Revised", nil))
		require.NoError(t, err)

		updated, err := worksetRepo.Get(nil, query_option.FilterByID(entity.WorksetTable, worksets[1].ID))
		require.NoError(t, err)
		assert.Equal(t, "Volume One Revised", updated.Name)
		assert.Equal(t, "", updated.Description)

		updatedRow := fetchWorksetRow(t, executor, worksets[1].ID)
		assert.False(t, updatedRow.Description.Valid)

		err = worksetRepo.Delete(nil, worksets[0].ID)
		require.NoError(t, err)

		_, err = worksetRepo.Get(nil, query_option.FilterByID(entity.WorksetTable, worksets[0].ID))
		assert.ErrorIs(t, err, ErrRecordNotFound)

		remaining, err := worksetRepo.List(nil,
			query_option.WorksetQuery().FilterByTeamID(teams[0].ID),
			query_option.WorksetQuery().OrderByIndexAsc(),
		)
		require.NoError(t, err)
		require.Len(t, remaining, 1)
		assert.Equal(t, worksets[1].ID, remaining[0].ID)
	})
}
