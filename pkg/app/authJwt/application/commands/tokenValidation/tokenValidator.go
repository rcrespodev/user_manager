package tokenValidation

import (
	"fmt"
	jwtDomain "github.com/rcrespodev/user_manager/pkg/app/authJwt/domain"
	returnLog "github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain/message"
)

type TokenValidator struct {
	jwt           *jwtDomain.Jwt
	jwtRepository jwtDomain.JwtRepository
	log           *returnLog.ReturnLog
	token         string
}

func NewTokenValidator(jwt *jwtDomain.Jwt, jwtRepository jwtDomain.JwtRepository) *TokenValidator {
	return &TokenValidator{
		jwt:           jwt,
		jwtRepository: jwtRepository,
	}
}

func (t *TokenValidator) Exec(command Command, log *returnLog.ReturnLog) {
	t.token = command.token
	t.log = log

	token, err := t.jwt.ValidateToken(t.token)
	if err != nil {
		t.unauthorizedResponse()
		return
	}

	tokenUuid := token["key"].(string)
	targetToken := t.jwtRepository.FindByUuid(jwtDomain.FindByUuidQuery{Uuid: tokenUuid})
	switch targetToken {
	case nil:
		t.unauthorizedResponse()
		return
	default:
		if !targetToken.IsValid {
			t.unauthorizedResponse()
			return
		}
	}

	return
}

func (t *TokenValidator) unauthorizedResponse() {
	t.log.LogError(returnLog.NewErrorCommand{
		NewMessageCommand: &message.NewMessageCommand{
			ObjectId:   fmt.Sprintf("%s", t.token),
			MessageId:  0,
			MessagePkg: message.AuthorizationPkg,
		},
	})
}
