package mock_repo

import (
	"context"
	"sort"

	"poprako-s/internal/domain/model"
	iface "poprako-s/internal/domain/repo"
)

type UnitRepo struct {
	Infos map[string]model.UnitInfo
}

func NewMockUnitRepo() *UnitRepo {
	return &UnitRepo{Infos: make(map[string]model.UnitInfo)}
}

func WithMockUnitRepo(cx context.Context, repo *UnitRepo) context.Context {
	return withMockRepo(cx, unitRepoCtxKey, repo)
}

func NewMockUnitRepoFromCx(cx context.Context) (iface.UnitRepo, error) {
	return getMockRepoFromCx[*UnitRepo](cx, unitRepoCtxKey, "unit")
}

func (r *UnitRepo) FromTxnCx(cx context.Context) (iface.UnitRepo, error) {
	return NewMockUnitRepoFromCx(cx)
}

func (r *UnitRepo) ensure() {
	if r.Infos == nil {
		r.Infos = make(map[string]model.UnitInfo)
	}
}

func (r *UnitRepo) List(opt model.UnitQueryOpt) ([]model.UnitInfo, error) {
	r.ensure()
	items := make([]model.UnitInfo, 0)
	for _, key := range sortedKeys(r.Infos) {
		info := r.Infos[key]
		if info.PageID != opt.PageID {
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

func (r *UnitRepo) CreateBatch(units []*model.UnitCreation) error {
	r.ensure()
	for _, unit := range units {
		r.Infos[unit.ID] = *unit
	}
	return nil
}

func (r *UnitRepo) PatchBatch(patches []*model.UnitPatch) error {
	r.ensure()
	for _, patch := range patches {
		info, ok := r.Infos[patch.ID]
		if !ok {
			continue
		}

		if patch.Index != nil {
			info.Index = *patch.Index
		}
		if patch.XCoord != nil {
			info.XCoord = *patch.XCoord
		}
		if patch.YCoord != nil {
			info.YCoord = *patch.YCoord
		}
		if patch.IsBubble != nil {
			info.IsBubble = *patch.IsBubble
		}
		if patch.TranslatedText != nil {
			info.TranslatedText = *patch.TranslatedText
		}
		if patch.TranslatorID != nil {
			info.TranslatorID = *patch.TranslatorID
		}
		if patch.TranslatorComment != nil {
			info.TranslatorComment = *patch.TranslatorComment
		}
		if patch.IsProofread != nil {
			info.IsProofread = *patch.IsProofread
		}
		if patch.ProofreadText != nil {
			info.ProofreadText = *patch.ProofreadText
		}
		if patch.ProofreaderID != nil {
			info.ProofreaderID = *patch.ProofreaderID
		}
		if patch.ProofreaderComment != nil {
			info.ProofreaderComment = *patch.ProofreaderComment
		}

		r.Infos[patch.ID] = info
	}
	return nil
}

func (r *UnitRepo) DeleteBatch(unitIDs []string) error {
	r.ensure()
	for _, id := range unitIDs {
		delete(r.Infos, id)
	}
	return nil
}
