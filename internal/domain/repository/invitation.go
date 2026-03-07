package repository

import "labelplus-next-web-be/internal/domain/model"

type InvitationRepository interface {
	Transactor
	List(executor Executor, options ...QueryOption) ([]model.InvitationInfo, error)
	Get(executor Executor, options ...QueryOption) (model.InvitationInfo, error)
	Create(executor Executor, creation model.InvitationCreation) (string, error) // 如果成功，返回主键 ID
	Update(executor Executor, update model.InvitationUpdate) error
	Invalidate(executor Executor, id string) error
	Delete(executor Executor, id string) error
}
