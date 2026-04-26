package svc

import (
	"poprako-s/internal/domain/model/aggr"
	"poprako-s/pkg/util"
)

type MemberSvc struct{}

func NewMemberSvc() MemberSvc {
	return MemberSvc{}
}

func (MemberSvc) NewMemberCre(userId string, teamId string, roleMask aggr.RoleMask) *aggr.MemberCre {
	id := util.GenId("member")

	return &aggr.MemberCre{
		Id:       id,
		UserId:   userId,
		TeamId:   teamId,
		RoleMask: roleMask,
	}
}
