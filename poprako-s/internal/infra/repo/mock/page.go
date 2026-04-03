package mock_repo

import (
	"context"
	"sort"
	"time"

	"poprako-s/internal/domain/model"
	iface "poprako-s/internal/domain/repo"
)

type PageRepo struct {
	Infos map[string]model.PageInfo
}

func NewMockPageRepo() *PageRepo {
	return &PageRepo{Infos: make(map[string]model.PageInfo)}
}

func WithMockPageRepo(cx context.Context, repo *PageRepo) context.Context {
	return withMockRepo(cx, pageRepoCtxKey, repo)
}

func NewMockPageRepoFromCx(cx context.Context) (iface.PageRepo, error) {
	return getMockRepoFromCx[*PageRepo](cx, pageRepoCtxKey, "page")
}

func (r *PageRepo) FromTxnCx(cx context.Context) (iface.PageRepo, error) {
	return NewMockPageRepoFromCx(cx)
}

func (r *PageRepo) ensure() {
	if r.Infos == nil {
		r.Infos = make(map[string]model.PageInfo)
	}
}

func (r *PageRepo) GetByID(id string) (*model.PageInfo, error) {
	r.ensure()
	info, ok := r.Infos[id]
	if !ok {
		return nil, errNotFound
	}
	copy := info
	return &copy, nil
}

func (r *PageRepo) List(opt model.PageQueryOpt) ([]model.PageInfo, error) {
	r.ensure()
	items := make([]model.PageInfo, 0)
	for _, key := range sortedKeys(r.Infos) {
		info := r.Infos[key]
		if opt.ID != nil && info.ID != *opt.ID {
			continue
		}
		if opt.ChapterID != nil && info.ChapterID != *opt.ChapterID {
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

func (r *PageRepo) GetStatsByID(pageID string) (*model.PageStats, error) {
	r.ensure()
	info, ok := r.Infos[pageID]
	if !ok {
		return nil, errNotFound
	}
	stats := model.PageStats{
		PageID:              pageID,
		TotalUnitCount:      info.TotalUnitCount,
		TranslatedUnitCount: info.TranslatedUnitCount,
		ProofreadUnitCount:  info.ProofreadUnitCount,
	}
	return &stats, nil
}

func (r *PageRepo) CreateBatch(pages []*model.PageCreation) error {
	r.ensure()
	now := time.Now()
	for _, page := range pages {
		r.Infos[page.ID] = model.PageInfo{
			ID:        page.ID,
			ChapterID: page.ChapterID,
			Index:     page.Index,
			OSSKey:    page.OSSKey,
			CreatorID: page.CreatorID,
			CreatedAt: now,
			UpdatedAt: now,
		}
	}
	return nil
}

func (r *PageRepo) Update(u *model.PageUpdate) error {
	r.ensure()
	info, ok := r.Infos[u.ID]
	if !ok {
		return errNotFound
	}

	info.Index = u.Index
	info.OSSKey = u.OSSKey
	info.IsUploaded = u.IsUploaded
	info.TotalUnitCount = u.TotalUnitCount
	info.TranslatedUnitCount = u.TranslatedUnitCount
	info.ProofreadUnitCount = u.ProofreadUnitCount
	info.UpdatedAt = time.Now()
	r.Infos[u.ID] = info

	return nil
}

func (r *PageRepo) UpdateStats(stats *model.PageStats) error {
	r.ensure()
	info, ok := r.Infos[stats.PageID]
	if !ok {
		return errNotFound
	}
	info.TotalUnitCount = stats.TotalUnitCount
	info.TranslatedUnitCount = stats.TranslatedUnitCount
	info.ProofreadUnitCount = stats.ProofreadUnitCount
	info.UpdatedAt = time.Now()
	r.Infos[stats.PageID] = info
	return nil
}

func (r *PageRepo) Delete(id string) error {
	r.ensure()
	delete(r.Infos, id)
	return nil
}

func (r *PageRepo) DeleteBatch(ids []string) error {
	r.ensure()
	for _, id := range ids {
		delete(r.Infos, id)
	}
	return nil
}
