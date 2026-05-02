package svc

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// LabelPlus page header and unit header regexp used in import parser.
var (
	lpPageHeaderRegex = regexp.MustCompile(`^>>>>>>>>\[.+\]<<<<<<<<$`)
	lpUnitHeaderRegex = regexp.MustCompile(`^----------------\[(\d+)\]----------------\[(-?[\d.]+),(-?[\d.]+),([12])\]$`)
	lpCommentRegex    = regexp.MustCompile(`^#\[翻校注释\]：(.*)$`)
)

// Import format constants.
const (
	FormatLp  = "lp"
	FormatPrk = "prk"
)

// ChapterImportSvc provides stateless chapter import parsing logic.
type ChapterImportSvc struct{}

// ChapterImportPage is a parsed page ready for import.
type ChapterImportPage struct {
	Units []ChapterImportUnit
}

// ChapterImportUnit is a parsed unit ready for import.
type ChapterImportUnit struct {
	ID        string
	Index     int
	X         float64
	Y         float64
	IsBubble  bool
	MainText  *string

	TranslatedText *string
	ProofreadText  *string
	IsProofread    bool

	TranslatorComment  *string
	ProofreaderComment *string
}

// ParseLabelPlus parses LabelPlus text content into import pages.
func (ChapterImportSvc) ParseLabelPlus(content string) ([]ChapterImportPage, error) {
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
				return nil, fmt.Errorf("LabelPlus 文件格式错误 缺少页面头")
			}

			if currentUnit != nil {
				applyLabelPlusParsedData(currentUnit, mainTextBuffer, commentBuffer)
				currentPage.units = append(currentPage.units, *currentUnit)
				mainTextBuffer = nil
				commentBuffer = nil
			}

			index, err := strconv.Atoi(matches[1])
			if err != nil {
				return nil, fmt.Errorf("LabelPlus 文件格式错误 单元序号非法")
			}

			x, err := strconv.ParseFloat(matches[2], 64)
			if err != nil {
				return nil, fmt.Errorf("LabelPlus 文件格式错误 X 坐标非法")
			}

			y, err := strconv.ParseFloat(matches[3], 64)
			if err != nil {
				return nil, fmt.Errorf("LabelPlus 文件格式错误 Y 坐标非法")
			}

			currentUnit = &lpParsedUnit{index: index, x: x, y: y, isBubble: matches[4] == "1"}

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
		return nil, fmt.Errorf("读取 LabelPlus 内容失败")
	}

	result := make([]ChapterImportPage, len(pages))
	for i := range pages {
		page := pages[i]

		units := make([]ChapterImportUnit, len(page.units))
		for j := range page.units {
			u := page.units[j]

			mainText := normalizeStringToPtr(u.text)

			units[j] = ChapterImportUnit{
				Index:              u.index,
				X:                  normalizeCoord(u.x),
				Y:                  normalizeCoord(u.y),
				IsBubble:           u.isBubble,
				MainText:           mainText,
				TranslatorComment:  u.translatorComment,
				ProofreaderComment: u.proofreaderComment,
			}
		}

		sort.Slice(units, func(a, b int) bool { return units[a].Index < units[b].Index })

		result[i] = ChapterImportPage{Units: units}
	}

	return result, nil
}

// ParsePoprako parses Poprako JSON content into import pages.
func (ChapterImportSvc) ParsePoprako(content string) ([]ChapterImportPage, error) {
	var project poprakoImportProject

	if err := json.Unmarshal([]byte(content), &project); err != nil {
		return nil, fmt.Errorf("Poprako JSON 解析失败")
	}

	if strings.TrimSpace(project.Author) == "" {
		return nil, fmt.Errorf("Poprako JSON 字段 author 不能为空")
	}

	if strings.TrimSpace(project.Title) == "" {
		return nil, fmt.Errorf("Poprako JSON 字段 title 不能为空")
	}

	if project.Pages == nil {
		return nil, fmt.Errorf("Poprako JSON 字段 pages 不能为空")
	}

	result := make([]ChapterImportPage, len(project.Pages))
	for i := range project.Pages {
		page := project.Pages[i]

		if strings.TrimSpace(page.ImageFilename) == "" {
			return nil, fmt.Errorf("Poprako JSON 第 %d 页 image_filename 不能为空", i+1)
		}

		seenIndex := make(map[uint32]struct{}, len(page.Units))
		units := make([]ChapterImportUnit, len(page.Units))

		for j := range page.Units {
			u := page.Units[j]

			if strings.TrimSpace(u.ID) == "" {
				return nil, fmt.Errorf("Poprako JSON 第 %d 页第 %d 个单元 id 不能为空", i+1, j+1)
			}

			if u.IndexInPage < 1 {
				return nil, fmt.Errorf("Poprako JSON 第 %d 页第 %d 个单元 index_in_page 必须大于等于 1", i+1, j+1)
			}

			if !isFiniteNumber(u.X) || !isFiniteNumber(u.Y) {
				return nil, fmt.Errorf("Poprako JSON 第 %d 页第 %d 个单元坐标非法", i+1, j+1)
			}

			if _, exists := seenIndex[u.IndexInPage]; exists {
				return nil, fmt.Errorf("Poprako JSON 第 %d 页出现重复 index_in_page %d", i+1, u.IndexInPage)
			}

			seenIndex[u.IndexInPage] = struct{}{}

			translatorComment, proofreaderComment := splitImportComment(derefString(normalizeOptionalText(u.Comment)))

			units[j] = ChapterImportUnit{
				ID:                 u.ID,
				Index:              int(u.IndexInPage),
				X:                  normalizeCoord(u.X),
				Y:                  normalizeCoord(u.Y),
				IsBubble:           u.IsInbox,
				TranslatedText:     normalizeOptionalText(u.TranslatedText),
				ProofreadText:      normalizeOptionalText(u.ProovedText),
				IsProofread:        u.IsProoved,
				TranslatorComment:  translatorComment,
				ProofreaderComment: proofreaderComment,
			}
		}

		sort.Slice(units, func(a, b int) bool { return units[a].Index < units[b].Index })

		result[i] = ChapterImportPage{Units: units}
	}

	return result, nil
}

// `lpParsedPage` stores parsed LabelPlus page data before normalization.
type lpParsedPage struct {
	units []lpParsedUnit
}

// `lpParsedUnit` stores parsed LabelPlus unit data before normalization.
type lpParsedUnit struct {
	index    int
	x        float64
	y        float64
	isBubble bool
	text     string

	translatorComment  *string
	proofreaderComment *string
}

// `poprakoImportProject` is the JSON root shape of Poprako import file.
type poprakoImportProject struct {
	Author string              `json:"author"`
	Title  string              `json:"title"`
	Pages  []poprakoImportPage `json:"pages"`
}

// `poprakoImportPage` is one page node in Poprako import file.
type poprakoImportPage struct {
	ImageFilename string              `json:"image_filename"`
	Units         []poprakoImportUnit `json:"units"`
}

// `poprakoImportUnit` is one unit node in Poprako import file.
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

// validateLabelPlusHeader validates LabelPlus file header block.
func validateLabelPlusHeader(scanner *bufio.Scanner) error {
	expectedLines := []string{"1,0", "-", "框内", "框外", "-"}

	for i := range expectedLines {
		expected := expectedLines[i]

		if !scanner.Scan() {
			return fmt.Errorf("LabelPlus 头部不完整 第 %d 行缺失", i+1)
		}

		if scanner.Text() != expected {
			return fmt.Errorf("LabelPlus 头部非法 第 %d 行应为 %s", i+1, expected)
		}
	}

	if !scanner.Scan() {
		return fmt.Errorf("LabelPlus 头部不完整 缺少导出备注")
	}

	if !scanner.Scan() {
		return fmt.Errorf("LabelPlus 头部不完整 缺少备注后的空行")
	}

	if scanner.Text() != "" {
		return fmt.Errorf("LabelPlus 头部非法 备注后必须有空行")
	}

	return nil
}

// applyLabelPlusParsedData writes parsed text and comment buffers back to one unit.
func applyLabelPlusParsedData(unit *lpParsedUnit, mainTextBuffer []string, commentBuffer []string) {
	if len(mainTextBuffer) > 0 {
		unit.text = strings.Join(mainTextBuffer, "\n")
	}

	if len(commentBuffer) > 0 {
		commentText := strings.Join(commentBuffer, "\n")
		unit.translatorComment, unit.proofreaderComment = splitImportComment(commentText)
	}
}

// splitImportComment splits merged comment text into translator and proofreader parts.
func splitImportComment(commentText string) (*string, *string) {
	if strings.TrimSpace(commentText) == "" {
		return nil, nil
	}

	lines := strings.Split(commentText, "\n")

	translatorLines := make([]string, 0, len(lines))
	proofreaderLines := make([]string, 0, len(lines))

	var currentTarget *[]string

	for i := range lines {
		line := lines[i]

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

	var translatorComment *string
	var proofreaderComment *string

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

// normalizeOptionalText normalizes optional text by trimming blank-only values.
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

// normalizeStringToPtr normalizes plain string to optional pointer.
func normalizeStringToPtr(s string) *string {
	trimmed := strings.TrimSpace(s)
	if trimmed == "" {
		return nil
	}

	value := s

	return &value
}

// normalizeCoord keeps coordinate normalization in one dedicated hook.
func normalizeCoord(v float64) float64 {
	return v
}

// isFiniteNumber validates whether number is finite.
func isFiniteNumber(v float64) bool {
	return !math.IsNaN(v) && !math.IsInf(v, 0)
}

// derefString safely dereferences optional string.
func derefString(s *string) string {
	if s == nil {
		return ""
	}

	return *s
}
