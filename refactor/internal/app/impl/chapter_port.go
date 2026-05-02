package app_impl

import (
	"context"
	"fmt"
	"strings"

	app_iface "poprako-s/internal/app"
	app_res "poprako-s/internal/app/res"
	app_util "poprako-s/internal/app/util"
	"poprako-s/internal/app/val"
	oss_iface "poprako-s/internal/domain/ext/oss"
	"poprako-s/internal/domain/model/aggr"
	"poprako-s/internal/domain/model/enum"
	repo_iface "poprako-s/internal/domain/repo"
	"poprako-s/internal/domain/svc"
	"poprako-s/pkg/util"

	"go.uber.org/zap"
)

// `chapterPortAppImpl` is default implementation of `ChapterPortApp`.
type chapterPortAppImpl struct {
	txnCtrl repo_iface.TxnCtrl

	unitSvc   svc.UnitSvc
	exportSvc svc.ChapterExportSvc
	importSvc svc.ChapterImportSvc

	chapterRepo    repo_iface.ChapterRepo
	comicRepo      repo_iface.ComicRepo
	pageRepo       repo_iface.PageRepo
	unitRepo       repo_iface.UnitRepo
	assignmentRepo repo_iface.AssignmentRepo

	ossSigner oss_iface.Signer
	errClsf   repo_iface.ErrClsf
}

// `chapterPortExportBundle` groups all loaded aggregates needed by export.
type chapterPortExportBundle struct {
	chapter *aggr.Chapter
	pages   []*aggr.Page
	units   map[string][]*aggr.Unit
}

// `NewChapterPortApp` creates one `ChapterPortApp` implementation.
func NewChapterPortApp(
	txnCtrl repo_iface.TxnCtrl,
	chapterRepo repo_iface.ChapterRepo,
	comicRepo repo_iface.ComicRepo,
	pageRepo repo_iface.PageRepo,
	unitRepo repo_iface.UnitRepo,
	assignmentRepo repo_iface.AssignmentRepo,
	unitSvc svc.UnitSvc,
	exportSvc svc.ChapterExportSvc,
	importSvc svc.ChapterImportSvc,
	ossSigner oss_iface.Signer,
	errClsf repo_iface.ErrClsf,
) app_iface.ChapterPortApp {
	if txnCtrl == nil || chapterRepo == nil || comicRepo == nil || pageRepo == nil || unitRepo == nil || assignmentRepo == nil || ossSigner == nil || errClsf == nil {
		zap.L().Panic(
			"[NewChapterPortApp] nil dependency",
			zap.Bool("txnCtrl", txnCtrl == nil),
			zap.Bool("chapterRepo", chapterRepo == nil),
			zap.Bool("comicRepo", comicRepo == nil),
			zap.Bool("pageRepo", pageRepo == nil),
			zap.Bool("unitRepo", unitRepo == nil),
			zap.Bool("assignmentRepo", assignmentRepo == nil),
			zap.Bool("ossSigner", ossSigner == nil),
			zap.Bool("errClsf", errClsf == nil),
		)
	}

	return &chapterPortAppImpl{
		txnCtrl:        txnCtrl,
		unitSvc:        unitSvc,
		exportSvc:      exportSvc,
		importSvc:      importSvc,
		chapterRepo:    chapterRepo,
		comicRepo:      comicRepo,
		pageRepo:       pageRepo,
		unitRepo:       unitRepo,
		assignmentRepo: assignmentRepo,
		ossSigner:      ossSigner,
		errClsf:        errClsf,
	}
}

// `Export` exports one chapter into JSON-safe value object.
func (a *chapterPortAppImpl) Export(cx context.Context, currUid string, chapterId string) app_res.AppRes[val.ChapterExportVal] {
	lgr := app_util.TakeLgr(cx)

	if re := vfyChapterPortExportChapterId(chapterId); re.IsReject() {
		return app_res.Reject[val.ChapterExportVal](re.Code(), re.Msg())
	}

	if re := a.unitSvc.CanListPageUnits(currUid, chapterId, a.assignmentRepo, a.errClsf); re.IsReject() {
		return app_res.Reject[val.ChapterExportVal](app_res.ErrCode(re.Code()), re.Msg())
	}

	bundle, err := a.fetchChapterPortExportBundle(chapterId)
	if err != nil {
		lgr.Error("[chapterPortAppImpl.Export] failed to load chapter export bundle", zap.Error(err))

		return app_res.Reject[val.ChapterExportVal](app_res.ServerError, "导出章节失败")
	}

	exportVal := a.asmChapterPortExportVal(bundle)

	return app_res.Accept(&exportVal)
}

// `ExportLp` exports one chapter into LabelPlus text payload.
func (a *chapterPortAppImpl) ExportLp(cx context.Context, currUid string, chapterId string) app_res.AppRes[string] {
	lgr := app_util.TakeLgr(cx)

	if re := vfyChapterPortExportChapterId(chapterId); re.IsReject() {
		return app_res.Reject[string](re.Code(), re.Msg())
	}

	if re := a.unitSvc.CanListPageUnits(currUid, chapterId, a.assignmentRepo, a.errClsf); re.IsReject() {
		return app_res.Reject[string](app_res.ErrCode(re.Code()), re.Msg())
	}

	bundle, err := a.fetchChapterPortExportBundle(chapterId)
	if err != nil {
		lgr.Error("[chapterPortAppImpl.ExportLp] failed to load chapter export bundle", zap.Error(err))

		return app_res.Reject[string](app_res.ServerError, "导出章节失败")
	}

	exportText := a.exportSvc.MakeLabelPlus(bundle.pages, bundle.units)

	return app_res.Accept(&exportText)
}

// `Import` imports chapter unit text into existing pages.
func (a *chapterPortAppImpl) Import(cx context.Context, currUid string, args *val.ImportChapterArgs) app_res.AppRes[val.ImportChapterRes] {
	lgr := app_util.TakeLgr(cx)

	if re := vfyImportChapterArgs(args); re.IsReject() {
		return app_res.Reject[val.ImportChapterRes](re.Code(), re.Msg())
	}

	// Verify the caller has translator/proofreader assignment on the target chapter.
	assignment, err := a.assignmentRepo.GetByChapterUserId(args.ChapterId, currUid)
	if err != nil {
		if a.errClsf.IsNotFound(err) {
			return app_res.Reject[val.ImportChapterRes](app_res.Forbidden, "仅当前章节的翻译或校对可导入")
		}

		lgr.Error("[chapterPortAppImpl.Import] failed to check assignment", zap.Error(err))

		return app_res.Reject[val.ImportChapterRes](app_res.ServerError, "导入章节失败")
	}

	if assignment == nil || !assignment.HasAnyRole(enum.RoleTranslator, enum.RoleProofreader) {
		return app_res.Reject[val.ImportChapterRes](app_res.Forbidden, "仅当前章节的翻译或校对可导入")
	}

	chapter, err := a.chapterRepo.GetById(args.ChapterId)
	if err != nil {
		return app_res.Reject[val.ImportChapterRes](app_res.BadRequest, "章节不存在")
	}

	format := strings.ToLower(strings.TrimSpace(args.Format))

	var parsedPages []svc.ChapterImportPage

	switch format {
	case svc.FormatLp:
		parsedPages, err = a.importSvc.ParseLabelPlus(args.Content)
	case svc.FormatPrk:
		parsedPages, err = a.importSvc.ParsePoprako(args.Content)
	default:
		return app_res.Reject[val.ImportChapterRes](app_res.BadRequest, "不支持的导入格式")
	}
	if err != nil {
		return app_res.Reject[val.ImportChapterRes](app_res.BadRequest, err.Error())
	}

	// Load all chapter pages explicitly so large chapters remain fully supported.
	pages, err := listAllPages(a.pageRepo, args.ChapterId)
	if err != nil {
		lgr.Error("[chapterPortAppImpl.Import] failed to list chapter pages", zap.Error(err))

		return app_res.Reject[val.ImportChapterRes](app_res.ServerError, "导入章节失败")
	}

	if len(parsedPages) != len(pages) {
		return app_res.Reject[val.ImportChapterRes](app_res.BadRequest, fmt.Sprintf("页面数量不匹配 文件为 %d 页 章节为 %d 页", len(parsedPages), len(pages)))
	}

	re, err := repo_iface.RunWithTxn[app_res.AppRes[val.ImportChapterRes]](a.txnCtrl, func(prov repo_iface.Prov) (app_res.AppRes[val.ImportChapterRes], error) {
		comicRepo := prov.ComicRepo()
		chapterRepo := prov.ChapterRepo()
		pageRepo := prov.PageRepo()
		unitRepo := prov.UnitRepo()

		result := &val.ImportChapterRes{}

		for i := range parsedPages {
			page := pages[i]
			parsedPage := parsedPages[i]

			existingUnits, err := unitRepo.ListByPage(page.Id)
			if err != nil {
				return app_res.Reject[val.ImportChapterRes](app_res.ServerError, "导入章节失败"), err
			}

			existingById := make(map[string]*aggr.Unit, len(existingUnits))
			existingByIndex := make(map[int]*aggr.Unit, len(existingUnits))
			for j := range existingUnits {
				u := existingUnits[j]
				existingById[u.Id] = u
				existingByIndex[u.Index] = u
			}

			ops := make([]aggr.UnitOp, 0, len(parsedPage.Units))
			candOrder := make([]string, 0, len(parsedPage.Units))

			for j := range parsedPage.Units {
				parsedUnit := parsedPage.Units[j]

				unitId := strings.TrimSpace(parsedUnit.ID)
				if unitId == "" {
					if existed := existingByIndex[parsedUnit.Index]; existed != nil {
						unitId = existed.Id
					} else {
						unitId = util.GenId("unit")
					}
				}

				existed := existingById[unitId]

				save := &aggr.UnitSave{
					Id:       unitId,
					PageId:   page.Id,
					IsBubble: parsedUnit.IsBubble,
					XCoord:   parsedUnit.X,
					YCoord:   parsedUnit.Y,
				}

				if existed != nil {
					save.TranslatedText = existed.TranslatedText
					save.TranslatorComment = existed.TranslatorComment
					save.LastTranslatorId = existed.LastTranslatorId
					save.ProofreadText = existed.ProofreadText
					save.IsProofread = existed.IsProofread
					save.ProofreaderComment = existed.ProofreaderComment
					save.LastProofreaderId = existed.LastProofreaderId
				}

				if parsedUnit.TranslatorComment != nil {
					save.TranslatorComment = parsedUnit.TranslatorComment
				}

				if parsedUnit.ProofreaderComment != nil {
					save.ProofreaderComment = parsedUnit.ProofreaderComment
				}

				if format == svc.FormatLp {
					if assignment.HasAnyRole(enum.RoleProofreader) {
						save.ProofreadText = parsedUnit.MainText
						if parsedUnit.MainText != nil || parsedUnit.IsProofread {
							save.IsProofread = true
							save.LastProofreaderId = &currUid
						}
					} else {
						save.TranslatedText = parsedUnit.MainText
						if parsedUnit.MainText != nil {
							save.LastTranslatorId = &currUid
						}
					}
				} else {
					if parsedUnit.TranslatedText != nil {
						save.TranslatedText = parsedUnit.TranslatedText
						save.LastTranslatorId = &currUid
					}

					if assignment.HasAnyRole(enum.RoleProofreader) {
						if parsedUnit.ProofreadText != nil {
							save.ProofreadText = parsedUnit.ProofreadText
							save.IsProofread = true
							save.LastProofreaderId = &currUid
						} else if parsedUnit.IsProofread {
							save.IsProofread = true
							save.LastProofreaderId = &currUid
						}
					}
				}

				ops = append(ops, save)
				candOrder = append(candOrder, unitId)
			}

			diff := &aggr.UnitDiff{PageId: page.Id, Ops: ops, CandOrder: candOrder}
			if re := a.unitSvc.ApplyOps(diff, unitRepo); re.IsReject() {
				return app_res.Reject[val.ImportChapterRes](app_res.ErrCode(re.Code()), re.Msg()), app_res.DefErr()
			}

			newTotal, newTranslated, newProofread, err := unitRepo.CountByPage(page.Id)
			if err != nil {
				return app_res.Reject[val.ImportChapterRes](app_res.ServerError, "导入章节失败"), err
			}

			if err := pageRepo.SetUnitCounts(page.Id, newTotal, newTranslated, newProofread); err != nil {
				return app_res.Reject[val.ImportChapterRes](app_res.ServerError, "导入章节失败"), err
			}

			deltaTotal := newTotal - page.TotalUnitCount
			deltaTranslated := newTranslated - page.TranslatedUnitCount
			deltaProofread := newProofread - page.ProofreadUnitCount

			if err := chapterRepo.AdjustUnitCounts(page.ChapterId, deltaTotal, deltaTranslated, deltaProofread); err != nil {
				return app_res.Reject[val.ImportChapterRes](app_res.ServerError, "导入章节失败"), err
			}

			result.ImportedPageCount++
			result.ImportedUnitCount += len(parsedPage.Units)
		}

		if err := comicRepo.TouchLastActive(chapter.ComicId); err != nil {
			return app_res.Reject[val.ImportChapterRes](app_res.ServerError, "导入章节失败"), err
		}

		return app_res.Accept(result), nil
	})
	if err != nil {
		lgr.Error("[chapterPortAppImpl.Import] failed to import chapter in transaction", zap.Error(err))

		return re
	}

	return re
}

// `fetchChapterPortExportBundle` loads all chapter export dependencies.
func (a *chapterPortAppImpl) fetchChapterPortExportBundle(chapterId string) (*chapterPortExportBundle, error) {
	chapter, err := a.chapterRepo.GetById(chapterId, enum.ChapterInclComic)
	if err != nil {
		return nil, err
	}

	pages, err := listAllPages(a.pageRepo, chapterId)
	if err != nil {
		return nil, err
	}

	unitsByPage := make(map[string][]*aggr.Unit, len(pages))
	for i := range pages {
		units, uerr := a.unitRepo.ListByPage(pages[i].Id)
		if uerr != nil {
			return nil, uerr
		}

		unitsByPage[pages[i].Id] = units
	}

	return &chapterPortExportBundle{chapter: chapter, pages: pages, units: unitsByPage}, nil
}

// `asmChapterPortExportVal` assembles chapter export output object.
func (a *chapterPortAppImpl) asmChapterPortExportVal(bundle *chapterPortExportBundle) val.ChapterExportVal {
	chapter := bundle.chapter

	subtitle := chapter.Subtitle
	if strings.TrimSpace(subtitle) == "" {
		subtitle = ""
	}

	var subtitlePtr *string
	if subtitle != "" {
		subtitlePtr = &subtitle
	}

	pages := make([]val.PageExportVal, 0, len(bundle.pages))
	for i := range bundle.pages {
		page := bundle.pages[i]

		imageUrl := ""
		if page.ImageUploaded && page.ImageKey != nil && *page.ImageKey != "" {
			url, err := a.ossSigner.GenGetUrl(*page.ImageKey)
			if err == nil {
				imageUrl = url
			}
		}

		rawUnits := bundle.units[page.Id]
		unitVals := make([]val.UnitExportVal, 0, len(rawUnits))
		for j := range rawUnits {
			u := rawUnits[j]
			unitVals = append(unitVals, val.UnitExportVal{
				UnitId:             u.Id,
				UnitIndex:          u.Index,
				PageId:             page.Id,
				PageIndex:          page.Index,
				XCoord:             u.XCoord,
				YCoord:             u.YCoord,
				IsBubble:           u.IsBubble,
				TranslatedText:     u.TranslatedText,
				TranslatorId:       u.LastTranslatorId,
				TranslatorComment:  u.TranslatorComment,
				IsProofread:        u.IsProofread,
				ProofreadText:      u.ProofreadText,
				ProofreaderId:      u.LastProofreaderId,
				ProofreaderComment: u.ProofreaderComment,
			})
		}

		pages = append(pages, val.PageExportVal{
			PageId:     page.Id,
			PageIndex:  page.Index,
			ImageUrl:   imageUrl,
			IsUploaded: page.ImageUploaded,
			Units:      unitVals,
		})
	}

	comicTitle := ""
	if chapter.Comic != nil {
		comicTitle = chapter.Comic.Title
	}

	return val.ChapterExportVal{
		ChapterId:       chapter.Id,
		ChapterIndex:    chapter.Index,
		ChapterSubtitle: subtitlePtr,
		ComicId:         chapter.ComicId,
		ComicTitle:      comicTitle,
		Pages:           pages,
	}
}

// `vfyChapterPortExportChapterId` validates chapter id input for export methods.
func vfyChapterPortExportChapterId(chapterId string) app_res.AppRes[app_res.None] {
	if strings.TrimSpace(chapterId) == "" {
		return app_res.Reject[app_res.None](app_res.BadRequest, "chapter_id 不能为空")
	}

	return app_res.Accept(&app_res.None{})
}

// `vfyImportChapterArgs` validates and normalizes chapter import arguments.
func vfyImportChapterArgs(args *val.ImportChapterArgs) app_res.AppRes[app_res.None] {
	if args == nil {
		return app_res.Reject[app_res.None](app_res.BadRequest, "导入参数不能为空")
	}

	args.ChapterId = strings.TrimSpace(args.ChapterId)
	args.Format = strings.TrimSpace(args.Format)

	if args.ChapterId == "" || args.Format == "" || args.Content == "" {
		return app_res.Reject[app_res.None](app_res.BadRequest, "chapter_id、format、content 不能为空")
	}

	return app_res.Accept(&app_res.None{})
}
