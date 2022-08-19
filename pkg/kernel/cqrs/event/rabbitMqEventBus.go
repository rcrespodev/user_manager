package event

import (
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/amqp/rabbitMq"
	"log"
)

type RabbitMqEventBus struct {
	handlersMap     HandlersMap
	client          *rabbitMq.Client
	publishedEvents publishedEvents
}

type HandlersMap map[Id][]Handler
type publishedEvents map[string]Event

type Schema struct {
	EventUuid   string `json:"event_uuid"`
	CommandUuid string `json:"command_uuid"`
	AggregateId string `json:"aggregate_id"`
	EventId     string `json:"event_id"`
	OccurredOn  string `json:"occurred_on"`
}

func NewRabbitMqEventBus(client *rabbitMq.Client) *RabbitMqEventBus {
	return &RabbitMqEventBus{
		handlersMap:     map[Id][]Handler{},
		client:          client,
		publishedEvents: map[string]Event{},
	}
}

func (r *RabbitMqEventBus) Subscribe(eventId Id, handler Handler) {
	handlers, ok := r.handlersMap[eventId]
	if !ok {
		r.handlersMap[eventId] = append(handlers, handler)
		err := r.client.DeclareQueue(string(eventId))
		if err != nil {
			log.Fatal(err)
		}
	}

	messages, err := r.client.ConsumeQueue(string(eventId))
	if err != nil {
		log.Fatal(err)
	}

	go r.messagesHandler(messages, eventId)
	return
}

func (r *RabbitMqEventBus) Publish(events []Event) {
	for _, event := range events {
		eventId := event.BaseEvent().eventId
		_, ok := r.handlersMap[eventId]
		if !ok {
			err := r.client.DeclareQueue(string(eventId))
			if err != nil {
				log.Fatal(err)
			}
		}
		schema := &Schema{
			EventUuid:   event.BaseEvent().eventUuid.String(),
			CommandUuid: event.BaseEvent().commandUuid.String(),
			AggregateId: event.BaseEvent().aggregateId.String(),
			EventId:     string(event.BaseEvent().eventId),
			OccurredOn:  event.BaseEvent().occurredOn.String(),
		}
		bodyMessage, err := json.Marshal(schema)
		if err != nil {
			log.Fatal(err)
		}
		err = r.client.PublishMessage(schema.EventId, schema.EventUuid, bodyMessage)
		if err != nil {
			log.Fatal(err)
		}
		r.publishedEvents[schema.EventUuid] = event
	}
	return
}

func (r *RabbitMqEventBus) messagesHandler(messages <-chan amqp.Delivery, eventId Id) {
	handlers, ok := r.handlersMap[eventId]
	if !ok {
		return
	}

	for message := range messages {
		var schema Schema
		err := json.Unmarshal(message.Body, &schema)
		if err != nil {
			log.Fatal(err)
		}
		event, ok := r.publishedEvents[message.MessageId]
		if !ok {
			continue
		}
		for _, handler := range handlers {
			handler.Handle(event)
		}
		err = message.Ack(false)
	}
}
