package service

import (
	"time"

	"labelplus-next-web-be/internal/domain/model"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	GenerateAccessToken(
		userID string,
		secretKey []byte,
		expirationHours int,
	) (string, error)
	ParseAccessToken(
		tokenString string,
		secretKey []byte,
	) (*model.TokenClaims, error)
	VerifyPassword(
		plainPassword,
		hashedPassword string,
	) bool
	HashPassword(
		plainPassword string,
	) (string, error)
}

type userService struct{}

func NewUserService() UserService {
	return &userService{}
}

func (us *userService) GenerateAccessToken(
	userID string,
	secretKey []byte,
	expirationHours int,
) (string, error) {
	claims := &model.TokenClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expirationHours) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "labelplus-next-web",
			Subject:   "user-auth",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		zap.L().Error("GenerateAccessToken: 生成访问令牌失败", zap.Error(err))
		return "", err
	}

	return tokenString, nil
}

func (us *userService) ParseAccessToken(
	tokenString string,
	secretKey []byte,
) (*model.TokenClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&model.TokenClaims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return secretKey, nil
		})
	if err != nil {
		zap.L().Error("ParseAccessToken: 解析访问令牌失败", zap.Error(err))
		return nil, err
	}

	claims, ok := token.Claims.(*model.TokenClaims)
	if !ok || !token.Valid {
		zap.L().Error("ParseAccessToken: 无效的访问令牌")
		return nil, jwt.ErrInvalidKey
	}

	return claims, nil
}

func (us *userService) VerifyPassword(plainPassword, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))

	return err == nil
}

func (us *userService) HashPassword(plainPassword string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		zap.L().Error("HashPassword: 哈希密码失败", zap.Error(err))
		return "", err
	}

	return string(hash), nil
}
