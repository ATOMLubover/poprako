package app_impl

import (
	"fmt"

	app_res "poprako-s/internal/app/res"
	"poprako-s/internal/app/val"
	"poprako-s/internal/domain/model/aggr"
)

// `asmUnitVal` converts one `Unit` aggregate to app-facing `UnitVal`.
func asmUnitVal(unit *aggr.Unit) val.UnitVal {
	return val.UnitVal{
		Id:                 unit.Id,
		PageId:             unit.PageId,
		Index:              unit.Index,
		IsBubble:           unit.IsBubble,
		IsProofread:        unit.IsProofread,
		XCoord:             unit.XCoord,
		YCoord:             unit.YCoord,
		TranslatedText:     unit.TranslatedText,
		TranslatorComment:  unit.TranslatorComment,
		LastTranslatorId:   unit.LastTranslatorId,
		ProofreadText:      unit.ProofreadText,
		ProofreaderComment: unit.ProofreaderComment,
		LastProofreaderId:  unit.LastProofreaderId,
		CreatedAt:          unit.CreatedAt.UnixMilli(),
		UpdatedAt:          unit.UpdatedAt.UnixMilli(),
	}
}

// `vfyListPageUnitsArgs` validates page unit list arguments.
func vfyListPageUnitsArgs(args *val.ListPageUnitsArgs) app_res.AppRes[app_res.None] {
	if args == nil || args.PageId == "" {
		return app_res.Reject[app_res.None](app_res.BadRequest, "page_id 不能为空")
	}

	return app_res.Accept(&app_res.None{})
}

// `vfySavePageUnitsArgs` validates save arguments and converts them to `UnitDiff`.
func vfySavePageUnitsArgs(args *val.SavePageUnitsArgs) (*aggr.UnitDiff, app_res.AppRes[app_res.None]) {
	if args == nil || args.PageId == "" {
		return nil, app_res.Reject[app_res.None](app_res.BadRequest, "page_id 不能为空")
	}

	if args.Diff == nil {
		return nil, app_res.Reject[app_res.None](app_res.BadRequest, "diff 不能为空")
	}

	if args.Diff.PageId == "" || args.Diff.PageId != args.PageId {
		return nil, app_res.Reject[app_res.None](app_res.BadRequest, "diff.page_id 必须与 path page_id 一致")
	}

	// Convert transport-safe ops into domain ops before starting the transaction.
	diff, msg := asmUnitDiff(args.Diff)
	if msg != "" {
		return nil, app_res.Reject[app_res.None](app_res.BadRequest, msg)
	}

	return diff, app_res.Accept(&app_res.None{})
}

// `asmUnitDiff` converts transport-safe `UnitDiffVal` into domain `UnitDiff`.
func asmUnitDiff(diff *val.UnitDiffVal) (*aggr.UnitDiff, string) {
	ops := make([]aggr.UnitOp, 0, len(diff.Ops))
	for i := range diff.Ops {
		op, msg := asmUnitOp(&diff.Ops[i], diff.PageId)
		if msg != "" {
			return nil, fmt.Sprintf("diff.ops[%d] %s", i, msg)
		}

		ops = append(ops, op)
	}

	return &aggr.UnitDiff{PageId: diff.PageId, Ops: ops, CandOrder: diff.CandOrder}, ""
}

// `asmUnitOp` converts one transport-safe unit op into a domain op.
func asmUnitOp(op *val.UnitOpVal, pageId string) (aggr.UnitOp, string) {
	if op == nil {
		return nil, "不能为空"
	}

	// Treat payloads with `local_id` as create ops.
	if op.LocalId != "" {
		if op.Id != "" {
			return nil, "create op 不得同时包含 id"
		}

		if op.IsBubble == nil || op.IsProofread == nil || op.XCoord == nil || op.YCoord == nil {
			return nil, "create op 缺少完整几何或状态字段"
		}

		return &aggr.UnitCre{
			LocalId:            op.LocalId,
			PageId:             pageId,
			IsBubble:           *op.IsBubble,
			IsProofread:        *op.IsProofread,
			XCoord:             *op.XCoord,
			YCoord:             *op.YCoord,
			TranslatedText:     op.TranslatedText,
			TranslatorComment:  op.TranslatorComment,
			LastTranslatorId:   op.LastTranslatorId,
			ProofreadText:      op.ProofreadText,
			ProofreaderComment: op.ProofreaderComment,
			LastProofreaderId:  op.LastProofreaderId,
		}, ""
	}

	if op.Id == "" {
		return nil, "缺少 id 或 local_id"
	}

	// Treat payloads with mutable fields as save ops, otherwise as delete ops.
	if hasUnitSavePayload(op) {
		if op.IsBubble == nil || op.IsProofread == nil || op.XCoord == nil || op.YCoord == nil {
			return nil, "save op 缺少完整几何或状态字段"
		}

		return &aggr.UnitSave{
			Id:                 op.Id,
			PageId:             pageId,
			IsBubble:           *op.IsBubble,
			IsProofread:        *op.IsProofread,
			XCoord:             *op.XCoord,
			YCoord:             *op.YCoord,
			TranslatedText:     op.TranslatedText,
			TranslatorComment:  op.TranslatorComment,
			LastTranslatorId:   op.LastTranslatorId,
			ProofreadText:      op.ProofreadText,
			ProofreaderComment: op.ProofreaderComment,
			LastProofreaderId:  op.LastProofreaderId,
		}, ""
	}

	return &aggr.UnitDel{Id: op.Id}, ""
}

// `hasUnitSavePayload` returns whether one transport op carries mutable save fields.
func hasUnitSavePayload(op *val.UnitOpVal) bool {
	return op.IsBubble != nil ||
		op.IsProofread != nil ||
		op.XCoord != nil ||
		op.YCoord != nil ||
		op.TranslatedText != nil ||
		op.TranslatorComment != nil ||
		op.LastTranslatorId != nil ||
		op.ProofreadText != nil ||
		op.ProofreaderComment != nil ||
		op.LastProofreaderId != nil
}
