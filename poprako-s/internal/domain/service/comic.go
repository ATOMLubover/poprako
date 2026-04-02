package service

import (
	"errors"

	"poprako-s/internal/domain/model"
	"poprako-s/internal/domain/repo"
)

// ComicService 定义漫画领域相关的业务能力
type ComicService interface {
	// NewCreation 根据业务参数创建 ComicCreation 领域模型（ID 由 service 内部生成）
	// index 通常由调用方在事务中通过 Count 获取
	// 仅团队管理员可以创建漫画
	NewCreation(currUserID, worksetID string, index int, title, author, description, creatorID string) (*model.ComicCreation, error)
}

type comicServiceImpl struct {
	memberRepo  repo.MemberRepo
	worksetRepo repo.WorksetRepo
	comicRepo   repo.ComicRepo
}

// NewComicService 返回 ComicService 的默认实现
func NewComicService(memberRepo repo.MemberRepo, worksetRepo repo.WorksetRepo, comicRepo repo.ComicRepo) ComicService {
	return &comicServiceImpl{
		memberRepo:  memberRepo,
		worksetRepo: worksetRepo,
		comicRepo:   comicRepo,
	}
}

// NewCreation 构造一个带有 service 生成 ID 的 ComicCreation
// 仅团队管理员可以创建漫画，通过 workset 解析所属团队
func (s *comicServiceImpl) NewCreation(
	currUserID, worksetID string,
	index int,
	title, author, description, creatorID string,
) (*model.ComicCreation, error) {
	workset, err := s.worksetRepo.GetByID(worksetID)
	if err != nil {
		return nil, errors.New("无法获取漫画所属作品集，创建失败")
	}
	member, err := s.memberRepo.Get(model.MemberQueryOpt{
		UserID: &currUserID,
		TeamID: &workset.TeamID,
	})
	if err != nil || !member.HasAnyRole(model.RoleAdmin) {
		return nil, errors.New("仅团队管理员可以创建漫画")
	}

	return &model.ComicCreation{
		ID:          GenID("comic"),
		WorksetID:   worksetID,
		Index:       index,
		Title:       title,
		Author:      author,
		Description: description,
		CreatorID:   creatorID,
	}, nil
}
