package repository

import "labelplus-next-web-be/internal/domain/model"

type TeamRepository interface {
	Transactor
	List(executor Executor, options ...QueryOption) ([]model.TeamInfo, error)
	Create(executor Executor, creation model.TeamCreation) (string, error) // 如果成功，返回主键 ID
	Update(executor Executor, update model.TeamUpdate) error
	Delete(executor Executor, id string) error
}
