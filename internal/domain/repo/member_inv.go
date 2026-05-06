package repo_iface

import (
	"poprako-s/internal/domain/model/aggr"
	"poprako-s/internal/domain/model/enum"
	"poprako-s/internal/domain/model/query"
)

// `MemberInvRepo` defines persistence contract for member invitations.
type MemberInvRepo interface {
	// `GetById` returns one invitation by id.
	GetById(id string, inc ...enum.MemberInvIncl) (*aggr.MemberInv, RepoErr)

	// `GetPendingByInviteeQid` returns pending invitation by invitee qid.
	GetPendingByInviteeQid(qid string, inc ...enum.MemberInvIncl) (*aggr.MemberInv, RepoErr)

	// `List` returns invitations by query options.
	// Relation loading is controlled by typed `MemberInvIncl` variadic arguments.
	List(opt query.ListMemberInvOpt, inc ...enum.MemberInvIncl) ([]aggr.MemberInv, RepoErr)

	// `Create` inserts one member invitation.
	Create(cre *aggr.MemberInvCre) (*aggr.MemberInv, RepoErr)

	// `Update` applies put-style role update to one member invitation.
	Update(upd *aggr.MemberInvUpd) RepoErr

	// `Delete` executes a **hard** delete on given member invitation id.
	Delete(id string) RepoErr

	// `MarkCompleted` marks the member invitation as completed without removing it from the database.
	MarkCompleted(id string) RepoErr
}
