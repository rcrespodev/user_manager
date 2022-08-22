package userLoggedOut

import (
	"fmt"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/command"
	returnLog "github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain"
)

type CommandHandler struct {
	userLoggerOut *UserLoggerOut
	command       *Command
}

func NewCommandHandler(userLoggerOut *UserLoggerOut) *CommandHandler {
	return &CommandHandler{userLoggerOut: userLoggerOut}
}

func (c CommandHandler) Handle(command command.Command, log *returnLog.ReturnLog, done chan bool) {
	cmd, ok := command.(*Command)
	if !ok {
		log.LogError(returnLog.NewErrorCommand{
			Error: fmt.Errorf("invalid type assertion"),
		})
		done <- true
		return
	}
	c.command = cmd
	go c.userLoggerOut.Exec(c.command, log, done)
	<-done
}
