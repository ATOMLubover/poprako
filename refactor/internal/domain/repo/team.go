package repo_iface

import "poprako-s/internal/domain/model/aggr"

type TeamRepo interface {
	// `GetById` retrieves a team by its unique identifier.
	// An error is returned if the team does not exist.
	GetById(id string) (*aggr.Team, RepoErr)
}
