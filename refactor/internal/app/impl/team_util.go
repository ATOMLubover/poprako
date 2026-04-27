package app_impl

import (
	"poprako-s/internal/app/val"
	oss_iface "poprako-s/internal/domain/ext/oss"
	"poprako-s/internal/domain/model/aggr"
)

// `asmTeamVal` converts `Team` aggregate to app value object.
func asmTeamVal(team *aggr.Team, signer oss_iface.Signer) (*val.TeamVal, error) {
	avatarUrl := ""
	if team.AvatarUploaded && team.AvatarKey != "" {
		url, err := signer.GenGetUrl(team.AvatarKey)
		if err != nil {
			return nil, err
		}

		avatarUrl = url
	}

	return &val.TeamVal{
		Id:             team.Id,
		Name:           team.Name,
		Desc:           team.Desc,
		AvatarUrl:      avatarUrl,
		AvatarUploaded: team.AvatarUploaded,
		CreatedAt:      team.CreatedAt.UnixMilli(),
		UpdatedAt:      team.UpdatedAt.UnixMilli(),
	}, nil
}
