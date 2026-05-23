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

// `AnnouncementSvc` provides stateless helpers for team announcement aggregate.
type AnnouncementSvc struct{}

// `NewAnnouncementSvc` returns one ready-to-use `AnnouncementSvc`.
func NewAnnouncementSvc() AnnouncementSvc {
	return AnnouncementSvc{}
}

// `CanListAnnouncement` validates whether current user can list team announcements.
func (AnnouncementSvc) CanListAnnouncement(currUid string, teamId string, memberRepo repo_iface.MemberRepo, clsf repo_iface.ErrClsf) svc_res.SvcRes {
	ok, err := memberRepo.ExistByUserTeamId(currUid, teamId)
	if err != nil {
		zap.L().Error(
			"[AnnouncementSvc.CanListAnnouncement] failed to verify member",
			zap.String("currUid", currUid),
			zap.String("teamId", teamId),
			zap.Error(err),
		)

		return classifyRepoErr(err, clsf, svc_res.Forbidden, "仅团队成员可查看公告列表", "权限校验失败", "权限校验服务暂不可用", "权限校验失败")
	}

	if !ok {
		return svc_res.Reject(svc_res.Forbidden, "仅团队成员可查看公告列表")
	}

	return svc_res.Accept()
}

// `CanCreateAnnouncement` validates whether current user can create team announcements.
func (AnnouncementSvc) CanCreateAnnouncement(currUid string, teamId string, memberRepo repo_iface.MemberRepo, clsf repo_iface.ErrClsf) svc_res.SvcRes {
	currMember, err := memberRepo.GetByUserTeamId(currUid, teamId)
	if err != nil {
		zap.L().Error(
			"[AnnouncementSvc.CanCreateAnnouncement] failed to verify admin member",
			zap.String("currUid", currUid),
			zap.String("teamId", teamId),
			zap.Error(err),
		)

		return classifyRepoErr(err, clsf, svc_res.Forbidden, "仅团队管理员可发布公告", "权限校验失败", "权限校验服务暂不可用", "权限校验失败")
	}

	if currMember == nil || !currMember.HasAnyRole(enum.RoleAdmin) {
		return svc_res.Reject(svc_res.Forbidden, "仅团队管理员可发布公告")
	}

	return svc_res.Accept()
}

// `NewAnnouncementCre` builds one announcement create payload.
func (AnnouncementSvc) NewAnnouncementCre(currUid string, teamId string, title string, content string) (*aggr.AnnouncementCre, error) {
	if title == "" {
		return nil, errors.New("title 不能为空")
	}

	if content == "" {
		return nil, errors.New("content 不能为空")
	}

	return &aggr.AnnouncementCre{
		Id:      util.GenId("announcement"),
		TeamId:  teamId,
		UserId:  currUid,
		Title:   title,
		Content: content,
	}, nil
}
