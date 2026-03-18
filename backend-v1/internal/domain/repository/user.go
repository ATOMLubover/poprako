package repository

import "labelplus-next-web-be/internal/domain/model"

type UserRepository interface {
	Transactor
	List(executor Executor, options ...QueryOption) ([]model.UserInfo, error)
	Get(executor Executor, options ...QueryOption) (model.UserInfo, error)
	GetStats(executor Executor, options ...QueryOption) (model.UserStats, error)
	CreateStats(executor Executor, creation model.UserStatsCreation) error
	IncrementStats(executor Executor, delta model.UserStatsDelta) error
	GetCredentials(executor Executor, options ...QueryOption) (model.UserCredentials, error)
	Create(executor Executor, creation model.UserCreation) (string, error) // 如果成功，返回主键 ID
	ReserveAvatar(executor Executor, id string, avatarOSSKey string) error
	ConfirmAvatarUploaded(executor Executor, id string) error
	Update(executor Executor, update model.UserUpdate) error
	Delete(executor Executor, id string) error
}
