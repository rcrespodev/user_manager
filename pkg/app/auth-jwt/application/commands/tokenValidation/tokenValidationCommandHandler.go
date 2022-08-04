package tokenValidation

import (
	"fmt"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/command"
	returnLog "github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain"
)

type CommandHandler struct {
	command        Command
	tokenValidator *TokenValidator
}

func NewCommandHandler(tokenValidator *TokenValidator) *CommandHandler {
	return &CommandHandler{tokenValidator: tokenValidator}
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
	c.command = *cmd
	c.tokenValidator.Exec(c.command, log)
	done <- true
}
