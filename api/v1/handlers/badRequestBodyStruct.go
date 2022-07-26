package handlers

import (
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain/message"
	"time"
)

func BodyRequestBadType() message.MessageData {
	return message.MessageData{
		ObjectId:        "",
		MessageId:       1,
		MessagePkg:      "http handler",
		Variables:       message.Variables{},
		Text:            "body request has´t correct type",
		Time:            time.Now(),
		ClientErrorType: message.ClientErrorBadRequest,
	}
}
