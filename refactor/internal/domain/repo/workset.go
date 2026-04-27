package repo_iface

import (
	"poprako-s/internal/domain/model/aggr"
	"poprako-s/internal/domain/model/enum"
	"poprako-s/internal/domain/model/query"
)

// `WorksetRepo` defines the persistence contract for the `Workset` aggregate.
type WorksetRepo interface {
	// `GetById` retrieves an active workset by its unique identifier.
	// Returns an error if the workset does not exist or has been soft-deleted.
	GetById(id string, inc ...enum.WorksetIncl) (*aggr.Workset, RepoErr)

	// `List` returns all active worksets that match the given options,
	// ordered by `index` ascending.
	List(opt *query.ListWorksetOpt, inc ...enum.WorksetIncl) ([]*aggr.Workset, RepoErr)

	// `Count` returns the number of active worksets matching the given options.
	// Used to determine the next available `index` before a new workset is created.
	// NOTE: optimistic lock on `index` is the responsibility of the caller to ensure consistency.
	Count(opt *query.ListWorksetOpt) (int64, RepoErr)

	// `Create` persists a new workset from the given creation input and returns
	// the fully populated aggregate.
	Create(cre *aggr.WorksetCre) (*aggr.Workset, RepoErr)

	// `Update` applies the mutable fields in `upd` to an existing workset.
	Update(upd *aggr.WorksetUpd) RepoErr

	// `Remove` marks the workset as deleted by setting its `deleted_at`
	// timestamp without physically removing the row.
	Remove(id string) RepoErr
}
