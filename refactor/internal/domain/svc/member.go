package svc

import (
	"poprako-s/internal/domain/model/aggr"
	"poprako-s/pkg/util"
)

// `MemberSvc` provides domain services for member creation
type MemberSvc struct{}

// `NewMemberSvc` creates a new `MemberSvc`
func NewMemberSvc() MemberSvc {
	return MemberSvc{}
}

// `NewMemberCre` builds a `MemberCre` aggregate with a generated id and the given role mask
func (MemberSvc) NewMemberCre(userId string, teamId string, roleMask aggr.RoleMask) *aggr.MemberCre {
	id := util.GenId("member")

	return &aggr.MemberCre{
		Id:       id,
		UserId:   userId,
		TeamId:   teamId,
		RoleMask: roleMask,
	}
}
