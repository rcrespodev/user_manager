package delete

import "github.com/rcrespodev/user_manager/pkg/kernel/cqrs/command"

type DeleteUserCommand struct {
	userUuid    string
	baseCommand *command.BaseCommand
}

type ClientArgs struct {
	UserUuid string `json:"user_uuid"`
}

func NewDeleteUserCommand(userUuid string) *DeleteUserCommand {
	return &DeleteUserCommand{
		userUuid:    userUuid,
		baseCommand: command.NewBaseCommand(userUuid, command.DeleteUser),
	}
}

func (d DeleteUserCommand) UserUuid() string {
	return d.userUuid
}

func (d DeleteUserCommand) BaseCommand() *command.BaseCommand {
	return d.baseCommand
}
