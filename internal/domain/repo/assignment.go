package repo_iface

import (
	"poprako-s/internal/domain/model/aggr"
	"poprako-s/internal/domain/model/query"
)

// `AssignmentRepo` defines persistence contract for assignment.
type AssignmentRepo interface {
	// `GetById` returns one assignment by id.
	GetById(id string) (*aggr.Assignment, RepoErr)

	// `GetByChapterUserId` returns one assignment by chapter and user.
	GetByChapterUserId(chapterId string, userId string) (*aggr.Assignment, RepoErr)

	// `List` returns assignment list by query options.
	List(opt *query.ListAssignmentOpt) ([]*aggr.Assignment, RepoErr)

	// `Create` inserts one assignment.
	Create(cre *aggr.AssignmentCre) (*aggr.Assignment, RepoErr)

	// `Put` overwrites all role timestamp fields.
	Put(put *aggr.AssignmentPut) RepoErr

	// `Delete` executes hard delete on one assignment by id.
	Delete(id string) RepoErr

	// `DeleteByChapterUserId` executes hard delete by chapter and user.
	DeleteByChapterUserId(chapterId string, userId string) RepoErr

	// `DeleteByChapterId` executes hard delete by chapter.
	DeleteByChapterId(chapterId string) RepoErr
}
