package domain

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

type Jwt struct {
	publicKey         []byte
	privateKey        []byte
	signMethod        jwt.SigningMethod
	expirationMinutes time.Duration
}

func NewJwt(publicKey, privateKey []byte, expirationMinutes time.Duration) *Jwt {
	return &Jwt{
		publicKey:         publicKey,
		privateKey:        privateKey,
		signMethod:        jwt.SigningMethodHS256,
		expirationMinutes: time.Minute * expirationMinutes,
	}
}

func (j Jwt) CreateNewToken(key string) (string, error) {
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(j.privateKey)
	if err != nil {
		return "", err
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"key": key,
		"exp": time.Now().Add(j.expirationMinutes).Unix(),
		"iat": time.Now().Unix(),
	}).SignedString(privateKey)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (j Jwt) ValidateToken(tokenString string) (interface{}, error) {
	key, err := jwt.ParseRSAPublicKeyFromPEM(j.publicKey)
	if err != nil {
		return nil, err
	}

	tok, err := jwt.Parse(tokenString, func(jwtToken *jwt.Token) (interface{}, error) {
		if _, ok := jwtToken.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected method: %s", jwtToken.Header["alg"])
		}

		return key, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := tok.Claims.(jwt.MapClaims)
	if !ok || !tok.Valid {
		return nil, fmt.Errorf("validate: invalid claims type")
	}

	return claims, nil
}
