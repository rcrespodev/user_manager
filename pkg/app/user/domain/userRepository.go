package domain

import (
	"github.com/google/uuid"
	returnLog "github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain"
	"sync"
)

type UserRepository interface {
	SaveUser(user *User, log *returnLog.ReturnLog, wg *sync.WaitGroup)
	FindUserById(command FindByIdCommand) *User
	FindUserByEmail(command FindByEmailCommand) *User
	FindUserByAlias(command FindByAliasCommand) *User
}

type FindUserCommand struct {
	Password string
	Log      *returnLog.ReturnLog
	Wg       *sync.WaitGroup
}

type FindByIdCommand struct {
	Uuid            uuid.UUID
	FindUserCommand FindUserCommand
}

type FindByEmailCommand struct {
	Email           *UserEmail
	FindUserCommand FindUserCommand
}

type FindByAliasCommand struct {
	Alias           *UserAlias
	FindUserCommand FindUserCommand
}
