package svc

import (
	"time"

	"poprako-s/internal/domain/model/aggr"
	"poprako-s/internal/domain/model/enum"
	"poprako-s/pkg/util"
)

// `AssignmentSvc` provides stateless domain services for assignment.
type AssignmentSvc struct{}

// `NewAssignmentSvc` returns a ready-to-use `AssignmentSvc`.
func NewAssignmentSvc() AssignmentSvc {
	return AssignmentSvc{}
}

// `NewAssignmentCre` builds one assignment create payload from role mask.
func (AssignmentSvc) NewAssignmentCre(chapterId string, userId string, roleMask aggr.RoleMask) *aggr.AssignmentCre {
	now := time.Now()
	roles := mkTimedRolesFromMask(roleMask, now)

	return &aggr.AssignmentCre{
		Id:         util.GenId("assignment"),
		ChapterId:  chapterId,
		UserId:     userId,
		TimedRoles: roles,
	}
}

// `NewAssignmentPut` builds put payload that preserves existing role timestamps.
func (AssignmentSvc) NewAssignmentPut(curr *aggr.Assignment, roleMask aggr.RoleMask) *aggr.AssignmentPut {
	now := time.Now()
	roles := mkTimedRolesFromCurr(curr, roleMask, now)

	return &aggr.AssignmentPut{
		Id:         curr.Id,
		TimedRoles: roles,
	}
}

func mkTimedRolesFromMask(mask aggr.RoleMask, now time.Time) aggr.TimedRoles {
	roles := aggr.TimedRoles{}

	if mask.HasAnyRole(enum.RoleRawProvider) {
		t := now
		roles.AssignedRawProviderAt = &t
	}
	if mask.HasAnyRole(enum.RoleTranslator) {
		t := now
		roles.AssignedTranslatorAt = &t
	}
	if mask.HasAnyRole(enum.RoleProofreader) {
		t := now
		roles.AssignedProofreaderAt = &t
	}
	if mask.HasAnyRole(enum.RoleTypesetter) {
		t := now
		roles.AssignedTypesetterAt = &t
	}
	if mask.HasAnyRole(enum.RoleRedrawer) {
		t := now
		roles.AssignedRedrawerAt = &t
	}
	if mask.HasAnyRole(enum.RoleReviewer) {
		t := now
		roles.AssignedReviewerAt = &t
	}
	if mask.HasAnyRole(enum.RolePublisher) {
		t := now
		roles.AssignedPublisherAt = &t
	}

	return roles
}

func mkTimedRolesFromCurr(curr *aggr.Assignment, mask aggr.RoleMask, now time.Time) aggr.TimedRoles {
	keepOrNew := func(has bool, prev *time.Time) *time.Time {
		if !has {
			return nil
		}

		if prev != nil {
			t := *prev
			return &t
		}

		t := now
		return &t
	}

	return aggr.TimedRoles{
		AssignedRawProviderAt: keepOrNew(mask.HasAnyRole(enum.RoleRawProvider), curr.AssignedRawProviderAt),
		AssignedTranslatorAt:  keepOrNew(mask.HasAnyRole(enum.RoleTranslator), curr.AssignedTranslatorAt),
		AssignedProofreaderAt: keepOrNew(mask.HasAnyRole(enum.RoleProofreader), curr.AssignedProofreaderAt),
		AssignedTypesetterAt:  keepOrNew(mask.HasAnyRole(enum.RoleTypesetter), curr.AssignedTypesetterAt),
		AssignedRedrawerAt:    keepOrNew(mask.HasAnyRole(enum.RoleRedrawer), curr.AssignedRedrawerAt),
		AssignedReviewerAt:    keepOrNew(mask.HasAnyRole(enum.RoleReviewer), curr.AssignedReviewerAt),
		AssignedPublisherAt:   keepOrNew(mask.HasAnyRole(enum.RolePublisher), curr.AssignedPublisherAt),
	}
}
