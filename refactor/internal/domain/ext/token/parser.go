package token_iface

import "poprako-s/internal/domain/model/aggr"

// `SignedToken` is a type alias for a string that represents a signed user token.
type SignedToken string

// `Parser` is an interface for signing and parsing user tokens.
type Parser interface {
	// `GenToken` generates a signed token string from the given user token.
	GenToken(raw *aggr.UserToken) (SignedToken, error)

	// `ParseToken` parses a signed token string and returns the corresponding user token.
	ParseToken(sgn SignedToken) (*aggr.UserToken, error)
}
