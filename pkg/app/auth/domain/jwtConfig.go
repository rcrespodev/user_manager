package domain

import (
	"github.com/golang-jwt/jwt/v4"
	"strconv"
	"time"
)

type JwtConfig struct {
	secret         string
	claimsKey      [2]string
	signMethod     jwt.SigningMethod
	expirationTime time.Duration
}

func NewJwtConfig(secret, expirationTime string) *JwtConfig {
	var tokenMinutesDuration time.Duration
	duration, err := strconv.Atoi(expirationTime)
	if err != nil {
		duration = 15
	}

	tokenMinutesDuration = time.Duration(duration)

	return &JwtConfig{
		secret:         secret,
		claimsKey:      [2]string{"uuid", "time"},
		signMethod:     jwt.SigningMethodHS256,
		expirationTime: time.Minute * tokenMinutesDuration,
	}
}
