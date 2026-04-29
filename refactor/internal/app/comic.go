package app_iface

import (
	"context"

	"poprako-s/internal/app/res"
	"poprako-s/internal/app/val"
)

// `ComicApp` defines application use-cases for `Comic`
type ComicApp interface {
	// `List` returns active comics under one workset
	// `currUid` must be a member of the owning team
	List(cx context.Context, currUid string, args *val.ListComicArgs) app_res.AppRes[[]val.ComicVal]

	// `Create` creates one new comic in a workset
	// `currUid` must be an admin of the owning team
	Create(cx context.Context, currUid string, args *val.CreateComicArgs) app_res.AppRes[val.ComicCreatedRes]

	// `Update` applies put-style update to one comic
	// `currUid` must be an admin of the owning team
	Update(cx context.Context, currUid string, args *val.ComicUpdArgs) app_res.AppRes[app_res.None]

	// `GetById` returns one comic by id
	// `currUid` must be a member of the owning team
	GetById(cx context.Context, currUid string, comicId string) app_res.AppRes[val.ComicVal]

	// `Remove` soft-deletes one comic by id
	// `currUid` must be an admin of the owning team
	Remove(cx context.Context, currUid string, comicId string) app_res.AppRes[app_res.None]
}
