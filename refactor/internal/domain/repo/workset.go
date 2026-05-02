package repo_iface

import (
	"poprako-s/internal/domain/model/aggr"
	"poprako-s/internal/domain/model/enum"
	"poprako-s/internal/domain/model/query"
)

// `WorksetRepo` defines the persistence contract for the `Workset` aggregate.
type WorksetRepo interface {
	// `GetById` retrieves an active workset by its unique identifier.
	// Returns an error if the workset does not exist.
	GetById(id string, inc ...enum.WorksetIncl) (*aggr.Workset, RepoErr)

	// `List` returns all worksets that match the given options,
	// ordered by `index` ascending.
	List(opt *query.ListWorksetOpt, inc ...enum.WorksetIncl) ([]*aggr.Workset, RepoErr)

	// `Count` returns the number of worksets matching the given options.
	Count(opt *query.ListWorksetOpt) (int64, RepoErr)

	// `Create` persists a new workset from the given creation input and returns
	// the fully populated aggregate.
	Create(cre *aggr.WorksetCre) (*aggr.Workset, RepoErr)

	// `Update` applies the mutable fields in `upd` to an existing workset.
	Update(upd *aggr.WorksetUpd) RepoErr

	// `UpdateComicCount` applies delta to one workset comic counter.
	// The counter update is usually called inside app transaction flows.
	UpdateComicCount(id string, delta int) RepoErr

	// `IncrementComicNextIndex` allocates one comic index from workset-scoped sequence.
	// It returns the allocated index value for immediate insert use.
	IncrementComicNextIndex(id string) (int, RepoErr)

	// `Delete` hard-deletes one workset row by id.
	Delete(id string) RepoErr
}
