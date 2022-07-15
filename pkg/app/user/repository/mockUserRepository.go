package repository

import (
	"github.com/rcrespodev/user_manager/pkg/app/user/domain"
	returnLog "github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain"
)

type MockUserRepository struct {
	userMockData UserMockData
}

func (m *MockUserRepository) FindUser(command domain.FindUserCommand) *domain.User {
	var user *domain.User
	for _, args := range command.Where {
		value := args.Value
		switch args.Field {
		case "uuid":
			user = m.findUserById(value)
		case "email":
			user = m.findUserByEmail(value)
		case "alias":
			user = m.findUserByAlias(value)
		}
		if user != nil {
			break
		}
	}

	return user
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

func (m *MockUserRepository) SaveUser(user *domain.User, log *returnLog.ReturnLog) {
	m.userMockData.users = append(m.userMockData.users, user)
}

func (m MockUserRepository) findUserById(uuid string) *domain.User {
	for _, user := range m.userMockData.users {
		if uuid == user.Uuid().String() {
			return user
		}
	}
	return nil
}

func (m MockUserRepository) findUserByEmail(email string) *domain.User {
	for _, user := range m.userMockData.users {
		if email == user.Email().Address() {
			return user
		}
	}
	return nil
}

func (m MockUserRepository) findUserByAlias(alias string) *domain.User {
	for _, user := range m.userMockData.users {
		if alias == user.Alias().Alias() {
			return user
		}
	}
	return nil
}
