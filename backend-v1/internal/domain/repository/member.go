package repository

import "labelplus-next-web-be/internal/domain/model"

type MemberRepository interface {
	Transactor
	List(executor Executor, options ...QueryOption) ([]model.MemberInfo, error)
	Get(executor Executor, option ...QueryOption) (model.MemberInfo, error)
	GetProfile(executor Executor, options ...QueryOption) (model.MemberInfo, error)
	Exist(executor Executor, options ...QueryOption) (bool, error)
	Create(executor Executor, creation model.MemberCreation) (string, error) // 如果成功，返回主键 ID
	Update(executor Executor, update model.MemberUpdate) error
	Delete(executor Executor, id string) error
}
