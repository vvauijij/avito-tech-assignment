package token

import (
	"crypto/rsa"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/vvauijij/avito-tech-assignment/server/internal/env"
)

type ClientImpl struct {
	JWTPublic *rsa.PublicKey
}

func NewClient(publicFile string) Client {
	absolutePublicFile, err := filepath.Abs(publicFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	public, err := os.ReadFile(absolutePublicFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	JWTPublic, err := jwt.ParseRSAPublicKeyFromPEM(public)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	return &ClientImpl{JWTPublic: JWTPublic}
}

func (c *ClientImpl) Verify(token string) AccessRight {
	if env.IsTestENV() && token == "test_token" {
		return Admin
	}

	userClaims := UserClaims{}
	_, err := jwt.ParseWithClaims(
		token,
		&userClaims,
		func(token *jwt.Token) (interface{}, error) {
			return c.JWTPublic, nil
		},
	)
	if err != nil || userClaims.ExpiresAt.Compare(time.Now()) < 0 {
		return None
	}

	switch userClaims.AccessRight {
	case "user":
		return User
	case "admin":
		return Admin
	default:
		return None
	}
}
