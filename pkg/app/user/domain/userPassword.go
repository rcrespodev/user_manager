package domain

import (
	returnLog "github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain"
	"golang.org/x/crypto/bcrypt"
)

type UserPassword struct {
	stringPassword string
	hashPassword   []byte
}

func NewUserPassword(password string, log *returnLog.ReturnLog) *UserPassword {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.LogError(returnLog.NewErrorCommand{
			Error:             err,
			NewMessageCommand: nil,
		})
		return nil
	}
	return &UserPassword{
		stringPassword: password,
		hashPassword:   hashPassword,
	}
}

func (u UserPassword) String() string {
	return u.stringPassword
}

func (u UserPassword) Hash() []byte {
	return u.hashPassword
}
