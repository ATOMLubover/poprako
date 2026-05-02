package svc

import (
	"sort"

	"poprako-s/internal/domain/model/aggr"
	"poprako-s/internal/domain/model/enum"
	repo_iface "poprako-s/internal/domain/repo"
	svc_res "poprako-s/internal/domain/svc/res"
	"poprako-s/pkg/util"

	"go.uber.org/zap"
)

// `UnitSvc` contains stateless domain logic for unit operations.
type UnitSvc struct{}

// `CanListPageUnits` validates chapter assignment permission for unit listing.
func (UnitSvc) CanListPageUnits(currUid string, chapterId string, assignmentRepo repo_iface.AssignmentRepo, clsf repo_iface.ErrClsf) svc_res.SvcRes {
	// Validate whether caller has any assignment on the chapter.
	assignment, err := assignmentRepo.GetByChapterUserId(chapterId, currUid)
	if err != nil {
		return classifyRepoErr(err, clsf, svc_res.Forbidden, "仅当前章节参与者可查看 unit", "权限校验失败", "权限校验服务暂不可用", "权限校验失败")
	}

	// Reject users who are not assigned on the chapter.
	if assignment == nil {
		return svc_res.Reject(svc_res.Forbidden, "仅当前章节参与者可查看 unit")
	}

	return svc_res.Accept()
}

// `CanEditPageUnits` validates chapter assignment permission for unit editing.
func (UnitSvc) CanEditPageUnits(currUid string, chapterId string, assignmentRepo repo_iface.AssignmentRepo, clsf repo_iface.ErrClsf) svc_res.SvcRes {
	// Validate whether caller has an assignment on the chapter.
	assignment, err := assignmentRepo.GetByChapterUserId(chapterId, currUid)
	if err != nil {
		return classifyRepoErr(err, clsf, svc_res.Forbidden, "仅当前章节的翻译或校对可编辑 unit", "权限校验失败", "权限校验服务暂不可用", "权限校验失败")
	}

	// Only translator or proofreader roles can mutate units.
	if assignment == nil || !assignment.HasAnyRole(enum.RoleTranslator, enum.RoleProofreader) {
		return svc_res.Reject(svc_res.Forbidden, "仅当前章节的翻译或校对可编辑 unit")
	}

	return svc_res.Accept()
}

// `validateCandOrder` checks whether `diff.CandOrder` follows the client-side
// protocol required by `UnitSvc.ApplyOps`.
func validateCandOrder(diff *aggr.UnitDiff) svc_res.SvcRes {
	if diff == nil {
		return svc_res.Reject(svc_res.BadRequest, "unit diff 不能为空")
	}

	// Build the set of ids explicitly listed in `CandOrder` and reject duplicates.
	candSet := make(map[string]struct{}, len(diff.CandOrder))
	for _, id := range diff.CandOrder {
		if id == "" {
			return svc_res.Reject(svc_res.BadRequest, "cand_order 不得包含空 unit id")
		}

		if _, dup := candSet[id]; dup {
			return svc_res.Reject(svc_res.BadRequest, "cand_order 不得包含重复 unit id")
		}

		candSet[id] = struct{}{}
	}

	// Collect ids that must exist in `CandOrder`, and ids that must not.
	requiredSet := make(map[string]struct{})
	deletedSet := make(map[string]struct{})
	for _, rawOp := range diff.Ops {
		switch op := rawOp.(type) {
		case *aggr.UnitCre:
			if op.LocalId == "" {
				return svc_res.Reject(svc_res.BadRequest, "unit create 缺少 local_id")
			}

			requiredSet[op.LocalId] = struct{}{}

		case *aggr.UnitSave:
			if op.Id == "" {
				return svc_res.Reject(svc_res.BadRequest, "unit save 缺少 id")
			}

			requiredSet[op.Id] = struct{}{}

		case *aggr.UnitDel:
			if op.Id == "" {
				return svc_res.Reject(svc_res.BadRequest, "unit delete 缺少 id")
			}

			deletedSet[op.Id] = struct{}{}

		default:
			return svc_res.Reject(svc_res.BadRequest, "unit diff 包含无效操作类型")
		}
	}

	// Reject impossible diffs where one id is both deleted and submitted.
	for id := range deletedSet {
		if _, submitted := requiredSet[id]; submitted {
			return svc_res.Reject(svc_res.BadRequest, "unit diff 不得同时提交 save/create 与 delete 到同一 unit")
		}

		if _, exists := candSet[id]; exists {
			return svc_res.Reject(svc_res.BadRequest, "cand_order 不得包含待删除 unit id")
		}
	}

	// Require every `UnitCre` and `UnitSave` id to appear in `CandOrder`.
	for id := range requiredSet {
		if _, exists := candSet[id]; !exists {
			return svc_res.Reject(svc_res.BadRequest, "cand_order 必须包含所有 save/create unit id")
		}
	}

	return svc_res.Accept()
}

// `ApplyOps` executes all ops in `diff` and then reindexes the page's units to
// match the sender client's suggested ordering (`diff.CandOrder`) as closely as
// possible while preserving units unknown to the sender client near their
// original neighbors.
func (UnitSvc) ApplyOps(diff *aggr.UnitDiff, unitRepo repo_iface.UnitRepo) svc_res.SvcRes {
	// Validate the `CandOrder` protocol before mutating any op payloads.
	if res := validateCandOrder(diff); res.IsReject() {
		return res
	}

	// `localToReal` maps a client-side temporary `LocalId` (sent in `UnitCre`)
	// to the real server-generated id assigned during this call.
	localToReal := make(map[string]string)

	// Apply every op in order.
	for _, op := range diff.Ops {
		switch op := op.(type) {
		case *aggr.UnitCre:
			// Preserve the client's temporary id before overwriting it, so that
			// `CandOrder` entries referencing the temp id can be resolved later.
			clientLocalId := op.LocalId
			op.LocalId = util.GenId("unit")

			localToReal[clientLocalId] = op.LocalId

			if err := unitRepo.Create(op); err != nil {
				zap.L().Error(
					"[UnitSvc] failed to create unit",
					zap.Error(err),
					zap.Any("op", op),
				)

				return svc_res.Reject(svc_res.ServerError, "创建 unit 失败")
			}

		case *aggr.UnitSave:
			if err := unitRepo.Save(op); err != nil {
				zap.L().Error(
					"[UnitSvc] failed to save unit",
					zap.Error(err),
					zap.Any("op", op),
				)

				return svc_res.Reject(svc_res.ServerError, "保存 unit 失败")
			}

		case *aggr.UnitDel:
			if err := unitRepo.Delete(op.Id); err != nil {
				zap.L().Error(
					"[UnitSvc] failed to delete unit",
					zap.Error(err),
					zap.Any("op", op),
				)

				return svc_res.Reject(svc_res.ServerError, "删除 unit 失败")
			}

		default:
			zap.L().Error(
				"[UnitSvc] invalid unit op type",
				zap.Any("op", op),
			)

			return svc_res.Reject(svc_res.BadRequest, "unit diff 包含无效操作类型")
		}
	}

	// Reindex the page's units to match `CandOrder` as closely as possible.
	indices, err := unitRepo.ListIndicesByPage(diff.PageId)
	if err != nil {
		zap.L().Error(
			"[UnitSvc] failed to list indices for reindex",
			zap.Error(err),
			zap.String("pageId", diff.PageId),
		)

		return svc_res.Reject(svc_res.ServerError, "获取 unit 排序信息失败")
	}

	orderedIds := reorderUnits(diff.CandOrder, localToReal, indices)

	newIndices := make([]*aggr.UnitIndex, len(orderedIds))
	for i, id := range orderedIds {
		newIndices[i] = &aggr.UnitIndex{Id: id, Index: i}
	}

	if err := unitRepo.Reindex(newIndices); err != nil {
		zap.L().Error(
			"[UnitSvc] failed to reindex units",
			zap.Error(err),
			zap.String("pageId", diff.PageId),
		)

		return svc_res.Reject(svc_res.ServerError, "重排 unit 失败")
	}

	return svc_res.Accept()
}

// `reorderUnits` computes the final id ordering for a page's units.
//
// Units present in `candOrder` are placed first in that order (after resolving
// any client temp ids via `localToReal`). Remaining units are treated as
// concurrent unknowns and inserted immediately after the anchor slot whose
// position is closest to the unknown unit's current rank among all surviving
// units. This preserves cluster locality while honouring the sender's intent.
func reorderUnits(
	candOrder []string,
	localToReal map[string]string,
	allIndices []*aggr.UnitIndex,
) []string {
	nAll := len(allIndices)
	if nAll == 0 {
		return nil
	}

	// Sort `allIndices` by the persisted `Index` value because repo return order
	// is not part of the contract.
	sortedIndices := make([]*aggr.UnitIndex, len(allIndices))
	copy(sortedIndices, allIndices)
	sort.Slice(sortedIndices, func(i, j int) bool {
		if sortedIndices[i].Index == sortedIndices[j].Index {
			return sortedIndices[i].Id < sortedIndices[j].Id
		}

		return sortedIndices[i].Index < sortedIndices[j].Index
	})

	// Build a set of all surviving unit ids so that stale anchors can be
	// excluded before slot mapping.
	allSet := make(map[string]struct{}, nAll)
	for _, ui := range sortedIndices {
		allSet[ui.Id] = struct{}{}
	}

	// Resolve temp ids in `candOrder` to real server ids.
	resolved := make([]string, 0, len(candOrder))
	anchorSet := make(map[string]struct{}, len(candOrder))
	for _, id := range candOrder {
		resolvedId := id
		if real, ok := localToReal[id]; ok {
			resolvedId = real
		}

		if _, exists := allSet[resolvedId]; !exists {
			continue
		}

		if _, dup := anchorSet[resolvedId]; dup {
			continue
		}

		anchorSet[resolvedId] = struct{}{}
		resolved = append(resolved, resolvedId)
	}

	nCand := len(resolved)

	// `slotExtras` maps each cand slot index to the list of unknown unit ids
	// that should be inserted immediately after that slot.
	// Unknown units are sorted by their current DB index (`sortedIndices` is
	// ordered ascending by index), so relative order inside each slot is kept.
	slotExtras := make(map[int][]string)
	for rank, ui := range sortedIndices {
		if _, inCand := anchorSet[ui.Id]; inCand {
			continue
		}

		// Map the unknown unit's rank to the nearest cand slot.
		slot := 0
		if nCand > 0 {
			slot = rank * nCand / nAll
			if slot >= nCand {
				slot = nCand - 1
			}
		}

		slotExtras[slot] = append(slotExtras[slot], ui.Id)
	}

	// Assemble the final order: for each cand slot, emit the cand id (if it
	// still exists in the DB) followed by any unknowns assigned to that slot.
	final := make([]string, 0, nAll)
	for p, id := range resolved {
		if _, exists := allSet[id]; exists {
			final = append(final, id)
		}

		final = append(final, slotExtras[p]...)
	}

	// Append all surviving ids if no anchors exist at all.
	if nCand == 0 {
		for _, ui := range sortedIndices {
			final = append(final, ui.Id)
		}
	}

	return final
}
