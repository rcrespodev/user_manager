package handlers_test

import (
	"github.com/joho/godotenv"
	"github.com/rcrespodev/user_manager/pkg/kernel"
	"github.com/rcrespodev/user_manager/test/integration"
	handlers "github.com/rcrespodev/user_manager/test/integration/v1/handlers/newTest"
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
		log.Printf("container %s removed", mySqlPool.DockerResource.Container.Name)

		if err := redisPool.DockerPool.Purge(redisPool.DockerResource); err != nil {
			os.Exit(3)
		}
		os.Exit(code)
	}()
}

func TestCheckHealthStatus(t *testing.T) {
	handlers.TestCheckStatusGinHandlerFunc(t)
}

func TestGetUser(t *testing.T) {
	handlers.TestGetUserGinHandlerFunc(t)
}

func TestPepito(t *testing.T) {
	handlers.NonTest(t)
}
