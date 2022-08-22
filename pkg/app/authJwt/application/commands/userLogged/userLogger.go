package userLogged

import (
	jwtDomain "github.com/rcrespodev/user_manager/pkg/app/authJwt/domain"
	returnLog "github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain"
)

type UserLogger struct {
	jwt           *jwtDomain.Jwt
	jwtRepository jwtDomain.JwtRepository
}

func NewUserLogger(jwt *jwtDomain.Jwt, jwtRepository jwtDomain.JwtRepository) *UserLogger {
	return &UserLogger{
		jwt:           jwt,
		jwtRepository: jwtRepository,
	}
}

func (u *UserLogger) Exec(command *Command, log *returnLog.ReturnLog, done chan bool) {
	defer func() {
		done <- true
	}()

	token, err := u.jwt.CreateNewToken(command.UserUuid().String())
	if err != nil {
		log.LogError(returnLog.NewErrorCommand{
			Error: err,
		})
		return
	}

	u.jwtRepository.Update(jwtDomain.UpdateCommand{Command: &jwtDomain.JwtSchema{
		Uuid:    command.UserUuid().String(),
		IsValid: true,
		Token:   token,
	}}, log)

	return
}
