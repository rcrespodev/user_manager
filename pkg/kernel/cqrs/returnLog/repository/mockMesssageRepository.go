package repository

import (
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain/message"
)

type MockMessageRepository struct {
	data []MockData
}

type MockData struct {
	Id   message.MessageId
	Pkg  string
	Text string
}

func NewMockMessageRepository(data []MockData) *MockMessageRepository {
	return &MockMessageRepository{data: data}
}

func (m MockMessageRepository) GetMessageText(id message.MessageId, messagePkg string) string {
	for _, data := range m.data {
		if data.Id == id && data.Pkg == messagePkg {
			return data.Text
		}
	}
	return ""
}
