package mock_repo

import (
	"context"
	"time"

	"poprako-s/internal/domain/model"
	iface "poprako-s/internal/domain/repo"
)

type TeamRepo struct {
	Infos map[string]model.TeamInfo
}

func NewMockTeamRepo() *TeamRepo {
	return &TeamRepo{Infos: make(map[string]model.TeamInfo)}
}

func WithMockTeamRepo(cx context.Context, repo *TeamRepo) context.Context {
	return withMockRepo(cx, teamRepoCtxKey, repo)
}

func NewMockTeamRepoFromCx(cx context.Context) (iface.TeamRepo, error) {
	return getMockRepoFromCx[*TeamRepo](cx, teamRepoCtxKey, "team")
}

func (r *TeamRepo) FromTxnCx(cx context.Context) (iface.TeamRepo, error) {
	return NewMockTeamRepoFromCx(cx)
}

func (r *TeamRepo) ensure() {
	if r.Infos == nil {
		r.Infos = make(map[string]model.TeamInfo)
	}
}

func (r *TeamRepo) GetByID(id string) (*model.TeamInfo, error) {
	r.ensure()
	info, ok := r.Infos[id]
	if !ok {
		return nil, errNotFound
	}
	copy := info
	return &copy, nil
}

func (r *TeamRepo) List(opt model.TeamQueryOpt) ([]model.TeamInfo, error) {
	r.ensure()
	items := make([]model.TeamInfo, 0)
	for _, key := range sortedKeys(r.Infos) {
		info := r.Infos[key]
		if opt.ID != nil && info.ID != *opt.ID {
			continue
		}
		items = append(items, info)
	}
	return items, nil
}

func (r *TeamRepo) Create(c *model.TeamCreation) (*model.TeamInfo, error) {
	r.ensure()
	now := time.Now()
	info := model.TeamInfo{
		ID:        c.ID,
		Name:      c.Name,
		Desc:      c.Desc,
		CreatedAt: now,
		UpdatedAt: now,
	}
	r.Infos[info.ID] = info
	copy := info
	return &copy, nil
}

func (r *TeamRepo) Update(u *model.TeamUpdate) error {
	r.ensure()
	info, ok := r.Infos[u.ID]
	if !ok {
		return errNotFound
	}
	info.Name = u.Name
	info.Desc = u.Desc
	info.UpdatedAt = time.Now()
	r.Infos[u.ID] = info
	return nil
}

func (r *TeamRepo) Delete(id string) error {
	r.ensure()
	delete(r.Infos, id)
	return nil
}

func (r *TeamRepo) PreFillAvatarOSSKey(id string, avatarOSSKey string) error {
	r.ensure()
	info, ok := r.Infos[id]
	if !ok {
		return errNotFound
	}
	info.AvatarOSSKey = avatarOSSKey
	info.IsAvatarUploaded = false
	info.UpdatedAt = time.Now()
	r.Infos[id] = info
	return nil
}

func (r *TeamRepo) ConfirmAvatarUploaded(id string) error {
	r.ensure()
	info, ok := r.Infos[id]
	if !ok {
		return errNotFound
	}
	info.IsAvatarUploaded = true
	info.UpdatedAt = time.Now()
	r.Infos[id] = info
	return nil
}
