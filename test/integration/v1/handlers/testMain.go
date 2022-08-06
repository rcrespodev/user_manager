package handlers

import (
	"github.com/joho/godotenv"
	"github.com/rcrespodev/user_manager/pkg/kernel"
	"github.com/rcrespodev/user_manager/test/integration"
	"log"
	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	if err := godotenv.Load("./../.env"); err != nil {
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
		time.Sleep(time.Second * 5)
		os.Exit(code)
	}()
}
