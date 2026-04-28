package event_infra

import (
	"context"
	"fmt"
	"time"

	"poprako-s/internal/domain/model/aggr"
	"poprako-s/internal/domain/model/event"
	repo_iface "poprako-s/internal/domain/repo"
	event_iface "poprako-s/internal/event"
	"poprako-s/pkg/util"

	"go.uber.org/zap"
)

type UpdateUserActiveHandler struct {
	userRepo repo_iface.UserRepo
}

func NewUpdateUserActiveHandler(userRepo repo_iface.UserRepo) *UpdateUserActiveHandler {
	return &UpdateUserActiveHandler{userRepo: userRepo}
}

func (h *UpdateUserActiveHandler) EvTyp() event_iface.EvTyp {
	return event.EvUserLogin
}

func (h *UpdateUserActiveHandler) Handle(_ context.Context, ev event_iface.Event) {
	payload, ok := ev.Payload().(*event.UserLoginEv)
	if !ok || payload == nil {
		zap.L().Error(
			"[UpdateUserActiveHandler.Handle] invalid event payload for UpdateUserActiveHandler",
			zap.Any("payload", ev.Payload()),
		)

		return
	}

	now := time.Now()

	if err := h.userRepo.Refresh(payload.UserId, now); err != nil {
		zap.L().Error(
			"[UpdateUserActiveHandler.Handle] failed to update user active time",
			zap.String("user_id", payload.UserId),
			zap.Error(err),
		)

		return
	}
}

type NotifyInvitorHandler struct {
	teamRepo    repo_iface.TeamRepo
	sysMailRepo repo_iface.SysMailRepo
}

func NewNotifyInvitorHandler(sysMailRepo repo_iface.SysMailRepo) *NotifyInvitorHandler {
	return &NotifyInvitorHandler{sysMailRepo: sysMailRepo}
}

func (h *NotifyInvitorHandler) EvTyp() event_iface.EvTyp {
	return event.EvUserReg
}

func (h *NotifyInvitorHandler) Handle(_ context.Context, ev event_iface.Event) {
	payload, ok := ev.Payload().(*event.UserRegEv)
	if !ok || payload == nil {
		zap.L().Error(
			"[NotifyInvitorHandler.Handle] invalid event payload for NotifyInvitorHandler",
			zap.Any("payload", ev.Payload()),
		)
		return
	}

	const TITLE = "你的邀请码已被使用"

	team, err := h.teamRepo.GetById(payload.TeamId)
	if err != nil {
		zap.L().Error(
			"[NotifyInvitorHandler.Handle] failed to get team info",
			zap.String("team_id", payload.TeamId),
			zap.Error(err),
		)
		return
	}

	cont := fmt.Sprintf(
		"你的邀请码已被使用，%s 已加入团队 %s",
		payload.InviteeQid,
		team.Name,
	)

	id := util.GenId("sys_mail")

	cre := &aggr.SysMailCre{
		Id:      id,
		RcvId:   payload.InvitorId,
		Title:   TITLE,
		Content: cont,
	}

	if err := h.sysMailRepo.Send(cre); err != nil {
		zap.L().Error(
			"[NotifyInvitorHandler.Handle] failed to create invitee reg mail",
			zap.String("invitor_id", payload.InvitorId),
			zap.String("invitee_qid", payload.InviteeQid),
			zap.Error(err),
		)
		return
	}
}
