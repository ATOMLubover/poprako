package service

import (
	"errors"
	"strings"

	"poprako-s/internal/domain/model"
	"poprako-s/internal/domain/repo"
)

// ComicService 定义漫画领域相关的业务能力
type ComicService interface {
	// NewCreation 根据业务参数创建 ComicCreation 领域模型（ID 由 service 内部生成）
	// index 通常由调用方在事务中通过 Count 获取
	// 仅汉化组管理员可以创建漫画
	NewCreation(
		mr repo.MemberRepo,
		wr repo.WorksetRepo,
		currUserID string,
		worksetID string,
		index int,
		title string,
		author string,
		description string,
		creatorID string,
	) (*model.ComicCreation, error)

	// GenCoverOSSKey 根据漫画 ID 生成封面的 OSS Key
	GenCoverOSSKey(
		comicID string,
	) string
}

// comicServiceImpl 是 ComicService 的具体实现 无内禀状态
type comicServiceImpl struct{}

// NewComicService 返回 ComicService 的默认实现
func NewComicService() ComicService {
	// 返回无状态实现
	return &comicServiceImpl{}
}

// NewCreation 构造一个带有 service 生成 ID 的 ComicCreation
// 仅汉化组管理员可以创建漫画，通过 workset 解析所属汉化组
func (s *comicServiceImpl) NewCreation(
	mr repo.MemberRepo,
	wr repo.WorksetRepo,
	currUserID string,
	worksetID string,
	index int,
	title string,
	author string,
	description string,
	creatorID string,
) (*model.ComicCreation, error) {
	// 查询作品集 解析所属汉化组
	workset, err := wr.GetByID(worksetID)
	if err != nil {
		// 返回查询失败错误
		return nil, errors.New("无法获取漫画所属作品集 创建失败")
	}

	// 查询当前用户在汉化组中的成员记录 用于鉴权
	member, err := mr.Get(model.MemberQueryOpt{
		UserID: &currUserID,
		TeamID: &workset.TeamID,
	})
	if err != nil || !member.HasAnyRole(model.RoleAdmin) {
		// 返回权限错误
		return nil, errors.New("仅汉化组管理员可以创建漫画")
	}

	// 返回创建载荷
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

// GenCoverOSSKey 以固定前缀拼接漫画 ID 作为封面的 OSS Key
func (s *comicServiceImpl) GenCoverOSSKey(
	comicID string,
) string {
	// 返回拼接结果
	return strings.Join([]string{"comic-cover", comicID}, "_")
}
