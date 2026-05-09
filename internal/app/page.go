package app_iface

import (
	"context"

	app_res "poprako-s/internal/app/res"
	"poprako-s/internal/app/val"
)

// `PageApp` defines application use-cases for page.
type PageApp interface {
	// `ResvChapterPages` reserves upload slots for chapter pages.
	ResvChapterPages(cx context.Context, currUid string, args *val.ResvChapterPagesArgs) app_res.AppRes[val.ResvChapterPagesRes]

	// `ResvChapterPage` reserves upload slot for one chapter page,
	// It is only used when **re-uploading** one page, and it will invalidate previous reservation if exists.
	ResvChapterPage(cx context.Context, currUid string, args *val.ResvChapterPageArgs) app_res.AppRes[val.ResvChapterPageRes]

	// `List` returns pages under one chapter.
	List(cx context.Context, currUid string, args *val.ListChapterPageArgs) app_res.AppRes[[]val.PageVal]

	// `MarkImageUploaded` confirms one page image upload.
	MarkImageUploaded(cx context.Context, currUid string, args *val.MarkPageImageUploadedArgs) app_res.AppRes[app_res.None]

	// `DeleteByChapterId` deletes all pages under one chapter.
	DeleteByChapterId(cx context.Context, currUid string, chapterId string) app_res.AppRes[app_res.None]
}
