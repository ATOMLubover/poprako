package model

import "time"

type AssignmentInfo struct {
	ID string

	ChapterID string
	UserID    string

	AssignedRawProviderAt *time.Time
	AssignedTranslatorAt  *time.Time
	AssignedProofreaderAt *time.Time
	AssignedTypesetterAt  *time.Time
	AssignedRedrawerAt    *time.Time
	AssignedReviewerAt    *time.Time
	AssignedPublisherAt   *time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewAssignmentInfo(
	id string,
	chapterID string,
	userID string,
	assignedRawProviderAt *time.Time,
	assignedTranslatorAt *time.Time,
	assignedProofreaderAt *time.Time,
	assignedTypesetterAt *time.Time,
	assignedRedrawerAt *time.Time,
	assignedReviewerAt *time.Time,
	assignedPublisherAt *time.Time,
	createdAt time.Time,
	updatedAt time.Time,
) AssignmentInfo {
	return AssignmentInfo{
		ID:                    id,
		ChapterID:             chapterID,
		UserID:                userID,
		AssignedRawProviderAt: assignedRawProviderAt,
		AssignedTranslatorAt:  assignedTranslatorAt,
		AssignedProofreaderAt: assignedProofreaderAt,
		AssignedTypesetterAt:  assignedTypesetterAt,
		AssignedRedrawerAt:    assignedRedrawerAt,
		AssignedReviewerAt:    assignedReviewerAt,
		AssignedPublisherAt:   assignedPublisherAt,
		CreatedAt:             createdAt,
		UpdatedAt:             updatedAt,
	}
}

func (a *AssignmentInfo) HasAnyRole(roles ...RoleFlag) bool {
	for _, role := range roles {
		switch role {
		case RoleRawProvider:
			if a.AssignedRawProviderAt != nil {
				return true
			}
		case RoleTranslator:
			if a.AssignedTranslatorAt != nil {
				return true
			}
		case RoleProofreader:
			if a.AssignedProofreaderAt != nil {
				return true
			}
		case RoleTypesetter:
			if a.AssignedTypesetterAt != nil {
				return true
			}
		case RoleReviewer:
			if a.AssignedReviewerAt != nil {
				return true
			}
		case RolePublisher:
			if a.AssignedPublisherAt != nil {
				return true
			}
		}
	}

	return false
}

func (a *AssignmentInfo) RoleMask() RoleMask {
	mask := RoleMask(0)
	if a.AssignedRawProviderAt != nil {
		mask |= RoleMask(RoleRawProvider)
	}
	if a.AssignedTranslatorAt != nil {
		mask |= RoleMask(RoleTranslator)
	}
	if a.AssignedProofreaderAt != nil {
		mask |= RoleMask(RoleProofreader)
	}
	if a.AssignedTypesetterAt != nil {
		mask |= RoleMask(RoleTypesetter)
	}
	if a.AssignedReviewerAt != nil {
		mask |= RoleMask(RoleReviewer)
	}
	if a.AssignedPublisherAt != nil {
		mask |= RoleMask(RolePublisher)
	}
	return mask
}

// AssignmentCreation 用于创建分配记录，UUID 由 repository 层生成。
type AssignmentCreation struct {
	ChapterID string
	UserID    string

	AssignedRawProviderAt *time.Time
	AssignedTranslatorAt  *time.Time
	AssignedProofreaderAt *time.Time
	AssignedTypesetterAt  *time.Time
	AssignedRedrawerAt    *time.Time
	AssignedReviewerAt    *time.Time
	AssignedPublisherAt   *time.Time
}

func NewAssignmentCreation(chapterID, userID string, roles RoleMask) AssignmentCreation {
	now := time.Now()

	toAssign := func(role RoleFlag) *time.Time {
		if roles&RoleMask(role) == 0 {
			return nil
		}
		t := now
		return &t
	}

	return AssignmentCreation{
		ChapterID:             chapterID,
		UserID:                userID,
		AssignedRawProviderAt: toAssign(RoleRawProvider),
		AssignedTranslatorAt:  toAssign(RoleTranslator),
		AssignedProofreaderAt: toAssign(RoleProofreader),
		AssignedTypesetterAt:  toAssign(RoleTypesetter),
		AssignedRedrawerAt:    nil,
		AssignedReviewerAt:    toAssign(RoleReviewer),
		AssignedPublisherAt:   toAssign(RolePublisher),
	}
}

// AssignmentUpdate 用于 PUT 语义的全量角色替换，保留已有角色的时间戳。
type AssignmentUpdate struct {
	ID string

	AssignedRawProviderAt *time.Time
	AssignedTranslatorAt  *time.Time
	AssignedProofreaderAt *time.Time
	AssignedTypesetterAt  *time.Time
	AssignedRedrawerAt    *time.Time
	AssignedReviewerAt    *time.Time
	AssignedPublisherAt   *time.Time
}

func NewAssignmentUpdate(id string, current AssignmentInfo, targetRoles RoleMask) AssignmentUpdate {
	now := time.Now()

	resolveAt := func(currentAt *time.Time, role RoleFlag) *time.Time {
		if targetRoles&RoleMask(role) == 0 {
			return nil
		}
		if currentAt != nil {
			t := *currentAt
			return &t
		}
		t := now
		return &t
	}

	return AssignmentUpdate{
		ID:                    id,
		AssignedRawProviderAt: resolveAt(current.AssignedRawProviderAt, RoleRawProvider),
		AssignedTranslatorAt:  resolveAt(current.AssignedTranslatorAt, RoleTranslator),
		AssignedProofreaderAt: resolveAt(current.AssignedProofreaderAt, RoleProofreader),
		AssignedTypesetterAt:  resolveAt(current.AssignedTypesetterAt, RoleTypesetter),
		AssignedRedrawerAt:    current.AssignedRedrawerAt,
		AssignedReviewerAt:    resolveAt(current.AssignedReviewerAt, RoleReviewer),
		AssignedPublisherAt:   resolveAt(current.AssignedPublisherAt, RolePublisher),
	}
}

// AssignmentWithUserInfo 用于列出某章节的所有分配（含用户信息）。
type AssignmentWithUserInfo struct {
	ID        string
	ChapterID string
	User      UserInfo

	AssignedRawProviderAt *time.Time
	AssignedTranslatorAt  *time.Time
	AssignedProofreaderAt *time.Time
	AssignedTypesetterAt  *time.Time
	AssignedRedrawerAt    *time.Time
	AssignedReviewerAt    *time.Time
	AssignedPublisherAt   *time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewAssignmentWithUserInfo(
	id string,
	chapterID string,
	user UserInfo,
	assignedRawProviderAt *time.Time,
	assignedTranslatorAt *time.Time,
	assignedProofreaderAt *time.Time,
	assignedTypesetterAt *time.Time,
	assignedRedrawerAt *time.Time,
	assignedReviewerAt *time.Time,
	assignedPublisherAt *time.Time,
	createdAt time.Time,
	updatedAt time.Time,
) AssignmentWithUserInfo {
	return AssignmentWithUserInfo{
		ID:                    id,
		ChapterID:             chapterID,
		User:                  user,
		AssignedRawProviderAt: assignedRawProviderAt,
		AssignedTranslatorAt:  assignedTranslatorAt,
		AssignedProofreaderAt: assignedProofreaderAt,
		AssignedTypesetterAt:  assignedTypesetterAt,
		AssignedRedrawerAt:    assignedRedrawerAt,
		AssignedReviewerAt:    assignedReviewerAt,
		AssignedPublisherAt:   assignedPublisherAt,
		CreatedAt:             createdAt,
		UpdatedAt:             updatedAt,
	}
}

func (a *AssignmentWithUserInfo) RoleMask() RoleMask {
	mask := RoleMask(0)
	if a.AssignedRawProviderAt != nil {
		mask |= RoleMask(RoleRawProvider)
	}
	if a.AssignedTranslatorAt != nil {
		mask |= RoleMask(RoleTranslator)
	}
	if a.AssignedProofreaderAt != nil {
		mask |= RoleMask(RoleProofreader)
	}
	if a.AssignedTypesetterAt != nil {
		mask |= RoleMask(RoleTypesetter)
	}
	if a.AssignedReviewerAt != nil {
		mask |= RoleMask(RoleReviewer)
	}
	if a.AssignedPublisherAt != nil {
		mask |= RoleMask(RolePublisher)
	}
	return mask
}

// AssignmentWithChapterInfo 用于列出某用户的所有分配（含章节+漫画信息）。
type AssignmentWithChapterInfo struct {
	ID      string
	Chapter ChapterWithComicInfo
	UserID  string

	AssignedRawProviderAt *time.Time
	AssignedTranslatorAt  *time.Time
	AssignedProofreaderAt *time.Time
	AssignedTypesetterAt  *time.Time
	AssignedRedrawerAt    *time.Time
	AssignedReviewerAt    *time.Time
	AssignedPublisherAt   *time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewAssignmentWithChapterInfo(
	id string,
	chapter ChapterWithComicInfo,
	userID string,
	assignedRawProviderAt *time.Time,
	assignedTranslatorAt *time.Time,
	assignedProofreaderAt *time.Time,
	assignedTypesetterAt *time.Time,
	assignedRedrawerAt *time.Time,
	assignedReviewerAt *time.Time,
	assignedPublisherAt *time.Time,
	createdAt time.Time,
	updatedAt time.Time,
) AssignmentWithChapterInfo {
	return AssignmentWithChapterInfo{
		ID:                    id,
		Chapter:               chapter,
		UserID:                userID,
		AssignedRawProviderAt: assignedRawProviderAt,
		AssignedTranslatorAt:  assignedTranslatorAt,
		AssignedProofreaderAt: assignedProofreaderAt,
		AssignedTypesetterAt:  assignedTypesetterAt,
		AssignedRedrawerAt:    assignedRedrawerAt,
		AssignedReviewerAt:    assignedReviewerAt,
		AssignedPublisherAt:   assignedPublisherAt,
		CreatedAt:             createdAt,
		UpdatedAt:             updatedAt,
	}
}

func (a *AssignmentWithChapterInfo) RoleMask() RoleMask {
	mask := RoleMask(0)
	if a.AssignedRawProviderAt != nil {
		mask |= RoleMask(RoleRawProvider)
	}
	if a.AssignedTranslatorAt != nil {
		mask |= RoleMask(RoleTranslator)
	}
	if a.AssignedProofreaderAt != nil {
		mask |= RoleMask(RoleProofreader)
	}
	if a.AssignedTypesetterAt != nil {
		mask |= RoleMask(RoleTypesetter)
	}
	if a.AssignedReviewerAt != nil {
		mask |= RoleMask(RoleReviewer)
	}
	if a.AssignedPublisherAt != nil {
		mask |= RoleMask(RolePublisher)
	}
	return mask
}
