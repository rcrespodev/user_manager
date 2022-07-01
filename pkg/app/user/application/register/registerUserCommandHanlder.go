package register

import (
	"fmt"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/command"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain"
)

type RegisterUserCommandHanlder struct {
	userRegistration *UserRegistration
	cmd              *RegisterUserCommand
}

func NewRegisterUserCommandHanlder(userRegistration *UserRegistration) *RegisterUserCommandHanlder {
	return &RegisterUserCommandHanlder{userRegistration: userRegistration}
}

func (r RegisterUserCommandHanlder) Handle(command command.Command, log *domain.ReturnLog, done chan bool) {
	registerUserCmd, ok := command.Args().(*RegisterUserCommand)
	if !ok {
		log.LogError(domain.NewErrorCommand{
			Error: fmt.Errorf("invalid type assertion"),
		})
		done <- true
		return
	}
	r.userRegistration.Exec(*registerUserCmd, log)
	done <- true
}
