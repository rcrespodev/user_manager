package userLoggedOut

import "github.com/rcrespodev/user_manager/pkg/kernel/cqrs/command"

type Command struct {
	userUuid    string
	baseCommand *command.BaseCommand
}

func NewCommand(userUuid string) *Command {
	return &Command{
		userUuid:    userUuid,
		baseCommand: command.NewBaseCommand(userUuid, command.UserLoggedOut),
	}
}

func (c Command) UserUuid() string {
	return c.userUuid
}

func (c Command) BaseCommand() *command.BaseCommand {
	return c.baseCommand
}
