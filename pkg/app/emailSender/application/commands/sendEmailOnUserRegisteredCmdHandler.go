package commands

import (
	"fmt"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/command"
	returnLog "github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain"
)

type SendEmailOnUserRegisteredCmdHandler struct {
	sendEmailOnUserRegistered *SendEmailOnUserRegistered
}

func NewSendEmailOnUserRegisteredCmdHandler(sendEmailOnUserRegistered *SendEmailOnUserRegistered) *SendEmailOnUserRegisteredCmdHandler {
	return &SendEmailOnUserRegisteredCmdHandler{sendEmailOnUserRegistered: sendEmailOnUserRegistered}
}

func (s SendEmailOnUserRegisteredCmdHandler) Handle(cmd command.Command, log *returnLog.ReturnLog, done chan bool) {
	sendEmailOnUserRegisteredCommand, ok := cmd.(*SendEmailOnUserRegisteredCommand)
	if !ok {
		log.LogError(returnLog.NewErrorCommand{
			Error: fmt.Errorf("invalid type assertion"),
		})
		done <- true
		return
	}
	s.sendEmailOnUserRegistered.Exec(sendEmailOnUserRegisteredCommand, log)
	done <- true
}
