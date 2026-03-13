package value

import (
	"errors"
	"unicode/utf8"

	"labelplus-next-web-be/internal/domain/model"
	"labelplus-next-web-be/internal/util"
)

type ChapterInfo struct {
	ID string `json:"id"`

	ComicID string     `json:"comic_id"`
	Comic   *ComicInfo `json:"comic,omitempty"`

	Index     int    `json:"index"`
	ChapterNo string `json:"chapter_no"`

	PageCount           int `json:"page_count"`
	TotalUnitCount      int `json:"total_unit_count"`
	TranslatedUnitCount int `json:"translated_unit_count"`
	ProofreadUnitCount  int `json:"proofread_unit_count"`

	UploadedAt     *int64 `json:"uploaded_at,omitempty"`
	TransalatingAt *int64 `json:"transalating_at,omitempty"`
	TranslatedAt   *int64 `json:"translated_at,omitempty"`
	ProofreadingAt *int64 `json:"proofreading_at,omitempty"`
	ProofreadAt    *int64 `json:"proofread_at,omitempty"`
	TypesettingAt  *int64 `json:"typesetting_at,omitempty"`
	TypesetAt      *int64 `json:"typeset_at,omitempty"`
	ReviewedAt     *int64 `json:"reviewed_at,omitempty"`
	PublishedAt    *int64 `json:"published_at,omitempty"`

	CreatorID   string    `json:"creator_id"`
	CreatorInfo *UserInfo `json:"creator_info,omitempty"`

	CreatedAt int64 `json:"created_at"`
	UpdatedAt int64 `json:"updated_at"`
}

func NewChapterInfoFromModel(chapterDetail model.ChapterDetail) ChapterInfo {
	return ChapterInfo{
		ID:                  chapterDetail.ID,
		ComicID:             chapterDetail.ComicID,
		Index:               chapterDetail.Index,
		ChapterNo:           chapterDetail.ChapterNo,
		PageCount:           chapterDetail.PageCount,
		TotalUnitCount:      chapterDetail.TotalUnitCount,
		TranslatedUnitCount: chapterDetail.TranslatedUnitCount,
		ProofreadUnitCount:  chapterDetail.ProofreadUnitCount,
		UploadedAt:          util.ToUnixPtr(chapterDetail.UploadedAt),
		TransalatingAt:      util.ToUnixPtr(chapterDetail.TransalatingAt),
		TranslatedAt:        util.ToUnixPtr(chapterDetail.TranslatedAt),
		ProofreadingAt:      util.ToUnixPtr(chapterDetail.ProofreadingAt),
		ProofreadAt:         util.ToUnixPtr(chapterDetail.ProofreadAt),
		TypesettingAt:       util.ToUnixPtr(chapterDetail.TypesettingAt),
		TypesetAt:           util.ToUnixPtr(chapterDetail.TypesetAt),
		ReviewedAt:          util.ToUnixPtr(chapterDetail.ReviewedAt),
		PublishedAt:         util.ToUnixPtr(chapterDetail.PublishedAt),
		CreatorID:           chapterDetail.CreatorID,
		CreatedAt:           chapterDetail.CreatedAt.Unix(),
		UpdatedAt:           chapterDetail.UpdatedAt.Unix(),
	}
}

type ListChapterArgs struct {
	ComicID string `url:"comic_id"`
	PaginationParams
}

func (args *ListChapterArgs) Validate() error {
	if args == nil {
		return errors.New("参数不能为空")
	}

	if args.ComicID == "" {
		return errors.New("漫画 ID 不能为空")
	}

	if err := args.PaginationParams.Validate(); err != nil {
		return errors.New("分页参数无效: " + err.Error())
	}

	return nil
}

type CreateChapterArgs struct {
	ComicID   string `json:"comic_id"`
	ChapterNo string `json:"chapter_no"`
}

func (args *CreateChapterArgs) Validate() error {
	if args == nil {
		return errors.New("参数不能为空")
	}

	if args.ComicID == "" {
		return errors.New("漫画 ID 不能为空")
	}

	if args.ChapterNo == "" {
		return errors.New("章节编号不能为空")
	}

	return nil
}

type CreateChapterResult struct {
	ID string `json:"id"`
}

func NewCreateChapterResultFromModel(chapterID string) CreateChapterResult {
	return CreateChapterResult{
		ID: chapterID,
	}
}

type UpdateChapterArgs struct {
	ChapterID string `json:"chapter_id"`

	// ChapterNo 章节编号，最多 10 字符；不传则不更新
	ChapterNo *string `json:"chapter_no,omitempty"`

	// UploadStatus 上传状态，可取值：pending（待上传）、completed（已上传）；不传则不更新
	UploadStatus *model.WorkflowStatus `json:"upload_status,omitempty"`
	// TranslateStatus 翻译状态，可取值：pending（待翻译）、in_progress（翻译中）、completed（已翻译）；不传则不更新
	TranslateStatus *model.WorkflowStatus `json:"translate_status,omitempty"`
	// ProofreadStatus 校对状态，可取值：pending（待校对）、in_progress（校对中）、completed（已校对）；不传则不更新
	ProofreadStatus *model.WorkflowStatus `json:"proofread_status,omitempty"`
	// TypesetStatus 排版状态，可取值：pending（待排版）、in_progress（排版中）、completed（已排版）；不传则不更新
	TypesetStatus *model.WorkflowStatus `json:"typeset_status,omitempty"`
	// ReviewStatus 审阅状态，可取值：pending（待审阅）、completed（已审阅）；不传则不更新
	ReviewStatus *model.WorkflowStatus `json:"review_status,omitempty"`
	// PublishStatus 发布状态，可取值：pending（待发布）、completed（已发布）；不传则不更新
	PublishStatus *model.WorkflowStatus `json:"publish_status,omitempty"`
}

func (args *UpdateChapterArgs) Validate() error {
	if args == nil {
		return errors.New("参数不能为空")
	}

	if args.ChapterID == "" {
		return errors.New("章节 ID 不能为空")
	}

	if args.ChapterNo != nil {
		chapterNoLen := utf8.RuneCountInString(*args.ChapterNo)
		if chapterNoLen == 0 {
			return errors.New("章节编号不能为空字符串")
		}
		if chapterNoLen > 10 {
			return errors.New("章节编号长度不能超过 10 字符")
		}
	}

	if args.UploadStatus != nil && !model.IsValidWorkflowCombination(model.WorkflowUploading, *args.UploadStatus) {
		return errors.New("upload_status 取值无效")
	}
	if args.TranslateStatus != nil && !model.IsValidWorkflowCombination(model.WorkflowTranslating, *args.TranslateStatus) {
		return errors.New("translate_status 取值无效")
	}
	if args.ProofreadStatus != nil && !model.IsValidWorkflowCombination(model.WorkflowProofreading, *args.ProofreadStatus) {
		return errors.New("proofread_status 取值无效")
	}
	if args.TypesetStatus != nil && !model.IsValidWorkflowCombination(model.WorkflowTypesetting, *args.TypesetStatus) {
		return errors.New("typeset_status 取值无效")
	}
	if args.ReviewStatus != nil && !model.IsValidWorkflowCombination(model.WorkflowReviewing, *args.ReviewStatus) {
		return errors.New("review_status 取值无效")
	}
	if args.PublishStatus != nil && !model.IsValidWorkflowCombination(model.WorkflowPublishing, *args.PublishStatus) {
		return errors.New("publish_status 取值无效")
	}

	return nil
}
