package svc

import (
	"errors"

	"poprako-s/internal/domain/model/aggr"
	"poprako-s/internal/domain/model/enum"
	repo_iface "poprako-s/internal/domain/repo"
	svc_res "poprako-s/internal/domain/svc/res"
	"poprako-s/pkg/util"

	"go.uber.org/zap"
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

// `VfyInviteeNotAssigned` validates the invitee has no existing assignment for the target chapter.
func (AssignmentInvSvc) VfyInviteeNotAssigned(inviteeQid string, chapterId string, userRepo repo_iface.UserRepo, assignmentRepo repo_iface.AssignmentRepo, clsf repo_iface.ErrClsf) svc_res.SvcRes {
	// Resolve the invitee qid to a user id.
	inviteeUser, err := userRepo.GetByQid(inviteeQid)
	if err != nil {
		if clsf.IsNotFound(err) {
			return svc_res.Accept()
		}

		zap.L().Error(
			"[AssignmentInvSvc.VfyInviteeNotAssigned] failed to get invitee user",
			zap.String("inviteeQid", inviteeQid),
			zap.String("chapterId", chapterId),
			zap.Error(err),
		)

		return classifyRepoErr(err, clsf, svc_res.BadRequest, "用户不存在", "校验失败", "校验服务暂不可用", "校验失败")
	}

	// Verify the invitee has no existing assignment for this chapter.
	_, err = assignmentRepo.GetByChapterUserId(chapterId, inviteeUser.Id)
	if err == nil {
		return svc_res.Reject(svc_res.Conflict, "该用户已被分配至当前章节")
	}

	if !clsf.IsNotFound(err) {
		zap.L().Error(
			"[AssignmentInvSvc.VfyInviteeNotAssigned] failed to check existing assignment",
			zap.String("inviteeId", inviteeUser.Id),
			zap.String("chapterId", chapterId),
			zap.Error(err),
		)

		return classifyRepoErr(err, clsf, svc_res.BadRequest, "该用户已被分配至当前章节", "校验失败", "校验服务暂不可用", "校验失败")
	}

	return svc_res.Accept()
}

func genAssignmentInvCode() string {
	id := util.GenId("assignment_inv_code")

	if len(id) >= 6 {
		return id[len(id)-6:]
	}

	return id
}
