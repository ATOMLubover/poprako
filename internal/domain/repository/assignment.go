package repository

import "labelplus-next-web-be/internal/domain/model"

type AssignmentRepository interface {
	Transactor
	Get(executor Executor, options ...QueryOption) (model.AssignmentInfo, error)
	Exist(executor Executor, options ...QueryOption) (bool, error)
	ListWithUserInfo(executor Executor, options ...QueryOption) ([]model.AssignmentInfo, error)
	ListWithChapterInfo(executor Executor, options ...QueryOption) ([]model.AssignmentInfo, error)
	Create(executor Executor, creation model.AssignmentCreation) (string, error)
	Update(executor Executor, update model.AssignmentUpdate) error
	Delete(executor Executor, assignmentID string) error
	LockByComicID(executor Executor, comicID string) error
}
