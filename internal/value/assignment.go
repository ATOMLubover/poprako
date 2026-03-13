package value

import (
	"errors"

	"labelplus-next-web-be/internal/domain/model"
	"labelplus-next-web-be/internal/util"
)

// 标准的章节分配信息，可以递归地添加 include 字段
type AssignmentInfo struct {
	ID string `json:"id"`

	// 关联的用户信息
	UserID string    `json:"user_id"`
	User   *UserInfo `json:"user,omitempty"`

	// 关联的章节信息
	ChapterID string       `json:"chapter_id"`
	Chapter   *ChapterInfo `json:"chapter,omitempty"`

	AssignedRawProviderAt int64 `json:"assigned_raw_provider_at"`
	AssignedTranslatorAt  int64 `json:"assigned_translator_at"`
	AssignedProofreaderAt int64 `json:"assigned_proofreader_at"`
	AssignedTypesetterAt  int64 `json:"assigned_typesetter_at"`
	AssignedRedrawerAt    int64 `json:"assigned_redrawer_at"`
	AssignedReviewerAt    int64 `json:"assigned_reviewer_at"`
	AssignedPublisherAt   int64 `json:"assigned_publisher_at"`

	CreatedAt int64 `json:"created_at"`
	UpdatedAt int64 `json:"updated_at"`
}

// ListAssignmentArgs 用于查询当前用户的所有分配，支持分页。
type ListAssignmentArgs struct {
	// 嵌套关联字段
	Includes []string `url:"includes"`

	// 查询条件
	ChapterID string `url:"chapter_id,omitempty"`
	UserID    string `url:"user_id,omitempty"`

	PaginationParams
}

func (a *ListAssignmentArgs) Validate() error {
	if a == nil {
		return errors.New("参数不能为空")
	}

	if err := a.PaginationParams.Validate(); err != nil {
		return err
	}

	if len(a.Includes) > 0 {
	}

	return nil
}

type CreateChapterAssignmentArgs struct {
	ChapterID string         `json:"chapter_id"`
	UserID    string         `json:"user_id"`
	Role      model.RoleMask `json:"role"`
}

func (a *CreateChapterAssignmentArgs) Validate() error {
	if a == nil {
		return errors.New("参数不能为空")
	}

	if a.ChapterID == "" {
		return errors.New("章节 ID 不能为空")
	}

	if a.UserID == "" {
		return errors.New("用户 ID 不能为空")
	}

	if a.Role == 0 {
		return errors.New("分工角色不能为空")
	}

	return nil
}

type CreateChapterAssignmentResult struct {
	ID string `json:"id"`
}

func NewCreateChapterAssignmentResult(id string) CreateChapterAssignmentResult {
	return CreateChapterAssignmentResult{ID: id}
}

// UpdateAssignmentArgs 采用 PUT 语义的全量替换。
type UpdateAssignmentArgs struct {
	ID   string         `json:"id"`
	Role model.RoleMask `json:"role"`
}

func (a *UpdateAssignmentArgs) Validate() error {
	if a == nil {
		return errors.New("参数不能为空")
	}

	if a.ID == "" {
		return errors.New("分配 ID 不能为空")
	}

	if a.Role == 0 {
		return errors.New("分工角色不能为空")
	}

	return nil
}

// ListChapterAssignmentArgs 查询某章节的所有分配
type ListChapterAssignmentArgs struct {
	ChapterID string `url:"chapter_id"`
	PaginationParams
}

func (a *ListChapterAssignmentArgs) Validate() error {
	if a == nil {
		return errors.New("参数不能为空")
	}

	if a.ChapterID == "" {
		return errors.New("章节 ID 不能为空")
	}

	if err := a.PaginationParams.Validate(); err != nil {
		return errors.New("分页参数无效: " + err.Error())
	}

	return nil
}

// AssignmentWithUserInfo 用于列出某章节的所有分配（含用户信息）
type AssignmentWithUserInfo struct {
	ID        string   `json:"id"`
	ChapterID string   `json:"chapter_id"`
	User      UserInfo `json:"user"`

	AssignedRawProviderAt *int64 `json:"assigned_raw_provider_at,omitempty"`
	AssignedTranslatorAt  *int64 `json:"assigned_translator_at,omitempty"`
	AssignedProofreaderAt *int64 `json:"assigned_proofreader_at,omitempty"`
	AssignedTypesetterAt  *int64 `json:"assigned_typesetter_at,omitempty"`
	AssignedReviewerAt    *int64 `json:"assigned_reviewer_at,omitempty"`
	AssignedPublisherAt   *int64 `json:"assigned_publisher_at,omitempty"`

	CreatedAt int64 `json:"created_at"`
	UpdatedAt int64 `json:"updated_at"`
}

func NewAssignmentWithUserInfo(source model.AssignmentWithUserInfo, avatarURL string) AssignmentWithUserInfo {
	return AssignmentWithUserInfo{
		ID:                    source.ID,
		ChapterID:             source.ChapterID,
		User:                  NewUserInfoFromModel(source.User, avatarURL),
		AssignedRawProviderAt: util.ToUnixPtr(source.AssignedRawProviderAt),
		AssignedTranslatorAt:  util.ToUnixPtr(source.AssignedTranslatorAt),
		AssignedProofreaderAt: util.ToUnixPtr(source.AssignedProofreaderAt),
		AssignedTypesetterAt:  util.ToUnixPtr(source.AssignedTypesetterAt),
		AssignedReviewerAt:    util.ToUnixPtr(source.AssignedReviewerAt),
		AssignedPublisherAt:   util.ToUnixPtr(source.AssignedPublisherAt),
		CreatedAt:             source.CreatedAt.UnixMilli(),
		UpdatedAt:             source.UpdatedAt.UnixMilli(),
	}
}

// ChapterWithComicInfo 用于列出分配时嵌套的章节+漫画摘要信息
type ChapterWithComicInfo struct {
	ID    string    `json:"id"`
	Comic ComicInfo `json:"comic"`

	Index     int    `json:"index"`
	ChapterNo string `json:"chapter_no"`

	CoverURL string `json:"cover_url"`

	PageCount           int `json:"page_count"`
	TotalUnitCount      int `json:"total_unit_count"`
	TranslatedUnitCount int `json:"translated_unit_count"`
	ProofreadUnitCount  int `json:"proofread_unit_count"`

	CreatedAt int64 `json:"created_at"`
	UpdatedAt int64 `json:"updated_at"`
}

func NewChapterWithComicInfo(source model.ChapterWithComicInfo, coverURL string) ChapterWithComicInfo {
	return ChapterWithComicInfo{
		ID:                  source.ID,
		Comic:               NewComicInfoFromModel(source.Comic),
		Index:               source.Index,
		ChapterNo:           source.ChapterNo,
		CoverURL:            coverURL,
		PageCount:           source.PageCount,
		TotalUnitCount:      source.TotalUnitCount,
		TranslatedUnitCount: source.TranslatedUnitCount,
		ProofreadUnitCount:  source.ProofreadUnitCount,
		CreatedAt:           source.CreatedAt.UnixMilli(),
		UpdatedAt:           source.UpdatedAt.UnixMilli(),
	}
}

// AssignmentWithChapterInfo 用于列出某用户的所有分配（含章节+漫画信息）
type AssignmentWithChapterInfo struct {
	ID      string               `json:"id"`
	UserID  string               `json:"user_id"`
	Chapter ChapterWithComicInfo `json:"chapter"`

	AssignedRawProviderAt *int64 `json:"assigned_raw_provider_at,omitempty"`
	AssignedTranslatorAt  *int64 `json:"assigned_translator_at,omitempty"`
	AssignedProofreaderAt *int64 `json:"assigned_proofreader_at,omitempty"`
	AssignedTypesetterAt  *int64 `json:"assigned_typesetter_at,omitempty"`
	AssignedReviewerAt    *int64 `json:"assigned_reviewer_at,omitempty"`
	AssignedPublisherAt   *int64 `json:"assigned_publisher_at,omitempty"`

	CreatedAt int64 `json:"created_at"`
	UpdatedAt int64 `json:"updated_at"`
}

func NewAssignmentWithChapterInfo(source model.AssignmentWithChapterInfo, coverURL string) AssignmentWithChapterInfo {
	return AssignmentWithChapterInfo{
		ID:                    source.ID,
		UserID:                source.UserID,
		Chapter:               NewChapterWithComicInfo(source.Chapter, coverURL),
		AssignedRawProviderAt: util.ToUnixPtr(source.AssignedRawProviderAt),
		AssignedTranslatorAt:  util.ToUnixPtr(source.AssignedTranslatorAt),
		AssignedProofreaderAt: util.ToUnixPtr(source.AssignedProofreaderAt),
		AssignedTypesetterAt:  util.ToUnixPtr(source.AssignedTypesetterAt),
		AssignedReviewerAt:    util.ToUnixPtr(source.AssignedReviewerAt),
		AssignedPublisherAt:   util.ToUnixPtr(source.AssignedPublisherAt),
		CreatedAt:             source.CreatedAt.UnixMilli(),
		UpdatedAt:             source.UpdatedAt.UnixMilli(),
	}
}
