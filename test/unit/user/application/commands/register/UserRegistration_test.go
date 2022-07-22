package register

import (
	uuid "github.com/google/uuid"
	register2 "github.com/rcrespodev/user_manager/pkg/app/user/application/commands/register"
	userDomain "github.com/rcrespodev/user_manager/pkg/app/user/domain"
	userRepository "github.com/rcrespodev/user_manager/pkg/app/user/repository"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/command"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain/message"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain/valueObjects"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/repository"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

var messageRepository = repository.NewMockMessageRepository([]repository.MockData{
	{
		Id:              001,
		Pkg:             "user",
		Text:            "user %v created successful",
		ClientErrorType: 0,
	},
	{
		Id:              14,
		Pkg:             "user",
		Text:            "user with component: %v and value: %v already exists",
		ClientErrorType: message.ClientErrorBadRequest,
	},
})

func TestUserRegistration(t *testing.T) {
	mockRepository := userRepository.NewMockUserRepository()
	setMockUserData(mockRepository)

	type args struct {
		alias      string
		name       string
		secondName string
		email      string
		password   string
	}
	type want struct {
		status         valueObjects.Status
		httpCodeReturn valueObjects.HttpCodeReturn
		error          error
		errorMessage   *message.MessageData
		successMessage *message.MessageData
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "correct request",
			args: args{
				alias:      "martin_fowler",
				name:       "martin",
				secondName: "fowler",
				email:      "foo@test.com",
				password:   "Linux648$",
			},
			want: want{
				status:         valueObjects.Success,
				httpCodeReturn: 200,
				error:          nil,
				errorMessage:   nil,
				successMessage: &message.MessageData{
					ObjectId:        "martin_fowler",
					MessageId:       1,
					MessagePkg:      "user",
					Variables:       message.Variables{"martin_fowler"},
					Text:            "user martin_fowler created successful",
					Time:            time.Time{},
					ClientErrorType: 0,
				},
			},
		},
		{
			name: "bad request - user alias already exists",
			args: args{
				alias:      "user_exists",
				name:       "martin",
				secondName: "fowler",
				email:      "foo@test.com",
				password:   "Linux648$",
			},
			want: want{
				status:         valueObjects.Error,
				httpCodeReturn: 400,
				error:          nil,
				errorMessage: &message.MessageData{
					ObjectId:        "user_exists",
					MessageId:       14,
					MessagePkg:      "user",
					Variables:       message.Variables{"alias", "user_exists"},
					Text:            "user with component: alias and value: user_exists already exists",
					Time:            time.Time{},
					ClientErrorType: message.ClientErrorBadRequest,
				},
				successMessage: nil,
			},
		},
		{
			name: "bad request - user email already exists",
			args: args{
				alias:      "martin_fowler_2",
				name:       "martin",
				secondName: "fowler",
				email:      "email_exists@test.com",
				password:   "Linux648$",
			},
			want: want{
				status:         valueObjects.Error,
				httpCodeReturn: 400,
				error:          nil,
				errorMessage: &message.MessageData{
					ObjectId:        "martin_fowler_2",
					MessageId:       14,
					MessagePkg:      "user",
					Variables:       message.Variables{"email", "email_exists@test.com"},
					Text:            "user with component: email and value: email_exists@test.com already exists",
					Time:            time.Time{},
					ClientErrorType: message.ClientErrorBadRequest,
				},
				successMessage: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			done := make(chan bool)
			cmdUuid := uuid.New()

			registerUserCommand := register2.NewRegisterUserCommand(register2.ClientArgs{
				Uuid:       cmdUuid.String(),
				Alias:      tt.args.alias,
				Name:       tt.args.name,
				SecondName: tt.args.secondName,
				Email:      tt.args.email,
				Password:   tt.args.password,
			})

			cmd := command.NewCommand(command.RegisterUser, cmdUuid, registerUserCommand)
			userRegistration := register2.NewUserRegistration(mockRepository)
			handler := register2.NewRegisterUserCommandHandler(userRegistration)
			retLog := domain.NewReturnLog(cmd.Uuid(), messageRepository, "user")
			go handler.Handle(*cmd, retLog, done)

			println(<-done)

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

			if tt.want.status == valueObjects.Success {
				findUserCmd := userDomain.FindUserQuery{
					Log: retLog,
					Where: []userDomain.WhereArgs{
						{
							Field: "uuid",
							Value: cmdUuid.String(),
						},
					},
				}

				if user := mockRepository.FindUser(findUserCmd); user == nil {
					t.Errorf("FindUserById()\n\t- User not found in repository!!")
				}
			}
		})
	}
}

func setMockUserData(userRepository *userRepository.MockUserRepository) {
	mockDataArgs := []userDomain.NewUserCommand{
		{
			Uuid:       uuid.New().String(),
			Alias:      "user_exists",
			Name:       "martin",
			SecondName: "fowler",
			Email:      "email_exists@test.com",
			Password:   "Linux648$",
		},
	}
	var users []*userDomain.User
	for _, arg := range mockDataArgs {
		user := userDomain.NewUser(arg, domain.NewReturnLog(uuid.New(), messageRepository, "user"))
		users = append(users, user)
	}
	userRepository.SetMockData(users)
}
