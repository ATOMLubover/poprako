package app

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"regexp"
	"sort"
	"strconv"
	"strings"

	event_handler "poprako-s/internal/app/event_handler"
	"poprako-s/internal/app/val"
	"poprako-s/internal/domain/event"
	"poprako-s/internal/domain/model"
	"poprako-s/internal/domain/repo"
	"poprako-s/internal/domain/service"

	"go.uber.org/zap"
)

var (
	lpPageHeaderRegex = regexp.MustCompile(`^>>>>>>>>\[.+\]<<<<<<<<$`)
	lpUnitHeaderRegex = regexp.MustCompile(`^----------------\[(\d+)\]----------------\[(-?[\d.]+),(-?[\d.]+),([12])\]$`)
	lpCommentRegex    = regexp.MustCompile(`^#\[翻校注释\]：(.*)$`)
)

// chapterImportAppImpl 代表 ChapterImportApp 的真实实现
type chapterImportAppImpl struct {
	unitSvc        service.UnitService
	eventBus       event.EventBus
	assignmentRepo repo.AssignmentRepo
	chapterRepo    repo.ChapterRepo
	pageRepo       repo.PageRepo
	unitRepo       repo.UnitRepo
	txnMgr         repo.TxnMgr
}

// NewChapterImportApp 创建章节导入应用服务
func NewChapterImportApp(
	unitSvc service.UnitService,
	eventBus event.EventBus,
	assignmentRepo repo.AssignmentRepo,
	chapterRepo repo.ChapterRepo,
	pageRepo repo.PageRepo,
	unitRepo repo.UnitRepo,
	txnMgr repo.TxnMgr,
) ChapterImportApp {
	if unitSvc == nil || eventBus == nil || assignmentRepo == nil || chapterRepo == nil || pageRepo == nil || unitRepo == nil || txnMgr == nil {
		zap.L().Panic(
			"NewChapterImportApp: 依赖项不能为空",
			zap.Bool("unitSvc_nil", unitSvc == nil),
			zap.Bool("eventBus_nil", eventBus == nil),
			zap.Bool("assignmentRepo_nil", assignmentRepo == nil),
			zap.Bool("chapterRepo_nil", chapterRepo == nil),
			zap.Bool("pageRepo_nil", pageRepo == nil),
			zap.Bool("unitRepo_nil", unitRepo == nil),
			zap.Bool("txnMgr_nil", txnMgr == nil),
		)
	}

	return &chapterImportAppImpl{
		unitSvc:        unitSvc,
		eventBus:       eventBus,
		assignmentRepo: assignmentRepo,
		chapterRepo:    chapterRepo,
		pageRepo:       pageRepo,
		unitRepo:       unitRepo,
		txnMgr:         txnMgr,
	}
}

type chapterImportPage struct {
	units []chapterImportUnit
}

type chapterImportUnit struct {
	id string

	index int
	x     int
	y     int

	isBubble bool

	mainText       *string
	translatedText *string
	proofreadText  *string
	isProofread    bool

	translatorComment  *string
	proofreaderComment *string
}

func (a *chapterImportAppImpl) ImportChapter(
	cx context.Context,
	currUserID string,
	args *val.ImportChapterArgs,
) (*val.ImportChapterRes, error) {
	lgr := retrieveLgr(cx)

	if args == nil || args.ChapterID == "" || args.Format == "" || args.Content == "" {
		return nil, errors.New("参数不合法")
	}

	assignment, err := a.assignmentRepo.Get(model.AssignmentQueryOpt{
		ChapterID: &args.ChapterID,
		UserID:    &currUserID,
	})
	if err != nil {
		lgr.Warn("导入章节失败 权限不足", zap.String("curr_user_id", currUserID), zap.String("chapter_id", args.ChapterID))

		return nil, errors.New("权限不足")
	}

	isTranslator := assignment.HasAnyRole(model.RoleTranslator)

	isProofreader := assignment.HasAnyRole(model.RoleProofreader)

	if !isTranslator && !isProofreader {
		return nil, errors.New("权限不足")
	}

	format := strings.ToLower(strings.TrimSpace(args.Format))

	var parsedPages []chapterImportPage

	switch format {
	case "lp":
		parsedPages, err = parseLabelPlusImport(args.Content)
	case "prk":
		parsedPages, err = parsePoprakoImport(args.Content)
	default:
		return nil, errors.New("不支持的导入格式")
	}
	if err != nil {
		return nil, err
	}

	chapterID := args.ChapterID

	dbPages, err := a.pageRepo.List(model.PageQueryOpt{ChapterID: &chapterID})
	if err != nil {
		lgr.Error("导入章节失败 获取页面失败", zap.String("chapter_id", chapterID), zap.Error(err))

		return nil, errors.New("无法获取章节页面")
	}

	sort.Slice(dbPages, func(i, j int) bool {
		return dbPages[i].Index < dbPages[j].Index
	})

	if len(parsedPages) != len(dbPages) {
		return nil, fmt.Errorf("页面数量不匹配 文件为 %d 页 章节为 %d 页", len(parsedPages), len(dbPages))
	}

	result := &val.ImportChapterRes{}

	err = a.txnMgr.RunInTxn(func(txCx context.Context) error {
		unitRepoTxn, txErr := a.unitRepo.FromTxnCx(txCx)
		if txErr != nil {
			return txErr
		}

		pageRepoTxn, txErr := a.pageRepo.FromTxnCx(txCx)
		if txErr != nil {
			return txErr
		}

		chapterRepoTxn, txErr := a.chapterRepo.FromTxnCx(txCx)
		if txErr != nil {
			return txErr
		}

		for i := range parsedPages {
			dbPage := dbPages[i]

			parsedPage := parsedPages[i]

			existingUnits, listErr := unitRepoTxn.List(model.UnitQueryOpt{PageID: dbPage.ID})
			if listErr != nil {
				return errors.New("无法获取页面单元")
			}

			existingByID := make(map[string]model.UnitInfo, len(existingUnits))

			existingByIndex := make(map[int]model.UnitInfo, len(existingUnits))

			for _, existed := range existingUnits {
				existingByID[existed.ID] = existed
				existingByIndex[existed.Index] = existed
			}

			upserts := make([]*model.UnitCreation, 0, len(parsedPage.units))

			insertCount := 0

			patchCount := 0

			totalDelta := 0

			translatedDelta := 0

			proofreadDelta := 0

			for _, parsedUnit := range parsedPage.units {
				unitID := parsedUnit.id

				if unitID == "" {
					if existed, ok := existingByIndex[parsedUnit.index]; ok {
						unitID = existed.ID
					} else {
						unitID = service.GenID("unit")
					}
				}

				existed, exists := existingByID[unitID]

				upsertUnit := model.UnitCreation{
					ID:       unitID,
					PageID:   dbPage.ID,
					Index:    parsedUnit.index,
					XCoord:   parsedUnit.x,
					YCoord:   parsedUnit.y,
					IsBubble: parsedUnit.isBubble,
				}

				if exists {
					upsertUnit.TranslatedText = existed.TranslatedText
					upsertUnit.TranslatorID = existed.TranslatorID
					upsertUnit.TranslatorComment = existed.TranslatorComment
					upsertUnit.IsProofread = existed.IsProofread
					upsertUnit.ProofreadText = existed.ProofreadText
					upsertUnit.ProofreaderID = existed.ProofreaderID
					upsertUnit.ProofreaderComment = existed.ProofreaderComment
				}

				if parsedUnit.translatorComment != nil {
					upsertUnit.TranslatorComment = parsedUnit.translatorComment
				}

				if parsedUnit.proofreaderComment != nil {
					upsertUnit.ProofreaderComment = parsedUnit.proofreaderComment
				}

				if format == "lp" {
					if isProofreader {
						upsertUnit.ProofreadText = parsedUnit.mainText

						if parsedUnit.mainText != nil || parsedUnit.isProofread {
							upsertUnit.IsProofread = true
							upsertUnit.ProofreaderID = &currUserID
						}
					} else {
						upsertUnit.TranslatedText = parsedUnit.mainText

						if parsedUnit.mainText != nil {
							upsertUnit.TranslatorID = &currUserID
						}
					}
				} else {
					if parsedUnit.translatedText != nil {
						upsertUnit.TranslatedText = parsedUnit.translatedText
						upsertUnit.TranslatorID = &currUserID
					}

					if isProofreader {
						if parsedUnit.proofreadText != nil {
							upsertUnit.ProofreadText = parsedUnit.proofreadText
							upsertUnit.IsProofread = true
							upsertUnit.ProofreaderID = &currUserID
						} else if parsedUnit.isProofread {
							upsertUnit.IsProofread = true
							upsertUnit.ProofreaderID = &currUserID
						}
					}
				}

				if !exists {
					insertCount++
					totalDelta++

					if upsertUnit.TranslatedText != nil {
						translatedDelta++
					}

					if upsertUnit.IsProofread {
						proofreadDelta++
					}
				} else {
					patchCount++

					translatedBefore := existed.TranslatedText != nil
					translatedAfter := upsertUnit.TranslatedText != nil

					if !translatedBefore && translatedAfter {
						translatedDelta++
					}

					if translatedBefore && !translatedAfter {
						translatedDelta--
					}

					proofreadBefore := existed.IsProofread
					proofreadAfter := upsertUnit.IsProofread

					if !proofreadBefore && proofreadAfter {
						proofreadDelta++
					}

					if proofreadBefore && !proofreadAfter {
						proofreadDelta--
					}
				}

				upserts = append(upserts, &upsertUnit)
			}

			if len(upserts) > 0 {
				if upsertErr := unitRepoTxn.UpsertBatch(upserts); upsertErr != nil {
					return errors.New("写入翻译单元失败")
				}
			}

			eventCx := event_handler.WithChapterRepoTxn(
				event_handler.WithPageRepoTxn(txCx, pageRepoTxn),
				chapterRepoTxn,
			)

			saveEvent := a.unitSvc.NewSaveEvent(
				dbPage.ID,
				dbPage.ChapterID,
				insertCount,
				patchCount,
				0,
				totalDelta,
				translatedDelta,
				proofreadDelta,
			)

			if pubErr := a.eventBus.Pub(eventCx, []event.Event{saveEvent}); pubErr != nil {
				return errors.New("导入后更新统计失败")
			}

			result.ImportedPageCount++
			result.ImportedUnitCount += len(upserts)
		}

		return nil
	})
	if err != nil {
		lgr.Error("导入章节失败", zap.String("chapter_id", args.ChapterID), zap.Error(err))

		return nil, err
	}

	return result, nil
}

type lpParsedPage struct {
	units []lpParsedUnit
}

type lpParsedUnit struct {
	index int
	x     float64
	y     float64

	isBubble bool

	text string

	translatorComment  *string
	proofreaderComment *string
}

func parseLabelPlusImport(content string) ([]chapterImportPage, error) {
	scanner := bufio.NewScanner(strings.NewReader(content))

	if err := validateLabelPlusHeader(scanner); err != nil {
		return nil, err
	}

	var pages []lpParsedPage

	var currentPage *lpParsedPage

	var currentUnit *lpParsedUnit

	var mainTextBuffer []string

	var commentBuffer []string

	for scanner.Scan() {
		line := scanner.Text()

		if lpPageHeaderRegex.MatchString(line) {
			if currentUnit != nil {
				applyLabelPlusParsedData(currentUnit, mainTextBuffer, commentBuffer)

				currentPage.units = append(currentPage.units, *currentUnit)

				currentUnit = nil
				mainTextBuffer = nil
				commentBuffer = nil
			}

			if currentPage != nil {
				pages = append(pages, *currentPage)
			}

			currentPage = &lpParsedPage{}
			continue
		}

		if matches := lpUnitHeaderRegex.FindStringSubmatch(line); matches != nil {
			if currentPage == nil {
				return nil, errors.New("LabelPlus 文件格式错误 缺少页面头")
			}

			if currentUnit != nil {
				applyLabelPlusParsedData(currentUnit, mainTextBuffer, commentBuffer)

				currentPage.units = append(currentPage.units, *currentUnit)
				mainTextBuffer = nil
				commentBuffer = nil
			}

			index, err := strconv.Atoi(matches[1])
			if err != nil {
				return nil, errors.New("LabelPlus 文件格式错误 单元序号非法")
			}

			x, err := strconv.ParseFloat(matches[2], 64)
			if err != nil {
				return nil, errors.New("LabelPlus 文件格式错误 X 坐标非法")
			}

			y, err := strconv.ParseFloat(matches[3], 64)
			if err != nil {
				return nil, errors.New("LabelPlus 文件格式错误 Y 坐标非法")
			}

			currentUnit = &lpParsedUnit{
				index:    index,
				x:        x,
				y:        y,
				isBubble: matches[4] == "1",
			}
			continue
		}

		if matches := lpCommentRegex.FindStringSubmatch(line); matches != nil {
			commentBuffer = append(commentBuffer, matches[1])
			continue
		}

		if currentUnit != nil && line != "" {
			if len(commentBuffer) > 0 {
				commentBuffer = append(commentBuffer, line)
			} else {
				mainTextBuffer = append(mainTextBuffer, line)
			}
		}
	}

	if currentUnit != nil {
		applyLabelPlusParsedData(currentUnit, mainTextBuffer, commentBuffer)

		currentPage.units = append(currentPage.units, *currentUnit)
	}

	if currentPage != nil {
		pages = append(pages, *currentPage)
	}

	if err := scanner.Err(); err != nil {
		return nil, errors.New("读取 LabelPlus 内容失败")
	}

	result := make([]chapterImportPage, len(pages))

	for i, p := range pages {
		units := make([]chapterImportUnit, len(p.units))

		for j, u := range p.units {
			mainText := normalizeStringToPtr(u.text)

			units[j] = chapterImportUnit{
				index:              u.index,
				x:                  normalizeCoord(u.x),
				y:                  normalizeCoord(u.y),
				isBubble:           u.isBubble,
				mainText:           mainText,
				translatorComment:  u.translatorComment,
				proofreaderComment: u.proofreaderComment,
			}
		}

		sort.Slice(units, func(a, b int) bool {
			return units[a].index < units[b].index
		})

		result[i] = chapterImportPage{units: units}
	}

	return result, nil
}

func validateLabelPlusHeader(scanner *bufio.Scanner) error {
	expectedLines := []string{"1,0", "-", "框内", "框外", "-"}

	for i, expected := range expectedLines {
		if !scanner.Scan() {
			return fmt.Errorf("LabelPlus 头部不完整 第 %d 行缺失", i+1)
		}

		if scanner.Text() != expected {
			return fmt.Errorf("LabelPlus 头部非法 第 %d 行应为 %s", i+1, expected)
		}
	}

	if !scanner.Scan() {
		return errors.New("LabelPlus 头部不完整 缺少导出备注")
	}

	if !scanner.Scan() {
		return errors.New("LabelPlus 头部不完整 缺少备注后的空行")
	}

	if scanner.Text() != "" {
		return errors.New("LabelPlus 头部非法 备注后必须有空行")
	}

	return nil
}

func applyLabelPlusParsedData(
	unit *lpParsedUnit,
	mainTextBuffer []string,
	commentBuffer []string,
) {
	if len(mainTextBuffer) > 0 {
		unit.text = strings.Join(mainTextBuffer, "\n")
	}

	if len(commentBuffer) > 0 {
		commentText := strings.Join(commentBuffer, "\n")

		unit.translatorComment, unit.proofreaderComment = splitImportComment(commentText)
	}
}

type poprakoImportProject struct {
	Author string              `json:"author"`
	Title  string              `json:"title"`
	Pages  []poprakoImportPage `json:"pages"`
}

type poprakoImportPage struct {
	ImageFilename string              `json:"image_filename"`
	Units         []poprakoImportUnit `json:"units"`
}

type poprakoImportUnit struct {
	ID             string  `json:"id"`
	X              float64 `json:"x"`
	Y              float64 `json:"y"`
	IndexInPage    uint32  `json:"index_in_page"`
	IsInbox        bool    `json:"is_inbox"`
	TranslatedText *string `json:"translated_text,omitempty"`
	ProovedText    *string `json:"prooved_text,omitempty"`
	IsProoved      bool    `json:"is_prooved"`
	Comment        *string `json:"comment,omitempty"`
	IsLocal        bool    `json:"is_local"`
}

func parsePoprakoImport(content string) ([]chapterImportPage, error) {
	var project poprakoImportProject

	if err := json.Unmarshal([]byte(content), &project); err != nil {
		return nil, errors.New("Poprako JSON 解析失败")
	}

	if strings.TrimSpace(project.Author) == "" {
		return nil, errors.New("Poprako JSON 字段 author 不能为空")
	}

	if strings.TrimSpace(project.Title) == "" {
		return nil, errors.New("Poprako JSON 字段 title 不能为空")
	}

	if project.Pages == nil {
		return nil, errors.New("Poprako JSON 字段 pages 不能为空")
	}

	result := make([]chapterImportPage, len(project.Pages))

	for i, page := range project.Pages {
		if strings.TrimSpace(page.ImageFilename) == "" {
			return nil, fmt.Errorf("Poprako JSON 第 %d 页 image_filename 不能为空", i+1)
		}

		seenIndex := make(map[uint32]struct{}, len(page.Units))

		units := make([]chapterImportUnit, len(page.Units))

		for j, unit := range page.Units {
			if strings.TrimSpace(unit.ID) == "" {
				return nil, fmt.Errorf("Poprako JSON 第 %d 页第 %d 个单元 id 不能为空", i+1, j+1)
			}

			if unit.IndexInPage < 1 {
				return nil, fmt.Errorf("Poprako JSON 第 %d 页第 %d 个单元 index_in_page 必须大于等于 1", i+1, j+1)
			}

			if !isFiniteNumber(unit.X) || !isFiniteNumber(unit.Y) {
				return nil, fmt.Errorf("Poprako JSON 第 %d 页第 %d 个单元坐标非法", i+1, j+1)
			}

			if _, exists := seenIndex[unit.IndexInPage]; exists {
				return nil, fmt.Errorf("Poprako JSON 第 %d 页出现重复 index_in_page %d", i+1, unit.IndexInPage)
			}

			seenIndex[unit.IndexInPage] = struct{}{}

			translatorComment, proofreaderComment := splitImportComment(derefString(normalizeOptionalText(unit.Comment)))

			units[j] = chapterImportUnit{
				id:                 unit.ID,
				index:              int(unit.IndexInPage),
				x:                  normalizeCoord(unit.X),
				y:                  normalizeCoord(unit.Y),
				isBubble:           unit.IsInbox,
				translatedText:     normalizeOptionalText(unit.TranslatedText),
				proofreadText:      normalizeOptionalText(unit.ProovedText),
				isProofread:        unit.IsProoved,
				translatorComment:  translatorComment,
				proofreaderComment: proofreaderComment,
			}
		}

		sort.Slice(units, func(a, b int) bool {
			return units[a].index < units[b].index
		})

		result[i] = chapterImportPage{units: units}
	}

	return result, nil
}

func splitImportComment(commentText string) (*string, *string) {
	var translatorComment *string

	var proofreaderComment *string

	if strings.TrimSpace(commentText) == "" {
		return nil, nil
	}

	lines := strings.Split(commentText, "\n")

	translatorLines := make([]string, 0, len(lines))

	proofreaderLines := make([]string, 0, len(lines))

	var currentTarget *[]string

	for _, line := range lines {
		if strings.HasPrefix(line, "【翻译】") {
			content := strings.TrimPrefix(line, "【翻译】")
			translatorLines = append(translatorLines, content)
			currentTarget = &translatorLines
			continue
		}

		if strings.HasPrefix(line, "【校对】") {
			content := strings.TrimPrefix(line, "【校对】")
			proofreaderLines = append(proofreaderLines, content)
			currentTarget = &proofreaderLines
			continue
		}

		if currentTarget != nil {
			*currentTarget = append(*currentTarget, line)
		} else {
			translatorLines = append(translatorLines, line)
			currentTarget = &translatorLines
		}
	}

	if len(translatorLines) > 0 {
		text := strings.Join(translatorLines, "\n")
		translatorComment = &text
	}

	if len(proofreaderLines) > 0 {
		text := strings.Join(proofreaderLines, "\n")
		proofreaderComment = &text
	}

	return translatorComment, proofreaderComment
}

func normalizeOptionalText(s *string) *string {
	if s == nil {
		return nil
	}

	trimmed := strings.TrimSpace(*s)

	if trimmed == "" {
		return nil
	}

	value := *s

	return &value
}

func normalizeStringToPtr(s string) *string {
	trimmed := strings.TrimSpace(s)

	if trimmed == "" {
		return nil
	}

	value := s

	return &value
}

func normalizeCoord(v float64) int {
	return int(math.Round(v))
}

func isFiniteNumber(v float64) bool {
	return !math.IsNaN(v) && !math.IsInf(v, 0)
}

func derefString(s *string) string {
	if s == nil {
		return ""
	}

	return *s
}

// logChapterImportAppImpl 代表 ChapterImportApp 的日志包装实现
type logChapterImportAppImpl struct {
	app ChapterImportApp
}

// NewLogChapterImportApp 创建日志包装实现
func NewLogChapterImportApp(app ChapterImportApp) ChapterImportApp {
	if app == nil {
		zap.L().Panic(
			"NewLogChapterImportApp: 依赖项不能为空",
			zap.Bool("app_nil", app == nil),
		)
	}

	return &logChapterImportAppImpl{app: app}
}

func (a *logChapterImportAppImpl) ImportChapter(
	cx context.Context,
	currUserID string,
	args *val.ImportChapterArgs,
) (*val.ImportChapterRes, error) {
	if a == nil || a.app == nil {
		return nil, errors.New("ChapterImportApp 不可用")
	}

	if currUserID == "" || args == nil || args.ChapterID == "" || args.Format == "" || args.Content == "" {
		return nil, errors.New("参数不合法")
	}

	lgr := retrieveLgr(cx).With(
		zap.String("method", "ImportChapter"),
		zap.String("curr_user_id", currUserID),
		zap.String("chapter_id", args.ChapterID),
		zap.String("format", args.Format),
	)

	cx = injectLgr(cx, lgr)

	lgr.Info("[logChapterImportAppImpl.ImportChapter] CALL")

	return a.app.ImportChapter(cx, currUserID, args)
}
