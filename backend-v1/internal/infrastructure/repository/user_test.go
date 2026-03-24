package repository

import (
	"labelplus-next-web-be/internal/domain/model"
	"labelplus-next-web-be/internal/infrastructure/repository/entity"
	"labelplus-next-web-be/internal/infrastructure/repository/query_option"
	"sort"
	"testing"

	intf "labelplus-next-web-be/internal/domain/repository"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type createdUser struct {
	ID   string
	Name string
	QQ   string
}

func createUsers(t *testing.T, repo intf.UserRepository, creations []model.UserCreation) []createdUser {
	t.Helper()

	users := make([]createdUser, 0, len(creations))
	for _, creation := range creations {
		id, err := repo.Create(nil, creation)
		require.NoError(t, err)
		users = append(users, createdUser{
			ID:   id,
			Name: creation.Name,
			QQ:   creation.QQ,
		})
	}

	return users
}

func fetchUserRow(t *testing.T, executor intf.Executor, id string) entity.UserInfoRow {
	t.Helper()

	var row entity.UserInfoRow
	err := executor.Table(entity.UserTable).Where("id = ?", id).First(&row).Error
	require.NoError(t, err)

	return row
}

func fetchStatsRow(t *testing.T, executor intf.Executor, userID string) entity.UserStatsRow {
	t.Helper()

	var row entity.UserStatsRow
	err := executor.Table(entity.UserStatsTable).Where("user_id = ?", userID).First(&row).Error
	require.NoError(t, err)

	return row
}

func TestUserRepository_UserLifecycleAndBulkQueries(t *testing.T) {
	runRepositoryTest(t, func(executor intf.Executor) {
		repo := NewUserRepository(executor)

		users := createUsers(t, repo, []model.UserCreation{
			{Name: "Alice Wonder", QQ: "200001", PasswordHash: "hash-alice"},
			{Name: "Bob Stone", QQ: "200002", PasswordHash: "hash-bob"},
			{Name: "Charlie Wave", QQ: "200003", PasswordHash: "hash-charlie"},
		})

		for _, user := range users {
			assert.NotEmpty(t, user.ID)
		}

		allUsers, err := repo.List(nil)
		require.NoError(t, err)
		require.Len(t, allUsers, 3)

		sort.Slice(allUsers, func(i, j int) bool {
			return allUsers[i].QQ < allUsers[j].QQ
		})
		assert.Equal(t, []string{"200001", "200002", "200003"}, []string{allUsers[0].QQ, allUsers[1].QQ, allUsers[2].QQ})

		alice, err := repo.Get(nil, query_option.FilterByID(entity.UserTable, users[0].ID))
		require.NoError(t, err)
		assert.Equal(t, "Alice Wonder", alice.Name)
		assert.Equal(t, "200001", alice.QQ)

		bobByQQ, err := repo.Get(nil, query_option.UserQuery().FilterByQQ("200002"))
		require.NoError(t, err)
		assert.Equal(t, users[1].ID, bobByQQ.ID)

		fuzzyUsers, err := repo.List(nil, query_option.UserQuery().FilterByFuzzyName("Wave"))
		require.NoError(t, err)
		require.Len(t, fuzzyUsers, 1)
		assert.Equal(t, users[2].ID, fuzzyUsers[0].ID)

		credentials, err := repo.GetCredentials(nil, query_option.UserQuery().FilterByQQ("200001"))
		require.NoError(t, err)
		assert.Equal(t, users[0].ID, credentials.UserID)
		assert.Equal(t, "hash-alice", credentials.PasswordHash)

		err = repo.Update(nil, model.UserUpdate{
			ID:           users[1].ID,
			Name:         "Bob Updated",
			QQ:           "299999",
			PasswordHash: "hash-bob-updated",
		})
		require.NoError(t, err)

		updatedBob, err := repo.Get(nil, query_option.FilterByID(entity.UserTable, users[1].ID))
		require.NoError(t, err)
		assert.Equal(t, "Bob Updated", updatedBob.Name)
		assert.Equal(t, "299999", updatedBob.QQ)

		err = repo.ReserveAvatar(nil, users[0].ID, "oss://avatars/alice.png")
		require.NoError(t, err)

		reservedAlice := fetchUserRow(t, executor, users[0].ID)
		assert.Equal(t, "oss://avatars/alice.png", reservedAlice.AvatarOSSKey)
		assert.False(t, reservedAlice.IsAvatarUploaded)

		err = repo.ConfirmAvatarUploaded(nil, users[0].ID)
		require.NoError(t, err)

		confirmedAlice, err := repo.Get(nil, query_option.FilterByID(entity.UserTable, users[0].ID))
		require.NoError(t, err)
		assert.True(t, confirmedAlice.IsAvatarUploaded)

		err = repo.Delete(nil, users[2].ID)
		require.NoError(t, err)

		_, err = repo.Get(nil, query_option.FilterByID(entity.UserTable, users[2].ID))
		assert.ErrorIs(t, err, ErrRecordNotFound)

		remainingUsers, err := repo.List(nil)
		require.NoError(t, err)
		require.Len(t, remainingUsers, 2)

		_, err = repo.GetCredentials(nil, query_option.UserQuery().FilterByQQ("does-not-exist"))
		assert.ErrorIs(t, err, ErrRecordNotFound)
	})
}

func TestUserRepository_StatsLifecycle(t *testing.T) {
	runRepositoryTest(t, func(executor intf.Executor) {
		repo := NewUserRepository(executor)

		users := createUsers(t, repo, []model.UserCreation{
			{Name: "Stats One", QQ: "300001", PasswordHash: "stats-one"},
			{Name: "Stats Two", QQ: "300002", PasswordHash: "stats-two"},
		})

		err := repo.CreateStats(nil, model.UserStatsCreation{
			UserID:                  users[0].ID,
			TotalAssignmentCount:    10,
			ActiveAssignmentCount:   4,
			FinishedAssignmentCount: 6,
		})
		require.NoError(t, err)

		err = repo.CreateStats(nil, model.UserStatsCreation{
			UserID:                  users[1].ID,
			TotalAssignmentCount:    2,
			ActiveAssignmentCount:   2,
			FinishedAssignmentCount: 0,
		})
		require.NoError(t, err)

		statsOne, err := repo.GetStats(nil, query_option.UserStatsQuery().FilterByUserID(users[0].ID))
		require.NoError(t, err)
		assert.Equal(t, 10, statsOne.TotalAssignmentCount)
		assert.Equal(t, 4, statsOne.ActiveAssignmentCount)
		assert.Equal(t, 6, statsOne.FinishedAssignmentCount)

		err = repo.IncrementStats(nil, model.UserStatsDelta{
			UserID:                       users[0].ID,
			TotalAssignmentCountDelta:    3,
			ActiveAssignmentCountDelta:   -1,
			FinishedAssignmentCountDelta: 4,
		})
		require.NoError(t, err)

		incremented := fetchStatsRow(t, executor, users[0].ID)
		assert.Equal(t, 13, incremented.TotalAssignmentCount)
		assert.Equal(t, 3, incremented.ActiveAssignmentCount)
		assert.Equal(t, 10, incremented.FinishedAssignmentCount)

		statsTwo, err := repo.GetStats(nil, query_option.UserStatsQuery().FilterByUserID(users[1].ID))
		require.NoError(t, err)
		assert.Equal(t, 2, statsTwo.TotalAssignmentCount)
		assert.Equal(t, 2, statsTwo.ActiveAssignmentCount)
		assert.Equal(t, 0, statsTwo.FinishedAssignmentCount)

		_, err = repo.GetStats(nil, query_option.UserStatsQuery().FilterByUserID("missing-user"))
		assert.ErrorIs(t, err, ErrRecordNotFound)
	})
}
