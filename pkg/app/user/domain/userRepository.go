package domain

import (
	"github.com/google/uuid"
	returnLog "github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain"
)

type UserRepository interface {
	SaveUser(*User) error
	FindUserById(uuid uuid.UUID, log *returnLog.ReturnLog) *User
}
