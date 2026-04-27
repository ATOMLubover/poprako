package token_infra

import (
	"os"
	"strconv"
	"time"

	token_iface "poprako-s/internal/domain/ext/token"
	"poprako-s/internal/domain/model/aggr"

	"github.com/golang-jwt/jwt/v5"
)

type jwtParser struct {
	secret []byte
	exp    time.Duration
}

func NewJwtParser() token_iface.Parser {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		panic("[NewJwtParser] JWT_SECRET env variable is not set")
	}

	expStr := os.Getenv("JWT_EXPIRATION_HOURS")
	if expStr == "" {
		panic("[NewJwtParser] JWT_EXPIRATION_HOURS env variable is not set")
	}

	expHrs, err := strconv.Atoi(expStr)
	if err != nil {
		panic("[NewJwtParser] JWT_EXPIRATION_HOURS env variable must be a valid integer")
	}

	return &jwtParser{
		secret: []byte(secret),
		exp:    time.Duration(expHrs) * time.Hour,
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
