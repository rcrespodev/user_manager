package api

import (
	returnLog "github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain/message"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain/valueObjects"
)

type QueryResponse struct {
	Message message.MessageData `json:"message"`
	Data    interface{}         `json:"data"`
}

type CommandResponse struct {
	Message message.MessageData `json:"message"`
}

func NewCommandResponse(log *returnLog.ReturnLog) *CommandResponse {
	response := &CommandResponse{}
	switch log.Status() {
	case valueObjects.Error:
		if log.Error().InternalError() != nil {
			response.Message = message.MessageData{}
		} else {
			response.Message = *log.Error().Message()
		}
	case valueObjects.Success:
		response.Message = *log.Success().MessageData()
	}
	return response
}

func NewQueryResponse(log *returnLog.ReturnLog, data interface{}) *QueryResponse {
	response := &QueryResponse{
		Data: data,
	}
	switch log.Status() {
	case valueObjects.Error:
		if log.Error().InternalError() != nil {
			response.Message = message.MessageData{}
		} else {
			response.Message = *log.Error().Message()
		}
		//case valueObjects.Success:
		//	response.Message = *log.Success().MessageData()
	}
	return response
}
