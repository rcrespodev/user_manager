package rabbitMq

import (
	"context"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rcrespodev/user_manager/pkg/kernel/config"
	"log"
	"time"
)

type Client struct {
	client *amqp.Connection
	chanel *amqp.Channel
	queues queues
}

type queues map[string]amqp.Queue

func NewRabbitMqClient(connection *amqp.Connection) *Client {
	if connection == nil {
		var err error
		envs := config.Conf.Rabbit
		connection, err = amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s/",
			envs.User, envs.Password, envs.Host))
		if err != nil {
			log.Fatal(err)
		}
	}
	ch, err := connection.Channel()
	if err != nil {
		log.Fatal(err)
	}
	return &Client{
		client: connection,
		chanel: ch,
		queues: map[string]amqp.Queue{},
	}
}

func (c *Client) DeclareQueue(queueName string) error {
	queue, err := c.chanel.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err == nil {
		c.queues[queueName] = queue
	}
	return err
}

func (c *Client) PublishMessage(queueName, messageId string, message []byte) error {
	if _, ok := c.queues[queueName]; !ok {
		return fmt.Errorf("plis, declare %s queue first", queueName)
	}
	err := c.chanel.PublishWithContext(
		context.Background(),
		"",
		queueName,
		false,
		false,
		amqp.Publishing{
			Headers:         nil,
			ContentType:     "application/json",
			ContentEncoding: "",
			DeliveryMode:    0,
			Priority:        0,
			CorrelationId:   "",
			ReplyTo:         "",
			Expiration:      "",
			MessageId:       messageId,
			Timestamp:       time.Now(),
			Type:            "",
			UserId:          "",
			AppId:           "",
			Body:            message,
		})
	return err
}

func (c *Client) ConsumeQueue(queueName string) (<-chan amqp.Delivery, error) {
	if _, ok := c.queues[queueName]; !ok {
		return nil, fmt.Errorf("queue %s is not declared yet", queueName)
	}
	return c.chanel.Consume(
		queueName,
		"",
		false,
		false,
		false,
		false,
		nil)
}

func (c *Client) Chanel() *amqp.Channel {
	return c.chanel
}
