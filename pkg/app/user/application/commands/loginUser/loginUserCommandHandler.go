package loginUser

import (
	"fmt"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/command"
	returnLog "github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain"
)

type CommandHandler struct {
	userLogger *Service
	cmd        *Command
}

func NewLoginUserCommandHandler(userLogger *Service) *CommandHandler {
	return &CommandHandler{userLogger: userLogger}
}

func (l *CommandHandler) Handle(command command.Command, log *returnLog.ReturnLog, done chan bool) {
	cmd, ok := command.(*Command)
	if !ok {
		log.LogError(returnLog.NewErrorCommand{
			Error: fmt.Errorf("invalid type assertion"),
		})
		done <- true
		return
	}
	l.cmd = cmd
	l.userLogger.Exec(l.cmd, log)
	done <- true
}
