package kernel

import (
	"github.com/go-redis/redis/v8"
	jwtDomain "github.com/rcrespodev/user_manager/pkg/app/auth-jwt/domain"
	jwtRepository "github.com/rcrespodev/user_manager/pkg/app/auth-jwt/repository"
	"github.com/rcrespodev/user_manager/pkg/kernel/config"
	"log"
	"os"
	"strconv"
	"time"
)

func jwtFactory(redisClient *redis.Client) (*jwtDomain.Jwt, jwtDomain.JwtRepository) {
	certPublicKey, err := os.ReadFile(config.Conf.Jwt.Key.Public)
	if err != nil {
		log.Fatal(err)
	}

	certPrivateKey, err := os.ReadFile(config.Conf.Jwt.Key.Private)
	if err != nil {
		log.Fatal(err)
	}

	expirationTime := config.Conf.Jwt.ExpirationTime
	expirationTimeInt, err := strconv.Atoi(expirationTime)
	if err != nil {
		log.Fatal(err)
	}

	jwt := jwtDomain.NewJwt(certPublicKey, certPrivateKey, time.Duration(expirationTimeInt))

	repository := jwtRepository.NewRedisJwtRepository(redisClient)

	return jwt, repository
}
