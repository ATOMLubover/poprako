package svc

import (
	"fmt"
	"path/filepath"
	"strings"

	"poprako-s/internal/domain/model/aggr"
)

// ChapterExportSvc provides stateless chapter export formatting logic.
type ChapterExportSvc struct{}

// MakeLabelPlus converts chapter pages and units to LabelPlus text format.
func (ChapterExportSvc) MakeLabelPlus(pages []*aggr.Page, unitsByPage map[string][]*aggr.Unit) string {
	var sb strings.Builder

	sb.WriteString("1,0\n")
	sb.WriteString("-\n")
	sb.WriteString("框内\n")
	sb.WriteString("框外\n")
	sb.WriteString("-\n")
	sb.WriteString("Exported by PopRaKo Web\n")

	for i := range pages {
		page := pages[i]

		sb.WriteString("\n\n>>>>>>>>[")
		sb.WriteString(lpImageName(page))
		sb.WriteString("]<<<<<<<<\n")

		units := unitsByPage[page.Id]
		for j := range units {
			u := units[j]

			g := 2
			if u.IsBubble {
				g = 1
			}

			sb.WriteString(fmt.Sprintf("----------------[%d]----------------[%.4f,%.4f,%d]\n", j+1, u.XCoord, u.YCoord, g))

			mainText := lpSelectText(u.ProofreadText, u.TranslatedText)
			if mainText != "" {
				sb.WriteString(mainText)
				sb.WriteString("\n")
			}

			comment := lpFormatComment(u.TranslatorComment, u.ProofreaderComment)
			if comment != "" {
				sb.WriteString("\n#[翻校注释]：")
				sb.WriteString(comment)
				sb.WriteString("\n")
			}

			sb.WriteString("\n")
		}
	}

	return sb.String()
}

// lpImageName derives the standardised LabelPlus image relative path from a page.
// Returns `images/{index:03d}.{ext}` so the official LP parser can locate
// image files by zero-padded index.
func lpImageName(page *aggr.Page) string {
	ext := ".jpg"
	if page.ImageKey != nil && *page.ImageKey != "" {
		if e := filepath.Ext(*page.ImageKey); e != "" {
			ext = e
		}
	}

	return fmt.Sprintf("%03d%s", page.Index, ext)
}

// lpSelectText selects export main text by proofread-first order.
func lpSelectText(proofreadText *string, translatedText *string) string {
	if proofreadText != nil && *proofreadText != "" {
		return *proofreadText
	}

	if translatedText != nil && *translatedText != "" {
		return *translatedText
	}

	return ""
}

// lpFormatComment merges translator and proofreader comments in LabelPlus style.
func lpFormatComment(translatorComment *string, proofreaderComment *string) string {
	parts := make([]string, 0, 2)

	if translatorComment != nil && *translatorComment != "" {
		parts = append(parts, "【翻译】"+*translatorComment)
	}

	if proofreaderComment != nil && *proofreaderComment != "" {
		parts = append(parts, "【校对】"+*proofreaderComment)
	}

	return strings.Join(parts, "\n")
}
