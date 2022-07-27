package loginUser

import (
	"github.com/rcrespodev/user_manager/pkg/kernel"
	"github.com/rcrespodev/user_manager/test/integration"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	mySqlPool := integration.NewDockerTestMySql()
	redisPool := integration.NewDockerTestRedis()

	kernel.NewPrdKernel(mySqlPool.MySqlClient, redisPool.RedisClient)

	code := m.Run()

	defer func() {
		if err := mySqlPool.DockerPool.Purge(mySqlPool.DockerResource); err != nil {
			os.Exit(3)
		}

		if err := redisPool.DockerPool.Purge(redisPool.DockerResource); err != nil {
			os.Exit(3)
		}
		os.Exit(code)
	}()
}
