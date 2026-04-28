package app_impl

import (
	"context"

	app_iface "poprako-s/internal/app"
	"poprako-s/internal/app/res"
	app_util "poprako-s/internal/app/util"
	"poprako-s/internal/app/val"
	repo_iface "poprako-s/internal/domain/repo"
	repo_infra "poprako-s/internal/infra/repo"

	"go.uber.org/zap"
)

// `userStatsAppImpl` is the default implementation of `UserStatsApp`.
type userStatsAppImpl struct {
	userStatsRepo repo_iface.UserStatsRepo
}

// `NewUserStatsApp` creates one `UserStatsApp` implementation.
func NewUserStatsApp(userStatsRepo repo_iface.UserStatsRepo) app_iface.UserStatsApp {
	if userStatsRepo == nil {
		zap.L().Panic("[NewUserStatsApp] nil dependency", zap.Bool("userStatsRepo", userStatsRepo == nil))
	}

	return &userStatsAppImpl{userStatsRepo: userStatsRepo}
}

// `GetStats` returns current user stats.
func (a *userStatsAppImpl) GetStats(cx context.Context, currUid string) res.AppRes[val.UserStats] {
	lgr := app_util.TakeLgr(cx)

	if currUid == "" {
		return res.Reject[val.UserStats](res.BadRequest, "curr_uid 不能为空")
	}

	stats, err := a.userStatsRepo.EnsureGet(currUid)
	if err != nil {
		if repo_infra.IsTimeout(err) {
			return res.Reject[val.UserStats](res.Timeout, "获取用户统计超时")
		}

		if repo_infra.IsUnavailable(err) {
			return res.Reject[val.UserStats](res.Unavailable, "统计服务暂不可用")
		}

		lgr.Error("[userStatsAppImpl.GetStats] failed to get user stats", zap.Error(err))
		return res.Reject[val.UserStats](res.ServerError, "获取用户统计失败")
	}

	return res.Accept(&val.UserStats{
		UserId:                stats.UserId,
		TotalAssignmentCnt:    stats.TotalAssignmentCnt,
		ActiveAssignmentCnt:   stats.ActiveAssignmentCnt,
		FinishedAssignmentCnt: stats.FinishedAssignmentCnt,
	})
}
