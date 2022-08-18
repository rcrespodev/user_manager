package register

import (
	"fmt"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/command"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/event"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain/valueObjects"
)

type RegisterUserCommandHanlder struct {
	userRegistration *UserRegistration
	cmd              *RegisterUserCommand
	eventBus         event.Bus
}

func NewRegisterUserCommandHandler(eventBus event.Bus, userRegistration *UserRegistration) *RegisterUserCommandHanlder {
	return &RegisterUserCommandHanlder{
		userRegistration: userRegistration,
		eventBus:         eventBus,
	}
}

func (r RegisterUserCommandHanlder) Handle(command command.Command, log *domain.ReturnLog, done chan bool) {
	cmd, ok := command.(*RegisterUserCommand)
	if !ok {
		log.LogError(domain.NewErrorCommand{
			Error: fmt.Errorf("invalid type assertion"),
		})
		done <- true
		return
	}
	r.cmd = cmd
	r.userRegistration.Exec(*r.cmd, log)
	if log.Status() == valueObjects.Success {
		e := NewUserRegistered(*r.cmd)
		r.eventBus.Publish([]event.Event{e})
	}
	done <- true

}
