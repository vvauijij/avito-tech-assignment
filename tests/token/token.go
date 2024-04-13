package token

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var Generator = NewTokenGenerator(os.Getenv("PRIVATE"))

func NewTokenGenerator(privateFile string) *TokenGenerator {
	absoluteprivateFile, err := filepath.Abs(privateFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	private, err := os.ReadFile(absoluteprivateFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	JWTPrivate, err := jwt.ParseRSAPrivateKeyFromPEM(private)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	return &TokenGenerator{
		JWTPrivate: JWTPrivate,
	}
}

func (g *TokenGenerator) NewUserToken() string {
	t := jwt.New(jwt.GetSigningMethod("RS256"))
	t.Claims = &UserClaims{
		RegisteredClaims: &jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 60)),
		},
		AccessRight: "user",
	}
	JWTToken, _ := t.SignedString(g.JWTPrivate)
	return JWTToken
}

func (g *TokenGenerator) NewAdminToken() string {
	t := jwt.New(jwt.GetSigningMethod("RS256"))
	t.Claims = &UserClaims{
		RegisteredClaims: &jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 60)),
		},
		AccessRight: "admin",
	}
	JWTToken, _ := t.SignedString(g.JWTPrivate)
	return JWTToken
}

func (g *TokenGenerator) NewInvalidToken() string {
	return "invalid_token"
}
