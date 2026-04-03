package mock_repo

import (
	"context"
	"time"

	"poprako-s/internal/domain/model"
	iface "poprako-s/internal/domain/repo"
)

type InvitationRepo struct {
	Infos map[string]model.InvitationInfo
}

func NewMockInvitationRepo() *InvitationRepo {
	return &InvitationRepo{Infos: make(map[string]model.InvitationInfo)}
}

func WithMockInvitationRepo(cx context.Context, repo *InvitationRepo) context.Context {
	return withMockRepo(cx, invitationRepoCtxKey, repo)
}

func NewMockInvitationRepoFromCx(cx context.Context) (iface.InvitationRepo, error) {
	return getMockRepoFromCx[*InvitationRepo](cx, invitationRepoCtxKey, "invitation")
}

func (r *InvitationRepo) FromTxnCx(cx context.Context) (iface.InvitationRepo, error) {
	return NewMockInvitationRepoFromCx(cx)
}

func (r *InvitationRepo) ensure() {
	if r.Infos == nil {
		r.Infos = make(map[string]model.InvitationInfo)
	}
}

func (r *InvitationRepo) GetByInviteeQQ(qq string) (*model.InvitationInfo, error) {
	r.ensure()
	for _, key := range sortedKeys(r.Infos) {
		info := r.Infos[key]
		if info.InviteeQQ == qq {
			copy := info
			return &copy, nil
		}
	}
	return nil, errNotFound
}

func (r *InvitationRepo) List(opt model.InvitationQueryOpt) ([]model.InvitationInfo, error) {
	r.ensure()
	items := make([]model.InvitationInfo, 0)
	for _, key := range sortedKeys(r.Infos) {
		info := r.Infos[key]
		if opt.TeamID != nil && info.TeamID != *opt.TeamID {
			continue
		}
		if opt.InvitationCode != nil && info.InvitationCode != *opt.InvitationCode {
			continue
		}
		if opt.Pending && !info.Pending {
			continue
		}
		items = append(items, info)
	}
	return items, nil
}

func (r *InvitationRepo) Create(c *model.InvitationCreation) (*model.InvitationInfo, error) {
	r.ensure()
	info := model.InvitationInfo{
		ID:              c.ID,
		InvitorID:       c.InvitorID,
		InviteeQQ:       c.InviteeQQ,
		TeamID:          c.TargetTeamID,
		InvitationCode:  c.InvitationCode,
		Pending:         true,
		ToBeRawProvider: c.ToBeRawProvider,
		ToBeTranslator:  c.ToBeTranslator,
		ToBeProofreader: c.ToBeProofreader,
		ToBeTypesetter:  c.ToBeTypesetter,
		ToBeReviewer:    c.ToBeReviewer,
		ToBePublisher:   c.ToBePublisher,
		ToBeAdmin:       c.ToBeAdmin,
		CreatedAt:       time.Now(),
	}
	r.Infos[info.ID] = info
	copy := info
	return &copy, nil
}

func (r *InvitationRepo) Update(u *model.InvitationUpdate) error {
	r.ensure()
	info, ok := r.Infos[u.ID]
	if !ok {
		return errNotFound
	}

	info.ToBeRawProvider = u.ToBeRawProvider
	info.ToBeTranslator = u.ToBeTranslator
	info.ToBeProofreader = u.ToBeProofreader
	info.ToBeTypesetter = u.ToBeTypesetter
	info.ToBeReviewer = u.ToBeReviewer
	info.ToBePublisher = u.ToBePublisher
	info.ToBeAdmin = u.ToBeAdmin
	r.Infos[u.ID] = info

	return nil
}

func (r *InvitationRepo) Invalidate(id string) error {
	r.ensure()
	info, ok := r.Infos[id]
	if !ok {
		return errNotFound
	}
	info.Pending = false
	r.Infos[id] = info
	return nil
}

func (r *InvitationRepo) Delete(id string) error {
	r.ensure()
	delete(r.Infos, id)
	return nil
}
