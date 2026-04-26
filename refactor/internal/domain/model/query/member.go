package query

type ListMemberOpt struct {
	// `UserId` and `TeamId` are optional filters. If both are provided,
	// the query will return the member that matches both criteria.
	UserId *string
	TeamId *string
}
