package domain

import (
	returnLog "github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain"
)

type UserRepository interface {
	SaveUser(user *User, log *returnLog.ReturnLog)
	FindUser(Command FindUserCommand) *UserSchema
}

type FindUserCommand struct {
	Password string
	Log      *returnLog.ReturnLog
	Where    []WhereArgs
}

type WhereArgs struct {
	Field string
	Value string
}

type UserSchema struct {
	Uuid           string
	Alias          string
	Name           string
	SecondName     string
	Email          string
	HashedPassword string
}
