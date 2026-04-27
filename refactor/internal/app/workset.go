package app_iface

import (
	"context"

	"poprako-s/internal/app/res"
	"poprako-s/internal/app/val"
)

// `WorksetApp` defines application use-cases for `Workset`.
type WorksetApp interface {
	// `List` returns all active worksets for a team.
	// `currUid` must be a member of the team; otherwise the request is rejected.
	List(cx context.Context, currUid string, args *val.ListWorksetArgs) res.AppRes[[]val.WorksetVal]

	// `Create` creates a new workset inside a team.
	// `currUid` must be an admin of the team.
	Create(cx context.Context, currUid string, args *val.CreateWorksetArgs) res.AppRes[val.WorksetCreatedRes]

	// `Update` updates the name and/or description of an existing workset.
	// `currUid` must be an admin of the workset's owning team.
	Update(cx context.Context, currUid string, args *val.WorksetUpdArgs) res.AppRes[res.None]

	// `Remove` soft-deletes a workset by id.
	// `currUid` must be an admin of the workset's owning team.
	Remove(cx context.Context, currUid string, worksetId string) res.AppRes[res.None]
}
