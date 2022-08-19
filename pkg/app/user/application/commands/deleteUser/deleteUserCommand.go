package deleteUser

import "github.com/rcrespodev/user_manager/pkg/kernel/cqrs/command"

type Command struct {
	userUuid    string
	baseCommand *command.BaseCommand
}

type ClientArgs struct {
	UserUuid string `json:"user_uuid"`
}

func NewDeleteUserCommand(userUuid string) *Command {
	return &Command{
		userUuid:    userUuid,
		baseCommand: command.NewBaseCommand(userUuid, command.DeleteUser),
	}
}

func (d Command) UserUuid() string {
	return d.userUuid
}

func (d Command) BaseCommand() *command.BaseCommand {
	return d.baseCommand
}
