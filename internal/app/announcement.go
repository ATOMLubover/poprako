package app_iface

import (
	"context"

	app_res "poprako-s/internal/app/res"
	"poprako-s/internal/app/val"
)

// `AnnouncementApp` defines application use-cases for team announcements.
type AnnouncementApp interface {
	// `List` lists team announcements.
	List(cx context.Context, currUid string, args *val.ListAnnouncementArgs) app_res.AppRes[[]val.AnnouncementVal]

	// `Create` creates one team announcement.
	Create(cx context.Context, currUid string, args *val.CreateAnnouncementArgs) app_res.AppRes[val.AnnouncementCreatedRes]
}
