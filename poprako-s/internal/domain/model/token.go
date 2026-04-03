package model

import "github.com/golang-jwt/jwt/v5"

// UserClaims 是 JWT 载荷中包含的用户信息
type UserClaims struct {
	jwt.RegisteredClaims
	UserID string `json:"user_id"`
}
