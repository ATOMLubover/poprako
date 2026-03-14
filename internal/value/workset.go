package value

import (
	"errors"
)

type WorksetInfo struct {
	ID string `json:"id"`

	TeamID string    `json:"team_id"`
	Team   *TeamInfo `json:"team,omitempty"`
	Index  int       `json:"index"`

	Name        string `json:"name"`
	Description string `json:"description"`
	ComicCount  int    `json:"comic_count"`

	CreatedAt int64 `json:"created_at"`
	UpdatedAt int64 `json:"updated_at"`
}

// 所有的 list 都聚合在一个 List 参数下
type ListWorksetArgs struct {
	TeamID   string   `url:"team_id"`
	Includes []string `url:"includes[]"` // 可选，指定是否包含关联信息，如 team
	PaginationParams
}

func (lwa *ListWorksetArgs) Validate() error {
	if lwa == nil {
		return errors.New("参数不能为空")
	}

	if lwa.TeamID == "" {
		return errors.New("汉化组 ID 不能为空")
	}

	if err := lwa.PaginationParams.Validate(); err != nil {
		return errors.New("分页参数无效: " + err.Error())
	}

	return nil
}

type CreateWorksetArgs struct {
	TeamID      string  `json:"team_id"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
}

func (cwa *CreateWorksetArgs) Validate() error {
	if cwa == nil {
		return errors.New("参数不能为空")
	}

	if cwa.TeamID == "" {
		return errors.New("汉化组 ID 不能为空")
	}

	if cwa.Name == "" {
		return errors.New("工作集名称不能为空")
	}

	return nil
}

type CreateWorksetResult struct {
	ID string `json:"id"`
}

// PUT 语义的 update 参数，nil 值视为强制置为 NULL，非 nil 值视为更新为该值，空字符串视为更新为 ""
type UpdateWorksetArgs struct {
	ID          string  `json:"id"`
	Name        string  `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
}

func (uwa *UpdateWorksetArgs) Validate() error {
	if uwa == nil {
		return errors.New("参数不能为空")
	}

	if uwa.ID == "" {
		return errors.New("工作集 ID 不能为空")
	}

	return nil
}
