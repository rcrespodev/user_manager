package registerUser

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/rcrespodev/user_manager/api"
	"github.com/rcrespodev/user_manager/api/v1/endpoints"
	"github.com/rcrespodev/user_manager/api/v1/handlers/registerUser"
	domain2 "github.com/rcrespodev/user_manager/pkg/app/auth/domain"
	"github.com/rcrespodev/user_manager/pkg/app/user/application/commands/register"
	"github.com/rcrespodev/user_manager/pkg/app/user/domain"
	"github.com/rcrespodev/user_manager/pkg/kernel"
	returnLog "github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain/message"
	"github.com/rcrespodev/user_manager/test/integration"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
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
	userRepositoryInstance = kernel.Instance.UserRepository()

	tableUsersSetup()

	mockGinSrv := integration.NewTestServerHttpGin(endpoints.Endpoints{
		registerUserRelPath: endpoints.Endpoint{
			HttpMethod: http.MethodPost,
			Handler:    registerUser.RegisterUserGinHandlerFunc(),
		},
	})

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
				email:      "user_registred@test.com",
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
			name: "user alias already exists",
			args: args{
				alias:      "user_exists_alias",
				name:       "martin",
				secondName: "fowler",
				email:      "test@test.com.ar",
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

			switch response.HttpCode {
			case 200:
				// Token validation
				require.NotNil(t, response.Header.Get("Token"))
				err := domain2.ParseJwt(response.Header.Get("Token"), kernel.Instance.JwtConfig())
				require.NoError(t, err)

				// Database validation
				retLog := returnLog.NewReturnLog(cmdUuid, kernel.Instance.MessageRepository(), "user")

				actualUser := userRepositoryInstance.FindUser(domain.FindUserQuery{
					Log: retLog,
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

				require.EqualValues(t, expectedUser.Uuid().String(), actualUser.Uuid)
				require.EqualValues(t, expectedUser.Alias().Alias(), actualUser.Alias)
				require.EqualValues(t, expectedUser.Name().Name(), actualUser.Name)
				require.EqualValues(t, expectedUser.SecondName().Name(), actualUser.SecondName)
				require.EqualValues(t, expectedUser.Email().Address(), actualUser.Email)
				// password are stored in hash format in DB.
				err = bcrypt.CompareHashAndPassword(actualUser.HashedPassword, []byte(expectedUser.Password().String()))
				require.NoError(t, err)
			default:
				require.EqualValues(t, "", response.Header.Get("Token"))
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
