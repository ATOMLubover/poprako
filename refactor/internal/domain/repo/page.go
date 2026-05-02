package repo_iface

import (
	"poprako-s/internal/domain/model/aggr"
	"poprako-s/internal/domain/model/query"
)

// `PageRepo` defines persistence contract for page aggregate.
type PageRepo interface {
	// `GetById` returns one page by id.
	GetById(id string) (*aggr.Page, RepoErr)

	// `List` returns pages matching query options.
	List(opt *query.ListPageOpt) ([]*aggr.Page, RepoErr)

	// `Count` returns page count matching query options.
	Count(opt *query.ListPageOpt) (int64, RepoErr)

	// `CreateBatch` inserts many pages in one batch.
	CreateBatch(cre []*aggr.PageCre) RepoErr

	// `MarkImageUploaded` marks one page image as uploaded.
	MarkImageUploaded(id string) RepoErr

	// `SetUnitCounts` overwrites unit count fields of one page.
	SetUnitCounts(id string, total int, translated int, proofread int) RepoErr

	// `DeleteByChapterId` hard-deletes all pages under one chapter.
	DeleteByChapterId(chapterId string) RepoErr
}
