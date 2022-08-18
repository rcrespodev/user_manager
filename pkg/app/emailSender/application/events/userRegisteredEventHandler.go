package events

import (
	"fmt"
	"github.com/rcrespodev/user_manager/pkg/app/emailSender/application/commands"
	userEvents "github.com/rcrespodev/user_manager/pkg/app/user/application/commands/register"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/command"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/event"
	returnLog "github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain"
)

type UserRegisteredEventHandler struct {
	sendEmailOnUserRegistered *commands.SendEmailOnUserRegistered
	commandBus                command.BusInterface
}

func (u UserRegisteredEventHandler) Handle(events []event.Event, log *returnLog.ReturnLog, done chan bool) {
	for _, e := range events {
		userRegisteredEvent, ok := e.(*userEvents.UserRegistered)
		if !ok {
			log.LogError(returnLog.NewErrorCommand{
				Error: fmt.Errorf("invalid type assertion"),
			})
			done <- true
			return
		}
		aggregateId := userRegisteredEvent.BaseEvent().AggregateId().String()
		sendEmailCommand := commands.NewSendEmailOnUserRegisteredCommand(aggregateId)
		u.commandBus.Exec(sendEmailCommand, log)
	}
}
