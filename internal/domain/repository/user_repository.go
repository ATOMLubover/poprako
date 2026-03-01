package repository

import "labelplus-next-web-be/internal/domain/model"

type UserRepository interface {
	GetInfoByID(executor Executor, userID string) (*model.UserInfo, error)
	GetCredentialsByQQ(executor Executor, qq string) (*model.UserCredentials, error)
	ExistsByQQ(executor Executor, qq string) (bool, error)
	Create(executor Executor, registration *model.UserRegistration) (string, error) // 如果成功，返回主键 ID
	List(executor Executor, options ...QueryOption) ([]model.UserInfo, error)
	DeleteByID(executor Executor, userID string) error // 根据 ID 软删除用户
}
