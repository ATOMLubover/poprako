package app_iface

import (
	"context"

	"poprako-s/internal/app/res"
	"poprako-s/internal/app/val"
)

// `SysMailApp` defines application use-cases for user system mails.
type SysMailApp interface {
	// `List` returns unread system mails for current user with pagination.
	List(cx context.Context, currUid string, args *val.ListSysMailArgs) res.AppRes[[]val.SysMailVal]

	// `MarkRead` marks one system mail as read for current user.
	MarkRead(cx context.Context, currUid string, id string) res.AppRes[res.None]
}
