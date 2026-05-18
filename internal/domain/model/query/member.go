package query

import (
	"poprako-s/internal/domain/model/aggr"
	"poprako-s/internal/domain/model/enum"
)

// `ListMemberOpt` defines filters and pagination for member listing.
type ListMemberOpt struct {
	// `UserId` and `TeamId` are optional filters. If both are provided,
	// the query will return the member that matches both criteria.
	UserId *string
	TeamId *string
	// `UserNicknameKeyword` applies fuzzy search on `user_nickname`.
	UserNicknameKeyword *string
	// `Role` filters members that hold the given single role.
	// Must be a single bit value (exactly one `enum.Role` constant).
	Role *aggr.RoleMask

	Pagi PagiOpt

	Includes []enum.MemberIncl
}
