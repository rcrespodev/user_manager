package handlers

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/rcrespodev/user_manager/api"
	"github.com/rcrespodev/user_manager/api/v1/endpoints"
	"github.com/rcrespodev/user_manager/api/v1/handlers/logOutUser"
	jwtDomain "github.com/rcrespodev/user_manager/pkg/app/authJwt/domain"
	"github.com/rcrespodev/user_manager/pkg/kernel"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain/message"
	"github.com/rcrespodev/user_manager/test/integration"
	"github.com/stretchr/testify/require"
	"log"
	"net/http"
	"testing"
	"time"
)

var jwtRepository jwtDomain.JwtRepository

func TestLogOutUserGinHandlerFunc(t *testing.T) {
	jwtRepository = kernel.Instance.JwtRepository()

	jwtRepositorySetup()

	type args struct {
		token string
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
			name: "token and user exists",
			args: args{
				token: "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE2NTkzODQxOTAsImtleSI6IjEyM2U0NTY3LWU4OWItMTJkMy1hNDU2LTQyNjYxNDE3NDAwMCJ9.sqCpS8ISEnmvVCMQMNkWvmb7OAhM8HIuXE6L4rl2X2uvn5wXbfRD79FmjaST3nFekSeNtr9SwicVxU9VRnylKK5xexDHaoPoBZ4h7W6KYsYcJGfF-SN2dKAUGkwHUExFbkwGEgny9dFZKXMArBSLkzG8wWx5QSn38j3a4D1fvBKvLazTAEhP3SAfm2idlgnSGV3BvwATsIDlO_WCpToBgTre9pumKJXHGm43Tay2Zk6vYLyR1-2J2he9jkZec6pS4Dxe52PuVQnGiC6RFv6ttqEP42H9t_dTzfxcz_cwbj0fd5zURK_5F5zTI0aZfd7siACwleB7qAF9jgw6-vyh0kad25S0q5n2NGOK5YULfCbMGg3zDHOKhKEKuWu5aqKEVSUtTowZCl-LGsuU1a3ngKEBvEpCyS7lDRRA8zOCMIWTYbaJCYGsSWdb7d4f0N_jzHTj_8HCPwl8AxBtpCUYFFKEKthVRUz5sh3YGfTOVJVvXoFYUdU-jN4CpDrQiPbSEcFa8NzVdylch5WT06f_KMzJ-_yKW_DPG_6zYKJmZM5l7VtnPzgXRyVXQfOlsE0WWdigZv_2nE7WyhlRwHtrZknqLT9Dr5LA3nn1nN4EEAk90R2pKFfXOlIwbNaO4bufbEgRN6TkQB0eiu80-4didZtJEJMJVRpz1Qfn7fwpcsM",
			},
			want: want{
				response: &api.CommandResponse{Message: message.MessageData{
					ObjectId:        "123e4567-e89b-12d3-a456-426614174000",
					MessageId:       16,
					MessagePkg:      "user",
					Text:            "user logged out successful",
					ClientErrorType: 0,
				}},
				httpStatusCode: 200,
			},
		},
		{
			name: "token not exists",
			args: args{
				token: "token.not.exists",
			},
			want: want{
				response:       nil,
				httpStatusCode: 401,
			},
		},
		{
			name: "token is valid but not exists in db (user is already logged out)",
			args: args{
				token: "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE2NTkzODcyNjIsImtleSI6IjEyM2U0NTY3LWU4OWItMTJkMy1hNDU2LTQyNjYxNDE3NDAwMSJ9.KGzLvI6GZg261KDLmWLpNf75d-7i_i1upvuvW6-_knqzRhvItcNR9nAyFO4AXM9Q2CQWnreaEqhrX5R49Q6vzIS19ykt5qpf7FgBrRVQQ5M7f05JkrqzvHBMe0yjGKPA3sdED1stLiwScBn8xhK7dEU7sj6WrYL7HxdNHyymZ5n7Neywpzp_bS_ssWYtAcxQq7waBR7V_eb633VgcOWVl0VpKQ4YfhqKN8i5YSUSASQUAF86MQ-f28dCuKYM4EVLenguDjl8-50mwBzwYKvFMRID0ys-EFlcm2jXvF3J6_QgoFeElGUIL7Hw4gCDwXrcmhRbJv1jN13O1zmQdL8cL5YVnoJBKLLggET2J8WmC5p3C0ILDQ2P60FcTtqLYV0uGOiGHthXvvq6MzXvcBebEGceV3XZvzouWTeowWvIrjZk-_ObfLxu5MctN7NFm1pKybjnwHVH_ToF3sf0pQoGEDQ-lRY0iOKIAlFBhYTpkFfs4VaNiRnR0GnVPog09JU9jFJN6YfCY5jjd21DOamS3hGHIts1vvJtRdGnUxCwFsnoaSgqnv76MMo9O9lm_G1CaBohdmz9uF-NsBf3N36zYNDuyqKxnaXI2_4fuLDnJ4qQ-6-y8gcFn9UQN3kX2nsekeqxgPHooaKgLNuvIUQRew-7oGp5acl9_eTiMMD38XY",
			},
			want: want{
				response: &api.CommandResponse{Message: message.MessageData{
					ObjectId:        "123e4567-e89b-12d3-a456-426614174001",
					MessageId:       0,
					MessagePkg:      "Authorization",
					Text:            "Unauthorized",
					ClientErrorType: 2,
				}},
				httpStatusCode: 401,
			},
		},
		{
			name: "token exists and is invalid in db",
			args: args{
				token: "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE2NTk2Njk2NTAsImtleSI6IjEyM2U0NTY3LWU4OWItMTJkMy1hNDU2LTQyNjYxNDE3NDAwMiJ9.DL6hU_aH5nuMFMol7WUuZMF9PSbW-BGzZJaVS7UuNQvBhZ_KQ9yMqhHFYr92WylVnAaCTVj7nT9pjKQHGvjfoWMb6mCFyEw1xaztxku1wTN4wChJbyUBENop6JifWMdkY2rOHUoB2PSPA1uaGY0LVVm76AeX2oguOG1PXx_3g_Xj4BJ2xzfshLXnf-f77WnKkLRS28eiv5rBjktvlyTvbp9swZ8jywO5k2zdSWrU7cGXNU7fo-yZGBXmZowTE6ftnFCpu55lV-p11X81E8thb68-3QqnChrF5It36d9k-HUBUsf7OAm8acDkuVqAtlJ2YufPmfEH9COyvcRAz_Z2PNKnpUbaElFNTio169hVJXfXAeirl_TRj3eQ6YR8kORWez_zuav6viuds3OEdN9ENTYu4k1F4i8IhlNyQCsKCo1-BB9QWsagU0Mfq37x_q66RoTlrvHnKkr6duXEcC8G-3jSEZ8HL2TjdPCMQAFi66joBUd_j0wVPLqIAnIyqIyw-LbIWi4t9VE0lmZ9InggXo9ogYbEGbv8RzpKNN83oujx42zHYNKmxGguvM94Ps5ByrWxQmuflDdexUxxqw6NQxM_SzTHpm11DN6nhaI0u2gGye1Ydp7t85gLSFPrMlQj-YCwndgUx-de5ZUmYGOAI-M7-2Kt2OpzBRwsnDPj4us",
			},
			want: want{
				response: &api.CommandResponse{Message: message.MessageData{
					ObjectId:        "123e4567-e89b-12d3-a456-426614174002",
					MessageId:       1,
					MessagePkg:      "Authorization",
					Text:            "user is not logged",
					ClientErrorType: 2,
				}},
				httpStatusCode: 401,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockServer := integration.NewTestServerHttpGin(endpoints.Endpoints{
				endpoints.BuildEndpointKey(endpoints.EndpointUserLogin, http.MethodDelete): {
					HttpMethod:     http.MethodPost,
					RelPath:        endpoints.EndpointUserLogin,
					Handler:        logOutUser.LogOutUserGinHandlerFunc(),
					AuthValidation: true,
				},
			})

			response := mockServer.DoRequest(integration.DoRequestCommand{
				RelativePath: endpoints.EndpointUserLogin,
				Method:       http.MethodDelete,
				Token:        tt.args.token,
			})

			// Header check
			require.EqualValues(t, tt.want.httpStatusCode, response.HttpCode)

			// Body check
			var gotRespBody *api.CommandResponse
			if len(response.Body) > 0 {
				if err := json.Unmarshal(response.Body, &gotRespBody); err != nil {
					log.Fatal(err)
				}
				gotRespBody.Message.Time = time.Time{}
			}

			require.EqualValues(t, tt.want.response, gotRespBody)

			if tt.want.httpStatusCode == http.StatusOK {
				jwt := jwtRepository.FindByUuid(jwtDomain.FindByUuidQuery{Uuid: tt.want.response.Message.ObjectId})
				require.False(t, jwt.IsValid)
			}

		})
	}
}

func jwtRepositorySetup() {
	retLog := domain.NewReturnLog(uuid.New(), kernel.Instance.MessageRepository(), "authorization")
	commands := []jwtDomain.UpdateCommand{
		{
			Command: &jwtDomain.JwtSchema{
				Uuid:     "123e4567-e89b-12d3-a456-426614174000",
				IsValid:  true,
				Token:    "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE2NTkzODQxOTAsImtleSI6IjEyM2U0NTY3LWU4OWItMTJkMy1hNDU2LTQyNjYxNDE3NDAwMCJ9.sqCpS8ISEnmvVCMQMNkWvmb7OAhM8HIuXE6L4rl2X2uvn5wXbfRD79FmjaST3nFekSeNtr9SwicVxU9VRnylKK5xexDHaoPoBZ4h7W6KYsYcJGfF-SN2dKAUGkwHUExFbkwGEgny9dFZKXMArBSLkzG8wWx5QSn38j3a4D1fvBKvLazTAEhP3SAfm2idlgnSGV3BvwATsIDlO_WCpToBgTre9pumKJXHGm43Tay2Zk6vYLyR1-2J2he9jkZec6pS4Dxe52PuVQnGiC6RFv6ttqEP42H9t_dTzfxcz_cwbj0fd5zURK_5F5zTI0aZfd7siACwleB7qAF9jgw6-vyh0kad25S0q5n2NGOK5YULfCbMGg3zDHOKhKEKuWu5aqKEVSUtTowZCl-LGsuU1a3ngKEBvEpCyS7lDRRA8zOCMIWTYbaJCYGsSWdb7d4f0N_jzHTj_8HCPwl8AxBtpCUYFFKEKthVRUz5sh3YGfTOVJVvXoFYUdU-jN4CpDrQiPbSEcFa8NzVdylch5WT06f_KMzJ-_yKW_DPG_6zYKJmZM5l7VtnPzgXRyVXQfOlsE0WWdigZv_2nE7WyhlRwHtrZknqLT9Dr5LA3nn1nN4EEAk90R2pKFfXOlIwbNaO4bufbEgRN6TkQB0eiu80-4didZtJEJMJVRpz1Qfn7fwpcsM",
				Duration: 5,
			},
		},
		{
			Command: &jwtDomain.JwtSchema{
				Uuid:     "123e4567-e89b-12d3-a456-426614174002",
				IsValid:  false,
				Token:    "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE2NTk2Njk2NTAsImtleSI6IjEyM2U0NTY3LWU4OWItMTJkMy1hNDU2LTQyNjYxNDE3NDAwMiJ9.DL6hU_aH5nuMFMol7WUuZMF9PSbW-BGzZJaVS7UuNQvBhZ_KQ9yMqhHFYr92WylVnAaCTVj7nT9pjKQHGvjfoWMb6mCFyEw1xaztxku1wTN4wChJbyUBENop6JifWMdkY2rOHUoB2PSPA1uaGY0LVVm76AeX2oguOG1PXx_3g_Xj4BJ2xzfshLXnf-f77WnKkLRS28eiv5rBjktvlyTvbp9swZ8jywO5k2zdSWrU7cGXNU7fo-yZGBXmZowTE6ftnFCpu55lV-p11X81E8thb68-3QqnChrF5It36d9k-HUBUsf7OAm8acDkuVqAtlJ2YufPmfEH9COyvcRAz_Z2PNKnpUbaElFNTio169hVJXfXAeirl_TRj3eQ6YR8kORWez_zuav6viuds3OEdN9ENTYu4k1F4i8IhlNyQCsKCo1-BB9QWsagU0Mfq37x_q66RoTlrvHnKkr6duXEcC8G-3jSEZ8HL2TjdPCMQAFi66joBUd_j0wVPLqIAnIyqIyw-LbIWi4t9VE0lmZ9InggXo9ogYbEGbv8RzpKNN83oujx42zHYNKmxGguvM94Ps5ByrWxQmuflDdexUxxqw6NQxM_SzTHpm11DN6nhaI0u2gGye1Ydp7t85gLSFPrMlQj-YCwndgUx-de5ZUmYGOAI-M7-2Kt2OpzBRwsnDPj4us",
				Duration: 5,
			},
		},
	}
	for _, command := range commands {
		jwtRepository.Update(command, retLog)
		if retLog.Error() != nil {
			log.Fatal(retLog.Error())
		}
	}

}
