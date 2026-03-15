package value

import (
	"errors"
	"unicode/utf8"
)

type ComicInfo struct {
	ID string `json:"id"`

	WorksetID string       `json:"workset_id"`
	Workset   *WorksetInfo `json:"workset,omitempty"`

	Index       int    `json:"index"`
	Title       string `json:"title"`
	Author      string `json:"author"`
	Description string `json:"description"`

	CoverURL string `json:"cover_url"`

	ChapterCount int       `json:"chapter_count"`
	CreatorID    string    `json:"creator_id"`
	Creator      *UserInfo `json:"creator,omitempty"`

	LastActiveAt int64 `json:"last_active_at"`

	CreatedAt int64 `json:"created_at"`
	UpdatedAt int64 `json:"updated_at"`
}

// ListComicArgs 查询指定工作集下的漫画列表参数。
type ListComicArgs struct {
	WorksetID string   `url:"workset_id"`
	Includes  []string `url:"includes[]"`
	PaginationParams
}

func (args *ListComicArgs) Validate() error {
	if args == nil {
		return errors.New("参数不能为空")
	}

	if args.WorksetID == "" {
		return errors.New("工作集 ID 不能为空")
	}

	if err := args.PaginationParams.Validate(); err != nil {
		return errors.New("分页参数无效: " + err.Error())
	}

	return nil
}

type CreateComicArgs struct {
	WorksetID   string `json:"workset_id"`
	Title       string `json:"title"`
	Author      string `json:"author"`
	Description string `json:"description"`
}

func (cca *CreateComicArgs) Validate() error {
	if cca == nil {
		return errors.New("参数不能为空")
	}

	if cca.WorksetID == "" {
		return errors.New("工作集 ID 不能为空")
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
