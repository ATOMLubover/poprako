package svc

import (
	"errors"

	"poprako-s/internal/domain/model/aggr"
	repo_iface "poprako-s/internal/domain/repo"
	svc_res "poprako-s/internal/domain/svc/res"
	"poprako-s/pkg/util"

	"go.uber.org/zap"
)

// `CommentSvc` provides stateless helpers for team comment aggregate.
type CommentSvc struct{}

// `NewCommentSvc` returns one ready-to-use `CommentSvc`.
func NewCommentSvc() CommentSvc {
	return CommentSvc{}
}

// `CanListComment` validates whether current user can list team comments.
func (CommentSvc) CanListComment(currUid string, teamId string, memberRepo repo_iface.MemberRepo, clsf repo_iface.ErrClsf) svc_res.SvcRes {
	ok, err := memberRepo.ExistByUserTeamId(currUid, teamId)
	if err != nil {
		zap.L().Error(
			"[CommentSvc.CanListComment] failed to verify member",
			zap.String("currUid", currUid),
			zap.String("teamId", teamId),
			zap.Error(err),
		)

		return classifyRepoErr(err, clsf, svc_res.Forbidden, "仅团队成员可查看留言列表", "权限校验失败", "权限校验服务暂不可用", "权限校验失败")
	}

	if !ok {
		return svc_res.Reject(svc_res.Forbidden, "仅团队成员可查看留言列表")
	}

	return svc_res.Accept()
}

// `CanCreateComment` validates whether current user can create team comments.
func (CommentSvc) CanCreateComment(currUid string, teamId string, memberRepo repo_iface.MemberRepo, clsf repo_iface.ErrClsf) svc_res.SvcRes {
	ok, err := memberRepo.ExistByUserTeamId(currUid, teamId)
	if err != nil {
		zap.L().Error(
			"[CommentSvc.CanCreateComment] failed to verify member",
			zap.String("currUid", currUid),
			zap.String("teamId", teamId),
			zap.Error(err),
		)

		return classifyRepoErr(err, clsf, svc_res.Forbidden, "仅团队成员可发布留言", "权限校验失败", "权限校验服务暂不可用", "权限校验失败")
	}

	if !ok {
		return svc_res.Reject(svc_res.Forbidden, "仅团队成员可发布留言")
	}

	return svc_res.Accept()
}

// `NewCommentCre` builds one team comment create payload.
func (CommentSvc) NewCommentCre(currUid string, teamId string, content string) (*aggr.CommentCre, error) {
	if content == "" {
		return nil, errors.New("content 不能为空")
	}

	return &aggr.CommentCre{
		Id:      util.GenId("comment"),
		TeamId:  teamId,
		UserId:  currUid,
		Content: content,
	}, nil
}
