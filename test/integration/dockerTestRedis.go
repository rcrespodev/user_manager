package integration

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"log"
)

func NewDockerTestRedis() *redis.Client {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	redisOptions := dockertest.RunOptions{
		Name:         "test_app_redis",
		Repository:   "redis",
		Tag:          "7.0.0-alpine",
		ExposedPorts: []string{"6379"},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"6379": {
				{HostIP: "0.0.0.0", HostPort: "6379"},
			},
		},
	}
	redisResource, err := pool.RunWithOptions(&redisOptions, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{
			Name: "no",
		}
	})
	if err != nil {
		log.Fatalf("Could not start redisResource: %s", err.Error())
	}

	var redisClient *redis.Client
	if err = pool.Retry(func() error {
		redisClient = redis.NewClient(&redis.Options{
			Addr: fmt.Sprintf("localhost:%s", redisResource.GetPort("6379/tcp")),
		})

		return redisClient.Ping(context.Background()).Err()
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	return redisClient
}
