package repository

import (
	emailSenderDomain "github.com/rcrespodev/user_manager/pkg/app/emailSender/domain"
	returnLog "github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain"
)

type MockSentEmailRepository struct {
	Schemas []*emailSenderDomain.SentEmailSchema
}

func NewMockSentEmailRepository(schemas []*emailSenderDomain.SentEmailSchema) *MockSentEmailRepository {
	return &MockSentEmailRepository{Schemas: schemas}
}

func (m *MockSentEmailRepository) Save(schema *emailSenderDomain.SentEmailSchema, log *returnLog.ReturnLog) {
	m.Schemas = append(m.Schemas, schema)
}

func (m *MockSentEmailRepository) Get(userUuid string) []*emailSenderDomain.SentEmailSchema {
	var schemas []*emailSenderDomain.SentEmailSchema
	if m.Schemas == nil {
		return nil
	}

	for _, schema := range m.Schemas {
		if schema.UserUuid == userUuid {
			schemas = append(schemas, schema)
		}
	}
	return schemas
}
