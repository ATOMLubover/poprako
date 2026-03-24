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

type createdAssignment struct {
	ID string
}

func createAssignments(t *testing.T, repo intf.AssignmentRepository, creations []model.AssignmentCreation) []createdAssignment {
	t.Helper()

	assignments := make([]createdAssignment, 0, len(creations))
	for _, creation := range creations {
		id, err := repo.Create(nil, creation)
		require.NoError(t, err)
		assignments = append(assignments, createdAssignment{ID: id})
	}

	return assignments
}

func TestAssignmentRepository_LifecycleAndIncludes(t *testing.T) {
	runRepositoryTest(t, func(executor intf.Executor) {
		userRepo := NewUserRepository(executor)
		teamRepo := NewTeamRepository(executor)
		worksetRepo := NewWorksetRepository(executor)
		comicRepo := NewComicRepository(executor)
		assignmentRepo := NewAssignmentRepository(executor)

		users := createUsers(t, userRepo, []model.UserCreation{
			model.NewUserRegistration("Assignment User One", "540001", "pw-1"),
			model.NewUserRegistration("Assignment User Two", "540002", "pw-2"),
			model.NewUserRegistration("Assignment Creator", "540003", "pw-3"),
		})
		teams := createTeams(t, teamRepo, []model.TeamCreation{
			*model.NewTeamCreation("Assignment Team", "Repository assignment tests"),
		})
		worksets := createWorksets(t, worksetRepo, []model.WorksetCreation{
			*model.NewWorksetCreation(teams[0].ID, 0, "Assignment Shelf", "assignment workset"),
		})
		comics := createComics(t, comicRepo, []model.ComicCreation{
			*model.NewComicCreation(worksets[0].ID, 0, "Assignment Comic", "Assignment Author", "assignment desc", users[2].ID),
		})

		firstChapterID := createChapterRow(t, executor, comics[0].ID, 0, "Assignment Chapter A", users[2].ID)
		secondChapterID := createChapterRow(t, executor, comics[0].ID, 1, "Assignment Chapter B", users[2].ID)

		transactionExecutor := assignmentRepo.BeginTransaction()
		require.NoError(t, transactionExecutor.Error)
		err := assignmentRepo.LockByComicID(transactionExecutor, comics[0].ID)
		require.NoError(t, err)
		require.NoError(t, transactionExecutor.Rollback().Error)

		assignments := createAssignments(t, assignmentRepo, []model.AssignmentCreation{
			model.NewAssignmentCreation(firstChapterID, users[0].ID, model.MaskRoles([]model.RoleFlag{model.RoleTranslator, model.RoleReviewer})),
			model.NewAssignmentCreation(firstChapterID, users[1].ID, model.MaskRoles([]model.RoleFlag{model.RoleProofreader})),
			model.NewAssignmentCreation(secondChapterID, users[0].ID, model.MaskRoles([]model.RoleFlag{model.RoleRawProvider})),
		})

		exists, err := assignmentRepo.Exist(nil,
			query_option.AssignmentQuery().FilterByChapterID(firstChapterID),
			query_option.AssignmentQuery().FilterByUserID(users[0].ID),
		)
		require.NoError(t, err)
		assert.True(t, exists)

		notExists, err := assignmentRepo.Exist(nil,
			query_option.AssignmentQuery().FilterByChapterID(secondChapterID),
			query_option.AssignmentQuery().FilterByUserID(users[1].ID),
		)
		require.NoError(t, err)
		assert.False(t, notExists)

		plainAssignment, err := assignmentRepo.Get(nil, query_option.FilterByID(entity.AssignmentTable, assignments[0].ID))
		require.NoError(t, err)
		assert.Equal(t, firstChapterID, plainAssignment.ChapterID)
		assert.Equal(t, users[0].ID, plainAssignment.UserID)
		assert.NotNil(t, plainAssignment.AssignedTranslatorAt)
		assert.NotNil(t, plainAssignment.AssignedReviewerAt)
		assert.Nil(t, plainAssignment.User)
		assert.Nil(t, plainAssignment.Chapter)

		includedAssignments, err := assignmentRepo.List(nil,
			query_option.AssignmentQuery().FilterByChapterID(firstChapterID),
			query_option.AssignmentQuery().IncludeRelationInfo(true, true, true, true),
			query_option.UpdatedAtDesc(entity.AssignmentTable),
		)
		require.NoError(t, err)
		require.Len(t, includedAssignments, 2)

		includedAssignmentByID := make(map[string]model.AssignmentInfo, len(includedAssignments))
		for _, includedAssignment := range includedAssignments {
			includedAssignmentByID[includedAssignment.ID] = includedAssignment
		}

		require.NotNil(t, includedAssignmentByID[assignments[0].ID].User)
		assert.Equal(t, "Assignment User One", includedAssignmentByID[assignments[0].ID].User.Name)
		require.NotNil(t, includedAssignmentByID[assignments[0].ID].Chapter)
		assert.Equal(t, "Assignment Chapter A", includedAssignmentByID[assignments[0].ID].Chapter.Subtitle)
		require.NotNil(t, includedAssignmentByID[assignments[0].ID].Chapter.Comic)
		assert.Equal(t, "Assignment Comic", includedAssignmentByID[assignments[0].ID].Chapter.Comic.Title)
		require.NotNil(t, includedAssignmentByID[assignments[0].ID].Chapter.Creator)
		assert.Equal(t, "Assignment Creator", includedAssignmentByID[assignments[0].ID].Chapter.Creator.Name)

		userOnlyAssignments, err := assignmentRepo.List(nil,
			query_option.AssignmentQuery().FilterByUserID(users[0].ID),
			query_option.AssignmentQuery().IncludeUserInfo(),
		)
		require.NoError(t, err)
		require.Len(t, userOnlyAssignments, 2)
		for _, userOnlyAssignment := range userOnlyAssignments {
			require.NotNil(t, userOnlyAssignment.User)
			assert.Nil(t, userOnlyAssignment.Chapter)
		}

		currentAssignment, err := assignmentRepo.Get(nil, query_option.FilterByID(entity.AssignmentTable, assignments[0].ID))
		require.NoError(t, err)
		originalTranslatorAt := currentAssignment.AssignedTranslatorAt

		assignmentUpdate := model.NewAssignmentUpdate(
			assignments[0].ID,
			currentAssignment,
			model.MaskRoles([]model.RoleFlag{model.RoleTranslator, model.RolePublisher}),
		)
		err = assignmentRepo.Update(nil, assignmentUpdate)
		require.NoError(t, err)

		updatedAssignment, err := assignmentRepo.Get(nil, query_option.FilterByID(entity.AssignmentTable, assignments[0].ID))
		require.NoError(t, err)
		require.NotNil(t, updatedAssignment.AssignedTranslatorAt)
		require.NotNil(t, originalTranslatorAt)
		assert.True(t, updatedAssignment.AssignedTranslatorAt.Equal(*originalTranslatorAt))
		assert.Nil(t, updatedAssignment.AssignedReviewerAt)
		assert.NotNil(t, updatedAssignment.AssignedPublisherAt)

		transactionExecutor = assignmentRepo.BeginTransaction()
		require.NoError(t, transactionExecutor.Error)
		err = assignmentRepo.LockByComicID(transactionExecutor, comics[0].ID)
		require.NoError(t, err)
		require.NoError(t, transactionExecutor.Rollback().Error)

		err = assignmentRepo.Delete(nil, assignments[1].ID)
		require.NoError(t, err)

		_, err = assignmentRepo.Get(nil, query_option.FilterByID(entity.AssignmentTable, assignments[1].ID))
		assert.ErrorIs(t, err, ErrRecordNotFound)

		exists, err = assignmentRepo.Exist(nil,
			query_option.AssignmentQuery().FilterByChapterID(firstChapterID),
			query_option.AssignmentQuery().FilterByUserID(users[1].ID),
		)
		require.NoError(t, err)
		assert.False(t, exists)
	})
}
