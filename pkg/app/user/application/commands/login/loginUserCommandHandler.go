package login

import (
	"fmt"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/command"
	returnLog "github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain"
)

type LoginUserCommandHandler struct {
	userLogger *UserLogger
	cmd        *LoginUserCommand
}

func NewLoginUserCommandHandler(userLogger *UserLogger) *LoginUserCommandHandler {
	return &LoginUserCommandHandler{userLogger: userLogger}
}

func (l *LoginUserCommandHandler) Handle(command command.Command, log *returnLog.ReturnLog, done chan bool) {
	cmd, ok := command.Args().(*LoginUserCommand)
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
