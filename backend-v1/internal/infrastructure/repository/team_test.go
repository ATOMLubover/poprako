package repository

import (
	"sort"
	"testing"

	"labelplus-next-web-be/internal/domain/model"
	intf "labelplus-next-web-be/internal/domain/repository"
	"labelplus-next-web-be/internal/infrastructure/repository/entity"
	"labelplus-next-web-be/internal/infrastructure/repository/query_option"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type createdTeam struct {
	ID   string
	Name string
}

func createTeams(t *testing.T, repo intf.TeamRepository, creations []model.TeamCreation) []createdTeam {
	t.Helper()

	teams := make([]createdTeam, 0, len(creations))
	for _, creation := range creations {
		id, err := repo.Create(nil, creation)
		require.NoError(t, err)
		teams = append(teams, createdTeam{ID: id, Name: creation.Name})
	}

	return teams
}

func fetchTeamRow(t *testing.T, executor intf.Executor, id string) entity.TeamInfoRow {
	t.Helper()

	var row entity.TeamInfoRow
	err := executor.Table(entity.TeamTable).Where("id = ?", id).First(&row).Error
	require.NoError(t, err)

	return row
}

func TestTeamRepository_LifecycleAndBulkQueries(t *testing.T) {
	runRepositoryTest(t, func(executor intf.Executor) {
		repo := NewTeamRepository(executor)

		teams := createTeams(t, repo, []model.TeamCreation{
			{Name: "Alpha Team", Description: "First team"},
			{Name: "Bravo Team", Description: "Second team"},
			{Name: "Gamma Team", Description: "Third team"},
		})

		listed, err := repo.List(nil)
		require.NoError(t, err)
		require.Len(t, listed, 3)

		sort.Slice(listed, func(i, j int) bool {
			return listed[i].Name < listed[j].Name
		})
		assert.Equal(t, []string{"Alpha Team", "Bravo Team", "Gamma Team"}, []string{listed[0].Name, listed[1].Name, listed[2].Name})

		selected, err := repo.List(nil, query_option.FilterByID(entity.TeamTable, teams[0].ID))
		require.NoError(t, err)
		require.Len(t, selected, 1)
		assert.Equal(t, teams[0].ID, selected[0].ID)

		err = repo.Update(nil, model.TeamUpdate{
			ID:          teams[1].ID,
			Name:        "Bravo Team Updated",
			Description: "Updated description",
		})
		require.NoError(t, err)

		updated := fetchTeamRow(t, executor, teams[1].ID)
		assert.Equal(t, "Bravo Team Updated", updated.Name)
		assert.Equal(t, "Updated description", updated.Description)

		err = repo.ReserveAvatar(nil, teams[0].ID, "oss://teams/alpha.png")
		require.NoError(t, err)

		reserved := fetchTeamRow(t, executor, teams[0].ID)
		assert.Equal(t, "oss://teams/alpha.png", reserved.AvatarOSSKey)
		assert.False(t, reserved.IsAvatarUploaded)

		err = repo.ConfirmAvatarUploaded(nil, teams[0].ID)
		require.NoError(t, err)

		confirmed := fetchTeamRow(t, executor, teams[0].ID)
		assert.True(t, confirmed.IsAvatarUploaded)

		err = repo.Delete(nil, teams[2].ID)
		require.NoError(t, err)

		remaining, err := repo.List(nil)
		require.NoError(t, err)
		require.Len(t, remaining, 2)

		deletedRows, err := repo.List(nil, query_option.FilterByID(entity.TeamTable, teams[2].ID))
		require.NoError(t, err)
		assert.Empty(t, deletedRows)
	})
}
