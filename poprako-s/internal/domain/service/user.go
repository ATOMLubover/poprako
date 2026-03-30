package service

import (
	"poprako-s/internal/domain/model"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

// UserService 是用户服务的接口，定义了用户相关的业务逻辑
// 所有无法归纳到 **单个聚合根** 的业务逻辑都应该放在这里
// 理论上它应该是无状态的
type UserService interface{}

// userServiceImpl 是 UserService 的具体实现
// 为了方便起见，依然需要注入一个 logger 来记录日志，这个 logger 携带上下文信息
type userServiceImpl struct {
	lgr zap.Logger
}

func NewUserService(l zap.Logger) UserService {
	return &userServiceImpl{
		lgr: l,
	}
}

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

func (s *userServiceImpl) HashPwd(pwd string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}
