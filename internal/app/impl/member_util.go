package app_impl

import (
	"strings"

	app_res "poprako-s/internal/app/res"
	app_util "poprako-s/internal/app/util"
	"poprako-s/internal/app/val"
	oss_iface "poprako-s/internal/domain/ext/oss"
	"poprako-s/internal/domain/model/aggr"
	"poprako-s/internal/domain/model/enum"
	"poprako-s/internal/domain/model/query"
	repo_iface "poprako-s/internal/domain/repo"
)

// `asmMemberVal` converts one member aggregate into app value object.
func asmMemberVal(member *aggr.Member, signer oss_iface.Signer) (*val.MemberVal, error) {
	memberVal := &val.MemberVal{
		Id:           member.Id,
		UserId:       member.UserId,
		UserNickname: member.UserNickname,
		TeamId:       member.TeamId,
		RoleMask:     member.ToRoleMask(),
		CreatedAt:    member.CreatedAt.UnixMilli(),
		UpdatedAt:    member.UpdatedAt.UnixMilli(),
	}

	if member.User != nil {
		userVal, err := asmUserVal(member.User, signer)
		if err != nil {
			return nil, err
		}

		memberVal.User = userVal
	}

	if member.Team != nil {
		teamVal, err := asmTeamVal(member.Team, signer)
		if err != nil {
			return nil, err
		}

		memberVal.Team = teamVal
	}

	return memberVal, nil
}

// `mkListMemberOptByTeam` creates member list query option by team.
func mkListMemberOptByTeam(teamId string, userNicknameKeyword string, includes []enum.MemberIncl, offset int, limit int) *query.ListMemberOpt {
	trimmedNicknameKeyword := strings.TrimSpace(userNicknameKeyword)

	var keywordOpt *string
	if trimmedNicknameKeyword != "" {
		keywordOpt = &trimmedNicknameKeyword
	}

	return &query.ListMemberOpt{
		TeamId:              &teamId,
		UserNicknameKeyword: keywordOpt,
		Pagi: query.PagiOpt{
			Offset: offset,
			Limit:  limit,
		},
		Includes: includes,
	}
}

// `mkListMemberOptByUser` creates member list query option by user.
func mkListMemberOptByUser(userId string, includes []enum.MemberIncl, offset int, limit int) *query.ListMemberOpt {
	return &query.ListMemberOpt{
		UserId: &userId,
		Pagi: query.PagiOpt{
			Offset: offset,
			Limit:  limit,
		},
		Includes: includes,
	}
}

// `vfyListMemberArgs` validates and normalizes member list pagination args.
func vfyListMemberArgs(offset int, limit *int) app_res.AppRes[app_res.None] {
	if re := app_util.ClampOffsetLimit(offset, limit); re.IsReject() {
		return re
	}

	return app_res.Accept(&app_res.None{})
}

// `listAllMembersByUser` loads all memberships by one user.
func listAllMembersByUser(memberRepo repo_iface.MemberRepo, userId string) ([]*aggr.Member, error) {
	const pageLimit = 200

	offset := 0
	result := make([]*aggr.Member, 0)

	for {
		batch, err := memberRepo.List(&query.ListMemberOpt{
			UserId: &userId,
			Pagi: query.PagiOpt{
				Offset: offset,
				Limit:  pageLimit,
			},
		})
		if err != nil {
			return nil, err
		}

		result = append(result, batch...)

		if len(batch) < pageLimit {
			break
		}

		offset += len(batch)
	}

	return result, nil
}
