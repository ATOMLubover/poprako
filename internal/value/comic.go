package value

import (
	"errors"
	"unicode/utf8"

	"labelplus-next-web-be/internal/domain/model"
)

type ComicInfo struct {
	ID string `json:"id"`

	Index       int    `json:"index"`
	Title       string `json:"title"`
	Author      string `json:"author"`
	Description string `json:"description"`

	CoverURL string `json:"coverUrl"`

	ChapterCount int    `json:"chapter_count"`
	CreatorID    string `json:"creator_id"`

	CreatedAt int64 `json:"created_at"`
	UpdatedAt int64 `json:"updated_at"`
}

func NewComicInfoFromModel(comicInfo *model.ComicInfo) *ComicInfo {
	return &ComicInfo{
		ID:           comicInfo.ID,
		Index:        comicInfo.Index,
		Title:        comicInfo.Title,
		Author:       comicInfo.Author,
		Description:  comicInfo.Description,
		CoverURL:     comicInfo.CoverURL,
		ChapterCount: comicInfo.ChapterCount,
		CreatorID:    comicInfo.CreatorID,
		CreatedAt:    comicInfo.CreatedAt.Unix(),
		UpdatedAt:    comicInfo.UpdatedAt.Unix(),
	}
}

type ListTeamComicArgs struct {
	TeamID string `url:"team_id"`
	PaginationParams
}

func (ltca *ListTeamComicArgs) Validate() error {
	if ltca == nil {
		return errors.New("参数不能为空")
	}

	if ltca.TeamID == "" {
		return errors.New("汉化组 ID 不能为空")
	}

	if err := ltca.PaginationParams.Validate(); err != nil {
		return errors.New("分页参数无效: " + err.Error())
	}

	return nil
}

type CreateComicArgs struct {
	TeamID      string `json:"team_id"`
	Title       string `json:"title"`
	Author      string `json:"author"`
	Description string `json:"description"`
}

func (cca *CreateComicArgs) Validate() error {
	if cca == nil {
		return errors.New("参数不能为空")
	}

	if cca.TeamID == "" {
		return errors.New("汉化组 ID 不能为空")
	}

	titleLen := utf8.RuneCountInString(cca.Title)
	if titleLen <= 0 || titleLen > 100 {
		return errors.New("漫画标题长度必须在 1~100 字符之间")
	}

	authorLen := utf8.RuneCountInString(cca.Author)
	if authorLen <= 0 || authorLen > 50 {
		return errors.New("作者名长度必须在 1~50 字符之间")
	}

	descLen := utf8.RuneCountInString(cca.Description)
	if descLen > 500 {
		return errors.New("漫画描述长度不能超过 500 字符")
	}

	return nil
}

type CreateComicResult struct {
	ID string `json:"id"`
}

type UpdateComicArgs struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Author      string `json:"author"`
	Description string `json:"description"`
}

func (uca *UpdateComicArgs) Validate() error {
	if uca == nil {
		return errors.New("参数不能为空")
	}

	if uca.ID == "" {
		return errors.New("漫画 ID 不能为空")
	}

	titleLen := utf8.RuneCountInString(uca.Title)
	if titleLen <= 0 || titleLen > 100 {
		return errors.New("漫画标题长度必须在 1~100 字符之间")
	}

	authorLen := utf8.RuneCountInString(uca.Author)
	if authorLen <= 0 || authorLen > 50 {
		return errors.New("作者名长度必须在 1~50 字符之间")
	}

	descLen := utf8.RuneCountInString(uca.Description)
	if descLen > 500 {
		return errors.New("漫画描述长度不能超过 500 字符")
	}

	return nil
}
