package service

import (
	"errors"
	"time"

	"poprako-s/internal/domain/event"
	"poprako-s/internal/domain/model"
	"poprako-s/internal/domain/repo"
)

// AssignmentService 定义分配记录领域相关的业务能力
type AssignmentService interface {
	// NewCreation 根据章节 ID、用户 ID 和角色掩码创建 AssignmentCreation
	// 被分配的角色会记录当前时间作为分配时间
	// 仅章节的监修可以创建分配
	NewCreation(
		ar repo.AssignmentRepo,
		currUserID string,
		chapterID string,
		userID string,
		roles model.RoleMask,
	) (*model.AssignmentCreation, error)
	// NewInitReviewerCreation 为章节创建者生成初始监修分配载荷
	// 无需鉴权，仅用于章节刚创建时的引导性分配
	NewInitReviewerCreation(
		chapterID string,
		creatorID string,
	) *model.AssignmentCreation
	// NewRemovalEvent 根据被删除分配记录的信息构造 AssignmentRemovedEvent
	NewRemovalEvent(
		assignment *model.AssignmentInfo,
		wasPublished bool,
	) *event.AssignmentRemovedEvent
	// NewUpdate 根据当前分配信息和目标角色掩码生成 AssignmentUpdate
	// 已有角色保留原时间戳，新增角色使用当前时间，移除的角色清空时间戳
	// 仅章节的监修可以更新分配
	NewUpdate(
		ar repo.AssignmentRepo,
		currUserID string,
		id string,
		curr *model.AssignmentInfo,
		targetRoles model.RoleMask,
	) (*model.AssignmentUpdate, error)
}

// assignmentServiceImpl 是 AssignmentService 的具体实现 无内禀状态
type assignmentServiceImpl struct{}

// NewAssignmentService 返回 AssignmentService 的默认实现
func NewAssignmentService() AssignmentService {
	// 返回无状态实现
	return &assignmentServiceImpl{}
}

// NewInitReviewerCreation 为章节创建者生成初始监修分配载荷
// 无需鉴权，仅用于章节刚创建时的引导性分配
func (s *assignmentServiceImpl) NewInitReviewerCreation(
	chapterID string,
	creatorID string,
) *model.AssignmentCreation {
	// 记录当前时间作为监修分配时间
	now := time.Now()

	// 返回仅含监修角色的创建载荷
	return &model.AssignmentCreation{
		ID:                    GenID("assignment"),
		ChapterID:             chapterID,
		UserID:                creatorID,
		AssignedRawProviderAt: nil,
		AssignedTranslatorAt:  nil,
		AssignedProofreaderAt: nil,
		AssignedTypesetterAt:  nil,
		AssignedRedrawerAt:    nil,
		AssignedReviewerAt:    &now,
		AssignedPublisherAt:   nil,
	}
}

// NewRemovalEvent 根据被删除分配记录的信息构造 AssignmentRemovedEvent
func (s *assignmentServiceImpl) NewRemovalEvent(
	assignment *model.AssignmentInfo,
	wasPublished bool,
) *event.AssignmentRemovedEvent {
	// 返回组装好的删除事件
	return &event.AssignmentRemovedEvent{
		UserID:       assignment.UserID,
		WasPublished: wasPublished,
	}
}

// NewCreation 构造一个带有 service 生成 ID 的 AssignmentCreation
// 仅章节的监修可以创建分配
func (s *assignmentServiceImpl) NewCreation(
	ar repo.AssignmentRepo,
	currUserID string,
	chapterID string,
	userID string,
	roles model.RoleMask,
) (*model.AssignmentCreation, error) {
	// 查询当前用户在章节中的分配 用于鉴权
	currAssignment, err := ar.Get(model.AssignmentQueryOpt{
		ChapterID: &chapterID,
		UserID:    &currUserID,
	})
	if err != nil || !currAssignment.HasAnyRole(model.RoleReviewer) {
		// 返回权限错误
		return nil, errors.New("仅章节监修可以创建分配")
	}

	// TODO: 应该检查被分配用户是否有资质担任目标角色

	// 记录当前时间 作为新分配角色时间戳
	now := time.Now()

	// 依据目标角色掩码生成角色时间戳
	toAssign := func(role model.Role) *time.Time {
		if roles&model.RoleMask(role) == 0 {
			return nil
		}

		t := now

		return &t
	}

	// 返回创建载荷
	c := &model.AssignmentCreation{
		ID:                    GenID("assignment"),
		ChapterID:             chapterID,
		UserID:                userID,
		AssignedRawProviderAt: toAssign(model.RoleRawProvider),
		AssignedTranslatorAt:  toAssign(model.RoleTranslator),
		AssignedProofreaderAt: toAssign(model.RoleProofreader),
		AssignedTypesetterAt:  toAssign(model.RoleTypesetter),
		AssignedRedrawerAt:    nil,
		AssignedReviewerAt:    toAssign(model.RoleReviewer),
		AssignedPublisherAt:   toAssign(model.RolePublisher),
	}

	// 分配创建时推送同步统计事件
	c.PushEvent(&event.AssignmentCreatedEvent{
		UserID:    userID,
		ChapterID: chapterID,
	})

	return c, nil
}

// NewUpdate 根据目标角色掩码和当前分配信息构造 AssignmentUpdate
// 仅章节的监修可以更新分配
func (s *assignmentServiceImpl) NewUpdate(
	ar repo.AssignmentRepo,
	currUserID string,
	id string,
	curr *model.AssignmentInfo,
	targetRoles model.RoleMask,
) (*model.AssignmentUpdate, error) {
	// 查询当前用户在章节中的分配 用于鉴权
	currAssignment, err := ar.Get(model.AssignmentQueryOpt{
		ChapterID: &curr.ChapterID,
		UserID:    &currUserID,
	})
	if err != nil || !currAssignment.HasAnyRole(model.RoleReviewer) {
		// 返回权限错误
		return nil, errors.New("仅章节监修可以更新分配")
	}

	// 记录当前时间 用于新增角色
	now := time.Now()

	// 解析目标角色对应时间戳
	resolve := func(currAt *time.Time, role model.Role) *time.Time {
		if targetRoles&model.RoleMask(role) == 0 {
			return nil
		}

		if currAt != nil {
			t := *currAt
			return &t
		}

		t := now

		return &t
	}

	// 返回更新载荷
	return &model.AssignmentUpdate{
		ID:                    id,
		AssignedRawProviderAt: resolve(curr.AssignedRawProviderAt, model.RoleRawProvider),
		AssignedTranslatorAt:  resolve(curr.AssignedTranslatorAt, model.RoleTranslator),
		AssignedProofreaderAt: resolve(curr.AssignedProofreaderAt, model.RoleProofreader),
		AssignedTypesetterAt:  resolve(curr.AssignedTypesetterAt, model.RoleTypesetter),
		AssignedRedrawerAt:    curr.AssignedRedrawerAt,
		AssignedReviewerAt:    resolve(curr.AssignedReviewerAt, model.RoleReviewer),
		AssignedPublisherAt:   resolve(curr.AssignedPublisherAt, model.RolePublisher),
	}, nil
}
