package delete

import (
	"fmt"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/command"
	returnLog "github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain"
)

type DeleteUserCommandHandler struct {
	userDeleter *UserDeleter
	cmd         *DeleteUserCommand
}

func NewDeleteUserCommandHandler(userDeleter *UserDeleter) *DeleteUserCommandHandler {
	return &DeleteUserCommandHandler{userDeleter: userDeleter}
}

func (d DeleteUserCommandHandler) Handle(command command.Command, log *returnLog.ReturnLog, done chan bool) {
	cmd, ok := command.Args().(*DeleteUserCommand)
	if !ok {
		log.LogError(returnLog.NewErrorCommand{
			Error: fmt.Errorf("invalid type assertion"),
		})
		done <- true
		return
	}
	d.cmd = cmd
	d.userDeleter.Exec(d.cmd, log)
	done <- true
}
