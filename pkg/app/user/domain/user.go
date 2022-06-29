package domain

import (
	"github.com/google/uuid"
	returnLog "github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain/message"
	"reflect"
	"strings"
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
	log.SetObjectId(cmd.Alias)

	checkMandatory(cmd, log)
	if log.Error() != nil {
		return nil
	}

	userUuid, err := uuid.Parse(cmd.Uuid)
	if err != nil {
		log.LogError(returnLog.NewErrorCommand{
			Error: nil,
			NewMessageCommand: &message.NewMessageCommand{
				MessageId: 004,
				Variables: message.Variables{cmd.Uuid, "uuid"},
			},
		})
		return nil
	}

	userAlias := NewUserAlias(strings.TrimSpace(cmd.Alias), log)
	if log.Error() != nil {
		return nil
	}

	userName := NewUserName(strings.TrimSpace(cmd.Name), log)
	if log.Error() != nil {
		return nil
	}

	userSecondName := NewUserName(strings.TrimSpace(cmd.SecondName), log)
	if log.Error() != nil {
		return nil
	}

	userEmail := NewUserEmail(strings.TrimSpace(cmd.Email), log)
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

func checkMandatory(command NewUserCommand, log *returnLog.ReturnLog) {
	mandatory := []string{"Uuid", "Alias", "Name", "SecondName", "Email", "Password"}
	valueOf := reflect.ValueOf(command)
	for i := 0; i < valueOf.NumField(); i++ {
		fieldName := valueOf.Type().Field(i).Name
		fieldValue := valueOf.Field(i).Interface()
		if fieldName == mandatory[i] && fieldValue == "" {
			log.LogError(returnLog.NewErrorCommand{
				Error: nil,
				NewMessageCommand: &message.NewMessageCommand{
					MessageId: 005,
					Variables: message.Variables{fieldName},
				},
			})
			return
		}
	}
}
