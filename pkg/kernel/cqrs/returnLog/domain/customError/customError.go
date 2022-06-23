package customError

import (
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain/message"
	"runtime"
)

type CustomError struct {
	internalError *InternalError
	message       *message.Message
}

func NewExternalError(command message.NewMessageCommand, repository message.MessageRepository) *CustomError {
	msg, err := message.NewMessage(command, repository)
	if err != nil {
		return NewInternalError(err, 3)
	}
	return &CustomError{
		internalError: nil,
		message:       msg,
	}
}

func NewInternalError(error error, caller int) *CustomError {
	_, file, line, _ := runtime.Caller(caller)
	return &CustomError{
		internalError: &InternalError{
			error: error,
			file:  file,
			line:  line,
		},
		message: nil,
	}
}

func (c CustomError) InternalError() *InternalError {
	return c.internalError
}

func (c CustomError) Message() *message.MessageData {
	if c.message == nil {
		return nil
	}
	return c.message.MessageData()
}
