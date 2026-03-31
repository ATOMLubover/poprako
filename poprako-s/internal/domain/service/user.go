package service

import (
	"strings"

	"poprako-s/internal/domain/model"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

// UserService 定义用户领域相关的业务能力
type UserService interface {
	// ParseToken 解析 JWT token 字符串，返回用户声明信息
	ParseToken(tokenStr string, secretKey []byte) (*model.UserClaims, error)
	// HashPwd 将明文密码哈希后返回
	HashPwd(pwd string) (string, error)
	// NewCreation 根据业务参数创建 UserCreation 领域模型（ID 由 service 内部生成）
	NewCreation(name, qq, pwdHash string) *model.UserCreation
	// GenAvatarOSSKey 根据用户 ID 生成头像的 OSS Key
	GenAvatarOSSKey(userID string) string
}

// userServiceImpl 是 UserService 的具体实现
type userServiceImpl struct {
	lgr zap.Logger
}

// NewUserService 返回 UserService 的默认实现
func NewUserService(l zap.Logger) UserService {
	return &userServiceImpl{
		lgr: l,
	}
}

// ParseToken 解析 JWT token 字符串，校验签名算法并返回用户声明
func (s *userServiceImpl) ParseToken(
	tokenStr string,
	secretKey []byte,
) (*model.UserClaims, error) {
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
		return nil, err
	}

	claims, ok := token.Claims.(*model.UserClaims)
	if !ok || !token.Valid {
		return nil, jwt.ErrInvalidKey
	}

	return claims, nil
}

// HashPwd 使用 bcrypt 对明文密码进行哈希
func (s *userServiceImpl) HashPwd(pwd string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

// NewCreation 构造一个带有 service 生成 ID 的 UserCreation
func (s *userServiceImpl) NewCreation(name, qq, pwdHash string) *model.UserCreation {
	return &model.UserCreation{
		ID:      GenID("user"),
		Name:    name,
		QQ:      qq,
		PwdHash: pwdHash,
	}
}

// GenAvatarOSSKey 以固定前缀拼接用户 ID 作为头像的 OSS Key
func (s *userServiceImpl) GenAvatarOSSKey(userID string) string {
	return strings.Join([]string{"user-avatar", userID}, "_")
}
