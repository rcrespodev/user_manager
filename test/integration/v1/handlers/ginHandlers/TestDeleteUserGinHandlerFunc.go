package handlers

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/rcrespodev/user_manager/api"
	"github.com/rcrespodev/user_manager/api/v1/endpoints"
	"github.com/rcrespodev/user_manager/api/v1/handlers/deleteUser"
	delete "github.com/rcrespodev/user_manager/pkg/app/user/application/commands/deleteUser"
	"github.com/rcrespodev/user_manager/pkg/app/user/domain"
	"github.com/rcrespodev/user_manager/pkg/kernel"
	returnLog "github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain/message"
	"github.com/rcrespodev/user_manager/test/integration"
	"github.com/rcrespodev/user_manager/test/integration/v1/handlers/ginHandlers/utils"
	"github.com/stretchr/testify/require"
	"log"
	"net/http"
	"testing"
	"time"
)

func TestDeleteUserGinHandlerFunc(t *testing.T) {
	userRepositoryInstance = kernel.Instance.UserRepository()
	require.NoError(t, deleteUserSetup(userRepositoryInstance))

	mockGinSrv := integration.NewTestServerHttpGin(endpoints.Endpoints{
		endpoints.BuildEndpointKey(endpoints.EndpointUser, http.MethodDelete): {
			HttpMethod:     http.MethodDelete,
			RelPath:        endpoints.EndpointUser,
			Handler:        deleteUser.DeleteUserGinHandlerFunc(),
			AuthValidation: true,
		},
	})

	type args struct {
		userUuid string
		token    string
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
			name: "correct request",
			args: args{
				userUuid: "123e4567-e89b-12d3-a456-426614174005",
				token:    "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE2NjAyNDU1NzcsImtleSI6IjEyM2U0NTY3LWU4OWItMTJkMy1hNDU2LTQyNjYxNDE3NDAwNSJ9.OSENP2-ug-dUd3GLSvbi9jZ3AYlpTAWp051y7WlYwqC1L__d229OJHYTxmevCJRhCGsEGNQdP98QvOsvIIuX8uchldO4XNzk7ZCH3CQ33Y99GN2aPXtUY7lKM_RMelOvSnukaHQ14FJ9bry3z_SHznLg6YiURq0dX3C_VfUvD8Jy6ARcOhbiFz96nvvnyuwdA1A2ok0FEr7hfDHduEEquW_ZMeEspgQCjJJ4NO_dTmBn3COk6N0B9vn74SDBJ57RpSTuQbyCrPyHDMOjuUHitQShELeQc0WjOgw8eqsJB4fwF6glT5N66Nph6aIrz1FPEwfr_TVYpwTPm94fpmBEzCWQG4OFBHyy3LARmZVifpGaKcdD_gMblpsFlw-LtzfyYZsKmolKe9bkjMpRf1vwIonImvZsqn_-bIHcP4m5Gi5y1SuXai25IJVkRvBgORE-HeYpBlKcqFzc0TNJBXOCDt7mZxIuRE5izYoF8nX-rJCOF-uCyTvdjJMCOcpDFfyF4Y0qZXkpw61Pea66VzwwbgljNKhjgNa9Uk4tZ8gbSj5lb1-IUbbBfxJplmujliyjqIiq2T5AAO2-WU-QzQq07xdoZvaZLYkxfFBrlWC_8L3mA33pCshMwGHSGBMZwL5eB8a4n9CYEZhXHVKDmSCq0xUgN-B4Gdbgc68IXTXcOa0",
			},
			want: want{
				response: &api.CommandResponse{
					Message: message.MessageData{
						ObjectId:   "123e4567-e89b-12d3-a456-426614174005",
						MessageId:  3,
						MessagePkg: "user",
						Variables:  message.Variables{"user_exists_alias"},
						Text:       "user user_exists_alias deleted successful",
					},
				},
				httpStatusCode: 200,
			},
		},
		{
			name: "user exists but token uuid not match with user uuid",
			args: args{
				userUuid: "123e4567-e89b-12d3-a456-426614174006",
				token:    "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE2NjAyNDU1NzcsImtleSI6IjEyM2U0NTY3LWU4OWItMTJkMy1hNDU2LTQyNjYxNDE3NDAwNSJ9.OSENP2-ug-dUd3GLSvbi9jZ3AYlpTAWp051y7WlYwqC1L__d229OJHYTxmevCJRhCGsEGNQdP98QvOsvIIuX8uchldO4XNzk7ZCH3CQ33Y99GN2aPXtUY7lKM_RMelOvSnukaHQ14FJ9bry3z_SHznLg6YiURq0dX3C_VfUvD8Jy6ARcOhbiFz96nvvnyuwdA1A2ok0FEr7hfDHduEEquW_ZMeEspgQCjJJ4NO_dTmBn3COk6N0B9vn74SDBJ57RpSTuQbyCrPyHDMOjuUHitQShELeQc0WjOgw8eqsJB4fwF6glT5N66Nph6aIrz1FPEwfr_TVYpwTPm94fpmBEzCWQG4OFBHyy3LARmZVifpGaKcdD_gMblpsFlw-LtzfyYZsKmolKe9bkjMpRf1vwIonImvZsqn_-bIHcP4m5Gi5y1SuXai25IJVkRvBgORE-HeYpBlKcqFzc0TNJBXOCDt7mZxIuRE5izYoF8nX-rJCOF-uCyTvdjJMCOcpDFfyF4Y0qZXkpw61Pea66VzwwbgljNKhjgNa9Uk4tZ8gbSj5lb1-IUbbBfxJplmujliyjqIiq2T5AAO2-WU-QzQq07xdoZvaZLYkxfFBrlWC_8L3mA33pCshMwGHSGBMZwL5eB8a4n9CYEZhXHVKDmSCq0xUgN-B4Gdbgc68IXTXcOa0",
			},
			want: want{
				response: &api.CommandResponse{
					Message: message.MessageData{
						ObjectId:        "123e4567-e89b-12d3-a456-426614174006",
						MessageId:       0,
						MessagePkg:      "Authorization",
						Text:            "Unauthorized",
						ClientErrorType: message.ClientErrorUnauthorized,
					},
				},
				httpStatusCode: 401,
			},
		},
		{
			name: "token uuid and token user match but user not exists in db",
			args: args{
				userUuid: "123e4567-e89b-12d3-a456-426614174007",
				token:    "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE2NjAyNTMwODIsImtleSI6IjEyM2U0NTY3LWU4OWItMTJkMy1hNDU2LTQyNjYxNDE3NDAwNyJ9.K62bdLDOerhcFT6bLG3WguDic5XpMEzz8UdpsQojfVbRQyxV754NW9FsohN4nqUK4VKDnjeGYAqxppHT_bUMIm6N5iMzAYGxvrEHv9D2tgjixVNwj0zDoJxq7wDeNrWelJFzorJ3hb5S6dkSbJmTRivzHkCgN7epY_4zf_KQF4pR35wwwAca6vU2UNteSjG4MSeHk8uPrA0lv0XaqokVAJOg0yNlU2NEFTHzsF1P6EMDpLxfFx1xef3xv4WTWHCNQg98S7k91BamfIVaYfRQTqqtbLnPPEth46dNro5KKqbxH2QvmhWVIHPJZ-_7mFQgXumw8cpyqLtFBNlulwQslmxOrpfPjoYkocDeter3Oris7BuX6x1B6Pwu40O_5My9lv12kBlW05FdACmqMjhmMhitXrAJpoPdHWxpHBjHcNy6D8Z07Fcoe9a8spAWTdTXTP5afND_1dk_oSsy8nYJzZy9aRuSvE0BYAwM-7KsjSW8jgEHMfSCrZn-_zifDiyVqq5Tk_HoTmVUmrIf4VpfGlWFvTuTYkbGg9x7xMJI3y0VxzVVEtqVXDQPLQl83rkb_atZ64q-_5mV9RHicllcDtyOtlv1cDgGJ_xiYvhmA1FV9mbNMWMBwrmPJKncsuSLGD9Xf5Ot0puUBk1FcU5xT5DtndDNc0jp6sLV9Tt9Zms",
			},
			want: want{
				response: &api.CommandResponse{
					Message: message.MessageData{
						ObjectId:        "123e4567-e89b-12d3-a456-426614174007",
						MessageId:       17,
						MessagePkg:      "user",
						Text:            "none of the input values correspond to a registered user",
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
			cmd := delete.ClientArgs{UserUuid: tt.args.userUuid}

			bytesCmd, err := json.Marshal(cmd)
			if err != nil {
				log.Panic(err)
			}

			response := mockGinSrv.DoRequest(integration.DoRequestCommand{
				BodyRequest:  bytesCmd,
				Method:       http.MethodDelete,
				RelativePath: endpoints.EndpointUser,
				Token:        tt.args.token,
			})

			// Header check
			require.EqualValues(t, tt.want.httpStatusCode, response.HttpCode)

			// Body check
			var gotRespBody *api.CommandResponse
			if err := json.Unmarshal(response.Body, &gotRespBody); err != nil {
				log.Fatal(err)
			}
			gotRespBody.Message.Time = time.Time{}
			require.EqualValues(t, tt.want.response, gotRespBody)

			switch response.HttpCode {
			case 200:
				// Token validation
				utils.TokenValidationForTesting(t, response.Header)

				// Database validation
				retLog := returnLog.NewReturnLog(cmdUuid, kernel.Instance.MessageRepository(), "user")

				deletedUser := userRepositoryInstance.FindUser(domain.FindUserQuery{
					Log: retLog,
					Where: []domain.WhereArgs{
						{
							Field: "uuid",
							Value: tt.args.userUuid,
						},
					},
				})

				require.Nil(t, deletedUser)
			}
		})
	}
}

func deleteUserSetup(repository domain.UserRepository) error {
	newUsersCommands := []*domain.NewUserCommand{
		{
			Uuid:       "123e4567-e89b-12d3-a456-426614174005",
			Alias:      "user_exists_alias",
			Name:       "martin",
			SecondName: "fowler",
			Email:      "foo@test.com.ar",
			Password:   "Linux648$",
		},
		{
			Uuid:       "123e4567-e89b-12d3-a456-426614174006",
			Alias:      "user_exists_alias_2",
			Name:       "martin",
			SecondName: "fowler",
			Email:      "foo_2@test.com.ar",
			Password:   "Linux648$",
		},
	}
	return utils.TableUsersSetup(newUsersCommands, repository)
}
