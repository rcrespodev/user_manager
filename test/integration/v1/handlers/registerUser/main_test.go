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
	//pool, err := dockertest.NewPool("")
	//if err != nil {
	//	log.Fatalf("Could not connect to docker: %s", err)
	//}
	//
	//mySqlOptions := dockertest.RunOptions{
	//	Name:       "test_app_mysql",
	//	Repository: "mysql",
	//	Tag:        "5.7",
	//	Env: []string{
	//		//"POSTGRES_USER=" + "root",
	//		"MYSQL_ROOT_PASSWORD=" + "my_secret",
	//		"MYSQL_DATABASE=" + "user_manager",
	//	},
	//	ExposedPorts: []string{"3306"},
	//	PortBindings: map[docker.Port][]docker.PortBinding{
	//		"3306": {
	//			{HostIP: "0.0.0.0", HostPort: "3306"},
	//		},
	//	},
	//}
	//mySqlResource, err := pool.RunWithOptions(&mySqlOptions, func(config *docker.HostConfig) {
	//	config.AutoRemove = true
	//	config.RestartPolicy = docker.RestartPolicy{
	//		Name: "no",
	//	}
	//})
	//if err != nil {
	//	log.Fatalf("Could not start mySqlResource: %s", err.Error())
	//}

	//if err = os.Setenv("MYSQL_PORT", mySqlResource.GetPort("3306/tcp")); err != nil {
	//	log.Fatalf("%v", err)
	//}
	//
	//redisOptions := dockertest.RunOptions{
	//	Name:         "test_app_redis",
	//	Repository:   "redis",
	//	Tag:          "7.0.0-alpine",
	//	ExposedPorts: []string{"6379"},
	//	PortBindings: map[docker.Port][]docker.PortBinding{
	//		"6379": {
	//			{HostIP: "0.0.0.0", HostPort: "6379"},
	//		},
	//	},
	//}
	//redisResource, err := pool.RunWithOptions(&redisOptions, func(config *docker.HostConfig) {
	//	config.AutoRemove = true
	//	config.RestartPolicy = docker.RestartPolicy{
	//		Name: "no",
	//	}
	//})
	//if err != nil {
	//	log.Fatalf("Could not start redisResource: %s", err.Error())
	//}
	//
	////dsn = fmt.Sprintf(dsn, user, password, port, db)
	//var db *sql.DB
	//if err = pool.Retry(func() error {
	//	datasource := fmt.Sprintf("root:my_secret@(localhost:%s)/user_manager", mySqlResource.GetPort("3306/tcp"))
	//	_ = os.Setenv("DATASOURCE", datasource)
	//	db, err = sql.Open("mysql", datasource)
	//	if err != nil {
	//		return err
	//	}
	//	return db.Ping()
	//}); err != nil {
	//	log.Fatalf("Could not connect to mysql: %s", err.Error())
	//}
}
