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
	"reflect"
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

	type UserSchema struct {
		Uuid       string
		Alias      string
		Name       string
		SecondName string
		Email      string
		Password   string
	}
	type args struct {
		alias      string
		name       string
		secondName string
		email      string
		password   string
	}
	type want struct {
		response       *api.CommandResponse
		httpStatusCode int
		//user           *UserSchema
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "good request",
			args: args{
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
			name: "user alias already exists",
			args: args{
				alias:      "user_exists",
				name:       "martin",
				secondName: "fowler",
				email:      "foo@test.com.ar",
				password:   "Linux648$",
			},
			want: want{
				response: &api.CommandResponse{
					Message: message.MessageData{
						ObjectId:        "user_exists",
						MessageId:       14,
						MessagePkg:      "user",
						Variables:       message.Variables{"alias", "user_exists"},
						Text:            "user with component: alias and value: user_exists already exists",
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
				log.Fatalln(err)
			}

			mockGinSrv := integration.NewTestServerHttpGin(&routes.Routes{Routes: []routes.Route{
				{
					HttpMethod:   http.MethodPost,
					RelativePath: registerUserRelPath,
					Handler:      registerUser.RegisterUserGinHandlerFunc(),
				},
			}})

			response := mockGinSrv.DoRequest(integration.DoRequestCommand{
				BodyRequest:  bytesCmd,
				RelativePath: registerUserRelPath,
			})

			if gotHttpCode := response.HttpCode; !reflect.DeepEqual(gotHttpCode, tt.want.httpStatusCode) {
				t.Errorf("HttpCode()\n\t- got: %v\n\t- want: %v", gotHttpCode, tt.want.httpStatusCode)
			}

			var gotRespBody *api.CommandResponse
			if err := json.Unmarshal(response.Body, &gotRespBody); err != nil {
				log.Fatalln(err)
			}

			gotRespBody.Message.Time = time.Time{}
			if !reflect.DeepEqual(gotRespBody, tt.want.response) {
				t.Errorf("Body()\n\t- got: %v\n\t- want: %v", gotRespBody, tt.want.response)
			}

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

				//userRepositoryInstance.FindUserById(domain.FindByIdCommand{
				//	Uuid: cmdUuid,
				//	FindUserCommand: domain.FindUserCommand{
				//		Password: tt.args.password,
				//		Log:      retLog,
				//		//Wg:       wg,
				//	},
				//}, actualUser)
				//wg.Wait()
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
			Alias:      "user_exists",
			Name:       "martin",
			SecondName: "fowler",
			Email:      "foo@test.com.ar",
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
