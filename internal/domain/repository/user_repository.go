package repository

import "labelplus-next-web-be/internal/domain/model"

type UserRepository interface {
	GetInfoByID(executor Executor, userID string) (*model.UserInfo, error)
}
