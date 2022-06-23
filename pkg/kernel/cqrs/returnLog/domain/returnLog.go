package domain

import (
	"github.com/google/uuid"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain/customError"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain/message"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain/valueObjects"
)

const (
	defaultCaller = 2
)

type ReturnLog struct {
	uuid           uuid.UUID           //uuid of command/query
	status         valueObjects.Status //Success or Error
	defaultPkg     string
	error          *customError.CustomError
	success        *message.Message
	httpCodeReturn valueObjects.HttpCodeReturn
	repository     message.MessageRepository
	caller         int
}

// NewErrorCommand
// If Error != nil, returnLog exclude the message data and only log the error as internal Error
// If Error == nil, returnLog use the messageData and log the error
type NewErrorCommand struct {
	Error error
	message.NewMessageCommand
}

type NewSuccessCommand message.NewMessageCommand

func NewReturnLog(uuid uuid.UUID, repository message.MessageRepository, defaultPkg string) *ReturnLog {
	return &ReturnLog{
		uuid:           uuid,
		defaultPkg:     defaultPkg,
		repository:     repository,
		caller:         defaultCaller,
		httpCodeReturn: valueObjects.HttpCodeSuccess,
	}
}

func (r *ReturnLog) LogError(command NewErrorCommand) {
	defer func() {
		r.updateInternalData()
	}()

	if command.Error != nil {
		r.error = customError.NewInternalError(command.Error, r.caller)
		return
	}

	pkg := r.getPkg(command.MessagePkg)
	r.error = customError.NewExternalError(message.NewMessageCommand{
		ObjectId:   command.ObjectId,
		MessageId:  command.MessageId,
		MessagePkg: pkg,
		Variables:  command.Variables,
	}, r.repository)
}

func (r ReturnLog) Error() *customError.CustomError {
	if r.error == nil {
		return nil
	}
	return r.error
}

func (r *ReturnLog) LogSuccess(command NewSuccessCommand) {
	defer func() {
		r.updateInternalData()
	}()

	pkg := r.getPkg(command.MessagePkg)
	msg, err := message.NewMessage(message.NewMessageCommand{
		ObjectId:   command.ObjectId,
		MessageId:  command.MessageId,
		MessagePkg: pkg,
		Variables:  command.Variables,
	}, r.repository)

	if err != nil {
		r.caller = 3
		r.LogError(NewErrorCommand{
			Error: err,
		})
		r.caller = defaultCaller
		return
	}

	r.success = msg
}

func (r ReturnLog) Success() *message.Message {
	if r.success == nil {
		return nil
	}
	return r.success
}

func (r ReturnLog) HttpCode() valueObjects.HttpCodeReturn {
	return r.httpCodeReturn
}

func (r *ReturnLog) Status() valueObjects.Status {
	return r.status
}

func (r ReturnLog) getPkg(commandPkg string) string {
	if commandPkg == "" {
		return r.defaultPkg
	}
	return commandPkg
}

func (r *ReturnLog) updateInternalData() {
	if r.error != nil {
		if r.status != valueObjects.Error {
			r.status = valueObjects.Error
		}
		if r.error.InternalError() != nil {
			r.httpCodeReturn = valueObjects.HttpCodeInternalError
			return
		}
		if r.error.Message() != nil {
			r.httpCodeReturn = valueObjects.HttpCodeBadRequest
			return
		}
	}

	if r.status != valueObjects.Success {
		r.status = valueObjects.Success
	}
	if r.httpCodeReturn != valueObjects.HttpCodeSuccess {
		r.httpCodeReturn = valueObjects.HttpCodeSuccess
	}
}
