package service

import (
	"errors"
	"strconv"

	"poprako-s/internal/domain/model"
	"poprako-s/internal/domain/repo"
)

// PageService 定义页面领域相关的业务能力
type PageService interface {
	// NewCreation 根据业务参数创建 PageCreation 领域模型（ID 由 service 内部生成）
	// 仅章节的监修或图源可以创建页面
	NewCreation(currUserID, chapterID string, index int, ossKey, creatorID string) (*model.PageCreation, error)
	// GenOSSKey 根据页面序号生成 OSS Key
	GenOSSKey(index int) string
}

type pageServiceImpl struct {
	assignmentRepo repo.AssignmentRepo
}

// NewPageService 返回 PageService 的默认实现
func NewPageService(
	assignmentRepo repo.AssignmentRepo,
) PageService {
	return &pageServiceImpl{
		assignmentRepo: assignmentRepo,
	}
}

// NewCreation 构造一个带有 service 生成 ID 的 PageCreation
// 仅章节的监修或图源可以创建页面
func (s *pageServiceImpl) NewCreation(
	currUserID, chapterID string,
	index int,
	ossKey, creatorID string,
) (*model.PageCreation, error) {
	currAssignment, err := s.assignmentRepo.Get(model.AssignmentQueryOpt{
		ChapterID: &chapterID,
		UserID:    &currUserID,
	})
	if err != nil || !currAssignment.HasAnyRole(model.RoleReviewer, model.RoleRawProvider) {
		return nil, errors.New("仅章节监修或图源可以创建页面")
	}

	return &model.PageCreation{
		ID:        GenID("page"),
		ChapterID: chapterID,
		Index:     index,
		OSSKey:    ossKey,
		CreatorID: creatorID,
	}, nil
}

// GenOSSKey 以固定前缀拼接页面序号作为 OSS Key
func (s *pageServiceImpl) GenOSSKey(index int) string {
	return "page_" + strconv.Itoa(index)
}
