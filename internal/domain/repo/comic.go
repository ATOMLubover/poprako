package repo_iface

import (
	"poprako-s/internal/domain/model/aggr"
	"poprako-s/internal/domain/model/enum"
	"poprako-s/internal/domain/model/query"
)

// `ComicRepo` defines the persistence contract for `Comic` aggregate
// List style methods support explicit pagination through query options
// Relation loading is controlled by typed `ComicIncl` variadic arguments
// This interface contains only semantic contracts without infra details
type ComicRepo interface {
	// `GetById` retrieves one comic by id
	GetById(id string, inc ...enum.ComicIncl) (*aggr.Comic, RepoErr)

	// `List` returns comics matching the query options
	List(opt *query.ListComicOpt, inc ...enum.ComicIncl) ([]*aggr.Comic, RepoErr)

	// `Count` returns comic count matching query options
	Count(opt *query.ListComicOpt) (int64, RepoErr)

	// `Create` inserts one comic and returns the created aggregate
	Create(cre *aggr.ComicCre) (*aggr.Comic, RepoErr)

	// `Update` applies put-style mutable fields to one comic
	Update(upd *aggr.ComicUpd) RepoErr

	// `UpdateChapterCount` applies delta to chapter counter of one comic.
	UpdateChapterCount(id string, delta int) RepoErr

	// `TouchLastActive` refreshes `last_active_at` and `updated_at` of one comic.
	TouchLastActive(id string) RepoErr

	// `PrefillCoverKey` writes one reserved cover key before client upload starts.
	PrefillCoverKey(id string, key string) RepoErr

	// `MarkCoverUploaded` marks one comic cover upload as completed.
	MarkCoverUploaded(id string) RepoErr

	// `IncrementChapterNextIndex` allocates one chapter index from comic-scoped sequence.
	// It returns the allocated index value for immediate insert use.
	IncrementChapterNextIndex(id string) (int, RepoErr)

	// `Delete` hard-deletes one comic row by id.
	Delete(id string) RepoErr

	// `MarkCompleted` sets `is_completed` to true for one comic.
	MarkCompleted(id string) RepoErr
}
