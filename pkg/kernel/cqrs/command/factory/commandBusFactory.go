package factory

import (
	"github.com/rcrespodev/user_manager/pkg/app/user/application/register"
	"github.com/rcrespodev/user_manager/pkg/app/user/domain"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/command"
)

type NewCommandBusCommand struct {
	UserRepository domain.UserRepository
}

func NewCommandBusInstance(busCommand NewCommandBusCommand) *command.Bus {
	registerUserHandler := register.NewRegisterUserCommandHandler(
		register.NewUserRegistration(busCommand.UserRepository))

	return command.NewBus(command.HandlersMap{
		command.RegisterUser: registerUserHandler,
	})
}
