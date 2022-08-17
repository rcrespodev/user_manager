package event

import (
	"github.com/google/uuid"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/command"
	"time"
)

type Event interface {
	BaseEvent() *BaseEvent
}

type BaseEvent struct {
	eventUuid   uuid.UUID
	commandUuid uuid.UUID
	aggregateId command.AggregateId // key of aggregate created/updated/deleted.
	eventId     Id                  // examples of eventId: register_new_user, update_user, etc.
	occurredOn  time.Time
}

type Id uint8

const (
	UserRegistered Id = iota + 1
)

func NewBaseEvent(aggregateId command.AggregateId, eventId Id, commandUuid uuid.UUID) *BaseEvent {
	return &BaseEvent{
		eventUuid:   uuid.New(),
		commandUuid: commandUuid,
		aggregateId: aggregateId,
		eventId:     eventId,
		occurredOn:  time.Now(),
	}
}

func (b BaseEvent) EventUuid() uuid.UUID {
	return b.eventUuid
}

func (b BaseEvent) CommandUuid() uuid.UUID {
	return b.commandUuid
}

func (b BaseEvent) AggregateId() command.AggregateId {
	return b.aggregateId
}

func (b BaseEvent) EventId() Id {
	return b.eventId
}

func (b BaseEvent) OccurredOn() time.Time {
	return b.occurredOn
}
