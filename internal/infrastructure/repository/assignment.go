package repository

import (
	"labelplus-next-web-be/internal/domain/model"
	intf "labelplus-next-web-be/internal/domain/repository"
	"labelplus-next-web-be/internal/infrastructure/repository/entity"
)

type assignmentRepository struct {
	executor intf.Executor
}

func NewAssignmentRepository(executor intf.Executor) intf.AssignmentRepository {
	return &assignmentRepository{executor: executor}
}

func (r *assignmentRepository) withTransaction(executor intf.Executor) intf.Executor {
	if executor != nil {
		return executor
	}
	return r.executor
}

func (r *assignmentRepository) BeginTransaction() intf.Executor {
	return r.executor.Begin()
}

func (r *assignmentRepository) Get(executor intf.Executor, options ...intf.QueryOption) (model.AssignmentInfo, error) {
	executor = r.withTransaction(executor)
	executor = executor.Table(entity.AssignmentTable)

	for _, opt := range options {
		executor = opt(executor)
	}

	var row entity.AssignmentInfoRow
	if err := executor.First(&row).Error; err != nil {
		return model.AssignmentInfo{}, err
	}

	return entity.ToAssignmentInfo(row), nil
}

func (r *assignmentRepository) LockByComicID(executor intf.Executor, comicID string) error {
	executor = r.withTransaction(executor)

	var lockedIDs []string

	return executor.
		Table(entity.AssignmentTable).
		Where("chapter_id IN (SELECT id FROM chapter_table WHERE comic_id = ?)", comicID).
		Pluck("id", &lockedIDs).Error
}
