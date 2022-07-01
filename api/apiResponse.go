package api

import "github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain/message"

type QueryResponse struct {
	Message message.MessageData `json:"message"`
	Data    interface{}         `json:"data"`
}

type CommandResponse struct {
	Message message.MessageData `json:"message"`
}
