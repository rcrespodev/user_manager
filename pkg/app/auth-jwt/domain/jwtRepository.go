package domain

import (
	returnLog "github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain"
)

type JwtRepository interface {
	Update(command UpdateCommand, log *returnLog.ReturnLog)
	FindByUuid(query FindByUuidQuery) *JwtSchema
}

type JwtSchema struct {
	Uuid    string
	IsValid bool
	Token   string
}

type UpdateCommand struct {
	Command *JwtSchema
}

type FindByUuidQuery struct {
	Uuid string
}
