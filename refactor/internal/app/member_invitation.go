package app_iface

import (
	"context"

	app_res "poprako-s/internal/app/res"
	"poprako-s/internal/app/val"
)

// `MemberInvApp` defines application use-cases for team invitation.
type MemberInvApp interface {
	// `List` lists invitations under one team.
	List(cx context.Context, currUid string, args *val.ListMemberInvArgs) app_res.AppRes[[]val.MemberInvVal]

	// `Create` creates one invitation under one team.
	Create(cx context.Context, currUid string, args *val.CreateMemberInvArgs) app_res.AppRes[val.CreateMemberInvRes]

	// `Update` updates invitation role mask by put semantics.
	Update(cx context.Context, currUid string, args *val.MemberInvUpdArgs) app_res.AppRes[app_res.None]

	// `Delete` deletes one invitation by id.
	Delete(cx context.Context, currUid string, invId string) app_res.AppRes[app_res.None]
}
