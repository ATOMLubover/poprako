package repo_iface

import (
	"poprako-s/internal/domain/model/aggr"
	"poprako-s/internal/domain/model/query"
)

type MemberInvRepo interface {
	GetPendingByInviteeQid(qid string) (*aggr.MemberInv, RepoErr)
	List(opt query.ListMemberInvOpt) ([]aggr.MemberInv, RepoErr)

	Create(cre *aggr.MemberInvCre) (*aggr.MemberInv, RepoErr)

	// `Delete` executes a **hard** delete on given member invitation id.
	Delete(id string) RepoErr

	// `MarkCmpl` marks the member invitation as completed without removing it from the database.
	MarkCmpl(id string) RepoErr
}
