package commandBusFactory

import (
	"github.com/rcrespodev/user_manager/pkg/app/authJwt/application/commands/userLogged"
	"github.com/rcrespodev/user_manager/pkg/app/authJwt/application/commands/userLoggedOut"
	jwtDomain "github.com/rcrespodev/user_manager/pkg/app/authJwt/domain"
	"github.com/rcrespodev/user_manager/pkg/app/emailSender/application/commands"
	emailSenderDomain "github.com/rcrespodev/user_manager/pkg/app/emailSender/domain"
	delete "github.com/rcrespodev/user_manager/pkg/app/user/application/commands/delete"
	"github.com/rcrespodev/user_manager/pkg/app/user/application/commands/login"
	"github.com/rcrespodev/user_manager/pkg/app/user/application/commands/register"
	"github.com/rcrespodev/user_manager/pkg/app/user/domain"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/command"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/event"
)

type NewCommandBusCommand struct {
	User struct {
		UserRepository domain.UserRepository
	}
	Jwt struct {
		Jwt           *jwtDomain.Jwt
		JwtRepository jwtDomain.JwtRepository
	}
	EmailSender struct {
		EmailSender         emailSenderDomain.EmailSender
		SentEmailRepository emailSenderDomain.SentEmailRepository
		WelcomeTemplatePath string
	}
	EventBus event.Bus
}

func NewCommandBusInstance(cmd NewCommandBusCommand) *command.Bus {
	registerUserCommandHandler := register.NewRegisterUserCommandHandler(
		cmd.EventBus, register.NewUserRegistration(cmd.User.UserRepository))

	loginUserCommandHandler := login.NewLoginUserCommandHandler(
		login.NewUserLogger(cmd.User.UserRepository))

	userLoggedCommandHandler := userLogged.NewCommandHandler(
		userLogged.NewUserLogger(cmd.Jwt.Jwt, cmd.Jwt.JwtRepository))

	userLoggedOutCommandHandler := userLoggedOut.NewCommandHandler(
		userLoggedOut.NewUserLoggerOut(cmd.Jwt.JwtRepository))

	deleteUserCommandHandler := delete.NewDeleteUserCommandHandler(
		delete.NewUserDeleter(cmd.User.UserRepository))

	sendEmailOnUserRegisteredCmdHandler := commands.NewSendEmailOnUserRegisteredCmdHandler(
		commands.NewSendEmailOnUserRegistered(commands.SendEmailOnUserRegisteredDependencies{
			UserRepository:      cmd.User.UserRepository,
			EmailSender:         cmd.EmailSender.EmailSender,
			SentEmailRepository: cmd.EmailSender.SentEmailRepository,
			WelcomeTemplatePath: cmd.EmailSender.WelcomeTemplatePath,
		}))

	return command.NewBus(command.HandlersMap{
		command.RegisterUser:            registerUserCommandHandler,
		command.LoginUser:               loginUserCommandHandler,
		command.UserLogged:              userLoggedCommandHandler,
		command.UserLoggedOut:           userLoggedOutCommandHandler,
		command.DeleteUser:              deleteUserCommandHandler,
		command.SendEmailUserRegistered: sendEmailOnUserRegisteredCmdHandler,
	})
}
