package app

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"

	"poprako-s/internal/app/val"
	"poprako-s/internal/domain/ext/oss"
	"poprako-s/internal/domain/model"
	"poprako-s/internal/domain/repo"

	"go.uber.org/zap"
)

// chapterExportAppImpl 代表 ChapterExportApp 的真实业务实现
type chapterExportAppImpl struct {
	chapterRepo    repo.ChapterRepo
	comicRepo      repo.ComicRepo
	pageRepo       repo.PageRepo
	unitRepo       repo.UnitRepo
	assignmentRepo repo.AssignmentRepo
	urlSigner      oss.URLSigner
}

// NewChapterExportApp 返回 ChapterExportApp 的真实业务实现
func NewChapterExportApp(
	chapterRepo repo.ChapterRepo,
	comicRepo repo.ComicRepo,
	pageRepo repo.PageRepo,
	unitRepo repo.UnitRepo,
	assignmentRepo repo.AssignmentRepo,
	urlSigner oss.URLSigner,
) ChapterExportApp {
	// 校验构造函数依赖
	if chapterRepo == nil || comicRepo == nil || pageRepo == nil || unitRepo == nil || assignmentRepo == nil || urlSigner == nil {
		zap.L().Panic(
			"NewChapterExportApp: 依赖项不能为空",
			zap.Bool("chapterRepo_nil", chapterRepo == nil),
			zap.Bool("comicRepo_nil", comicRepo == nil),
			zap.Bool("pageRepo_nil", pageRepo == nil),
			zap.Bool("unitRepo_nil", unitRepo == nil),
			zap.Bool("assignmentRepo_nil", assignmentRepo == nil),
			zap.Bool("urlSigner_nil", urlSigner == nil),
		)
	}

	// 返回真实业务实现
	return &chapterExportAppImpl{
		chapterRepo:    chapterRepo,
		comicRepo:      comicRepo,
		pageRepo:       pageRepo,
		unitRepo:       unitRepo,
		assignmentRepo: assignmentRepo,
		urlSigner:      urlSigner,
	}
}

// chapterExportBundle 代表导出所需的全部原始域数据
type chapterExportBundle struct {
	// chapter 表示章节信息
	chapter model.ChapterInfo
	// comic 表示所属漫画信息
	comic model.ComicInfo
	// pages 表示按序号排列的页面列表
	pages []model.PageInfo
	// units 表示按页面 ID 索引的翻译单元列表
	units map[string][]model.UnitInfo
}

// fetchBundle 加载指定章节的全部导出原始数据
func (a *chapterExportAppImpl) fetchBundle(
	cx context.Context,
	chapterID string,
) (*chapterExportBundle, error) {
	lgr := retrieveLgr(cx)

	// 获取章节信息
	chapter, err := a.chapterRepo.GetByID(chapterID)
	if err != nil {
		lgr.Error("获取章节信息失败", zap.String("chapterID", chapterID), zap.Error(err))

		return nil, err
	}

	// 获取所属漫画信息
	comic, err := a.comicRepo.GetByID(chapter.ComicID)
	if err != nil {
		lgr.Error("获取漫画信息失败", zap.String("comicID", chapter.ComicID), zap.Error(err))

		return nil, err
	}

	// 获取页面列表并按序号排列
	pages, err := a.pageRepo.List(model.PageQueryOpt{ChapterID: &chapterID})
	if err != nil {
		lgr.Error("获取页面列表失败", zap.String("chapterID", chapterID), zap.Error(err))

		return nil, err
	}

	sort.Slice(pages, func(i, j int) bool {
		return pages[i].Index < pages[j].Index
	})

	// 获取各页面的翻译单元
	unitsByPage := make(map[string][]model.UnitInfo, len(pages))

	for _, page := range pages {
		units, unitErr := a.unitRepo.List(model.UnitQueryOpt{PageID: page.ID})
		if unitErr != nil {
			lgr.Error("获取翻译单元列表失败", zap.String("pageID", page.ID), zap.Error(unitErr))

			return nil, unitErr
		}

		unitsByPage[page.ID] = units
	}

	return &chapterExportBundle{
		chapter: *chapter,
		comic:   *comic,
		pages:   pages,
		units:   unitsByPage,
	}, nil
}

func (a *chapterExportAppImpl) ExportChapter(
	cx context.Context,
	currUserID string,
	chapterID string,
) (*val.ChapterExport, error) {
	lgr := retrieveLgr(cx)

	// 校验当前用户在该章节中存在分配记录
	hasAssignment, err := a.assignmentRepo.Exist(model.AssignmentQueryOpt{
		ChapterID: &chapterID,
		UserID:    &currUserID,
	})
	if err != nil {
		lgr.Error("检查分配记录失败", zap.String("chapterID", chapterID), zap.String("currUserID", currUserID), zap.Error(err))

		return nil, errors.New("权限校验失败")
	}

	if !hasAssignment {
		return nil, errors.New("权限不足")
	}

	// 加载导出所需的全部域数据
	bundle, err := a.fetchBundle(cx, chapterID)
	if err != nil {
		return nil, err
	}

	// 逐页组装导出数据
	pageExports := make([]val.PageExport, 0, len(bundle.pages))

	for _, page := range bundle.pages {
		// 若页面已上传则生成图片预签名 URL
		imageURL := ""

		if page.IsUploaded && page.OSSKey != "" {
			url, urlErr := a.urlSigner.GenerateGetPresignedURL(page.OSSKey)
			if urlErr != nil {
				lgr.Warn("生成图片预签名 URL 失败", zap.String("pageID", page.ID), zap.Error(urlErr))
			} else {
				imageURL = url
			}
		}

		// 将域模型翻译单元转换为导出 VO
		rawUnits := bundle.units[page.ID]

		unitExports := make([]val.UnitExport, 0, len(rawUnits))

		for _, unit := range rawUnits {
			unitExports = append(unitExports, val.UnitExport{
				UnitID:             unit.ID,
				UnitIndex:          unit.Index,
				PageID:             page.ID,
				PageIndex:          page.Index,
				XCoord:             unit.XCoord,
				YCoord:             unit.YCoord,
				IsBubble:           unit.IsBubble,
				TranslatedText:     unit.TranslatedText,
				TranslatorID:       unit.TranslatorID,
				TranslatorComment:  unit.TranslatorComment,
				IsProofread:        unit.IsProofread,
				ProofreadText:      unit.ProofreadText,
				ProofreaderID:      unit.ProofreaderID,
				ProofreaderComment: unit.ProofreaderComment,
			})
		}

		pageExports = append(pageExports, val.PageExport{
			PageID:     page.ID,
			PageIndex:  page.Index,
			ImageURL:   imageURL,
			IsUploaded: page.IsUploaded,
			Units:      unitExports,
		})
	}

	// 副标题为空字符串时不导出
	var subtitle *string

	if bundle.chapter.Subtitle != "" {
		subtitle = &bundle.chapter.Subtitle
	}

	// 组装并返回章节导出 VO
	return &val.ChapterExport{
		ChapterID:       bundle.chapter.ID,
		ChapterIndex:    bundle.chapter.Index,
		ChapterSubtitle: subtitle,
		ComicID:         bundle.comic.ID,
		ComicTitle:      bundle.comic.Title,
		Pages:           pageExports,
	}, nil
}

func (a *chapterExportAppImpl) ExportChapterLp(
	cx context.Context,
	currUserID string,
	chapterID string,
) (string, error) {
	lgr := retrieveLgr(cx)

	// 校验当前用户在该章节中存在分配记录
	hasAssignment, err := a.assignmentRepo.Exist(model.AssignmentQueryOpt{
		ChapterID: &chapterID,
		UserID:    &currUserID,
	})
	if err != nil {
		lgr.Error("检查分配记录失败", zap.String("chapterID", chapterID), zap.String("currUserID", currUserID), zap.Error(err))

		return "", errors.New("权限校验失败")
	}

	if !hasAssignment {
		return "", errors.New("权限不足")
	}

	// 加载导出所需的全部域数据
	bundle, err := a.fetchBundle(cx, chapterID)
	if err != nil {
		lgr.Error("导出 LabelPlus 时加载数据失败", zap.String("chapterID", chapterID), zap.Error(err))

		return "", err
	}

	// 构建 LabelPlus 固定文件头
	var sb strings.Builder

	sb.WriteString("1,0\n")
	sb.WriteString("-\n")
	sb.WriteString("框内\n")
	sb.WriteString("框外\n")
	sb.WriteString("-\n")
	sb.WriteString("Exported by PopRaKo Web\n")

	// 逐页写入 LabelPlus 内容
	for _, page := range bundle.pages {
		imgName := lpImageName(page)

		sb.WriteString("\n\n>>>>>>>>[")
		sb.WriteString(imgName)
		sb.WriteString("]<<<<<<<<\n")

		rawUnits := bundle.units[page.ID]

		for i, unit := range rawUnits {
			// G=1 表示框内（气泡），G=2 表示框外
			g := 2

			if unit.IsBubble {
				g = 1
			}

			// 写入单元定位头
			sb.WriteString(fmt.Sprintf(
				"----------------[%d]----------------[%.4f,%.4f,%d]\n",
				i+1,
				float64(unit.XCoord),
				float64(unit.YCoord),
				g,
			))

			// 写入正文，校对文本优先于翻译文本
			mainText := lpSelectText(unit.ProofreadText, unit.TranslatedText)

			if mainText != "" {
				sb.WriteString(mainText)
				sb.WriteString("\n")
			}

			// 写入翻校备注
			comment := lpFormatComment(unit.TranslatorComment, unit.ProofreaderComment)

			if comment != "" {
				sb.WriteString("\n#[翻校注释]：")
				sb.WriteString(comment)
				sb.WriteString("\n")
			}

			sb.WriteString("\n")
		}
	}

	return sb.String(), nil
}

// lpImageName 从页面 OSSKey 提取 LabelPlus 所需的图片文件名
func lpImageName(page model.PageInfo) string {
	if page.OSSKey == "" {
		return fmt.Sprintf("page_%d.jpg", page.Index)
	}

	// 取 OSSKey 最后一段作为文件名
	parts := strings.Split(page.OSSKey, "/")

	name := parts[len(parts)-1]

	if name == "" {
		return fmt.Sprintf("page_%d.jpg", page.Index)
	}

	return name
}

// lpSelectText 按校对优先返回 LabelPlus 正文
func lpSelectText(proofreadText, translatedText *string) string {
	if proofreadText != nil && *proofreadText != "" {
		return *proofreadText
	}

	if translatedText != nil && *translatedText != "" {
		return *translatedText
	}

	return ""
}

// lpFormatComment 合并翻译者与校对者备注为 LabelPlus 注释格式
func lpFormatComment(translatorComment, proofreaderComment *string) string {
	var parts []string

	if translatorComment != nil && *translatorComment != "" {
		parts = append(parts, "【翻译】"+*translatorComment)
	}

	if proofreaderComment != nil && *proofreaderComment != "" {
		parts = append(parts, "【校对】"+*proofreaderComment)
	}

	return strings.Join(parts, "\n")
}

// logChapterExportAppImpl 代表 ChapterExportApp 的日志包装实现
type logChapterExportAppImpl struct {
	app ChapterExportApp
}

// NewLogChapterExportApp 返回 ChapterExportApp 的日志包装实现
func NewLogChapterExportApp(
	app ChapterExportApp,
) ChapterExportApp {
	// 校验依赖
	if app == nil {
		zap.L().Panic(
			"NewLogChapterExportApp: 依赖项不能为空",
			zap.Bool("app_nil", app == nil),
		)
	}

	// 返回日志包装实现
	return &logChapterExportAppImpl{app: app}
}

func (a *logChapterExportAppImpl) ExportChapter(
	cx context.Context,
	currUserID string,
	chapterID string,
) (*val.ChapterExport, error) {
	// 校验包装器是否可用
	if a == nil || a.app == nil {
		return nil, errors.New("ChapterExportApp 不可用")
	}

	// 校验必填参数
	if currUserID == "" || chapterID == "" {
		return nil, errors.New("参数不合法")
	}

	// 注入日志上下文并转发调用
	lgr := retrieveLgr(cx).With(zap.String("method", "ExportChapter"), zap.String("curr_user_id", currUserID), zap.String("chapter_id", chapterID))

	cx = injectLgr(cx, lgr)

	lgr.Info("[logChapterExportAppImpl.ExportChapter] CALL")

	return a.app.ExportChapter(cx, currUserID, chapterID)
}

func (a *logChapterExportAppImpl) ExportChapterLp(
	cx context.Context,
	currUserID string,
	chapterID string,
) (string, error) {
	// 校验包装器是否可用
	if a == nil || a.app == nil {
		return "", errors.New("ChapterExportApp 不可用")
	}

	// 校验必填参数
	if currUserID == "" || chapterID == "" {
		return "", errors.New("参数不合法")
	}

	// 注入日志上下文并转发调用
	lgr := retrieveLgr(cx).With(zap.String("method", "ExportChapterLp"), zap.String("curr_user_id", currUserID), zap.String("chapter_id", chapterID))

	cx = injectLgr(cx, lgr)

	lgr.Info("[logChapterExportAppImpl.ExportChapterLp] CALL")

	return a.app.ExportChapterLp(cx, currUserID, chapterID)
}
