package domain

import (
	"github.com/google/uuid"
	returnLog "github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain"
)

type User struct {
	uuid       uuid.UUID
	alias      *UserAlias
	name       *UserName
	secondName *UserName
	email      *UserEmail
	password   *UserPassword
}

type NewUserCommand struct {
	Uuid       string
	Alias      string
	Name       string
	SecondName string
	Email      string
	Password   string
}

func NewUser(cmd NewUserCommand, log *returnLog.ReturnLog) *User {
	userUuid, err := uuid.Parse(cmd.Uuid)
	if err != nil {
		log.LogError(returnLog.NewErrorCommand{
			Error: err,
		})
		return nil
	}

	userAlias := NewUserAlias(cmd.Alias, log)
	if log.Error() != nil {
		return nil
	}

	userName := NewUserName(cmd.Name, log)
	if log.Error() != nil {
		return nil
	}

	userSecondName := NewUserName(cmd.SecondName, log)
	if log.Error() != nil {
		return nil
	}

	userEmail := NewUserEmail(cmd.Email, log)
	if log.Error() != nil {
		return nil
	}

	userPassword := NewUserPassword(cmd.Password, log)
	if log.Error() != nil {
		return nil
	}

	return &User{
		uuid:       userUuid,
		alias:      userAlias,
		name:       userName,
		secondName: userSecondName,
		email:      userEmail,
		password:   userPassword,
	}
}

func (u User) Uuid() uuid.UUID {
	return u.uuid
}

func (u User) Alias() *UserAlias {
	return u.alias
}

func (u User) Name() *UserName {
	return u.name
}

func (u User) SecondName() *UserName {
	return u.secondName
}

func (u User) Email() *UserEmail {
	return u.email
}

func (u User) Password() *UserPassword {
	return u.password
}
