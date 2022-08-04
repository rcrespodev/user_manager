package userLogged

import (
	"fmt"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/command"
	returnLog "github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain"
)

type CommandHandler struct {
	userLogger *UserLogger
	command    *Command
}

func NewCommandHandler(userLogger *UserLogger) *CommandHandler {
	return &CommandHandler{userLogger: userLogger}
}

func (c CommandHandler) Handle(command command.Command, log *returnLog.ReturnLog, done chan bool) {
	cmd, ok := command.Args().(*Command)
	if !ok {
		log.LogError(returnLog.NewErrorCommand{
			Error: fmt.Errorf("invalid type assertion"),
		})
		done <- true
		return
	}
	c.command = cmd
	go c.userLogger.Exec(c.command, log, done)
	<-done
}
