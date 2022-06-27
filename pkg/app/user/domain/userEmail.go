package domain

import (
	returnLog "github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain"
	"net/mail"
)

type UserEmail struct {
	address string
}

func NewUserEmail(emailAddress string, log *returnLog.ReturnLog) *UserEmail {
	address, err := mail.ParseAddress(emailAddress)
	if err != nil {
		log.LogError(returnLog.NewErrorCommand{
			Error: err,
		})
		return nil
	}
	return &UserEmail{address: address.Address}
}

func (e UserEmail) Address() string {
	return e.address
}

//func (e UserEmail) MessageData() *messageData {
//	return e.messageData
//}
