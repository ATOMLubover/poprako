package token_infra

import "github.com/golang-jwt/jwt/v5"

type UserClaims struct {
	// Embedded for standard claims like exp, iat, etc.
	jwt.RegisteredClaims
	UserId string `json:"user_id"`
}
