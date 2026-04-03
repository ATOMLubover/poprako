package service

import (
	"strings"

	"poprako-s/internal/domain/event"
	"poprako-s/internal/domain/model"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// UserService 定义用户领域相关的业务能力
type UserService interface {
	// ParseToken 解析 JWT token 字符串，校验签名算法并返回用户声明信息
	ParseToken(
		tokenStr string,
		secretKey []byte,
	) (*model.UserClaims, error)

	// HashPwd 将明文密码用 bcrypt 哈希后返回
	HashPwd(
		pwd string,
	) (string, error)

	// NewCreation 根据业务参数创建 UserCreation 领域模型，ID 由 service 内部生成
	NewCreation(
		name string,
		pwd string,
		i *model.InvitationInfo,
	) (*model.UserCreation, error)

	// GenAvatarOSSKey 根据用户 ID 生成头像的 OSS Key
	GenAvatarOSSKey(
		userID string,
	) string
}

// userServiceImpl 是 UserService 的具体实现，无内禀状态
type userServiceImpl struct{}

// NewUserService 返回 UserService 的默认实现
func NewUserService() UserService {
	// 返回无状态实现
	return &userServiceImpl{}
}

// NewCreation 构造一个带 service 生成 ID 的 UserCreation
func (s *userServiceImpl) NewCreation(
	name string,
	pwd string,
	i *model.InvitationInfo,
) (*model.UserCreation, error) {
	// 生成唯一用户 ID
	id := GenID("user")

	// 对明文密码进行哈希处理
	pwdHash, err := s.HashPwd(pwd)
	if err != nil {
		// 直接返回哈希失败错误
		return nil, err
	}

	// 构造用户创建载荷
	c := &model.UserCreation{
		ID:      id,
		Name:    name,
		QQ:      i.InviteeQQ,
		PwdHash: pwdHash,
	}

	// 追加用户创建领域事件
	c.PushEvent(&event.UserCreatedEvent{
		InvitorID:     i.InvitorID,
		CreatedUserID: id,
	})

	// 返回构造结果
	return c, nil
}

// ParseToken 解析 JWT token 字符串，校验签名算法并返回用户声明
func (s *userServiceImpl) ParseToken(
	tokenStr string,
	secretKey []byte,
) (*model.UserClaims, error) {
	// 解析 token 并校验签名方法
	token, err := jwt.ParseWithClaims(
		tokenStr,
		&model.UserClaims{},
		func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}

			return secretKey, nil
		},
	)
	if err != nil {
		// 返回解析或签名校验错误
		return nil, err
	}

	// 断言 claims 类型并校验 token 有效性
	claims, ok := token.Claims.(*model.UserClaims)
	if !ok || !token.Valid {
		// 返回 claims 类型断言失败错误
		return nil, jwt.ErrInvalidKey
	}

	// 返回校验通过的声明
	return claims, nil
}

// HashPwd 使用 bcrypt 对明文密码进行哈希
func (s *userServiceImpl) HashPwd(
	pwd string,
) (string, error) {
	// 使用 bcrypt 默认代价因子生成哈希
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		// 返回哈希生成失败错误
		return "", err
	}

	// 返回哈希字符串
	return string(hash), nil
}

// GenAvatarOSSKey 以固定前缀拼接用户 ID 作为头像的 OSS Key
func (s *userServiceImpl) GenAvatarOSSKey(
	userID string,
) string {
	// 返回拼接后的对象 Key
	return strings.Join([]string{"user-avatar", userID}, "_")
}
