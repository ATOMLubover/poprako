package repo

import (
	"context"

	"poprako-s/internal/domain/model"
)

// ChapterInvitationRepo 是章节邀请仓库接口
type ChapterInvitationRepo interface {
	// List 根据筛选条件返回章节邀请列表
	List(opt model.ChapterInvitationQueryOpt) ([]model.ChapterInvitationInfo, error)
	// Create 持久化一个新的章节邀请
	Create(c *model.ChapterInvitationCreation) (*model.ChapterInvitationInfo, error)
	// Invalidate 使章节邀请失效
	Invalidate(id string) error

	// FromTxnCx 从上下文中获取事务并返回带事务的仓库实例
	FromTxnCx(cx context.Context) (ChapterInvitationRepo, error)
}
