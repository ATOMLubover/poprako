package app_iface

import (
	"context"

	"poprako-s/internal/app/res"
	"poprako-s/internal/app/val"
)

// `ChapterApp` defines application use-cases for chapter.
type ChapterApp interface {
	// `List` returns chapter list under one comic.
	List(cx context.Context, currUid string, args *val.ListChapterArgs) res.AppRes[[]val.ChapterVal]

	// `GetPinned` returns pinned chapter under one comic.
	GetPinned(cx context.Context, currUid string, comicId string) res.AppRes[val.ChapterVal]

	// `Create` creates one chapter under comic.
	Create(cx context.Context, currUid string, args *val.CreateChapterArgs) res.AppRes[val.ChapterCreatedRes]

	// `Update` updates mutable chapter fields.
	Update(cx context.Context, currUid string, args *val.ChapterUpdArgs) res.AppRes[res.None]

	// `Remove` soft-deletes one chapter.
	Remove(cx context.Context, currUid string, chapterId string) res.AppRes[res.None]
}
