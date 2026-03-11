package value

import (
	"errors"

	"labelplus-next-web-be/internal/domain/model"
)

type ReserveChapterPagesArgs struct {
	ChapterID string `json:"chapter_id"`
	PageCount int    `json:"page_count"`
}

func (args *ReserveChapterPagesArgs) Validate() error {
	if args == nil {
		return errors.New("参数不能为空")
	}

	if args.ChapterID == "" {
		return errors.New("章节 ID 不能为空")
	}

	if args.PageCount <= 0 {
		return errors.New("页面数量必须大于 0")
	}

	return nil
}

type ReserveChapterPagesResult struct {
	Creations []PageCreationResult `json:"creations"`
}

func NewCreateChapterPagesResult(creations []PageCreationResult) ReserveChapterPagesResult {
	return ReserveChapterPagesResult{
		Creations: creations,
	}
}

type PageCreationResult struct {
	PageID string `json:"page_id"`
	PutURL string `json:"put_url"`
}

func NewPageCreationResult(pageID string, putURL string) PageCreationResult {
	return PageCreationResult{
		PageID: pageID,
		PutURL: putURL,
	}
}

type UpdatePageArgs struct {
	ID string `json:"id"`
	// 由于采用预签名方式让客户端上传，所以在客户端上传后，
	// 必须让其主动调用一次更新接口来告诉后端该页面已经上传完成了
	IsUploaded bool `json:"is_uploaded"`
}

func (args *UpdatePageArgs) Validate() error {
	if args == nil {
		return errors.New("参数不能为空")
	}

	if args.ID == "" {
		return errors.New("页面 ID 不能为空")
	}

	return nil
}

type PageInfo struct {
	ID string `json:"id"`

	ChapterID string `json:"chapter_id"`
	Index     int    `json:"index"`

	ImageURL string `json:"image_url"`

	TotalUnitCount      int `json:"total_unit_count"`
	TranslatedUnitCount int `json:"translated_unit_count"`
	ProofreadUnitCount  int `json:"proofread_unit_count"`

	CreatedAt int64 `json:"created_at"`
	UpdatedAt int64 `json:"updated_at"`
}

func NewPageInfoFromModel(pageInfo model.PageInfo, imageURL string) PageInfo {
	return PageInfo{
		ID:                  pageInfo.ID,
		ChapterID:           pageInfo.ChapterID,
		Index:               pageInfo.Index,
		ImageURL:            imageURL,
		TotalUnitCount:      pageInfo.TotalUnitCount,
		TranslatedUnitCount: pageInfo.TranslatedUnitCount,
		ProofreadUnitCount:  pageInfo.ProofreadUnitCount,
		CreatedAt:           pageInfo.CreatedAt.Unix(),
		UpdatedAt:           pageInfo.UpdatedAt.Unix(),
	}
}
