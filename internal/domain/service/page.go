package service

import (
	"errors"

	"poprako-s/internal/domain/model"
	"poprako-s/internal/domain/repo"
)

// PageService 定义页面领域相关的业务能力
type PageService interface {
	// NewCreation 根据业务参数创建 PageCreation 领域模型（ID 由 service 内部生成）
	// 仅章节的监修或图源可以创建页面
	NewCreation(
		ar repo.AssignmentRepo,
		currUserID string,
		chapterID string,
		index int,
		ossKey string,
		creatorID string,
	) (*model.PageCreation, error)
	// GenOSSKey 根据章节 ID 与页面 ID 生成 OSS Key，格式为 chapter_{chapterID}/page_{pageID}
	GenOSSKey(
		chapterID string,
		pageID string,
	) string
}

// pageServiceImpl 是 PageService 的具体实现 无内禀状态
type pageServiceImpl struct{}

// NewPageService 返回 PageService 的默认实现
func NewPageService() PageService {
	// 返回无状态实现
	return &pageServiceImpl{}
}

// NewCreation 构造一个带有 service 生成 ID 的 PageCreation
// 仅章节的监修或图源可以创建页面
func (s *pageServiceImpl) NewCreation(
	ar repo.AssignmentRepo,
	currUserID string,
	chapterID string,
	index int,
	ossKey string,
	creatorID string,
) (*model.PageCreation, error) {
	// 查询当前用户在章节中的分配 用于鉴权
	currAssignment, err := ar.Get(model.AssignmentQueryOpt{
		ChapterID: &chapterID,
		UserID:    &currUserID,
	})
	if err != nil || !currAssignment.HasAnyRole(model.RoleReviewer, model.RoleRawProvider) {
		// 返回权限错误
		return nil, errors.New("仅章节监修或图源可以创建页面")
	}

	// 返回创建载荷
	return &model.PageCreation{
		ID:        GenID("page"),
		ChapterID: chapterID,
		Index:     index,
		OSSKey:    ossKey,
		CreatorID: creatorID,
	}, nil
}

// GenOSSKey 以目录式命名生成页面图片的 OSS Key，格式为 chapter_{chapterID}/page_{pageID}
func (s *pageServiceImpl) GenOSSKey(
	chapterID string,
	pageID string,
) string {
	// 必须包含 pageID 以避免重排/重拍时 key 变化或覆盖
	return "chapter_" + chapterID + "/page_" + pageID
}
