package registerUser

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/rcrespodev/user_manager/pkg/kernel"
	"github.com/rcrespodev/user_manager/test/integration"
	"log"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	mySqlPool := integration.NewDockerTestMySql()
	redisPool := integration.NewDockerTestRedis()

	kernel.NewPrdKernel(mySqlPool.MySqlClient, redisPool.RedisClient)

	code := m.Run()

	if err := mySqlPool.DockerPool.Purge(mySqlPool.DockerResource); err != nil {
		log.Fatal(err)
	}

	if err := redisPool.DockerPool.Purge(redisPool.DockerResource); err != nil {
		log.Fatal(err)
	}

	os.Exit(code)
}
