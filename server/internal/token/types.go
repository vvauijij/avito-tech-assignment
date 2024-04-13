package token

import "github.com/golang-jwt/jwt/v5"

type AccessRight int

const (
	None AccessRight = iota
	User
	Admin
)

type (
	Client interface {
		Verify(token string) AccessRight
	}

	UserClaims struct {
		*jwt.RegisteredClaims
		AccessRight string
	}
)
