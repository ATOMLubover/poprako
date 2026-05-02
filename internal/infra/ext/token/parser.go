package token_infra

import (
	"time"

	token_iface "poprako-s/internal/domain/ext/token"
	"poprako-s/internal/domain/model/aggr"

	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

type jwtParser struct {
	secret []byte
	exp    time.Duration
}

func NewJwtParser() token_iface.Parser {
	secret := viper.GetString("JWT_SECRET")
	if secret == "" {
		panic("[NewJwtParser] JWT_SECRET env variable is not set")
	}

	exp := viper.GetInt("JWT_EXPIRATION_HOURS")

	return &jwtParser{
		secret: []byte(secret),
		exp:    time.Duration(exp) * time.Hour,
	}
}

func (p *jwtParser) GenToken(raw *aggr.UserToken) (token_iface.SignedToken, error) {
	const ISSUER = "poprako-s"
	const SUBJECT = "user-authorization"

	now := time.Now()

	cl := &UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(p.exp)),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    ISSUER,
			Subject:   SUBJECT,
		},
		UserId: raw.UserId,
	}

	tk := jwt.NewWithClaims(jwt.SigningMethodHS256, cl)

	sgn, err := tk.SignedString(p.secret)
	if err != nil {
		return "", err
	}

	return token_iface.SignedToken(sgn), nil
}

func (p *jwtParser) ParseToken(sgn token_iface.SignedToken) (*aggr.UserToken, error) {
	tk, err := jwt.ParseWithClaims(
		string(sgn),
		&UserClaims{},
		func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return p.secret, nil
		},
	)
	if err != nil {
		return nil, err
	}

	cl, ok := tk.Claims.(*UserClaims)
	if !ok || !tk.Valid {
		return nil, jwt.ErrInvalidKey
	}

	return &aggr.UserToken{
		UserId: cl.UserId,
	}, nil
}
