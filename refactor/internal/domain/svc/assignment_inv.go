package svc

import (
	"errors"

	"poprako-s/internal/domain/model/aggr"
	"poprako-s/internal/domain/model/enum"
	"poprako-s/pkg/util"
)

// `AssignmentInvSvc` provides stateless domain services for assignment invitation.
type AssignmentInvSvc struct{}

// `NewAssignmentInvSvc` returns a ready-to-use `AssignmentInvSvc`.
func NewAssignmentInvSvc() AssignmentInvSvc {
	return AssignmentInvSvc{}
}

// `NewAssignmentInvCre` builds one assignment invitation create payload.
func (AssignmentInvSvc) NewAssignmentInvCre(inviterId string, chapterId string, inviteeQid string, roleMask aggr.RoleMask) (*aggr.AssignmentInvCre, error) {
	if roleMask == 0 {
		return nil, errors.New("至少指定一个角色")
	}

	if roleMask.HasAnyRole(enum.RoleAdmin) {
		return nil, errors.New("章节邀请不支持管理员角色")
	}

	return &aggr.AssignmentInvCre{
		Id:         util.GenId("assignment_invitation"),
		ChapterId:  chapterId,
		InviterId:  inviterId,
		InviteeQid: inviteeQid,
		InvCode:    genAssignmentInvCode(),
		RoleMask:   roleMask,
	}, nil
}

func genAssignmentInvCode() string {
	id := util.GenId("assignment_inv_code")

	if len(id) >= 6 {
		return id[len(id)-6:]
	}

	return id
}
