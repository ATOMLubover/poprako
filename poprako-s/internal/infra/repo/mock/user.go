package mock_repo

import (
	"context"
	"time"

	"poprako-s/internal/domain/model"
	iface "poprako-s/internal/domain/repo"
)

type UserRepo struct {
	Infos map[string]model.UserInfo
	Creds map[string]model.UserCreds
	Stats map[string]model.UserStats
}

func NewMockUserRepo() *UserRepo {
	return &UserRepo{
		Infos: make(map[string]model.UserInfo),
		Creds: make(map[string]model.UserCreds),
		Stats: make(map[string]model.UserStats),
	}
}

func WithMockUserRepo(cx context.Context, repo *UserRepo) context.Context {
	return withMockRepo(cx, userRepoCtxKey, repo)
}

func NewMockUserRepoFromCx(cx context.Context) (iface.UserRepo, error) {
	return getMockRepoFromCx[*UserRepo](cx, userRepoCtxKey, "user")
}

func (r *UserRepo) FromTxnCx(cx context.Context) (iface.UserRepo, error) {
	return NewMockUserRepoFromCx(cx)
}

func (r *UserRepo) ensure() {
	if r.Infos == nil {
		r.Infos = make(map[string]model.UserInfo)
	}
	if r.Creds == nil {
		r.Creds = make(map[string]model.UserCreds)
	}
	if r.Stats == nil {
		r.Stats = make(map[string]model.UserStats)
	}
}

func (r *UserRepo) GetCredsByQQ(qq string) (*model.UserCreds, error) {
	r.ensure()
	creds, ok := r.Creds[qq]
	if !ok {
		return nil, errNotFound
	}
	copy := creds
	return &copy, nil
}

func (r *UserRepo) GetByID(id string) (*model.UserInfo, error) {
	r.ensure()
	info, ok := r.Infos[id]
	if !ok {
		return nil, errNotFound
	}
	copy := info
	return &copy, nil
}

func (r *UserRepo) GetByQQ(qq string) (*model.UserInfo, error) {
	r.ensure()
	for _, key := range sortedKeys(r.Infos) {
		info := r.Infos[key]
		if info.QQ == qq {
			copy := info
			return &copy, nil
		}
	}
	return nil, errNotFound
}

func (r *UserRepo) List(opt model.UserQueryOpt) ([]model.UserInfo, error) {
	r.ensure()
	items := make([]model.UserInfo, 0)
	for _, key := range sortedKeys(r.Infos) {
		info := r.Infos[key]
		if opt.ID != nil && info.ID != *opt.ID {
			continue
		}
		if opt.QQ != nil && info.QQ != *opt.QQ {
			continue
		}
		items = append(items, info)
	}
	return items, nil
}

func (r *UserRepo) Create(c *model.UserCreation) (*model.UserInfo, error) {
	r.ensure()
	now := time.Now()
	info := model.UserInfo{
		ID:          c.ID,
		Name:        c.Name,
		QQ:          c.QQ,
		LastLoginAt: now,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	r.Infos[info.ID] = info
	r.Creds[c.QQ] = model.UserCreds{QQ: c.QQ, PwdHash: c.PwdHash}
	copy := info
	return &copy, nil
}

func (r *UserRepo) Update(u *model.UserUpdate) error {
	r.ensure()
	info, ok := r.Infos[u.ID]
	if !ok {
		return errNotFound
	}
	info.Name = u.Name
	info.QQ = u.QQ
	info.UpdatedAt = time.Now()
	r.Infos[u.ID] = info
	return nil
}

func (r *UserRepo) Remove(id string) error {
	r.ensure()
	delete(r.Infos, id)
	for qq, creds := range r.Creds {
		user, err := r.GetByQQ(qq)
		if err != nil || user.ID != id {
			continue
		}
		delete(r.Creds, creds.QQ)
	}
	delete(r.Stats, id)
	return nil
}

func (r *UserRepo) RefreshLastLogin(qq string, t time.Time) error {
	r.ensure()
	for key, info := range r.Infos {
		if info.QQ != qq {
			continue
		}
		info.LastLoginAt = t
		info.UpdatedAt = time.Now()
		r.Infos[key] = info
		return nil
	}
	return errNotFound
}

func (r *UserRepo) PreFillAvatarOSSKey(id string, avatarOSSKey string) error {
	r.ensure()
	info, ok := r.Infos[id]
	if !ok {
		return errNotFound
	}
	info.AvatarKey = avatarOSSKey
	info.IsAvatarUploaded = false
	info.UpdatedAt = time.Now()
	r.Infos[id] = info
	return nil
}

func (r *UserRepo) ConfirmAvatarUploaded(id string) error {
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

func (r *UserRepo) GetOrCreateStats(userID string) (*model.UserStats, error) {
	r.ensure()
	stats, ok := r.Stats[userID]
	if !ok {
		stats = model.UserStats{UserID: userID}
		r.Stats[userID] = stats
	}
	copy := stats
	return &copy, nil
}

func (r *UserRepo) PatchStats(stats *model.UserStatsPatch) error {
	r.ensure()
	current, ok := r.Stats[stats.UserID]
	if !ok {
		current = model.UserStats{UserID: stats.UserID}
	}
	current.TotalAssignmentCount += stats.TotalAssignmentCountDelta
	current.ActiveAssignmentCount += stats.ActiveAssignmentCountDelta
	current.FinishedAssignmentCount += stats.FinishedAssignmentCountDelta
	r.Stats[stats.UserID] = current
	return nil
}
