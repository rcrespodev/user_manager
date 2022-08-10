package handlers

import (
	"encoding/json"
	"github.com/rcrespodev/user_manager/api"
	"github.com/rcrespodev/user_manager/api/v1/endpoints"
	"github.com/rcrespodev/user_manager/api/v1/handlers/loginUser"
	"github.com/rcrespodev/user_manager/pkg/app/user/application/commands/login"
	"github.com/rcrespodev/user_manager/pkg/app/user/domain"
	"github.com/rcrespodev/user_manager/pkg/kernel"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain/message"
	"github.com/rcrespodev/user_manager/test/integration"
	"github.com/rcrespodev/user_manager/test/integration/v1/handlers/ginHandlers/utils"
	"github.com/stretchr/testify/require"
	"log"
	"net/http"
	"testing"
	"time"
)

const (
	relPath = endpoints.EndpointLogin
)

func TestLoginUserGinHandlerFunc(t *testing.T) {
	userRepository := kernel.Instance.UserRepository()
	require.NoError(t, loginUserSetup(userRepository))

	mockGinSrv := integration.NewTestServerHttpGin(endpoints.Endpoints{
		relPath: endpoints.Endpoint{
			HttpMethod: http.MethodPost,
			Handler:    loginUser.LoginUserGinHandlerFunc(),
		},
	})

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

			// Header check
			require.EqualValues(t, tt.want.httpStatusCode, response.HttpCode)

			// Body check
			var gotRespBody *api.CommandResponse
			if err := json.Unmarshal(response.Body, &gotRespBody); err != nil {
				log.Panicln(err)
			}
			gotRespBody.Message.Time = time.Time{}

			require.EqualValues(t, tt.want.response, gotRespBody)

			// Jwt Check
			if response.HttpCode == 200 {
				utils.TokenValidationForTesting(t, response.Header)
			}
		})
	}
}

func loginUserSetup(repository domain.UserRepository) error {
	newUsersCommands := []*domain.NewUserCommand{
		{
			Uuid:       "123e4567-e89b-12d3-a456-426614174000",
			Alias:      "martin_fowler",
			Name:       "martin",
			SecondName: "fowler",
			Email:      "martin_fowler@gmail.com",
			Password:   "Linux648$",
		},
	}
	return utils.TableUsersSetup(newUsersCommands, repository)
}
