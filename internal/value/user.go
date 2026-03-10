package value

import (
	"errors"
	"unicode/utf8"

	"labelplus-next-web-be/internal/domain/model"
	"labelplus-next-web-be/internal/util"

	"go.uber.org/zap"
)

type LoginUserArgs struct {
	QQ       string `json:"qq"`
	Password string `json:"password"`
}

func (lua *LoginUserArgs) Validate() error {
	if lua == nil {
		return errors.New("参数不能为空")
	}

	if lua.QQ == "" {
		return errors.New("QQ 不能为空")
	}

	if lua.Password == "" {
		return errors.New("密码不能为空")
	}

	return nil
}

type LoginUserResult struct {
	UserID      string `json:"user_id"`
	AccessToken string `json:"access_token"`
}

func NewLoginUserResult(userID string, accessToken string) LoginUserResult {
	return LoginUserResult{
		UserID:      userID,
		AccessToken: accessToken,
	}
}

type RegisterUserArgs struct {
	QQ             string `json:"qq"`
	Password       string `json:"password"`
	Name           string `json:"name"`
	InvitationCode string `json:"invitation_code"`
}

func (rua *RegisterUserArgs) Validate() error {
	if rua == nil {
		return errors.New("参数不能为空")
	}

	qqLen := utf8.RuneCountInString(rua.QQ)

	if rua.QQ == "" || qqLen < 5 || qqLen > 20 {
		return errors.New("QQ 长度必须在 5 到 20 个字符之间")
	}

	passwordLen := utf8.RuneCountInString(rua.Password)

	if rua.Password == "" || passwordLen < 6 || passwordLen > 30 {
		return errors.New("密码长度必须在 6 到 30 个字符之间")
	}

	nameLen := utf8.RuneCountInString(rua.Name)

	if rua.Name == "" || nameLen < 2 || nameLen > 20 {
		return errors.New("名字长度必须在 2 到 20 个字之间")
	}

	if rua.InvitationCode == "" {
		return errors.New("邀请码不能为空")
	}

	return nil
}

type RegisterUserResult struct {
	UserID      string `json:"user_id"`
	AccessToken string `json:"access_token"`
}

func NewRegisterUserResult(userID string, accessToken string) RegisterUserResult {
	return RegisterUserResult{
		UserID:      userID,
		AccessToken: accessToken,
	}
}

// UserInfo 是应用层对外暴露的用户信息值对象
type UserInfo struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	QQ   string `json:"qq"`

	AvatarURL string `json:"avatar_url"`

	CreatedAt int64 `json:"created_at"`
	UpdatedAt int64 `json:"updated_at"`
}

func NewUserInfoFromModel(user model.UserInfo) UserInfo {
	if user.ID == "" {
		zap.L().Warn("NewUserInfoFromModel: user 为空")
		return UserInfo{}
	}

	return UserInfo{
		ID:        user.ID,
		Name:      user.Name,
		QQ:        user.QQ,
		AvatarURL: user.AvatarURL,
		CreatedAt: user.CreatedAt.UnixMilli(),
		UpdatedAt: user.UpdatedAt.UnixMilli(),
	}
}

type ListUserArgs struct {
	PaginationParams

	QQ        string `json:"qq"`
	FuzzyName string `json:"fuzzy_name"`
}

func (lua *ListUserArgs) Validate() error {
	if lua == nil {
		return errors.New("参数不能为空")
	}

	if lua.Offset < 0 {
		return errors.New("offset 不能为负数")
	}

	if lua.Limit <= 0 || lua.Limit > 30 {
		return errors.New("limit 必须大于 0 且不超过 30")
	}

	return nil
}

type UpdateUserArgs struct {
	UserID string `json:"user_id"`

	Name     util.Option[string] `json:"name"`
	QQ       util.Option[string] `json:"qq"`
	Password util.Option[string] `json:"password"`
}

func (uua *UpdateUserArgs) Validate() error {
	if uua == nil {
		return errors.New("参数不能为空")
	}

	if uua.Name.State() == util.OptionSome {
		nameLen := utf8.RuneCountInString(uua.Name.Unwrap())

		if nameLen < 2 || nameLen > 20 {
			return errors.New("名字长度必须在 2 到 20 个字之间")
		}
	}

	if uua.QQ.State() == util.OptionSome {
		if uua.QQ.Unwrap() == "" {
			return errors.New("QQ 不能为空")
		}

		qqLen := utf8.RuneCountInString(uua.QQ.Unwrap())

		if qqLen < 5 || qqLen > 20 {
			return errors.New("QQ 长度必须在 5 到 20 个字符之间")
		}
	}

	if uua.Password.State() == util.OptionSome {
		if uua.Password.Unwrap() == "" {
			return errors.New("密码不能为空，或包含非数字字母的特殊字符")
		}

		nameLen := utf8.RuneCountInString(uua.Password.Unwrap())

		if nameLen < 6 || nameLen > 30 {
			return errors.New("密码长度必须在 6 到 30 个字符之间")
		}
	}

	return nil
}
