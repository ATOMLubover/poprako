package mock_repo

import (
	"context"
	"sort"
	"time"

	"poprako-s/internal/domain/model"
	iface "poprako-s/internal/domain/repo"
)

type WorksetRepo struct {
	Infos map[string]model.WorksetInfo
}

func NewMockWorksetRepo() *WorksetRepo {
	return &WorksetRepo{Infos: make(map[string]model.WorksetInfo)}
}

func WithMockWorksetRepo(cx context.Context, repo *WorksetRepo) context.Context {
	return withMockRepo(cx, worksetRepoCtxKey, repo)
}

func NewMockWorksetRepoFromCx(cx context.Context) (iface.WorksetRepo, error) {
	return getMockRepoFromCx[*WorksetRepo](cx, worksetRepoCtxKey, "workset")
}

func (r *WorksetRepo) FromTxnCx(cx context.Context) (iface.WorksetRepo, error) {
	return NewMockWorksetRepoFromCx(cx)
}

func (r *WorksetRepo) ensure() {
	if r.Infos == nil {
		r.Infos = make(map[string]model.WorksetInfo)
	}
}

func (r *WorksetRepo) GetByID(id string) (*model.WorksetInfo, error) {
	r.ensure()
	info, ok := r.Infos[id]
	if !ok {
		return nil, errNotFound
	}
	copy := info
	return &copy, nil
}

func (r *WorksetRepo) List(opt model.WorksetQueryOpt) ([]model.WorksetInfo, error) {
	r.ensure()
	items := make([]model.WorksetInfo, 0)
	for _, key := range sortedKeys(r.Infos) {
		info := r.Infos[key]
		if opt.ID != nil && info.ID != *opt.ID {
			continue
		}
		if opt.TeamID != nil && info.TeamID != *opt.TeamID {
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

func (r *WorksetRepo) Count(opt model.WorksetQueryOpt) (int64, error) {
	items, err := r.List(opt)
	if err != nil {
		return 0, err
	}
	return int64(len(items)), nil
}

func (r *WorksetRepo) Create(c *model.WorksetCreation) (*model.WorksetInfo, error) {
	r.ensure()
	now := time.Now()
	info := model.WorksetInfo{
		ID:          c.ID,
		TeamID:      c.TeamID,
		Index:       c.Index,
		Name:        c.Name,
		Description: c.Description,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	r.Infos[info.ID] = info
	copy := info
	return &copy, nil
}

func (r *WorksetRepo) Update(u *model.WorksetUpdate) error {
	r.ensure()
	info, ok := r.Infos[u.ID]
	if !ok {
		return errNotFound
	}
	info.Name = u.Name
	if u.Description != nil {
		info.Description = *u.Description
	}
	info.UpdatedAt = time.Now()
	r.Infos[u.ID] = info
	return nil
}

func (r *WorksetRepo) UpdateComicCount(id string, delta int) error {
	r.ensure()
	info, ok := r.Infos[id]
	if !ok {
		return errNotFound
	}
	info.ComicCount += delta
	info.UpdatedAt = time.Now()
	r.Infos[id] = info
	return nil
}

func (r *WorksetRepo) Delete(id string) error {
	r.ensure()
	delete(r.Infos, id)
	return nil
}
