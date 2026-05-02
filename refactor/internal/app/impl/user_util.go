package app_impl

import (
	app_res "poprako-s/internal/app/res"
	"poprako-s/internal/app/val"
	oss_iface "poprako-s/internal/domain/ext/oss"
	"poprako-s/internal/domain/model/aggr"
)

func asmUserVal(user *aggr.User, signer oss_iface.Signer) (*val.UserVal, error) {
	// avatarUrl is defaulted to empty string.
	avatarUrl := ""
	if user.AvatarUploaded && user.AvatarKey != "" {
		url, err := signer.GenGetUrl(user.AvatarKey)
		if err != nil {
			return nil, err
		}

		avatarUrl = url
	}

	return &val.UserVal{
		Id:             user.Id,
		Nickname:       user.Nickname,
		Qid:            user.Qid,
		AvatarUrl:      avatarUrl,
		AvatarUploaded: user.AvatarUploaded,
		IsSuperAdmin:   user.IsSuperAdmin,
		LastActiveAt:   user.LastActiveAt.UnixMilli(),
		CreatedAt:      user.CreatedAt.UnixMilli(),
		UpdatedAt:      user.UpdatedAt.UnixMilli(),
	}, nil
}

// `vfyUserUpdArgs` validates user put update args.
func vfyUserUpdArgs(args *val.UserUpdArgs) app_res.AppRes[app_res.None] {
	if args == nil {
		return app_res.Reject[app_res.None](app_res.BadRequest, "更新参数不能为空")
	}

	if args.Id == "" || args.Name == "" || args.Qid == "" {
		return app_res.Reject[app_res.None](app_res.BadRequest, "id name qq 不能为空")
	}

	return app_res.Accept(&app_res.None{})
}
