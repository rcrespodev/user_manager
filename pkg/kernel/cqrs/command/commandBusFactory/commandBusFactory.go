package commandBusFactory

import (
	"github.com/rcrespodev/user_manager/pkg/app/auth-jwt/application/commands/userLogged"
	"github.com/rcrespodev/user_manager/pkg/app/auth-jwt/application/commands/userLoggedOut"
	jwtDomain "github.com/rcrespodev/user_manager/pkg/app/auth-jwt/domain"
	"github.com/rcrespodev/user_manager/pkg/app/user/application/commands/login"
	"github.com/rcrespodev/user_manager/pkg/app/user/application/commands/register"
	"github.com/rcrespodev/user_manager/pkg/app/user/domain"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/command"
)

type NewCommandBusCommand struct {
	UserRepository domain.UserRepository
	Jwt            *jwtDomain.Jwt
	JwtRepository  jwtDomain.JwtRepository
}

func NewCommandBusInstance(busCommand NewCommandBusCommand) *command.Bus {
	registerUserCommandHandler := register.NewRegisterUserCommandHandler(
		register.NewUserRegistration(busCommand.UserRepository))

	loginUserCommandHandler := login.NewLoginUserCommandHandler(
		login.NewUserLogger(busCommand.UserRepository))

	userLoggedCommandHandler := userLogged.NewCommandHandler(
		userLogged.NewUserLogger(busCommand.Jwt, busCommand.JwtRepository))

	userLoggedOutCommandHandler := userLoggedOut.NewCommandHandler(
		userLoggedOut.NewUserLoggerOut(busCommand.JwtRepository))

	return command.NewBus(command.HandlersMap{
		command.RegisterUser:  registerUserCommandHandler,
		command.LoginUser:     loginUserCommandHandler,
		command.UserLogged:    userLoggedCommandHandler,
		command.UserLoggedOut: userLoggedOutCommandHandler,
	})
}
