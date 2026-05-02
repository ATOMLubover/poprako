package app_iface

import (
	"context"

	app_res "poprako-s/internal/app/res"
	"poprako-s/internal/app/val"
)

// `ChapterPortApp` defines chapter import and export use-cases.
type ChapterPortApp interface {
	// `Export` returns chapter export payload in JSON-safe object format.
	Export(cx context.Context, currUid string, chapterId string) app_res.AppRes[val.ChapterExportVal]

	// `ExportLp` returns chapter export payload in LabelPlus text format.
	ExportLp(cx context.Context, currUid string, chapterId string) app_res.AppRes[string]

	// `Import` imports chapter text data from one supported format.
	Import(cx context.Context, currUid string, args *val.ImportChapterArgs) app_res.AppRes[val.ImportChapterRes]
}
