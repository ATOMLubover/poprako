package query

import "poprako-s/internal/domain/model/enum"

// `ListMemberOpt` defines filters and pagination for member listing.
type ListMemberOpt struct {
	// `UserId` and `TeamId` are optional filters. If both are provided,
	// the query will return the member that matches both criteria.
	UserId *string
	TeamId *string
	// `UserNicknameKeyword` applies fuzzy search on `user_nickname`.
	UserNicknameKeyword *string

	Pagi PagiOpt

	Includes []enum.MemberIncl
}
