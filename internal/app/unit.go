package app_iface

import (
	"context"

	app_res "poprako-s/internal/app/res"
	"poprako-s/internal/app/val"
)

// `UnitApp` defines application use-cases for page units.
type UnitApp interface {
	// `ListByPage` returns units under one page.
	ListByPage(cx context.Context, currUid string, args *val.ListPageUnitsArgs) app_res.AppRes[val.ListPageUnitsRes]

	// `SaveByPage` applies one page unit diff and synchronizes counters.
	SaveByPage(cx context.Context, currUid string, args *val.SavePageUnitsArgs) app_res.AppRes[val.SavePageUnitsRes]
}
