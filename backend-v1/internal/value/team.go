package value

import (
	"errors"
	"unicode/utf8"
)

type CreateTeamArgs struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (cta *CreateTeamArgs) Validate() error {
	if cta == nil {
		return errors.New("参数不能为空")
	}

	nameLen := utf8.RuneCountInString(cta.Name)
	if nameLen <= 0 || nameLen > 20 {
		return errors.New("汉化组名称长度必须在 1~20 字符之间")
	}

	descriptionLen := utf8.RuneCountInString(cta.Description)
	if descriptionLen > 100 {
		return errors.New("汉化组描述长度不能超过 100 字符")
	}

	return nil
}

type CreateTeamResult struct {
	ID string `json:"id"`
}

type TeamInfo struct {
	ID string `json:"id"`

	Name             string `json:"name"`
	Description      string `json:"description"`
	AvatarURL        string `json:"avatar_url"`
	IsAvatarUploaded bool   `json:"is_avatar_uploaded"`

	CreatedAt int64 `json:"created_at"`
	UpdatedAt int64 `json:"updated_at"`
}

type ReserveTeamAvatarArgs struct {
	Extension string `json:"extension"`
}

func (rtaa *ReserveTeamAvatarArgs) Validate() error {
	if rtaa == nil {
		return errors.New("参数不能为空")
	}

	if rtaa.Extension == "" {
		return errors.New("文件扩展名不能为空")
	}

	if rtaa.Extension != "jpg" && rtaa.Extension != "jpeg" && rtaa.Extension != "png" && rtaa.Extension != "webp" {
		return errors.New("不支持的文件扩展名，仅支持 jpg/jpeg、png 和 webp")
	}

	return nil
}

type ReserveTeamAvatarResult struct {
	PutURL string `json:"put_url"`
}

type UpdateTeamArgs struct {
	ID string

	Name        string `json:"name"`
	Description string `json:"description"`
}

func (uta *UpdateTeamArgs) Validate() error {
	if uta == nil {
		return errors.New("参数不能为空")
	}

	if uta.ID == "" {
		return errors.New("汉化组 ID 不能为空")
	}

	nameLen := utf8.RuneCountInString(uta.Name)
	if nameLen <= 0 || nameLen > 20 {
		return errors.New("汉化组名称长度必须在 1~20 字符之间")
	}

	descriptionLen := utf8.RuneCountInString(uta.Description)
	if descriptionLen > 100 {
		return errors.New("汉化组描述长度不能超过 100 字符")
	}

	return nil
}

type ListTeamArgs struct {
	PaginationParams
}

func (args *ListTeamArgs) Validate() error {
	if args == nil {
		return errors.New("参数不能为空")
	}

	if err := args.PaginationParams.Validate(); err != nil {
		return errors.New("分页参数无效: " + err.Error())
	}

	return nil
}

type ListMyTeamArgs struct {
	PaginationParams
}

func (args *ListMyTeamArgs) Validate() error {
	if args == nil {
		return errors.New("参数不能为空")
	}

	if err := args.PaginationParams.Validate(); err != nil {
		return errors.New("分页参数无效: " + err.Error())
	}

	return nil
}
