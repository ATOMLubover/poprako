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

func TestInvitationRepository_LifecycleAndFilters(t *testing.T) {
	runRepositoryTest(t, func(executor intf.Executor) {
		userRepo := NewUserRepository(executor)
		teamRepo := NewTeamRepository(executor)
		invitationRepo := NewInvitationRepository(executor)

		users := createUsers(t, userRepo, []model.UserCreation{{Name: "Invitor User", QQ: "510001", PasswordHash: "pw-invitor"}})
		teams := createTeams(t, teamRepo, []model.TeamCreation{{Name: "Invite Team", Description: "Invitation target"}})

		invitationIDs := make([]string, 0, 2)
		for _, creation := range []model.InvitationCreation{
			{InvitorID: users[0].ID, TargetTeamID: teams[0].ID, InviteeQQ: "600001", InvitationCode: "INV-A", ToBeTranslator: true, ToBeReviewer: true},
			{InvitorID: users[0].ID, TargetTeamID: teams[0].ID, InviteeQQ: "600002", InvitationCode: "INV-B", ToBeProofreader: true},
		} {
			id, err := invitationRepo.Create(nil, creation)
			require.NoError(t, err)
			invitationIDs = append(invitationIDs, id)
		}

		listed, err := invitationRepo.List(nil, query_option.InvitationQuery().FilterByTeamID(teams[0].ID))
		require.NoError(t, err)
		require.Len(t, listed, 2)

		selected, err := invitationRepo.Get(nil,
			query_option.InvitationQuery().FilterByInviteeQQ("600001"),
			query_option.InvitationQuery().FilterByCode("INV-A"),
			query_option.InvitationQuery().IncludeInvitorInfo(),
		)
		require.NoError(t, err)
		assert.Equal(t, invitationIDs[0], selected.ID)
		require.NotNil(t, selected.Invitor)
		assert.Equal(t, "Invitor User", selected.Invitor.Name)
		assert.True(t, selected.ToBeTranslator)
		assert.True(t, selected.ToBeReviewer)

		err = invitationRepo.Update(nil, model.InvitationUpdate{
			ID:              invitationIDs[1],
			ToBeRawProvider: true,
			ToBeTranslator:  true,
			ToBeProofreader: false,
			ToBeTypesetter:  false,
			ToBeReviewer:    false,
			ToBePublisher:   true,
			ToBeAdmin:       false,
		})
		require.NoError(t, err)

		updated, err := invitationRepo.Get(nil, query_option.FilterByID(entity.InvitationTable, invitationIDs[1]))
		require.NoError(t, err)
		assert.True(t, updated.ToBeRawProvider)
		assert.True(t, updated.ToBeTranslator)
		assert.False(t, updated.ToBeProofreader)
		assert.True(t, updated.ToBePublisher)

		err = invitationRepo.Invalidate(nil, invitationIDs[0])
		require.NoError(t, err)

		pendingInvitations, err := invitationRepo.List(nil, query_option.InvitationQuery().FilterPending(true))
		require.NoError(t, err)
		require.Len(t, pendingInvitations, 1)
		assert.Equal(t, invitationIDs[1], pendingInvitations[0].ID)

		closedInvitations, err := invitationRepo.List(nil, query_option.InvitationQuery().FilterPending(false))
		require.NoError(t, err)
		require.Len(t, closedInvitations, 1)
		assert.Equal(t, invitationIDs[0], closedInvitations[0].ID)

		err = invitationRepo.Delete(nil, invitationIDs[1])
		require.NoError(t, err)

		_, err = invitationRepo.Get(nil, query_option.FilterByID(entity.InvitationTable, invitationIDs[1]))
		assert.ErrorIs(t, err, ErrRecordNotFound)
	})
}
