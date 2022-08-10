package userRepository

import (
	"github.com/rcrespodev/user_manager/pkg/app/user/domain"
	returnLog "github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain/message"
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
			HashedPassword: user.Password().Hash(),
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
		HashedPassword: user.Password().Hash(),
	}
	m.userMockData.users = append(m.userMockData.users, userSchema)
}

func (m *MockUserRepository) FindUser(query domain.FindUserQuery) *domain.UserSchema {
	var userSchema *domain.UserSchema
	for _, args := range query.Where {
		value := args.Value
		switch args.Field {
		case "uuid":
			userSchema = m.findUserById(value)
		case "email":
			userSchema = m.findUserByEmail(value)
		case "alias":
			userSchema = m.findUserByAlias(value)
		case "name":
			userSchema = m.findUserByName(value)
		case "second_name":
			userSchema = m.findUserBySecondName(value)
		}
		if userSchema != nil {
			break
		}
	}

	if userSchema == nil {
		query.Log.LogError(returnLog.NewErrorCommand{
			NewMessageCommand: &message.NewMessageCommand{
				MessageId:  17,
				MessagePkg: "user",
			},
		})
	}

	return userSchema
}

func (m *MockUserRepository) findUserById(uuid string) *domain.UserSchema {
	for _, userSchema := range m.userMockData.users {
		if uuid == userSchema.Uuid {
			return userSchema
		}
	}
	return nil
}

func (m *MockUserRepository) findUserByEmail(email string) *domain.UserSchema {
	for _, userSchema := range m.userMockData.users {
		if email == userSchema.Email {
			return userSchema
		}
	}
	return nil
}

func (m *MockUserRepository) findUserByAlias(alias string) *domain.UserSchema {
	for _, userSchema := range m.userMockData.users {
		if alias == userSchema.Alias {
			return userSchema
		}
	}
	return nil
}

func (m *MockUserRepository) findUserByName(name string) *domain.UserSchema {
	for _, userSchema := range m.userMockData.users {
		if name == userSchema.Name {
			return userSchema
		}
	}
	return nil
}

func (m *MockUserRepository) findUserBySecondName(secondName string) *domain.UserSchema {
	for _, userSchema := range m.userMockData.users {
		if secondName == userSchema.SecondName {
			return userSchema
		}
	}
	return nil
}
