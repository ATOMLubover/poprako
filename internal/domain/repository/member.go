package repository

import "labelplus-next-web-be/internal/domain/model"

type MemberRepository interface {
	Transactor
	List(executor Executor, options ...QueryOption) ([]model.MemberProfile, error)
	ListWithUserInfo(executor Executor, options ...QueryOption) ([]model.MemberProfile, error)
	Exist(executor Executor, query_option ...QueryOption) (bool, error)
	GetByID(executor Executor, memberID string) (*model.MemberProfile, error)
	Create(executor Executor, creation *model.MemberCreation) (string, error) // 如果成功，返回主键 ID
	Update(executor Executor, update *model.MemberUpdate) error
	DeleteByID(executor Executor, memberID string) error
}
