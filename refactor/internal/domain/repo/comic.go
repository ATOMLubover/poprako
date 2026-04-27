package repo_iface

import (
	"poprako-s/internal/domain/model/aggr"
	"poprako-s/internal/domain/model/enum"
	"poprako-s/internal/domain/model/query"
)

// `ComicRepo` defines the persistence contract for `Comic` aggregate
// Repo methods return active records only deleted rows are filtered out
// List style methods support explicit pagination through query options
// Relation loading is controlled by typed `ComicIncl` variadic arguments
// This interface contains only semantic contracts without infra details
type ComicRepo interface {
	// `GetById` retrieves one active comic by id
	GetById(id string, inc ...enum.ComicIncl) (*aggr.Comic, RepoErr)

	// `List` returns active comics matching the query options
	List(opt *query.ListComicOpt, inc ...enum.ComicIncl) ([]*aggr.Comic, RepoErr)

	// `Count` returns active comic count matching query options
	Count(opt *query.ListComicOpt) (int64, RepoErr)

	// `Create` inserts one comic and returns the created aggregate
	Create(cre *aggr.ComicCre) (*aggr.Comic, RepoErr)

	// `Update` applies put-style mutable fields to one comic
	Update(upd *aggr.ComicUpd) RepoErr

	// `UpdateChapterCount` applies delta to chapter counter of one comic.
	UpdateChapterCount(id string, delta int) RepoErr

	// `TouchLastActive` refreshes `last_active_at` and `updated_at` of one comic.
	TouchLastActive(id string) RepoErr

	// `Remove` executes soft delete by setting `deleted_at`
	Remove(id string) RepoErr
}
