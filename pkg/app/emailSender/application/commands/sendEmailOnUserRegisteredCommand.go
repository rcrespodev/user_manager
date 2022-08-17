package commands

import "github.com/rcrespodev/user_manager/pkg/kernel/cqrs/command"

type SendEmailOnUserRegisteredCommand struct {
	baseCommand *command.BaseCommand
}

func NewSendEmailOnUserRegisteredCommand(aggregateId string) *SendEmailOnUserRegisteredCommand {
	return &SendEmailOnUserRegisteredCommand{
		baseCommand: command.NewBaseCommand(aggregateId, command.SendEmailUserRegistered),
	}
}

func (s SendEmailOnUserRegisteredCommand) BaseCommand() *command.BaseCommand {
	return s.baseCommand
}
