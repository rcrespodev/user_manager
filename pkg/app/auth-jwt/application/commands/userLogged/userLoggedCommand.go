package userLogged

import (
	"github.com/google/uuid"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/command"
)

type Command struct {
	userUuid    uuid.UUID
	baseCommand *command.BaseCommand
}

func NewCommand(userUuid uuid.UUID) *Command {
	return &Command{
		userUuid:    userUuid,
		baseCommand: command.NewBaseCommand(userUuid.String(), command.UserLogged),
	}
}

func (c Command) UserUuid() uuid.UUID {
	return c.userUuid
}

func (c Command) BaseCommand() *command.BaseCommand {
	return c.baseCommand
}
