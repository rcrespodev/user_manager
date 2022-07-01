package domain

import (
	"github.com/google/uuid"
	returnLog "github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain"
	"sync"
)

type UserRepository interface {
	SaveUser(user *User, log *returnLog.ReturnLog, wg *sync.WaitGroup)
	FindUserById(uuid uuid.UUID, log *returnLog.ReturnLog, wg *sync.WaitGroup) *User
	FindUserByEmail(email *UserEmail, log *returnLog.ReturnLog, wg *sync.WaitGroup) *User
	FindUserByAlias(alias *UserAlias, log *returnLog.ReturnLog, wg *sync.WaitGroup) *User
}
