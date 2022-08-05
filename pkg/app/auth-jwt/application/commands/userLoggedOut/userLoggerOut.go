package userLoggedOut

import (
	jwtDomain "github.com/rcrespodev/user_manager/pkg/app/auth-jwt/domain"
	returnLog "github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain/message"
)

type UserLoggerOut struct {
	jwtRepository jwtDomain.JwtRepository
}

func NewUserLoggerOut(jwtRepository jwtDomain.JwtRepository) *UserLoggerOut {
	return &UserLoggerOut{jwtRepository: jwtRepository}
}

func (u *UserLoggerOut) Exec(command *Command, log *returnLog.ReturnLog, done chan bool) {
	defer func() {
		done <- true
	}()

	targetJwt := u.jwtRepository.FindByUuid(jwtDomain.FindByUuidQuery{Uuid: command.UserUuid()})
	switch targetJwt {
	case nil:
		log.LogError(returnLog.NewErrorCommand{
			NewMessageCommand: &message.NewMessageCommand{
				ObjectId:   command.UserUuid(),
				MessageId:  0,
				MessagePkg: message.AuthorizationPkg,
			},
		})
	default:
		if !targetJwt.IsValid {
			log.LogError(returnLog.NewErrorCommand{
				NewMessageCommand: &message.NewMessageCommand{
					ObjectId:   command.UserUuid(),
					MessageId:  1,
					MessagePkg: message.AuthorizationPkg,
				},
			})
			return
		}
		u.jwtRepository.Update(jwtDomain.UpdateCommand{Command: &jwtDomain.JwtSchema{
			Uuid:    command.UserUuid(),
			IsValid: false,
			Token:   targetJwt.Token,
		}}, log)
		if log.Error() != nil {
			return
		}
		log.LogSuccess(&message.NewMessageCommand{
			ObjectId:   command.UserUuid(),
			MessageId:  16,
			MessagePkg: "user",
		})
	}
}
