package domain

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"time"
)

// SignJwt build a jwt and return the string token.
// The claims are composed by: uuid.UUID and time.Time.
// The proposal is pass the user uuid to sign the token.
func SignJwt(uuid uuid.UUID, loggedOn time.Time, config *JwtConfig) (string, error) {
	key0 := config.claimsKey[0]
	key1 := config.claimsKey[1]
	token := jwt.NewWithClaims(config.signMethod, jwt.MapClaims{
		key0:  uuid.String(),
		key1:  loggedOn.Unix(),
		"exp": time.Now().Add(config.expirationTime).Unix(),
	})
	tokenString, err := token.SignedString([]byte(config.secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
