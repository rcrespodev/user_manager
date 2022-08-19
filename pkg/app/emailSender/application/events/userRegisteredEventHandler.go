package events

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/rcrespodev/user_manager/pkg/app/emailSender/application/commands"
	userEvents "github.com/rcrespodev/user_manager/pkg/app/user/application/commands/registerUser"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/command"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/event"
	returnLog "github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain/message"
)

type UserRegisteredEventHandler struct {
	messageRepository         message.MessageRepository
	sendEmailOnUserRegistered *commands.SendEmailOnUserRegistered
	commandBus                command.Bus
}

func NewUserRegisteredEventHandler(sendEmailOnUserRegistered *commands.SendEmailOnUserRegistered,
	commandBus command.Bus, messageRepository message.MessageRepository) *UserRegisteredEventHandler {
	return &UserRegisteredEventHandler{
		sendEmailOnUserRegistered: sendEmailOnUserRegistered,
		commandBus:                commandBus,
		messageRepository:         messageRepository,
	}
}

func (u UserRegisteredEventHandler) Handle(events event.Event) {
	log := returnLog.NewReturnLog(uuid.New(), u.messageRepository, "user")
	userRegisteredEvent, ok := events.(*userEvents.UserRegistered)
	if !ok {
		log.LogError(returnLog.NewErrorCommand{
			Error: fmt.Errorf("invalid type assertion"),
		})
		return
	}
	aggregateId := userRegisteredEvent.BaseEvent().AggregateId().String()
	sendEmailCommand := commands.NewSendEmailOnUserRegisteredCommand(aggregateId)
	u.commandBus.Exec(sendEmailCommand, log)
	return
}
