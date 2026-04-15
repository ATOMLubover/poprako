package app

import (
	"context"
	"errors"

	event_handler "poprako-s/internal/app/event_handler"
	"poprako-s/internal/app/val"
	"poprako-s/internal/domain/event"
	"poprako-s/internal/domain/model"
	"poprako-s/internal/domain/repo"
	"poprako-s/internal/domain/service"

	"go.uber.org/zap"
)

type UnitApp interface {
	// List 获取指定页面的翻译单元列表
	List(
		cx context.Context,
		currUserID string,
		pageID string,
	) ([]*val.UnitInfo, error)

	// Save 以 diff 语义保存翻译单元（insert / patch / delete），
	// 并更新所属 Page 和 Chapter 的统计字段
	Save(
		cx context.Context,
		currUserID string,
		args *val.SavePageUnitArgs,
	) error
}

type unitAppImpl struct {
	unitSvc  service.UnitService
	eventBus event.EventBus

	userRepo       repo.UserRepo
	pageRepo       repo.PageRepo
	chapterRepo    repo.ChapterRepo
	assignmentRepo repo.AssignmentRepo
	unitRepo       repo.UnitRepo
	txnMgr         repo.TxnMgr
}

func NewUnitApp(
	unitSvc service.UnitService,
	eventBus event.EventBus,
	userRepo repo.UserRepo,
	pageRepo repo.PageRepo,
	chapterRepo repo.ChapterRepo,
	assignmentRepo repo.AssignmentRepo,
	unitRepo repo.UnitRepo,
	txnMgr repo.TxnMgr,
) UnitApp {
	// 校验构造函数依赖
	if unitSvc == nil ||
		eventBus == nil ||
		userRepo == nil ||
		pageRepo == nil ||
		chapterRepo == nil ||
		assignmentRepo == nil ||
		unitRepo == nil ||
		txnMgr == nil {
		zap.L().Panic(
			"NewUnitApp: 依赖项不能为空",
			zap.Bool("unitSvc_nil", unitSvc == nil),
			zap.Bool("eventBus_nil", eventBus == nil),
			zap.Bool("userRepo_nil", userRepo == nil),
			zap.Bool("pageRepo_nil", pageRepo == nil),
			zap.Bool("chapterRepo_nil", chapterRepo == nil),
			zap.Bool("assignmentRepo_nil", assignmentRepo == nil),
			zap.Bool("unitRepo_nil", unitRepo == nil),
			zap.Bool("txnMgr_nil", txnMgr == nil),
		)
	}

	// 返回真实业务实现
	return &unitAppImpl{
		unitSvc:        unitSvc,
		eventBus:       eventBus,
		userRepo:       userRepo,
		pageRepo:       pageRepo,
		chapterRepo:    chapterRepo,
		assignmentRepo: assignmentRepo,
		unitRepo:       unitRepo,
		txnMgr:         txnMgr,
	}
}

func (a *unitAppImpl) List(
	cx context.Context,
	currUserID string,
	pageID string,
) ([]*val.UnitInfo, error) {
	// 获取上下文中的日志记录器
	lgr := retrieveLgr(cx)

	// 查询页面信息以获取所属章节
	targetPage, err := a.pageRepo.GetByID(pageID)
	if err != nil {
		// 记录查询失败
		lgr.Error(
			"获取翻译单元列表失败：获取页面信息失败",
			zap.String("page_id", pageID),
			zap.Error(err),
		)

		// 返回客户端可展示的错误
		return nil, errors.New("无法获取页面信息")
	}

	// 鉴权：检查当前用户是否为该章节的分配人员
	_, err = a.assignmentRepo.Get(model.AssignmentQueryOpt{
		ChapterID: &targetPage.ChapterID,
		UserID:    &currUserID,
	})
	if err != nil {
		// 记录权限校验失败
		lgr.Warn(
			"获取翻译单元列表失败：权限不足",
			zap.String("curr_user_id", currUserID),
			zap.String("page_id", pageID),
		)

		// 返回客户端可展示的错误
		return nil, errors.New("权限不足")
	}

	// 查询翻译单元列表
	units, err := a.unitRepo.List(model.UnitQueryOpt{
		PageID: pageID,
	})
	if err != nil {
		// 记录查询失败
		lgr.Error(
			"获取翻译单元列表失败",
			zap.String("page_id", pageID),
			zap.Error(err),
		)

		// 返回客户端可展示的错误
		return nil, errors.New("获取翻译单元列表失败")
	}

	// 组装为 app 层值对象列表
	result := make([]*val.UnitInfo, len(units))

	for i, unit := range units {
		result[i] = assembleUnitInfo(&unit)
	}

	// 返回翻译单元列表
	return result, nil
}

func (a *unitAppImpl) Save(
	cx context.Context,
	currUserID string,
	args *val.SavePageUnitArgs,
) error {
	// 获取上下文中的日志记录器
	lgr := retrieveLgr(cx)

	// 获取当前用户信息
	currUser, err := a.userRepo.GetByID(currUserID)
	if err != nil {
		// 记录获取用户信息失败
		lgr.Error(
			"保存翻译单元失败：获取当前用户信息失败",
			zap.String("curr_user_id", currUserID),
			zap.Error(err),
		)

		// 返回客户端可展示的错误
		return errors.New("无法获取用户信息")
	}

	// 通过领域服务将 val 层类型转换为 model 层类型
	insertUnits := make([]model.UnitCreation, len(args.UnitDiff.Insert))

	for i, c := range args.UnitDiff.Insert {
		insertUnits[i] = a.unitSvc.NewCreation(
			c.ID,
			args.PageID,
			c.Index,
			c.XCoord,
			c.YCoord,
			c.IsBubble,
			c.TranslatedText,
			c.TranslatorID,
			c.TranslatorComment,
			c.IsProofread,
			c.ProofreadText,
			c.ProofreaderID,
			c.ProofreaderComment,
		)
	}

	patchUnits := make([]model.UnitPatch, len(args.UnitDiff.Patch))

	for i, patch := range args.UnitDiff.Patch {
		patchUnits[i] = model.UnitPatch{
			ID:                 patch.ID,
			Index:              patch.Index,
			XCoord:             patch.XCoord,
			YCoord:             patch.YCoord,
			IsBubble:           patch.IsBubble,
			TranslatedText:     patch.TranslatedText,
			TranslatorID:       patch.TranslatorID,
			TranslatorComment:  patch.TranslatorComment,
			IsProofread:        patch.IsProofread,
			ProofreadText:      patch.ProofreadText,
			ProofreaderID:      patch.ProofreaderID,
			ProofreaderComment: patch.ProofreaderComment,
		}
	}

	// 通过领域服务校验权限和构建 diff（含角色校验）
	_, err = a.unitSvc.NewDiff(
		a.pageRepo,
		a.assignmentRepo,
		currUser,
		args.PageID,
		insertUnits,
		patchUnits,
		args.UnitDiff.Delete,
	)
	if err != nil {
		// 记录权限校验失败
		lgr.Warn(
			"保存翻译单元失败：权限不足",
			zap.String("curr_user_id", currUserID),
			zap.Error(err),
		)

		// 返回领域服务返回的错误
		return err
	}

	// 获取页面信息以获取所属章节 ID
	pageInfo, err := a.pageRepo.GetByID(args.PageID)
	if err != nil {
		// 记录查询失败
		lgr.Error(
			"保存翻译单元失败：获取页面信息失败",
			zap.String("page_id", args.PageID),
			zap.Error(err),
		)

		// 返回客户端可展示的错误
		return errors.New("页面不存在")
	}

	if err := a.txnMgr.RunInTxn(func(cx context.Context) error {
		// 从事务上下文中构造事务版 Repo
		unitRepoTxn, err := a.unitRepo.FromTxnCx(cx)
		if err != nil {
			return err
		}

		// 仅提取需要对比前后状态的受影响单元
		affectedUnitIDMap := make(map[string]struct{}, len(patchUnits)+len(args.UnitDiff.Delete))

		for _, patchUnit := range patchUnits {
			affectedUnitIDMap[patchUnit.ID] = struct{}{}
		}

		for _, unitID := range args.UnitDiff.Delete {
			affectedUnitIDMap[unitID] = struct{}{}
		}

		existingUnits, err := unitRepoTxn.List(model.UnitQueryOpt{PageID: args.PageID})
		if err != nil {
			lgr.Error(
				"保存翻译单元失败：获取现有翻译单元失败",
				zap.String("page_id", args.PageID),
				zap.Error(err),
			)

			return errors.New("保存翻译单元失败")
		}

		existingUnitByID := make(map[string]model.UnitInfo, len(affectedUnitIDMap))

		for _, unit := range existingUnits {
			if _, ok := affectedUnitIDMap[unit.ID]; ok {
				existingUnitByID[unit.ID] = unit
			}
		}

		totalUnitCountDelta := 0
		translatedUnitCountDelta := 0
		proofreadUnitCountDelta := 0

		for _, insertUnit := range insertUnits {
			totalUnitCountDelta++

			if insertUnit.TranslatedText != nil {
				translatedUnitCountDelta++
			}

			if insertUnit.IsProofread {
				proofreadUnitCountDelta++
			}
		}

		for _, patchUnit := range patchUnits {
			existingUnit, ok := existingUnitByID[patchUnit.ID]
			if !ok {
				continue
			}

			translatedBefore := existingUnit.TranslatedText != nil
			if patchUnit.TranslatedText != nil {
				existingUnit.TranslatedText = *patchUnit.TranslatedText
			}
			translatedAfter := existingUnit.TranslatedText != nil

			if !translatedBefore && translatedAfter {
				translatedUnitCountDelta++
			}
			if translatedBefore && !translatedAfter {
				translatedUnitCountDelta--
			}

			proofreadBefore := existingUnit.IsProofread
			if patchUnit.IsProofread != nil {
				existingUnit.IsProofread = *patchUnit.IsProofread
			}
			proofreadAfter := existingUnit.IsProofread

			if !proofreadBefore && proofreadAfter {
				proofreadUnitCountDelta++
			}
			if proofreadBefore && !proofreadAfter {
				proofreadUnitCountDelta--
			}
		}

		for _, unitID := range args.UnitDiff.Delete {
			existingUnit, ok := existingUnitByID[unitID]
			if !ok {
				continue
			}

			totalUnitCountDelta--

			if existingUnit.TranslatedText != nil {
				translatedUnitCountDelta--
			}

			if existingUnit.IsProofread {
				proofreadUnitCountDelta--
			}
		}

		if len(insertUnits) > 0 {
			insertPtrs := make([]*model.UnitCreation, len(insertUnits))

			for i := range insertUnits {
				insertPtrs[i] = &insertUnits[i]
			}

			if err := unitRepoTxn.CreateBatch(insertPtrs); err != nil {
				lgr.Error(
					"保存翻译单元失败：批量创建失败",
					zap.String("page_id", args.PageID),
					zap.Error(err),
				)

				return errors.New("保存翻译单元失败")
			}
		}

		if len(patchUnits) > 0 {
			patchPtrs := make([]*model.UnitPatch, len(patchUnits))

			for i := range patchUnits {
				patchPtrs[i] = &patchUnits[i]
			}

			if err := unitRepoTxn.PatchBatch(patchPtrs); err != nil {
				lgr.Error(
					"保存翻译单元失败：批量修改失败",
					zap.String("page_id", args.PageID),
					zap.Error(err),
				)

				return errors.New("保存翻译单元失败")
			}
		}

		if len(args.UnitDiff.Delete) > 0 {
			if err := unitRepoTxn.DeleteBatch(args.UnitDiff.Delete); err != nil {
				lgr.Error(
					"保存翻译单元失败：批量删除失败",
					zap.String("page_id", args.PageID),
					zap.Error(err),
				)

				return errors.New("保存翻译单元失败")
			}
		}

		pageRepoTxn, err := a.pageRepo.FromTxnCx(cx)
		if err != nil {
			return err
		}

		chapterRepoTxn, err := a.chapterRepo.FromTxnCx(cx)
		if err != nil {
			return err
		}

		eventCx := event_handler.WithChapterRepoTxn(
			event_handler.WithPageRepoTxn(cx, pageRepoTxn),
			chapterRepoTxn,
		)

		saveEvent := a.unitSvc.NewSaveEvent(
			args.PageID,
			pageInfo.ChapterID,
			len(insertUnits),
			len(patchUnits),
			len(args.UnitDiff.Delete),
			totalUnitCountDelta,
			translatedUnitCountDelta,
			proofreadUnitCountDelta,
		)

		if err := a.eventBus.Pub(eventCx, []event.Event{saveEvent}); err != nil {
			lgr.Error(
				"保存翻译单元失败：发布事件失败",
				zap.String("page_id", args.PageID),
				zap.Error(err),
			)

			return errors.New("保存翻译单元失败")
		}

		return nil
	}); err != nil {
		return err
	}

	// 返回保存成功
	return nil
}

// assembleUnitInfo 将领域层翻译单元信息转换为 app 层值对象
func assembleUnitInfo(info *model.UnitInfo) *val.UnitInfo {
	return &val.UnitInfo{
		ID:                 info.ID,
		PageID:             info.PageID,
		Index:              info.Index,
		XCoord:             info.XCoord,
		YCoord:             info.YCoord,
		IsBubble:           info.IsBubble,
		TranslatedText:     info.TranslatedText,
		TranslatorID:       info.TranslatorID,
		TranslatorComment:  info.TranslatorComment,
		IsProofread:        info.IsProofread,
		ProofreadText:      info.ProofreadText,
		ProofreaderID:      info.ProofreaderID,
		ProofreaderComment: info.ProofreaderComment,
	}
}

// logUnitAppImpl 是 UnitApp 的日志包装实现
type logUnitAppImpl struct {
	app UnitApp
}

func NewLogUnitApp(
	app UnitApp,
) UnitApp {
	if app == nil {
		zap.L().Panic(
			"NewLogUnitApp: 依赖项不能为空",
			zap.Bool("app_nil", app == nil),
		)
	}

	return &logUnitAppImpl{app: app}
}

func (a *logUnitAppImpl) List(
	cx context.Context,
	currUserID string,
	pageID string,
) ([]*val.UnitInfo, error) {
	if a == nil || a.app == nil {
		return nil, errors.New("UnitApp 不可用")
	}

	if pageID == "" {
		return nil, errors.New("页面 ID 不能为空")
	}

	lgr := retrieveLgr(cx).With(zap.String("method", "List"), zap.String("curr_user_id", currUserID), zap.String("page_id", pageID))

	cx = injectLgr(cx, lgr)

	lgr.Info("[logUnitAppImpl.List] CALL")

	return a.app.List(cx, currUserID, pageID)
}

func (a *logUnitAppImpl) Save(
	cx context.Context,
	currUserID string,
	args *val.SavePageUnitArgs,
) error {
	if a == nil || a.app == nil {
		return errors.New("UnitApp 不可用")
	}

	if args == nil || args.PageID == "" {
		return errors.New("参数不合法")
	}

	lgr := retrieveLgr(cx).With(zap.String("method", "Save"), zap.String("curr_user_id", currUserID), zap.String("page_id", args.PageID))

	cx = injectLgr(cx, lgr)

	lgr.Info("[logUnitAppImpl.Save] CALL")

	return a.app.Save(cx, currUserID, args)
}
