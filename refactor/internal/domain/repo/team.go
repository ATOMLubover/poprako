package repo_iface

import (
	"poprako-s/internal/domain/model/aggr"
	"poprako-s/internal/domain/model/query"
)

// `TeamRepo` defines persistence contract for team aggregate.
type TeamRepo interface {
	// `GetById` retrieves a team by its unique identifier.
	// An error is returned if the team does not exist.
	GetById(id string) (*aggr.Team, RepoErr)

	// `List` returns teams by list filters and pagination.
	List(opt *query.ListTeamOpt) ([]*aggr.Team, RepoErr)

	// `Create` inserts one team row and returns created aggregate.
	Create(cre *aggr.TeamCre) (*aggr.Team, RepoErr)

	// `Update` applies put-style team update to one row.
	Update(upd *aggr.TeamUpd) RepoErr

	// `Delete` executes hard delete on one team row.
	Delete(id string) RepoErr

	// `PrefillAvatarKey` writes reserved avatar key before client upload.
	PrefillAvatarKey(id string, key string) RepoErr

	// `MarkAvatarUploaded` marks team avatar upload as completed.
	MarkAvatarUploaded(id string) RepoErr

	// `IncrementWorksetNextIndex` allocates one workset index from team-scoped sequence.
	// It returns the allocated index value for immediate insert use.
	IncrementWorksetNextIndex(id string) (int, RepoErr)
}
