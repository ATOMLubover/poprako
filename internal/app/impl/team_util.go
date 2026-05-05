package app_impl

import (
	app_res "poprako-s/internal/app/res"
	app_util "poprako-s/internal/app/util"
	"poprako-s/internal/app/val"
	oss_iface "poprako-s/internal/domain/ext/oss"
	"poprako-s/internal/domain/model/aggr"
	"poprako-s/internal/domain/model/query"
	repo_iface "poprako-s/internal/domain/repo"
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

// `vfyListTeamArgs` validates and normalizes list team args.
func vfyListTeamArgs(args *val.ListTeamArgs) app_res.AppRes[app_res.None] {
	if args == nil {
		return app_res.Reject[app_res.None](app_res.BadRequest, "分页参数不能为空")
	}

	if re := app_util.ClampOffsetLimit(args.Offset, &args.Limit); re.IsReject() {
		return re
	}

	return app_res.Accept(&app_res.None{})
}

// `vfyListByUserArgs` validates and normalizes list team by user args.
func vfyListByUserArgs(args *val.ListTeamArgs) app_res.AppRes[app_res.None] {
	if args == nil {
		return app_res.Reject[app_res.None](app_res.BadRequest, "分页参数不能为空")
	}

	if re := app_util.ClampOffsetLimit(args.Offset, &args.Limit); re.IsReject() {
		return re
	}

	return app_res.Accept(&app_res.None{})
}

// `asmTeamVals` assembles a team aggregate list into app value list.
func asmTeamVals(teams []*aggr.Team, signer oss_iface.Signer) ([]val.TeamVal, error) {
	teamsVal := make([]val.TeamVal, len(teams))
	for i := range teams {
		teamVal, err := asmTeamVal(teams[i], signer)
		if err != nil {
			return nil, err
		}

		teamsVal[i] = *teamVal
	}

	return teamsVal, nil
}

// `listUserTeamIds` loads all team ids where current user has membership.
func listUserTeamIds(memberRepo repo_iface.MemberRepo, userId string) ([]string, error) {
	members, err := memberRepo.List(&query.ListMemberOpt{UserId: &userId})
	if err != nil {
		return nil, err
	}

	teamIds := make([]string, 0, len(members))
	seen := make(map[string]struct{}, len(members))
	for i := range members {
		if _, ok := seen[members[i].TeamId]; ok {
			continue
		}

		seen[members[i].TeamId] = struct{}{}
		teamIds = append(teamIds, members[i].TeamId)
	}

	return teamIds, nil
}

// `listTeamsByIds` loads teams by ids with stable order from input ids.
func listTeamsByIds(teamRepo repo_iface.TeamRepo, teamIds []string) ([]*aggr.Team, error) {
	teams := make([]*aggr.Team, 0, len(teamIds))
	for i := range teamIds {
		team, err := teamRepo.GetById(teamIds[i])
		if err != nil {
			return nil, err
		}

		teams = append(teams, team)
	}

	return teams, nil
}
