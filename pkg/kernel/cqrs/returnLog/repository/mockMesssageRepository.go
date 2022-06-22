package repository

import (
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain/message"
)

type MockMessageRepository struct {
}

func (m MockMessageRepository) GetMessageText(id message.MessageId, messagePkg string) string {
	return ""
}
