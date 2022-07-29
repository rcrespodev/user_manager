package login

import (
	"github.com/google/uuid"
	login2 "github.com/rcrespodev/user_manager/pkg/app/user/application/commands/login"
	userDomain "github.com/rcrespodev/user_manager/pkg/app/user/domain"
	userRepository "github.com/rcrespodev/user_manager/pkg/app/user/repository"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/command"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain/message"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain/valueObjects"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/repository"
	"github.com/stretchr/testify/require"
	"log"
	"testing"
	"time"
)

var messageRepository = repository.NewMockMessageRepository([]repository.MockData{
	{
		Id:              0,
		Pkg:             "user",
		Text:            "user %v logged successful",
		ClientErrorType: 0,
	},
	{
		Id:              4,
		Pkg:             "user",
		Text:            "value %v is invalid as %v attribute",
		ClientErrorType: message.ClientErrorBadRequest,
	},
	{
		Id:              7,
		Pkg:             "user",
		Text:            "attribute %v can´t contain special characters (%v)",
		ClientErrorType: message.ClientErrorBadRequest,
	},
	{
		Id:              15,
		Pkg:             "user",
		Text:            "email, alias or password are not correct. Repeat the access data.",
		ClientErrorType: message.ClientErrorBadRequest,
	},
})

var mockUserRepository = userRepository.NewMockUserRepository()

func TestUserLogin(t *testing.T) {
	setUserMockData(mockUserRepository)
	type want struct {
		status         valueObjects.Status
		httpCodeReturn valueObjects.HttpCodeReturn
		error          error
		errorMessage   *message.MessageData
		successMessage *message.MessageData
	}
	tests := []struct {
		name string
		args login2.ClientArgs
		want want
	}{
		{
			name: "invalid user or email",
			args: login2.ClientArgs{
				AliasOrEmail: "test.test$",
				Password:     "Linux638$01",
			},
			want: want{
				status:         valueObjects.Error,
				httpCodeReturn: 400,
				error:          nil,
				errorMessage: &message.MessageData{
					ObjectId:        "test.test$",
					MessageId:       15,
					MessagePkg:      "user",
					Variables:       message.Variables{},
					Text:            "email, alias or password are not correct. Repeat the access data.",
					Time:            time.Time{},
					ClientErrorType: message.ClientErrorBadRequest,
				},
				successMessage: nil,
			},
		},
		{
			name: "invalid password",
			args: login2.ClientArgs{
				AliasOrEmail: "test@test.com",
				Password:     "without_numbers_and_special_chars",
			},
			want: want{
				status:         valueObjects.Error,
				httpCodeReturn: 400,
				error:          nil,
				errorMessage: &message.MessageData{
					ObjectId:        "test@test.com",
					MessageId:       15,
					MessagePkg:      "user",
					Variables:       message.Variables{},
					Text:            "email, alias or password are not correct. Repeat the access data.",
					Time:            time.Time{},
					ClientErrorType: message.ClientErrorBadRequest,
				},
				successMessage: nil,
			},
		},
		{
			name: "user not exists",
			args: login2.ClientArgs{
				AliasOrEmail: "user_not_found@test.com",
				Password:     "Linux638$01",
			},
			want: want{
				status:         valueObjects.Error,
				httpCodeReturn: 400,
				error:          nil,
				errorMessage: &message.MessageData{
					ObjectId:        "user_not_found@test.com",
					MessageId:       15,
					MessagePkg:      "user",
					Variables:       message.Variables{},
					Text:            "email, alias or password are not correct. Repeat the access data.",
					Time:            time.Time{},
					ClientErrorType: message.ClientErrorBadRequest,
				},
				successMessage: nil,
			},
		},
		{
			name: "user alias exists but password is incorrect",
			args: login2.ClientArgs{
				AliasOrEmail: "user_alias_exists",
				Password:     "Linux638$01",
			},
			want: want{
				status:         valueObjects.Error,
				httpCodeReturn: 400,
				error:          nil,
				errorMessage: &message.MessageData{
					ObjectId:        "user_alias_exists",
					MessageId:       15,
					MessagePkg:      "user",
					Variables:       message.Variables{},
					Text:            "email, alias or password are not correct. Repeat the access data.",
					Time:            time.Time{},
					ClientErrorType: message.ClientErrorBadRequest,
				},
				successMessage: nil,
			},
		},
		{
			name: "user email exists but password is incorrect",
			args: login2.ClientArgs{
				AliasOrEmail: "email_exists@gmail.com",
				Password:     "Linux638$01",
			},
			want: want{
				status:         valueObjects.Error,
				httpCodeReturn: 400,
				error:          nil,
				errorMessage: &message.MessageData{
					ObjectId:        "email_exists@gmail.com",
					MessageId:       15,
					MessagePkg:      "user",
					Variables:       message.Variables{},
					Text:            "email, alias or password are not correct. Repeat the access data.",
					Time:            time.Time{},
					ClientErrorType: message.ClientErrorBadRequest,
				},
				successMessage: nil,
			},
		},
		{
			name: "correct login with user alias",
			args: login2.ClientArgs{
				AliasOrEmail: "user_alias_exists",
				Password:     "Linux648$",
			},
			want: want{
				status:         valueObjects.Success,
				httpCodeReturn: 200,
				error:          nil,
				errorMessage:   nil,
				successMessage: &message.MessageData{
					ObjectId:   "user_alias_exists",
					MessageId:  0,
					MessagePkg: "user",
					Variables:  message.Variables{"user_alias_exists"},
					Text:       "user user_alias_exists logged successful",
					Time:       time.Time{},
				},
			},
		},
		{
			name: "correct login with user email",
			args: login2.ClientArgs{
				AliasOrEmail: "email_exists@gmail.com",
				Password:     "Linux648$",
			},
			want: want{
				status:         valueObjects.Success,
				httpCodeReturn: 200,
				error:          nil,
				errorMessage:   nil,
				successMessage: &message.MessageData{
					ObjectId:   "email_exists@gmail.com",
					MessageId:  0,
					MessagePkg: "user",
					Variables:  message.Variables{"user_alias_exists"},
					Text:       "user user_alias_exists logged successful",
					Time:       time.Time{},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userLoginCmd := login2.NewLoginUserCommand(tt.args)
			cmd := command.NewCommand(command.LoginUser, uuid.New(), userLoginCmd)
			retLog := domain.NewReturnLog(cmd.Uuid(), messageRepository, "user")
			userLogger := login2.NewUserLogger(mockUserRepository)
			cmdHandler := login2.NewLoginUserCommandHandler(userLogger)

			done := make(chan bool)
			go cmdHandler.Handle(*cmd, retLog, done)
			<-done

			// ReturnLog check
			require.EqualValues(t, tt.want.status, retLog.Status())
			require.EqualValues(t, tt.want.httpCodeReturn, retLog.HttpCode())

			// Check Internal error
			switch tt.want.error {
			case nil:
				if retLog.Error() != nil {
					require.Nil(t, retLog.Error().InternalError())
				}
			default:
				require.EqualValues(t, tt.want.error, retLog.Error().InternalError().Error())
			}

			// Check Client error messages
			switch tt.want.errorMessage {
			case nil:
				if retLog.Error() != nil {
					require.Nil(t, retLog.Error().Message())
				}
			default:
				gotMessage := retLog.Error().Message()
				gotMessage.Time = time.Time{}
				require.EqualValues(t, tt.want.errorMessage, gotMessage)
			}

			// Check Success message
			switch tt.want.successMessage {
			case nil:
				require.Nil(t, retLog.Success())
			default:
				gotMessage := retLog.Success().MessageData()
				gotMessage.Time = tt.want.successMessage.Time
				require.EqualValues(t, tt.want.successMessage, gotMessage)
			}
		})
	}
}

func setUserMockData(userRepository *userRepository.MockUserRepository) {
	mockDataArgs := []userDomain.NewUserCommand{
		{
			Uuid:       uuid.New().String(),
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
