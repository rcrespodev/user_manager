package integration

import (
	"fmt"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

type RabbitMqPoolConnection struct {
	resource   *dockertest.Resource
	connection *amqp.Connection
	pool       *dockertest.Pool
}

func NewRabbitPoolConnection() *RabbitMqPoolConnection {
	const (
		user     = "my_user"
		password = "my_password"
		host     = "0.0.0.0"
		port     = "5672"
	)

	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	redisOptions := dockertest.RunOptions{
		Name:         "test_app_rabbitmq",
		Repository:   "rabbitmq",
		Tag:          "3.10.7-management",
		ExposedPorts: []string{"5672", "15672"},
		Env:          []string{"RABBITMQ_DEFAULT_USER=my_user", "RABBITMQ_DEFAULT_PASS=my_password"},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"5672": {
				{HostIP: "0.0.0.0", HostPort: "5672"},
			},
			"15672": {
				{HostIP: "0.0.0.0", HostPort: "15672"},
			},
		},
	}
	resource, err := pool.RunWithOptions(&redisOptions, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{
			Name: "no",
		}
	})
	if err != nil {
		log.Fatalf("Could not start Rabbit MQ resource: %s", err.Error())
	}

	var rabbitMqClient *amqp.Connection
	if err = pool.Retry(func() error {
		rabbitMqClient, err = amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/",
			user, password, host, port))
		return err
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	return &RabbitMqPoolConnection{
		pool:       Pool,
		resource:   resource,
		connection: rabbitMqClient,
	}
}

func (r RabbitMqPoolConnection) Resource() *dockertest.Resource {
	return r.resource
}

func (r RabbitMqPoolConnection) Connection() *amqp.Connection {
	return r.connection
}

func (r RabbitMqPoolConnection) Pool() *dockertest.Pool {
	return r.pool
}
