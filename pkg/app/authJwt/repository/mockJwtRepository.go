package repository

import (
	"github.com/rcrespodev/user_manager/pkg/app/authJwt/domain"
	returnLog "github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain"
)

type MockJwtRepository struct {
	MockData MockData
}

type MockData map[string]*domain.JwtSchema

func NewMockJwtRepository(mockData MockData) *MockJwtRepository {
	return &MockJwtRepository{MockData: mockData}
}
func (m *MockJwtRepository) Update(command domain.UpdateCommand, log *returnLog.ReturnLog) {
	m.MockData[command.Command.Uuid] = command.Command
}

func (m *MockJwtRepository) FindByUuid(query domain.FindByUuidQuery) *domain.JwtSchema {
	v, ok := m.MockData[query.Uuid]
	if !ok {
		return nil
	}

	return v
}
