package value

import (
	"errors"
	"unicode/utf8"

	"labelplus-next-web-be/internal/domain/model"

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

	if rua.QQ == "" {
		return errors.New("QQ 不能为空")
	}

	if rua.Password == "" {
		return errors.New("密码不能为空")
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

func NewRegisterUserResult(userID string, accessToken string) *RegisterUserResult {
	return &RegisterUserResult{
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

	Roles model.RoleMask `json:"roles"`

	CreatedAt int64 `json:"created_at"`
	UpdatedAt int64 `json:"updated_at"`
}

func NewUserInfoFromModel(user *model.UserInfo) *UserInfo {
	if user == nil {
		zap.L().Warn("NewUserInfoFromModel: user 为空")
		return nil
	}

	return &UserInfo{
		ID:        user.ID,
		Name:      user.Name,
		QQ:        user.QQ,
		AvatarURL: user.AvatarURL,
		Roles:     user.MaskRoles(),
		CreatedAt: user.CreatedAt.UnixMilli(),
		UpdatedAt: user.UpdatedAt.UnixMilli(),
	}
}

type ListUsersArgs struct {
	PaginationParams

	QQ        string `json:"qq"`
	FuzzyName string `json:"fuzzy_name"`
	// 仅支持单选
	Role model.RoleFlag `json:"role"`
}

func (lua *ListUsersArgs) Validate() error {
	if lua == nil {
		return errors.New("参数不能为空")
	}

	if lua.Offset < 0 {
		return errors.New("offset 不能为负数")
	}

	if lua.Limit <= 0 {
		return errors.New("limit 必须大于 0")
	}

	if lua.Role != 0 && (lua.Role != model.RolePictureSource &&
		lua.Role != model.RoleTranslator &&
		lua.Role != model.RoleProofreader &&
		lua.Role != model.RoleTypesetter &&
		lua.Role != model.RoleReviewer &&
		lua.Role != model.RoleAdmin &&
		lua.Role != model.RoleSuperAdmin) {
		return errors.New("role 参数无效")
	}

	return nil
}
