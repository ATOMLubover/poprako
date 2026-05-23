package repo_iface

import (
	"poprako-s/internal/domain/model/aggr"
	"poprako-s/internal/domain/model/enum"
	"poprako-s/internal/domain/model/query"
)

// `CommentRepo` defines persistence contract for team board comments.
type CommentRepo interface {
	// `List` returns team comments by query options.
	// Relation loading is controlled by typed `CommentIncl` variadic arguments.
	List(opt query.ListCommentOpt, inc ...enum.CommentIncl) ([]aggr.Comment, RepoErr)

	// `Create` inserts one team board comment.
	Create(cre *aggr.CommentCre) (*aggr.Comment, RepoErr)
}
