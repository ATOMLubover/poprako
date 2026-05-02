package repo_iface

import "poprako-s/internal/domain/model/aggr"

type UnitRepo interface {
	// `ListByPage` returns units under one page ordered by `Index` ascending.
	ListByPage(pageId string) ([]*aggr.Unit, RepoErr)

	ListIndicesByPage(pageId string) ([]*aggr.UnitIndex, RepoErr)

	// `CountByPage` returns total, translated, and proofread unit counts of one page.
	CountByPage(pageId string) (int, int, int, RepoErr)

	Create(cre *aggr.UnitCre) RepoErr

	Save(sv *aggr.UnitSave) RepoErr

	Reindex(indices []*aggr.UnitIndex) RepoErr

	Delete(id string) RepoErr
}
