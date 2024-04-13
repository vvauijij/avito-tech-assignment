package token

import (
	"crypto/rsa"

	"github.com/golang-jwt/jwt/v5"
)

type (
	TokenGenerator struct {
		JWTPrivate *rsa.PrivateKey
	}

	UserClaims struct {
		*jwt.RegisteredClaims
		AccessRight string
	}
)
