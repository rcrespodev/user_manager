package message

import (
	"fmt"
	"time"
)

type Message struct {
	objectId   string
	messageId  MessageId
	messagePkg string
	variables  Variables
	text       string
	time       time.Time
}

type NewMessageCommand struct {
	ObjectId   string
	MessageId  MessageId
	MessagePkg string
	Variables  Variables
}

type MessageId uint16

type Variables [4]string

type MessageData struct {
	ObjectId   string
	MessageId  MessageId
	MessagePkg string
	Variables  Variables
	Text       string
	Time       time.Time
}

func NewMessage(command NewMessageCommand, repository MessageRepository) (*Message, error) {
	text := repository.GetMessageText(command.MessageId, command.MessagePkg)

	msg := &Message{
		objectId:   command.ObjectId,
		messageId:  command.MessageId,
		messagePkg: command.MessagePkg,
		variables:  command.Variables,
		text:       text,
		time:       time.Now(),
	}

	if msg.text == "" {
		return nil, fmt.Errorf("message %v not found in pkg %v", command.MessageId, command.MessagePkg)
	}

	msg.text = fmt.Sprintf(msg.text, msg.variables)

	return msg, nil
}

func (m Message) MessageData() *MessageData {
	return &MessageData{
		ObjectId:   m.objectId,
		MessageId:  m.messageId,
		MessagePkg: m.messagePkg,
		Variables:  m.variables,
		Text:       m.text,
		Time:       m.time,
	}
}
