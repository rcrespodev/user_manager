package userSessionRepository

import (
	"github.com/rcrespodev/user_manager/pkg/app/user/domain"
	returnLog "github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain"
	"time"
)

type MockUserSessionRepository struct {
	mockData MockData
}

type MockData map[string]SessionData

type SessionData struct {
	IsLogged     bool
	LastLoginOn  time.Time
	LastLogoutOn time.Time
}

func NewMockUserSessionRepository(mockData MockData) *MockUserSessionRepository {
	return &MockUserSessionRepository{mockData: mockData}
}

func (m *MockUserSessionRepository) UpdateUserSession(command domain.UpdateUserSessionCommand, log *returnLog.ReturnLog) {
	m.mockData[command.UserUuid] = SessionData{
		IsLogged:     command.IsLogged,
		LastLoginOn:  command.LastLoginOn,
		LastLogoutOn: command.LastLogoutOn,
	}
}

func (m MockUserSessionRepository) GetUserSession(userUuid string) *domain.UserSessionSchema {
	session, ok := m.mockData[userUuid]
	if !ok {
		return nil
	}
	return &domain.UserSessionSchema{
		UserUuid:     userUuid,
		IsLogged:     session.IsLogged,
		LastLoginOn:  session.LastLoginOn,
		LastLogoutOn: session.LastLogoutOn,
	}
}

func (m *MockUserSessionRepository) ClearAll() {
	m.mockData = nil
}
