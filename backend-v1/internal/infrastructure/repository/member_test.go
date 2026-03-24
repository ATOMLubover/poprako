package repository

import (
	"labelplus-next-web-be/internal/domain/model"
	"labelplus-next-web-be/internal/infrastructure/repository/entity"
	"labelplus-next-web-be/internal/infrastructure/repository/query_option"
	"testing"
	"time"

	intf "labelplus-next-web-be/internal/domain/repository"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMemberRepository_LifecycleAndIncludes(t *testing.T) {
	runRepositoryTest(t, func(executor intf.Executor) {
		userRepo := NewUserRepository(executor)
		teamRepo := NewTeamRepository(executor)
		memberRepo := NewMemberRepository(executor)

		users := createUsers(t, userRepo, []model.UserCreation{
			{Name: "Member User One", QQ: "410001", PasswordHash: "pw-1"},
			{Name: "Member User Two", QQ: "410002", PasswordHash: "pw-2"},
		})
		teams := createTeams(t, teamRepo, []model.TeamCreation{{Name: "Member Team", Description: "Team for members"}})

		memberIDs := make([]string, 0, 2)
		for _, creation := range []model.MemberCreation{
			{UserID: users[0].ID, TeamID: teams[0].ID, ToBeTranslator: true, ToBeReviewer: true},
			{UserID: users[1].ID, TeamID: teams[0].ID, ToBeProofreader: true},
		} {
			id, err := memberRepo.Create(nil, creation)
			require.NoError(t, err)
			memberIDs = append(memberIDs, id)
		}

		listed, err := memberRepo.List(nil, query_option.MemberQuery().FilterByTeamID(teams[0].ID))
		require.NoError(t, err)
		require.Len(t, listed, 2)

		exists, err := memberRepo.Exist(nil,
			query_option.MemberQuery().FilterByTeamID(teams[0].ID),
			query_option.MemberQuery().FilterByUserID(users[0].ID),
		)
		require.NoError(t, err)
		assert.True(t, exists)

		plain, err := memberRepo.Get(nil, query_option.FilterByID(entity.MemberTable, memberIDs[0]))
		require.NoError(t, err)
		assert.Equal(t, users[0].ID, plain.UserID)
		assert.Equal(t, teams[0].ID, plain.TeamID)
		assert.NotNil(t, plain.AssignedTranslatorAt)
		assert.NotNil(t, plain.AssignedReviewerAt)
		assert.Nil(t, plain.User)
		assert.Nil(t, plain.Team)

		profile, err := memberRepo.GetProfile(nil, query_option.FilterByID(entity.MemberTable, memberIDs[0]))
		require.NoError(t, err)
		assert.Equal(t, users[0].ID, profile.UserID)
		assert.Equal(t, teams[0].ID, profile.TeamID)

		withUser, err := memberRepo.List(nil,
			query_option.MemberQuery().FilterByTeamID(teams[0].ID),
			query_option.MemberQuery().FilterByUserID(users[0].ID),
			query_option.MemberQuery().IncludeUserInfo(),
		)
		require.NoError(t, err)
		require.Len(t, withUser, 1)
		require.NotNil(t, withUser[0].User)
		assert.Equal(t, "Member User One", withUser[0].User.Name)

		withTeam, err := memberRepo.List(nil,
			query_option.MemberQuery().FilterByUserID(users[0].ID),
			query_option.MemberQuery().IncludeTeamInfo(),
		)
		require.NoError(t, err)
		require.Len(t, withTeam, 1)
		require.NotNil(t, withTeam[0].Team)
		assert.Equal(t, "Member Team", withTeam[0].Team.Name)

		now := time.Now().UTC().Truncate(time.Microsecond)
		err = memberRepo.Update(nil, model.MemberUpdate{
			ID:                    memberIDs[1],
			AssignedRawProviderAt: &now,
			AssignedTranslatorAt:  &now,
			AssignedProofreaderAt: nil,
			AssignedTypesetterAt:  nil,
			AssignedReviewerAt:    nil,
			AssignedPublisherAt:   nil,
			AssignedAdminAt:       &now,
		})
		require.NoError(t, err)

		updated, err := memberRepo.Get(nil, query_option.FilterByID(entity.MemberTable, memberIDs[1]))
		require.NoError(t, err)
		assert.NotNil(t, updated.AssignedRawProviderAt)
		assert.NotNil(t, updated.AssignedTranslatorAt)
		assert.Nil(t, updated.AssignedProofreaderAt)
		assert.NotNil(t, updated.AssignedAdminAt)

		qqMatch, err := memberRepo.Exist(nil,
			query_option.MemberQuery().JoinUser(),
			query_option.MemberQuery().FilterByTeamID(teams[0].ID),
			query_option.MemberQuery().FilterOnUserQQ("410002"),
		)
		require.NoError(t, err)
		assert.True(t, qqMatch)

		err = memberRepo.Delete(nil, memberIDs[0])
		require.NoError(t, err)

		afterDelete, err := memberRepo.List(nil, query_option.MemberQuery().FilterByTeamID(teams[0].ID))
		require.NoError(t, err)
		require.Len(t, afterDelete, 1)
		assert.Equal(t, memberIDs[1], afterDelete[0].ID)

		exists, err = memberRepo.Exist(nil,
			query_option.MemberQuery().FilterByTeamID(teams[0].ID),
			query_option.MemberQuery().FilterByUserID(users[0].ID),
		)
		require.NoError(t, err)
		assert.False(t, exists)
	})
}
