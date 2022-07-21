package domain

import (
	returnLog "github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain"
	"golang.org/x/crypto/bcrypt"
)

func LoginUser(password *UserPassword, schema *UserSchema, log *returnLog.ReturnLog) *User {
	if err := bcrypt.CompareHashAndPassword(schema.HashedPassword,
		[]byte(password.stringPassword)); err != nil {
		return nil
	}
	return NewUser(NewUserCommand{
		Uuid:       schema.Uuid,
		Alias:      schema.Alias,
		Name:       schema.Name,
		SecondName: schema.SecondName,
		Email:      schema.Email,
		Password:   password.stringPassword,
		IgnorePass: false,
	}, log)
}
