package repo_iface

import (
	"poprako-s/internal/domain/model/aggr"
	"poprako-s/internal/domain/model/enum"
	"poprako-s/internal/domain/model/query"
)

type MemberRepo interface {
	GetById(id string, inc ...enum.MemberIncl) (*aggr.Member, RepoErr)
	List(opt *query.ListMemberOpt, inc ...enum.MemberIncl) ([]*aggr.Member, RepoErr)
	ExistByUserTeamId(userId string, teamId string) (bool, RepoErr)

	Create(cre *aggr.MemberCre) (*aggr.Member, RepoErr)

	UpdateRoles(upd *aggr.MemberRoleUpd) RepoErr

	// `Delete` executes a **hard** delete on given member id.
	Delete(id string) RepoErr
}
