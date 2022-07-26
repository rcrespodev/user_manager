package domain

import "github.com/golang-jwt/jwt/v4"

type JwtConfig struct {
	secret     string
	claimsKey  [2]string
	signMethod jwt.SigningMethod
}

func NewJwtConfig(secret string) *JwtConfig {
	return &JwtConfig{
		secret:     secret,
		claimsKey:  [2]string{"uuid", "time"},
		signMethod: jwt.SigningMethodHS256,
	}
}
