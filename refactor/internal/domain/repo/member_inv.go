package repo_iface

import (
	"poprako-s/internal/domain/model/aggr"
	"poprako-s/internal/domain/model/query"
)

type MemberInvRepo interface {
	GetById(id string) (*aggr.MemberInv, RepoErr)

	GetPendingByInviteeQid(qid string) (*aggr.MemberInv, RepoErr)
	List(opt query.ListMemberInvOpt) ([]aggr.MemberInv, RepoErr)

	Create(cre *aggr.MemberInvCre) (*aggr.MemberInv, RepoErr)

	// `Update` applies put-style role update to one member invitation.
	Update(upd *aggr.MemberInvUpd) RepoErr

	// `Delete` executes a **hard** delete on given member invitation id.
	Delete(id string) RepoErr

	// `MarkCompleted` marks the member invitation as completed without removing it from the database.
	MarkCompleted(id string) RepoErr
}
