package svc

import (
	"poprako-s/internal/domain/model/aggr"
	"poprako-s/pkg/util"
)

// `WorksetSvc` provides stateless domain services for the `Workset` aggregate.
type WorksetSvc struct{}

// `NewWorksetSvc` returns a ready-to-use `WorksetSvc`.
func NewWorksetSvc() WorksetSvc {
	return WorksetSvc{}
}

// `NewWorksetCre` builds a `WorksetCre` aggregate with a generated id.
// The `index` must be pre-computed by the caller from a transactional count
// of active worksets for the given team.
func (WorksetSvc) NewWorksetCre(teamId string, index int, name string, desc *string) *aggr.WorksetCre {
	return &aggr.WorksetCre{
		Id:     util.GenId("workset"),
		TeamId: teamId,
		Index:  index,
		Name:   name,
		Desc:   desc,
	}
}
