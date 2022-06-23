package repository

import (
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain/message"
)

type MockMessageRepository struct {
	data []MockData
}

type MockData struct {
	Id              message.MessageId
	Pkg             string
	Text            string
	ClientErrorType message.ClientErrorType
}

func NewMockMessageRepository(data []MockData) *MockMessageRepository {
	return &MockMessageRepository{data: data}
}

func (m MockMessageRepository) GetMessageData(id message.MessageId, messagePkg string) (text string, clientErrorType message.ClientErrorType) {
	for _, data := range m.data {
		if data.Id == id && data.Pkg == messagePkg {
			return data.Text, data.ClientErrorType
		}
	}
	return "", 0
}
