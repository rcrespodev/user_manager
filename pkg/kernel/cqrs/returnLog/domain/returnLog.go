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

// ReturnLog
// Every command/query can be Success or Error. Log the status in the ReturnLog
// using LogError() and LogSuccess. If log has one Error and one Success, the error
// is more important.
// Also, if error has one internal error and one external error, the internal is more
// important
type ReturnLog struct {
	uuid            uuid.UUID           //uuid of command/query
	status          valueObjects.Status //Success or Error
	defaultPkg      string
	error           *customError.CustomError
	success         *message.Message
	httpCodeReturn  valueObjects.HttpCodeReturn
	currentObjectId string
	repository      message.MessageRepository
	caller          int
}

// NewErrorCommand
// If Error != nil, returnLog exclude the message data and only log the error as internal Error
// If Error == nil, returnLog use the messageData and log the error
type NewErrorCommand struct {
	Error error
	*message.NewMessageCommand
	Caller      int
	Overwritten bool
}

type NewSuccessCommand *message.NewMessageCommand

func NewReturnLog(uuid uuid.UUID, repository message.MessageRepository, defaultPkg string) *ReturnLog {
	return &ReturnLog{
		uuid:            uuid,
		defaultPkg:      defaultPkg,
		repository:      repository,
		caller:          defaultCaller,
		error:           nil,
		success:         nil,
		currentObjectId: "",
		status:          "",
		httpCodeReturn:  000,
	}
}

// LogError
// See NewErrorCommand usage.
// The log can only contain one error. If the log already contain an error, no new
// errors are logged, that is, the original error is not overwritten.
// If you want overwritten the original message set NewErrorCommand.overwritten=true
func (r *ReturnLog) LogError(command NewErrorCommand) {
	defer func() {
		r.updateInternalData()
	}()

	var caller int
	if command.Caller != 0 {
		caller = command.Caller
	} else {
		caller = r.caller
	}

	if command.Error != nil && command.Overwritten == false {
		r.error = customError.NewInternalError(command.Error, caller)
		return
	}

	pkg := r.getPkg(command.MessagePkg)
	objectId := r.getObjectId(command.ObjectId)
	r.error = customError.NewExternalError(message.NewMessageCommand{
		ObjectId:   objectId,
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

// LogSuccess
// The log can only contain one Success. If the log already contain a success, no new
// success are logged, that is, the original success is not overwritten.
func (r *ReturnLog) LogSuccess(command NewSuccessCommand) {
	defer func() {
		r.updateInternalData()
	}()

	if r.success != nil {
		return
	}

	pkg := r.getPkg(command.MessagePkg)
	objectId := r.getObjectId(command.ObjectId)
	msg, err := message.NewMessage(message.NewMessageCommand{
		ObjectId:   objectId,
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
		if r.success != nil {
			r.success = nil
		}
		if r.status != valueObjects.Error {
			r.status = valueObjects.Error
		}
		if r.error.InternalError() != nil {
			r.httpCodeReturn = valueObjects.HttpCodeInternalError
			return
		}
		if r.error.Message() != nil {
			httpCode, err := r.httpCodeReturn.MapClientErrorToHttpCode(r.error.Message().ClientErrorType)
			//if httpCode == 0 {
			//	return
			//}
			if err != nil {
				r.LogError(NewErrorCommand{Error: err})
				return
			}
			r.httpCodeReturn = httpCode
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

func (r *ReturnLog) SetObjectId(objectId string) {
	r.currentObjectId = objectId
}

func (r ReturnLog) getObjectId(cmdObjectId string) string {
	if r.currentObjectId != "" {
		return r.currentObjectId
	}

	return cmdObjectId
}

func (r *ReturnLog) CurrentObjectId() string {
	return r.currentObjectId
}
