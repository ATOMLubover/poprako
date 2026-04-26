package mock_repo

import (
	"context"
	"sort"
	"time"

	"poprako-s/internal/domain/model"
	iface "poprako-s/internal/domain/repo"
)

type AssignmentRepo struct {
	Infos map[string]model.AssignmentInfo
}

func NewMockAssignmentRepo() *AssignmentRepo {
	return &AssignmentRepo{Infos: make(map[string]model.AssignmentInfo)}
}

func WithMockAssignmentRepo(cx context.Context, repo *AssignmentRepo) context.Context {
	return withMockRepo(cx, assignmentRepoCtxKey, repo)
}

func NewMockAssignmentRepoFromCx(cx context.Context) (iface.AssignmentRepo, error) {
	return getMockRepoFromCx[*AssignmentRepo](cx, assignmentRepoCtxKey, "assignment")
}

func (r *AssignmentRepo) FromTxnCx(cx context.Context) (iface.AssignmentRepo, error) {
	return NewMockAssignmentRepoFromCx(cx)
}

func (r *AssignmentRepo) ensure() {
	if r.Infos == nil {
		r.Infos = make(map[string]model.AssignmentInfo)
	}
}

func (r *AssignmentRepo) GetByID(id string) (*model.AssignmentInfo, error) {
	r.ensure()
	info, ok := r.Infos[id]
	if !ok {
		return nil, errNotFound
	}
	copy := info
	return &copy, nil
}

func (r *AssignmentRepo) Get(opt model.AssignmentQueryOpt) (*model.AssignmentInfo, error) {
	items, err := r.List(opt)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, errNotFound
	}
	return &items[0], nil
}

func (r *AssignmentRepo) List(opt model.AssignmentQueryOpt) ([]model.AssignmentInfo, error) {
	r.ensure()
	items := make([]model.AssignmentInfo, 0)
	for _, key := range sortedKeys(r.Infos) {
		info := r.Infos[key]
		if opt.ChapterID != nil && info.ChapterID != *opt.ChapterID {
			continue
		}
		if opt.UserID != nil && info.UserID != *opt.UserID {
			continue
		}
		items = append(items, info)
	}

	sort.Slice(items, func(i, j int) bool {
		if items[i].CreatedAt.Equal(items[j].CreatedAt) {
			return items[i].ID < items[j].ID
		}
		return items[i].CreatedAt.Before(items[j].CreatedAt)
	})

	return items, nil
}

func (r *AssignmentRepo) Exist(opt model.AssignmentQueryOpt) (bool, error) {
	items, err := r.List(opt)
	if err != nil {
		return false, err
	}
	return len(items) > 0, nil
}

func (r *AssignmentRepo) Create(c *model.AssignmentCreation) (*model.AssignmentInfo, error) {
	r.ensure()
	now := time.Now()
	info := model.AssignmentInfo{
		ID:                    c.ID,
		ChapterID:             c.ChapterID,
		UserID:                c.UserID,
		AssignedRawProviderAt: cloneTimePtr(c.AssignedRawProviderAt),
		AssignedTranslatorAt:  cloneTimePtr(c.AssignedTranslatorAt),
		AssignedProofreaderAt: cloneTimePtr(c.AssignedProofreaderAt),
		AssignedTypesetterAt:  cloneTimePtr(c.AssignedTypesetterAt),
		AssignedRedrawerAt:    cloneTimePtr(c.AssignedRedrawerAt),
		AssignedReviewerAt:    cloneTimePtr(c.AssignedReviewerAt),
		AssignedPublisherAt:   cloneTimePtr(c.AssignedPublisherAt),
		CreatedAt:             now,
		UpdatedAt:             now,
	}
	r.Infos[info.ID] = info
	copy := info
	return &copy, nil
}

func (r *AssignmentRepo) UpsertCreate(c *model.AssignmentCreation) (*model.AssignmentInfo, error) {
	r.ensure()
	now := time.Now()

	for _, info := range r.Infos {
		if info.ChapterID == c.ChapterID && info.UserID == c.UserID {
			info.AssignedRawProviderAt = cloneTimePtr(c.AssignedRawProviderAt)
			info.AssignedTranslatorAt = cloneTimePtr(c.AssignedTranslatorAt)
			info.AssignedProofreaderAt = cloneTimePtr(c.AssignedProofreaderAt)
			info.AssignedTypesetterAt = cloneTimePtr(c.AssignedTypesetterAt)
			info.AssignedRedrawerAt = cloneTimePtr(c.AssignedRedrawerAt)
			info.AssignedReviewerAt = cloneTimePtr(c.AssignedReviewerAt)
			info.AssignedPublisherAt = cloneTimePtr(c.AssignedPublisherAt)
			info.UpdatedAt = now
			r.Infos[info.ID] = info
			copy := info
			return &copy, nil
		}
	}

	info := model.AssignmentInfo{
		ID:                    c.ID,
		ChapterID:             c.ChapterID,
		UserID:                c.UserID,
		AssignedRawProviderAt: cloneTimePtr(c.AssignedRawProviderAt),
		AssignedTranslatorAt:  cloneTimePtr(c.AssignedTranslatorAt),
		AssignedProofreaderAt: cloneTimePtr(c.AssignedProofreaderAt),
		AssignedTypesetterAt:  cloneTimePtr(c.AssignedTypesetterAt),
		AssignedRedrawerAt:    cloneTimePtr(c.AssignedRedrawerAt),
		AssignedReviewerAt:    cloneTimePtr(c.AssignedReviewerAt),
		AssignedPublisherAt:   cloneTimePtr(c.AssignedPublisherAt),
		CreatedAt:             now,
		UpdatedAt:             now,
	}
	r.Infos[info.ID] = info
	copy := info
	return &copy, nil
}

func (r *AssignmentRepo) Delete(id string) error {
	r.ensure()
	delete(r.Infos, id)
	return nil
}
