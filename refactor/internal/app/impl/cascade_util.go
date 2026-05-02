package app_impl

import (
	"poprako-s/internal/domain/model/aggr"
	"poprako-s/internal/domain/model/enum"
	"poprako-s/internal/domain/model/query"
	repo_iface "poprako-s/internal/domain/repo"
	"poprako-s/internal/domain/svc"
)

// Helper list batch-size constants used by app-layer pagination loops.
const (
	// `listBatchSize` is the default batch size used by helper pagination loops.
	listBatchSize = 500
)

// `listAllComics` loads all comics under one workset with explicit pagination.
func listAllComics(comicRepo repo_iface.ComicRepo, worksetId string) ([]*aggr.Comic, repo_iface.RepoErr) {
	items := make([]*aggr.Comic, 0)
	offset := 0

	for {
		batch, err := comicRepo.List(&query.ListComicOpt{
			WorksetId: &worksetId,
			Pagi: query.PagiOpt{
				Offset: offset,
				Limit:  listBatchSize,
			},
		})
		if err != nil {
			return nil, err
		}

		items = append(items, batch...)

		if len(batch) < listBatchSize {
			return items, nil
		}

		offset += len(batch)
	}
}

// `listAllChapters` loads all chapters under one comic with explicit pagination.
func listAllChapters(chapterRepo repo_iface.ChapterRepo, comicId string) ([]*aggr.Chapter, repo_iface.RepoErr) {
	items := make([]*aggr.Chapter, 0)
	offset := 0

	for {
		batch, err := chapterRepo.List(&query.ListChapterOpt{
			ComicId: &comicId,
			Pagi: query.PagiOpt{
				Offset: offset,
				Limit:  listBatchSize,
			},
		})
		if err != nil {
			return nil, err
		}

		items = append(items, batch...)

		if len(batch) < listBatchSize {
			return items, nil
		}

		offset += len(batch)
	}
}

// `listAllPages` loads all pages under one chapter with explicit pagination.
func listAllPages(pageRepo repo_iface.PageRepo, chapterId string) ([]*aggr.Page, repo_iface.RepoErr) {
	items := make([]*aggr.Page, 0)
	offset := 0

	for {
		batch, err := pageRepo.List(&query.ListPageOpt{
			ChapterId: &chapterId,
			Pagi: query.PagiOpt{
				Offset: offset,
				Limit:  listBatchSize,
			},
		})
		if err != nil {
			return nil, err
		}

		items = append(items, batch...)

		if len(batch) < listBatchSize {
			return items, nil
		}

		offset += len(batch)
	}
}

// `listAllAssignments` loads all assignments under one chapter with explicit pagination.
func listAllAssignments(assignmentRepo repo_iface.AssignmentRepo, chapterId string) ([]*aggr.Assignment, repo_iface.RepoErr) {
	items := make([]*aggr.Assignment, 0)
	offset := 0

	for {
		batch, err := assignmentRepo.List(&query.ListAssignmentOpt{
			ChapterId: &chapterId,
			Pagi: query.PagiOpt{
				Offset: offset,
				Limit:  listBatchSize,
			},
		})
		if err != nil {
			return nil, err
		}

		items = append(items, batch...)

		if len(batch) < listBatchSize {
			return items, nil
		}

		offset += len(batch)
	}
}

// `collectAssignedUserIds` extracts assignment user ids in result order.
func collectAssignedUserIds(assignments []*aggr.Assignment) []string {
	userIds := make([]string, 0, len(assignments))

	for i := range assignments {
		userIds = append(userIds, assignments[i].UserId)
	}

	return userIds
}

// `savePageImageDeleteMsgs` enqueues OSS delete messages for all stored page images.
func savePageImageDeleteMsgs(
	ossMsgSvc svc.OssMsgSvc,
	ossMsgRepo repo_iface.OssMsgRepo,
	pages []*aggr.Page,
) error {
	for i := range pages {
		if pages[i].ImageKey == nil || *pages[i].ImageKey == "" {
			continue
		}

		if err := ossMsgSvc.SavePendingDel(
			ossMsgRepo,
			enum.OssResPageImage,
			pages[i].Id,
			[]string{*pages[i].ImageKey},
		); err != nil {
			return err
		}
	}

	return nil
}

// `deleteChapterCascade` removes one chapter and all descendant rows atomically.
func deleteChapterCascade(
	chapter *aggr.Chapter,
	pageRepo repo_iface.PageRepo,
	assignmentRepo repo_iface.AssignmentRepo,
	chapterRepo repo_iface.ChapterRepo,
	ossMsgRepo repo_iface.OssMsgRepo,
	ossMsgSvc svc.OssMsgSvc,
) ([]string, error) {
	// Load assignments before deleting rows so event payload stays complete.
	assignments, err := listAllAssignments(assignmentRepo, chapter.Id)
	if err != nil {
		return nil, err
	}

	// Load pages before deleting rows so OSS cleanup stays complete.
	pages, err := listAllPages(pageRepo, chapter.Id)
	if err != nil {
		return nil, err
	}

	// Enqueue remote page-image cleanup before deleting local rows.
	if err := savePageImageDeleteMsgs(ossMsgSvc, ossMsgRepo, pages); err != nil {
		return nil, err
	}

	// Delete descendant rows before removing the chapter row itself.
	if err := pageRepo.DeleteByChapterId(chapter.Id); err != nil {
		return nil, err
	}

	if err := assignmentRepo.DeleteByChapterId(chapter.Id); err != nil {
		return nil, err
	}

	if err := chapterRepo.Delete(chapter.Id); err != nil {
		return nil, err
	}

	return collectAssignedUserIds(assignments), nil
}
