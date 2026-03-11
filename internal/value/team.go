package value

import (
	"errors"
	"unicode/utf8"

	"labelplus-next-web-be/internal/domain/model"

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

func NewCreateTeamResult(teamID string) CreateTeamResult {
	return CreateTeamResult{
		ID: teamID,
	}
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

func NewTeamInfoFromModel(team model.TeamInfo, avatarURL string) TeamInfo {
	if team.ID == "" {
		zap.L().Warn("NewTeamInfoFromModel: team 为空")
		return TeamInfo{}
	}

	return TeamInfo{
		ID:               team.ID,
		Name:             team.Name,
		Description:      team.Description,
		AvatarURL:        avatarURL,
		IsAvatarUploaded: team.IsAvatarUploaded,
		CreatedAt:        team.CreatedAt.UnixMilli(),
		UpdatedAt:        team.UpdatedAt.UnixMilli(),
	}
}

type ReserveTeamAvatarResult struct {
	AvatarOSSKey string `json:"avatar_oss_key"`
	PutURL       string `json:"put_url"`
}

func NewReserveTeamAvatarResult(avatarOSSKey string, putURL string) ReserveTeamAvatarResult {
	return ReserveTeamAvatarResult{
		AvatarOSSKey: avatarOSSKey,
		PutURL:       putURL,
	}
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
