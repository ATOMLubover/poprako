package value

import (
	"errors"
	"unicode/utf8"

	"labelplus-next-web-be/internal/domain/model"
	"labelplus-next-web-be/internal/util"

	"go.uber.org/zap"
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
		return errors.New("团队名称长度必须在 1~20 字符之间")
	}

	descriptionLen := utf8.RuneCountInString(cta.Description)
	if descriptionLen > 100 {
		return errors.New("团队描述长度不能超过 100 字符")
	}

	return nil
}

type TeamInfo struct {
	ID string `json:"id"`

	Name        string `json:"name"`
	Description string `json:"description"`

	CreatedAt int64 `json:"created_at"`
	UpdatedAt int64 `json:"updated_at"`
}

func NewTeamInfoFromModel(team *model.TeamInfo) *TeamInfo {
	if team == nil {
		zap.L().Warn("NewTeamInfoFromModel: team 为空")
		return nil
	}

	return &TeamInfo{
		ID:          team.ID,
		Name:        team.Name,
		Description: team.Description,
		CreatedAt:   team.CreatedAt.UnixMilli(),
		UpdatedAt:   team.UpdatedAt.UnixMilli(),
	}
}

type UpdateTeamArgs struct {
	ID string

	Name        util.Option[string] `json:"name"`
	Description util.Option[string] `json:"description"`
}

func (uta *UpdateTeamArgs) Validate() error {
	if uta == nil {
		return errors.New("参数不能为空")
	}

	if uta.Name.State() == util.OptionSome {
		nameLen := utf8.RuneCountInString(uta.Name.Unwrap())
		if nameLen <= 0 || nameLen > 20 {
			return errors.New("团队名称长度必须在 1~20 字符之间")
		}
	}

	if uta.Description.State() == util.OptionSome {
		descriptionLen := utf8.RuneCountInString(uta.Description.Unwrap())

		if descriptionLen > 100 {
			return errors.New("团队描述长度不能超过 100 字符")
		}
	}

	return nil
}
