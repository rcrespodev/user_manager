package domain

import (
	returnLog "github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain"
)

type UserRepository interface {
	SaveUser(user *User, log *returnLog.ReturnLog)
	FindUser(Command FindUserCommand) *User
	//FindUserById(command FindByIdCommand, user chan *User)
	//FindUserByEmail(command FindByEmailCommand, user chan *User)
	//FindUserByAlias(command FindByAliasCommand, user chan *User)
}

type FindUserCommand struct {
	Password string
	Log      *returnLog.ReturnLog
	Where    []WhereArgs
	//Wg       *sync.WaitGroup
}

type WhereArgs struct {
	Field string
	Value string
}

//type FindByIdCommand struct {
//	Uuid            uuid.UUID
//	FindUserCommand FindUserCommand
//}
//
//type FindByEmailCommand struct {
//	Email           *UserEmail
//	FindUserCommand FindUserCommand
//}
//
//type FindByAliasCommand struct {
//	Alias           *UserAlias
//	FindUserCommand FindUserCommand
//}
