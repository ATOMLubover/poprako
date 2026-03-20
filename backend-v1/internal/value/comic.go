package value

import (
	"errors"
	"unicode/utf8"

	"labelplus-next-web-be/internal/domain/model"
)

type ComicInfo struct {
	ID string `json:"id"`

	WorksetID string       `json:"workset_id"`
	Workset   *WorksetInfo `json:"workset,omitempty"`

	Index       int    `json:"index"`
	Title       string `json:"title"`
	Author      string `json:"author"`
	Description string `json:"description"`

	ChapterCount int       `json:"chapter_count"`
	CreatorID    string    `json:"creator_id"`
	Creator      *UserInfo `json:"creator,omitempty"`

	// FIXME: update?
	LastActiveAt int64 `json:"last_active_at"`

	CreatedAt int64 `json:"created_at"`
	UpdatedAt int64 `json:"updated_at"`
}

// ListComicArgs 查询指定工作集下的漫画列表参数。
type ListComicArgs struct {
	WorksetID string   `url:"workset_id"`
	Includes  []string `url:"includes[]"`

	FuzzyTitle string `url:"fuzzy_title,omitempty"`

	UploadStatus    *model.WorkflowStatus `url:"upload_status,omitempty"`
	TranslateStatus *model.WorkflowStatus `url:"translate_status,omitempty"`
	ProofreadStatus *model.WorkflowStatus `url:"proofread_status,omitempty"`
	TypesetStatus   *model.WorkflowStatus `url:"typeset_status,omitempty"`
	ReviewStatus    *model.WorkflowStatus `url:"review_status,omitempty"`
	PublishStatus   *model.WorkflowStatus `url:"publish_status,omitempty"`

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

	if args.UploadStatus != nil && !model.IsValidWorkflowCombination(model.WorkflowUploading, *args.UploadStatus) {
		return errors.New("上传状态参数无效")
	}

	if args.TranslateStatus != nil && !model.IsValidWorkflowCombination(model.WorkflowTranslating, *args.TranslateStatus) {
		return errors.New("翻译状态参数无效")
	}

	if args.ProofreadStatus != nil && !model.IsValidWorkflowCombination(model.WorkflowProofreading, *args.ProofreadStatus) {
		return errors.New("校对状态参数无效")
	}

	if args.TypesetStatus != nil && !model.IsValidWorkflowCombination(model.WorkflowTypesetting, *args.TypesetStatus) {
		return errors.New("嵌字状态参数无效")
	}

	if args.ReviewStatus != nil && !model.IsValidWorkflowCombination(model.WorkflowReviewing, *args.ReviewStatus) {
		return errors.New("审核状态参数无效")
	}

	if args.PublishStatus != nil && !model.IsValidWorkflowCombination(model.WorkflowPublishing, *args.PublishStatus) {
		return errors.New("发布状态参数无效")
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

type ComicCoverResult struct {
	CoverURL string `json:"cover_url"`
}

func NewComicCoverResult(coverURL string) ComicCoverResult {
	return ComicCoverResult{CoverURL: coverURL}
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
