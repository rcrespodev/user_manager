package domain

import (
	returnLog "github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain/message"
	"net/mail"
)

type UserEmail struct {
	address string
}

func NewUserEmail(emailAddress string, log *returnLog.ReturnLog) *UserEmail {
	address, err := mail.ParseAddress(emailAddress)
	if err != nil {
		log.LogError(returnLog.NewErrorCommand{
			Error: nil,
			NewMessageCommand: &message.NewMessageCommand{
				MessageId: 004,
				Variables: message.Variables{emailAddress, "email"},
			},
		})
		return nil
	}
	return &UserEmail{address: address.Address}
}

func (e UserEmail) Address() string {
	return e.address
}
