package tokenValidation

import (
	"fmt"
	jwtDomain "github.com/rcrespodev/user_manager/pkg/app/auth-jwt/domain"
	returnLog "github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain/message"
)

type TokenValidator struct {
	jwt *jwtDomain.Jwt
}

func NewTokenValidator(jwt *jwtDomain.Jwt) *TokenValidator {
	return &TokenValidator{jwt: jwt}
}

func (t TokenValidator) Exec(command Command, log *returnLog.ReturnLog) {
	if _, err := t.jwt.ValidateToken(command.token); err != nil {
		log.LogError(returnLog.NewErrorCommand{
			NewMessageCommand: &message.NewMessageCommand{
				ObjectId:   fmt.Sprintf("%s", command.token),
				MessageId:  0,
				MessagePkg: message.AuthorizationPkg,
			},
		})
	}
	return
}
