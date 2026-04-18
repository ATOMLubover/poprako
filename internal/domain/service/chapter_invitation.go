package service

import (
	"errors"

	"poprako-s/internal/domain/model"
)

// ChapterInvitationService 定义章节邀请领域能力
type ChapterInvitationService interface {
	// NewCreation 根据参数构造章节邀请创建载荷
	NewCreation(
		inviterID string,
		chapterID string,
		inviteeQQ string,
		roles ...model.Role,
	) (*model.ChapterInvitationCreation, error)
	// GenInvCode 生成邀请代码
	GenInvCode() string
}

type chapterInvitationServiceImpl struct{}

// NewChapterInvitationService 返回默认章节邀请服务实现
func NewChapterInvitationService() ChapterInvitationService {
	return &chapterInvitationServiceImpl{}
}

func (s *chapterInvitationServiceImpl) NewCreation(
	inviterID string,
	chapterID string,
	inviteeQQ string,
	roles ...model.Role,
) (*model.ChapterInvitationCreation, error) {
	if len(roles) == 0 {
		return nil, errors.New("至少指定一个角色")
	}

	c := &model.ChapterInvitationCreation{
		ID:             GenID("chapter_invitation"),
		ChapterID:      chapterID,
		InviterID:      inviterID,
		InviteeQQ:      inviteeQQ,
		InvitationCode: s.GenInvCode(),
	}

	for _, role := range roles {
		switch role {
		case model.RoleRawProvider:
			c.ToBeRawProvider = true
		case model.RoleTranslator:
			c.ToBeTranslator = true
		case model.RoleProofreader:
			c.ToBeProofreader = true
		case model.RoleTypesetter:
			c.ToBeTypesetter = true
		case model.RoleReviewer:
			c.ToBeReviewer = true
		case model.RolePublisher:
			c.ToBePublisher = true
		case model.RoleAdmin:
			return nil, errors.New("章节邀请不支持管理员角色")
		}
	}

	if !c.ToBeRawProvider &&
		!c.ToBeTranslator &&
		!c.ToBeProofreader &&
		!c.ToBeTypesetter &&
		!c.ToBeReviewer &&
		!c.ToBePublisher {
		return nil, errors.New("至少指定一个有效角色")
	}

	return c, nil
}

func (s *chapterInvitationServiceImpl) GenInvCode() string {
	id := GenID("chapter_invitation")

	if len(id) >= 6 {
		return id[len(id)-6:]
	}

	return id
}
