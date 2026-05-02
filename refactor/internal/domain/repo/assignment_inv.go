package repo_iface

import (
	"poprako-s/internal/domain/model/aggr"
	"poprako-s/internal/domain/model/query"
)

// `AssignmentInvRepo` defines persistence contract for assignment invitation.
type AssignmentInvRepo interface {
	// `GetById` returns one invitation by id.
	GetById(id string) (*aggr.AssignmentInv, RepoErr)

	// `List` returns chapter invitations by list options.
	List(opt query.ListAssignmentInvOpt) ([]aggr.AssignmentInv, RepoErr)

	// `ListPendingByInviteeQid` returns pending invitations of one invitee qid.
	ListPendingByInviteeQid(qid string) ([]aggr.AssignmentInv, RepoErr)

	// `Create` inserts one invitation.
	Create(cre *aggr.AssignmentInvCre) (*aggr.AssignmentInv, RepoErr)

	// `Delete` executes hard delete on one invitation.
	Delete(id string) RepoErr

	// `MarkCompleted` marks one invitation as completed.
	MarkCompleted(id string) RepoErr
}
