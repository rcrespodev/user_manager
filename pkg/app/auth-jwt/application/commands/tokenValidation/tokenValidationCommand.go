package tokenValidation

import "github.com/rcrespodev/user_manager/pkg/kernel/cqrs/command"

type Command struct {
	token       string
	baseCommand *command.BaseCommand
}

func NewCommand(token string) *Command {
	return &Command{
		token:       token,
		baseCommand: command.NewBaseCommand(token, command.TokenValidation),
	}
}

func (c Command) BaseCommand() *command.BaseCommand {
	return c.baseCommand
}
