package mock_repo

import (
	"context"
	"sort"
	"time"

	"poprako-s/internal/domain/model"
	iface "poprako-s/internal/domain/repo"
)

type MemberRepo struct {
	Infos map[string]model.MemberInfo
}

func NewMockMemberRepo() *MemberRepo {
	return &MemberRepo{Infos: make(map[string]model.MemberInfo)}
}

func WithMockMemberRepo(cx context.Context, repo *MemberRepo) context.Context {
	return withMockRepo(cx, memberRepoCtxKey, repo)
}

func NewMockMemberRepoFromCx(cx context.Context) (iface.MemberRepo, error) {
	return getMockRepoFromCx[*MemberRepo](cx, memberRepoCtxKey, "member")
}

func (r *MemberRepo) FromTxnCx(cx context.Context) (iface.MemberRepo, error) {
	return NewMockMemberRepoFromCx(cx)
}

func (r *MemberRepo) ensure() {
	if r.Infos == nil {
		r.Infos = make(map[string]model.MemberInfo)
	}
}

func (r *MemberRepo) GetByID(id string) (*model.MemberInfo, error) {
	r.ensure()
	info, ok := r.Infos[id]
	if !ok {
		return nil, errNotFound
	}
	copy := info
	return &copy, nil
}

func (r *MemberRepo) Get(opt model.MemberQueryOpt) (*model.MemberInfo, error) {
	items, err := r.List(opt)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, errNotFound
	}
	return &items[0], nil
}

func (r *MemberRepo) List(opt model.MemberQueryOpt) ([]model.MemberInfo, error) {
	r.ensure()
	items := make([]model.MemberInfo, 0)
	for _, key := range sortedKeys(r.Infos) {
		info := r.Infos[key]
		if opt.ID != nil && info.ID != *opt.ID {
			continue
		}
		if opt.UserID != nil && info.UserID != *opt.UserID {
			continue
		}
		if opt.TeamID != nil && info.TeamID != *opt.TeamID {
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

func (r *MemberRepo) Exist(opt model.MemberQueryOpt) (bool, error) {
	items, err := r.List(opt)
	if err != nil {
		return false, err
	}
	return len(items) > 0, nil
}

func (r *MemberRepo) Create(c *model.MemberCreation) (*model.MemberInfo, error) {
	r.ensure()
	now := time.Now()
	info := model.MemberInfo{
		ID:        c.ID,
		UserID:    c.UserID,
		TeamID:    c.TeamID,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if c.ToBeRawProvider {
		info.AssignedRawProviderAt = cloneTimePtr(&now)
	}
	if c.ToBeTranslator {
		info.AssignedTranslatorAt = cloneTimePtr(&now)
	}
	if c.ToBeProofreader {
		info.AssignedProofreaderAt = cloneTimePtr(&now)
	}
	if c.ToBeTypesetter {
		info.AssignedTypesetterAt = cloneTimePtr(&now)
	}
	if c.ToBeRedrawer {
		info.AssignedRedrawerAt = cloneTimePtr(&now)
	}
	if c.ToBeReviewer {
		info.AssignedReviewerAt = cloneTimePtr(&now)
	}
	if c.ToBePublisher {
		info.AssignedPublisherAt = cloneTimePtr(&now)
	}
	if c.ToBeAdmin {
		info.AssignedAdminAt = cloneTimePtr(&now)
	}
	r.Infos[info.ID] = info
	copy := info
	return &copy, nil
}

func (r *MemberRepo) Update(u *model.MemberUpdate) error {
	r.ensure()
	info, ok := r.Infos[u.ID]
	if !ok {
		return errNotFound
	}

	info.AssignedRawProviderAt = cloneTimePtr(u.AssignedRawProviderAt)
	info.AssignedTranslatorAt = cloneTimePtr(u.AssignedTranslatorAt)
	info.AssignedProofreaderAt = cloneTimePtr(u.AssignedProofreaderAt)
	info.AssignedTypesetterAt = cloneTimePtr(u.AssignedTypesetterAt)
	info.AssignedRedrawerAt = cloneTimePtr(u.AssignedRedrawerAt)
	info.AssignedReviewerAt = cloneTimePtr(u.AssignedReviewerAt)
	info.AssignedPublisherAt = cloneTimePtr(u.AssignedPublisherAt)
	info.AssignedAdminAt = cloneTimePtr(u.AssignedAdminAt)
	info.UpdatedAt = time.Now()
	r.Infos[u.ID] = info

	return nil
}

func (r *MemberRepo) Delete(id string) error {
	r.ensure()
	delete(r.Infos, id)
	return nil
}
