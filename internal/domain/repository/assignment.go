package repository

import "labelplus-next-web-be/internal/domain/model"

type AssignmentRepository interface {
	Transactor
	Get(executor Executor, options ...QueryOption) (model.AssignmentInfo, error)
	LockByComicID(executor Executor, comicID string) error
}
