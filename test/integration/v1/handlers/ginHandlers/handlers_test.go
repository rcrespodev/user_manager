package handlers_test

import (
	"github.com/joho/godotenv"
	jwtRepository "github.com/rcrespodev/user_manager/pkg/app/auth-jwt/repository"
	userRepository "github.com/rcrespodev/user_manager/pkg/app/user/repository/userRepository"
	"github.com/rcrespodev/user_manager/pkg/kernel"
	"github.com/rcrespodev/user_manager/test/integration"
	handlers "github.com/rcrespodev/user_manager/test/integration/v1/handlers/ginHandlers"
	"github.com/stretchr/testify/require"
	"log"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal(err)
	}

	mySqlPool := integration.NewDockerTestMySql()
	redisPool := integration.NewDockerTestRedis()

	kernel.NewPrdKernel(mySqlPool.MySqlClient, redisPool.RedisClient)

	code := m.Run()

	defer func() {
		if err := mySqlPool.DockerPool.Purge(mySqlPool.DockerResource); err != nil {
			os.Exit(3)
		}
		log.Printf(" === Docker Container %s removed", mySqlPool.DockerResource.Container.Name)

		if err := redisPool.DockerPool.Purge(redisPool.DockerResource); err != nil {
			os.Exit(3)
		}
		log.Printf(" === Docker Container %s removed", redisPool.DockerResource.Container.Name)
		os.Exit(code)
	}()
}

func TestCheckHealthStatus(t *testing.T) {
	handlers.TestCheckStatusGinHandlerFunc(t)
}

func TestRegisterUser(t *testing.T) {
	clearRepositories(t)
	handlers.TestRegisterUserGinHandlerFunc(t)
}

func TestLoginUser(t *testing.T) {
	clearRepositories(t)
	handlers.TestLoginUserGinHandlerFunc(t)
}

func TestLogOutUser(t *testing.T) {
	clearRepositories(t)
	handlers.TestLogOutUserGinHandlerFunc(t)
}

func TestGetUser(t *testing.T) {
	clearRepositories(t)
	handlers.TestGetUserGinHandlerFunc(t)
}

func clearRepositories(t *testing.T) {
	mySqlUserRepository, ok := kernel.Instance.UserRepository().(*userRepository.MySqlUserRepository)
	if ok {
		_, err := mySqlUserRepository.ClearAll()
		require.NoError(t, err)
	}

	jwtRedisRepository, ok := kernel.Instance.JwtRepository().(*jwtRepository.RedisJwtRepository)
	if ok {
		jwtRedisRepository.ClearAll()
	}

}
