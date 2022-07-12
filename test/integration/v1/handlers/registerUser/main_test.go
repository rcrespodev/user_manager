package registerUser

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/rcrespodev/user_manager/pkg/kernel"
	"github.com/rcrespodev/user_manager/test/integration"
	"log"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("%v", err)
	}

	mySqlClient := integration.NewDockerTestMySql()
	redisClient := integration.NewDockerTestRedis()

	kernel.NewPrdKernel(mySqlClient, redisClient)

	code := m.Run()

	os.Exit(code)
}
