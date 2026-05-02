package app_iface

import (
	"context"

	app_res "poprako-s/internal/app/res"
	"poprako-s/internal/app/val"
)

// `MemberApp` defines application use-cases for team member.
type MemberApp interface {
	// `Create` creates one member under one team.
	Create(cx context.Context, currUid string, args *val.CreateMemberArgs) app_res.AppRes[val.CreateMemberRes]

	// `ListByTeam` lists members under one team.
	ListByTeam(cx context.Context, currUid string, args *val.ListMemberByTeamArgs) app_res.AppRes[[]val.MemberVal]

	// `ListMine` lists memberships of current user.
	ListMine(cx context.Context, currUid string, args *val.ListMyMemberArgs) app_res.AppRes[[]val.MemberVal]

	// `UpdateRole` updates one member role mask by put semantics.
	UpdateRole(cx context.Context, currUid string, args *val.MemberRoleUpdArgs) app_res.AppRes[app_res.None]

	// `Delete` deletes one member by id.
	Delete(cx context.Context, currUid string, memberId string) app_res.AppRes[app_res.None]

	// `JoinTeam` creates one membership by invitation code.
	JoinTeam(cx context.Context, currUid string, args *val.JoinTeamArgs) app_res.AppRes[app_res.None]
}
