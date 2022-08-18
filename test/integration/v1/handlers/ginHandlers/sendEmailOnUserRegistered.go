package handlers

import (
	"github.com/google/uuid"
	"github.com/rcrespodev/user_manager/pkg/app/emailSender/application/commands"
	userDomain "github.com/rcrespodev/user_manager/pkg/app/user/domain"
	"github.com/rcrespodev/user_manager/pkg/kernel"
	retLogDomain "github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain/message"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain/valueObjects"
	utils2 "github.com/rcrespodev/user_manager/test/integration/v1/handlers/ginHandlers/utils"
	"github.com/rcrespodev/user_manager/test/unit/app/user/application/utils"
	"github.com/stretchr/testify/require"
	"log"
	"testing"
)

func TestSendEmailOnUserRegistered(t *testing.T) {
	require.NoError(t, setUserMockData(kernel.Instance.UserRepository()))

	type args struct {
		aggregateId  string // user Uuid
		templatePath string
	}
	tests := []struct {
		name string
		args args
		want *utils.Want
	}{
		{
			name: "good request - user uuid exists and email send successful",
			args: args{
				aggregateId:  "123e4567-e89b-12d3-a456-426614174008",
				templatePath: "welcomeTemplate.txt",
			},
			want: &utils.Want{
				Status:         valueObjects.Success,
				HttpCodeReturn: 200,
				Error:          nil,
				ErrorMessage:   nil,
				SuccessMessage: &message.MessageData{
					ObjectId:        "123e4567-e89b-12d3-a456-426614174008",
					MessageId:       0,
					MessagePkg:      "email_sender",
					Variables:       message.Variables{"foo_test@gmail.com"},
					Text:            "welcome email send successful to user email foo_test@gmail.com",
					ClientErrorType: 0,
				},
			},
		},
		{
			name: "bad request - user uuid not exists",
			args: args{
				aggregateId: "123e4567-e89b-12d3-a456-426614174000",
			},
			want: &utils.Want{
				Status:         valueObjects.Error,
				HttpCodeReturn: 400,
				Error:          nil,
				ErrorMessage: &message.MessageData{
					ObjectId:        "123e4567-e89b-12d3-a456-426614174000",
					MessageId:       17,
					MessagePkg:      "user",
					Text:            "none of the input values correspond to a registered user",
					ClientErrorType: message.ClientErrorBadRequest,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := commands.NewSendEmailOnUserRegisteredCommand(tt.args.aggregateId)
			returnLog := retLogDomain.NewReturnLog(cmd.BaseCommand().CommandUuid(), kernel.Instance.MessageRepository(), "email_sender")
			kernel.Instance.CommandBus().Exec(cmd, returnLog)

			//if tt.want.Error != nil {
			//	tt.want.Error = returnLog.Error().InternalError().Error()
			//}
			tt.want.TestResponse(t, returnLog)

			if returnLog.Status() == valueObjects.Success {
				schemas := kernel.Instance.SentEmailRepository().Get(tt.args.aggregateId)
				require.NotNil(t, schemas)
				for _, schema := range schemas {
					require.EqualValues(t, tt.args.aggregateId, schema.UserUuid)
					require.EqualValues(t, true, schema.Sent)
				}
			}

			if returnLog.Status() == valueObjects.Error {
				schemas := kernel.Instance.SentEmailRepository().Get(tt.args.aggregateId)
				require.Nil(t, schemas)
			}
		})
	}
}

func setUserMockData(userRepository userDomain.UserRepository) error {
	mockDataArgs := []*userDomain.NewUserCommand{
		{
			Uuid:       "123e4567-e89b-12d3-a456-426614174008",
			Alias:      "user_alias_exists",
			Name:       "martin",
			SecondName: "fowler",
			Email:      "foo_test@gmail.com",
			Password:   "Linux648$",
		},
		{
			Uuid:       "123e4567-e89b-12d3-a456-426614174009",
			Alias:      "user_alias_exists_2",
			Name:       "martin",
			SecondName: "fowler",
			Email:      "foo_test_2@gmail.com",
			Password:   "Linux648$",
		},
	}
	var users []*userDomain.User
	for _, arg := range mockDataArgs {
		user := userDomain.NewUser(*arg, retLogDomain.NewReturnLog(uuid.New(), kernel.Instance.MessageRepository(), "user"))
		if user == nil {
			log.Fatal("invalid user data")
		}
		users = append(users, user)
	}
	return utils2.TableUsersSetup(mockDataArgs, userRepository)
}
