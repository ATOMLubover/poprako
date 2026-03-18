package value

import (
	"errors"
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

type PageCreationResult struct {
	PageID string `json:"page_id"`
	PutURL string `json:"put_url"`
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

	ChapterID string    `json:"chapter_id"`
	Index     int       `json:"index"`
	CreatorID string    `json:"creator_id"`
	Creator   *UserInfo `json:"creator,omitempty"`

	ImageURL string `json:"image_url"`

	TotalUnitCount      int `json:"total_unit_count"`
	TranslatedUnitCount int `json:"translated_unit_count"`
	ProofreadUnitCount  int `json:"proofread_unit_count"`

	CreatedAt int64 `json:"created_at"`
	UpdatedAt int64 `json:"updated_at"`
}

type ListChapterPageArgs struct {
	ChapterID string   `url:"chapter_id"`
	Includes  []string `url:"includes[]"`
	PaginationParams
}

func (args *ListChapterPageArgs) Validate() error {
	if args == nil {
		return errors.New("参数不能为空")
	}

	if args.ChapterID == "" {
		return errors.New("章节 ID 不能为空")
	}

	if err := args.PaginationParams.Validate(); err != nil {
		return errors.New("分页参数无效: " + err.Error())
	}

	return nil
}
