package domain

import (
	returnLog "github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain"
	"time"
)

type UserSessionRepository interface {
	UpdateUserSession(command UpdateUserSessionCommand, log *returnLog.ReturnLog)
	GetUserSession(userUuid string) *UserSessionSchema
}

type UpdateUserSessionCommand UserSessionSchema

type UserSessionSchema struct {
	UserUuid     string
	IsLogged     bool
	LastLoginOn  time.Time
	LastLogoutOn time.Time
}
