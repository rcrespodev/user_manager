package customError

import (
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain/message"
	"log"
	"runtime"
)

type CustomError struct {
	internalError *InternalError
	clientError   *message.Message
}

func NewExternalError(command message.NewMessageCommand, repository message.MessageRepository) *CustomError {
	msg, err := message.NewMessage(command, repository)
	if err != nil {
		return NewInternalError(err, 3)
	}
	return &CustomError{
		internalError: nil,
		clientError:   msg,
	}
}

func NewInternalError(error error, caller int) *CustomError {
	_, file, line, _ := runtime.Caller(caller)
	customerErr := &CustomError{
		internalError: &InternalError{
			error: error,
			file:  file,
			line:  line,
		},
		clientError: nil,
	}

	// see app logs
	log.Printf("error:%v, file:%v, line:%v", customerErr.internalError.error, customerErr.internalError.file, customerErr.internalError.line)
	return customerErr
}

func (c CustomError) InternalError() *InternalError {
	return c.internalError
}

func (c CustomError) Message() *message.MessageData {
	if c.clientError == nil {
		return nil
	}
	return c.clientError.MessageData()
}
