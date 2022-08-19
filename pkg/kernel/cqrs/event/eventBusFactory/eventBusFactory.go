package eventBusFactory

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/amqp/rabbitMq"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/event"
	"log"
)

type NewEventBusCommand struct {
	RabbitMqConnection *amqp.Connection
}

func NewEventBusInstance(command NewEventBusCommand) event.Bus {
	rabbitClient := rabbitMq.NewRabbitMqClient(command.RabbitMqConnection)
	if err := rabbitClient.DeclareQueue(string(event.UserRegistered)); err != nil {
		log.Fatal(err)
	}
	return event.NewRabbitMqEventBus(rabbitMq.NewRabbitMqClient(command.RabbitMqConnection))
}
