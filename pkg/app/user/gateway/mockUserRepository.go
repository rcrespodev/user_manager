package gateway

import (
	"github.com/google/uuid"
	"github.com/rcrespodev/user_manager/pkg/app/user/domain"
	returnLog "github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain"
	"sync"
)

type MockUserRepository struct {
	userMockData UserMockData
}

type UserMockData struct {
	users []*domain.User
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{}
}

func (m *MockUserRepository) SetMockData(users []*domain.User) {
	m.userMockData.users = users
}

func (m *MockUserRepository) SaveUser(user *domain.User, log *returnLog.ReturnLog, wg *sync.WaitGroup) {
	m.userMockData.users = append(m.userMockData.users, user)
	wg.Done()
}

func (m MockUserRepository) FindUserById(uuid uuid.UUID, log *returnLog.ReturnLog, wg *sync.WaitGroup) *domain.User {
	for _, user := range m.userMockData.users {
		if uuid == user.Uuid() {
			wg.Done()
			return user
		}
	}
	wg.Done()
	return nil
}

func (m MockUserRepository) FindUserByEmail(email *domain.UserEmail, log *returnLog.ReturnLog, wg *sync.WaitGroup) *domain.User {
	for _, user := range m.userMockData.users {
		if email.Address() == user.Email().Address() {
			wg.Done()
			return user
		}
	}
	wg.Done()
	return nil
}

func (m MockUserRepository) FindUserByAlias(alias *domain.UserAlias, log *returnLog.ReturnLog, wg *sync.WaitGroup) *domain.User {
	for _, user := range m.userMockData.users {
		if alias.Alias() == user.Alias().Alias() {
			wg.Done()
			return user
		}
	}
	wg.Done()
	return nil
}
