package mock_repo

import (
	"context"
	"sort"
	"strings"
	"time"

	"poprako-s/internal/domain/model"
	iface "poprako-s/internal/domain/repo"
)

type ComicRepo struct {
	Infos map[string]model.ComicInfo
}

func NewMockComicRepo() *ComicRepo {
	return &ComicRepo{Infos: make(map[string]model.ComicInfo)}
}

func WithMockComicRepo(cx context.Context, repo *ComicRepo) context.Context {
	return withMockRepo(cx, comicRepoCtxKey, repo)
}

func NewMockComicRepoFromCx(cx context.Context) (iface.ComicRepo, error) {
	return getMockRepoFromCx[*ComicRepo](cx, comicRepoCtxKey, "comic")
}

func (r *ComicRepo) FromTxnCx(cx context.Context) (iface.ComicRepo, error) {
	return NewMockComicRepoFromCx(cx)
}

func (r *ComicRepo) ensure() {
	if r.Infos == nil {
		r.Infos = make(map[string]model.ComicInfo)
	}
}

func (r *ComicRepo) GetByID(id string) (*model.ComicInfo, error) {
	r.ensure()
	info, ok := r.Infos[id]
	if !ok {
		return nil, errNotFound
	}
	copy := info
	return &copy, nil
}

func (r *ComicRepo) List(opt model.ComicQueryOpt) ([]model.ComicInfo, error) {
	r.ensure()
	items := make([]model.ComicInfo, 0)
	for _, key := range sortedKeys(r.Infos) {
		info := r.Infos[key]
		if opt.ID != nil && info.ID != *opt.ID {
			continue
		}
		if opt.WorksetID != "" && info.WorksetID != opt.WorksetID {
			continue
		}
		if opt.FuzzyTitle != nil && !strings.Contains(info.ComposeComicTitle(), *opt.FuzzyTitle) && !strings.Contains(info.Title, *opt.FuzzyTitle) {
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

func (r *ComicRepo) Count(opt model.ComicQueryOpt) (int64, error) {
	items, err := r.List(opt)
	if err != nil {
		return 0, err
	}
	return int64(len(items)), nil
}

func (r *ComicRepo) Create(c *model.ComicCreation) (*model.ComicInfo, error) {
	r.ensure()
	now := time.Now()
	info := model.ComicInfo{
		ID:           c.ID,
		WorksetID:    c.WorksetID,
		Index:        c.Index,
		Title:        c.Title,
		Author:       c.Author,
		Desc:         c.Desc,
		CreatorID:    c.CreatorID,
		LastActiveAt: now,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	r.Infos[info.ID] = info
	copy := info
	return &copy, nil
}

func (r *ComicRepo) Update(u *model.ComicUpdate) error {
	r.ensure()
	info, ok := r.Infos[u.ID]
	if !ok {
		return errNotFound
	}

	info.Title = u.Title
	info.Author = u.Author
	info.Desc = u.Desc
	info.UpdatedAt = time.Now()
	r.Infos[u.ID] = info

	return nil
}

func (r *ComicRepo) UpdateChapterCount(id string, delta int) error {
	r.ensure()
	info, ok := r.Infos[id]
	if !ok {
		return errNotFound
	}
	info.ChapterCount += delta
	info.UpdatedAt = time.Now()
	r.Infos[id] = info
	return nil
}

func (r *ComicRepo) Delete(id string) error {
	r.ensure()
	delete(r.Infos, id)
	return nil
}

func (r *ComicRepo) PreFillCoverOSSKey(id string, coverOSSKey string) error {
	r.ensure()
	info, ok := r.Infos[id]
	if !ok {
		return errNotFound
	}

	info.CoverOSSKey = coverOSSKey
	info.IsCoverUploaded = false
	info.UpdatedAt = time.Now()
	r.Infos[id] = info

	return nil
}

func (r *ComicRepo) ConfirmCoverUploaded(id string) error {
	r.ensure()
	info, ok := r.Infos[id]
	if !ok {
		return errNotFound
	}

	info.IsCoverUploaded = true
	info.UpdatedAt = time.Now()
	r.Infos[id] = info

	return nil
}
