package consumersFactory

import (
	"github.com/rcrespodev/user_manager/pkg/app/emailSender/application/commands"
	"github.com/rcrespodev/user_manager/pkg/app/emailSender/application/events"
	emailSenderDomain "github.com/rcrespodev/user_manager/pkg/app/emailSender/domain"
	"github.com/rcrespodev/user_manager/pkg/app/user/domain"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/command"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/event"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain/message"
)

type SubscribeConsumersCommand struct {
	EventBus            event.Bus
	CommandBus          command.BusInterface
	UserRepository      domain.UserRepository
	EmailSender         emailSenderDomain.EmailSender
	SentEmailRepository emailSenderDomain.SentEmailRepository
	MessageRepository   message.MessageRepository
	WelcomeTemplatePath string
}

func SubscribeConsumers(command SubscribeConsumersCommand) {
	sendEmailOnUserRegistered := commands.NewSendEmailOnUserRegistered(commands.SendEmailOnUserRegisteredDependencies{
		UserRepository:      command.UserRepository,
		EmailSender:         command.EmailSender,
		SentEmailRepository: command.SentEmailRepository,
		WelcomeTemplatePath: command.WelcomeTemplatePath,
	})
	userRegisteredEventHandler := events.NewUserRegisteredEventHandler(
		sendEmailOnUserRegistered, command.CommandBus, command.MessageRepository)
	go command.EventBus.Subscribe(event.UserRegistered, userRegisteredEventHandler)
}
