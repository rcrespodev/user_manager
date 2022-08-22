package domain

import (
	returnLog "github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain"
)

type UserRepository interface {
	SaveUser(user *User, log *returnLog.ReturnLog)
	FindUser(Query FindUserQuery) *UserSchema
	DeleteUser(user *User, log *returnLog.ReturnLog)
}

type FindUserQuery struct {
	Log   *returnLog.ReturnLog
	Where []WhereArgs
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
	HashedPassword []byte
}
