package deleteUser

import (
	"fmt"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/command"
	returnLog "github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain"
)

type CommandHandler struct {
	userDeleter *Service
	cmd         *Command
}

func NewDeleteUserCommandHandler(userDeleter *Service) *CommandHandler {
	return &CommandHandler{userDeleter: userDeleter}
}

func (d CommandHandler) Handle(command command.Command, log *returnLog.ReturnLog, done chan bool) {
	cmd, ok := command.(*Command)
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
