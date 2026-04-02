package app

import (
	"context"
	"errors"

	"poprako-s/internal/app/val"
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
	unitSvc service.UnitService

	userRepo       repo.UserRepo
	pageRepo       repo.PageRepo
	chapterRepo    repo.ChapterRepo
	assignmentRepo repo.AssignmentRepo
	unitRepo       repo.UnitRepo
}

func NewUnitApp(
	unitSvc service.UnitService,
	userRepo repo.UserRepo,
	pageRepo repo.PageRepo,
	chapterRepo repo.ChapterRepo,
	assignmentRepo repo.AssignmentRepo,
	unitRepo repo.UnitRepo,
) UnitApp {
	// 校验构造函数依赖
	if unitSvc == nil ||
		userRepo == nil ||
		pageRepo == nil ||
		chapterRepo == nil ||
		assignmentRepo == nil ||
		unitRepo == nil {
		zap.L().Panic(
			"NewUnitApp: 依赖项不能为空",
			zap.Bool("unitSvc_nil", unitSvc == nil),
			zap.Bool("userRepo_nil", userRepo == nil),
			zap.Bool("pageRepo_nil", pageRepo == nil),
			zap.Bool("chapterRepo_nil", chapterRepo == nil),
			zap.Bool("assignmentRepo_nil", assignmentRepo == nil),
			zap.Bool("unitRepo_nil", unitRepo == nil),
		)
	}

	// 返回真实业务实现
	return &unitAppImpl{
		unitSvc:        unitSvc,
		userRepo:       userRepo,
		pageRepo:       pageRepo,
		chapterRepo:    chapterRepo,
		assignmentRepo: assignmentRepo,
		unitRepo:       unitRepo,
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

	// 将 val 层类型转换为 model 层类型
	insertUnits := make([]model.UnitCreation, len(args.UnitDiff.Insert))

	for i, creation := range args.UnitDiff.Insert {
		insertUnits[i] = model.UnitCreation{
			ID:                 creation.ID,
			PageID:             args.PageID,
			Index:              creation.Index,
			XCoord:             creation.XCoord,
			YCoord:             creation.YCoord,
			IsBubble:           creation.IsBubble,
			TranslatedText:     creation.TranslatedText,
			TranslatorID:       creation.TranslatorID,
			TranslatorComment:  creation.TranslatorComment,
			IsProofread:        creation.IsProofread,
			ProofreadText:      creation.ProofreadText,
			ProofreaderID:      creation.ProofreaderID,
			ProofreaderComment: creation.ProofreaderComment,
		}
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

	// 计算受影响的现有翻译单元，用于统计 delta
	affectedUnitIDMap := make(map[string]struct{})

	for _, patchUnit := range patchUnits {
		affectedUnitIDMap[patchUnit.ID] = struct{}{}
	}

	for _, unitID := range args.UnitDiff.Delete {
		affectedUnitIDMap[unitID] = struct{}{}
	}

	// 获取当前页面所有翻译单元用于计算 delta
	existingUnits, err := a.unitRepo.List(model.UnitQueryOpt{
		PageID: args.PageID,
	})
	if err != nil {
		// 记录查询失败
		lgr.Error(
			"保存翻译单元失败：获取现有翻译单元失败",
			zap.String("page_id", args.PageID),
			zap.Error(err),
		)

		// 返回客户端可展示的错误
		return errors.New("保存翻译单元失败")
	}

	// 构建 ID 映射用于快速查找
	existingUnitByID := make(map[string]model.UnitInfo, len(existingUnits))

	for _, unit := range existingUnits {
		existingUnitByID[unit.ID] = unit
	}

	// 计算统计 delta
	totalUnitCountDelta := 0
	translatedUnitCountDelta := 0
	proofreadUnitCountDelta := 0

	// 新增操作增加计数
	for _, insertUnit := range insertUnits {
		totalUnitCountDelta++

		if insertUnit.TranslatedText != nil {
			translatedUnitCountDelta++
		}

		if insertUnit.IsProofread {
			proofreadUnitCountDelta++
		}
	}

	// 修改操作计算翻译和校对状态变化
	for _, patchUnit := range patchUnits {
		existingUnit, ok := existingUnitByID[patchUnit.ID]
		if !ok {
			continue
		}

		// 检查翻译状态变化
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

		// 检查校对状态变化
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

	// 删除操作减少计数
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

	// 执行批量操作
	if len(insertUnits) > 0 {
		insertPtrs := make([]*model.UnitCreation, len(insertUnits))

		for i := range insertUnits {
			insertPtrs[i] = &insertUnits[i]
		}

		if err := a.unitRepo.CreateBatch(insertPtrs); err != nil {
			// 记录创建失败
			lgr.Error(
				"保存翻译单元失败：批量创建失败",
				zap.String("page_id", args.PageID),
				zap.Error(err),
			)

			// 返回客户端可展示的错误
			return errors.New("保存翻译单元失败")
		}
	}

	if len(patchUnits) > 0 {
		patchPtrs := make([]*model.UnitPatch, len(patchUnits))

		for i := range patchUnits {
			patchPtrs[i] = &patchUnits[i]
		}

		if err := a.unitRepo.PatchBatch(patchPtrs); err != nil {
			// 记录修改失败
			lgr.Error(
				"保存翻译单元失败：批量修改失败",
				zap.String("page_id", args.PageID),
				zap.Error(err),
			)

			// 返回客户端可展示的错误
			return errors.New("保存翻译单元失败")
		}
	}

	if len(args.UnitDiff.Delete) > 0 {
		if err := a.unitRepo.DeleteBatch(args.UnitDiff.Delete); err != nil {
			// 记录删除失败
			lgr.Error(
				"保存翻译单元失败：批量删除失败",
				zap.String("page_id", args.PageID),
				zap.Error(err),
			)

			// 返回客户端可展示的错误
			return errors.New("保存翻译单元失败")
		}
	}

	// 获取当前页面统计数据并应用 delta
	pageStats, err := a.pageRepo.GetStatsByID(args.PageID)
	if err != nil {
		// 记录获取统计失败
		lgr.Error(
			"保存翻译单元失败：获取页面统计数据失败",
			zap.String("page_id", args.PageID),
			zap.Error(err),
		)

		// 返回客户端可展示的错误
		return errors.New("保存翻译单元失败")
	}

	// 应用 delta 到页面统计
	updatedPageStats := &model.PageStats{
		PageID:              args.PageID,
		TotalUnitCount:      pageStats.TotalUnitCount + totalUnitCountDelta,
		TranslatedUnitCount: pageStats.TranslatedUnitCount + translatedUnitCountDelta,
		ProofreadUnitCount:  pageStats.ProofreadUnitCount + proofreadUnitCountDelta,
	}

	// 校验统计值不为负
	if updatedPageStats.TotalUnitCount < 0 ||
		updatedPageStats.TranslatedUnitCount < 0 ||
		updatedPageStats.ProofreadUnitCount < 0 {
		// 记录异常
		lgr.Error(
			"保存翻译单元失败：页面统计值出现负数",
			zap.Any("page_stats", updatedPageStats),
		)

		// 返回客户端可展示的错误
		return errors.New("保存翻译单元失败")
	}

	// 更新页面统计
	if err := a.pageRepo.UpdateStats(updatedPageStats); err != nil {
		// 记录更新失败
		lgr.Error(
			"保存翻译单元失败：更新页面统计失败",
			zap.String("page_id", args.PageID),
			zap.Error(err),
		)

		// 返回客户端可展示的错误
		return errors.New("保存翻译单元失败")
	}

	// 获取章节信息以获取当前统计数据
	chapterInfo, err := a.chapterRepo.GetByID(pageInfo.ChapterID)
	if err != nil {
		// 记录查询失败
		lgr.Error(
			"保存翻译单元失败：获取章节信息失败",
			zap.String("chapter_id", pageInfo.ChapterID),
			zap.Error(err),
		)

		// 返回客户端可展示的错误
		return errors.New("保存翻译单元失败")
	}

	// 更新章节统计
	chapterStats := &model.ChapterStats{
		ChapterID:           pageInfo.ChapterID,
		TotalUnitCount:      chapterInfo.TotalUnitCount + totalUnitCountDelta,
		TranslatedUnitCount: chapterInfo.TranslatedUnitCount + translatedUnitCountDelta,
		ProofreadUnitCount:  chapterInfo.ProofreadUnitCount + proofreadUnitCountDelta,
	}

	if err := a.chapterRepo.UpdateStats(chapterStats); err != nil {
		// 记录更新失败
		lgr.Error(
			"保存翻译单元失败：更新章节统计失败",
			zap.String("chapter_id", pageInfo.ChapterID),
			zap.Error(err),
		)

		// 返回客户端可展示的错误
		return errors.New("保存翻译单元失败")
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
	lgr *zap.Logger
	app UnitApp
}

func NewLogUnitApp(
	lgr *zap.Logger,
	app UnitApp,
) UnitApp {
	if lgr == nil || app == nil {
		zap.L().Panic(
			"NewLogUnitApp: 依赖项不能为空",
			zap.Bool("lgr_nil", lgr == nil),
			zap.Bool("app_nil", app == nil),
		)
	}

	return &logUnitAppImpl{lgr: lgr, app: app}
}

func (a *logUnitAppImpl) List(
	cx context.Context,
	currUserID string,
	pageID string,
) ([]*val.UnitInfo, error) {
	if a == nil || a.app == nil || a.lgr == nil {
		return nil, errors.New("UnitApp 不可用")
	}

	if pageID == "" {
		return nil, errors.New("页面 ID 不能为空")
	}

	lgr := a.lgr.With(zap.String("method", "List"), zap.String("curr_user_id", currUserID), zap.String("page_id", pageID))

	cx = injectLgr(cx, lgr)

	lgr.Info("[logUnitAppImpl.List] CALL")

	return a.app.List(cx, currUserID, pageID)
}

func (a *logUnitAppImpl) Save(
	cx context.Context,
	currUserID string,
	args *val.SavePageUnitArgs,
) error {
	if a == nil || a.app == nil || a.lgr == nil {
		return errors.New("UnitApp 不可用")
	}

	if args == nil || args.PageID == "" {
		return errors.New("参数不合法")
	}

	lgr := a.lgr.With(zap.String("method", "Save"), zap.String("curr_user_id", currUserID), zap.String("page_id", args.PageID))

	cx = injectLgr(cx, lgr)

	lgr.Info("[logUnitAppImpl.Save] CALL")

	return a.app.Save(cx, currUserID, args)
}
