package service

import "poprako-s/internal/domain/model"

// ComicService 定义漫画领域相关的业务能力
type ComicService interface {
	// NewCreation 根据业务参数创建 ComicCreation 领域模型（ID 由 service 内部生成）
	// index 通常由调用方在事务中通过 Count 获取
	NewCreation(worksetID string, index int, title, author, description, creatorID string) *model.ComicCreation
}

type comicServiceImpl struct{}

// NewComicService 返回 ComicService 的默认实现
func NewComicService() ComicService {
	return &comicServiceImpl{}
}

// NewCreation 构造一个带有 service 生成 ID 的 ComicCreation
func (s *comicServiceImpl) NewCreation(
	worksetID string,
	index int,
	title, author, description, creatorID string,
) *model.ComicCreation {
	return &model.ComicCreation{
		ID:          GenID("comic"),
		WorksetID:   worksetID,
		Index:       index,
		Title:       title,
		Author:      author,
		Description: description,
		CreatorID:   creatorID,
	}
}
