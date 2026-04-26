package app_iface

import (
	"context"

	"poprako-s/internal/app/res"
	"poprako-s/internal/app/val"
)

type UserApp interface {
	// `GetInfo` is used to get user info by user id.
	// NOTE: it can be used as GetMyInfo.
	GetInfo(cx context.Context, id string) res.AppRes[val.UserVal]

	Login(cx context.Context, args *val.UserLoginArgs) res.AppRes[val.UserLoginRes]
	Reg(cx context.Context, args *val.UserRegArgs) res.AppRes[val.UserRegRes]

	// Update(cx context.Context, args *val.UserUpdateArgs) res.AppRes[res.None]

	ResvAvatar(cx context.Context, args *val.ResvUserAvatarArgs) res.AppRes[val.ResvUserAvatarRes]
	MarkAvatarUploaded(cx context.Context, currUid string) res.AppRes[res.None]

	// Delete(cx context.Context, id string) res.AppRes[res.None]
}
