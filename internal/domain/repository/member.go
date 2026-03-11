package repository

import "labelplus-next-web-be/internal/domain/model"

type MemberRepository interface {
	Transactor
	ListProfiles(executor Executor, options ...QueryOption) ([]model.MemberWithUserInfo, error)
	ListWithTeamInfo(executor Executor, options ...QueryOption) ([]model.MemberWithTeamInfo, error)
	Get(executor Executor, option ...QueryOption) (model.MemberInfo, error)
	GetProfile(executor Executor, options ...QueryOption) (model.MemberWithUserInfo, error)
	Exist(executor Executor, options ...QueryOption) (bool, error)
	Create(executor Executor, creation model.MemberCreation) (string, error) // 如果成功，返回主键 ID
	Update(executor Executor, update model.MemberUpdate) error
	Delete(executor Executor, id string) error
}
