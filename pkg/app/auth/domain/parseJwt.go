package domain

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
)

func IsValidJwt(tokenString string, config *JwtConfig) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(config.secret), nil
	})
	if err != nil {
		return fmt.Errorf("can´t parse string. The token is invalid")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return err
	}

	for _, key := range config.claimsKey {
		if _, ok := claims[key]; !ok {
			return err
		}
	}
	return nil
}
