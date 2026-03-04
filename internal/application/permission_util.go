package application

import (
	"errors"

	"labelplus-next-web-be/internal/domain/model"
	"labelplus-next-web-be/internal/domain/service"
	"labelplus-next-web-be/internal/repository/query_option"
	"labelplus-next-web-be/internal/util"

	"go.uber.org/zap"
)

func (ia *invitationApplication) checkCurrentUserPermissionInTeam(
	scope util.TraceScope,
	currentUserID string,
	targetTeamID string,
	permission model.Permission,
) error {
	// 鉴权：检查当前用户在指定的汉化组是否有权限执行相关操作
	currentUserMemberships, err := ia.memberRepository.List(
		nil,
		query_option.MemberQuery().FilterByTeamID(targetTeamID),
		query_option.MemberQuery().FilterByUserID(currentUserID),
	)
	if err != nil {
		scope.Logger().Error("检查用户权限失败: 获取当前用户的成员信息失败", zap.Error(err))
		return errors.New("无法获取成员信息")
	}

	if !service.CheckInvitationPermission(
		targetTeamID,
		currentUserMemberships,
		permission,
	) {
		return errors.New("没有权限执行该操作")
	}

	return nil
}
