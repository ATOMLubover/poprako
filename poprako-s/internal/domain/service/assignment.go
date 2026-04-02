package service

import (
	"errors"
	"time"

	"poprako-s/internal/domain/model"
	"poprako-s/internal/domain/repo"
)

// AssignmentService 定义分配记录领域相关的业务能力
type AssignmentService interface {
	// NewCreation 根据章节 ID、用户 ID 和角色掩码创建 AssignmentCreation
	// 被分配的角色会记录当前时间作为分配时间
	// 仅章节的监修可以创建分配
	NewCreation(currUserID, chapterID, userID string, roles model.RoleMask) (*model.AssignmentCreation, error)
	// NewUpdate 根据当前分配信息和目标角色掩码生成 AssignmentUpdate
	// 已有角色保留原时间戳，新增角色使用当前时间，移除的角色清空时间戳
	// 仅章节的监修可以更新分配
	NewUpdate(currUserID string, id string, current *model.AssignmentInfo, targetRoles model.RoleMask) (*model.AssignmentUpdate, error)
}

type assignmentServiceImpl struct {
	assignmentRepo repo.AssignmentRepo
}

// NewAssignmentService 返回 AssignmentService 的默认实现
func NewAssignmentService(
	assignmentRepo repo.AssignmentRepo,
) AssignmentService {
	return &assignmentServiceImpl{
		assignmentRepo: assignmentRepo,
	}
}

// NewCreation 构造一个带有 service 生成 ID 的 AssignmentCreation
// 仅章节的监修可以创建分配
func (s *assignmentServiceImpl) NewCreation(
	currUserID, chapterID, userID string,
	roles model.RoleMask,
) (*model.AssignmentCreation, error) {
	currAssignment, err := s.assignmentRepo.Get(model.AssignmentQueryOpt{
		ChapterID: &chapterID,
		UserID:    &currUserID,
	})
	if err != nil || !currAssignment.HasAnyRole(model.RoleReviewer) {
		return nil, errors.New("仅章节监修可以创建分配")
	}
	now := time.Now()

	toAssign := func(role model.Role) *time.Time {
		if roles&model.RoleMask(role) == 0 {
			return nil
		}
		t := now
		return &t
	}

	return &model.AssignmentCreation{
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
	}, nil
}

// NewUpdate 根据目标角色掩码和当前分配信息构造 AssignmentUpdate
// 仅章节的监修可以更新分配
func (s *assignmentServiceImpl) NewUpdate(
	currUserID string,
	id string,
	current *model.AssignmentInfo,
	targetRoles model.RoleMask,
) (*model.AssignmentUpdate, error) {
	currAssignment, err := s.assignmentRepo.Get(model.AssignmentQueryOpt{
		ChapterID: &current.ChapterID,
		UserID:    &currUserID,
	})
	if err != nil || !currAssignment.HasAnyRole(model.RoleReviewer) {
		return nil, errors.New("仅章节监修可以更新分配")
	}
	now := time.Now()

	resolve := func(currentAt *time.Time, role model.Role) *time.Time {
		if targetRoles&model.RoleMask(role) == 0 {
			return nil
		}
		if currentAt != nil {
			t := *currentAt
			return &t
		}
		t := now
		return &t
	}

	return &model.AssignmentUpdate{
		ID:                    id,
		AssignedRawProviderAt: resolve(current.AssignedRawProviderAt, model.RoleRawProvider),
		AssignedTranslatorAt:  resolve(current.AssignedTranslatorAt, model.RoleTranslator),
		AssignedProofreaderAt: resolve(current.AssignedProofreaderAt, model.RoleProofreader),
		AssignedTypesetterAt:  resolve(current.AssignedTypesetterAt, model.RoleTypesetter),
		AssignedRedrawerAt:    current.AssignedRedrawerAt,
		AssignedReviewerAt:    resolve(current.AssignedReviewerAt, model.RoleReviewer),
		AssignedPublisherAt:   resolve(current.AssignedPublisherAt, model.RolePublisher),
	}, nil
}
