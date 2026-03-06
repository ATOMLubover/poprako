package repository

import "labelplus-next-web-be/internal/domain/model"

type UserRepository interface {
	Transactor
	List(executor Executor, options ...QueryOption) ([]model.UserInfo, error)
	GetInfoByID(executor Executor, userID string) (*model.UserInfo, error)
	GetCredentialsByQQ(executor Executor, qq string) (*model.UserCredentials, error)
	Create(executor Executor, registration *model.UserRegistration) (string, error) // 如果成功，返回主键 ID
	DeleteByID(executor Executor, userID string) error
}
