package adapter

import (
	"labelplus-next-web-be/internal/domain/model"
	"labelplus-next-web-be/internal/domain/repository"
	"labelplus-next-web-be/internal/domain/service"
	repository_infra "labelplus-next-web-be/internal/repository"
	"labelplus-next-web-be/internal/repository/query_option"
)

func HandleLoadUserInfo(userRepository repository.UserRepository) service.OnLoadUserInfo {
	return func(userID string) (model.UserInfo, error) {
		return userRepository.GetByID(nil, query_option.FilterByID(repository_infra.UserTable, userID))
	}
}

func HandleLoadMemberInfo(memberRepository repository.MemberRepository) service.OnLoadMemberInfo {
	return func(teamID string, userID string) (model.MemberInfo, error) {
		return memberRepository.Get(
			nil,
			query_option.MemberQuery().FilterByUserID(userID),
			query_option.MemberQuery().FilterByTeamID(teamID),
		)
	}
}
