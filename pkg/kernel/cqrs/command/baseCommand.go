package command

import (
	"github.com/google/uuid"
	"time"
)

type Command interface {
	BaseCommand() *BaseCommand
}

type BaseCommand struct {
	commandUuid uuid.UUID
	aggregateId AggregateId // key of aggregate created/updated/deleted.
	commandId   Id          // examples of commandId: register_new_user, update_user, etc.
	occurredOn  time.Time
}

type Id uint8

const (
	RegisterUser Id = 0 + iota
	LoginUser
	UpdateUser
	DeleteUser
	TokenValidation
	UserLogged
	UserLoggedOut
	SendEmailUserRegistered
)

type AggregateId string

func (a AggregateId) String() string {
	return string(a)
}

func NewBaseCommand(aggregateId string, commandId Id) *BaseCommand {
	return &BaseCommand{
		commandUuid: uuid.New(),
		aggregateId: AggregateId(aggregateId),
		commandId:   commandId,
		occurredOn:  time.Now(),
	}
}

func (b BaseCommand) CommandUuid() uuid.UUID {
	return b.commandUuid
}

func (b BaseCommand) AggregateId() AggregateId {
	return b.aggregateId
}

func (b BaseCommand) CommandId() Id {
	return b.commandId
}

func (b BaseCommand) OccurredOn() time.Time {
	return b.occurredOn
}
