package mock_repo

import (
	"context"
	"sort"
	"time"

	"poprako-s/internal/domain/model"
	iface "poprako-s/internal/domain/repo"
)

type ChapterRepo struct {
	Infos map[string]model.ChapterInfo
}

func NewMockChapterRepo() *ChapterRepo {
	return &ChapterRepo{Infos: make(map[string]model.ChapterInfo)}
}

func WithMockChapterRepo(cx context.Context, repo *ChapterRepo) context.Context {
	return withMockRepo(cx, chapterRepoCtxKey, repo)
}

func NewMockChapterRepoFromCx(cx context.Context) (iface.ChapterRepo, error) {
	return getMockRepoFromCx[*ChapterRepo](cx, chapterRepoCtxKey, "chapter")
}

func (r *ChapterRepo) FromTxnCx(cx context.Context) (iface.ChapterRepo, error) {
	return NewMockChapterRepoFromCx(cx)
}

func (r *ChapterRepo) ensure() {
	if r.Infos == nil {
		r.Infos = make(map[string]model.ChapterInfo)
	}
}

func (r *ChapterRepo) GetByID(id string) (*model.ChapterInfo, error) {
	r.ensure()
	info, ok := r.Infos[id]
	if !ok {
		return nil, errNotFound
	}
	copy := info
	return &copy, nil
}

func (r *ChapterRepo) FindPinnedByComicID(comicID string) (*model.ChapterInfo, error) {
	items, err := r.List(model.ChapterQueryOpt{ComicID: &comicID})
	if err != nil {
		return nil, err
	}
	for _, item := range items {
		if item.IsPinned {
			copy := item
			return &copy, nil
		}
	}
	return nil, errNotFound
}

func (r *ChapterRepo) List(opt model.ChapterQueryOpt) ([]model.ChapterInfo, error) {
	r.ensure()
	items := make([]model.ChapterInfo, 0)
	for _, key := range sortedKeys(r.Infos) {
		info := r.Infos[key]
		if opt.ID != nil && info.ID != *opt.ID {
			continue
		}
		if opt.ComicID != nil && info.ComicID != *opt.ComicID {
			continue
		}
		if opt.CreatorID != nil && info.CreatorID != *opt.CreatorID {
			continue
		}
		items = append(items, info)
	}

	sort.Slice(items, func(i, j int) bool {
		if items[i].Index == items[j].Index {
			return items[i].ID < items[j].ID
		}
		return items[i].Index < items[j].Index
	})

	return items, nil
}

func (r *ChapterRepo) Count(opt model.ChapterQueryOpt) (int64, error) {
	items, err := r.List(opt)
	if err != nil {
		return 0, err
	}
	return int64(len(items)), nil
}

func (r *ChapterRepo) Create(c *model.ChapterCreation) (*model.ChapterInfo, error) {
	r.ensure()
	now := time.Now()
	subtitle := ""
	if c.Subtitle != nil {
		subtitle = *c.Subtitle
	}
	isPinned := false
	if c.IsPinned != nil {
		isPinned = *c.IsPinned
	}

	info := model.ChapterInfo{
		ID:        c.ID,
		ComicID:   c.ComicID,
		IsPinned:  isPinned,
		Index:     c.Index,
		Subtitle:  subtitle,
		CreatorID: c.CreatorID,
		CreatedAt: now,
		UpdatedAt: now,
	}
	r.Infos[info.ID] = info
	copy := info
	return &copy, nil
}

func (r *ChapterRepo) Update(u *model.ChapterUpdate) error {
	r.ensure()
	info, ok := r.Infos[u.ID]
	if !ok {
		return errNotFound
	}

	if u.Subtitle != nil {
		info.Subtitle = *u.Subtitle
	}
	if u.IsPinned != nil {
		info.IsPinned = *u.IsPinned
	}
	info.UploadedAt = cloneTimePtr(u.UploadedAt)
	info.TransalatingAt = cloneTimePtr(u.TransalatingAt)
	info.TranslatedAt = cloneTimePtr(u.TranslatedAt)
	info.ProofreadingAt = cloneTimePtr(u.ProofreadingAt)
	info.ProofreadAt = cloneTimePtr(u.ProofreadAt)
	info.TypesettingAt = cloneTimePtr(u.TypesettingAt)
	info.TypesetAt = cloneTimePtr(u.TypesetAt)
	info.ReviewedAt = cloneTimePtr(u.ReviewedAt)
	info.PublishedAt = cloneTimePtr(u.PublishedAt)
	info.UpdatedAt = time.Now()
	r.Infos[u.ID] = info

	return nil
}

func (r *ChapterRepo) Remove(id string) error {
	r.ensure()
	delete(r.Infos, id)
	return nil
}

func (r *ChapterRepo) LockByID(id string) error {
	r.ensure()
	if _, ok := r.Infos[id]; !ok {
		return errNotFound
	}
	return nil
}

func (r *ChapterRepo) UpdateStats(id string, totalDelta, translatedDelta, proofreadDelta int) error {
	r.ensure()
	info, ok := r.Infos[id]
	if !ok {
		return errNotFound
	}

	info.TotalUnitCount += totalDelta
	info.TranslatedUnitCount += translatedDelta
	info.ProofreadUnitCount += proofreadDelta
	info.UpdatedAt = time.Now()
	r.Infos[id] = info

	return nil
}

func (r *ChapterRepo) UpdatePageCount(id string, delta int) error {
	r.ensure()
	info, ok := r.Infos[id]
	if !ok {
		return errNotFound
	}

	info.PageCount += delta
	info.UpdatedAt = time.Now()
	r.Infos[id] = info

	return nil
}
