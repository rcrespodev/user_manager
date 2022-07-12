package integration

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"log"
	"os"
)

func NewDockerTestMySql() *sql.DB {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	mySqlOptions := dockertest.RunOptions{
		Name:       "test_app_mysql",
		Repository: "mysql",
		Tag:        "5.7",
		Env: []string{
			"MYSQL_ROOT_PASSWORD=" + "my_secret",
			"MYSQL_DATABASE=" + "user_manager",
		},
		ExposedPorts: []string{"3306"},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"3306": {
				{HostIP: "0.0.0.0", HostPort: "3306"},
			},
		},
	}
	mySqlResource, err := pool.RunWithOptions(&mySqlOptions, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{
			Name: "no",
		}
	})
	if err != nil {
		log.Fatalf("Could not start mySqlResource: %s", err.Error())
	}

	if err = os.Setenv("MYSQL_PORT", mySqlResource.GetPort("3306/tcp")); err != nil {
		log.Fatalf("%v", err)
	}

	var mySqlClient *sql.DB
	if err = pool.Retry(func() error {
		datasource := fmt.Sprintf("root:my_secret@(localhost:%s)/user_manager", mySqlResource.GetPort("3306/tcp"))
		mySqlClient, err = sql.Open("mysql", datasource)
		if err != nil {
			return err
		}
		return mySqlClient.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to mysql: %s", err.Error())
	}
	return mySqlClient
}
