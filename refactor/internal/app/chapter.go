package app_iface

import (
	"context"

	"poprako-s/internal/app/res"
	"poprako-s/internal/app/val"
)

// `ChapterApp` defines application use-cases for chapter.
type ChapterApp interface {
	// `List` returns chapter list under one comic.
	List(cx context.Context, currUid string, args *val.ListChapterArgs) app_res.AppRes[[]val.ChapterVal]

	// `GetById` returns one chapter by id.
	// `currUid` must be a member of the owning team.
	GetById(cx context.Context, currUid string, chapterId string) app_res.AppRes[val.ChapterVal]

	// `GetPinned` returns pinned chapter under one comic.
	GetPinned(cx context.Context, currUid string, comicId string) app_res.AppRes[val.ChapterVal]

	// `Create` creates one chapter under comic.
	Create(cx context.Context, currUid string, args *val.CreateChapterArgs) app_res.AppRes[val.ChapterCreatedRes]

	// `Update` updates mutable chapter fields.
	Update(cx context.Context, currUid string, args *val.ChapterUpdArgs) app_res.AppRes[app_res.None]

	// `Delete` hard-deletes one chapter.
	Delete(
		cx context.Context,
		currUid string,
		chapterId string,
	) app_res.AppRes[app_res.None]
}
