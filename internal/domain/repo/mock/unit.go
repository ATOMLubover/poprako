package repo_mock

import (
	"sort"

	"poprako-s/internal/domain/model/aggr"
	repo_iface "poprako-s/internal/domain/repo"
)

type mockUnitRepo struct {
	units map[string]*aggr.Unit
}

// `NewMockUnitRepo` returns an in-memory `UnitRepo` for tests.
func NewMockUnitRepo(initUnits []*aggr.Unit) repo_iface.UnitRepo {
	units := make(map[string]*aggr.Unit)
	for _, u := range initUnits {
		units[u.Id] = u
	}

	return &mockUnitRepo{
		units: units,
	}
}

// `ListMockUnitsByPage` returns cloned units of one page sorted by `Index`.
func ListMockUnitsByPage(repo repo_iface.UnitRepo, pageId string) []*aggr.Unit {
	mockRepo, ok := repo.(*mockUnitRepo)
	if !ok {
		return nil
	}

	units := make([]*aggr.Unit, 0)
	for _, unit := range mockRepo.units {
		if unit.PageId != pageId {
			continue
		}

		clone := *unit
		units = append(units, &clone)
	}

	sort.Slice(units, func(i, j int) bool {
		if units[i].Index == units[j].Index {
			return units[i].Id < units[j].Id
		}

		return units[i].Index < units[j].Index
	})

	return units
}

// `GetMockUnitById` returns a cloned unit by id from the in-memory repo.
func GetMockUnitById(repo repo_iface.UnitRepo, id string) *aggr.Unit {
	mockRepo, ok := repo.(*mockUnitRepo)
	if !ok {
		return nil
	}

	unit, exists := mockRepo.units[id]
	if !exists {
		return nil
	}

	clone := *unit

	return &clone
}

// `CountMockUnitsByPage` returns unit counts for one page from the in-memory repo.
func CountMockUnitsByPage(repo repo_iface.UnitRepo, pageId string) (int, int, int) {
	mockRepo, ok := repo.(*mockUnitRepo)
	if !ok {
		return 0, 0, 0
	}

	total := 0
	translated := 0
	proofread := 0
	for _, unit := range mockRepo.units {
		if unit.PageId != pageId {
			continue
		}

		total++

		if unit.TranslatedText != nil {
			translated++
		}

		if unit.ProofreadText != nil {
			proofread++
		}
	}

	return total, translated, proofread
}

// `ListByPage` returns cloned page units sorted by `Index`.
func (m *mockUnitRepo) ListByPage(pageId string) ([]*aggr.Unit, repo_iface.RepoErr) {
	return ListMockUnitsByPage(m, pageId), nil
}

func (m *mockUnitRepo) ListIndicesByPage(pageId string) ([]*aggr.UnitIndex, repo_iface.RepoErr) {
	indices := make([]*aggr.UnitIndex, 0)

	for _, u := range m.units {
		if u.PageId == pageId {
			indices = append(indices, &aggr.UnitIndex{
				Id:    u.Id,
				Index: u.Index,
			})
		}
	}

	return indices, nil
}

// `CountByPage` returns total, translated, and proofread counts for one page.
func (m *mockUnitRepo) CountByPage(pageId string) (int, int, int, repo_iface.RepoErr) {
	total, translated, proofread := CountMockUnitsByPage(m, pageId)

	return total, translated, proofread, nil
}

func (m *mockUnitRepo) Create(cre *aggr.UnitCre) repo_iface.RepoErr {
	unit := &aggr.Unit{
		Id:                 cre.LocalId,
		PageId:             cre.PageId,
		IsBubble:           cre.IsBubble,
		IsProofread:        cre.IsProofread,
		XCoord:             cre.XCoord,
		YCoord:             cre.YCoord,
		TranslatedText:     cre.TranslatedText,
		TranslatorComment:  cre.TranslatorComment,
		LastTranslatorId:   cre.LastTranslatorId,
		ProofreadText:      cre.ProofreadText,
		ProofreaderComment: cre.ProofreaderComment,
		LastProofreaderId:  cre.LastProofreaderId,
	}

	m.units[unit.Id] = unit

	return nil
}

func (m *mockUnitRepo) Save(sv *aggr.UnitSave) repo_iface.RepoErr {
	// Here it simply mocks the ON CONFLICT DO UPDATE behavior of SQL by checking if the unit exists in the map.
	unit, ok := m.units[sv.Id]
	if !ok {
		cre := &aggr.UnitCre{
			LocalId:            sv.Id,
			PageId:             sv.PageId,
			IsBubble:           sv.IsBubble,
			IsProofread:        sv.IsProofread,
			XCoord:             sv.XCoord,
			YCoord:             sv.YCoord,
			TranslatedText:     sv.TranslatedText,
			TranslatorComment:  sv.TranslatorComment,
			LastTranslatorId:   sv.LastTranslatorId,
			ProofreadText:      sv.ProofreadText,
			ProofreaderComment: sv.ProofreaderComment,
			LastProofreaderId:  sv.LastProofreaderId,
		}

		return m.Create(cre)
	}

	unit.PageId = sv.PageId
	unit.IsBubble = sv.IsBubble
	unit.IsProofread = sv.IsProofread
	unit.XCoord = sv.XCoord
	unit.YCoord = sv.YCoord
	unit.TranslatedText = sv.TranslatedText
	unit.TranslatorComment = sv.TranslatorComment
	unit.LastTranslatorId = sv.LastTranslatorId
	unit.ProofreadText = sv.ProofreadText
	unit.ProofreaderComment = sv.ProofreaderComment
	unit.LastProofreaderId = sv.LastProofreaderId

	return nil
}

func (m *mockUnitRepo) Reindex(indices []*aggr.UnitIndex) repo_iface.RepoErr {
	for _, idx := range indices {
		if unit, ok := m.units[idx.Id]; ok {
			unit.Index = idx.Index
		}
	}

	return nil
}

func (m *mockUnitRepo) Delete(id string) repo_iface.RepoErr {
	delete(m.units, id)
	return nil
}
