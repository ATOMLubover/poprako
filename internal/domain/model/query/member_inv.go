package query

import "poprako-s/internal/domain/model/enum"

// `ListMemberInvOpt` defines filters and pagination for member invitation listing.
type ListMemberInvOpt struct {
	// `TeamId` is required, as we only list invitations of a team.
	TeamId string
	// `Pending` is optional, if not specified,
	// we list all invitations regardless of their pending status.
	Pending *bool

	// `Includes` controls relation preloads supported by the repo.
	Includes []enum.MemberInvIncl

	Pagi PagiOpt
}
