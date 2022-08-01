package factory

import (
	"github.com/rcrespodev/user_manager/pkg/app/user/application/commands/login"
	"github.com/rcrespodev/user_manager/pkg/app/user/application/commands/register"
	"github.com/rcrespodev/user_manager/pkg/app/user/domain"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/command"
)

type NewCommandBusCommand struct {
	UserRepository domain.UserRepository
}

func NewCommandBusInstance(busCommand NewCommandBusCommand) *command.Bus {
	registerUserCommandHandler := register.NewRegisterUserCommandHandler(
		register.NewUserRegistration(busCommand.UserRepository))

	loginUserCommandHandler := login.NewLoginUserCommandHandler(
		login.NewUserLogger(busCommand.UserRepository))

	return command.NewBus(command.HandlersMap{
		command.RegisterUser: registerUserCommandHandler,
		command.LoginUser:    loginUserCommandHandler,
	})
}
