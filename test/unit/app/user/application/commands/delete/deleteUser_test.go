package delete

import (
	"github.com/google/uuid"
	deleteApp "github.com/rcrespodev/user_manager/pkg/app/user/application/commands/delete"
	userDomain "github.com/rcrespodev/user_manager/pkg/app/user/domain"
	"github.com/rcrespodev/user_manager/pkg/app/user/repository/userRepository"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/command"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain/message"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain/valueObjects"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/repository"
	"github.com/rcrespodev/user_manager/test/unit/app/user/application/utils"
	"log"
	"testing"
)

var messageRepository = repository.NewMockMessageRepository([]repository.MockData{
	{
		Id:              3,
		Pkg:             "user",
		Text:            "user %v deleted successful",
		ClientErrorType: 0,
	},
	{
		Id:              4,
		Pkg:             "user",
		Text:            "value %v is invalid as %v attribute",
		ClientErrorType: 1,
	},
	{
		Id:              17,
		Pkg:             "user",
		Text:            "none of the input values correspond to a registered user",
		ClientErrorType: 1,
	},
})

var mockUserRepository = userRepository.NewMockUserRepository()

func TestDeleteUser(t *testing.T) {
	setUserMockData(mockUserRepository)

	type args struct {
		userUuid string
	}
	tests := []struct {
		name string
		args args
		want utils.Want
	}{
		{
			name: "user deleted successful",
			args: args{
				userUuid: "123e4567-e89b-12d3-a456-426614174000",
			},
			want: utils.Want{
				Status:         valueObjects.Success,
				HttpCodeReturn: 200,
				Error:          nil,
				ErrorMessage:   nil,
				SuccessMessage: &message.MessageData{
					ObjectId:   "123e4567-e89b-12d3-a456-426614174000",
					MessageId:  3,
					MessagePkg: "user",
					Variables:  message.Variables{"user_alias_exists"},
					Text:       "user user_alias_exists deleted successful",
				},
			},
		},
		{
			name: "invalid user uuid - not exists in db",
			args: args{
				userUuid: "123e4567-e89b-12d3-a456-426614174",
			},
			want: utils.Want{
				Status:         valueObjects.Error,
				HttpCodeReturn: 400,
				ErrorMessage: &message.MessageData{
					ObjectId:        "123e4567-e89b-12d3-a456-426614174",
					MessageId:       17,
					MessagePkg:      "user",
					Text:            "none of the input values correspond to a registered user",
					ClientErrorType: 1,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			deleteUserCommand := deleteApp.NewDeleteUserCommand(tt.args.userUuid)
			cmd := command.NewCommand(command.DeleteUser, uuid.New(), deleteUserCommand)
			retLog := domain.NewReturnLog(cmd.Uuid(), messageRepository, "user")
			userDeleter := deleteApp.NewUserDeleter(mockUserRepository)
			cmdHandler := deleteApp.NewDeleteUserCommandHandler(userDeleter)

			done := make(chan bool)
			go cmdHandler.Handle(*cmd, retLog, done)
			<-done
			tt.want.TestResponse(t, retLog)
		})
	}
}

func setUserMockData(userRepository *userRepository.MockUserRepository) {
	mockDataArgs := []userDomain.NewUserCommand{
		{
			Uuid:       "123e4567-e89b-12d3-a456-426614174000",
			Alias:      "user_alias_exists",
			Name:       "martin",
			SecondName: "fowler",
			Email:      "email_exists@gmail.com",
			Password:   "Linux648$",
		},
	}
	var users []*userDomain.User
	for _, arg := range mockDataArgs {
		user := userDomain.NewUser(arg, domain.NewReturnLog(uuid.New(), messageRepository, "user"))
		if user == nil {
			log.Fatal("invalid user data")
		}
		users = append(users, user)
	}
	userRepository.SetMockData(users)
}
