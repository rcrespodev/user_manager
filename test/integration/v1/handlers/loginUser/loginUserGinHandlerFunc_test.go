package loginUser

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/rcrespodev/user_manager/api"
	"github.com/rcrespodev/user_manager/api/v1/endpoints"
	"github.com/rcrespodev/user_manager/api/v1/handlers/loginUser"
	"github.com/rcrespodev/user_manager/pkg/app/user/application/commands/login"
	"github.com/rcrespodev/user_manager/pkg/app/user/domain"
	"github.com/rcrespodev/user_manager/pkg/kernel"
	returnLog "github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain/message"
	"github.com/rcrespodev/user_manager/test/integration"
	"github.com/stretchr/testify/require"
	"log"
	"net/http"
	"testing"
	"time"
)

const (
	relPath = endpoints.EndpointLogin
)

var userRepositoryInstance domain.UserRepository

func TestLoginUserGinHandlerFunc(t *testing.T) {
	userRepositoryInstance = kernel.Instance.UserRepository()

	tableUsersSetup()

	mockGinSrv := integration.NewTestServerHttpGin(&endpoints.Endpoints{Endpoints: []endpoints.Endpoint{
		{
			HttpMethod:   http.MethodPost,
			RelativePath: relPath,
			Handler:      loginUser.LoginUserGinHandlerFunc(),
		},
	}})

	type args struct {
		aliasOrEmail string
		password     string
	}
	type want struct {
		response       *api.CommandResponse
		httpStatusCode int
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "good request with alias",
			args: args{
				aliasOrEmail: "martin_fowler",
				password:     "Linux648$",
			},
			want: want{
				response: &api.CommandResponse{
					Message: message.MessageData{
						ObjectId:   "martin_fowler",
						MessageId:  0,
						MessagePkg: "user",
						Variables:  message.Variables{"martin_fowler"},
						Text:       "user martin_fowler logged successful",
						Time:       time.Time{},
					},
				},
				httpStatusCode: 200,
			},
		},
		{
			name: "good request with email",
			args: args{
				aliasOrEmail: "martin_fowler@gmail.com",
				password:     "Linux648$",
			},
			want: want{
				response: &api.CommandResponse{
					Message: message.MessageData{
						ObjectId:   "martin_fowler@gmail.com",
						MessageId:  0,
						MessagePkg: "user",
						Variables:  message.Variables{"martin_fowler"},
						Text:       "user martin_fowler logged successful",
						Time:       time.Time{},
					},
				},
				httpStatusCode: 200,
			},
		},
		{
			name: "bad request - user not exists",
			args: args{
				aliasOrEmail: "user_not_found@gmail.com",
				password:     "Linux648$",
			},
			want: want{
				response: &api.CommandResponse{
					Message: message.MessageData{
						ObjectId:        "user_not_found@gmail.com",
						MessageId:       15,
						MessagePkg:      "user",
						Variables:       message.Variables{""},
						Text:            "email, alias or password are not correct. Repeat the access data.",
						Time:            time.Time{},
						ClientErrorType: message.ClientErrorBadRequest,
					},
				},
				httpStatusCode: 400,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := login.ClientArgs{
				AliasOrEmail: tt.args.aliasOrEmail,
				Password:     tt.args.password,
			}

			bytesCmd, err := json.Marshal(cmd)
			if err != nil {
				log.Panic(err)
			}

			response := mockGinSrv.DoRequest(integration.DoRequestCommand{
				BodyRequest:  bytesCmd,
				RelativePath: relPath,
			})

			require.EqualValues(t, tt.want.httpStatusCode, response.HttpCode)

			var gotRespBody *api.CommandResponse
			if err := json.Unmarshal(response.Body, &gotRespBody); err != nil {
				log.Panicln(err)
			}
			gotRespBody.Message.Time = time.Time{}

			require.EqualValues(t, tt.want.response, gotRespBody)
		})
	}
}

func tableUsersSetup() {
	const defaultPkf = "user"

	newUsersCommands := []*domain.NewUserCommand{
		{
			Uuid:       uuid.NewString(),
			Alias:      "martin_fowler",
			Name:       "martin",
			SecondName: "fowler",
			Email:      "martin_fowler@gmail.com",
			Password:   "Linux648$",
		},
	}
	for _, command := range newUsersCommands {
		cmdUuid, err := uuid.Parse(command.Uuid)
		if err != nil {
			log.Fatalln(err)
		}

		retLog := returnLog.NewReturnLog(cmdUuid, kernel.Instance.MessageRepository(), defaultPkf)
		user := domain.NewUser(*command, retLog)
		userRepositoryInstance.SaveUser(user, retLog)
		if retLog.Error() != nil {
			log.Fatalf("%v", retLog.Error())
		}
	}

}
