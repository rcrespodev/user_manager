package registerUser

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/rcrespodev/user_manager/api"
	"github.com/rcrespodev/user_manager/api/v1/handlers/registerUser"
	"github.com/rcrespodev/user_manager/api/v1/routes"
	"github.com/rcrespodev/user_manager/pkg/app/user/application/register"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain/message"
	"github.com/rcrespodev/user_manager/test/integration"
	"log"
	"net/http"
	"reflect"
	"testing"
	"time"
)

const (
	registerUserRelPath = "/register_user"
)

func TestRegisterUserGinHandlerFunc(t *testing.T) {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal(err)
	}
	//kernel.NewPrdKernel()
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

			//response := integration.NewHttpRequest(integration.NewHttpRequestCommand{
			//	Method: http.MethodPost,
			//	//Host:   "app",
			//	Host:        "0.0.0.0",
			//	Port:        "8080",
			//	Path:        registerUserRelPath,
			//	Body:        bytes.NewReader(bytesCmd),
			//	ContentType: "application/json",
			//})

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
				return
			}
			var gotRespBody *api.CommandResponse
			if err := json.Unmarshal(response.Body, &gotRespBody); err != nil {
				log.Fatalln(err)
			}
			gotRespBody.Message.Time = time.Time{}
			if !reflect.DeepEqual(gotRespBody, tt.want.response) {
				t.Errorf("Body()\n\t- got: %v\n\t- want: %v", gotRespBody, tt.want.response)
				return
			}
		})
	}
}
