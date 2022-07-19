package registerUser

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/rcrespodev/user_manager/api"
	"github.com/rcrespodev/user_manager/api/v1/handlers/registerUser"
	"github.com/rcrespodev/user_manager/api/v1/routes"
	"github.com/rcrespodev/user_manager/pkg/app/user/application/register"
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
	registerUserRelPath = "/register_user"
)

var userRepositoryInstance domain.UserRepository

func TestRegisterUserGinHandlerFunc(t *testing.T) {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal(err)
	}

	userRepositoryInstance = kernel.Instance.UserRepository()

	tableUsersSetup()

	mockGinSrv := integration.NewTestServerHttpGin(&routes.Routes{Routes: []routes.Route{
		{
			HttpMethod:   http.MethodPost,
			RelativePath: registerUserRelPath,
			Handler:      registerUser.RegisterUserGinHandlerFunc(),
		},
	}})

	type args struct {
		uuid       string
		alias      string
		name       string
		secondName string
		email      string
		password   string
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
			name: "good request",
			args: args{
				uuid:       "123e4567-e89b-12d3-a456-426614174000",
				alias:      "martin_fowler",
				name:       "martin",
				secondName: "fowler",
				email:      "foo@test.com",
				password:   "Linux648$",
			},
			want: want{
				response: &api.CommandResponse{
					Message: message.MessageData{
						ObjectId:        "martin_fowler",
						MessageId:       1,
						MessagePkg:      "user",
						Variables:       message.Variables{"martin_fowler"},
						Text:            "user martin_fowler created successful",
						Time:            time.Time{},
						ClientErrorType: 0,
					},
				},
				httpStatusCode: 200,
			},
		},
		{
			name: "user uuid already exists",
			args: args{
				uuid:       "123e4567-e89b-12d3-a456-426614174000",
				alias:      "linus_torvalds",
				name:       "linus",
				secondName: "torvalds",
				email:      "linus@test.com.ar",
				password:   "Linux648$",
			},
			want: want{
				response: &api.CommandResponse{
					Message: message.MessageData{
						ObjectId:        "linus_torvalds",
						MessageId:       14,
						MessagePkg:      "user",
						Variables:       message.Variables{"uuid", "linus_torvalds"},
						Text:            "user with component: uuid and value: 123e4567-e89b-12d3-a456-426614174000 already exists",
						Time:            time.Time{},
						ClientErrorType: message.ClientErrorBadRequest,
					},
				},
				httpStatusCode: 400,
			},
		},
		{
			name: "user alias already exists",
			args: args{
				alias:      "user_exists_alias",
				name:       "martin",
				secondName: "fowler",
				email:      "foo@test.com.ar",
				password:   "Linux648$",
			},
			want: want{
				response: &api.CommandResponse{
					Message: message.MessageData{
						ObjectId:        "user_exists_alias",
						MessageId:       14,
						MessagePkg:      "user",
						Variables:       message.Variables{"alias", "user_exists_alias"},
						Text:            "user with component: alias and value: user_exists_alias already exists",
						Time:            time.Time{},
						ClientErrorType: message.ClientErrorBadRequest,
					},
				},
				httpStatusCode: 400,
			},
		},
		{
			name: "user email already exists",
			args: args{
				alias:      "martin_fowler_2",
				name:       "martin",
				secondName: "fowler",
				email:      "email_exists@test.com.ar",
				password:   "Linux648$",
			},
			want: want{
				response: &api.CommandResponse{
					Message: message.MessageData{
						ObjectId:        "martin_fowler_2",
						MessageId:       14,
						MessagePkg:      "user",
						Variables:       message.Variables{"email", "email_exists@test.com.ar"},
						Text:            "user with component: email and value: email_exists@test.com.ar already exists",
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
			cmdUuid := uuid.New()
			cmd := register.ClientArgs{
				Uuid:       cmdUuid.String(),
				Alias:      tt.args.alias,
				Name:       tt.args.name,
				SecondName: tt.args.secondName,
				Email:      tt.args.email,
				Password:   tt.args.password,
			}

			bytesCmd, err := json.Marshal(cmd)
			if err != nil {
				log.Panic(err)
			}

			response := mockGinSrv.DoRequest(integration.DoRequestCommand{
				BodyRequest:  bytesCmd,
				RelativePath: registerUserRelPath,
			})

			var gotRespBody *api.CommandResponse
			if err := json.Unmarshal(response.Body, &gotRespBody); err != nil {
				log.Panicln(err)
			}
			gotRespBody.Message.Time = time.Time{}

			require.EqualValues(t, tt.want.httpStatusCode, response.HttpCode)
			require.EqualValues(t, tt.want.response, gotRespBody)

			if response.HttpCode == 200 {
				retLog := returnLog.NewReturnLog(cmdUuid, kernel.Instance.MessageRepository(), "user")

				actualUser := userRepositoryInstance.FindUser(domain.FindUserCommand{
					Password: tt.args.password,
					Log:      retLog,
					Where: []domain.WhereArgs{
						{
							Field: "uuid",
							Value: cmdUuid.String(),
						},
					},
				})

				expectedUser := domain.NewUser(domain.NewUserCommand{
					Uuid:       cmdUuid.String(),
					Alias:      tt.args.alias,
					Name:       tt.args.name,
					SecondName: tt.args.secondName,
					Email:      tt.args.email,
					Password:   tt.args.password,
				}, retLog)
				if retLog.Error() != nil {
					panic(retLog.Error())
				}

				require.EqualValues(t, expectedUser.Uuid(), actualUser.Uuid())
				require.EqualValues(t, expectedUser.Alias(), actualUser.Alias())
				require.EqualValues(t, expectedUser.Name(), actualUser.Name())
				require.EqualValues(t, expectedUser.SecondName(), actualUser.SecondName())
				require.EqualValues(t, expectedUser.Email(), actualUser.Email())
				require.EqualValues(t, expectedUser.Password().String(), actualUser.Password().String())
			}
		})
	}
}

func tableUsersSetup() {
	const defaultPkf = "user"

	newUsersCommands := []*domain.NewUserCommand{
		{
			Uuid:       uuid.NewString(),
			Alias:      "user_exists_alias",
			Name:       "martin",
			SecondName: "fowler",
			Email:      "foo@test.com.ar",
			Password:   "Linux648$",
		},
		{
			Uuid:       uuid.NewString(),
			Alias:      "user_exists_email",
			Name:       "martin",
			SecondName: "fowler",
			Email:      "email_exists@test.com.ar",
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
