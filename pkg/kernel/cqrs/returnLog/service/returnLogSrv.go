package service

import (
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain/customError"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain/message"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain/valueObjects"
)

type ReturnLogSrv interface {
	LogError(command domain.NewErrorCommand)
	Error() *customError.CustomError
	LogSuccess(command domain.NewSuccessCommand)
	Success() *message.Message
	HttpCode() valueObjects.HttpCodeReturn
	Status() valueObjects.Status
}
