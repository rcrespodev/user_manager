package message

import (
	"fmt"
	"strings"
	"time"
)

const (
	ClientErrorBadRequest ClientErrorType = 1 + iota
	ClientErrorUnauthorized
	ClientErrorFordibben
	ClientErrorNotFound
)

type Message struct {
	objectId        string
	messageId       MessageId
	messagePkg      string
	variables       Variables
	text            string
	time            time.Time
	clientErrorType ClientErrorType
}

type NewMessageCommand struct {
	ObjectId   string
	MessageId  MessageId
	MessagePkg string
	Variables  Variables
}

type MessageData struct {
	ObjectId        string          `json:"object_id"`
	MessageId       MessageId       `json:"message_id"`
	MessagePkg      string          `json:"message_pkg"`
	Variables       Variables       `json:"variables"`
	Text            string          `json:"text"`
	Time            time.Time       `json:"time"`
	ClientErrorType ClientErrorType `json:"client_error_type"`
}

type MessageId uint16

type Variables [4]string

type ClientErrorType uint8

func NewMessage(command NewMessageCommand, repository MessageRepository) (*Message, error) {
	text, clientErrorType := repository.GetMessageData(command.MessageId, command.MessagePkg)

	msg := &Message{
		objectId:        command.ObjectId,
		messageId:       command.MessageId,
		messagePkg:      command.MessagePkg,
		variables:       command.Variables,
		text:            text,
		time:            time.Now(),
		clientErrorType: clientErrorType,
	}

	if msg.text == "" {
		return nil, fmt.Errorf("message %v not found in pkg %v", command.MessageId, command.MessagePkg)
	}

	for _, variable := range msg.variables {
		msg.text = strings.Replace(msg.text, "%v", variable, 1)
	}

	return msg, nil
}

func (m Message) MessageData() *MessageData {
	return &MessageData{
		ObjectId:        m.objectId,
		MessageId:       m.messageId,
		MessagePkg:      m.messagePkg,
		Variables:       m.variables,
		Text:            m.text,
		Time:            m.time,
		ClientErrorType: m.clientErrorType,
	}
}
