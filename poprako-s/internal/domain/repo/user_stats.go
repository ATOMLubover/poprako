package repo

import "poprako-s/internal/domain/model"

// UserStatsRepo 是用户统计信息仓库的接口，主要负责获取和更新用户任务的统计信息
type UserStatsRepo interface {
	// 获取或创建用户统计信息（注意，可能总是需要一个 UPDATE ON CONFLICT 的操作）
	GetOrCreate(userID string) (*model.UserStats, error)
	// 更新用户统计信息
	Patch(stats *model.UserStats) error
}
