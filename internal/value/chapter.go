package value

import (
	"errors"
	"unicode/utf8"

	"labelplus-next-web-be/internal/domain/model"
	"labelplus-next-web-be/internal/util"
)

type ChapterDetail struct {
	ID string `json:"id"`

	ComicID   string `json:"comic_id"`
	Index     int    `json:"index"`
	ChapterNo string `json:"chapter_no"`

	PageCount           int `json:"page_count"`
	TotalUnitCount      int `json:"total_unit_count"`
	TranslatedUnitCount int `json:"translated_unit_count"`
	ProofreadUnitCount  int `json:"proofread_unit_count"`

	CoverURL string `json:"cover_url"`

	UploadedAt     *int64 `json:"uploaded_at,omitempty"`
	TransalatingAt *int64 `json:"transalating_at,omitempty"`
	TranslatedAt   *int64 `json:"translated_at,omitempty"`
	ProofreadingAt *int64 `json:"proofreading_at,omitempty"`
	ProofreadAt    *int64 `json:"proofread_at,omitempty"`
	TypesettingAt  *int64 `json:"typesetting_at,omitempty"`
	TypesetAt      *int64 `json:"typeset_at,omitempty"`
	ReviewedAt     *int64 `json:"reviewed_at,omitempty"`
	PublishedAt    *int64 `json:"published_at,omitempty"`

	CreatorID string `json:"creator_id"`

	CreatedAt int64 `json:"created_at"`
	UpdatedAt int64 `json:"updated_at"`
}

func NewChapterInfoFromModel(chapterDetail model.ChapterInfo) ChapterDetail {
	return ChapterDetail{
		ID:                  chapterDetail.ID,
		ComicID:             chapterDetail.ComicID,
		Index:               chapterDetail.Index,
		ChapterNo:           chapterDetail.ChapterNo,
		PageCount:           chapterDetail.PageCount,
		TotalUnitCount:      chapterDetail.TotalUnitCount,
		TranslatedUnitCount: chapterDetail.TranslatedUnitCount,
		ProofreadUnitCount:  chapterDetail.ProofreadUnitCount,
		CoverURL:            chapterDetail.CoverURL,
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

type ListComicChapterArgs struct {
	ComicID string `json:"comic_id"`
	PaginationParams
}

func (args *ListComicChapterArgs) Validate() error {
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
	ChapterNo string `json:"chapter_no"`

}

func (args *UpdateChapterArgs) Validate() error {
	if args == nil {
		return errors.New("参数不能为空")
	}

	if args.ChapterID == "" {
		return errors.New("章节 ID 不能为空")
	}

	if args.ChapterNo == "" {
		return errors.New("章节编号不能为空")
	}

	chapterNoLen := utf8.RuneCountInString(args.ChapterNo)

	if chapterNoLen > 10 {
		return errors.New("章节编号长度不能超过 10 字符")
	}

	return nil
}
