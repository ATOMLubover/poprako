package application

import (
	"errors"

	"labelplus-next-web-be/internal/application/adapter"
	"labelplus-next-web-be/internal/application/assembler"
	"labelplus-next-web-be/internal/domain/model"
	"labelplus-next-web-be/internal/domain/repository"
	repository_infra "labelplus-next-web-be/internal/infrastructure/repository"
	"labelplus-next-web-be/internal/infrastructure/repository/query_option"
	"labelplus-next-web-be/internal/util"
	"labelplus-next-web-be/internal/value"

	"go.uber.org/zap"
)

type UnitApplication interface {
	// ListPageUnits 获取指定页面下的所有 unit，按 index 升序排列。
	ListPageUnits(
		scope util.TraceScope,
		currentUserID string,
		pageID string,
	) ([]value.UnitInfo, error)

	// SavePageUnits 以 diff 语义保存页面 units（insert / patch / delete），
	// 并在同一事务内更新所属 Page 和 Chapter 的统计字段。
	//
	// 注意：repo 层不对未匹配到的 Update / Delete 行报错，满足协作场景的幂等要求。
	SavePageUnits(
		scope util.TraceScope,
		currentUserID string,
		args value.SavePageUnitArgs,
	) error
}

type unitApplication struct {
	memberRepository     repository.MemberRepository
	comicRepository      repository.ComicRepository
	chapterRepository    repository.ChapterRepository
	pageRepository       repository.PageRepository
	assignmentRepository repository.AssignmentRepository
	unitRepository       repository.UnitRepository
}

func NewUnitApplication(
	memberRepository repository.MemberRepository,
	comicRepository repository.ComicRepository,
	chapterRepository repository.ChapterRepository,
	pageRepository repository.PageRepository,
	assignmentRepository repository.AssignmentRepository,
	unitRepository repository.UnitRepository,
) UnitApplication {
	if memberRepository == nil ||
		comicRepository == nil ||
		chapterRepository == nil ||
		pageRepository == nil ||
		assignmentRepository == nil ||
		unitRepository == nil {
		zap.L().Panic(
			"NewUnitApplication: 依赖项不能为空",
			zap.Bool("memberRepository_nil", memberRepository == nil),
			zap.Bool("comicRepository_nil", comicRepository == nil),
			zap.Bool("chapterRepository_nil", chapterRepository == nil),
			zap.Bool("pageRepository_nil", pageRepository == nil),
			zap.Bool("assignmentRepository_nil", assignmentRepository == nil),
			zap.Bool("unitRepository_nil", unitRepository == nil),
		)
	}

	return &unitApplication{
		memberRepository:     memberRepository,
		comicRepository:      comicRepository,
		chapterRepository:    chapterRepository,
		pageRepository:       pageRepository,
		assignmentRepository: assignmentRepository,
		unitRepository:       unitRepository,
	}
}

func (ua *unitApplication) ListPageUnits(
	scope util.TraceScope,
	currentUserID string,
	pageID string,
) ([]value.UnitInfo, error) {
	const fn = "UnitApplication.ListPageUnits"

	scope.
		WithFields(
			zap.String("current_user_id", currentUserID),
			zap.String("page_id", pageID),
		).
		Logger().
		Debug(fn + ": 被调用")

	if !model.PermUnitList().Check(
		currentUserID,
		pageID,
		adapter.HandleLoadPageInfo(ua.pageRepository),
		adapter.HandleLoadChapterInfo(ua.chapterRepository),
		adapter.HandleLoadAssignmentInfo(ua.assignmentRepository),
	) {
		scope.Logger().Warn(fn + ": 权限检查失败")
		return nil, errors.New("权限不足")
	}

	units, err := ua.unitRepository.List(
		nil,
		query_option.UnitQuery().FilterByPageID(pageID),
		query_option.UnitQuery().OrderByIndexAsc(),
	)
	if err != nil {
		scope.Logger().Error(fn+": 获取翻译单元列表失败", zap.Error(err))
		return nil, errors.New("获取翻译单元列表失败")
	}

	result := make([]value.UnitInfo, len(units))
	for i, unit := range units {
		result[i] = assembler.AssembleUnitInfo(unit)
	}

	return result, nil
}

func (ua *unitApplication) SavePageUnits(
	scope util.TraceScope,
	currentUserID string,
	args value.SavePageUnitArgs,
) error {
	const fn = "UnitApplication.SavePageUnits"

	scope.
		WithFields(
			zap.String("current_user_id", currentUserID),
			zap.String("page_id", args.PageID),
		).
		Logger().
		Debug(fn + ": 被调用")

	// 权限检查：译者或校对者才能保存 unit
	if !model.PermUnitSave().Check(
		currentUserID,
		args.PageID,
		adapter.HandleLoadPageInfo(ua.pageRepository),
		adapter.HandleLoadChapterInfo(ua.chapterRepository),
		adapter.HandleLoadAssignmentInfo(ua.assignmentRepository),
	) {
		scope.Logger().Warn(fn + ": 权限检查失败")
		return errors.New("权限不足")
	}

	// 获取 page → chapter 信息，用于事务结束后更新统计字段
	pageInfo, err := ua.pageRepository.Get(
		nil,
		query_option.FilterByID(repository_infra.PageTable, args.PageID),
	)
	if err != nil {
		scope.Logger().Warn(fn+": 获取页面信息失败", zap.Error(err))
		return errors.New("页面不存在")
	}

	// 将 value.UnitCreation 转换为 model.UnitCreation（事务外）
	insertUnits := make([]model.UnitCreation, len(args.UnitDiff.Insert))
	for i, creation := range args.UnitDiff.Insert {
		unitCreation := model.NewUnitCreation(
			creation.ID,
			creation.Index,
			creation.XCoord,
			creation.YCoord,
			creation.IsBubble,
			creation.TranslatedText,
			creation.TranslatorID,
			creation.TranslatorComment,
			creation.IsProofread,
			creation.ProofreadText,
			creation.ProofreaderID,
			creation.ProofreaderComment,
		)
		unitCreation.PageID = args.PageID
		insertUnits[i] = unitCreation
	}

	// 将 value.UnitPatch 转换为 model.UnitPatch（事务外）
	patchUnits := make([]model.UnitPatch, len(args.UnitDiff.Patch))
	for i, patch := range args.UnitDiff.Patch {
		patchUnits[i] = model.NewUnitPatch(
			patch.ID,
			util.NewSomeOption(patch.Index),
			util.NewSomeOption(patch.XCoord),
			util.NewSomeOption(patch.YCoord),
			util.NewSomeOption(patch.IsBubble),
			patch.TranslatedText,
			patch.TranslatorID,
			patch.TranslatorComment,
			patch.IsProofread,
			patch.ProofreadText,
			patch.ProofreaderID,
			patch.ProofreaderComment,
		)
	}

	transactionExecutor := ua.unitRepository.BeginTransaction()

	var transactionErr error

	defer func() {
		if transactionErr != nil {
			if rollbackErr := transactionExecutor.Rollback().Error; rollbackErr != nil {
				scope.Logger().Error(
					fn+": 事务回滚失败",
					zap.Error(transactionErr),
					zap.Error(rollbackErr),
				)
			}
		}
	}()

	// 悲观锁防并发冲突
	transactionErr = ua.pageRepository.LockByID(transactionExecutor, args.PageID)
	if transactionErr != nil {
		scope.Logger().Error(fn+": 页面加锁失败", zap.Error(transactionErr))
		return errors.New("保存翻译单元失败")
	}

	transactionErr = ua.chapterRepository.LockByID(transactionExecutor, pageInfo.ChapterID)
	if transactionErr != nil {
		scope.Logger().Error(fn+": 章节加锁失败", zap.Error(transactionErr))
		return errors.New("保存翻译单元失败")
	}

	transactionErr = ua.unitRepository.LockByPageID(transactionExecutor, args.PageID)
	if transactionErr != nil {
		scope.Logger().Error(fn+": 加锁失败", zap.Error(transactionErr))
		return errors.New("保存翻译单元失败")
	}

	affectedUnitIDMap := make(map[string]struct{})
	for _, patchUnit := range patchUnits {
		affectedUnitIDMap[patchUnit.ID] = struct{}{}
	}
	for _, unitID := range args.UnitDiff.Delete {
		affectedUnitIDMap[unitID] = struct{}{}
	}

	affectedUnitIDs := make([]string, 0, len(affectedUnitIDMap))
	for unitID := range affectedUnitIDMap {
		affectedUnitIDs = append(affectedUnitIDs, unitID)
	}

	existingAffectedUnits := make([]model.UnitInfo, 0)
	if len(affectedUnitIDs) > 0 {
		existingAffectedUnits, transactionErr = ua.unitRepository.List(
			transactionExecutor,
			query_option.UnitQuery().FilterByPageID(args.PageID),
			query_option.FilterByIDs(repository_infra.UnitTable, affectedUnitIDs),
		)
		if transactionErr != nil {
			scope.Logger().Error(fn+": 获取受影响翻译单元失败", zap.Error(transactionErr))
			return errors.New("保存翻译单元失败")
		}
	}

	existingAffectedUnitByID := make(map[string]model.UnitInfo, len(existingAffectedUnits))
	for _, existingAffectedUnit := range existingAffectedUnits {
		existingAffectedUnitByID[existingAffectedUnit.ID] = existingAffectedUnit
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
		existingUnit, ok := existingAffectedUnitByID[patchUnit.ID]
		if !ok {
			continue
		}

		translatedBeforePatch := existingUnit.TranslatedText != nil
		proofreadBeforePatch := existingUnit.IsProofread

		if patchUnit.TranslatedText.State() == util.OptionSome {
			existingUnit.TranslatedText = patchUnit.TranslatedText.Unwrap()
		}

		if patchUnit.IsProofread.State() == util.OptionSome {
			existingUnit.IsProofread = patchUnit.IsProofread.Unwrap()
		}

		translatedAfterPatch := existingUnit.TranslatedText != nil
		proofreadAfterPatch := existingUnit.IsProofread

		if !translatedBeforePatch && translatedAfterPatch {
			translatedUnitCountDelta++
		}
		if translatedBeforePatch && !translatedAfterPatch {
			translatedUnitCountDelta--
		}

		if !proofreadBeforePatch && proofreadAfterPatch {
			proofreadUnitCountDelta++
		}
		if proofreadBeforePatch && !proofreadAfterPatch {
			proofreadUnitCountDelta--
		}
	}

	for _, unitID := range args.UnitDiff.Delete {
		existingUnit, ok := existingAffectedUnitByID[unitID]
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

	transactionErr = ua.unitRepository.CreateBatch(transactionExecutor, insertUnits)
	if transactionErr != nil {
		scope.Logger().Error(fn+": 批量插入翻译单元失败", zap.Error(transactionErr))
		return errors.New("保存翻译单元失败")
	}

	transactionErr = ua.unitRepository.PatchBatch(transactionExecutor, patchUnits)
	if transactionErr != nil {
		scope.Logger().Error(fn+": 批量 patch翻译单元失败", zap.Error(transactionErr))
		return errors.New("保存翻译单元失败")
	}

	transactionErr = ua.unitRepository.DeleteBatch(transactionExecutor, args.UnitDiff.Delete)
	if transactionErr != nil {
		scope.Logger().Error(fn+": 批量删除翻译单元失败", zap.Error(transactionErr))
		return errors.New("保存翻译单元失败")
	}

	pageStats, transactionErr := ua.pageRepository.GetStatsByID(transactionExecutor, args.PageID)
	if transactionErr != nil {
		scope.Logger().Error(fn+": 获取页面统计字段失败", zap.Error(transactionErr))
		return errors.New("保存翻译单元失败")
	}

	updatedPageStats := model.NewPageStats(
		pageStats.PageID,
		pageStats.TotalUnitCount+totalUnitCountDelta,
		pageStats.TranslatedUnitCount+translatedUnitCountDelta,
		pageStats.ProofreadUnitCount+proofreadUnitCountDelta,
	)
	if updatedPageStats.TotalUnitCount < 0 ||
		updatedPageStats.TranslatedUnitCount < 0 ||
		updatedPageStats.ProofreadUnitCount < 0 {
		scope.Logger().Error(
			fn+": 页面统计字段出现负值",
			zap.Any("page_stats", updatedPageStats),
		)
		return errors.New("保存翻译单元失败")
	}

	transactionErr = ua.pageRepository.UpdateStats(transactionExecutor, updatedPageStats)
	if transactionErr != nil {
		scope.Logger().Error(fn+": 更新页面统计字段失败", zap.Error(transactionErr))
		return errors.New("保存翻译单元失败")
	}

	chapterStats, transactionErr := ua.chapterRepository.GetStatsByID(transactionExecutor, pageInfo.ChapterID)
	if transactionErr != nil {
		scope.Logger().Error(fn+": 获取章节统计字段失败", zap.Error(transactionErr))
		return errors.New("保存翻译单元失败")
	}

	updatedChapterStats := model.NewChapterStats(
		chapterStats.ChapterID,
		chapterStats.TotalUnitCount+totalUnitCountDelta,
		chapterStats.TranslatedUnitCount+translatedUnitCountDelta,
		chapterStats.ProofreadUnitCount+proofreadUnitCountDelta,
	)
	if updatedChapterStats.TotalUnitCount < 0 ||
		updatedChapterStats.TranslatedUnitCount < 0 ||
		updatedChapterStats.ProofreadUnitCount < 0 {
		scope.Logger().Error(
			fn+": 章节统计字段出现负值",
			zap.Any("chapter_stats", updatedChapterStats),
		)
		return errors.New("保存翻译单元失败")
	}

	transactionErr = ua.chapterRepository.UpdateStats(
		transactionExecutor,
		updatedChapterStats,
	)
	if transactionErr != nil {
		scope.Logger().Error(fn+": 更新章节统计字段失败", zap.Error(transactionErr))
		return errors.New("保存翻译单元失败")
	}

	if commitErr := transactionExecutor.Commit().Error; commitErr != nil {
		scope.Logger().Error(fn+": 事务提交失败", zap.Error(commitErr))
		return errors.New("保存翻译单元失败")
	}

	return nil
}
