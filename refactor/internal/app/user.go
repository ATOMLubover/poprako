package app_iface

import (
	"context"

	"poprako-s/internal/app/res"
	"poprako-s/internal/app/val"
)

type UserApp interface {
	// `GetInfo` is used to get user info by user id.
	// NOTE: it can be used as GetMyInfo.
	GetInfo(cx context.Context, id string) app_res.AppRes[val.UserVal]

	Login(cx context.Context, args *val.UserLoginArgs) app_res.AppRes[val.UserLoginRes]
	Register(cx context.Context, args *val.UserRegArgs) app_res.AppRes[val.UserRegRes]

	// `Update` updates user profile by put semantics.
	Update(cx context.Context, args *val.UserUpdArgs) app_res.AppRes[app_res.None]

	ResvAvatar(cx context.Context, args *val.ResvUserAvatarArgs) app_res.AppRes[val.ResvUserAvatarRes]
	MarkAvatarUploaded(cx context.Context, currUid string) app_res.AppRes[app_res.None]

	// Delete(cx context.Context, id string) app_res.AppRes[app_res.None]
}
