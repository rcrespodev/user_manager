package handlers_test

import (
	"github.com/joho/godotenv"
	userRepository2 "github.com/rcrespodev/user_manager/pkg/app/user/repository/userRepository"
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
	mySqlRepository, ok := kernel.Instance.UserRepository().(*userRepository2.MySqlUserRepository)
	if ok {
		_, err := mySqlRepository.ClearAll()
		require.NoError(t, err)
	}
	handlers.TestRegisterUserGinHandlerFunc(t)
}

func TestLoginUser(t *testing.T) {
	mySqlRepository, ok := kernel.Instance.UserRepository().(*userRepository2.MySqlUserRepository)
	if ok {
		_, err := mySqlRepository.ClearAll()
		require.NoError(t, err)
	}
	handlers.TestLoginUserGinHandlerFunc(t)
}

func TestGetUser(t *testing.T) {
	mySqlRepository, ok := kernel.Instance.UserRepository().(*userRepository2.MySqlUserRepository)
	if ok {
		_, err := mySqlRepository.ClearAll()
		require.NoError(t, err)
	}
	handlers.TestGetUserGinHandlerFunc(t)
}

func TestLogOutUser(t *testing.T) {
	mySqlRepository, ok := kernel.Instance.UserRepository().(*userRepository2.MySqlUserRepository)
	if ok {
		_, err := mySqlRepository.ClearAll()
		require.NoError(t, err)
	}
	handlers.TestLogOutUserGinHandlerFunc(t)
}
