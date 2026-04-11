package mock_repo

import (
	"context"
	"time"

	"poprako-s/internal/domain/model"
	iface "poprako-s/internal/domain/repo"
)

type ChapterInvitationRepo struct {
	Infos map[string]model.ChapterInvitationInfo
}

func NewMockChapterInvitationRepo() *ChapterInvitationRepo {
	return &ChapterInvitationRepo{Infos: make(map[string]model.ChapterInvitationInfo)}
}

func WithMockChapterInvitationRepo(cx context.Context, repo *ChapterInvitationRepo) context.Context {
	return withMockRepo(cx, chapterInvitationRepoCtxKey, repo)
}

func NewMockChapterInvitationRepoFromCx(cx context.Context) (iface.ChapterInvitationRepo, error) {
	return getMockRepoFromCx[*ChapterInvitationRepo](cx, chapterInvitationRepoCtxKey, "chapter invitation")
}

func (r *ChapterInvitationRepo) FromTxnCx(cx context.Context) (iface.ChapterInvitationRepo, error) {
	return NewMockChapterInvitationRepoFromCx(cx)
}

func (r *ChapterInvitationRepo) ensure() {
	if r.Infos == nil {
		r.Infos = make(map[string]model.ChapterInvitationInfo)
	}
}

func (r *ChapterInvitationRepo) List(opt model.ChapterInvitationQueryOpt) ([]model.ChapterInvitationInfo, error) {
	r.ensure()
	items := make([]model.ChapterInvitationInfo, 0)

	for _, key := range sortedKeys(r.Infos) {
		info := r.Infos[key]

		if opt.ChapterID != nil && info.ChapterID != *opt.ChapterID {
			continue
		}
		if opt.InvitationCode != nil && info.InvitationCode != *opt.InvitationCode {
			continue
		}
		if opt.InviteeQQ != nil && info.InviteeQQ != *opt.InviteeQQ {
			continue
		}
		if opt.OnlyPendingTrue && !info.Pending {
			continue
		}

		items = append(items, info)
	}

	return items, nil
}

func (r *ChapterInvitationRepo) Create(c *model.ChapterInvitationCreation) (*model.ChapterInvitationInfo, error) {
	r.ensure()
	now := time.Now()
	info := model.ChapterInvitationInfo{
		ID:              c.ID,
		ChapterID:       c.ChapterID,
		InviterID:       c.InviterID,
		InviteeQQ:       c.InviteeQQ,
		InvitationCode:  c.InvitationCode,
		Pending:         true,
		ToBeRawProvider: c.ToBeRawProvider,
		ToBeTranslator:  c.ToBeTranslator,
		ToBeProofreader: c.ToBeProofreader,
		ToBeTypesetter:  c.ToBeTypesetter,
		ToBeRedrawer:    c.ToBeRedrawer,
		ToBeReviewer:    c.ToBeReviewer,
		ToBePublisher:   c.ToBePublisher,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	r.Infos[info.ID] = info
	copy := info
	return &copy, nil
}

func (r *ChapterInvitationRepo) Invalidate(id string) error {
	r.ensure()
	info, ok := r.Infos[id]
	if !ok {
		return errNotFound
	}

	info.Pending = false
	info.UpdatedAt = time.Now()
	r.Infos[id] = info
	return nil
}
