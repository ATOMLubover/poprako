package repo_iface

import (
	"poprako-s/internal/domain/model/aggr"
	"poprako-s/internal/domain/model/enum"
	"poprako-s/internal/domain/model/query"
)

// `ChapterRepo` defines persistence contract for chapter aggregate.
type ChapterRepo interface {
	// `GetById` retrieves one active chapter by id.
	GetById(id string, inc ...enum.ChapterIncl) (*aggr.Chapter, RepoErr)

	// `FindPinnedByComicId` retrieves pinned active chapter of one comic.
	FindPinnedByComicId(comicId string, inc ...enum.ChapterIncl) (*aggr.Chapter, RepoErr)

	// `List` returns active chapters matching query options.
	List(opt *query.ListChapterOpt, inc ...enum.ChapterIncl) ([]*aggr.Chapter, RepoErr)

	// `Count` returns active chapter count matching query options.
	Count(opt *query.ListChapterOpt) (int64, RepoErr)

	// `Create` inserts one chapter and returns created aggregate.
	Create(cre *aggr.ChapterCre) (*aggr.Chapter, RepoErr)

	// `Update` applies mutable fields to one chapter.
	Update(upd *aggr.ChapterUpd) RepoErr

	// `SetPageCount` overwrites the page count of one chapter.
	SetPageCount(id string, count int) RepoErr

	// `Remove` soft-deletes one chapter.
	Remove(id string) RepoErr
}
