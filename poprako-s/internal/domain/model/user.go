package model

import (
	"time"

	"poprako-s/internal/domain/event"

	"golang.org/x/crypto/bcrypt"

	"github.com/golang-jwt/jwt/v5"
)

// UserInfo 是用户信息的 view，包含用户的基本信息
type UserInfo struct {
	ID string

	Name string
	QQ   string

	AvatarKey        string
	IsAvatarUploaded bool

	// 用户是否是 **整个系统** 的超级管理员
	IsSuperAdmin bool

	// 在刚注册时，LastLogin 应该是注册时间
	// 这个值与 UpdatedAt 不同，
	// UpdatedAt 是用户信息最后一次被修改的时间，而 LastLogin 是用户最后一次登录的时间
	LastLoginAt time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
}

// GenToken 生成一个 token（算法可变），用于校验用户身份
func (i *UserInfo) GenToken(
	secretKey string,
	expHrs int,
) (string, error) {
	now := time.Now()
	exp := now.Add(time.Duration(expHrs) * time.Hour)

	c := &UserClaims{
		UserID: i.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
			IssuedAt:  jwt.NewNumericDate(now),
			Issuer:    "poprako-s",
			Subject:   "poprako-user-auth",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)

	return token.SignedString([]byte(secretKey))
}

// UserCreds 是用户凭据的 view
type UserCreds struct {
	QQ      string
	PwdHash string

	event.EventBase
}

// 校验密码是否正确，如果正确则发布 UserLoginEvent 领域事件，且返回 nil
func (c *UserCreds) Authenticate(pwd string) error {
	// 比较传入的密码与存储的密码哈希，验证用户身份
	if err := bcrypt.CompareHashAndPassword(
		[]byte(c.PwdHash),
		[]byte(pwd),
	); err != nil {
		return err
	}

	// 验证通过，发布 UserLoginEvent 领域事件
	c.PushEvent(&event.UserLoginEvent{
		UserQQ: c.QQ,
	})

	return nil
}

// UserReg 是用户注册信息的 view
type UserCreation struct {
	// ID 由 domain service 生成，外部不提供
	ID string

	Name    string
	QQ      string
	PwdHash string

	event.EventBase
}

// UserUpdate 是用户更新信息的 view，是 PUT 语义的载荷
// 更改密码不走这个结构体
type UserUpdate struct {
	ID string

	// 以下字段都是必填
	Name string
	QQ   string
}

// UserQueryOpt 指定用户查询的可选筛选条件
// 所有字段均为可空，nil 表示不参与筛选
type UserQueryOpt struct {
	// ID 按用户 ID 筛选
	ID *string
	// QQ 按 QQ 号筛选
	QQ *string
}
