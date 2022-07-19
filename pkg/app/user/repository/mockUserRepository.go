package repository

import (
	"github.com/rcrespodev/user_manager/pkg/app/user/domain"
	returnLog "github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain"
)

type MockUserRepository struct {
	userMockData UserMockData
}

type UserMockData struct {
	users []*domain.UserSchema
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{}
}

func (m *MockUserRepository) SetMockData(users []*domain.User) {
	for _, user := range users {
		userSchema := &domain.UserSchema{
			Uuid:           user.Uuid().String(),
			Alias:          user.Alias().Alias(),
			Name:           user.Name().Name(),
			SecondName:     user.SecondName().Name(),
			Email:          user.Email().Address(),
			HashedPassword: user.Password().String(),
		}
		m.userMockData.users = append(m.userMockData.users, userSchema)
	}
}

func (m *MockUserRepository) SaveUser(user *domain.User, log *returnLog.ReturnLog) {
	userSchema := &domain.UserSchema{
		Uuid:           user.Uuid().String(),
		Alias:          user.Alias().Alias(),
		Name:           user.Name().Name(),
		SecondName:     user.SecondName().Name(),
		Email:          user.Email().Address(),
		HashedPassword: user.Password().String(),
	}
	m.userMockData.users = append(m.userMockData.users, userSchema)
}

func (m *MockUserRepository) FindUser(command domain.FindUserCommand) *domain.UserSchema {
	var userSchema *domain.UserSchema
	for _, args := range command.Where {
		value := args.Value
		switch args.Field {
		case "uuid":
			userSchema = m.findUserById(value)
		case "email":
			userSchema = m.findUserByEmail(value)
		case "alias":
			userSchema = m.findUserByAlias(value)
		}
		if userSchema != nil {
			break
		}
	}

	return userSchema
}

func (m MockUserRepository) findUserById(uuid string) *domain.UserSchema {
	for _, userSchema := range m.userMockData.users {
		if uuid == userSchema.Uuid {
			return userSchema
		}
	}
	return nil
}

func (m MockUserRepository) findUserByEmail(email string) *domain.UserSchema {
	for _, userSchema := range m.userMockData.users {
		if email == userSchema.Email {
			return userSchema
		}
	}
	return nil
}

func (m MockUserRepository) findUserByAlias(alias string) *domain.UserSchema {
	for _, userSchema := range m.userMockData.users {
		if alias == userSchema.Alias {
			return userSchema
		}
	}
	return nil
}
