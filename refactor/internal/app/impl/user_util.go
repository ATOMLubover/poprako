package app_impl

import (
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
