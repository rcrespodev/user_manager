package repository

import (
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

func (m MockUserRepository) FindUserById(command domain.FindByIdCommand) *domain.User {
	for _, user := range m.userMockData.users {
		if command.Uuid == user.Uuid() {
			command.FindUserCommand.Wg.Done()
			return user
		}
	}
	command.FindUserCommand.Wg.Done()
	return nil
}

func (m MockUserRepository) FindUserByEmail(command domain.FindByEmailCommand) *domain.User {
	for _, user := range m.userMockData.users {
		if command.Email.Address() == user.Email().Address() {
			command.FindUserCommand.Wg.Done()
			return user
		}
	}
	command.FindUserCommand.Wg.Done()
	return nil
}

func (m MockUserRepository) FindUserByAlias(command domain.FindByAliasCommand) *domain.User {
	for _, user := range m.userMockData.users {
		if command.Alias.Alias() == user.Alias().Alias() {
			command.FindUserCommand.Wg.Done()
			return user
		}
	}
	command.FindUserCommand.Wg.Done()
	return nil
}
